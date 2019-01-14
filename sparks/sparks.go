package sparks

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/joomcode/errorx"
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

func Load() error {
	log.Info("sparks load")
	return CurrentProduct.Load()
}

func Save() {
	log.Info("sparks save")
	CurrentProduct.Save()
}

func Build(sourceDirectory string, outputDirectory string) error {
	log.Info("sparks build " + sourceDirectory)
	checkParameters(sourceDirectory, outputDirectory)
	if err := Load(); err != nil {
		return errorx.Decorate(err, "could not load sparks project at %s", sourceDirectory)
	}
	log.Tracef("loaded product:%s%+v", utils.NewLine, CurrentProduct)
	createBuildDirectoryStructure()
	sparksSourceDirectory := filepath.Join(config.SDKDirectory, "src", config.SDKName)
	sparksPlayerSourceDirectory := filepath.Join(config.SDKDirectory, "src", config.PlayerName)
	generateLuaBindings(sparksSourceDirectory, config.SDKName)
	// TODO fix math constants utils.Sed(filename, regex, newContent)
	generateLuaBindings(sparksSourceDirectory, "SparksNetworksLua")
	// TODO the line below probably should stay like this to build other c++ projects
	// generateLuaBindings(sparksPlayerSourceDirectory, config.ProductName)
	generateLuaBindings(sparksPlayerSourceDirectory, config.PlayerName)
	generateIcons(filepath.Join(config.SDKDirectory, "Assets", "Icon"))
	generateIcons(filepath.Join(config.SDKDirectory, "Assets", "SparksPlayerIcon"))
	generateSplash(filepath.Join(config.SDKDirectory, "Assets", "Splash"))
	// iterating through all enabled platforms in all enabled configurations
	// to call Platform.Build
	for _, platformName := range PlatformNames {
		platform := Platforms[platformName]
		if platform != nil && platform.Enabled() {
			for _, configurationName := range ConfigurationNames {
				configuration := Configurations[configurationName]
				if configuration != nil && configuration.Enabled() {
					log.Infof("sparks build --%s --%s --name %s", platform.Name(), configuration.Name(), config.ProductName)
					if err := platform.Build(configuration); err != nil {
						return errorx.Decorate(err, "sparks build failed for %s-%s", platform.Title(), configuration.Title())
					}
					// TODO calculate build time and build size
					log.Infof("build completed successfully in %d seconds", 42)
					log.Infof("build size: %d Mb", 42)
				}
			}
		}
	}
	return nil
}

func checkParameters(sourceDirectory string, outputDirectory string) { // TODO Check output here
	file, err := os.Stat(sourceDirectory)
	if err != nil {
		errx.Fatalf(err, "could not stat source directory: "+sourceDirectory)
	}
	if !file.IsDir() {
		errx.Fatalf(err, "source directory is not a directory: "+sourceDirectory)
	}
	config.SourceDirectory = sourceDirectory
	config.OutputDirectory = outputDirectory
	config.SDKDirectory = sourceDirectory
	log.Debugf("SDK Directory: %s", config.SDKDirectory)
	log.Debugf("Source Directory: %s", config.SourceDirectory)
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

func generateLuaBindings(sourceDirectory string, packageName string) {

	log.Info("sparks lua bind " + packageName)
	toluapp := getToluaPath()
	toluaHooksPath := filepath.Join(config.SDKDirectory, "src", "Sparks", "tolua.hooks.lua")
	dofileWithCorrectPath := fmt.Sprintf("dofile(\"%s\")", toluaHooksPath)
	reflectionFile := filepath.Join(sourceDirectory, packageName+".Reflection.lua")
	utils.SedFile(reflectionFile, "dofile\\(.*\\)", dofileWithCorrectPath)
	packagePath := filepath.Join(sourceDirectory, packageName)
	output, err := utils.ExecuteEx(
		toluapp,
		sourceDirectory,
		true,
		"-L", packagePath+".Reflection.lua",
		"-n", packageName,
		"-o", packagePath+".tolua.cpp",
		"-H", packagePath+".tolua.h",
		packagePath+".pkg")
	if err != nil {
		errx.Fatalf(err, "tolua++ execution failed: "+output)
	}
}

func getToluaPath() string {
	toluapp := filepath.Join(config.SDKDirectory, "dependencies", "toluapp", "bin")
	os, _ := utils.GetOs()
	switch os {
	case utils.Osx:
		toluapp = filepath.Join(toluapp, "toluapp_osx")
	case utils.Linux:
		toluapp = filepath.Join(config.SDKDirectory, "scripts", "bin", "tolua++")
	case utils.Windows:
		toluapp = filepath.Join(toluapp, "toluapp_script.exe")
	}
	return toluapp
}
func generateIcons(directory string) {
	log.Info("sparks icon generate")
}

func generateSplash(directory string) {
	log.Info("sparks splash generate")
}
