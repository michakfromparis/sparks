package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/logger"
	"github.com/michaKFromParis/sparks/sparks"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	PersistentPreRun: PreRunAllCommands,
	Use:              "sparks",
	Short:            "command line interface to the sparks framework",
	Long:             `Sparks command line interface`,
}

// PreRunAllCommands is called by cobra for every commands just before executing a command
func PreRunAllCommands(cmd *cobra.Command, args []string) {
	// reinitialized here to take the --verbose command line flag into account
	logger.Init()
}

// SetVersion passes the current version to the cli
func SetVersion(version string, commit string, date string) {
	rootCmd.Version = version
}

// Execute is the main blocking entry point, usually called from main
// it chains all the commands into a tree and execcutes all sub commands
func Execute() error {
	rootCmd.SilenceErrors = true
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

// Init is calling all cobra command functions here to avoid a race condition
// with the loading of the sparks platforms / configurations first
func Init() {
	enabledPlatforms = make([]bool, len(sparks.Platforms))
	enabledConfigurations = make([]bool, len(sparks.Configurations))
	initRoot()
	initGet()
	initClean()
	initBuild()
	initCode()
	initRun()
}

func initRoot() {
	log.Trace("root init")
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "set log verbose level to debug")
	rootCmd.PersistentFlags().BoolVarP(&config.VeryVerbose, "v", "", false, "set log verbose level to trace")
}
