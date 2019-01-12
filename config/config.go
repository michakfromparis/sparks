package config

import (
	"fmt"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/platform"
	"github.com/michaKFromParis/sparks/utils"
)

var SourceDirectory string
var OutputDirectory string
var ProductName string
var Platforms map[string]bool
var Debug = false
var Release = false
var Shipping = false

func Init() error {
	log.Info("Initializing config")
	ProductName = "Sparks"
	SourceDirectory = utils.Pwd()
	OutputDirectory = path.Join(SourceDirectory, "build")
	Platforms = make(map[string]bool)
	return nil
}

func String() string {
	platforms := ""
	for _, name := range platform.PlatformNames {
		if Platforms[name] {
			platforms += name + " "
		}
	}
	return fmt.Sprintf(
		`SourceDirectory: %s
		 OutputDirectory: %s
		 Debug: %t
		 Release: %t
		 Shipping: %t
		 Platforms: %s`, SourceDirectory, OutputDirectory, Debug, Release, Shipping, platforms)
}
