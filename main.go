package main

import (
	"github.com/michakfromparis/sparks/cmd"
	"github.com/michakfromparis/sparks/config"
	"github.com/michakfromparis/sparks/errx"
	"github.com/michakfromparis/sparks/platform"
)

func main() {
	platform.RegisterPlatforms()
	config.RegisterConfigurations()
	cmd.Init()
	cmd.SetVersion(version, commit, date)
	if err := cmd.Execute(); err != nil {
		errx.Fatal(err)
	}
}
