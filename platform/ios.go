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
	iosSysRoot = strings.TrimSpace(iosSysRoot)
	log.Tracef("ios sysroot: %s", iosSysRoot)

	cmakeToolchainFile := filepath.Join(config.SDKDirectory, "scripts", "CMake", "toolchains", "iOS.cmake")
	libraryPath := filepath.Join(config.OutputDirectory, "lib", i.Title()+"-"+configuration.Title())

	cmake := sparks.NewCMake(i, configuration)

	cmake.AddArg(fmt.Sprintf("-DOS_IOS=1"))
	cmake.AddArg(fmt.Sprintf("-GXcode"))
	cmake.AddArg(fmt.Sprintf("-DXCODE_SIGNING_IDENTITY=\"%s\"", config.XCodeSigningIdentity))
	cmake.AddArg(fmt.Sprintf("-DCMAKE_TOOLCHAIN_FILE=%s", cmakeToolchainFile))
	cmake.AddArg(fmt.Sprintf("-DXCODE_PROVISIONING_PROFILE_UUID=%s", config.ProvisioningProfileUUID))
	cmake.AddArg(fmt.Sprintf("-DPRODUCT_BUNDLE_IDENTIFIER=%s", config.BundleIdentifier))
	cmake.AddArg(fmt.Sprintf("-DCMAKE_OSX_SYSROOT=%s", iosSysRoot))
	cmake.AddArg(fmt.Sprintf("-DCMAKE_IOS_SYSROOT=%s", iosSysRoot))

	// root directory of 2 different projects for iphoneos and iphonesimulator
	projectsPath := filepath.Join(config.OutputDirectory, "projects", i.Title()+"-"+configuration.Title())

	// calling cmake once for the iphone
	platform := "iphoneos"
	iphoneProjectPath := filepath.Join(projectsPath, platform)
	iphoneLibraryPath := filepath.Join(libraryPath, platform)
	var commonArgs = cmake.Args()
	cmake.AddArg(fmt.Sprintf("-DCMAKE_ARCHIVE_OUTPUT_DIRECTORY=\"%s\"", iphoneLibraryPath))
	cmake.AddArg(fmt.Sprintf("-DCMAKE_LIBRARY_OUTPUT_DIRECTORY=\"%s\"", iphoneLibraryPath))
	out, err := cmake.Run(iphoneProjectPath)
	if err != nil {
		errx.Fatalf(err, "sparks project generate failed")
	}
	log.Trace("cmake output" + out)

	// and once for the simulator
	platform = "iphonesimulator"
	iphoneSimulatorProjectPath := filepath.Join(projectsPath, platform)
	iphoneSimulatorLibraryPath := filepath.Join(libraryPath, platform)
	cmake.SetArgs(commonArgs)
	cmake.AddArg(fmt.Sprintf("-DCMAKE_ARCHIVE_OUTPUT_DIRECTORY=\"%s\"", iphoneSimulatorLibraryPath))
	cmake.AddArg(fmt.Sprintf("-DCMAKE_LIBRARY_OUTPUT_DIRECTORY=\"%s\"", iphoneSimulatorLibraryPath))
	out, err = cmake.Run(iphoneSimulatorProjectPath)
	if err != nil {
		errx.Fatalf(err, "sparks project generate failed")
	}
	log.Trace("cmake output" + out)
}

func (i *Ios) compile() {

}
func (i *Ios) postbuild() {

}
