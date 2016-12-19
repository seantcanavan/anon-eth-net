package logger

import (
	"bufio"
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/seantcanavan/utils"
)

// The file extension to use for all new log files that are created
const LOG_EXTENSION = ".log"

var Lgr *Logger

// Logger allows for aggressive log management in scenarios where disk space
// might be limited. You can limit based on log message count or duration and
// also prune log files when too many are saved on disk.
type Logger struct {
	MaxLogFileCount    uint64        // The maximum number of log files saved to disk before pruning occurs
	MaxLogMessageCount uint64        // The maximum number of bytes a log file can take up before it's cut off and a new one is created
	MaxLogDuration     uint64        // The maximum number of seconds a log can exist for before it's cut off and a new one is created
	baseLogName        string        // The beginning text to append to this log instance for naming and management purposes
	logFileCount       uint64        // The current number of logs that have been created
	logFileNames       list.List     // The list of log files we're currently holding on to
	logMessageCount    uint64        // The current number of messages that have been logged
	logDuration        uint64        // The duration, in seconds, that this log has been logging for
	logStamp           uint64        // The time when this log was last written to in unix time
	log                *os.File      // The file that we're logging to
	writer             *bufio.Writer // our writer we use to log to the current log file
	lock 			   sync.Mutex
}

// CustomLogger returns a logger with the given variables customized to your
// liking. Smaller values are better for devices with less free space and vice
// versa for devices with more free space.
func CustomLogger(logBaseName string, maxFileCount uint64, maxMessageCount uint64, maxDuration uint64) (*Logger, error) {

	lgr := &Logger{
		MaxLogFileCount:    maxFileCount,
		MaxLogMessageCount: maxMessageCount,
		MaxLogDuration:     maxDuration,
	}

	err := lgr.initLogger(logBaseName)
	if err != nil {
		return nil, err
	}

	fmt.Println(fmt.Sprintf("Successfully initialized custom logger: %+v", lgr))

	return lgr, nil
}

// StandardLogger will return a Logger struct which will hoard a massive
// amount of logs and messages. Recommended for systems with a healthy amount of
// free disk space. Logs can be left unchecked for up to 7 days before they're
// pruned. If you don't want to check them and you don't want to lose them then
// make sure you download them via REST otherwise you'll miss log data.
func StandardLogger(logBaseName string) error {

	lgr := &Logger{
		MaxLogFileCount:    1000,   // up to 1000 max log files simultaneously stored on disk
		MaxLogMessageCount: 10000,  // a new log file every 10,000 messages
		MaxLogDuration:     604800, // a new log file every 7 days
	}

	err := lgr.initLogger(logBaseName)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Successfully initialized standard logger: %+v", lgr))

	Lgr = lgr
	return nil
}

// CurrentLogContents returns the contents of the current log file that's being
// managed by the logger instance. The current log should be active thus
// multiple calls to CurrentLogContents() should give different results.
func (lgr *Logger) CurrentLogContents() ([]byte, error) {

	lgr.writer.Flush()

	fileBytes, readErr := ioutil.ReadFile(lgr.log.Name())
	if readErr != nil {
		return nil, readErr
	}

	Lgr.LogMessage("Successfully retrieved current log contents")

	return fileBytes, nil
}

// CurrentLogName returns the name of the file which contains the most current
// log output. This log will be active and likely changing frequently.
func (lgr *Logger) CurrentLogName() (string, error) {

	fileInfo, statErr := lgr.log.Stat()
	if statErr != nil {
		return "", statErr
	}

	Lgr.LogMessage("Successfully retrieved current log name: %v", fileInfo.Name())

	return fileInfo.Name(), nil
}

// CurrentLogFile returns a pointer to the current os.File representation of
// the current log file that is being written to. If this reference is held
// log enough it can become invalid if the log file is pruned from the disk.
func (lgr *Logger) CurrentLogFile() *os.File {
	return lgr.log
}

// initLogger will initialize all of the helper values required to maintain a
// circular array of log files. When you reach the end of the circle the log
// is 'pruned'.
func (lgr *Logger) initLogger(logBaseName string) error {

	logFileName := utils.TimeStampFileName(logBaseName, LOG_EXTENSION)

	filePtr, err := os.Create(logFileName)
	if err != nil {
		return err
	}

	// private variable
	lgr.baseLogName = logBaseName
	lgr.logFileCount = 0
	lgr.logDuration = 0
	lgr.logStamp = uint64(time.Now().Unix())
	lgr.log = filePtr
	lgr.writer = bufio.NewWriter(lgr.log)
	lgr.logFileNames.PushBack(logFileName)

	lgr.LogMessage("Successfully created initial log file: %v", filePtr.Name())

	return nil
}

// LogMessage will write the given string to the current active log file. It
// will then perform all the necessary checks to make sure that the max number
// of messages, the max duration of the log file, and the maximum number of
// overall log files has not been reached. If any of the above parameters have
// been tripped, action will be taken accordingly.
func (lgr *Logger) LogMessage(formatString string, values ...interface{}) {

	lgr.lock.Lock()
	defer lgr.lock.Unlock()

	// what time is it right now?
	now := uint64(time.Now().Unix())
	// write the logging message to the current log file
	fmt.Fprintln(lgr.writer, fmt.Sprintf(formatString, values...))
	// write the logging message to std.out for local watchers
	fmt.Println(fmt.Sprintf(formatString, values...))
	// manually flush for now... it ain't pretty but it works
	lgr.writer.Flush()

	lgr.logMessageCount++
	lgr.logDuration += now - lgr.logStamp
	lgr.logStamp = now

	if lgr.logMessageCount >= lgr.MaxLogMessageCount ||
		lgr.logDuration >= lgr.MaxLogDuration {
		lgr.newFile()
	}
}

// newFile generates a new log file to store the log messages within. It
// intelligently keeps track of the number of log files that have already been
// created so that you don't overload your disk with logs and can 'prune' extra
// logs as they pass the threshold to keep around.
func (lgr *Logger) newFile() error {

	logFileName := utils.TimeStampFileName(lgr.baseLogName, LOG_EXTENSION)

	filePtr, err := os.Create(logFileName)
	if err != nil {
		return err
	}

	Lgr.LogMessage("Created new log file: %v", filePtr.Name())

	lgr.log.Close()

	Lgr.LogMessage("Successfully closed the old log file: %v", lgr.CurrentLogFile().Name())

	lgr.log = filePtr
	lgr.writer = bufio.NewWriter(lgr.log)

	lgr.logMessageCount = 0
	lgr.logFileCount++
	lgr.logFileNames.PushBack(logFileName)

	if lgr.logFileCount >= lgr.MaxLogFileCount {
		if err := lgr.pruneFile(); err != nil {
			return err
		}
	}

	return nil
}

// pruneFile will remove the oldest file handle from the queue and delete the
// file from the local file system.
func (lgr *Logger) pruneFile() error {

	oldestLog := lgr.logFileNames.Remove(lgr.logFileNames.Front())
	logFileName := reflect.ValueOf(oldestLog).String()

	Lgr.LogMessage("Deleting oldest log file: %v", logFileName)
	return os.Remove(logFileName)
}
