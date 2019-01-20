package platform

import (
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sparks"
	"github.com/michaKFromParis/sparks/sys"
)

// Ios represents the iOS platform
type Ios struct {
	enabled         bool
	SigningIdentity string
}

// Name is the lowercase name of the platform
func (i *Ios) Name() string {
	return "ios"
}

// Title is name of the platform
func (i *Ios) Title() string {
	return "iOS"
}

// Opt is the short command line option of the platform
func (i *Ios) Opt() string {
	return "i"
}

func (i *Ios) String() string {
	return i.Title()
}

// Enabled returns true if the platform is enabled
func (i *Ios) Enabled() bool {
	return i.enabled
}

// SetEnabled allows to enable / disable the platform
func (i *Ios) SetEnabled(enabled bool) {
	i.enabled = enabled
}

// Get installs the platform dependencies
func (i *Ios) Get() error {
	log.Info("Installing dependencies for " + i.Title())
	return nil
}

// Clean cleans the platform build
func (i *Ios) Clean() error {
	return nil
}

// Build builds the platform
func (i *Ios) Build(configuration sparks.Configuration) error {
	projectDirectory := filepath.Join(config.OutputDirectory, "projects", i.Title()+"-"+configuration.Title())
	i.prebuild()
	i.generate(configuration, projectDirectory)
	i.compile(configuration, projectDirectory)
	i.postbuild()

	return nil
}

func (i *Ios) prebuild() {
	xcode := sparks.XCode{}
	var signing sparks.SigningType
	signing = sparks.IPhoneDeveloper
	log.Debugf("selecting an %s profile for signing", signing.String())
	xcode.DetectSigning()
	identity, err := xcode.SelectSigning(sparks.IPhoneDeveloper)
	if err != nil {
		log.Warnf("could not select a %s signing identity.", signing)
	}
	i.SigningIdentity = identity
	log.Debugf("signing identity: %s", identity)
}

func (i *Ios) generate(configuration sparks.Configuration, projectDirectory string) {
	log.Info("sparks project generate --ios")

	iosSysRoot, err := sys.ExecuteEx("xcodebuild", "", true, "-sdk", config.SparksiOSSDK, "-version", "Path")
	if err != nil {
		errx.Fatalf(err, "could not determine ios sysroot")
	}
	iosSysRoot = strings.TrimSpace(iosSysRoot)
	log.Tracef("ios sysroot: %s", iosSysRoot)

	cmakeToolchainFile := filepath.Join(config.SDKDirectory, "scripts", "CMake", "toolchains", "iOS.cmake")
	libraryPath := filepath.Join(config.OutputDirectory, "lib", i.Title()+"-"+configuration.Title())

	cmake := sparks.NewCMake(i, configuration)

	cmake.AddArg("-GXcode")
	cmake.AddDefine("OS_IOS", "1")
	cmake.AddDefine("XCODE_SIGNING_IDENTITY", i.SigningIdentity)
	cmake.AddDefine("CMAKE_TOOLCHAIN_FILE", cmakeToolchainFile)
	cmake.AddDefine("XCODE_PROVISIONING_PROFILE_UUID", config.ProvisioningProfileUUID)
	cmake.AddDefine("PRODUCT_BUNDLE_IDENTIFIER", config.BundleIdentifier)
	cmake.AddDefine("CMAKE_OSX_SYSROOT", iosSysRoot)
	cmake.AddDefine("CMAKE_IOS_SYSROOT", iosSysRoot)

	// calling cmake once for the iphone
	platform := "iphoneos"
	iphoneProjectPath := filepath.Join(projectDirectory, platform)
	iphoneLibraryPath := filepath.Join(libraryPath, platform)
	var commonArgs = cmake.Args()
	cmake.AddDefine("CMAKE_ARCHIVE_OUTPUT_DIRECTORY", iphoneLibraryPath)
	cmake.AddDefine("CMAKE_LIBRARY_OUTPUT_DIRECTORY", iphoneLibraryPath)
	out, err := cmake.Run(iphoneProjectPath)
	if err != nil {
		errx.Fatalf(nil, "sparks project generate failed")
	}
	log.Trace("cmake output" + out)

	// and once for the simulator
	platform = "iphonesimulator"
	iphoneSimulatorProjectPath := filepath.Join(projectDirectory, platform)
	iphoneSimulatorLibraryPath := filepath.Join(libraryPath, platform)
	cmake.SetArgs(commonArgs)
	cmake.AddDefine("CMAKE_ARCHIVE_OUTPUT_DIRECTORY", iphoneSimulatorLibraryPath)
	cmake.AddDefine("CMAKE_LIBRARY_OUTPUT_DIRECTORY", iphoneSimulatorLibraryPath)
	out, err = cmake.Run(iphoneSimulatorProjectPath)
	if err != nil {
		errx.Fatalf(nil, "sparks project generate failed")
	}
	log.Trace("cmake output" + out)
}

func (i *Ios) compile(configuration sparks.Configuration, projectDirectory string) {
	log.Info("sparks project compile --ios")
	xcode := sparks.NewXCode(i, configuration)
	// if [ $buildConfiguration = $debugConfiguration ]; then
	//   LogWarning "Only building armv7 arch in $debugConfiguration configuration"
	//   arch='-arch armv7'
	// else
	//   arch="-arch armv7 -arch armv7s -arch arm64"
	// fi

	archsIos := []string{"-arch", "armv7", "-arch", "armv7s", "-arch", "arm64"}
	archsSimulator := []string{ /* "-arch", "i386", */ "-arch", "x86_64"}
	archsDebug := []string{"-arch", "armv7"}
	if configuration.Name() == "debug" { // build only one arch to speed local dev builds
		archsIos = archsDebug
	}
	iphoneOs := "iphoneos"
	simulator := "iphonesimulator"
	err := xcode.Build(filepath.Join(projectDirectory, iphoneOs), archsIos...)
	if err != nil {
		errx.Fatalf(err, "sparks project compile failed for "+iphoneOs)
	}
	err = xcode.Build(filepath.Join(projectDirectory, simulator), archsSimulator...)
	if err != nil {
		errx.Fatalf(err, "sparks project compile failed for "+simulator)
	}
}

func (i *Ios) postbuild() {
}
