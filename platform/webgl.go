package platform

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

type WebGl struct {
}

func (p WebGl) Opt() string {
	return "e"
}

func (p WebGl) Title() string {
	return "WebGL"
}

func (p WebGl) Name() string {
	return "webgl"
}

func (p WebGl) Deps() {

}
func (p WebGl) Clean() {

}
func (p WebGl) Build(configuration Configuration) {
	log.Info(fmt.Sprintf("Building %s-%s", p.Title(), configuration))
	p.prebuild()
	p.generate()
	p.compile()
	p.postbuild()
}
func (p WebGl) prebuild() {

}
func (p WebGl) generate() {

}
func (p WebGl) compile() {

}
func (p WebGl) postbuild() {

}
