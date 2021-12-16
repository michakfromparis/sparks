package cmd

import (
	"os"

	"github.com/michakfromparis/sparks/errx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"leblox.com/sparks-cli/v2/conf"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean a sparks product",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: clean,
}

func clean(cmd *cobra.Command, args []string) {

	log.Info("Cleaning")
	if err := os.RemoveAll(conf.OutputDirectory); err != nil {
		errx.Fatalf(err, "failed to clean")
	}
}

func initClean() {
	log.Trace("clean init")
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
