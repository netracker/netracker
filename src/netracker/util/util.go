package util

import (
	"path/filepath"
	"runtime"
)

func currentDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

func RelativePath(path string) string {
	return filepath.Clean(currentDir() + path)
}
