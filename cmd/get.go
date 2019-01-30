package cmd

import (
	log "github.com/Sirupsen/logrus"
	conf "github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
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
	sparks.SetEnabledPlatforms(enabledPlatforms)
	if err := sparks.Get(); err != nil {
		errx.Fatal(err)
	}
	sparks.Shutdown()
}

func initGet() {
	log.Tracef("registered platforms: %d", len(sparks.Platforms))
	getCmd.Flags().SortFlags = false
	getCmd.Flags().BoolVarP(&conf.GetDependencies, "dependencies", "d", true, "get system dependencies")
	addCommandPlatforms(getCmd, "get")
	rootCmd.AddCommand(getCmd)
}
