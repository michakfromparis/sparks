package configuration

import (
	"github.com/michaKFromParis/sparks/sparks"
)

func RegisterConfigurations() {
	sparks.RegisterConfiguration(&Debug{})
	sparks.RegisterConfiguration(&Release{})
	sparks.RegisterConfiguration(&Shipping{})
}
