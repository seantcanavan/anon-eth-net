package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
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
	processes []LoaderProcess // the slice of LoaderProcesses which the loader will execute and keep an eye on
}

type LoaderProcess struct {
	Name      string
	Command   string
	Arguments []string
	Lgr       *logger.Logger
}

// NewLoader will initialize a new instance of the Loader struct and load the
// associated processes from the given file. It will wait the given amount of
// time after a process exists before restarting it and it will use the given
// config reference to initialize the logger. It will probably utilize the
// config object more heavily in the future.
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
	loadedProcesses, loadErr := getProcessesFromJSONFile(processesPath)

	if loadErr != nil {
		return nil, loadErr
	}

	l.processes = loadedProcesses
	return &l, nil
}

// getProcessesFromJSONFile will read in a set of JSON values which define both
// the canonical name of the process as well as the command and any associated
// parameters to successfully execute that command. A slice containing a
// LoaderProcess struct for each individual command to execute will be returned.
// Each individual LoaderProcess struct and associated process will be monitored
// and AEN will do its best to keep it running at all times.
func getProcessesFromJSONFile(processesPath string) ([]LoaderProcess, error) {

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
	for _, currentProcess := range l.processes {
		cmd := exec.Command(currentProcess.Command, currentProcess.Arguments...)
		lgr.LogMessage("Asynchronously executing LoaderProcess: %+v", currentProcess)
		localProcess := currentProcess
		go func() {
			output, err := cmd.CombinedOutput()
			if err != nil {
				lgr.LogMessage("LoaderProcess exited with error status: %+v\n %v", localProcess, err.Error())
			} else {
				lgr.LogMessage("LoaderProcess exited successfully: %+v", localProcess)
				localProcess.Lgr.LogMessage(string(output))
			}
			time.Sleep(time.Second * TIME_BETWEEN_SUCCESSIVE_ITERATIONS)
		}()
	}
	return l.processes
}

// StartSynchronous will execute all the processes that have been loaded into
// this specific instance of Loader. It will execute them synchronously and
// return a slice of pointers to instances of os.File. Each instance of os.File
// contains the log output from each command that was executed.
func (l *Loader) StartSynchronous() []LoaderProcess {
	for _, currentProcess := range l.processes {
		cmd := exec.Command(currentProcess.Command, currentProcess.Arguments...)
		lgr.LogMessage("Synchronously executing LoaderProcess: %+v", currentProcess)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(fmt.Sprintf("LoaderProcess exited with error status: %+v", currentProcess))
			lgr.LogMessage("LoaderProcess exited with error status: %+v", currentProcess)
			currentProcess.Lgr.LogMessage("LoaderProcess exited with error status: %+v", currentProcess)
		} else {
			fmt.Println(fmt.Sprintf("LoaderProcess exited successfully: %+v", currentProcess))
			lgr.LogMessage("LoaderProcess exited successfully: %+v", currentProcess)
			currentProcess.Lgr.LogMessage("LoaderProcess exited successfully: %+v", currentProcess)
		}

		fmt.Println(fmt.Sprintf("Command output:\n%v", string(output)))
		lgr.LogMessage("Command output:\n%v", string(output))
		currentProcess.Lgr.LogMessage("Command output: %v", string(output))
	}
	return l.processes
}
