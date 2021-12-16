package sparks

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	"leblox.com/sparks-cli/v2/conf"
	"leblox.com/sparks-cli/v2/errx"
	"leblox.com/sparks-cli/v2/sys"
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

// SigningIdentity is a class representing all the fields of an OSX / iOS signing identity
type SigningIdentity struct {
	ID     string
	Name   string
	TeamID string
}

func (s *SigningIdentity) String() string {
	return fmt.Sprintf("%s %s", s.ID, s.Name)
}

// SigningIdentities holds the list of detected signing identities mapped by SigningType
var SigningIdentities = map[SigningType]*SigningIdentity{}

// CurrentSigningIdentity holds a reference to the signing identity that will be used
var CurrentSigningIdentity *SigningIdentity

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
	xc.arguments = append(xc.arguments, "-project", conf.ProductName+".xcodeproj")
	xc.arguments = append(xc.arguments, "-target", conf.ProductName)
	if xc.platform.Name() != "ios" {
		xc.arguments = append(xc.arguments, "-sdk", conf.SparksOSXSDK)
	}
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
	s, err := sys.Execute("/usr/bin/env", "xcrun", "security", "find-identity", "-v", "-p", "codesigning")
	if err != nil {
		errx.Fatalf(err, "security find-identity failed")
		return
	}
	lines := strings.Split(s, sys.NewLine)
	lines = lines[:len(lines)-2]
	for _, line := range lines {
		re := regexp.MustCompile("\\s+[0-9]+\\)\\s([0-9a-fA-F]+)\\s\"(.*\\((.*)\\))")
		matched := re.MatchString(line)
		if !matched {
			log.Warnf("Could not parse xcode signing identity: " + line)
			continue
		}
		matches := re.FindStringSubmatch(line)
		if len(matches) != 4 {
			log.Warnf("Could not parse xcode signing identity: " + line)
			continue
		}
		ID := matches[1]
		name := matches[2]
		teamID := matches[3]
		identity := SigningIdentity{ID: ID, Name: name, TeamID: teamID}
		// fmt.Printf("%+v", identity.String())

		for i := 0; i < len(SigningTypeNames); i++ {
			if strings.Contains(identity.Name, SigningTypeNames[i]) {
				SigningIdentities[SigningType(i)] = &identity
				log.Tracef("detected identity: %s", identity.String())
			}
		}
	}
}

// SelectSigning returns the configured signing identity provided a signing type
func (xc *XCode) SelectSigning(signingType SigningType) (*SigningIdentity, error) {
	if SigningIdentities[signingType] == nil {
		return nil, fmt.Errorf("Could not select a %+v signing identity", signingType)
	}
	CurrentSigningIdentity = SigningIdentities[signingType]
	// log.Tracef("xcode signing identity: %s", SigningIdentity)
	return CurrentSigningIdentity, nil
}
