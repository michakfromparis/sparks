package platform

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joomcode/errorx"
	"github.com/michakfromparis/sparks/sys"

	log "github.com/sirupsen/logrus"
	"github.com/michakfromparis/sparks/conf"
	"github.com/michakfromparis/sparks/errx"
	"github.com/michakfromparis/sparks/sparks"
)

// WebGl represents the WebGl platform
type WebGl struct {
	enabled       bool
	configuration sparks.Configuration
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
	deps := sparks.Deps{}
	var err error

	if err = deps.Update(); err != nil {
		return err
	}
	if err = deps.Get("python", "unzip", "cmake", "git", "g++", "default-jre"); err != nil {
		return err
	}

	installNeeded := false
	installVersion := false

	if err = w.createEmscriptenSDKRoot(); err != nil {
		return err
	}
	if installNeeded, err = w.checkEmscriptenSDKInstallion(); err != nil {
		return err
	}
	if installNeeded {
		output, err := sys.Execute("git", "clone", "--depth", "1", "https://github.com/juj/emsdk.git", conf.EmscriptenSDKRoot)
		if err != nil {
			return errorx.Decorate(err, "failed to clone Emscripten repo: %s", output)
		}
	}
	log.Debug("Emscripten SDK installed at: " + conf.EmscriptenSDKRoot)

	if installVersion, err = w.checkEmscriptenSDKVersionInstallion(); err != nil {
		return err
	}
	if installVersion {
		if !installNeeded {
			output, err := sys.Execute("git", "pull", conf.EmscriptenSDKRoot)
			if err != nil {
				return errorx.Decorate(err, "failed to update Emscripten SDK: %s", output)
			}
		}
		emsdkPath := filepath.Join(conf.EmscriptenSDKRoot, "emsdk")
		output, err := sys.Execute(emsdkPath, "install", "--build=Release", "sdk-"+conf.EmscriptenVersion+"-64bit", "binaryen-master-64bit")
		if err != nil {
			return errorx.Decorate(err, "failed to emsdk install version %s: %s", conf.EmscriptenVersion, output)
		}
		output, err = sys.Execute(emsdkPath, "activate", "--build=Release", "sdk-"+conf.EmscriptenVersion+"-64bit", "binaryen-master-64bit")
		if err != nil {
			return errorx.Decorate(err, "failed to emsdk install version %s: %s", conf.EmscriptenVersion, output)
		}

		if err = w.createLatestSymlink(); err != nil {
			return err
		}

		log.Debug("Building a simple test to download and build SDL2 port")
		log.Debug("Otherwise, the build fails when calling make -j 8")
		if err = w.SetEnv(); err != nil {
			return errorx.Decorate(err, "failed to set emscripten environment")
		}
		emscriptenPath := os.Getenv("EMSCRIPTEN")
		output, err = sys.ExecuteEx("emcc", "/tmp", true, "-s", "USE_SDL=2", filepath.Join(emscriptenPath, "tests", "hello_world.c"))
		if err != nil {
			return fmt.Errorf("failed to build sdl test: %s", output)
		}
	}
	return nil
}

// Clean cleans the platform build
func (w *WebGl) Clean() error {
	return nil
}

// Code opens the code editor for the project
func (w *WebGl) Code(configuration sparks.Configuration) error {
	return nil
}

// Build builds the platform
func (w *WebGl) Build(configuration sparks.Configuration) error {
	w.configuration = configuration
	w.prebuild()
	w.generate()
	w.compile()
	w.postbuild()
	return nil
}

func (w *WebGl) prebuild() {

}

func (w *WebGl) generate() {
	log.Info("sparks project generate --webgl")

	cmakeToolchainFile := filepath.Join(conf.SDKDirectory, "scripts", "CMake", "toolchains", "Emscripten.cmake")
	cmake := sparks.NewCMake(w, w.configuration)
	cmake.AddDefine("OS_EMSCRIPTEN", "1")
	cmake.AddDefine("CMAKE_TOOLCHAIN_FILE", cmakeToolchainFile)
	cmake.AddDefine("EMSCRIPTEN_ROOT_PATH", filepath.Join(filepath.Join(conf.EmscriptenSDKRoot, "emscripten"), conf.EmscriptenVersion))
	projectsPath := filepath.Join(conf.OutputDirectory, "projects", w.Title()+"-"+w.configuration.Title())
	out, err := cmake.Run(projectsPath)
	if err != nil {
		errx.Fatalf(err, "sparks project generate failed: "+out)
	}
	//log.Trace("cmake output" + out)
}

func (w *WebGl) compile() {
	log.Info("sparks project compile --webgl")
	projectsPath := filepath.Join(conf.OutputDirectory, "projects", w.Title()+"-"+w.configuration.Title())
	out, err := sys.ExecuteEx("make", projectsPath, true, "-j8")
	if err != nil {
		errx.Fatalf(err, "sparks project compile failed: "+out)
	}
}

func (w *WebGl) postbuild() {
}
