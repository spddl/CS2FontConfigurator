package helper

import (
	"os"
)

// https://golang.org/src/path/path.go?s=4371:4399#L158
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
