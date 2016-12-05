package profiler

import (
	"fmt"
	"os"
	"testing"

	"github.com/seantcanavan/config"
)

func TestMain(m *testing.M) {

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
