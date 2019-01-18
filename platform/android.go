package platform

import (
	"fmt"
	"path/filepath"

	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sparks"
	"github.com/michaKFromParis/sparks/sys"

	log "github.com/Sirupsen/logrus"
)

type Android struct {
	enabled bool
}

func (a *Android) Name() string {
	return "android"
}

func (a *Android) Title() string {
	return "Android"
}

func (a *Android) Opt() string {
	return "a"
}

func (a *Android) String() string {
	return a.Title()
}

func (a *Android) Enabled() bool {
	return a.enabled
}

func (a *Android) SetEnabled(enabled bool) {
	a.enabled = enabled
}

func (a *Android) Get() error {
	log.Info("Installing dependencies for " + a.Title())
	return nil
}
func (a *Android) Clean() error {
	return nil
}
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
	cmakeToolchainFile := filepath.Join(config.SDKDirectory, "scripts", "CMake", "toolchains", "Android.cmake")

	// params := generateCmakeCommon(a, configuration)
	params := fmt.Sprintf("-DOS_ANDROID=1 ")
	params += fmt.Sprintf("\"-DCMAKE_TOOLCHAIN_FILE%s\" ", cmakeToolchainFile)
	params += fmt.Sprintf("\"-GEclipse CDT4 - Unix Makefiles\" ")
	params += fmt.Sprintf("-DNDK_CCACHE=\"%s\" ", ccachePath)
	params += fmt.Sprintf("-DANDROID_NDK_RELEASE=\"%s\" ", config.SpakrsAndroidNDKVersion)
	params += fmt.Sprintf("-DANDROID_NATIVE_API_LEVEL=android-\"%s\" ", config.SpakrsAndroidApiLevel)
	// "-DLIBRARY_OUTPUT_PATH_ROOT=${buildRoot}/lib/${platformName}-${buildConfiguration}"
}

func (a *Android) compile() {

}

func (a *Android) postbuild() {

}
