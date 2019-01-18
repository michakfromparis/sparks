package sparks

import (
	"fmt"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sys"
)

func getToluaPath() string {
	toluapp := filepath.Join(config.SDKDirectory, "dependencies", "toluapp", "bin")
	os, _ := sys.GetOs()
	switch os {
	case sys.Osx:
		toluapp = filepath.Join(toluapp, "toluapp_osx")
	case sys.Linux:
		toluapp = filepath.Join(config.SDKDirectory, "scripts", "bin", "tolua++")
	case sys.Windows:
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
	sys.SedFile(reflectionFile, "dofile\\(.*\\)", dofileWithCorrectPath)
	packagePath := filepath.Join(sourceDirectory, packageName)
	output, err := sys.ExecuteEx(
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
