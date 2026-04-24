package printer

import (
	"fmt"
	"gord/internal/model"

	"github.com/fatih/color"
)

func PrintConsole(result model.DictResult, err error) {

	errPen := color.New(color.FgHiRed, color.Bold)
	wordPen := color.New(color.FgHiCyan, color.Bold)
	phoneticPen := color.New(color.FgHiMagenta)
	bulletPen := color.New(color.FgHiYellow)
	sourcePen := color.New(color.FgHiBlack, color.Italic)

	if err != nil {
		errPen.Printf(
			"\n查询错误❌：%s、\n\n", err)
		return
	}
	fmt.Println()
	wordPen.Printf(
		"%s", result.Word)

	// 2. 打印音标 (如果有的话)
	if result.Phonetic != "" {
		phoneticPen.Printf(
			"[%s]\n", result.Phonetic)
	} else {
		fmt.Println()
	}

	// 分割线
	color.HiBlack(" ----------------------------------------")

	// 3. 遍历打印多重释义
	for _, meaning := range result.Meanings {
		bulletPen.Print("  • ")
		fmt.Printf(
			"%s\n", meaning)
	}

	fmt.Println()

	// 4. 右下角水印
	sourcePen.Printf(
		"    (Powered by %s)\n", result.Source)
	fmt.Println()
	fmt.Println() //顶部空格流出呼吸感

}
