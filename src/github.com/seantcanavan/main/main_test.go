package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	result := m.Run()
	// flush logs
	os.Exit(result)
}
