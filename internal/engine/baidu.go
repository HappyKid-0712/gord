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

type baiduRawResponse struct {
	ErrorCode   string `json:"error_code"`
	ErrorMsg    string `json:"error_msg"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
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

	if raw.ErrorCode != "" && raw.ErrorCode != "52000" { // 52000 代表成功
		return model.DictResult{}, fmt.Errorf("百度 API 报错, 错误码: %s, 信息: %s", raw.ErrorCode, raw.ErrorMsg)
	}

	if len(raw.TransResult) == 0 {
		return model.DictResult{}, fmt.Errorf("未找到翻译结果")
	}

	result := model.DictResult{
		Word:     raw.TransResult[0].Src,
		Meanings: []string{raw.TransResult[0].Dst}, // 百度标准版默认返回翻译句子，作为一条释义
		Source:   "Baidu Translate API",
	}

	return result, nil
}
