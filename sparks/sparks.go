package sparks

import (
	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
)

var CurrentProduct = Product{}

func Init() {
	log.Info("sparks init")
	if err := config.Init(); err != nil {
		errx.Fatalf(err, "Configuration initialization failed")
	}
}

func Shutdown() {
}

func Load() {
	log.Info("sparks load")
	CurrentProduct.Load()
	log.Trace("%+v", CurrentProduct)
}

func Save() {
	log.Info("sparks save")
	CurrentProduct.Save()
}

func Build() {
	log.Info("sparks build")
	Load()
	for _, platformName := range PlatformNames {
		platform := Platforms[platformName]
		if platform != nil && platform.Enabled() {
			for _, configurationName := range ConfigurationNames {
				configuration := Configurations[configurationName]
				if configuration != nil && configuration.Enabled() {
					Load()
					platform.Build(configuration)
				}
			}
		}
	}
}
