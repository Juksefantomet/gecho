package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current Gecho version",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile("VERSION")
		if err != nil {
			fmt.Println("✗ VERSION file not found")
			os.Exit(1)
		}
		fmt.Printf("Gecho v%s\n", string(data))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
