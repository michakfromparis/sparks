package sparks

import (
	"fmt"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/utils"
)

func getToluaPath() string {
	toluapp := filepath.Join(config.SDKDirectory, "dependencies", "toluapp", "bin")
	os, _ := utils.GetOs()
	switch os {
	case utils.Osx:
		toluapp = filepath.Join(toluapp, "toluapp_osx")
	case utils.Linux:
		toluapp = filepath.Join(config.SDKDirectory, "scripts", "bin", "tolua++")
	case utils.Windows:
		toluapp = filepath.Join(toluapp, "toluapp_script.exe")
	}
	return toluapp
}

func GenerateLuaBindings(sourceDirectory string, packageName string) {

	log.Info("sparks lua bind " + packageName)
	toluapp := getToluaPath()
	toluaHooksPath := filepath.Join(config.SDKDirectory, "src", "Sparks", "tolua.hooks.lua")
	dofileWithCorrectPath := fmt.Sprintf("dofile(\"%s\")", toluaHooksPath)
	reflectionFile := filepath.Join(sourceDirectory, packageName+".Reflection.lua")
	utils.SedFile(reflectionFile, "dofile\\(.*\\)", dofileWithCorrectPath)
	packagePath := filepath.Join(sourceDirectory, packageName)
	output, err := utils.ExecuteEx(
		toluapp,
		sourceDirectory,
		true,
		"-L", packagePath+".Reflection.lua",
		"-n", packageName,
		"-o", packagePath+".tolua.cpp",
		"-H", packagePath+".tolua.h",
		packagePath+".pkg")
	if err != nil {
		errx.Fatalf(err, "tolua++ execution failed: "+output)
	}
}
