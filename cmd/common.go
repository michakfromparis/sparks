package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/platform"
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
		i = 0
		for name, _ := range platform.Platforms {
			if name == "osx" {
				enabledPlatforms[i] = true
			}
			i++
		}
		config.Platforms["osx"] = true
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
	log.Info(config.Format())

}
func postcmd() {
}
