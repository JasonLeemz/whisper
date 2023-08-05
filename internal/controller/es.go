package controller

import (
	"fmt"
	"whisper/internal/logic"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

type ReqQuery struct {
	KeyWords string `json:"key_words"`
	Platform int    `json:"platform,omitempty"`
	Category int    `json:"category,omitempty"`
	Map      string `json:"map,omitempty"`
}

type RespQuery struct {
	Tips  string  `json:"tips,omitempty"`
	Lists []*list `json:"lists"`
}

type list struct {
	Name        string `json:"name,omitempty"`
	IconPath    string `json:"iconPath,omitempty"`
	Price       string `json:"price,omitempty"`
	Description string `json:"description,omitempty"`
	Plaintext   string `json:"plaintext,omitempty"`
	Sell        string `json:"sell,omitempty"`
	Total       string `json:"total,omitempty"`
	Maps        string `json:"maps,omitempty"`
}

func Query(ctx *context.Context) {

	req := &ReqQuery{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	params := logic.SearchParams{
		KeyWords: req.KeyWords,
		Platform: req.Platform,
		Category: req.Category,
		Map:      req.Map,
	}

	result, err := logic.EsSearch(ctx, &params)
	if err != nil {
		ctx.Reply(nil, errors.New(err))
	}

	resp := RespQuery{}
	total := result.Total.Value
	display := len(result.Hits)
	if total != display {
		resp.Tips += fmt.Sprintf("查找到%d条记录，篇幅有限只展示%d条。请补充搜索条件以便于精确查找<br>", total, display)
	}
	for _, hit := range result.Hits {
		t := list{
			Name:        hit.Source.Name,
			IconPath:    hit.Source.IconPath,
			Price:       hit.Source.Price,
			Description: hit.Source.Description,
			Plaintext:   hit.Source.Plaintext,
			Sell:        hit.Source.Sell,
			Total:       hit.Source.Total,
			Maps:        hit.Source.Maps,
		}

		resp.Lists = append(resp.Lists, &t)
	}
	ctx.Reply(resp, nil)
}

type ReqBuild struct {
	Index string `json:"index,omitempty" binding:"-"`
}

func Build(ctx *context.Context) {

	req := &ReqBuild{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	err := logic.BuildIndex(ctx, req.Index)

	ctx.Reply(nil, errors.New(err))
}