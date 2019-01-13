package utils

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/errx"
)

func Pwd() (string, error) {
	log.Trace("pwd")
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		errx.Fatalf(err, "Could not determine current working directory")
	}
	log.Debug("pwd: " + path)
	return path, err
}
