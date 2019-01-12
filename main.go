package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/cmd"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/logger"
)

func main() {
	logger.Init()
	config.Init()
	cmd.Execute()
	log.Info(config.Format())
}
