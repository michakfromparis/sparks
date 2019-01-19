package sparks

import (
	log "github.com/Sirupsen/logrus"
)

// Platform Interface used to represent a sparks supported platform
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

// Platforms is the map of all registered platforms
var Platforms = map[string]Platform{}

// PlatformNames is an ordered array of Platform keys used to iterate over Platforms
var PlatformNames = []string{
	"osx",
	"ios",
	"android",
	"windows",
	"linux",
	"webgl",
}

// RegisterPlatform allows external code to register a new platform as long as it respects the Platform interface
func RegisterPlatform(platform Platform) {
	log.Debug("registering platform: " + platform.Title())
	Platforms[platform.Name()] = platform
}

// SetEnabledPlatforms is used to enable / disable build platforms, platforms parameter comes ordered like PlatformNames
func SetEnabledPlatforms(platforms []bool) {
	i := 0
	for _, name := range PlatformNames {
		if i < len(platforms) && platforms[i] {
			Platforms[name].SetEnabled(true)
			log.Debugf("enabled platform %s", Platforms[name].Title())
		}
		i++
	}
}
