package platform

import (
	"fmt"
	"path/filepath"

	"github.com/michaKFromParis/sparks/config"
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

func (w *WebGl) Deps() error {
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

func (w *WebGl) generate(configuration sparks.Configuration) string {
	params := generateCmakeCommon(w, configuration)
	cmakeToolchainFile := filepath.Join(config.SDKDirectory, "scripts", "CMake", "toolchains", "Emscripten.cmake")
	params += fmt.Sprintf("-DOS_EMSCRIPTEN=1 ")
	params += fmt.Sprintf("\"-DCMAKE_TOOLCHAIN_FILE%s\" ", cmakeToolchainFile)
	params += fmt.Sprintf("-DEMSCRIPTEN_ROOT_PATH=${EmscriptenSDKRoot}/emscripten/${EmscriptenVersion} ")
	return params
}

func (w *WebGl) compile() {

}
func (w *WebGl) postbuild() {

}
