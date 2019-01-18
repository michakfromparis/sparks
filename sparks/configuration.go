package sparks

import (
	log "github.com/Sirupsen/logrus"
)

// Configuration Interfaxce
type Configuration interface {
	Name() string
	Title() string
	Opt() string
	Enabled() bool
	SetEnabled(bool)
}

// Map of all Configurations
var Configurations = map[string]Configuration{}

// Ordered array of Configuration keys
var ConfigurationNames = []string{
	"debug",
	"release",
	"shipping",
}

func RegisterConfiguration(configuration Configuration) {
	log.Debug("registering configuration: " + configuration.Title())
	Configurations[configuration.Name()] = configuration
}

func SetEnabledConfigurations(configurations []bool) {
	i := 0
	for _, name := range ConfigurationNames {
		if i < len(configurations) && configurations[i] == true {
			Configurations[name].SetEnabled(true)
			log.Debug("enabled configuration " + Configurations[name].Title())
		}
		i++
	}
}
