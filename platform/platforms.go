package platform

import "leblox.com/sparks-cli/v2/sparks"

// RegisterPlatforms registers all existing platforms into sparks
// TODO This should be replaced by a plugin system
func RegisterPlatforms() {
	sparks.RegisterPlatform(&Osx{})
	sparks.RegisterPlatform(&Ios{})
	sparks.RegisterPlatform(&Android{})
	sparks.RegisterPlatform(&Windows{})
	sparks.RegisterPlatform(&Linux{})
	sparks.RegisterPlatform(&WebGl{})
}
