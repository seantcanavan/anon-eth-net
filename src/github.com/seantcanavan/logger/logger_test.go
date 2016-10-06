package logger

import (
    "fmt"
    "testing"
)


func TestLogName(t *testing.T) {
    fmt.Println(LogFileHandle("test123"))
}
