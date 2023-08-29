package controller

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"strconv"
	"whisper/internal/dto"
	"whisper/internal/logic"
	"whisper/internal/service/mq"
	"whisper/pkg/context"
	"whisper/pkg/errors"
	"whisper/pkg/log"
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
	ID          interface{} `json:"id"`
	Name        string      `json:"name"`
	IconPath    string      `json:"iconPath"`
	Tags        []string    `json:"tags"`
	Description string      `json:"description"`
	Plaintext   string      `json:"plaintext"`
	Version     string      `json:"version"`
	Platform    int         `json:"platform"`
	ItemId      string      `json:"itemId"`
	Maps        string      `json:"maps"`
}

func Query(ctx *context.Context) {

	req := &ReqQuery{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	defer func() {
		mq.ProduceMessage(mq.Exchange, mq.RoutingKeySearchKey, []byte(req.KeyWords))
	}()

	way := make([]string, 0, 4)
	for _, w := range req.Way {
		if w == "name" {
			way = append(way, "name", "keywords")
		}

		if w == "description" {
			way = append(way, "description", "plaintext")
		}
	}

	params := logic.SearchParams{
		KeyWords: req.KeyWords,
		Platform: req.Platform,
		Category: req.Category,
		Way:      way,
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
		if name, ok := hit.Highlight["name"]; ok {
			if len(name) > 0 {
				hit.Source.Name = name[0]
			}
		}
		//if desc, ok := hit.Highlight["description"]; ok {
		//	if len(desc) > 0 {
		//		hit.Source.Description = desc[0]
		//	}
		//}
		t := list{
			ID:          hit.Id,
			ItemId:      hit.Source.ItemId,
			Name:        hit.Source.Name,
			IconPath:    hit.Source.IconPath,
			Description: prettyHeroDesc(ctx, hit.Source.Description, req.Platform, req.Category),
			Plaintext:   hit.Source.Plaintext,
			Tags:        hit.Source.Tags,
			Version:     hit.Source.Version,
			Platform:    cast.ToInt(hit.Source.Platform),
			Maps:        hit.Source.Maps,
		}

		resp.Lists = append(resp.Lists, &t)
	}
	ctx.Reply(resp, nil)
}

type ReqBuild struct {
	Index   string `json:"index,omitempty" binding:"-"`
	ReBuild bool   `json:"rebuild,omitempty"`
}

func Build(ctx *context.Context) {

	req := &ReqBuild{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	err := logic.BuildIndex(ctx, req.Index, req.ReBuild)

	ctx.Reply(nil, errors.New(err))
}

func prettyHeroDesc(ctx *context.Context, desc, platform, category string) string {
	if desc == "" {
		return ""
	}

	if platform != "0" && platform != "1" {
		return desc
	}

	if category != "lol_heroes" {
		return desc
	}

	//pretty := ""
	heroDesc := make([]*dto.HeroDescription, 0, 5)
	err := json.Unmarshal([]byte(desc), &heroDesc)
	if err != nil {
		return desc
	}

	mDesc := make(map[int]*dto.HeroDescription)
	for _, d := range heroDesc {
		mDesc[d.Sort] = d
	}
	sDesc := ""

	for i := 0; i < 5; i++ {
		if _, ok := mDesc[i]; !ok {
			log.Logger.Warn(ctx, "hero desc error", mDesc)
			return desc
		}
		sk := mDesc[i].SpellKey
		if i > 0 {
			if platform == "0" {
				sk = mDesc[i].SpellKey
			} else if platform == "1" {
				sk = strconv.Itoa(i)
			}

		}

		sDetail := ""
		if mDesc[i].Detail != "" {
			list := make([]string, 0, 4)
			err := json.Unmarshal([]byte(mDesc[i].Detail), &list)
			if err != nil {
				log.Logger.Error(ctx, "Unmarshal mDesc[i].Detail fail:", err)
			}

			for _, row := range list {
				sDetail += fmt.Sprintf(detailTPL, row)
			}
		}

		sDesc += fmt.Sprintf(descTPL,
			mDesc[i].AbilityIconPath,
			mDesc[i].Name,
			sk,
			sDetail,
			mDesc[i].Description,
		)
	}

	return sDesc
}

const descTPL = `
<ul>
	<li>
		<img src="%s" />
		<h6>%s</h6>
		<span>%s</span>
<component>
<a-space>
		<a-tooltip placement="topLeft" title="%s">
		<div>%s</div>
		</a-tooltip>
</a-space>
</component>

	</li>
</ul>
`
const detailTPL = `
%s
`
