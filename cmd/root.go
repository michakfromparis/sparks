package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sparks",
	Short: "command line interface to the sparks framework",
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
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose log level")
}
