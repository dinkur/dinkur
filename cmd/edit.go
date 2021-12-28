/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"e"},
	Short:   "Edit a task",
	Long: `Applies changes to the currently active task, or the latest task, or
a specific task using the --id or -i flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("edit called")
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}