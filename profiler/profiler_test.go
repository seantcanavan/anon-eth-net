package profiler

import (
	"fmt"
	"os"
	"testing"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/logger"
)

func TestMain(m *testing.M) {

	logErr := logger.StandardLogger("profiler_test")
	if logErr != nil {
		fmt.Println(fmt.Sprintf("Could not initialize logger: %v", logErr))
		return
	}

	configErr := config.FromFile()
	if configErr != nil {
		return
	}

	result := m.Run()
	os.Exit(result)
}

func TestSendFileArchiveAsAttachment(t *testing.T) {
	filePtr, err := SendArchiveProfileAsAttachment()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("Archive successfully created: %v", filePtr.Name()))
}
