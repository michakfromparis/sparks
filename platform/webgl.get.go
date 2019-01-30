package platform

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/joomcode/errorx"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/sys"
)

func (w *WebGl) createEmscriptenSDKRoot() error {
	file, err := os.Stat(conf.EmscriptenSDKRoot)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug("could not find Emscripten SDK directory at: " + conf.EmscriptenSDKRoot)
			log.Debug("Creating it.")
			if err := os.MkdirAll(conf.EmscriptenSDKRoot, os.ModePerm); err != nil {
				return errorx.Decorate(err, "Could not create Emscripten SDK directory: "+conf.EmscriptenSDKRoot)
			}
		} else {
			return errorx.Decorate(err, "Could not open Emscripten directory: "+conf.EmscriptenSDKRoot)
		}
	} else if !file.IsDir() {
		return errorx.Decorate(err, "Emscripten SDK path is not a directory: "+conf.EmscriptenSDKRoot)
	}
	return nil
}

func (w *WebGl) checkEmscriptenSDKInstallion() (bool, error) {
	emsdkPath := filepath.Join(conf.EmscriptenSDKRoot, "emsdk")
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
	sdkVersionPath := filepath.Join(conf.EmscriptenSDKRoot, "emscripten", conf.EmscriptenVersion)
	file, err := os.Stat(sdkVersionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, errorx.Decorate(err, "Could not open emscripten SDK version "+conf.EmscriptenVersion+"at: "+sdkVersionPath)
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
	symlink := filepath.Join(conf.EmscriptenSDKRoot, "emscripten", "latest")
	if err := os.Symlink(target, symlink); err != nil {
		return errorx.Decorate(err, "failed to create latest symlink")
	}
	return nil
}

// SetEnv generates the emsdk environment variable file, parses it and sets system env variables accordingly
func (w *WebGl) SetEnv() (rerr error) {
	log.Debug("Setting emscripten environment variables")
	emsdkEnv := filepath.Join(conf.EmscriptenSDKRoot, "emsdk_env.sh")
	output, err := sys.ExecuteEx("bash", "", true, "-c", emsdkEnv)
	if err != nil {
		return errorx.Decorate(err, "failed to call emsdk_env.sh: "+output)
	}
	emsdkSetEnv := filepath.Join(conf.EmscriptenSDKRoot, "emsdk_set_env.sh")
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
