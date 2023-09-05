package test

import (
	"fmt"
	"testing"
	"whisper/pkg/utils"
)

func TestCompareVersion(t *testing.T) {
	res := utils.CompareVersion("1.2a", "1.2b")

	fmt.Println("result:", res)

}

func TestStr2Int(t *testing.T) {
	text := "123456"
	i := utils.Str2Int(text)
	fmt.Println(i)

}
