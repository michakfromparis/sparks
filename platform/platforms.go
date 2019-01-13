package platform

import "github.com/michaKFromParis/sparks/sparks"

func RegisterPlatforms() {
	sparks.RegisterPlatform(Osx{})
	sparks.RegisterPlatform(Ios{})
	sparks.RegisterPlatform(WebGl{})
}

func SetEnabledPlatforms(enabledPlatforms []bool) {
	i := 0
	for _, name := range sparks.PlatformNames {
		if i < len(enabledPlatforms) && enabledPlatforms[i] {
			sparks.Platforms[name].SetEnabled(true)
		}
		i++
	}
}
