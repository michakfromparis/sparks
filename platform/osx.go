package platform

import (
	"fmt"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	version "github.com/hashicorp/go-version"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sparks"
	"github.com/michaKFromParis/sparks/utils"
)

type Osx struct {
	enabled bool
}

func (o *Osx) Name() string {
	return "osx"
}

func (o *Osx) Title() string {
	return "OSX"
}

func (o *Osx) Opt() string {
	return "o"
}

func (o *Osx) String() string {
	return o.Title()
}

func (o *Osx) Enabled() bool {
	return o.enabled
}

func (o *Osx) SetEnabled(enabled bool) {
	o.enabled = enabled
}

func (o *Osx) Deps() error {
	log.Info("Installing dependencies for " + o.Title())
	return nil
}
func (o *Osx) Clean() error {
	return nil
}
func (o *Osx) Build(configuration sparks.Configuration) error {
	o.prebuild()
	o.generate(configuration)
	o.compile()
	o.postbuild()
	return nil
}

func (o *Osx) prebuild() {
	o.guessXCodeSigningIdentity()
}

func (o *Osx) guessXCodeSigningIdentity() {

}

func (o *Osx) generate(configuration sparks.Configuration) {
	log.Info("sparks project generate --osx")
	log.Trace("determining osx sysroot")
	osxSysRoot, err := utils.ExecuteEx("xcodebuild", "", true, "-sdk", "macosx", "-version", "Path")
	if err != nil {
		errx.Fatalf(err, "could not determine osx sysroot")
	}
	osxSysRoot = strings.TrimSpace(osxSysRoot)
	log.Tracef("osx sysroot: %s", osxSysRoot)

	params := generateCmakeCommon(o, configuration)
	params += generateCmakeXcodeCommon()
	params += fmt.Sprintf("-DOS_OSX=1 ")
	params += fmt.Sprintf("-DCMAKE_OSX_ARCHITECTURES=\"%s\" ", config.SparksOSXArchitecture)
	params += fmt.Sprintf("-DCMAKE_OSX_DEPLOYMENT_TARGET=\"%s\" ", config.SparksOSXDeploymentTarget)
	params += fmt.Sprintf("-DCMAKE_OSX_SYSROOT=\"%s\" ", osxSysRoot)

	projectDirectory := filepath.Join(config.OutputDirectory, "projects", o.Title()+"-"+configuration.Title())

	utils.MkDir(projectDirectory)
	runCmake(projectDirectory, params)
}

func runCmake(directory string, params string) (string, error) {
	log.Debugf("running cmake in directory %s with these parameters: %s", directory, params)

	utils.MkDir(directory)
	cmakelistsPath := filepath.Join(config.SDKDirectory, "scripts", "CMake", "Sparks")
	params += fmt.Sprintf("\"%s\"", cmakelistsPath)
	var output string
	if output, err := utils.ExecuteEx("pwd", directory, true); err != nil {
		errx.Fatalf(err, "cmake execution failed: "+output)
	}
	log.Debug("cmake generation complete")
	return output, nil
}

func generateCmakeCommon(platform sparks.Platform, configuration sparks.Configuration) string {

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
	params += fmt.Sprintf("-DCMAKE_BUILD_TYPE=%s ", configuration.Title())
	if configuration.Name() == "shipping" {
		params += fmt.Sprintf("-DSHIPPING=ON ")
	}
	v, err := version.NewVersion(sparks.CurrentProduct.Version)
	if err != nil {
		errx.Fatalf(nil, "Failed to parse version: "+sparks.CurrentProduct.Version)
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
	widthAndHeight := strings.Split(sparks.CurrentProduct.View.Resolution, "x")
	if len(widthAndHeight) != 2 {
		errx.Fatalf(nil, "Failed to parse product resolution: "+sparks.CurrentProduct.View.Resolution)
	}
	width, height := widthAndHeight[0], widthAndHeight[1]
	params += fmt.Sprintf("-DPRODUCT_WIDTH=%s -DPRODUCT_HEIGHT=%s ", width, height)
	if sparks.CurrentProduct.View.Fullscreen == "yes" {
		params += fmt.Sprintf("--DPRODUCT_FULLSCREEN=ON ")
	} else {
		params += fmt.Sprintf("--DPRODUCT_FULLSCREEN=OFF ")
	}
	params += fmt.Sprintf("--DPRODUCT_DEFAULT_ORIENTATION=\"%s\" -DPRODUCT_SUPPORTED_ORIENTATIONS=\"%s\" ", sparks.CurrentProduct.View.DefaultOrientation, sparks.CurrentProduct.View.SupportedOrientations[0]) // TODO Proper listing of orientations
	params += fmt.Sprintf("-DCMAKE_INSTALL_PREFIX=\"%s\" ", filepath.Join(config.OutputDirectory, "lib", platform.Title()+"-"+configuration.Title()))
	log.Trace(params)
	return params
}

func generateCmakeXcodeCommon() string {
	return fmt.Sprintf("-GXcode -DXCODE_SIGNING_IDENTITY=\"%s\" ", config.XCodeSigningIdentity)
}

func (o *Osx) compile() {

}

func (o *Osx) postbuild() {

}
