package utils

import (
	"errors"
	"os"
	"runtime"

	"github.com/Sirupsen/logrus"
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

func (o Os) String() string {
	return OsNames[o]
}

func GetOs() (Os, error) {
	switch os := runtime.GOOS; os {
	case "darwin":
		return Osx, nil
	case "linux":
		return Linux, nil
	case "windows":
		return Windows, nil
	default:
		return Unknown, errors.New("Unsupported host os: " + os)
	}
}

// import (
// 	"github.com/sirupsen/logrus"
// 	"os"
// )

// Term is the terminal logger for the Genesis application.
var Term *logrus.Logger

func setupTerm() {
	Term = logrus.New()
	Term.Out = os.Stdout

}

func init() {
	setupTerm()
}
