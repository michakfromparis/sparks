package cmd

import (
	"github.com/michakfromparis/sparks/sparks"
	"github.com/spf13/cobra"
)

var enabledPlatforms []bool
var enabledConfigurations []bool

func addCommandPlatforms(cmd *cobra.Command, prefix string) {
	i := 0
	for _, name := range sparks.PlatformNames {
		p := sparks.Platforms[name]
		if p != nil && i < len(enabledPlatforms) {
			cmd.Flags().BoolVarP(&enabledPlatforms[i], p.Name(), p.Opt(), false, prefix+" for "+p.Title()+"")
		}
		i++
	}
}

func addCommandConfigurations(cmd *cobra.Command) {
	i := 0
	for _, name := range sparks.ConfigurationNames {
		c := sparks.Configurations[name]
		if c != nil && i < len(enabledConfigurations) {
			cmd.Flags().BoolVarP(&enabledConfigurations[i], c.Name(), c.Opt(), false, "build in "+c.Title()+" configuration")
		}
		i++
	}
}
