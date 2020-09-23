package main

import (
	"path/filepath"
	"runtime"
)

// RunDir returns the directory of the runtime.
func RunDir() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Dir(file)
}
