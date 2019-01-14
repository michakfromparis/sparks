package platform

import (
	"github.com/michaKFromParis/sparks/sparks"
)

type WebGl struct {
	enabled bool
}

func (w *WebGl) Name() string {
	return "webgl"
}

func (w *WebGl) Title() string {
	return "WebGl"
}

func (w *WebGl) Opt() string {
	return "e"
}

func (w *WebGl) Enabled() bool {
	return w.enabled
}

func (w *WebGl) SetEnabled(enabled bool) {
	w.enabled = enabled
}

func (w *WebGl) Deps() {

}
func (w *WebGl) Clean() {

}
func (w *WebGl) Build(configuration sparks.Configuration) {
	w.prebuild()
	w.generate()
	w.compile()
	w.postbuild()
}
func (w *WebGl) prebuild() {

}
func (w *WebGl) generate() {

}
func (w *WebGl) compile() {

}
func (w *WebGl) postbuild() {

}
