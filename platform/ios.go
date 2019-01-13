package platform

import (
	"github.com/michaKFromParis/sparks/sparks"

	"fmt"

	log "github.com/Sirupsen/logrus"
)

type Ios struct {
	enabled bool
}

func (i Ios) Name() string {
	return "ios"
}

func (i Ios) Title() string {
	return "iOS"
}

func (i Ios) Opt() string {
	return "i"
}

func (i Ios) String() string {
	return i.Title()
}

func (i Ios) Enabled() bool {
	return i.enabled
}

func (i Ios) SetEnabled(enabled bool) {
	i.enabled = enabled
}

func (i Ios) Deps() {

}
func (i Ios) Clean() {

}
func (i Ios) Build(configuration sparks.Configuration) {
	log.Info(fmt.Sprintf("Building %s-%s", i.Title(), configuration))
	i.prebuild()
	i.generate()
	i.compile()
	i.postbuild()
}

func (i Ios) prebuild() {

}
func (i Ios) generate() {

}
func (i Ios) compile() {

}
func (i Ios) postbuild() {

}
