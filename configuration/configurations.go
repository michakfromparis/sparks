package configuration

import "github.com/michaKFromParis/sparks/sparks"

func RegisterConfigurations() {
	sparks.RegisterConfiguration(Debug{})
	sparks.RegisterConfiguration(Release{})
	sparks.RegisterConfiguration(Shipping{})
}

func SetEnabledConfigurations(enabledConfigurations []bool) {
	i := 0
	for _, name := range sparks.ConfigurationNames {
		if i < len(enabledConfigurations) && enabledConfigurations[i] {
			sparks.Configurations[name].SetEnabled(true)
		}
		i++
	}
}
