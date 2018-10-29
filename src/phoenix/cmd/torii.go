package cmd

import (
	"../torii"
	"github.com/spf13/cobra"
)

// toriiCmd represents the torii command
var toriiCmd = &cobra.Command{
	Use:   "torii",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return torii.Run()
	},
}

func init() {
	rootCmd.AddCommand(toriiCmd)
}
