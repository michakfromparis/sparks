package sparks

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	version "github.com/hashicorp/go-version"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sys"
)

// CMake is a wrapper class around cmake
type CMake struct {
	command         string
	arguments       []string
	cmakelistsPath  string
	outputDirectory string
	platform        Platform
	configuration   Configuration
}

// NewCMake returns an instance of a new CMake class
func NewCMake(platform Platform, configuration Configuration) *CMake {
	cm := new(CMake)
	cm.command = "cmake"
	cm.platform = platform
	cm.configuration = configuration
	cm.generateArgs()
	return cm
}

// Args returns the array of strings passed as parameters to cmake
func (cm *CMake) Args() []string {
	return cm.arguments
}

// SetArgs sets the array of strings passed as parameters to cmake
func (cm *CMake) SetArgs(params []string) {
	cm.arguments = params
}

// AddArg adds a string to the array of strings passed as parameters to cmake
func (cm *CMake) AddArg(arg string) {
	cm.arguments = append(cm.arguments, arg)
}

// AddDefine adds a string to the array of strings passed as parameters to cmake
func (cm *CMake) AddDefine(key string, value string) {
	if value == "" {
		log.Warnf("cmake define %s is passed an empty value", key)
	}
	cm.arguments = append(cm.arguments, fmt.Sprintf("-D%s=%s", key, value))
}

// Run needs to be called once the parameters list is built.
// Run will output project files in outputDirectory
func (cm *CMake) Run(outputDirectory string, args ...[]string) (string, error) {
	log.Debugf("running cmake in %s", outputDirectory)
	cm.outputDirectory = outputDirectory
	var err error
	if err = sys.MkDir(outputDirectory); err != nil {
		return "could not create " + outputDirectory, err
	}
	cm.cmakelistsPath = filepath.Join(config.SDKDirectory, "scripts", "CMake", "Sparks")
	cm.AddArg(cm.cmakelistsPath)
	return sys.ExecuteEx(cm.command, cm.outputDirectory, true, cm.arguments[0:]...)
}

func (cm *CMake) generateArgs() {
	cm.AddDefine("SPARKS_ROOT", config.SDKDirectory)
	cm.AddDefine("BUILD_ROOT", config.OutputDirectory)
	cm.AddDefine("PRODUCT_ROOT", config.SourceDirectory)
	cm.AddDefine("PRODUCT_NAME", config.ProductName)

	// if config.VeryVerbose == true {
	// 	cm.AddArg("CMAKE_VERBOSE_MAKEFILE=ON --debug-output --trace")))
	// } else if config.Verbose {
	// 	cm.AddArg(fmt.Sprintf("--debug-output --trace")))
	// }

	if config.IncludeSparksSource {
		cm.AddDefine("INCLUDE_SPARKS_SOURCE", "ON")
	} else {
		cm.AddDefine("INCLUDE_SPARKS_SOURCE", "`OFF")
	}
	cm.AddDefine("CMAKE_BUILD_TYPE", cm.configuration.Title())
	if cm.configuration.Name() == "shipping" {
		cm.AddDefine("SHIPPING", "ON")
	}
	if CurrentProduct.View.Fullscreen == "yes" {
		cm.AddDefine("PRODUCT_FULLSCREEN", "ON")
	} else {
		cm.AddDefine("PRODUCT_FULLSCREEN", "OFF")
	}
	widthAndHeight := strings.Split(CurrentProduct.View.Resolution, "x")
	if len(widthAndHeight) != 2 {
		errx.Fatalf(nil, "Failed to parse product resolution: "+CurrentProduct.View.Resolution)
	}
	width, height := widthAndHeight[0], widthAndHeight[1]
	major, minor, patch, build := parseVersion() // TODO move this out of here

	log.Infof("sparks project version %d.%d.%d.%d", major, minor, patch, build)
	cm.AddDefine("SPARKS_VERSION_MAJOR", strconv.Itoa(major))
	cm.AddDefine("SPARKS_VERSION_MINOR", strconv.Itoa(minor))
	cm.AddDefine("SPARKS_VERSION_PATCH", strconv.Itoa(build))
	cm.AddDefine("SPARKS_VERSION_BUILD", strconv.Itoa(patch))
	cm.AddDefine("PRODUCT_WIDTH", width)
	cm.AddDefine("PRODUCT_HEIGHT", height)
	cm.AddDefine("PRODUCT_DEFAULT_ORIENTATION", CurrentProduct.View.DefaultOrientation)          // TODO Proper listing of orientations
	cm.AddDefine("PRODUCT_SUPPORTED_ORIENTATIONS", CurrentProduct.View.SupportedOrientations[0]) // TODO Proper listing of orientations
	cm.AddDefine("CMAKE_INSTALL_PREFIX", filepath.Join(config.OutputDirectory, "lib", cm.platform.Title()+"-"+cm.configuration.Title()))
}

func parseVersion() (int, int, int, int) {
	v, err := version.NewVersion(CurrentProduct.Version)
	if err != nil {
		errx.Fatalf(nil, "Failed to parse version: "+CurrentProduct.Version)
	}
	var major, minor, patch, build int
	segments := v.Segments()
	segmentsCount := len(segments)
	if segmentsCount > 0 {
		major = segments[0]
	}
	if segmentsCount > 1 {
		minor = segments[1]
	}
	if segmentsCount > 2 {
		patch = segments[2]
	}
	if segmentsCount > 3 {
		build = segments[3]
	}
	return major, minor, patch, build
}
