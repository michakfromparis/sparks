package sys

import (
	"errors"
	"os/user"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/joomcode/errorx"
)

// Os represents the operating system the code in running on
type Os int

const (
	// Unknown OS was not detected
	Unknown = iota
	// Osx Os
	Osx
	// Linux Os
	Linux
	// Windows Os
	Windows
)

var hostOs Os

// GetOs returns the Os enum representing the currently running operating system
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

// GetHome returns the fullpasth to the user's home directory
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
