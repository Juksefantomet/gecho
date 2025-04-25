package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current Gecho version",
	Run: func(cmd *cobra.Command, args []string) {
		info, ok := debug.ReadBuildInfo()
		if !ok || info.Main.Version == "" {
			fmt.Println("Gecho (unknown version)")
			return
		}
		fmt.Printf("Gecho %s\n", info.Main.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
