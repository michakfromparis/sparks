package platform

import (
	"fmt"
	"path/filepath"

	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sparks"

	log "github.com/Sirupsen/logrus"
)

// Linux represents the Linux platform
type Linux struct {
	enabled bool
}

// Name is the lowercase name of the platform
func (l *Linux) Name() string {
	return "linux"
}

// Title is name of the platform
func (l *Linux) Title() string {
	return "Linux"
}

// Opt is the short command line option of the platform
func (l *Linux) Opt() string {
	return "l"
}

func (l *Linux) String() string {
	return l.Title()
}

// Enabled returns true if the platform is enabled
func (l *Linux) Enabled() bool {
	return l.enabled
}

// SetEnabled allows to enable / disable the platform
func (l *Linux) SetEnabled(enabled bool) {
	l.enabled = enabled
}

// Get installs the platform dependencies
func (l *Linux) Get() error {
	log.Info("Installing dependencies for " + l.Title())
	return nil
}

// Clean cleans the platform build
func (l *Linux) Clean() error {
	return nil
}

// Build builds the platform
func (l *Linux) Build(configuration sparks.Configuration) error {
	l.prebuild()
	l.generate(configuration)
	l.compile()
	l.postbuild()
	return nil
}

func (l *Linux) prebuild() {
}

func (l *Linux) generate(configuration sparks.Configuration) {
	log.Info("sparks project generate --linux")

	cmake := sparks.NewCMake(l, configuration)
	params := fmt.Sprintf("-DOS_LINUX=1 ")
	params += fmt.Sprintf("\"-GCodeBlocks - Unix Makefiles\" ")
	projectsPath := filepath.Join(config.OutputDirectory, "projects", l.Title()+"-"+configuration.Title())
	out, err := cmake.Run(projectsPath)
	if err != nil {
		errx.Fatalf(err, "sparks project generate failed: "+out)
	}
	log.Trace("cmake output" + out)
}

func (l *Linux) compile() {

}

func (l *Linux) postbuild() {

}
