package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/errx"
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

func Execute(name string, args ...string) string {
	fullCommand := fmt.Sprintf("%s %s", name, strings.Join(args[:], " "))
	log.Debugf("executing %s", fullCommand)
	cmd := exec.Command(name, args...)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		errx.Fatalf(err, "failed to execute "+fullCommand)
	}
	out := string(bytes)
	log.Tracef("combined output:%s%s", LineBreak, out)
	return out
}
