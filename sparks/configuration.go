package sparks

import (
	log "github.com/Sirupsen/logrus"
)

// Configuration Interface used to represent a build configuration
type Configuration interface {
	Name() string
	Title() string
	Opt() string
	Enabled() bool
	SetEnabled(bool)
}

// Configurations is the map of all registered configurations
var Configurations = map[string]Configuration{}

// ConfigurationNames is an ordered array of Configuration keys used to iterate over Configurations
var ConfigurationNames = []string{
	"debug",
	"release",
	"shipping",
}

// RegisterConfiguration allows external code to register a new configuration as long as it respects the Configuration interface
func RegisterConfiguration(configuration Configuration) {
	log.Debug("registering configuration: " + configuration.Title())
	Configurations[configuration.Name()] = configuration
}

// SetEnabledConfigurations is used to enable / disable build configurations, configurations comes ordered like ConfigurationNames
func SetEnabledConfigurations(configurations []bool) {
	i := 0
	enabledOne := false
	for _, name := range ConfigurationNames {
		if i < len(configurations) && configurations[i] {
			Configurations[name].SetEnabled(true)
			enabledOne = true
			log.Debug("enabled configuration " + Configurations[name].Title())
		}
		i++
	}
	if !enabledOne {
		// os, _ := sys.GetOs()
		if Configurations["release"] != nil {
			Configurations["release"].SetEnabled(true)
		}
	}

}
