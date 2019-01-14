package sparks

import log "github.com/Sirupsen/logrus"

// Platform Interfaxce
type Platform interface {
	Name() string
	Title() string
	Opt() string
	Enabled() bool
	SetEnabled(bool)

	Deps() error
	Clean() error
	Build(Configuration) error
}

// Map of all Platforms
var Platforms = map[string]Platform{}

// Ordered array of Platform keys
var PlatformNames = []string{
	"osx",
	"ios",
	"webgl",
}

func RegisterPlatform(platform Platform) {
	log.Debug("registering platform: " + platform.Title())
	Platforms[platform.Name()] = platform
}
