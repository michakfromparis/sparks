package platform

import (
	"fmt"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
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
	projectDirectory := filepath.Join(config.OutputDirectory, "projects", o.Title()+"-"+configuration.Title())
	o.prebuild()
	o.generate(configuration, projectDirectory)
	o.compile(configuration, projectDirectory)
	o.postbuild()
	return nil
}

var SigningIdentity string

func (o *Osx) prebuild() {
	xcode := sparks.XCode{}
	xcode.DetectSigning()
	SigningIdentity = xcode.SigningIdentity(sparks.IPhoneDeveloper)
	if SigningIdentity == "" {
		errx.Fatalf(nil, "could not detect an xcode signing identity") // TODO explain how to obtain one
	}
}

func (o *Osx) generate(configuration sparks.Configuration, projectDirectory string) {
	log.Info("sparks project generate --osx")
	log.Trace("determining osx sysroot")
	osxSysRoot, err := utils.ExecuteEx("xcodebuild", "", true, "-sdk", "macosx", "-version", "Path")
	if err != nil {
		errx.Fatalf(err, "could not determine osx sysroot")
	}
	osxSysRoot = strings.TrimSpace(osxSysRoot)
	cc, err := utils.ExecuteEx("xcrun", "", true, "-find", "cc")
	if err != nil {
		errx.Fatalf(err, "could not determine C compiler")
	}
	cc = strings.TrimSpace(cc)
	cpp, err := utils.ExecuteEx("xcrun", "", true, "-find", "c++")
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
	cmake.AddArg(fmt.Sprintf("-DXCODE_SIGNING_IDENTITY=%s", SigningIdentity))
	cmake.AddArg(fmt.Sprintf("-DCMAKE_OSX_ARCHITECTURES=%s", config.SparksOSXArchitecture))
	cmake.AddArg(fmt.Sprintf("-DCMAKE_OSX_DEPLOYMENT_TARGET=%s", config.SparksOSXDeploymentTarget))

	out, err := cmake.Run(projectDirectory)
	log.Trace("cmake output" + out)
	if err != nil {
		errx.Fatalf(err, "sparks project generate failed")
	}
}

func (o *Osx) compile(configuration sparks.Configuration, projectDirectory string) {
	log.Info("sparks project compile --osx")
	xcode := sparks.NewXCode(o, configuration)
	xcode.Build(projectDirectory)
}

func (o *Osx) postbuild() {

}
