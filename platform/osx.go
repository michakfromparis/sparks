package platform

import (
	"github.com/michaKFromParis/sparks/sparks"

	log "github.com/Sirupsen/logrus"
)

type Osx struct {
	enabled bool
}

func (o *Osx) Name() string {
	return "osx"
}

func (o *Osx) Title() string {
	return "OSX"
}

func (o *Osx) Opt() string {
	return "o"
}

func (o *Osx) String() string {
	return o.Title()
}

func (o *Osx) Enabled() bool {
	return o.enabled
}

func (o *Osx) SetEnabled(enabled bool) {
	o.enabled = enabled
}

func (o *Osx) Deps() {
	log.Info("Installing dependencies for " + o.Title())
}
func (o *Osx) Clean() {

}
func (o *Osx) Build(configuration sparks.Configuration) {
	o.prebuild()
	o.generate()
	o.compile()
	o.postbuild()
}

func (o *Osx) prebuild() {

}
func (o *Osx) generate() {

}
func (o *Osx) compile() {

}
func (o *Osx) postbuild() {

}
