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
var ProductName string
var Debug = false
var Release = false
var Shipping = false

var Platforms map[string]bool
var Verbose = false

func Init() error {
	log.Info("initializing config")
	ProductName = "Sparks"

	var err error
	if SourceDirectory, err = utils.Pwd(); err != nil {
		errx.Fatalf(err, "could not determine current working directory")
	}
	OutputDirectory = path.Join(SourceDirectory, "build")
	Platforms = make(map[string]bool)
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
Verbose: %t
Platforms: %s`, ProductName, SourceDirectory, OutputDirectory, Debug, Release, Shipping, Verbose)
	// Platforms: %s`, ProductName, SourceDirectory, OutputDirectory, Debug, Release, Shipping, Verbose, platforms)

}
