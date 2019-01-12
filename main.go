package main

import (
	"github.com/michaKFromParis/sparks/cmd"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/logger"
)

func main() {
	if err := logger.Init(); err != nil {
		errx.Fatal(err, "Logger initialization failed")
	}
	if err := config.Init(); err != nil {
		errx.Fatal(err, "Configuration initialization failed")
	}
	cmd.Execute()
}
