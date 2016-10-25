package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/logger"
)

var log *logger.SeanLogger

// Loader represents a struct that will load a set of processes and watch over
// them. It knows the name of every process that it should be keeping an eye on
// as well as how to resurrect that process should it no longer be executing.
// The idea of the Loader is to make sure that all external process dependencies
// are executing and are in a healthy state as much as possible.
type Loader struct {
	processes        map[string]string // the map of process names to the command that's used to execute them. this data is utilized by the watchdog go routine to make sure that all required processes are running and executing as much as possible.
	exitProcessDelay time.Duration     // the delay after a process exits, either successfully or unsuccessfully, before being started up again
}

// NewLoader will initialize a new instance of the Loader struct and load the
// associated processes from the given file. It will wait the given amount of
// time after a process exists before restarting it and it will use the given
// config reference to initialize the logger. It will probably utilize the
// config object more heavily in the future.
func NewLoader(processesPath string, exitProcessDelay time.Duration) (*Loader, error) {
	if log == nil {
		newLogger, logError := logger.LoggerFromConservativeValue(config.Cfg.LoggingVolatility, "loader_package")
		if logError != nil {
			return nil, logError
		}
		log = newLogger
	}
	l := Loader{}
	processMap, mapErr := getProcessMapFromFile(processesPath)
	if mapErr != nil {
		return nil, mapErr
	}
	l.processes = processMap
	l.exitProcessDelay = exitProcessDelay
	return &l, nil
}

// getProcessMapFromFile will read in a set of JSON values which define both
// the canonical name of the process as well as the command and any associated
// parameters to successfully execute that command. This map will be used to
// continuously ensure that the write processes are running and that AEN is
// keeping them up and running as much as possible.
func getProcessMapFromFile(processesPath string) (map[string]string, error) {
	loadedMap := make(map[string]*json.RawMessage)
	returnedMap := make(map[string]string)
	fileBytes, readErr := ioutil.ReadFile(processesPath)
	if readErr != nil {
		return nil, readErr
	}

	mapErr1 := json.Unmarshal(fileBytes, &loadedMap)
	if mapErr1 != nil {
		return nil, mapErr1
	}

	for key, value := range loadedMap {
		var s string
		mapErr2 := json.Unmarshal(*value, &s)
		if mapErr2 != nil {
			return nil, mapErr2
		}
		returnedMap[key] = s
	}

	for key, value := range returnedMap {
		fmt.Println("key: " + key)
		fmt.Println("value: " + value)
	}

	return returnedMap, nil
}

// Start will execute all the processes that have been loaded on the local
// system. It will run each process in its own go routine.
func (l *Loader) Start() {
	for key, value := range l.processes {
		go l.runProcess(key, string(value))
	}
}

// runProcess handles the individual running of a single conceptual process or
// program. It will create a process, pass the associated arguments, and
// initialize it all within go. It will also monitor the stdout and stderr
// channels for output and return those when the program finishes executing.
// TODO(Canavan): figure out how to monitor the program output during execution
// so that it can be actively logged and reported back to the user.
func (l *Loader) runProcess(name string, command string) {
	for 1 == 1 {
		// split the incoming command into its primary command and its parameters which follow after
		commandParts := strings.Split(command, " ")
		command := commandParts[0]
		arguments := commandParts[1:]
		cmd := exec.Command(command, arguments...)
		fmt.Println("Loader: " + command + " is about to start.")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Loader: '" + command + "' exited with error status.")
			fmt.Println("Loader: '" + command + "' Error= '" + err.Error() + "'")
			fmt.Println("Loader: '" + command + "' output= '" + string(out) + "'")
		} else {
			fmt.Println(string(out))
		}
		time.Sleep(time.Second * l.exitProcessDelay)
	}
}
