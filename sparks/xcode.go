package sparks

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/joomcode/errorx"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sys"
)

type SigningType int

const (
	IPhoneDeveloper = 0
	IphoneDistribution
	MacDeveloper
	MacDistribution
)

var SigningTypeNames = []string{
	"iPhone Developer",
	"iPhone Distribution",
	"Mac Developer",
	"Mac Distribution",
}

func (st SigningType) String() string {
	return SigningTypeNames[st]
}

var SigningIdentities = map[SigningType]string{}
var SigningIdentity string

type XCode struct {
	command       string
	arguments     []string
	platform      Platform
	configuration Configuration
}

func NewXCode(platform Platform, configuration Configuration) *XCode {
	xc := new(XCode)
	xc.command = "xcodebuild"
	xc.platform = platform
	xc.configuration = configuration
	return xc
}

func (xc *XCode) Build(directory string, arg ...string) error {

	var args []string

	args = append(args, "-parallelizeTargets")
	args = append(args, "-verbose")
	args = append(args, "build")
	args = append(args, "-project", config.ProductName+".xcodeproj")
	args = append(args, "-target", config.ProductName)
	args = append(args, "-sdk", config.SparksOSXSDK)
	args = append(args, "-configuration", xc.configuration.Title())
	args = append(args, arg...)
	output, err := sys.ExecuteEx(xc.command, directory, true, args...)
	if err != nil {
		return errorx.Decorate(err, "xcode build failed: "+output)
	}
	return nil
}

func (xc *XCode) Clean() {

}

func (xc *XCode) DetectSigning() {
	log.Debug("detecting xcode signing identity")
	s, err := sys.Execute("security", "find-identity", "-v", "-p", "codesigning")
	if err != nil {
		errx.Fatalf(err, "security find-identity failed")
		return
	}
	lines := strings.Split(s, sys.NewLine)
	for _, line := range lines {
		parts := strings.Split(line, "\"")
		if len(parts) == 3 {
			identity := parts[1]
			for i := 0; i < len(SigningTypeNames); i++ {
				if strings.Contains(identity, SigningTypeNames[i]) {
					SigningIdentities[SigningType(i)] = identity
					log.Tracef("detected identity: %s", identity)
				}
			}
		}
	}
}

func (xc *XCode) SigningIdentity(signingType SigningType) string {
	return SigningIdentities[signingType]
}
