package spider

import (
	"encoding/json"
	"whisper/internal/dto"
	"whisper/pkg/context"
	"whisper/pkg/http"
	"whisper/pkg/log"
)

type Bilibili struct {
	ts            int64
	referer       []http.Header
	commonHeaders []http.Header
}

func (b *Bilibili) SearchKeywords(ctx *context.Context, url string) (interface{}, error) {

	// 发送 GetForm 请求
	sk := dto.SearchKeywords{}

	body, err := http.GetForm(ctx, url, b.commonHeaders...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &sk)
	if err != nil {
		log.Logger.Error(ctx, string(body), err)
	}
	return &sk, err
}
