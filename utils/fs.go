package utils

import (
	"os"
	"path/filepath"
)

func Pwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		// Fatal(err)
	}
	return dir
}
