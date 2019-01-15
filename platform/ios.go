package platform

import (
	"fmt"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sparks"
	"github.com/michaKFromParis/sparks/utils"
)

type Ios struct {
	enabled bool
}

func (i *Ios) Name() string {
	return "ios"
}

func (i *Ios) Title() string {
	return "iOS"
}

func (i *Ios) Opt() string {
	return "i"
}

func (i *Ios) String() string {
	return i.Title()
}

func (i *Ios) Enabled() bool {
	return i.enabled
}

func (i *Ios) SetEnabled(enabled bool) {
	i.enabled = enabled
}

func (i *Ios) Deps() error {
	return nil
}
func (i *Ios) Clean() error {
	return nil
}
func (i *Ios) Build(configuration sparks.Configuration) error {
	i.prebuild()
	i.generate(configuration)
	i.compile()
	i.postbuild()
	return nil
}

func (i *Ios) prebuild() {

}
func (i *Ios) generate(configuration sparks.Configuration) {
	log.Info("sparks project generate --ios")

	iosSysRoot, err := utils.ExecuteEx("xcodebuild", "", true, "-sdk", config.SparksiOSSDK, "-version", "Path")
	if err != nil {
		errx.Fatalf(err, "could not determine ios sysroot")
	}
	log.Tracef("ios sysroot: %s", iosSysRoot)

	cmakeToolchainFile := filepath.Join(config.SDKDirectory, "scripts", "CMake", "toolchains", "iOS.cmake")
	libraryPath := filepath.Join(config.OutputDirectory, "lib", i.Title()+"-"+configuration.Title())
	params := generateCmakeCommon(i, configuration)
	params += generateCmakeXcodeCommon()
	params += fmt.Sprintf("-DOS_IOS=1 ")
	params += fmt.Sprintf("\"-DCMAKE_TOOLCHAIN_FILE%s\" ", cmakeToolchainFile)
	params += fmt.Sprintf("-DXCODE_PROVISIONING_PROFILE_UUID=\"%s\" ", config.ProvisioningProfileUUID)
	params += fmt.Sprintf("-DPRODUCT_BUNDLE_IDENTIFIER=\"%s\" ", config.BundleIdentifier)
	params += fmt.Sprintf("-DCMAKE_OSX_SYSROOT=\"%s\" ", iosSysRoot)
	params += fmt.Sprintf("-DCMAKE_IOS_SYSROOT=\"%s\" ", iosSysRoot)

	projectsPath := filepath.Join(config.OutputDirectory, "projects", i.Title()+"-"+configuration.Title())

	// calling cmake once for the iphone
	platform := "iphoneos"
	iphoneProjectPath := filepath.Join(projectsPath, platform)
	iphoneLibraryPath := filepath.Join(libraryPath, platform)
	paramsiPhone := params + fmt.Sprintf("-DCMAKE_ARCHIVE_OUTPUT_DIRECTORY=\"%s\" -DCMAKE_LIBRARY_OUTPUT_DIRECTORY=\"%s\" ", iphoneLibraryPath, iphoneLibraryPath)
	if output, err := runCmake(iphoneProjectPath, paramsiPhone); err != nil {
		errx.Fatalf(err, fmt.Sprintf("cmake execution failed for %s: %s", platform, output))
	}

	// and once for the simulator
	platform = "iphonesimulator"
	iphoneSimulatorProjectPath := filepath.Join(projectsPath, platform)
	iphoneSimulatorLibraryPath := filepath.Join(libraryPath, platform)
	paramsiPhoneSimulator := params + fmt.Sprintf("-DCMAKE_ARCHIVE_OUTPUT_DIRECTORY=\"%s\" -DCMAKE_LIBRARY_OUTPUT_DIRECTORY=\"%s\" ", iphoneSimulatorLibraryPath, iphoneSimulatorLibraryPath)
	if output, err := runCmake(iphoneSimulatorProjectPath, paramsiPhoneSimulator); err != nil {
		errx.Fatalf(err, fmt.Sprintf("cmake execution failed for %s: %s", platform, output))
	}
}

func (i *Ios) compile() {

}
func (i *Ios) postbuild() {

}
