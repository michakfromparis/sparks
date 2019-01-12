package utils

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

func Pwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
