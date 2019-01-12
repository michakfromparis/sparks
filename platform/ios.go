package platform

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

type Ios struct {
}

func (p Ios) Opt() string {
	return "i"
}

func (p Ios) Title() string {
	return "iOS"
}

func (p Ios) Name() string {
	return "ios"
}

func (p Ios) Deps() {

}
func (p Ios) Clean() {

}
func (p Ios) Build(configuration Configuration) {
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
