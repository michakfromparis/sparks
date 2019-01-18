package sparks

import (
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/joomcode/errorx"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sys"
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

func Get() error {

	for _, platformName := range PlatformNames {
		platform := Platforms[platformName]
		if platform != nil && platform.Enabled() {
			return platform.Get()
		}
	}
	return nil
}

func Build(sourceDirectory string, outputDirectory string) error {
	log.Info("sparks build " + sourceDirectory)
	checkParameters(sourceDirectory, outputDirectory)
	if err := Load(); err != nil {
		return errorx.Decorate(err, "could not load sparks project at %s", sourceDirectory)
	}
	log.Tracef("loaded product:%s%+v", sys.NewLine, CurrentProduct)
	createBuildDirectoryStructure()
	sparksSourceDirectory := filepath.Join(config.SDKDirectory, "src", config.SDKName)
	sparksPlayerSourceDirectory := filepath.Join(config.SDKDirectory, "src", config.PlayerName)
	if config.GenerateLua {
		GenerateLuaBindings(sparksSourceDirectory, config.SDKName)
		// TODO fix math constants sys.Sed(filename, regex, newContent)
		GenerateLuaBindings(sparksSourceDirectory, "SparksNetworksLua")
		// TODO the line below probably should stay like this to build other c++ projects
		// GenerateLuaBindings(sparksPlayerSourceDirectory, config.ProductName)
		GenerateLuaBindings(sparksPlayerSourceDirectory, config.PlayerName)
	}
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
					start := time.Now()
					log.Infof("sparks build --%s --%s --name %s", platform.Name(), configuration.Name(), config.ProductName)

					if err := platform.Build(configuration); err != nil {
						return errorx.Decorate(err, "sparks build failed for %s-%s", platform.Title(), configuration.Title())
					}
					buildPath := filepath.Join(outputDirectory, "bin", platform.Title()+"-"+configuration.Title())
					libPath := filepath.Join(outputDirectory, "lib", platform.Title()+"-"+configuration.Title())
					stat, err := os.Stat(buildPath)
					if err != nil || !stat.IsDir() {
						return errorx.Decorate(err, "build directory does not exist: "+buildPath)
					}
					buildSize, _ := sys.DirSize(buildPath)
					libSize, _ := sys.DirSize(libPath)
					log.Infof("build completed successfully in %v", sys.FmtDuration(time.Since(start)))
					log.Infof("build size: %s (libraries: %s)", buildSize, libSize)
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

func generateIcons(directory string) {
	log.Info("sparks icon generate")
}

func generateSplash(directory string) {
	log.Info("sparks splash generate")
}
