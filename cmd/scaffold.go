package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Juksefantomet/gecho/internal/scaffold"
)

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold <name>",
	Short: "Generate model, routes, migrations, and query boilerplate",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		err := scaffold.Run(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ Scaffold failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(scaffoldCmd)
}
