package cmd

import (
	"github.com/spf13/cobra"
	"leblox.com/sparks-cli/v2/conf"
	"leblox.com/sparks-cli/v2/errx"
	"leblox.com/sparks-cli/v2/sparks"
)

// codeCmd represents the code command
var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Open the generated project",
	Long:  "Open the generated project",
	Run: func(cmd *cobra.Command, args []string) {
		sparks.Init()
		sparks.SetEnabledPlatforms(enabledPlatforms)
		sparks.SetEnabledConfigurations(enabledConfigurations)
		for _, sourceDirectory := range args {
			if err := sparks.Code(sourceDirectory, conf.OutputDirectory); err != nil {
				errx.Fatal(err)
			}
		}
		sparks.Shutdown()
	},
}

func initCode() {
	rootCmd.AddCommand(codeCmd)
	addCommandPlatforms(codeCmd, "code")
	addCommandConfigurations(codeCmd)
}
