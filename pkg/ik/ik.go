package ik

import (
	"fmt"
	"whisper/pkg/context"
	"whisper/pkg/es"
)

func Analyze(ctx *context.Context, text string) ([]string, error) {
	// 创建一个用于分词的请求
	res, err := es.ESClient.IndexAnalyze().Analyzer("ik_smart").Text(text).Do(ctx)
	if err != nil {
		return nil, err
	}
	if len(res.Tokens) != 3 {
		return nil, err
	}

	fmt.Println(res)
	return nil, nil
}
