package sparks

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/utils"
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
func createBuildDirectoryStructure() {
	log.Trace("creating build/bin, build/lib, build/projects")
	var binPath = filepath.Join(config.OutputDirectory, "bin")
	var libPath = filepath.Join(config.OutputDirectory, "lib")
	var projectsPath = filepath.Join(config.OutputDirectory, "projects")
	if err := os.MkdirAll(binPath, os.ModePerm); err != nil {
		errx.Error("failed to create build bin directory: " + binPath)
	}
	if err := os.MkdirAll(libPath, os.ModePerm); err != nil {
		errx.Error("failed to create build lib directory: " + libPath)
	}
	if err := os.MkdirAll(projectsPath, os.ModePerm); err != nil {
		errx.Error("failed to create build projects directory: " + projectsPath)
	}
}

func generateLuaBindings() {

	toluapp := filepath.Join(config.SDKDirectory, "dependencies", "toluapp", "bin")
	os, _ := utils.GetOs()
	switch os {
	case utils.Osx:
		toluapp = filepath.Join(toluapp, "toluapp_osx")
	case utils.Linux:
		toluapp = filepath.Join(toluapp, "toluapp_osx")
	case utils.Windows:
		toluapp = filepath.Join(toluapp, "tolua++")
	}
	utils.Execute(toluapp)
}

func Build() {
	log.Info("sparks build")
	Load()
	createBuildDirectoryStructure()
	generateLuaBindings()
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
