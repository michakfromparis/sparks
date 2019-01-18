package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var depsCmd = &cobra.Command{
	Run:   deps,
	Use:   "deps",
	Short: "Install system dependencies",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func deps(cmd *cobra.Command, args []string) {
	log.Info("installing dependencies")
}

func init_deps() {
	log.Trace("deps init")
	rootCmd.AddCommand(depsCmd)
}
