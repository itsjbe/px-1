package cmd

import (
	"../wph"
	"github.com/spf13/cobra"
)

var wphCmd = &cobra.Command{
	Use:   "wph",
	Short: "A brief description of wph command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return wph.Run()
	},
}

func init() {
	rootCmd.AddCommand(wphCmd)
}
