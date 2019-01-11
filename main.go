package main

import (
	"github.com/michaKFromParis/sparks/cmd"
	"github.com/michaKFromParis/sparks/log"
)

func main() {
	log.Init()
	cmd.Execute()
}
