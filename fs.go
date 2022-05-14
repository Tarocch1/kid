package kid

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Filesystem implement http.FileSystem
type FileSystem struct {
	// Index page filename, default: index.html
	Index string

	// Root dir, default: os.Getwd()
	Root string

	// Enable open file out of root path, default: false
	EnableOuter bool

	// FS that implement http.FileSystem, default: os
	FS http.FileSystem
}

func (fs *FileSystem) open(path string) (http.File, error) {
	if fs.FS != nil {
		return fs.FS.Open(path)
	}
	return os.Open(path)
}

// Open opens file at given path
func (fs *FileSystem) Open(path string) (http.File, error) {
	var err error
	index := _if(fs.Index != "", fs.Index, "index.html")
	root := fs.Root
	if root == "" {
		root, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	if !fs.EnableOuter && containsDotDot(path) {
		return nil, NewError(http.StatusBadRequest, "400 Bad Request: Invalid path")
	}

	var target string

	if !fs.EnableOuter {
		target = filepath.Join(root, path)
	} else {
		target = _if(filepath.IsAbs(path), path, filepath.Join(root, path))
	}

	file, err := fs.open(target)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	for {
		if !stat.IsDir() {
			break
		}
		file.Close()
		target = filepath.Join(target, index)

		file, err = fs.open(target)
		if err != nil {
			return nil, err
		}

		stat, err = file.Stat()
		if err != nil {
			file.Close()
			return nil, err
		}
	}

	return file, nil
}

func containsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}
	return false
}

func isSlashRune(r rune) bool { return r == '/' || r == '\\' }
