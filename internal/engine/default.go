package engine

import (
	"gord/internal/model"

	"github.com/spf13/viper"
)

type Translator interface {
	Search(word string) (model.DictResult, error)
}

func GetDefaultEngine(engineName string) Translator {
	switch engineName {
	case "dictapi":
		return NewDictAPI()
	case "baidu":
		appID := viper.GetString("api_keys.baidu_id")
		secret := viper.GetString("api_keys.baidu_secret")
		return NewBaiduEngine(appID, secret)
	default:
		return NewDictAPI()
	}
}
