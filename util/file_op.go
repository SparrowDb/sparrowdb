package util

import (
	"io/ioutil"
	"os"
)

// Exists check if dir/file exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// CreateDir creates directory
func CreateDir(x string) error {
	return os.MkdirAll(x, os.ModePerm)
}

// DeleteDir delete directory
func DeleteDir(x string) error {
	return os.RemoveAll(x)
}

// CreateEmptyFile creates empty file
func CreateEmptyFile(filepath string) {
	ioutil.WriteFile(filepath, []byte(""), 0644)
}
