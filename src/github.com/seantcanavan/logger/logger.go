package logger

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/seantcanavan/utils"
)

// SeanLogger allows for aggressive log management in scenarios where disk space might be limited.
// You can limit based on log message count or duration and also prune log files when too many are saved on disk.
type SeanLogger struct {
	MaxLogFileCount    uint64        // The maximum number of log files saved to disk before pruning occurs
	MaxLogMessageCount uint64        // The maximum number of bytes a log file can take up before it's cut off and a new one is created
	MaxLogDuration     int64         // The maximum number of seconds a log can exist for before it's cut off and a new one is created
	baseLogName        string        // The beginning text to append to this log instance for naming and management purposes
	logFileCount       uint64        // The current number of logs that have been created
	logFileNames       list.List     // The list of log files we're currently holding on to
	logMessageCount    uint64        // The current number of messages that have been logged
	logDuration        int64         // The duration, in seconds, that this log has been logging for
	logStamp           int64         // The time when this log was last written to in unix time
	log                *os.File      // The file that we're logging to
	writer             *bufio.Writer // our writer we use to log to the current log file
}

// LogFileHandle will generate a string name of a file based off of an initial
// string and append the date to the end to signify when it was created.
func LogFileHandle(logBaseName string) string {

	dts := utils.FullDateStringSafe()

	var nameBuffer bytes.Buffer
	nameBuffer.WriteString(logBaseName)
	nameBuffer.WriteString(" ")
	nameBuffer.WriteString(dts)
	nameBuffer.WriteString(".log")

	return nameBuffer.String()
}

// StartLog initializes all the log tracking variables
func (sl *SeanLogger) StartLog(logBaseName string) error {

	logFileName := LogFileHandle(logBaseName)

	filePtr, err := os.Create(logFileName)
	if err != nil {
		return err
	}

	// init / reset the log trackers
	sl.baseLogName = logBaseName
	sl.logFileCount = 0
	sl.logDuration = 0
	sl.logStamp = time.Now().Unix()
	sl.log = filePtr
	sl.writer = bufio.NewWriter(sl.log)
	sl.logFileNames.PushBack(logFileName)
	return nil
}

// LogMessage will write the given string to the log file. It will then perform
// all the necessary checks to make sure that the max number of messages, the
// max duration of the log file, and the maximum number of overall log files
// has not been reached. If any of the above parameters have been tripped,
// log cleanup will occur.
func (sl *SeanLogger) LogMessage(message string) {

	now := time.Now().Unix()

	fmt.Fprintln(sl.writer, message)

	sl.logMessageCount++
	sl.logDuration += now - sl.logStamp
	sl.logStamp = now

	if sl.logMessageCount >= sl.MaxLogMessageCount ||
		sl.logDuration >= sl.MaxLogDuration {
		sl.newFile()
	}
}

// newFile generates a new log file to store the log messages within. It
// intelligently keeps track of the number of log files that have already been
// created so that you don't overload your disk with logs and can 'prune' extra
// logs as necessary.
func (sl *SeanLogger) newFile() {

	logFileName := LogFileHandle(sl.baseLogName)

	filePtr, err := os.Create(logFileName)
	if err != nil {
		sl.handleCreateError()
		return
	}

	sl.writer.Flush()
	sl.log.Close()

	sl.log = filePtr
	sl.writer = bufio.NewWriter(sl.log)

	sl.logMessageCount = 0
	sl.logFileCount++
	sl.logFileNames.PushBack(logFileName)

	if sl.logFileCount >= sl.MaxLogFileCount {
		sl.pruneFile()
	}
}

func (sl *SeanLogger) pruneFile() {

	oldestLog := sl.logFileNames.Remove(sl.logFileNames.Front())
	logFileName := reflect.ValueOf(oldestLog).String()

	fmt.Println("Deleting oldest log file: %v", logFileName)

	os.Remove(logFileName)
}

func (sl *SeanLogger) handleCreateError() {
	// send last 3 log files, generate status report, email out update
}
