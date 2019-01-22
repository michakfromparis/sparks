package main

import (
	"github.com/michaKFromParis/sparks/cmd"
	"github.com/michaKFromParis/sparks/configuration"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/platform"
)

func main() {
	platform.RegisterPlatforms()
	configuration.RegisterConfigurations()
	cmd.Init()
	cmd.SetVersion(version, commit, date)
	if err := cmd.Execute(); err != nil {
		errx.Fatal(err)
	}
}
