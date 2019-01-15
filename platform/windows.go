package platform

import (
	"fmt"

	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/sparks"

	log "github.com/Sirupsen/logrus"
)

type Windows struct {
	enabled bool
}

func (w *Windows) Name() string {
	return "windows"
}

func (w *Windows) Title() string {
	return "Windows"
}

func (w *Windows) Opt() string {
	return "w"
}

func (w *Windows) String() string {
	return w.Title()
}

func (w *Windows) Enabled() bool {
	return w.enabled
}

func (w *Windows) SetEnabled(enabled bool) {
	w.enabled = enabled
}

func (w *Windows) Deps() error {
	log.Info("Installing dependencies for " + w.Title())
	return nil
}
func (w *Windows) Clean() error {
	return nil
}
func (w *Windows) Build(configuration sparks.Configuration) error {
	w.prebuild()
	w.generate(configuration)
	w.compile()
	w.postbuild()
	return nil
}

func (w *Windows) prebuild() {
}

func (w *Windows) generate(configuration sparks.Configuration) {
	log.Info("sparks project generate --windows")
	params := generateCmakeCommon(w, configuration)
	params += fmt.Sprintf("-DOS_WINDOWS=1")
	params += fmt.Sprintf("\"-G%s\" ", config.WindowsCompiler)
	params += "-DCMAKE_SYSTEM_NAME=Windows"
}

func (w *Windows) compile() {

}

func (w *Windows) postbuild() {

}
