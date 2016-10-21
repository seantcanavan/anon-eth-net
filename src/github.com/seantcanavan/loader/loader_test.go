package loader

import (
	"testing"
)

func TestProcessMap(t *testing.T) {
	loader, err := NewLoader("loader.json", 500)
	if err != nil {
		t.Error(err)
	}

	loader.Start()
}
