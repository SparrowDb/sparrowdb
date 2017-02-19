package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

// ListDirectory returns file list in directory
func ListDirectory(dir string) ([]string, error) {
	fileList := []string{}

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(fileList) > 1 {
		fileList = append(fileList[:0], fileList[1:]...)
	}

	return fileList, nil
}
