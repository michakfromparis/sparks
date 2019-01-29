package platform

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/michaKFromParis/sparks/sys"

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
	deps := sparks.Deps{}
	deps.Update()
	deps.Get("python")
	deps.Get("unzip")
	deps.Get("cmake")
	deps.Get("git")
	deps.Get("g++")

	file, err := os.Stat(config.EmscriptenSDKRoot)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug("could not find Emscripten SDK directory at: " + config.EmscriptenSDKRoot)
			log.Debug("Creating it.")
			if err := os.MkdirAll(config.EmscriptenSDKRoot, os.ModePerm); err != nil {
				errx.Fatalf(err, "Could not create Emscripten SDK directory: "+config.EmscriptenSDKRoot)
			}
		} else {
			errx.Fatalf(err, "Could not open Emscripten directory: "+config.EmscriptenSDKRoot)
		}
	} else if !file.IsDir() {
		errx.Fatalf(err, "Emscripten SDK path is not a directory: "+config.EmscriptenSDKRoot)
	}

	install := false
	emsdkPath := filepath.Join(config.EmscriptenSDKRoot, "emsdk")
	file, err = os.Stat(emsdkPath)
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			errx.Fatalf(err, "Could not open emsdk: "+emsdkPath)
		}
	} else if file.IsDir() {
		errx.Fatalf(err, "Emscripten emsdk is a directory: "+emsdkPath)
	}

	if install {
		output, err := sys.Execute("git", "clone", "--depth", "1", "https://github.com/juj/emsdk.git", config.EmscriptenSDKRoot)
		if err != nil {
			return fmt.Errorf("failed to clone Emscripten repo: %s", output)
		}
	}
	log.Debug("Emscripten SDK installed at: " + config.EmscriptenSDKRoot)

	installVersion := false
	sdkVersionPath := filepath.Join(config.EmscriptenSDKRoot, "emscripten", config.EmscriptenVersion)
	file, err = os.Stat(sdkVersionPath)
	if err != nil {
		if os.IsNotExist(err) {
			installVersion = true
		} else {
			errx.Fatalf(err, "Could not open emscripten SDK version "+config.EmscriptenVersion+"at: "+sdkVersionPath)
		}
	} else if !file.IsDir() {
		errx.Fatalf(err, "Emscripten sdk version is not a directory: "+sdkVersionPath)
	}

	if installVersion {
		if !install {
			output, err := sys.Execute("git", "pull", config.EmscriptenSDKRoot)
			if err != nil {
				return fmt.Errorf("failed to update Emscripten SDK: %s", output)
			}
		}
		output, err := sys.Execute(emsdkPath, "install", "--build=Release", "sdk-"+config.EmscriptenVersion+"-64bit", "binaryen-master-64bit")
		if err != nil {
			return fmt.Errorf("failed to emsdk install version %s: %s", config.EmscriptenVersion, output)
		}
		output, err = sys.Execute(emsdkPath, "activate", "--build=Release", "sdk-"+config.EmscriptenVersion+"-64bit", "binaryen-master-64bit")
		if err != nil {
			return fmt.Errorf("failed to emsdk install version %s: %s", config.EmscriptenVersion, output)
		}
		log.Debug("Building a simple test to download and build SDL2 port")
		log.Debug("Otherwise, the build fails when calling make -j 8")
		// source "$EmscriptenSDKRoot/emsdk_env.sh" > /dev/null

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
