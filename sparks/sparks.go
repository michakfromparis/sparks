package sparks

import (
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/joomcode/errorx"
	"github.com/michakfromparis/sparks/conf"
	"github.com/michakfromparis/sparks/errx"
	"github.com/michakfromparis/sparks/sys"
)

// CurrentProduct holds a reference to the currently loaded Sparks product from a .sparks file
var CurrentProduct = Product{}

// Init needs to be called once at the beggining of the program to Initialize Sparks
func Init() {
	log.Info("sparks init")
	if err := conf.Init(); err != nil {
		errx.Fatalf(err, "Configuration initialization failed")
	}
}

// Shutdown needs to be called once at the end of the program to Shutdown Sparks
func Shutdown() {
}

// Load loads a sparks product
func Load() error {
	log.Info("sparks load")
	return CurrentProduct.Load()
}

// Save saves a sparks product
func Save() {
	log.Info("sparks save")
	CurrentProduct.Save()
}

// Get is used to get packages and build dependencies
func Get() error {
	for _, platformName := range PlatformNames {
		platform := Platforms[platformName]
		if platform != nil && platform.Enabled() {
			if conf.GetDependencies {
				log.Info("sparks get --dependencies --" + platform.Name())
				err := platform.Get()
				if err != nil {
					return errorx.Decorate(err, "sparks get failed")
				}
			}
		}
	}
	return nil
}

// Code opens the default code editor for the defined & enabled platforms / configurations
func Code(sourceDirectory string, outputDirectory string) error {
	log.Info("sparks code " + sourceDirectory)
	checkParameters(sourceDirectory, outputDirectory)
	if err := Load(); err != nil {
		return errorx.Decorate(err, "could not load sparks project at %s", sourceDirectory)
	}
	log.Tracef("loaded product:%s%+v", sys.NewLine, CurrentProduct)
	createBuildDirectoryStructure()
	for _, platformName := range PlatformNames {
		platform := Platforms[platformName]
		if platform != nil && platform.Enabled() {
			for _, configurationName := range ConfigurationNames {
				configuration := Configurations[configurationName]
				if configuration != nil && configuration.Enabled() {
					log.Infof("sparks code --%s --%s --name %s", platform.Name(), configuration.Name(), conf.ProductName)
					if err := platform.Code(configuration); err != nil {
						return errorx.Decorate(err, "sparks code failed for %s-%s", platform.Title(), configuration.Title())
					}
				}
			}
		}
	}
	return nil
}

// Build builds a product in its defined & enabled platforms / configurations
func Build(sourceDirectory string, outputDirectory string) error {
	log.Info("sparks build " + sourceDirectory)
	checkParameters(sourceDirectory, outputDirectory)
	if err := Load(); err != nil {
		return errorx.Decorate(err, "could not load sparks project at %s", sourceDirectory)
	}
	log.Tracef("loaded product:%s%+v", sys.NewLine, CurrentProduct)
	createBuildDirectoryStructure()
	sparksSourceDirectory := filepath.Join(conf.SDKDirectory, "src", conf.SDKName)
	sparksPlayerSourceDirectory := filepath.Join(conf.SDKDirectory, "src", conf.PlayerName)
	if conf.GenerateLua {
		GenerateLuaBindings(sparksSourceDirectory, conf.SDKName)
		// TODO fix math constants sys.Sed(filename, regex, newContent)
		GenerateLuaBindings(sparksSourceDirectory, "SparksNetworksLua")
		// TODO the line below probably should stay like this to build other c++ projects
		// GenerateLuaBindings(sparksPlayerSourceDirectory, conf.ProductName)
		GenerateLuaBindings(sparksPlayerSourceDirectory, conf.PlayerName)
	}
	generateIcons(filepath.Join(conf.SDKDirectory, "Assets", "Icon"))
	generateIcons(filepath.Join(conf.SDKDirectory, "Assets", "SparksPlayerIcon"))
	generateSplash(filepath.Join(conf.SDKDirectory, "Assets", "Splash"))
	// iterating through all enabled platforms in all enabled configurations
	// to call Platform.Build
	for _, platformName := range PlatformNames {
		platform := Platforms[platformName]
		if platform != nil && platform.Enabled() {
			for _, configurationName := range ConfigurationNames {
				configuration := Configurations[configurationName]
				if configuration != nil && configuration.Enabled() {
					start := time.Now()
					log.Infof("sparks build --%s --%s --name %s", platform.Name(), configuration.Name(), conf.ProductName)

					if err := platform.Build(configuration); err != nil {
						return errorx.Decorate(err, "sparks build failed for %s-%s", platform.Title(), configuration.Title())
					}
					buildPath := filepath.Join(outputDirectory, "bin", platform.Title()+"-"+configuration.Title())
					libPath := filepath.Join(outputDirectory, "lib", platform.Title()+"-"+configuration.Title())
					_, err := os.Stat(buildPath)
					if err != nil {
						if os.IsNotExist(err) {
							if err := os.MkdirAll(buildPath, os.ModePerm); err != nil {
								return errorx.Decorate(err, "Could not create build directory: "+buildPath)
							}
						} else {
							return errorx.Decorate(err, "build directory does not exist: "+buildPath)
						}
					}
					buildSize, err := sys.DirSize(buildPath)
					if err != nil {
						log.Error("Could not determine build size: " + err.Error())
					}
					libSize, err := sys.DirSize(libPath)
					if err != nil {
						log.Error("Could not determine libraries size: " + err.Error())
					}
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
	conf.SourceDirectory = sourceDirectory
	conf.OutputDirectory = outputDirectory
	conf.SDKDirectory = sourceDirectory
	log.Debugf("SDK Directory: %s", conf.SDKDirectory)
	log.Debugf("Source Directory: %s", conf.SourceDirectory)
}

func createBuildDirectoryStructure() {
	log.Trace("creating build/bin, build/lib, build/projects")
	var binPath = filepath.Join(conf.OutputDirectory, "bin")
	var libPath = filepath.Join(conf.OutputDirectory, "lib")
	var projectsPath = filepath.Join(conf.OutputDirectory, "projects")
	if err := os.MkdirAll(binPath, os.ModePerm); err != nil {
		log.Warn("failed to create build bin directory: " + binPath)
	}
	if err := os.MkdirAll(libPath, os.ModePerm); err != nil {
		log.Warn("failed to create build lib directory: " + libPath)
	}
	if err := os.MkdirAll(projectsPath, os.ModePerm); err != nil {
		log.Warn("failed to create build projects directory: " + projectsPath)
	}
}

func generateIcons(directory string) {
	log.Info("sparks icon generate")
}

func generateSplash(directory string) {
	log.Info("sparks splash generate")
}
