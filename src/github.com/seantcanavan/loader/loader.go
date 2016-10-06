package loader

import ()

type Loader struct {
	Processes     map[string]string
	WatchdogDelay uint64
	logFilePath   string
}

func (l Loader) Monitor() {

}
