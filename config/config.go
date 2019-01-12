package config

import (
	"fmt"
	"path"

	"github.com/michaKFromParis/sparks/platform"
	"github.com/michaKFromParis/sparks/utils"
)

var SourceDirectory string
var OutputDirectory string
var Platforms map[string]bool
var Debug = false
var Release = false
var Shipping = false

func Init() {
	SourceDirectory = utils.Pwd()
	OutputDirectory = path.Join(SourceDirectory, "build")
	Platforms = make(map[string]bool)
}

func Format() string {
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
