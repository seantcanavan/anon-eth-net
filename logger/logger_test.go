package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/seantcanavan/utils"
)

func TestMain(m *testing.M) {

	logErr := StandardLogger("logger_test")
	if logErr != nil {
		fmt.Println(fmt.Sprintf("Could not initialize logger: %v", logErr))
		return
	}

	result := m.Run()
	os.Exit(result)
}

// TestLogger will test all logger functionality
func TestLogger(t *testing.T) {

	logBaseName := "logger_package"

	logNameAsset, assetErr := utils.AssetPath("logger_test.sample")
	if assetErr != nil {
		t.Error(assetErr)
	}

	testFileLines, readErr := utils.ReadLines(logNameAsset)
	if readErr != nil {
		t.Error(readErr)
	}

	maxFileCount := uint64(2)
	maxMessageCount := uint64(600)
	maxDuration := uint64(10)

	sl1, logErr := CustomLogger(logBaseName, maxFileCount, maxMessageCount, maxDuration)

	if logErr != nil {
		t.Error(logErr)
	}

	for _, currentLine := range testFileLines {
		sl1.LogMessage(currentLine)
	}
}
