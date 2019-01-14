package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
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

func Execute(filename string, args ...string) (string, error) {
	fullCommand := fmt.Sprintf("%s %s", filename, strings.Join(args[:], " "))
	log.Debugf("executing %s", fullCommand)
	cmd := exec.Command(filename, args...)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		errx.Fatalf(err, "failed to execute "+fullCommand)
	}
	out := string(bytes)
	log.Tracef("combined output:%s%s", NewLine, out)
	return out, nil
}

func ExecuteEx(filename string, directoryName string, environment bool, args ...string) (string, error) {
	fullCommand := fmt.Sprintf("%s %s", filename, strings.Join(args[:], " "))
	log.Debugf("executing %s in directory %s with environment: %t", fullCommand, directoryName, environment)
	cmd := exec.Command(filename, args...)
	cmd.Dir = directoryName
	if environment {
		cmd.Env = os.Environ()
	}
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		errx.Fatalf(err, "failed to execute "+fullCommand)
	}
	out := string(bytes)
	log.Tracef("combined output:%s%s", NewLine, out)
	return out, nil
}

func Sed(filename string, regex string, newContent string) {
	log.Tracef("sed %s %s in %s", regex, newContent, filename)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		errx.Fatalf(err, "Could not read: "+filename)
	}
	re := regexp.MustCompile(regex)
	// log.Trace(string(bytes))
	matched := re.MatchString(string(bytes))
	if !matched {
		log.Warnf("regex %s did not match content of %s", regex, filename)
	}
	data := re.ReplaceAllString(string(bytes), newContent)
	// log.Trace(data)
	err = ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		errx.Fatalf(err, "Could not write to: "+filename)
	}
}
