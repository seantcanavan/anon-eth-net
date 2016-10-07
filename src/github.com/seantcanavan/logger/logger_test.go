package logger

import (
	"testing"

	"github.com/seantcanavan/utils"
)

// TestLogger will test all logger functionality
func TestLogger(t *testing.T) {

	sl1 := &SeanLogger{
		MaxLogFileCount:    3,   // 10 max log files
		MaxLogMessageCount: 300, // 10 max messages per log
		MaxLogDuration:     10,  // 10 max seconds per log file
	}

	logBaseName := "test01"
	sampleLogFileName := "logger_test_sample.txt"
	testFileLines, err := utils.ReadLines(sampleLogFileName)

	if err != nil {
		t.Error(err)
	}

	sl1.StartLog(logBaseName)

	for _, currentLine := range testFileLines {
		sl1.LogMessage(currentLine)
	}

	// directoryFiles := utils.DirectoryList(".", logBaseName)
}
