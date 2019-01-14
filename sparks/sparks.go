package sparks

import (
	"fmt"
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
	log.Tracef("%+v", CurrentProduct)
}

func Save() {
	log.Info("sparks save")
	CurrentProduct.Save()
}

func Build(sourceDirectory string, outputDirectory string) {
	log.Info("sparks build " + sourceDirectory)
	checkParameters(sourceDirectory, outputDirectory)
	Load()
	createBuildDirectoryStructure()
	sparksSourceDirectory := filepath.Join(config.SDKDirectory, "src", config.SDKName)
	generateLuaBindings(sparksSourceDirectory, config.SDKName)
	// utils.Sed(filename, regex, newContent)
	generateLuaBindings(sparksSourceDirectory, "SparksNetworksLua")
	generateLuaBindings(sparksSourceDirectory, config.ProductName)

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

	log.Info("generating lua bindings for " + packageName)
	toluapp := filepath.Join(config.SDKDirectory, "dependencies", "toluapp", "bin")

	os, _ := utils.GetOs()
	switch os {
	case utils.Osx:
		toluapp = filepath.Join(toluapp, "toluapp_osx")
	case utils.Linux:
		toluapp = filepath.Join(toluapp, "tolua++")
	case utils.Windows:
		toluapp = filepath.Join(toluapp, "toluapp_script.exe")
	}
	toluaHooksPath := filepath.Join(config.SDKDirectory, "src", "Sparks", "tolua.hooks.lua")
	dofileWithCorrectPath := fmt.Sprintf("dofile(\"%s\")%s", toluaHooksPath, utils.LineBreak)
	reflectionFile := filepath.Join(sourceDirectory, packageName+".Reflection.lua")
	utils.Sed(reflectionFile, "dofile\\(.*\\)", dofileWithCorrectPath)
	packagePath := filepath.Join(sourceDirectory, packageName)
	out, err := utils.Execute(
		toluapp, "-L", packagePath+".Reflection.lua",
		"-n", packageName,
		"-o", packagePath+".tolua.cpp",
		"-H", packagePath+".tolua.h",
		packagePath+".pkg")
	if err != nil {
		errx.Fatalf(err, "tolua++ execution failed: "+out)
	}
}
