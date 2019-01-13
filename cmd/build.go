package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/platform"
	"github.com/michaKFromParis/sparks/configurations"
	"github.com/michaKFromParis/sparks/sparks"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Run:   build,
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
	platform.RegisterPlatforms()
	platform.SetEnabledPlatforms(enabledPlatforms)
	configurations.SetEnabledConfigurations(enabledConfigurations)
	sparks.Initialize()
	sparks.Build()
	sparks.Shutdown()
}

func init() {

	log.Info("QQQ")
	rootCmd.AddCommand(buildCmd)

	enabledPlatforms = make([]bool, len(sparks.Platforms))
	enabledConfigurations = make([]bool, len(sparks.Configurations))
	log.Infof("platforms: %d", len(sparks.Platforms))
	log.Infof("configurations: %d", len(sparks.Configurations))

	// i := 0
	// for _, name := range sparks.PlatformNames {
	// 	p := sparks.Platforms[name]
	// 	buildCmd.Flags().BoolVarP(&enabledPlatforms[i], p.Name(), p.Opt(), false, "Build "+p.Title()+" platform")
	// 	i++
	// }

	// i = 0
	// for _, name := range sparks.ConfigurationNames {
	// 	c := sparks.Configurations[name]
	// 	buildCmd.Flags().BoolVarP(&enabledConfigurations[i], c.Name(), c.Opt(), false, "Build "+c.Title()+" configuration")
	// 	i++
	// }

	// buildCmd.Flags().StringVarP(&config.OutputDirectory, "output", "", "", "Output directory for all selected builds")
	// buildCmd.Flags().BoolVarP(&config.Debug, "debug", "d", false, "Build debug configuration")
	// buildCmd.Flags().BoolVarP(&config.Release, "release", "r", false, "Build release configuration")
	// buildCmd.Flags().BoolVarP(&config.Shipping, "shipping", "s", false, "Build shipping configuration")
	// buildCmd.Flags().SortFlags = false
}
