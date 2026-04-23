package printer

import (
	"fmt"
	"gord/internal/model"
)

func PrintConsole(result model.DictResult, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result.Meanings)
	fmt.Println(result.Phonetic)
	fmt.Println(result.Source)
}
