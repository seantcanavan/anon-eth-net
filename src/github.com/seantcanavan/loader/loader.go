package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

// Loader represents a struct that will load a set of processes and watch over
// them. It knows the name of every process that it should be keeping an eye on
// as well as how to resurrect that process should it no longer be executing.
// The idea of the Loader is to make sure that all external process dependencies
// are executing and are in a healthy state as much as possible.
type Loader struct {
	processes        map[string]string // the map of process names to the command that's used to execute them. this data is utilized by the watchdog go routine to make sure that all required processes are running and executing as much as possible.
	ExitProcessDelay time.Duration     // the delay after a process exits, either successfully or unsuccessfully, before being started up again
}

func NewLoader(processesPath string, exitProcessDelay time.Duration) (*Loader, error) {
	l := Loader{}
	processMap, mapErr := getProcessMapFromFile(processesPath)
	if mapErr != nil {
		return nil, mapErr
	}
	l.processes = processMap
	l.ExitProcessDelay = exitProcessDelay
	return &l, nil
}

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

func (l *Loader) Start() {
	for key, value := range l.processes {
		go l.RunProcess(key, string(value))
	}
}

func (l *Loader) RunProcess(name string, command string) {
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
		time.Sleep(time.Second * l.ExitProcessDelay)
	}
}
