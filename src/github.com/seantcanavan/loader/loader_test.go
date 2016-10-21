package loader

import (
	"fmt"
	"testing"
	"time"
)

func TestProcessMap(t *testing.T) {
	loader, err := NewLoader("loader.json", 30)
	if err != nil {
		t.Error(err)
	}

	loader.Start()
	fmt.Println("about to sleep")
	time.Sleep(time.Second * 30)
	fmt.Println("done sleeping")
}
