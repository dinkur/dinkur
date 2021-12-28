package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// outCmd represents the out command
var outCmd = &cobra.Command{
	Use:     "out",
	Aliases: []string{"o", "end"},
	Short:   "Check out/end the currently active task",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("out called")
	},
}

func init() {
	rootCmd.AddCommand(outCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// outCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// outCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}