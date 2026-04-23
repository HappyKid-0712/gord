package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gord/internal/model" // 确保这里的路径对应你实际的项目
)

// 1. 定义引擎本身（它负责实现 Translator 接口）
type DictAPI struct{}

// NewDictAPI 提供一个方便的实例化方法
func NewDictAPI() *DictAPI {
	return &DictAPI{}
}

// 2. 私有的数据模具（首字母小写，外部不可见，只管接数据）
type dictApiRawResponse []struct {
	Word     string `json:"word"`
	Phonetic string `json:"phonetic"`
	Meanings []struct {
		Definitions []struct {
			Definition string `json:"definition"`
		} `json:"definitions"`
	} `json:"meanings"`
}

// 3. Search 方法绑定在 DictAPI 引擎上
func (d *DictAPI) Search(word string) (model.DictResult, error) {
	url := "https://api.dictionaryapi.dev/api/v2/entries/en/" + word

	resp, err := http.Get(url)
	if err != nil {
		// 【修复】遇到错误立刻 return 空结构体和 error
		return model.DictResult{}, errors.New("网络请求失败，请检查网络连接")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// 【修复】fmt.Errorf 自带格式化，不需要嵌套 fmt.Sprintf，并且立刻 return
		return model.DictResult{}, fmt.Errorf("查询失败，状态码：%d", resp.StatusCode)
	}

	var raw dictApiRawResponse
	if errDecode := json.NewDecoder(resp.Body).Decode(&raw); errDecode != nil {
		return model.DictResult{}, errors.New("JSON 解析失败")
	}

	// 【修复】安全检查：防止查不到单词时 API 返回空数组导致 raw[0] 越界崩溃
	if len(raw) == 0 {
		return model.DictResult{}, errors.New("未找到该单词的释义")
	}

	// 开始组装标准结果
	standardResult := model.DictResult{
		Word:     raw[0].Word,
		Phonetic: raw[0].Phonetic,
		Source:   "DictionaryAPI.dev",
	}

	// 提取嵌套的释义
	for _, m := range raw[0].Meanings {
		for _, def := range m.Definitions {
			standardResult.Meanings = append(standardResult.Meanings, def.Definition)
		}
	}

	return standardResult, nil
}
