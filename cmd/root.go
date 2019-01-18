package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	PersistentPreRun: PreRunAllCommands,
	Use:              "sparks",
	Short:            "command line interface to the sparks framework",
	Long: `
Sparks command line interface`,
}

// calling all cobra command init functions here to avoid a race condition
// with the loading of the sparks platforms / configurations first
func Init() {
	init_root()
	init_get()
	init_clean()
	init_build()
}

func PreRunAllCommands(cmd *cobra.Command, args []string) {
	// reinitialized here to take the --verbose command line flag into account
	logger.Init()
}

func Execute() error {
	rootCmd.SilenceErrors = true
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init_root() {
	log.Trace("root init")
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "set log verbose level to debug")
	rootCmd.PersistentFlags().BoolVarP(&config.VeryVerbose, "v", "", false, "set log verbose level to trace")
}
