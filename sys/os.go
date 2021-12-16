package sys

import (
	"errors"
	"os/user"
	"runtime"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
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

var osNames = []string{
	"unknown",
	"osx",
	"linux",
	"windows",
}

var osTitles = []string{
	"Unknown",
	"OSX",
	"Linux",
	"Windows",
}

// Name is the operating system lowercase name / identifier
func (o Os) Name() string {
	return osNames[o]
}

// Title is the operating system pretty name
func (o Os) Title() string {
	return osTitles[o]
}

func (o Os) String() string {
	return o.Title()
}

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
