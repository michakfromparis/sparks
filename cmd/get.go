package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/platform"
	"github.com/michaKFromParis/sparks/sparks"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Run:   get,
	Use:   "get",
	Short: "Get modules and system dependencies",
	Long:  `Get modules and system dependencies`,
}

func get(cmd *cobra.Command, args []string) {
	sparks.Init()
	platform.SetEnabledPlatforms(enabledPlatforms)
	if err := sparks.Get(); err != nil {
		errx.Fatal(err)
	}
	sparks.Shutdown()
}

func init_get() {
	enabledPlatforms = make([]bool, len(sparks.Platforms))
	enabledConfigurations = make([]bool, len(sparks.Configurations))
	log.Tracef("registered platforms: %d", len(sparks.Platforms))
	getCmd.Flags().SortFlags = false
	getCmd.Flags().BoolVarP(&config.GetDependencies, "dependencies", "d", true, "get system dependencies")
	addPlatforms(getCmd, "get")
	rootCmd.AddCommand(getCmd)
}
