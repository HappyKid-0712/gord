/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gord/internal/engine"
	"gord/internal/printer"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// qCmd represents the q command
var qCmd = &cobra.Command{
	Use:   "q [word]",
	Short: "使用q来查询可能为Command的单词",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		word := args[0]

		currEngine := viper.GetString("engine")
		//获得一个translator的接口
		transLator := engine.GetDefaultEngine(currEngine)
		//通过translator实现松耦合
		result, err := transLator.Search(word)
		//输出，顺便就处理错误了
		printer.PrintConsole(result, err)

	},
}

func init() {
	rootCmd.AddCommand(qCmd)

}
