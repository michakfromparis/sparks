package utils

import (
	"runtime"

	log "github.com/Sirupsen/logrus"
)

func GetOs() string {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "osx"
	case "linux":
		return os
	case "windows":
		return os
	default:
		log.Fatal("Unsupported host os: " + os)
		return ""
	}
}
