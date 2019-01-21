package config

import (
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sys"
)

// GetDependencies stores the command line flag of go get --dependencies
var GetDependencies bool

// SourceDirectory points to the directory containing the sources to be built
var SourceDirectory string

// OutputDirectory points to the directory where the product will be built
var OutputDirectory string

// SDKDirectory points to the directory containing the sparks sdk
var SDKDirectory string

// SDKName is the name of the sparks sdk
var SDKName = "Sparks"

// PlayerName is the name of the sparks player
var PlayerName = "SparksPlayer"

// ProductName is the name of the built product
var ProductName string

// Verbose flag to show verbose log output
var Verbose = false

// VeryVerbose flag to show very verbose log output // TODO temporary for debugging. need to set a proper verbose level
var VeryVerbose = false

// IncludeSparksSource determines wether the sparks source should be included when generating projects
var IncludeSparksSource = true

// GenerateLua determines wether the sparks lua bindings should be generated
var GenerateLua = false

// OSX specific

// SparksOSXSDK determines the OSX SDK sparks will be built with. run `xcodebuild -showsdks` to list valid values
var SparksOSXSDK = "macosx10.14"

// SparksOSXArchitecture determines the OSX built architectures
var SparksOSXArchitecture = "x86_64"

// SparksOSXDeploymentTarget determines the minimum deployment target for a Sparks OSX application
var SparksOSXDeploymentTarget = "10.10"

// XCodeSigningIdentity is usually set at runtime using XCode.DetectSigning
var XCodeSigningIdentity = ""

// iOS specific

// SparksiOSSDK determines the iOS SDK sparks will be built with. run `xcodebuild -showsdks` to list valid values
var SparksiOSSDK = "iphoneos12.1"

// sparksiOSSimulatorSDK determines the iPhone Simulator SDK sparks will be built with. run `xcodebuild -showsdks` to list valid values
var sparksiOSSimulatorSDK = "iphonesimulator12.1"

// DevelopmentSigningType is the type signing type used when building Debug or Release builds
var DevelopmentSigningType = "iPhone Developer"

// DistributionSigningType is the type signing type used when building Shipping builds
var DistributionSigningType = "iPhone Distribution"

// BundleIdentifier is the bundle identifier of the iOS ipa
var BundleIdentifier = "me.ivoltage.sparksplayer"

// ProvisioningProfileUUID is the identifier of the provisionning profile
var ProvisioningProfileUUID = ""

// Android specific

// AndroidSDKRoot points to the location where the Android sdk is installed
var AndroidSDKRoot = filepath.Join("Android", "sdk")

// AndroidNDKRoot points to the location where the Android ndk is installed
var AndroidNDKRoot = filepath.Join("Android", "ndk")

// AndroidSDKToolsVersion defines the android sdk tools version to use
var AndroidSDKToolsVersion = "25.2.3"

// SparksAndroidBuildToolsVersion defines the android sdk build tools version to use
var SparksAndroidBuildToolsVersion = "25.2.3"

// SparksAndroidAPILevel defines the minimum android os requirement to run a sparks application
var SparksAndroidAPILevel = 19

// SpakrsAndroidNDKVersion defines the minimum android NDK version required to build a sparks application
var SpakrsAndroidNDKVersion = "r10d"

// SparkzsAndroidPackagePrefix defines the package prefix used for the android bundle identifier
var SparkzsAndroidPackagePrefix = "me.ivoltage"

// SparksAndroidNDKUrlOSX specifies the url to download the android NDK for OSX
var SparksAndroidNDKUrlOSX = "http://dl.google.com/android/ndk/android-ndk-${spakrsAndroidNDKVersion}-darwin-x86_64.bin"

// SparksAndroidNDKUrlLinux specifies the url to download the android NDK for Linux
var SparksAndroidNDKUrlLinux = "http://dl.google.com/android/ndk/android-ndk-${spakrsAndroidNDKVersion}-linux-x86_64.bin"

// FacebookAndroidAPILevel specifies android api level for the Facebook SDK
var FacebookAndroidAPILevel = 9

// VolleyAndroidAPILevel specifies android api level for the Volley SDK
var VolleyAndroidAPILevel = 19

// Emscripten specific

// EmscriptenSDKRoot points to the location where the Emscripten sdk is installed
var EmscriptenSDKRoot = "Emscripten"

// EmscriptenVersion defines the version of the Emscripten SDK to use. Possible values: '1.27.0', '1.29.0' 'latest', 'master' or 'incoming' for the latest version
var EmscriptenVersion = "incoming"

// EmscriptenBrowser defines the browser started when using emrun to start an Emscripten application. Possible values: 'chrome' 'firefox' 'safari'
var EmscriptenBrowser = "firefox"

// Windows specific

// WindowsCompiler specifies the Windows compiler to use to build. Possible values are "VisualStudio" or "MSBuild"
var WindowsCompiler = "VisualStudio"

// Init initializes config
func Init() error {
	log.Info("sparks config init")
	if ProductName == "" {
		ProductName = PlayerName
	}
	var err error
	if SourceDirectory, err = sys.Pwd(); err != nil {
		errx.Fatalf(err, "could not determine current working directory")
	}
	SDKDirectory = SourceDirectory
	OutputDirectory = path.Join(SourceDirectory, "build")
	log.Debug("loaded product configuration for " + ProductName)
	return nil
}
