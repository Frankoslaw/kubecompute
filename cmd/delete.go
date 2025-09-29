package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use: "delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("delete called")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

}
