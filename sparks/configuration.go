package sparks

import log "github.com/Sirupsen/logrus"

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
	log.Info("Registering configuration: " + configuration.Name())
	Configurations[configuration.Name()] = configuration
}
