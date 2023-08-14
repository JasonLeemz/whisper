package utils

import (
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
