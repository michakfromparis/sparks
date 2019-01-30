package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
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
	Example: `
  sparks build
    will build the product in the current directory
    in Release configuration for the current platform

  sparks build ProductName
    will build ProductName in Release configuration for the current platform
    default Sparks Product Directory: $HOME/Sparks

  sparks build --clean --clean-assets --ios --shipping --threads=8 ProductPath
    will clean the build and the assets to build the product at ProductPath
    with a fresh assets build for iOS in Shipping configuration using 8 cores

  sparks build --no-assets --android --release --threads=8 ProductPath
    will incrementally build the product at ProductPath for Android in Release
    configuration using 8 cores`,
}

func build(cmd *cobra.Command, args []string) {
	sparks.Init()
	sparks.SetEnabledPlatforms(enabledPlatforms)
	sparks.SetEnabledConfigurations(enabledConfigurations)
	for _, sourceDirectory := range args {
		if err := sparks.Build(sourceDirectory, conf.OutputDirectory); err != nil {
			errx.Fatal(err)
		}
	}
	sparks.Shutdown()
}

func initBuild() {
	log.Trace("build init")
	rootCmd.AddCommand(buildCmd)
	log.Tracef("registered platforms: %d", len(sparks.Platforms))
	log.Tracef("registered configurations: %d", len(sparks.Configurations))
	buildCmd.Flags().SortFlags = false
	buildCmd.Flags().StringVarP(&conf.SourceDirectory, "source", "", "", "source directory to build")
	buildCmd.Flags().StringVarP(&conf.OutputDirectory, "output", "", "", "output directory for all selected builds")
	buildCmd.Flags().StringVarP(&conf.ProductName, "name", "", "", "set the product name / filename of the built binaries")
	buildCmd.Flags().BoolVarP(&conf.GenerateLua, "lua", "L", false, "generate lua bindings")
	addCommandPlatforms(buildCmd, "build")
	addCommandConfigurations(buildCmd)
}
