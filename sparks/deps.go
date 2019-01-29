package sparks

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/sys"
)

// Deps is a wrapper class around platform specific dependencies installer
type Deps struct {
}

// Update updates the dependencies sources list
func (d *Deps) Update() error {
	log.Debugf("sparks get -du")
	os, _ := sys.GetOs()
	switch os {
	case sys.Linux:
		output, err := sys.Execute("apt-get", "update")
		if err != nil {
			return fmt.Errorf("failed to update: %s", output)
		}
	case sys.Osx:
		output, err := sys.Execute("brew", "update")
		if err != nil {
			return fmt.Errorf("failed to update: %s", output)
		}
	}
	return nil
}

// Get install a package from its name
func (d *Deps) Get(name string) error {
	log.Debugf("sparks get %s", name)
	os, _ := sys.GetOs()
	switch os {
	case sys.Linux:
		output, err := sys.Execute("apt-get", "install", "-y", name)
		if err != nil {
			return fmt.Errorf("failed to install %s: %s", name, output)
		}
	case sys.Osx:
		output, err := sys.Execute("brew", "install", name)
		if err != nil {
			return fmt.Errorf("failed to install %s: %s", name, output)
		}
	}
	return nil
}
