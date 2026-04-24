/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gord/internal/engine"
	"gord/internal/printer"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gord [word]",
	Short: "查询一个你不认识的单词吧！grod为你服务！",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("engine", "e", "baidu", "选择翻译引擎（可选：dictapi，baidu）")
	viper.BindPFlag("engine", rootCmd.PersistentFlags().Lookup("engine"))

	viper.SetConfigName(".gord")
	viper.SetConfigType("yaml")
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(home)
	}
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("配置文件读取状态:", err)
	}
}
