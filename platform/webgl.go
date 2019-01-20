package platform

import (
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sparks"
)

// WebGl represents the WebGl platform
type WebGl struct {
	enabled bool
}

// Name is the lowercase name of the platform
func (w *WebGl) Name() string {
	return "webgl"
}

// Title is name of the platform
func (w *WebGl) Title() string {
	return "WebGl"
}

// Opt is the short command line option of the platform
func (w *WebGl) Opt() string {
	return "e"
}

// Enabled returns true if the platform is enabled
func (w *WebGl) Enabled() bool {
	return w.enabled
}

// SetEnabled allows to enable / disable the platform
func (w *WebGl) SetEnabled(enabled bool) {
	w.enabled = enabled
}

// Get installs the platform dependencies
func (w *WebGl) Get() error {
	log.Info("Installing dependencies for " + w.Title())
	return nil
}

// Clean cleans the platform build
func (w *WebGl) Clean() error {
	return nil
}

// Build builds the platform
func (w *WebGl) Build(configuration sparks.Configuration) error {
	w.prebuild()
	w.generate(configuration)
	w.compile()
	w.postbuild()
	return nil
}

func (w *WebGl) prebuild() {

}

func (w *WebGl) generate(configuration sparks.Configuration) {
	log.Info("sparks project generate --webgl")

	cmakeToolchainFile := filepath.Join(config.SDKDirectory, "scripts", "CMake", "toolchains", "Emscripten.cmake")

	cmake := sparks.NewCMake(w, configuration)
	cmake.AddDefine("OS_EMSCRIPTEN", "1")
	cmake.AddDefine("CMAKE_TOOLCHAIN_FILE", cmakeToolchainFile)
	cmake.AddDefine("EMSCRIPTEN_ROOT_PATH", filepath.Join(filepath.Join(config.EmscriptenSDKRoot, "emscripten"), config.EmscriptenVersion))
	projectsPath := filepath.Join(config.OutputDirectory, "projects", w.Title()+"-"+configuration.Title())
	out, err := cmake.Run(projectsPath)
	if err != nil {
		errx.Fatalf(err, "sparks project generate failed: "+out)
	}
	//log.Trace("cmake output" + out)
}

func (w *WebGl) compile() {
}

func (w *WebGl) postbuild() {
}
