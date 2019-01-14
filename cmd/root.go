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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Init() {
	init_root()
	init_deps()
	init_clean()
	init_build()
}

func PreRunAllCommands(cmd *cobra.Command, args []string) {
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
	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "verbose log level. One of (panic, fatal, error, warning, debug, trace)")
}
