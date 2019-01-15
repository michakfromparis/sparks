package platform

import (
	"fmt"

	"github.com/michaKFromParis/sparks/sparks"

	log "github.com/Sirupsen/logrus"
)

type Linux struct {
	enabled bool
}

func (l *Linux) Name() string {
	return "linux"
}

func (l *Linux) Title() string {
	return "Linux"
}

func (l *Linux) Opt() string {
	return "l"
}

func (l *Linux) String() string {
	return l.Title()
}

func (l *Linux) Enabled() bool {
	return l.enabled
}

func (l *Linux) SetEnabled(enabled bool) {
	l.enabled = enabled
}

func (l *Linux) Deps() error {
	log.Info("Installing dependencies for " + l.Title())
	return nil
}
func (l *Linux) Clean() error {
	return nil
}
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

	params := generateCmakeCommon(l, configuration)
	params += fmt.Sprintf("-DOS_LINUX=1 ")
	params += fmt.Sprintf("\"-GCodeBlocks - Unix Makefiles\" ")
}

func (l *Linux) compile() {

}

func (l *Linux) postbuild() {

}
