/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// qCmd represents the q command
var qCmd = &cobra.Command{
	Use:   "q [word]",
	Short: "使用q来查询可能为Command的单词",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("q called")
	},
}

func init() {
	rootCmd.AddCommand(qCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// qCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// qCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
