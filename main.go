package main

import (
	"leblox.com/sparks-cli/v2/cmd"
	"leblox.com/sparks-cli/v2/config"
	"leblox.com/sparks-cli/v2/errx"
	"leblox.com/sparks-cli/v2/platform"
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
