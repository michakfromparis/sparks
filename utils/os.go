package utils

import (
	"errors"
	"runtime"
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
