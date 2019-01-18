package sys

import (
	"errors"
	"os/user"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/joomcode/errorx"
)

type Os int

const (
	Unknown = iota
	Osx
	Linux
	Windows
)

var OsNames = []string{
	"unknown",
	"osx",
	"linux",
	"windows"}

var hostOs Os

func GetOs() (Os, error) {
	if hostOs != Unknown {
		return hostOs, nil
	}
	switch os := runtime.GOOS; os {
	case "darwin":
		hostOs = Osx
	case "linux":
		hostOs = Linux
	case "windows":
		hostOs = Windows
		return Windows, nil
	default:
		return Unknown, errors.New("Unsupported host os: " + os)
	}
	return hostOs, nil
}

var homeDirectory = ""

func GetHome() (string, error) {
	if homeDirectory != "" {
		return homeDirectory, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", errorx.Decorate(err, "could not access current user environment")
	}
	homeDirectory = usr.HomeDir
	log.Debugf("$HOME Directory: %s", homeDirectory)
	return homeDirectory, nil
}
