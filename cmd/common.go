package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/platform"
	"github.com/michaKFromParis/sparks/utils"
)

func precmd() {
	c := 0
	i := 0
	for range platform.Platforms {
		if enabledPlatforms[i] {
			c++
		}
		i++
	}
	// setting default platform if no platform selected
	if c == 0 {

		default_platform := utils.GetOs()
		i = 0
		for name, _ := range platform.Platforms {
			if name == default_platform {
				enabledPlatforms[i] = true
			}
			i++
		}
	}
	// setting default configuration if no configuration selected
	if !config.Debug && !config.Release && !config.Shipping {
		config.Release = true
	}

	// setting enabled platforms
	i = 0
	for _, p := range platform.Platforms {
		if enabledPlatforms[i] {
			config.Platforms[p.Name()] = true
		}
		i++
	}
	log.Info(config.String())

}

func LoadProjectConfig() {

}

func postcmd() {
}
