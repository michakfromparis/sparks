package main

import (
	"github.com/michaKFromParis/sparks/cmd"
	"github.com/michaKFromParis/sparks/configuration"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/logger"
	"github.com/michaKFromParis/sparks/platform"
)

func main() {
	logger.Init()
	platform.RegisterPlatforms()
	configuration.RegisterConfigurations()
	cmd.Init()
	if err := cmd.Execute(); err != nil {
		errx.Fatal(err)
	}
}
