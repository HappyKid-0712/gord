package model

type DictResult struct {
	Word     string
	Phonetic string
	Meanings []string
	Source   string
}
type Translator interface {
	Search(word string) (DictResult, error)
}
