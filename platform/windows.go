package platform

import (
	"fmt"
	"path/filepath"

	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sparks"

	log "github.com/Sirupsen/logrus"
)

// Windows represents the Windows platform
type Windows struct {
	enabled bool
}

// Name is the lowercase name of the platform
func (w *Windows) Name() string {
	return "windows"
}

// Title is name of the platform
func (w *Windows) Title() string {
	return "Windows"
}

// Opt is the short command line option of the platform
func (w *Windows) Opt() string {
	return "w"
}

func (w *Windows) String() string {
	return w.Title()
}

// Enabled returns true if the platform is enabled
func (w *Windows) Enabled() bool {
	return w.enabled
}

// SetEnabled allows to enable / disable the platform
func (w *Windows) SetEnabled(enabled bool) {
	w.enabled = enabled
}

// Get installs the platform dependencies
func (w *Windows) Get() error {
	log.Info("Installing dependencies for " + w.Title())
	return nil
}

// Clean cleans the platform build
func (w *Windows) Clean() error {
	return nil
}

// Build builds the platform
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

	cmake := sparks.NewCMake(w, configuration)
	params := fmt.Sprintf("-DOS_WINDOWS=1")
	params += fmt.Sprintf("\"-G%s\" ", config.WindowsCompiler)
	params += "-DCMAKE_SYSTEM_NAME=Windows"
	projectsPath := filepath.Join(config.OutputDirectory, "projects", w.Title()+"-"+configuration.Title())
	out, err := cmake.Run(projectsPath)
	if err != nil {
		errx.Fatalf(err, "sparks project generate failed: "+out)
	}
	log.Trace("cmake output" + out)
}

func (w *Windows) compile() {

}

func (w *Windows) postbuild() {

}
