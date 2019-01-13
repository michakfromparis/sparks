package sparks

import log "github.com/Sirupsen/logrus"

// Platform Interfaxce
type Platform interface {
	Name() string
	Title() string
	Opt() string
	Enabled() bool
	SetEnabled(bool)

	Deps()
	Clean()
	Build(Configuration)
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
	log.Debug("Registering platform: " + platform.Title())
	Platforms[platform.Name()] = platform
}
