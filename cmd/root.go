// Copyright (c) Inlets Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {

}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

// inletsCmd represents the base command when called without any sub commands.
var rootCmd = &cobra.Command{
	Use:   "lavoro",
	Short: "lavoro runs Kubernetes jobs",
	Long: `
lavoro run will run a Kubernetes job for you adn then print out the 
logs upon completion.`,
	Run: runRoot,
}

func runRoot(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func parseBaseCommand(_ *cobra.Command, _ []string) {
	os.Exit(0)
}
