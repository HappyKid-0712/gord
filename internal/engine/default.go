package engine

import (
	"gord/internal/model"
)

type Translator interface {
	Search(word string) (model.DictResult, error)
}

func GetDefaultEngine(engineName string) Translator {
	switch engineName {
	case "dictapi":
		return NewDictAPI()
	case "youdao":
		return NewDictAPI()
	default:
		return NewDictAPI()
	}
}
