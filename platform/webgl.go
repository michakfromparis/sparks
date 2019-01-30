package platform

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/joomcode/errorx"
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
	deps.Get("default-jre")

	var err error
	installNeeded := false
	installVersion := false

	if err = w.createEmscriptenSDKRoot(); err != nil {
		return err
	}
	if installNeeded, err = w.checkEmscriptenSDKInstallion(); err != nil {
		return err
	}
	if installNeeded {
		output, err := sys.Execute("git", "clone", "--depth", "1", "https://github.com/juj/emsdk.git", config.EmscriptenSDKRoot)
		if err != nil {
			return errorx.Decorate(err, "failed to clone Emscripten repo: %s", output)
		}
	}
	log.Debug("Emscripten SDK installed at: " + config.EmscriptenSDKRoot)

	if installVersion, err = w.checkEmscriptenSDKVersionInstallion(); err != nil {
		return err
	}
	if installVersion {
		if !installNeeded {
			output, err := sys.Execute("git", "pull", config.EmscriptenSDKRoot)
			if err != nil {
				return errorx.Decorate(err, "failed to update Emscripten SDK: %s", output)
			}
		}
		emsdkPath := filepath.Join(config.EmscriptenSDKRoot, "emsdk")
		output, err := sys.Execute(emsdkPath, "install", "--build=Release", "sdk-"+config.EmscriptenVersion+"-64bit", "binaryen-master-64bit")
		if err != nil {
			return errorx.Decorate(err, "failed to emsdk install version %s: %s", config.EmscriptenVersion, output)
		}
		output, err = sys.Execute(emsdkPath, "activate", "--build=Release", "sdk-"+config.EmscriptenVersion+"-64bit", "binaryen-master-64bit")
		if err != nil {
			return errorx.Decorate(err, "failed to emsdk install version %s: %s", config.EmscriptenVersion, output)
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

func (w *WebGl) createEmscriptenSDKRoot() error {
	file, err := os.Stat(config.EmscriptenSDKRoot)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug("could not find Emscripten SDK directory at: " + config.EmscriptenSDKRoot)
			log.Debug("Creating it.")
			if err := os.MkdirAll(config.EmscriptenSDKRoot, os.ModePerm); err != nil {
				return errorx.Decorate(err, "Could not create Emscripten SDK directory: "+config.EmscriptenSDKRoot)
			}
		} else {
			return errorx.Decorate(err, "Could not open Emscripten directory: "+config.EmscriptenSDKRoot)
		}
	} else if !file.IsDir() {
		return errorx.Decorate(err, "Emscripten SDK path is not a directory: "+config.EmscriptenSDKRoot)
	}
	return nil
}

func (w *WebGl) checkEmscriptenSDKInstallion() (bool, error) {
	emsdkPath := filepath.Join(config.EmscriptenSDKRoot, "emsdk")
	file, err := os.Stat(emsdkPath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, errorx.Decorate(err, "Could not open emsdk: "+emsdkPath)
	}
	if file.IsDir() {
		return false, errorx.Decorate(err, "emscripten emsdk is a directory: "+emsdkPath)
	}
	return true, nil
}

func (w *WebGl) checkEmscriptenSDKVersionInstallion() (bool, error) {
	sdkVersionPath := filepath.Join(config.EmscriptenSDKRoot, "emscripten", config.EmscriptenVersion)
	file, err := os.Stat(sdkVersionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, errorx.Decorate(err, "Could not open emscripten SDK version "+config.EmscriptenVersion+"at: "+sdkVersionPath)
	}
	if !file.IsDir() {
		return false, errorx.Decorate(err, "Emscripten sdk version is not a directory: "+sdkVersionPath)
	}
	return true, nil
}

func (w *WebGl) createLatestSymlink() error {
	log.Debug("creating emscripten latest symlink")
	w.SetEnv()
	target := os.Getenv("EMSCRIPTEN")
	symlink := filepath.Join(config.EmscriptenSDKRoot, "emscripten", "latest")
	if err := os.Symlink(target, symlink); err != nil {
		return errorx.Decorate(err, "failed to create latest symlink")
	}
	return nil
}

// SetEnv generates the emsdk environment variable file, parses it and sets system env variables accordingly
func (w *WebGl) SetEnv() (rerr error) {
	log.Debug("Setting emscripten environment variables")
	emsdkEnv := filepath.Join(config.EmscriptenSDKRoot, "emsdk_env.sh")
	output, err := sys.ExecuteEx("bash", "", true, "-c", emsdkEnv)
	if err != nil {
		return errorx.Decorate(err, "failed to call emsdk_env.sh: "+output)
	}
	emsdkSetEnv := filepath.Join(config.EmscriptenSDKRoot, "emsdk_set_env.sh")
	file, err := os.Open(emsdkSetEnv)
	if err != nil {
		return errorx.Decorate(err, "could not open emsdk_set_env.sh")
	}
	defer func() {
		err := file.Close()
		if err != nil {
			rerr = err
		}
	}()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile("export (.*)=\"(.*)\"")
	lineNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++
		if !re.MatchString(line) {
			return fmt.Errorf("could not parse line %d of emsdk_set_env.sh: %s", lineNumber, line)
		}
		envVarName := re.FindStringSubmatch(line)[1]
		envVarValue := re.FindStringSubmatch(line)[2]
		log.Debugf("%s=%s", envVarName, envVarValue)
		if err = os.Setenv(envVarName, envVarValue); err != nil {
			return errorx.Decorate(err, "failed to set environment variable %s to %s", envVarName, envVarValue)
		}
	}
	if err := scanner.Err(); err != nil {
		return errorx.Decorate(err, "could not read lines of emsdk_set_env.sh")
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
