package main

import (
	"github.com/michaKFromParis/sparks/cmd"
	"github.com/michaKFromParis/sparks/errx"
)

func main() {
	if err := cmd.Execute(); err != nil {
		errx.Fatal(err)
	}
}
