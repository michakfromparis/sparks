package configuration

import (
	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/sparks"
)

func RegisterConfigurations() {
	sparks.RegisterConfiguration(&Debug{})
	sparks.RegisterConfiguration(&Release{})
	sparks.RegisterConfiguration(&Shipping{})
}

func SetEnabledConfigurations(configurations []bool) {
	i := 0
	for _, name := range sparks.ConfigurationNames {
		if i < len(configurations) && configurations[i] == true {
			sparks.Configurations[name].SetEnabled(true)
			log.Debug("enabled configuration " + sparks.Configurations[name].Title())
		}
		i++
	}
}
