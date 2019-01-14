package platform

import (
	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/sparks"
)

func RegisterPlatforms() {
	sparks.RegisterPlatform(&Osx{})
	sparks.RegisterPlatform(&Ios{})
	sparks.RegisterPlatform(&WebGl{})
}

func SetEnabledPlatforms(platforms []bool) {
	i := 0
	for _, name := range sparks.PlatformNames {
		if i < len(platforms) && platforms[i] == true {
			sparks.Platforms[name].SetEnabled(true)
			log.Debugf("enabled platform %s", sparks.Platforms[name].Title())
		}
		i++
	}
}
