package loader

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/seantcanavan/anon-eth-net/logger"
)

// Loader represents a struct that will load a set of processes and watch over
// them. It knows the name of every process that it should be keeping an eye on
// as well as how to resurrect that process should it no longer be executing.
// The idea of the Loader is to make sure that all external process dependencies
// are executing and are in a healthy state as much as possible.
type Loader struct {
	Processes []LoaderProcess // the slice of LoaderProcesses which the loader will execute and keep an eye on
}

type LoaderProcess struct {
	Name      string
	Command   string
	Arguments []string
	Start     int64
	End       int64
	Duration  int64
	Lgr       *logger.Logger
}

// NewLoader will initialize a new instance of the Loader struct and execute the
// associated processes from the given file with the appropriate parameters.
// Each individual process will have its own logs.
func NewLoader(processesPath string) (*Loader, error) {

	loadedProcesses, loadErr := processesFromJSONFile(processesPath)
	if loadErr != nil {
		return nil, loadErr
	}

	logger.Lgr.LogMessage("Successfully loaded processes from file: %v", processesPath)
	logger.Lgr.LogMessage("Successfully instantiated loader from JSON:\n%+v", loadedProcesses)

	loader := &Loader{Processes: loadedProcesses}

	return loader, nil
}

// processesFromJSONFile will read in a set of JSON values which define both
// the canonical name of the process as well as the command and any associated
// parameters to successfully execute that command. A slice containing a
// LoaderProcess struct for each individual command to execute will be returned.
// Each individual LoaderProcess struct and associated process will be monitored
// and AEN will do its best to keep it running at all times.
func processesFromJSONFile(processesPath string) ([]LoaderProcess, error) {

	rawJSONMap := make(map[string]*json.RawMessage)
	var processList []LoaderProcess

	fileBytes, readErr := ioutil.ReadFile(processesPath)
	if readErr != nil {
		return nil, readErr
	}

	logger.Lgr.LogMessage("Successfully loaded process map bytes from file: %v", processesPath)

	mapErr1 := json.Unmarshal(fileBytes, &rawJSONMap)
	if mapErr1 != nil {
		return nil, mapErr1
	}

	logger.Lgr.LogMessage("Successfully unmarshalled JSON process file bytes into a map")

	for key, value := range rawJSONMap {
		var s string
		mapErr2 := json.Unmarshal(*value, &s)
		if mapErr2 != nil {
			return nil, mapErr2
		}

		commandParts := strings.Split(s, " ")
		lp := LoaderProcess{Name: key, Command: commandParts[0], Arguments: commandParts[1:]}

		logger.Lgr.LogMessage("Successfully created LoaderProcess instance: %v", lp.Name)

		logInstance, logError := logger.CustomLogger(lp.Name, 1, 50000, 604800)
		if logError != nil {
			return nil, logError
		}

		logger.Lgr.LogMessage("Successfully instantiated custom logger process for LoaderProcess: %v", lp.Name)

		lp.Lgr = logInstance
		processList = append(processList, lp)

		logger.Lgr.LogMessage("Successfully initialized one LoaderProcess instance: %+v", lp)
	}

	return processList, nil
}

// StartAsynchronous will execute all the processes that have been loaded into
// this specific instance of Loader asynchronously. It will capture their
// individual log output and put each specific process output in its own log
// file. It will also track how long each process runs for and return all this
// information inside an array of LoaderProcess.
func (ldr *Loader) StartAsynchronous() []LoaderProcess {

	var waitGroup sync.WaitGroup
	numProcesses := len(ldr.Processes)
	waitGroup.Add(numProcesses)

	logger.Lgr.LogMessage("Adding %d processes to the Asynchronous WaitGroup", numProcesses)

	for index := range ldr.Processes {

		go func(currentProcess *LoaderProcess) {

			defer waitGroup.Done()

			logger.Lgr.LogMessage("Asynchronously executing LoaderProcess: %+v", currentProcess)

			cmd := exec.Command(currentProcess.Command, currentProcess.Arguments...)
			cmd.Stdout = currentProcess.Lgr
			cmd.Stderr = currentProcess.Lgr

			currentProcess.Start = time.Now().Unix()
			err := cmd.Run()
			currentProcess.End = time.Now().Unix()
			currentProcess.Duration = currentProcess.End - currentProcess.Start

			if err != nil {
				currentProcess.Lgr.LogMessage("LoaderProcess:\n%+v\nexited with error status: %v", currentProcess, err.Error())
			} else {
				currentProcess.Lgr.LogMessage("LoaderProcess:\n%+v\nexited successfully", currentProcess)
			}

			logger.Lgr.LogMessage("Removing '%v' process from the Asynchronous WaitGroup. Execution took: %v", currentProcess.Name, currentProcess.Duration)

		}(&ldr.Processes[index]) // passing the current process using index
	}

	logger.Lgr.LogMessage("Waiting for %d processes to finish executing asynchronously", numProcesses)
	waitGroup.Wait()
	logger.Lgr.LogMessage("%d processes finished executing asynchronously. returning.", numProcesses)
	return ldr.Processes
}

// StartSynchronous will execute all the processes that have been loaded into
// this specific instance of Loader in series. It will capture their
// individual log output and put each specific process output in its own log
// file. It will also track how long each process runs for and return all this
// information inside an array of LoaderProcess.
func (ldr *Loader) StartSynchronous() []LoaderProcess {

	numProcesses := len(ldr.Processes)

	logger.Lgr.LogMessage("Executing %d processes in series", numProcesses)

	for _, currentProcess := range ldr.Processes {

		logger.Lgr.LogMessage("Synchronously executing LoaderProcess: %+v", currentProcess)

		cmd := exec.Command(currentProcess.Command, currentProcess.Arguments...)
		cmd.Stdout = currentProcess.Lgr
		cmd.Stderr = currentProcess.Lgr

		currentProcess.Start = time.Now().Unix()
		err := cmd.Run()
		currentProcess.End = time.Now().Unix()
		currentProcess.Duration = currentProcess.End - currentProcess.Start

		if err != nil {
			currentProcess.Lgr.LogMessage("LoaderProcess:\n%+v\nexited with error status: %v", currentProcess, err.Error())
		} else {
			currentProcess.Lgr.LogMessage("LoaderProcess:\n%+vexited successfully", currentProcess)
		}

		logger.Lgr.LogMessage("Finished executing one process out of %d", numProcesses)
	}

	logger.Lgr.LogMessage("%d processes finished executing synchronously. returning.", numProcesses)
	return ldr.Processes
}

// Run will continuously execute this specific instance of Loader indefinitely.
// Should only be called externally when all configuration options have been
// correctly setup and you wish to execute a set number of processes forever.
func (ldr *Loader) Run() {
	go func() {
		for 1 == 1 {
			ldr.StartAsynchronous()
		}
	}()
}
