package sparks

import (
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/logger"
)

func Initialize() {

	if err := logger.Init(); err != nil {
		errx.FatalF(err, "Logger initialization failed")
	}
	if err := config.Init(); err != nil {
		errx.FatalF(err, "Configuration initialization failed")
	}
	// if err := platform.Init(); err != nil {
	// 	errx.FatalF(err, "Platforms initialization failed")
	// }
}

func Shutdown() {
}

func LoadProjectConfig() {

}

func Build() {
	for _, platformName := range PlatformNames {
		platform := Platforms[platformName]
		if platform.Enabled() {
			for _, configurationName := range ConfigurationNames {
				configuration := Configurations[configurationName]
				if configuration.Enabled() {
					platform.Build(configuration)
				}
			}
		}
	}
}
