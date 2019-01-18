package sparks

import (
	log "github.com/Sirupsen/logrus"
)

// Platform Interfaxce
type Platform interface {
	Name() string
	Title() string
	Opt() string
	Enabled() bool
	SetEnabled(bool)

	Get() error
	Clean() error
	Build(Configuration) error
}

// Map of all Platforms
var Platforms = map[string]Platform{}

// Ordered array of Platform keys
var PlatformNames = []string{
	"osx",
	"ios",
	"android",
	"windows",
	"linux",
	"webgl",
}

func RegisterPlatform(platform Platform) {
	log.Debug("registering platform: " + platform.Title())
	Platforms[platform.Name()] = platform
}

func SetEnabledPlatforms(platforms []bool) {
	i := 0
	for _, name := range PlatformNames {
		if i < len(platforms) && platforms[i] == true {
			Platforms[name].SetEnabled(true)
			log.Debugf("enabled platform %s", Platforms[name].Title())
		}
		i++
	}
}
