package controller

import (
	"fmt"
	"whisper/internal/logic"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

type ReqQuery struct {
	KeyWords string   `json:"key_words" form:"key_words"`
	Platform string   `json:"platform,omitempty" form:"platform,omitempty"`
	Category string   `json:"category,omitempty" form:"category,omitempty"`
	Way      []string `json:"way,omitempty" form:"way,omitempty"`
	Map      []string `json:"map,omitempty" form:"map,omitempty"`
}

type RespQuery struct {
	Tips  string  `json:"tips,omitempty"`
	Lists []*list `json:"lists"`
}

type list struct {
	Name        string   `json:"name"`
	IconPath    string   `json:"iconPath"`
	Description string   `json:"description"`
	Plaintext   string   `json:"plaintext"`
	Tags        []string `json:"tags"`
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
		Way:      req.Way,
		Map:      req.Map,
	}

	result, err := logic.EsSearch(ctx, &params)
	if err != nil {
		ctx.Reply(nil, errors.New(err))
	}

	resp := RespQuery{}
	total := result.Total.Value
	display := len(result.Hits)
	resp.Tips = fmt.Sprintf("为您找到相关结果约%d个", total)
	if total != display {
		resp.Tips += fmt.Sprintf(",篇幅有限只展示%d条", display)
	}
	for _, hit := range result.Hits {
		t := list{
			Name:        hit.Source.Name,
			IconPath:    hit.Source.IconPath,
			Description: hit.Source.Description,
			Plaintext:   hit.Source.Plaintext,
			Tags:        hit.Source.Tags,
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
