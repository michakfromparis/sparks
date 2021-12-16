package cmd

import (
	"github.com/michakfromparis/sparks/errx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"leblox.com/sparks-cli/v2/conf"
	"leblox.com/sparks-cli/v2/sparks"
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
