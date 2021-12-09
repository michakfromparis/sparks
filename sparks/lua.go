package sparks

import (
	"fmt"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/michakfromparis/sparks/conf"
	"github.com/michakfromparis/sparks/errx"
	"github.com/michakfromparis/sparks/sys"
)

func getToluaPath() string {
	toluapp := filepath.Join(conf.SDKDirectory, "dependencies", "toluapp", "bin")
	os, err := sys.GetOs()
	if err != nil {
		errx.Fatal(err)
	}
	switch os {
	case sys.Osx:
		toluapp = filepath.Join(toluapp, "toluapp_osx")
	case sys.Linux:
		toluapp = filepath.Join(conf.SDKDirectory, "scripts", "bin", "tolua++")
	case sys.Windows:
		toluapp = filepath.Join(toluapp, "toluapp_script.exe")
	}
	return toluapp
}

// GenerateLuaBindings generates C/C++ code from the definition of a package
// defined with packageName. It invokes tolua++ and has complex toluahooks
// that also generates C++ class reflection
func GenerateLuaBindings(sourceDirectory string, packageName string) {

	log.Info("sparks lua bind " + packageName)
	toluapp := getToluaPath()
	toluaHooksPath := filepath.Join(conf.SDKDirectory, "src", "Sparks", "tolua.hooks.lua")
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
