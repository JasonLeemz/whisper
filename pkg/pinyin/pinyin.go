package pinyin

import (
	"github.com/mozillazg/go-pinyin"
	"strings"
)

func Trans(text string) (string, string) {
	// doc: https://pkg.go.dev/github.com/mozillazg/go-pinyin
	a := pinyin.NewArgs()
	result := strings.Join(pinyin.LazyPinyin(text, a), "")

	a.Style = pinyin.FirstLetter
	first := strings.Join(pinyin.LazyPinyin(text, a), "")

	return result, first
}
