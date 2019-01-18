package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/configuration"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/platform"
	"github.com/michaKFromParis/sparks/sparks"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Run:   build,
	Args:  cobra.ArbitraryArgs,
	Use:   "build",
	Short: "Build a sparks product",
	Long: `
build a sparks product for the selected platforms in the selected configurations
You can compose various platforms and configurations such as:

sparks build --osx --ios

will build the product in the current directory for iOS and OSX in Release configuration

sparks build --webgl --shipping $HOME/sparks/app

will build the product in $HOME/sparks/app for WebGL in Debug configuration`,
}

func build(cmd *cobra.Command, args []string) {
	sparks.Init()
	platform.SetEnabledPlatforms(enabledPlatforms)
	configuration.SetEnabledConfigurations(enabledConfigurations)
	for _, sourceDirectory := range args {
		if err := sparks.Build(sourceDirectory, config.OutputDirectory); err != nil {
			errx.Fatal(err)
		}
	}
	sparks.Shutdown()
}

func init_build() {

	log.Trace("build init")
	rootCmd.AddCommand(buildCmd)

	enabledPlatforms = make([]bool, len(sparks.Platforms))
	enabledConfigurations = make([]bool, len(sparks.Configurations))
	log.Tracef("registered platforms: %d", len(sparks.Platforms))
	log.Tracef("registered configurations: %d", len(sparks.Configurations))
	buildCmd.Flags().SortFlags = false
	buildCmd.Flags().StringVarP(&config.SourceDirectory, "source", "", "", "source directory to build")
	buildCmd.Flags().StringVarP(&config.OutputDirectory, "output", "", "", "output directory for all selected builds")
	buildCmd.Flags().StringVarP(&config.ProductName, "name", "", "", "set the product name / filename of the built binaries")
	buildCmd.Flags().BoolVarP(&config.GenerateLua, "lua", "L", false, "generate lua bindings")
	addPlatforms(buildCmd, "build")
	addConfigurations(buildCmd)
}
