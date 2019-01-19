package configuration

import (
	"github.com/michaKFromParis/sparks/sparks"
)

// RegisterConfigurations registers all existing configurations into sparks
// TODO This should be replaced by a plugin system
func RegisterConfigurations() {
	sparks.RegisterConfiguration(&Debug{})
	sparks.RegisterConfiguration(&Release{})
	sparks.RegisterConfiguration(&Shipping{})
}
