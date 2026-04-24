package engine

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gord/internal/model"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// 定义百度引擎的用户结构
type BaiduEngine struct {
	AppID     string
	AppSecret string
}

func NewBaiduEngine(appID string, secret string) *BaiduEngine {
	return &BaiduEngine{
		AppID:     appID,
		AppSecret: secret,
	}
}

// 1. 在原来的基础上，增加一个 Dict 字段，用来接收百度返回的“高级词典数据”
type baiduRawResponse struct {
	ErrorCode   string `json:"error_code"`
	ErrorMsg    string `json:"error_msg"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
	Dict string `json:"dict"` // 这是百度隐藏的词典数据（是一个字符串格式的JSON）
}

// 2. 专门为了解析上面那个 Dict 字符串而准备的私有模具
type baiduDictResult struct {
	WordResult struct {
		SimpleMeans struct {
			Symbols []struct {
				PhEn  string `json:"ph_en"` // 英式音标
				PhAm  string `json:"ph_am"` // 美式音标
				Parts []struct {
					Part  string   `json:"part"`  // 词性 (比如 n.)
					Means []string `json:"means"` // 释义数组 (比如 ["苹果", "苹果树"])
				} `json:"parts"`
			} `json:"symbols"`
		} `json:"simple_means"`
	} `json:"word_result"`
}

func (b *BaiduEngine) Search(word string) (model.DictResult, error) {
	salt := strconv.FormatInt(time.Now().Unix(), 10)
	signStr := b.AppID + word + salt + b.AppSecret
	hash := md5.Sum([]byte(signStr))
	sign := hex.EncodeToString(hash[:])

	params := url.Values{}
	params.Add("q", word)
	params.Add("from", "auto")
	params.Add("to", "zh")
	params.Add("appid", b.AppID)
	params.Add("salt", salt)
	params.Add("sign", sign)
	params.Add("dict", "1") // 👈 核心魔法：加上这个参数，强制要求百度返回词典和音标数据！

	apiURL := "http://api.fanyi.baidu.com/api/trans/vip/translate?" + params.Encode()

	resp, err := http.Get(apiURL)
	if err != nil {
		return model.DictResult{}, fmt.Errorf("网络请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var raw baiduRawResponse

	if err := json.Unmarshal(body, &raw); err != nil {
		return model.DictResult{}, fmt.Errorf("解析数据失败")
	}

	if raw.ErrorCode != "" && raw.ErrorCode != "52000" {
		return model.DictResult{}, fmt.Errorf("百度 API 报错, 错误码: %s, 信息: %s", raw.ErrorCode, raw.ErrorMsg)
	}

	if len(raw.TransResult) == 0 {
		return model.DictResult{}, fmt.Errorf("未找到翻译结果")
	}

	// 3. 先给一个保底的单句翻译结果（如果你查的是很长的句子，它就没有词典数据，只会走这里）
	result := model.DictResult{
		Word:     raw.TransResult[0].Src,
		Meanings: []string{raw.TransResult[0].Dst}, // 初始只有 1 个结果
		Source:   "Baidu Translate API",
	}

	// 4. 如果我们拿到了 Dict 数据（说明查的是单词），就用多重释义覆盖掉上面那个的结果
	if raw.Dict != "" {
		var dictData baiduDictResult
		// 二次解析那个字符串
		if err := json.Unmarshal([]byte(raw.Dict), &dictData); err == nil {
			symbols := dictData.WordResult.SimpleMeans.Symbols
			if len(symbols) > 0 {

				// 获取音标
				if symbols[0].PhEn != "" {
					result.Phonetic = symbols[0].PhEn
				} else if symbols[0].PhAm != "" {
					result.Phonetic = symbols[0].PhAm
				}

				// 获取多重词性和释义
				if len(symbols[0].Parts) > 0 {
					// 准备一个新的空切片
					var detailedMeanings []string

					for _, p := range symbols[0].Parts {
						// 拼接词性，比如 "n. "
						meansStr := p.Part + " "
						for i, m := range p.Means {
							if i > 0 {
								meansStr += ", "
							}
							meansStr += m
						}

						detailedMeanings = append(detailedMeanings, meansStr)
					}

					// 用我们收集到的“多重释义”，强行覆盖掉原先那个只有 1 个结果的切片！
					result.Meanings = detailedMeanings
				}
			}
		}
	}

	return result, nil
}
