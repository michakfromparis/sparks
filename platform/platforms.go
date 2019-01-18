package platform

import (
	"github.com/michaKFromParis/sparks/sparks"
)

func RegisterPlatforms() {
	sparks.RegisterPlatform(&Osx{})
	sparks.RegisterPlatform(&Ios{})
	sparks.RegisterPlatform(&Android{})
	sparks.RegisterPlatform(&Windows{})
	sparks.RegisterPlatform(&Linux{})
	sparks.RegisterPlatform(&WebGl{})
}
