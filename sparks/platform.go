package sparks

import (
	log "github.com/sirupsen/logrus"
	"leblox.com/sparks-cli/v2/sys"
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
	Code(Configuration) error
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
	enabledOne := false
	for _, name := range PlatformNames {
		if i < len(platforms) && platforms[i] {
			Platforms[name].SetEnabled(true)
			enabledOne = true
			log.Debugf("enabled platform %s", Platforms[name].Title())
		}
		i++
	}
	if !enabledOne {
		os, _ := sys.GetOs()
		if Platforms[os.Name()] != nil {
			Platforms[os.Name()].SetEnabled(true)
		}
	}
}
