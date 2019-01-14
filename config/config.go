package config

import (
	"fmt"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/utils"
)

var SourceDirectory string
var OutputDirectory string
var SDKDirectory string
var SDKName string
var ProductName string
var Debug = false
var Release = false
var Shipping = false

var Platforms map[string]bool
var Verbose = false
var VeryVerbose = false // TODO temporary for debugging. need to set a proper verbose level

func Init() error {
	log.Info("initializing sparks")
	SDKName = "Sparks"
	ProductName = "SparksPlayer"

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
Debug: %t
Release: %t
Shipping: %t
Verbose: %t`, ProductName, SourceDirectory, OutputDirectory, Debug, Release, Shipping, Verbose)
	// Platforms: %s`, ProductName, SourceDirectory, OutputDirectory, Debug, Release, Shipping, Verbose, platforms)

}
