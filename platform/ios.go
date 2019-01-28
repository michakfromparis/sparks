package platform

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/joomcode/errorx"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sparks"
	"github.com/michaKFromParis/sparks/sys"
)

// Ios represents the iOS platform
type Ios struct {
	enabled         bool
	signingIdentity *sparks.SigningIdentity
	configuration   sparks.Configuration
	provisioningID  string
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

// Code opens the code editor for the project
func (i *Ios) Code(configuration sparks.Configuration) error {
	i.configuration = configuration
	projectDirectory := filepath.Join(config.OutputDirectory, "projects", i.Title()+"-"+i.configuration.Title())
	projectFilename := filepath.Join(projectDirectory, "iphoneos", config.ProductName+".xcodeproj")
	i.prebuild()
	i.generate(projectDirectory)
	_, err := sys.Execute("open", projectFilename)
	if err != nil {
		return errorx.Decorate(err, "could not open xcode for project: "+projectFilename)
	}
	return nil
}

// Build builds the platform
func (i *Ios) Build(configuration sparks.Configuration) error {
	i.configuration = configuration
	projectDirectory := filepath.Join(config.OutputDirectory, "projects", i.Title()+"-"+configuration.Title())
	i.prebuild()
	i.generate(projectDirectory)
	i.compile(projectDirectory)
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
	i.signingIdentity = identity
	log.Debugf("signing identity: %s", identity)
	i.importProvisioningProfile()
}

func (i *Ios) importProvisioningProfile() {
	log.Debug("looking for provioning profile")
	keysDirectory := filepath.Join(config.SourceDirectory, "conf", "keys", "ios", "development")
	if i.configuration.Name() == "shipping" {
		keysDirectory = filepath.Join(config.SourceDirectory, "conf", "keys", "ios", "distribution")
	}
	provisioningFilename, err := sys.Execute("find", keysDirectory, "-name", "*.mobileprovision")
	provisioningFilename = strings.TrimSpace(provisioningFilename)
	stat, err2 := os.Stat(provisioningFilename)
	if err != nil || err2 != nil || stat.IsDir() {
		log.Warnf("could not find a provisioning profile inside the project at %s", keysDirectory)
		return
	}
	log.Debugf("provisioning profile found: %s", provisioningFilename)
	output, err := sys.Execute("grep", "UUID", "-A1", "-a", provisioningFilename)
	if err != nil {
		log.Error("Could not parse provisioning profile at " + provisioningFilename)
	}
	re := regexp.MustCompile("string>(.*)</string")
	matched := re.MatchString(output)
	if !matched {
		log.Warnf("Could not parse provisioning profile at " + provisioningFilename)
	}
	// basename := filepath.Base(provisioningFilename)
	// i.provisioningID = strings.TrimSuffix(basename, filepath.Ext(basename))
	i.provisioningID = re.FindStringSubmatch(output)[1]
	log.Debugf("Provisioning Profile Name: %s", i.provisioningID)
}

func (i *Ios) generate(projectDirectory string) {
	log.Info("sparks project generate --ios")

	iosSysRoot, err := sys.ExecuteEx("xcodebuild", "", true, "-sdk", config.SparksiOSSDK, "-version", "Path")
	if err != nil {
		errx.Fatalf(err, "could not determine ios sysroot")
	}
	iosSysRoot = strings.TrimSpace(iosSysRoot)
	log.Tracef("ios sysroot: %s", iosSysRoot)

	cmakeToolchainFile := filepath.Join(config.SDKDirectory, "scripts", "CMake", "toolchains", "iOS.cmake")
	libraryPath := filepath.Join(config.OutputDirectory, "lib", i.Title()+"-"+i.configuration.Title())

	cmake := sparks.NewCMake(i, i.configuration)

	var automaticSigning = false

	cmake.AddArg("-GXcode")
	cmake.AddDefine("OS_IOS", "1")
	cmake.AddDefine("IOS_PLATFORM", "OS")
	if automaticSigning {
	} else {
		cmake.AddDefine("XCODE_SIGNING_IDENTITY", i.signingIdentity.Name)
		cmake.AddDefine("XCODE_SIGNING_DEVELOPMENT_TEAM", i.signingIdentity.TeamID)
		cmake.AddDefine("XCODE_PROVISIONING_PROFILE_UUID", i.provisioningID)
	}
	cmake.AddDefine("CMAKE_TOOLCHAIN_FILE", cmakeToolchainFile)
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
	//log.Trace("cmake output: " + out)

	// and once for the simulator
	platform = "iphonesimulator"
	iphoneSimulatorProjectPath := filepath.Join(projectDirectory, platform)
	iphoneSimulatorLibraryPath := filepath.Join(libraryPath, platform)
	cmake.SetArgs(commonArgs)
	cmake.AddDefine("CMAKE_ARCHIVE_OUTPUT_DIRECTORY", iphoneSimulatorLibraryPath)
	cmake.AddDefine("CMAKE_LIBRARY_OUTPUT_DIRECTORY", iphoneSimulatorLibraryPath)
	out, err = cmake.Run(iphoneSimulatorProjectPath)
	if err != nil {
		errx.Fatalf(nil, "sparks project generate failed: "+out)
	}
	//log.Trace("cmake output: " + out)
}

func (i *Ios) compile(projectDirectory string) {
	log.Info("sparks project compile --ios")
	xcode := sparks.NewXCode(i, i.configuration)
	// if [ $buildConfiguration = $debugConfiguration ]; then
	//   LogWarning "Only building armv7 arch in $debugConfiguration configuration"
	//   arch='-arch armv7'
	// else
	//   arch="-arch armv7 -arch armv7s -arch arm64"
	// fi

	archsIos := []string{"-arch", "armv7", "-arch", "armv7s", "-arch", "arm64"}
	archsSimulator := []string{ /* "-arch", "i386", */ "-arch", "x86_64"}
	// archsDebug := []string{"-arch", "armv7"}
	if i.configuration.Name() == "debug" { // build only one arch to speed local dev builds
		// archsIos = archsDebug
	}
	iphoneOs := "iphoneos"
	simulator := "iphonesimulator"
	log.Infof("sparks project compile --ios %v", archsIos)
	err := xcode.Build(filepath.Join(projectDirectory, iphoneOs), archsIos...)
	if err != nil {
		errx.Fatalf(err, "sparks project compile failed for "+iphoneOs)
	}
	log.Infof("sparks project compile --ios %v", archsSimulator[:])
	err = xcode.Build(filepath.Join(projectDirectory, simulator), archsSimulator...)
	if err != nil {
		errx.Fatalf(err, "sparks project compile failed for "+simulator)
	}
}

func (i *Ios) postbuild() {
}
