package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gord/internal/model"
)

type DictAPI struct{}

func NewDictAPI() *DictAPI {
	return &DictAPI{}
}

type dictApiRawResponse []struct {
	Word     string `json:"word"`
	Phonetic string `json:"phonetic"`
	Meanings []struct {
		Definitions []struct {
			Definition string `json:"definition"`
		} `json:"definitions"`
	} `json:"meanings"`
}

func (d *DictAPI) Search(word string) (model.DictResult, error) {
	url := "https://api.dictionaryapi.dev/api/v2/entries/en/" + word

	resp, err := http.Get(url)
	if err != nil {
		return model.DictResult{}, errors.New("网络请求失败，请检查网络连接")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.DictResult{}, fmt.Errorf("查询失败，状态码：%d", resp.StatusCode)
	}

	var raw dictApiRawResponse
	if errDecode := json.NewDecoder(resp.Body).Decode(&raw); errDecode != nil {
		return model.DictResult{}, errors.New("JSON 解析失败")
	}

	if len(raw) == 0 {
		return model.DictResult{}, errors.New("未找到该单词的释义")
	}

	standardResult := model.DictResult{
		Word:     raw[0].Word,
		Phonetic: raw[0].Phonetic,
		Source:   "DictionaryAPI.dev",
	}

	for _, m := range raw[0].Meanings {
		for _, def := range m.Definitions {
			standardResult.Meanings = append(standardResult.Meanings, def.Definition)
		}
	}

	return standardResult, nil
}
