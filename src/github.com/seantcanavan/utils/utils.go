package utils

import (
	"fmt"
	"strconv"
	"time"
)

func FullDateString() string {
	return time.Now().String()
}

func UnixDateString() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// FullDateStringSafe returns the current time as a string with only file-name safe characters
func FullDateStringSafe() string {
	t := time.Now()
	return fmt.Sprintf("%v-%02d-%02d %02d_%02d_%02d.%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}
