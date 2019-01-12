// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/platform"
	"github.com/spf13/cobra"
)

var debug, release, shipping, osx, ios bool

var enabledPlatforms []bool

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: build,
}

func build(cmd *cobra.Command, args []string) {
	precmd()
	i := 0
	for _, name := range platform.PlatformNames {
		p := platform.Platforms[name]
		if enabledPlatforms[i] {
			if config.Debug {
				p.Build(platform.Debug)
			}
			if config.Release {
				p.Build(platform.Release)
			}
			if config.Shipping {
				p.Build(platform.Shipping)
			}
		}
		i++
	}

	postcmd()
}

func init() {

	rootCmd.AddCommand(buildCmd)

	enabledPlatforms = make([]bool, len(platform.Platforms))
	i := 0
	for _, name := range platform.PlatformNames {
		p := platform.Platforms[name]
		buildCmd.Flags().BoolVarP(&enabledPlatforms[i], p.Name(), p.Opt(), false, "Build "+p.Title()+" platform")
		i++
	}

	buildCmd.Flags().BoolVarP(&config.Debug, "debug", "d", false, "Build debug configuration")
	buildCmd.Flags().BoolVarP(&config.Release, "release", "r", false, "Build release configuration")
	buildCmd.Flags().BoolVarP(&config.Shipping, "shipping", "s", false, "Build shipping configuration")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
