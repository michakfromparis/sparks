package sparks

import (
	"fmt"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	version "github.com/hashicorp/go-version"
	"github.com/joomcode/errorx"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/utils"
)

type CMake struct {
	command         string
	arguments       []string
	cmakelistsPath  string
	outputDirectory string
	platform        Platform
	configuration   Configuration
}

func NewCMake(platform Platform, configuration Configuration) *CMake {
	cm := new(CMake)
	cm.command = "cmake"
	cm.platform = platform
	cm.configuration = configuration
	cm.generateParams()
	return cm
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

func (cm *CMake) generateParams() {
	var args []string
	args = append(args, fmt.Sprintf("-DSPARKS_ROOT=%s", config.SDKDirectory))
	args = append(args, fmt.Sprintf("-DBUILD_ROOT=%s", config.OutputDirectory))
	args = append(args, fmt.Sprintf("-DPRODUCT_ROOT=%s", config.SourceDirectory))
	args = append(args, fmt.Sprintf("-DPRODUCT_NAME=%s", config.ProductName))
	if config.IncludeSparksSource == true {
		args = append(args, fmt.Sprintf("-DINCLUDE_SPARKS_SOURCE=ON"))
	} else {
		args = append(args, fmt.Sprintf("-DINCLUDE_SPARKS_SOURCE=OFF"))
	}
	// if config.VeryVerbose == true {
	// 	args = append(args, fmt.Sprintf("-DCMAKE_VERBOSE_MAKEFILE=ON --debug-output --trace"))
	// } else if config.Verbose {
	// 	args = append(args, fmt.Sprintf("--debug-output --trace"))
	// }
	args = append(args, fmt.Sprintf("-DCMAKE_BUILD_TYPE=%s", cm.configuration.Title()))
	if cm.configuration.Name() == "shipping" {
		args = append(args, fmt.Sprintf("-DSHIPPING=ON"))
	}
	major, minor, patch, build := parseVersion() // TODO move this out of here
	widthAndHeight := strings.Split(CurrentProduct.View.Resolution, "x")
	if len(widthAndHeight) != 2 {
		errx.Fatalf(nil, "Failed to parse product resolution: "+CurrentProduct.View.Resolution)
	}

	log.Infof("sparks project version %d.%d.%d.%d", major, minor, patch, build)
	args = append(args, fmt.Sprintf("-DSPARKS_VERSION_MAJOR=%d", major))
	args = append(args, fmt.Sprintf("-DSPARKS_VERSION_MINOR=%d", minor))
	args = append(args, fmt.Sprintf("-DSPARKS_VERSION_PATCH=%d", build))
	args = append(args, fmt.Sprintf("-DSPARKS_VERSION_BUILD=%d", patch))
	width, height := widthAndHeight[0], widthAndHeight[1]
	args = append(args, fmt.Sprintf("-DPRODUCT_WIDTH=%s", width))
	args = append(args, fmt.Sprintf("-DPRODUCT_HEIGHT=%s", height))
	if CurrentProduct.View.Fullscreen == "yes" {
		args = append(args, fmt.Sprintf("-DPRODUCT_FULLSCREEN=ON"))
	} else {
		args = append(args, fmt.Sprintf("-DPRODUCT_FULLSCREEN=OFF"))
	}
	args = append(args, fmt.Sprintf("-DPRODUCT_DEFAULT_ORIENTATION=\"%s\"", CurrentProduct.View.DefaultOrientation))          // TODO Proper listing of orientations
	args = append(args, fmt.Sprintf("-DPRODUCT_SUPPORTED_ORIENTATIONS=\"%s\"", CurrentProduct.View.SupportedOrientations[0])) // TODO Proper listing of orientations
	args = append(args, fmt.Sprintf("-DCMAKE_INSTALL_PREFIX=%s", filepath.Join(config.OutputDirectory, "lib", cm.platform.Title()+"-"+cm.configuration.Title())))
	cm.arguments = args
	log.Trace(args)
}

func (cm *CMake) Params() []string {
	return cm.arguments
}

func (cm *CMake) SetParams(params []string) {
	cm.arguments = params
}

func (cm *CMake) AddParam(params string) {
	cm.arguments = append(cm.arguments, params)
}

func (cm *CMake) Run(outputDirectory string) (string, error) {
	log.Debugf("running cmake in %s", outputDirectory)
	if err := utils.MkDir(outputDirectory); err != nil {
		return "could not create " + outputDirectory, err
	}
	cmakelistsPath := filepath.Join(config.SDKDirectory, "scripts", "CMake", "Sparks")
	cm.AddParam(cmakelistsPath)
	var output string
	// parameters := fmt.Sprintf("%s %s %s", cm.arguments, params, cmakelistsPath)
	// parts := strings.Split(parameters, " ")
	if output, err := utils.ExecuteEx(cm.command, outputDirectory, true, cm.arguments[0:]...); err != nil {
		return output, errorx.Decorate(err, "cmake execution failed")
	}
	return output, nil
}
