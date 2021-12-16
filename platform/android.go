package platform

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"leblox.com/sparks-cli/v2/conf"
	"leblox.com/sparks-cli/v2/errx"
	"leblox.com/sparks-cli/v2/sparks"
	"leblox.com/sparks-cli/v2/sys"
)

// Android represents the Android platform
type Android struct {
	enabled bool
}

// Name is the lowercase name of the platform
func (a *Android) Name() string {
	return "android"
}

// Title is name of the platform
func (a *Android) Title() string {
	return "Android"
}

// Opt is the short command line option of the platform
func (a *Android) Opt() string {
	return "a"
}

func (a *Android) String() string {
	return a.Title()
}

// Enabled returns true if the platform is enabled
func (a *Android) Enabled() bool {
	return a.enabled
}

// SetEnabled allows to enable / disable the platform
func (a *Android) SetEnabled(enabled bool) {
	a.enabled = enabled
}

// Get installs the platform dependencies
func (a *Android) Get() error {
	log.Info("Installing dependencies for " + a.Title())
	return nil
}

// Clean cleans the platform build
func (a *Android) Clean() error {
	return nil
}

// Code opens the code editor for the project
func (a *Android) Code(configuration sparks.Configuration) error {
	return nil
}

// Build builds the platform
func (a *Android) Build(configuration sparks.Configuration) error {
	a.prebuild()
	a.generate(configuration)
	a.compile()
	a.postbuild()
	return nil
}

func (a *Android) prebuild() {
}

func (a *Android) generate(configuration sparks.Configuration) {
	log.Info("sparks project generate --android")
	ccachePath := ""
	ccachePath, err := sys.Execute("which", "ccache")
	if err != nil {
		errx.Fatal(err)
	}
	cmakeToolchainFile := filepath.Join(conf.SDKDirectory, "scripts", "CMake", "toolchains", "Android.cmake")

	cmake := sparks.NewCMake(a, configuration)
	cmake.AddArg("-GEclipse CDT4 - Unix Makefiles")
	cmake.AddDefine("OS_ANDROID", "1")
	cmake.AddDefine("CMAKE_TOOLCHAIN_FILE", cmakeToolchainFile)
	cmake.AddDefine("NDK_CCACHE", ccachePath)
	cmake.AddDefine("ANDROID_NDK_RELEASE", conf.SpakrsAndroidNDKVersion)
	cmake.AddDefine("ANDROID_NATIVE_API_LEVEL", "android-"+string(conf.SparksAndroidAPILevel))
	// "-DLIBRARY_OUTPUT_PATH_ROOT=${buildRoot}/lib/${platformName}-${buildConfiguration}"
}

func (a *Android) compile() {

}

func (a *Android) postbuild() {

}
