package jieba

import (
	"fmt"
	"github.com/yanyiwu/gojieba"
	"os"
	"whisper/pkg/context"
)

var dictPath = "plugin/dict/lol_equip.dict"

func updateDict(ctx *context.Context, words []string) error {
	// 打开文件以追加写入
	file, err := os.OpenFile(dictPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return err
	}
	defer file.Close()

	// 追加写入内容
	str := ""
	for _, word := range words {
		str += word + "\n"
	}
	_, err = file.Write([]byte(str))
	if err != nil {
		return err
	}
	return nil
}

func Analyzer(ctx *context.Context, text string, keyWords map[string][]string, stopWords []string) ([]string, error) {
	dict := make([]string, 0)
	// 添加字典用于提取关键词
	for _, kws := range keyWords {
		//log.Logger.Info(ctx, "add cate keywords:", cate, kws)
		for _, keyword := range kws {
			dict = append(dict, keyword)
		}
	}

	// 更新词库
	err := updateDict(ctx, dict)
	if err != nil {
		return nil, err
	}

	// 初始化分词器
	var seg = gojieba.NewJieba(dictPath)
	defer seg.Free()

	// 分词并提取关键词
	words := seg.Cut(text, true)

	hits := make([]string, 0)
	for _, word := range words {
		if len(word) > 1 {
			if contains(&keyWords, word) {
				hits = append(hits, word)
			}
		}
	}

	return hits, nil
}

// 判断切片中是否包含某个元素
func contains(m *map[string][]string, str string) bool {
	for _, slice := range *m {
		for _, s := range slice {
			if s == str {
				return true
			}
		}
	}
	return false
}
