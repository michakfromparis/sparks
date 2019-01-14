package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/configuration"
	"github.com/michaKFromParis/sparks/platform"
	"github.com/michaKFromParis/sparks/sparks"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Run:   build,
	Args:  cobra.ExactArgs(1),
	Use:   "build",
	Short: "Build a sparks product",
	Long: `
build a sparks product for the selected platforms in the selected configurations
You can compose various platforms and configurations such as:

sparks build --osx --ios

will build the product in the current directory for iOS and OSX in Release configuration

sparks build --webgl --debug $HOME/sparks/app

will build the product in $HOME/sparks/app for WebGL in Debug configuration`,
}

var enabledPlatforms []bool
var enabledConfigurations []bool

func build(cmd *cobra.Command, args []string) {
	platform.SetEnabledPlatforms(enabledPlatforms)
	configuration.SetEnabledConfigurations(enabledConfigurations)
	sparks.Init()
	sparks.Build(args[0])
	sparks.Shutdown()
}

func init_build() {

	log.Trace("build init")
	rootCmd.AddCommand(buildCmd)

	enabledPlatforms = make([]bool, len(sparks.Platforms))
	enabledConfigurations = make([]bool, len(sparks.Configurations))
	log.Tracef("registered platforms: %d", len(sparks.Platforms))
	log.Tracef("registered configurations: %d", len(sparks.Configurations))

	i := 0
	for _, name := range sparks.PlatformNames {
		p := sparks.Platforms[name]
		if p != nil && i < len(enabledPlatforms) {
			buildCmd.Flags().BoolVarP(&enabledPlatforms[i], p.Name(), p.Opt(), false, "Build "+p.Title()+" platform")
		}
		i++
	}

	i = 0
	for _, name := range sparks.ConfigurationNames {
		c := sparks.Configurations[name]
		if c != nil && i < len(enabledConfigurations) {
			buildCmd.Flags().BoolVarP(&enabledConfigurations[i], c.Name(), c.Opt(), false, "Build "+c.Title()+" configuration")
		}
		i++
	}

	buildCmd.Flags().StringVarP(&config.ProductName, "name", "", "", "Product name")
	buildCmd.Flags().StringVarP(&config.SourceDirectory, "source", "", "", "Source directory")
	buildCmd.Flags().StringVarP(&config.OutputDirectory, "output", "", "", "Output directory for all selected builds")
	buildCmd.Flags().SortFlags = false
}
