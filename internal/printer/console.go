package printer

import (
	"errors"
	"fmt"
	"gord/internal/engine"
)

func PrintConsole(word string) error {
	result, err := engine.NewDictAPI().Search(word) //这里其实是越权了
	if result.Word == "" {
		err = errors.New("查询失败！")
		fmt.Println(err)
		return err
	}
	fmt.Println(result.Meanings)
	fmt.Println(result.Phonetic)
	fmt.Printf("翻译引擎为：%s", result.Source)
	return nil
}
