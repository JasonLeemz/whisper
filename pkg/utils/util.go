package utils

import (
	"fmt"
	"github.com/spf13/cast"
	"regexp"
	"strconv"
	"strings"
)

func CompareVersion(version1 string, version2 string) int {
	v1 := strings.Split(version1, ".")
	v2 := strings.Split(version2, ".")
	l1 := len(v1)
	l2 := len(v2)
	maxLen := l1
	if l2 > l1 {
		maxLen = l2
	}

	v1Arr := make([]int, 0)
	v2Arr := make([]int, 0)
	// 补位&修订号转换成数字，补到一样的长度
	for i := 0; i < maxLen; i++ {
		if i > l1-1 {
			v1Arr = append(v1Arr, 0)
		} else {
			cv, _ := strconv.Atoi(v1[i])
			v1Arr = append(v1Arr, cv)
		}

		if i > l2-1 {
			v2Arr = append(v2Arr, 0)
		} else {
			cv, _ := strconv.Atoi(v2[i])
			v2Arr = append(v2Arr, cv)
		}
	}

	for j := 0; j < maxLen; j++ {
		if v1Arr[j] == v2Arr[j] {
			continue
		} else if v1Arr[j] > v2Arr[j] {
			return 1
		} else if v1Arr[j] < v2Arr[j] {
			return -1
		}
	}

	return 0
}

func ExtractKeywords(text string, re *regexp.Regexp) []string {
	matches := re.FindAllString(text, -1)

	uniqueMatches := make(map[string]bool)

	for _, match := range matches {
		uniqueMatches[match] = true
	}

	result := make([]string, 0, len(uniqueMatches))
	for match := range uniqueMatches {
		result = append(result, match)
	}

	return result
}

func CompileKeywordsRegex(keywords []string) *regexp.Regexp {
	// 使用或(|)运算符连接所有关键词，并使用转义字符\进行转义
	regexPattern := fmt.Sprintf("(%s)", strings.Join(keywords, "|"))

	// 使用正则表达式编译正则字符串
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		panic(err)
	}

	return re
}

func RemoveRepeatedBRTag(text string) string {
	re := regexp.MustCompile(`(?i)(?:(?:\r\n|\r|\n|\f|\x85)+)`)
	text = re.ReplaceAllString(text, "<br>")
	re = regexp.MustCompile(`(?i)(?:(?:<br>)+)`)
	text = re.ReplaceAllString(text, "<br>")

	return text
}

// Str2Int 这个不是将小数转int！！！
func Str2Int(text string) int {
	// 0.001234

	i := 0
	for _, s := range text {
		if s == '0' || s == '.' {
			i++
			continue
		}
		break
	}

	ns := text[i:]
	return cast.ToInt(ns)
}
