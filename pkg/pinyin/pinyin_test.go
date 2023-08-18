package pinyin

import (
	"fmt"
	"testing"
)

func TestPY(t *testing.T) {
	text := "国际化大都市"
	fmt.Println(Trans(text))

	text2 := "guojihuadadushi"
	fmt.Println(Trans(text2))

	text3 := "Hello World"
	fmt.Println(Trans(text3))
}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		text := "国际化大都市"
		fmt.Println(Trans(text))
	}
}
