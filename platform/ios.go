package platform

import (
	"github.com/michaKFromParis/sparks/sparks"

	"fmt"

	log "github.com/Sirupsen/logrus"
)

type Ios struct {
	enabled bool
}

func (p Ios) Name() string {
	return "ios"
}

func (p Ios) Title() string {
	return "iOS"
}

func (p Ios) Opt() string {
	return "i"
}

func (p Ios) Enabled() bool {
	return p.enabled
}

func (p Ios) SetEnabled(enabled bool) {
	p.enabled = enabled
}

func (p Ios) Deps() {

}
func (p Ios) Clean() {

}
func (p Ios) Build(configuration sparks.Configuration) {
	log.Info(fmt.Sprintf("Building %s-%s", p.Title(), configuration))
	p.prebuild()
	p.generate()
	p.compile()
	p.postbuild()
}

func (p Ios) prebuild() {

}
func (p Ios) generate() {

}
func (p Ios) compile() {

}
func (p Ios) postbuild() {

}
