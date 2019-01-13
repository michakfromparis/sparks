package platform

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/sparks"
)

type WebGl struct {
	enabled bool
}

func (p WebGl) Name() string {
	return "webgl"
}

func (p WebGl) Title() string {
	return "WebGl"
}

func (p WebGl) Opt() string {
	return "e"
}

func (p WebGl) Enabled() bool {
	return p.enabled
}

func (p WebGl) SetEnabled(enabled bool) {
	p.enabled = enabled
}

func (p WebGl) Deps() {

}
func (p WebGl) Clean() {

}
func (p WebGl) Build(configuration sparks.Configuration) {
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
