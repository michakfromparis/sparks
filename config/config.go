package config

import (
	"fmt"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/utils"
)

// Get specific
var GetDependencies bool

var SourceDirectory string
var OutputDirectory string
var SDKDirectory string
var SDKName = "Sparks"
var PlayerName = "SparksPlayer"
var ProductName string
var Verbose = false
var VeryVerbose = false // TODO temporary for debugging. need to set a proper verbose level

var IncludeSparksSource = true
var GenerateLua = false

// OSX specific
var SparksOSXSDK = "macosx10.14" // run `xcodebuild -showsdks` to list valid values
var SparksOSXArchitecture = "x86_64"
var SparksOSXDeploymentTarget = "10.10"
var XCodeSigningIdentity = "iPhone Developer: Michel Courtine"
var XCodeRecommendedVersion = "10.1"

// iOS specific
var SparksiOSSDK = "iphoneos12.1"                 // run `xcodebuild -showsdks` to list valid values
var sparksiOSSimulatorSDK = "iphonesimulator12.1" // run `xcodebuild -showsdks` to list valid values
var DevelopmentSigningType = "iPhone Developer"
var DistributionSigningType = "iPhone Distribution"
var BundleIdentifier = "me.ivoltage.sparksplayer"
var ProvisioningProfileUUID = ""

// Android specific
var AndroidSDKRoot = filepath.Join("Android", "sdk")
var AndroidNDKRoot = filepath.Join("Android", "ndk")
var AndroidSDKToolsVersion = "25.2.3"
var SparksAndroidBuildToolsVersion = "25.2.3"
var SpakrsAndroidApiLevel = 19
var SpakrsAndroidNDKVersion = "r10d"
var SndroidPackagePrefix = "me.ivoltage"
var SparksAndroidNDKUrlOSX = "http://dl.google.com/android/ndk/android-ndk-${spakrsAndroidNDKVersion}-darwin-x86_64.bin"
var SparksAndroidNDKUrlLinux = "http://dl.google.com/android/ndk/android-ndk-${spakrsAndroidNDKVersion}-linux-x86_64.bin"
var FacebookAndroidApiLevel = 9
var VolleyAndroidApiLevel = 19

// Emscripten specific
var EmscriptenSDKRoot = "Emscripten"
var EmscriptenVersion = "incoming" // possible values '1.27.0', '1.29.0' 'latest', 'master' or 'incoming' for the latest version
var EmscriptenBrowser = "firefox"  // possible values: 'chrome' 'firefox' 'safari'

// Windows specific
// var WindowsCompiler="MSBuild"
var WindowsCompiler = "VisualStudio"

func Init() error {
	log.Info("sparks config init")
	if ProductName == "" {
		ProductName = PlayerName
	}
	var err error
	if SourceDirectory, err = utils.Pwd(); err != nil {
		errx.Fatalf(err, "could not determine current working directory")
	}
	SDKDirectory = SourceDirectory
	OutputDirectory = path.Join(SourceDirectory, "build")
	log.Debug(String())
	return nil
}

func String() string {
	// platforms := ""
	// for _, name := range platform.PlatformNames {
	// 	if Platforms[name] {
	// 		platforms += name + " "
	// 	}
	// }
	return fmt.Sprintf(`Loaded Configuration:
ProductName: %s
SourceDirectory: %s
OutputDirectory: %s
Verbose: %t`, ProductName, SourceDirectory, OutputDirectory, Verbose)
	// Platforms: %s`, ProductName, SourceDirectory, OutputDirectory, Debug, Release, Shipping, Verbose, platforms)

}
