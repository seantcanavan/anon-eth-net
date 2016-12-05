package utils

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {

	result := m.Run()
	// flush logs
	os.Exit(result)
}

func TestGetAssetPathPass(t *testing.T) {
	relativePath, pathErr := AssetPath("config.json")
	if pathErr != nil {
		t.Error(pathErr)
	}
	fmt.Println(fmt.Sprintf("relativePath: %v", relativePath))
}

func TestGetSysAssetPathPass(t *testing.T) {
	relativePath, pathErr := SysAssetPath("loader_test.json")
	if pathErr != nil {
		t.Error(pathErr)
	}
	fmt.Println(fmt.Sprintf("relativePath: %v", relativePath))
}
