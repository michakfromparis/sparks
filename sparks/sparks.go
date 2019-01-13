package sparks

import (
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
)

func Init() {
	if err := config.Init(); err != nil {
		errx.FatalF(err, "Configuration initialization failed")
	}
}

func Shutdown() {
}

func LoadProjectConfig() {
}

func Build() {
	for _, platformName := range PlatformNames {
		platform := Platforms[platformName]
		if platform != nil && platform.Enabled() {
			for _, configurationName := range ConfigurationNames {
				configuration := Configurations[configurationName]
				if configuration != nil && configuration.Enabled() {
					platform.Build(configuration)
				}
			}
		}
	}
}
