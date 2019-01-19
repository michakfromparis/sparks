package platform

import (
	"fmt"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sparks"
	"github.com/michaKFromParis/sparks/sys"
)

// Osx represents the OSX platform
type Osx struct {
	enabled         bool
	SigningIdentity string
}

// Name is the lowercase name of the platform
func (o *Osx) Name() string {
	return "osx"
}

// Title is name of the platform
func (o *Osx) Title() string {
	return "OSX"
}

// Opt is the short command line option of the platform
func (o *Osx) Opt() string {
	return "o"
}

func (o *Osx) String() string {
	return o.Title()
}

// Enabled returns true if the platform is enabled
func (o *Osx) Enabled() bool {
	return o.enabled
}

// SetEnabled allows to enable / disable the platform
func (o *Osx) SetEnabled(enabled bool) {
	o.enabled = enabled
}

// Get installs the platform dependencies
func (o *Osx) Get() error {
	return nil
}

// Clean cleans the platform build
func (o *Osx) Clean() error {
	return nil
}

// Build builds the platform
func (o *Osx) Build(configuration sparks.Configuration) error {
	projectDirectory := filepath.Join(config.OutputDirectory, "projects", o.Title()+"-"+configuration.Title())
	o.prebuild()
	o.generate(configuration, projectDirectory)
	o.compile(configuration, projectDirectory)
	o.postbuild()
	return nil
}

func (o *Osx) prebuild() {
	xcode := sparks.XCode{}
	xcode.DetectSigning()
	o.SigningIdentity = xcode.SigningIdentity(sparks.MacDeveloper)
	if o.SigningIdentity == "" {
		errx.Fatalf(nil, "could not detect an xcode signing identity") // TODO explain how to obtain one
	}
}

func (o *Osx) generate(configuration sparks.Configuration, projectDirectory string) {
	log.Info("sparks project generate --osx")
	log.Trace("determining osx sysroot")
	osxSysRoot, err := sys.ExecuteEx("xcodebuild", "", true, "-sdk", "macosx", "-version", "Path")
	if err != nil {
		errx.Fatalf(err, "could not determine osx sysroot")
	}
	osxSysRoot = strings.TrimSpace(osxSysRoot)
	cc, err := sys.ExecuteEx("xcrun", "", true, "-find", "cc")
	if err != nil {
		errx.Fatalf(err, "could not determine C compiler")
	}
	cc = strings.TrimSpace(cc)
	cpp, err := sys.ExecuteEx("xcrun", "", true, "-find", "c++")
	if err != nil {
		errx.Fatalf(err, "could not determine C++ compiler")
	}
	cpp = strings.TrimSpace(cpp)
	log.Tracef("osx sysroot: %s", osxSysRoot)
	log.Tracef("C compiler: %s", cc)
	log.Tracef("C++ compiler: %s", cpp)

	cmake := sparks.NewCMake(o, configuration)
	cmake.AddArg(fmt.Sprintf("-DOS_OSX=1"))
	cmake.AddArg(fmt.Sprintf("-GXcode"))
	cmake.AddArg(fmt.Sprintf("-DCMAKE_OSX_SYSROOT=%s", osxSysRoot))
	cmake.AddArg(fmt.Sprintf("-DCMAKE_C_COMPILER=%s", cc))
	cmake.AddArg(fmt.Sprintf("-DCMAKE_CXX_COMPILER=%s", cpp))
	cmake.AddArg(fmt.Sprintf("-DXCODE_SIGNING_IDENTITY=%s", o.SigningIdentity))
	cmake.AddArg(fmt.Sprintf("-DCMAKE_OSX_ARCHITECTURES=%s", config.SparksOSXArchitecture))
	cmake.AddArg(fmt.Sprintf("-DCMAKE_OSX_DEPLOYMENT_TARGET=%s", config.SparksOSXDeploymentTarget))

	out, err := cmake.Run(projectDirectory)
	log.Trace("cmake output" + out)
	if err != nil {
		errx.Fatalf(nil, "sparks project generate failed")
	}
}

func (o *Osx) compile(configuration sparks.Configuration, projectDirectory string) {
	log.Info("sparks project compile --osx")
	xcode := sparks.NewXCode(o, configuration)
	err := xcode.Build(projectDirectory)
	if err != nil {
		errx.Fatalf(err, "sparks project compile failed")
	}

}

func (o *Osx) postbuild() {
}
