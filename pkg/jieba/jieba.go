package jieba

import (
	"fmt"
	"github.com/yanyiwu/gojieba"
)

func Analyzer(text string) {
	// 初始化分词器
	var seg = gojieba.NewJieba()
	defer seg.Free()

	// 分词并提取关键词
	words := seg.Cut(text, true)
	keywords := make(map[string]int)

	// 过滤掉停用词等不需要的词汇
	stopWords := []string{"的", "之"}
	for _, word := range words {
		if len(word) > 1 && !contains(stopWords, word) {
			keywords[word]++
		}
	}

	// 输出提取的关键词和出现次数
	fmt.Println("提取的关键词：")
	for keyword, count := range keywords {
		fmt.Printf("%s: %d\n", keyword, count)
	}
}

// 判断切片中是否包含某个元素
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}