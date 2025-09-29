package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var nukeCmd = &cobra.Command{
	Use: "nuke",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("nuke called")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(nukeCmd)

}
