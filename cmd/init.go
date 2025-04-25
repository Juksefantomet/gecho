package cmd

import (
	"github.com/spf13/cobra"
	"github.com/Juksefantomet/gecho/internal/setup"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Gecho-compatible project",
	Run: func(cmd *cobra.Command, args []string) {
		setup.RunInit()
	},
}
