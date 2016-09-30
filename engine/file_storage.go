package engine

import (
	"os"
	"path/filepath"
	"sync"
)

type fileLock interface {
	release() error
}

type fileStorage struct {
	path string
	mu   sync.RWMutex
}

// OpenFile opens file and returns file
func OpenFile(path string) (Storage, error) {
	if fi, err := os.Stat(path); err == nil {
		if !fi.IsDir() {
			return nil, err
		}
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	fs := fileStorage{
		path: path,
	}
	return &fs, nil
}

func (fs *fileStorage) Open(fd FileDesc) (Reader, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fpath := filepath.Join(fs.path, fd.Name())
	f, err := os.OpenFile(fpath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs *fileStorage) Create(fd FileDesc) (Writer, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fpath := filepath.Join(fs.path, fd.Name())
	f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs *fileStorage) Size(fd FileDesc) (int64, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fpath := filepath.Join(fs.path, fd.Name())
	stat, err := os.Stat(fpath)
	if err != nil {
		return 0, nil
	}
	return stat.Size(), nil
}

func (fs *fileStorage) Exists(fd FileDesc) bool {
	fpath := filepath.Join(fs.path, fd.Name())
	if _, err := os.Stat(fpath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (fs *fileStorage) Remove(fd FileDesc) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fpath := filepath.Join(fs.path, fd.Name())
	if err := os.Remove(fpath); err != nil {
		return nil
	}
	return nil
}

func (fs *fileStorage) Rename(ofd, nfd FileDesc) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	old := filepath.Join(fs.path, ofd.Name())
	new := filepath.Join(fs.path, nfd.Name())

	if err := os.Rename(old, new); err != nil {
		return nil
	}
	return nil
}

func (fs *fileStorage) Truncate(pos int64) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return os.Truncate(fs.path, pos)
}

func (fs *fileStorage) Close() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	return nil
}
