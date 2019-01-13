package sparks

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
)

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
	log.Tracef("opening %s", config.SourceDirectory)
	f, err := os.Open(config.SourceDirectory)
	if err != nil {
		errx.Fatalf(err, "Could not open SourceDirectory: "+config.SourceDirectory)
	}
	files, err := f.Readdir(-1)
	if err != nil {
		errx.Fatalf(err, "Could not read SourceDirectory: "+config.SourceDirectory)
	}
	if err = f.Close(); err != nil {
		errx.Fatalf(err, "Could not close SourceDirectory: "+config.SourceDirectory)
	}
	log.Trace("files in SourceDirectory:")
	for _, file := range files {
		log.Trace(file.Name())
	}
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
