package platform

import (
	"fmt"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
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

func (w *WebGl) Get() error {
	return nil
}
func (w *WebGl) Clean() error {
	return nil
}
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
	params := fmt.Sprintf("-DOS_EMSCRIPTEN=1")
	params += fmt.Sprintf("-DCMAKE_TOOLCHAIN_FILE=%s", cmakeToolchainFile)
	params += fmt.Sprintf("-DEMSCRIPTEN_ROOT_PATH=${EmscriptenSDKRoot}/emscripten/${EmscriptenVersion} ")
	projectsPath := filepath.Join(config.OutputDirectory, "projects", w.Title()+"-"+configuration.Title())
	out, err := cmake.Run(projectsPath)
	if err != nil {
		errx.Fatalf(err, "sparks project generate failed: "+out)
	}
	log.Trace("cmake output" + out)
}

func (w *WebGl) compile() {

}
func (w *WebGl) postbuild() {

}
