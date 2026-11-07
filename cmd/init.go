package cmd

import (
	"github.com/Juksefantomet/gecho/internal/setup"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Gecho-compatible project",
	Run: func(cmd *cobra.Command, args []string) {
		setup.RunInit()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
