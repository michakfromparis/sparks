package sparks

import (
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/utils"
)

type XCode struct {
	command       string
	arguments     []string
	platform      Platform
	configuration Configuration
}

func (xc *XCode) Build(directory string) {

	var args []string

	args = append(args, "-parallelizeTargets")
	args = append(args, "-verbose")
	args = append(args, "build")
	args = append(args, "-project", config.ProductName+".xcodeproj")
	args = append(args, "-target", config.ProductName)
	args = append(args, "-sdk", config.SparksOSXSDK)
	args = append(args, "-configuration", xc.configuration.Title())

	utils.ExecuteEx("xcodebuild", directory, true, args...)
}

func (t *XCode) Clean() {

}

func NewXCode(platform Platform, configuration Configuration) *XCode {
	xc := new(XCode)
	xc.command = "cmake"
	xc.platform = platform
	xc.configuration = configuration
	return xc
}
