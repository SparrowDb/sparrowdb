package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SparrowDb/sparrowdb/util"
)

// CreateSnapshot creates snapshot of database using hard link
func CreateSnapshot(srcDir, dstDir string) error {
	uTime := fmt.Sprintf("%v", time.Now().UnixNano())
	ndestDir := dstDir + uTime
	tmpDir := dstDir + uTime + "_temp"

	// create target directory
	if ok, err := util.Exists(tmpDir); err == nil {
		if !ok {
			util.CreateDir(tmpDir)
		}
	} else {
		return err
	}

	// list of files in source directory
	var fileList []string
	fileList = append(fileList, listFilesDir(srcDir)...)

	// create hard link for each file
	createLink(srcDir, tmpDir, fileList)

	rename(tmpDir, ndestDir)

	return nil
}

func listFilesDir(dir string) []string {
	var fileList []string
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	return fileList
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}

func rename(src, dst string) error {
	return os.Rename(src, dst)
}

func createHardLink(src, dst string) error {
	return os.Link(src, dst)
}

func createLink(srcDir, targetDir string, fileList []string) error {
	for _, src := range fileList {
		fname := strings.Replace(src, srcDir, "", -1)
		dst := filepath.Join(targetDir, fname)

		if ok, err := isDirectory(src); err == nil {
			if ok {
				if fok, ferr := exists(targetDir); ferr == nil {
					if !fok {
						util.CreateDir(dst)
					}
				} else {
					return err
				}
			} else {
				if fok, ferr := exists(targetDir); ferr == nil {
					if !fok {
						createHardLink(src, dst)
					}
				} else {
					return err
				}
			}
		} else {
			return err
		}
	}

	return nil
}
