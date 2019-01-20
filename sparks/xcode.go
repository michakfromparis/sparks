package sparks

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/joomcode/errorx"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sys"
)

// SigningType represents the type of signing used to sign an OSX or iOS application
type SigningType int

const (
	// IPhoneDeveloper signing type
	IPhoneDeveloper = iota
	// IphoneDistribution signing type
	IphoneDistribution
	// MacDeveloper signing type
	MacDeveloper
	// MacDistribution signing type
	MacDistribution
)

// SigningTypeNames is used to String a SigningType into a string
var SigningTypeNames = []string{
	"iPhone Developer",
	"iPhone Distribution",
	"Mac Developer",
	"Mac Distribution",
}

func (st SigningType) String() string {
	return SigningTypeNames[st]
}

// SigningIdentities holds the list of detected signing type identities
var SigningIdentities = map[SigningType]string{}

// SigningIdentity holds the name of the signing identity that will be used
var SigningIdentity string

// XCode is a wrapper class around the xcodebuild command
type XCode struct {
	command       string
	arguments     []string
	platform      Platform
	configuration Configuration
}

// NewXCode returns a new instance of XCode
func NewXCode(platform Platform, configuration Configuration) *XCode {
	xc := new(XCode)
	xc.command = "xcodebuild"
	xc.platform = platform
	xc.configuration = configuration
	return xc
}

// Build starts to build the project
func (xc *XCode) Build(directory string, arg ...string) error {
	xc.arguments = append(xc.arguments, "-parallelizeTargets")
	xc.arguments = append(xc.arguments, "-verbose")
	xc.arguments = append(xc.arguments, "build")
	xc.arguments = append(xc.arguments, "-project", config.ProductName+".xcodeproj")
	xc.arguments = append(xc.arguments, "-target", config.ProductName)
	xc.arguments = append(xc.arguments, "-sdk", config.SparksOSXSDK)
	xc.arguments = append(xc.arguments, "-configuration", xc.configuration.Title())
	xc.arguments = append(xc.arguments, arg...)
	output, err := sys.ExecuteEx(xc.command, directory, true, xc.arguments...)
	if err != nil {
		return errorx.Decorate(err, "xcode build failed: "+output)
	}
	return nil
}

// Clean cleans a built project
func (xc *XCode) Clean() {

}

// DetectSigning detects the currently imported Signing identities into OSX / XCode
// and fills the sparks.SigningIdentities array with them
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

// SelectSigning returns the configured signing identity provided a signing type
func (xc *XCode) SelectSigning(signingType SigningType) (string, error) {
	if SigningIdentities[signingType] == "" {
		return "", fmt.Errorf("Could not select a %+v signing identity", signingType)
	}
	SigningIdentity = SigningIdentities[signingType]
	// log.Debugf("xcode signing identity: %s", SigningIdentity)
	return SigningIdentity, nil
}
