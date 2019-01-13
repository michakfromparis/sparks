package platform

import (
	"github.com/michaKFromParis/sparks/sparks"

	"fmt"

	log "github.com/Sirupsen/logrus"
)

type Osx struct {
	enabled bool
}

func (p Osx) Name() string {
	return "osx"
}

func (p Osx) Title() string {
	return "OSX"
}

func (p Osx) Opt() string {
	return "o"
}

func (p Osx) Enabled() bool {
	return p.enabled
}

func (p Osx) SetEnabled(enabled bool) {
	p.enabled = enabled
}

func (p Osx) Deps() {
	log.Info("Installing dependencies for " + p.Title())
}
func (p Osx) Clean() {

}
func (p Osx) Build(configuration sparks.Configuration) {
	log.Info(fmt.Sprintf("Building %s-%s", p.Title(), configuration))
	p.prebuild()
	p.generate()
	p.compile()
	p.postbuild()
}

func (p Osx) prebuild() {

}
func (p Osx) generate() {

}
func (p Osx) compile() {

}
func (p Osx) postbuild() {

}
