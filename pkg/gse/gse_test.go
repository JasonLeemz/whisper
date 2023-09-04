package gse

import (
	"fmt"
	"testing"
)

func TestGse(t *testing.T) {
	text := "我是拖拉机学院手扶拖拉机专业的。不用多久，我就会升职加薪，当上CEO，走上人生巅峰。"
	aa := Analyzer(text)
	fmt.Println(aa)
}
