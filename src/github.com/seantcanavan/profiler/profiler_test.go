package profiler

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/utils"
)

func TestMain(m *testing.M) {
	flag.Parse()

	assetPath, assetErr := utils.AssetPath("config.json")
	if assetErr != nil {
		fmt.Println(assetErr)
		return
	}

	err := config.FromFile(assetPath)
	if err != nil {
		return
	}

	result := m.Run()
	// flush logs
	os.Exit(result)
}

func TestSendFileArchiveAsAttachment(t *testing.T) {
	filePtr, err := SendArchiveProfileAsAttachment()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("Archive successfully created: %v", filePtr.Name()))
}
