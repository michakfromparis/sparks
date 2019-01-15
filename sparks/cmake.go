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
	arguments       string
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

func (cm *CMake) generateParams() {

	params := fmt.Sprintf("-DSPARKS_ROOT=\"%s\" ", config.SDKDirectory)
	params += fmt.Sprintf("-DBUILD_ROOT=\"%s\" ", config.OutputDirectory)
	if config.IncludeSparksSource == true {
		params += fmt.Sprintf("-DINCLUDE_SPARKS_SOURCE=ON ")
	} else {
		params += fmt.Sprintf("-DINCLUDE_SPARKS_SOURCE=OFF ")
	}
	params += fmt.Sprintf("-DPRODUCT_ROOT=\"%s\" ", config.SourceDirectory)
	params += fmt.Sprintf("-DPRODUCT_NAME=\"%s\" ", config.ProductName)
	if config.VeryVerbose == true {
		params += fmt.Sprintf("-DCMAKE_VERBOSE_MAKEFILE=ON --debug-output --trace ")
	} else if config.Verbose {
		params += fmt.Sprintf("--debug-output --trace ")
	}
	params += fmt.Sprintf("-DCMAKE_BUILD_TYPE=%s ", cm.configuration.Title())
	if cm.configuration.Name() == "shipping" {
		params += fmt.Sprintf("-DSHIPPING=ON ")
	}
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
	log.Infof("sparks project version %d.%d.%d.%d", major, minor, patch, build)
	params += fmt.Sprintf("-DSPARKS_VERSION_MAJOR=\"%d\" -DSPARKS_VERSION_MINOR=\"%d\" -DSPARKS_VERSION_PATCH=\"%d\" -DSPARKS_VERSION_BUILD=\"%d\" ", major, minor, patch, build)
	widthAndHeight := strings.Split(CurrentProduct.View.Resolution, "x")
	if len(widthAndHeight) != 2 {
		errx.Fatalf(nil, "Failed to parse product resolution: "+CurrentProduct.View.Resolution)
	}
	width, height := widthAndHeight[0], widthAndHeight[1]
	params += fmt.Sprintf("-DPRODUCT_WIDTH=%s -DPRODUCT_HEIGHT=%s ", width, height)
	if CurrentProduct.View.Fullscreen == "yes" {
		params += fmt.Sprintf("--DPRODUCT_FULLSCREEN=ON ")
	} else {
		params += fmt.Sprintf("--DPRODUCT_FULLSCREEN=OFF ")
	}
	params += fmt.Sprintf("--DPRODUCT_DEFAULT_ORIENTATION=\"%s\" -DPRODUCT_SUPPORTED_ORIENTATIONS=\"%s\" ", CurrentProduct.View.DefaultOrientation, CurrentProduct.View.SupportedOrientations[0]) // TODO Proper listing of orientations
	params += fmt.Sprintf("-DCMAKE_INSTALL_PREFIX=\"%s\" ", filepath.Join(config.OutputDirectory, "lib", cm.platform.Title()+"-"+cm.configuration.Title()))
	cm.arguments = params
	log.Trace(params)
}

func (cm *CMake) Params() string {
	return cm.arguments
}

func (cm *CMake) SetParams(params string) {
	cm.arguments = params
}

func (cm *CMake) AddParams(params string) {
	cm.arguments += params
}

func (cm *CMake) Run(outputDirectory string) (string, error) {
	return cm.RunEx(outputDirectory, "")
}

func (cm *CMake) RunEx(outputDirectory string, params string) (string, error) {
	log.Debugf("running cmake in %s", outputDirectory)
	utils.MkDir(outputDirectory)
	cmakelistsPath := filepath.Join(config.SDKDirectory, "scripts", "CMake", "Sparks")
	cm.arguments += fmt.Sprintf("\"%s\"", cmakelistsPath)
	var output string
	if output, err := utils.ExecuteEx(cm.command, outputDirectory, true, cm.arguments+" "+params); err != nil {
		return output, errorx.Decorate(err, "cmake execution failed")
	}
	return output, nil
}
