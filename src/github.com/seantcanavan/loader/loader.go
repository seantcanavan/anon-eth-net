package loader

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/seantcanavan/logger"
)

var lgr *logger.Logger

const TIME_BETWEEN_SUCCESSIVE_ITERATIONS = 60

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

	if lgr == nil {
		newLogger, logError := logger.FromVolatilityValue("loader_package")
		if logError != nil {
			return nil, logError
		}
		lgr = newLogger
	}

	l := Loader{}
	var loadedProcesses []LoaderProcess
	loadedProcesses, loadErr := processesFromJSONFile(processesPath)

	if loadErr != nil {
		return nil, loadErr
	}

	l.Processes = loadedProcesses
	return &l, nil
}

// getProcessesFromJSONFile will read in a set of JSON values which define both
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

	mapErr1 := json.Unmarshal(fileBytes, &rawJSONMap)
	if mapErr1 != nil {
		return nil, mapErr1
	}

	for key, value := range rawJSONMap {
		var s string
		mapErr2 := json.Unmarshal(*value, &s)
		if mapErr2 != nil {
			return nil, mapErr2
		}

		lp := LoaderProcess{}
		lp.Name = key
		commandParts := strings.Split(s, " ")
		lp.Command = commandParts[0]
		lp.Arguments = commandParts[1:]

		logInstance, logError := logger.FromVolatilityValue(lp.Name)
		if logError != nil {
			lgr.LogMessage("LoaderProcess unsuccessfully initialized logger: %+v", lp)
			return nil, logError
		}

		lp.Lgr = logInstance
		lgr.LogMessage("Read process successfully from file: %+v", lp)
		processList = append(processList, lp)
	}
	return processList, nil
}

// StartAsynchronous will execute all the processes that have been loaded into
// this specific instance of Loader. It will execute them asynchronously and
// eventually in the future it will hopefully figure out a meaningful way of
// logging the output of each individual process...
func (l *Loader) StartAsynchronous() []LoaderProcess {
	var waitGroup sync.WaitGroup
	lgr.LogMessage("Adding %d processes to the Asynchronous WaitGroup", len(l.Processes))
	waitGroup.Add(len(l.Processes))

	for index := range l.Processes {
		go func(currentProcess *LoaderProcess) {
			defer waitGroup.Done()
			cmd := exec.Command(currentProcess.Command, currentProcess.Arguments...)
			lgr.LogMessage("Asynchronously executing LoaderProcess: %+v", currentProcess)
			localProcess := currentProcess
			currentProcess.Start = time.Now().Unix()
			output, err := cmd.CombinedOutput()
			currentProcess.End = time.Now().Unix()
			currentProcess.Duration = currentProcess.End - currentProcess.Start
			if err != nil {
				lgr.LogMessage("LoaderProcess exited with error status: %+v\n %v", localProcess, err.Error())
			} else {
				lgr.LogMessage("LoaderProcess exited successfully: %+v", localProcess)
			}
			currentProcess.Lgr.LogMessage("LoaderProcess: %+v", currentProcess)
			currentProcess.Lgr.LogMessage("Command output:\n%v", string(output))
			currentProcess.Lgr.Flush()
			lgr.LogMessage("Removing '%v' process from the Asynchronous WaitGroup. Execution took: %v", currentProcess.Name, currentProcess.Duration)
		}(&l.Processes[index]) // passing the current process using index
	}
	waitGroup.Wait()
	lgr.Flush()
	return l.Processes
}

// StartSynchronous will execute all the processes that have been loaded into
// this specific instance of Loader. It will execute them synchronously and
// return a slice of pointers to instances of os.File. Each instance of os.File
// contains the log output from each command that was executed.
func (l *Loader) StartSynchronous() []LoaderProcess {
	for _, currentProcess := range l.Processes {
		cmd := exec.Command(currentProcess.Command, currentProcess.Arguments...)
		lgr.LogMessage("Synchronously executing LoaderProcess: %+v", currentProcess)
		currentProcess.Start = time.Now().Unix()
		output, err := cmd.CombinedOutput()
		currentProcess.End = time.Now().Unix()
		currentProcess.Duration = currentProcess.End - currentProcess.Start
		if err != nil {
			lgr.LogMessage("LoaderProcess exited with error status: %+v", currentProcess)
		} else {
			lgr.LogMessage("LoaderProcess exited successfully: %+v", currentProcess)
		}
		currentProcess.Lgr.LogMessage("LoaderProcess: %+v", currentProcess)
		currentProcess.Lgr.LogMessage("Command output:\n%v", string(output))
		currentProcess.Lgr.Flush()
	}
	lgr.Flush()
	return l.Processes
}
