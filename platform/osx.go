package platform

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

type Osx struct {
}

func (p Osx) Opt() string {
	return "o"
}

func (p Osx) Title() string {
	return "OSX"
}

func (p Osx) Name() string {
	return "osx"
}

func (p Osx) Deps() {
	log.Info("Installing dependencies for " + p.Title())
}
func (p Osx) Clean() {

}
func (p Osx) Build(configuration Configuration) {
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
