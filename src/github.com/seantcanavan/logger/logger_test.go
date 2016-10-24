package logger

import (
	"testing"

	"github.com/seantcanavan/utils"
)

// TestLogger will test all logger functionality
func TestLogger(t *testing.T) {

	logBaseName := "test01"
	sampleLogFileName := "logger_test_sample.txt"
	testFileLines, err := utils.ReadLines(sampleLogFileName)

	if err != nil {
		t.Error(err)
	}

	sl1, logErr := CustomLogger(logBaseName, 3, 300, 10)

	if logErr != nil {
		t.Error(logErr)
	}

	for _, currentLine := range testFileLines {
		sl1.LogMessage(currentLine)
	}
}
