package controller

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"whisper/internal/dto"
	"whisper/internal/logic"
	"whisper/internal/service/mq"
	"whisper/pkg/context"
	"whisper/pkg/errors"
	"whisper/pkg/log"
)

type ReqQuery struct {
	KeyWords string `json:"key_words" form:"key_words"`
	Platform string `json:"platform,omitempty" form:"platform,omitempty"`
	Category string `json:"category,omitempty" form:"category,omitempty"`
	Switch   struct {
		Way struct {
			Name        bool `json:"name"`
			Description bool `json:"description"`
		} `json:"way"`
		Maps struct {
			M5v5 bool `json:"_5V5"`
			Mdld bool `json:"_dld"`
			M2v2 bool `json:"_2v2"`
		} `json:"maps"`
	} `json:"switch,omitempty" form:"switch,omitempty"`
	Map []string `json:"map,omitempty" form:"map,omitempty"`
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
	if req.Switch.Way.Name {
		way = append(way, "name", "keywords")
	}
	if req.Switch.Way.Description {
		way = append(way, "description", "plaintext")
	}

	maps := make([]string, 0, 3)
	if req.Switch.Maps.M5v5 {
		maps = append(maps, "召唤师峡谷")
	}
	if req.Switch.Maps.Mdld {
		maps = append(maps, "嚎哭深渊")
	}
	if req.Switch.Maps.M2v2 {
		maps = append(maps, "斗魂竞技场")
	}

	params := logic.SearchParams{
		KeyWords: req.KeyWords,
		Platform: req.Platform,
		Category: req.Category,
		Way:      way,
		Map:      maps,
	}

	result, err := logic.EsSearch(ctx, &params)
	if err != nil {
		ctx.Reply(nil, errors.New(err))
	}

	resp := dto.SearchResult{}
	total := result.Total.Value
	display := len(result.Hits)
	resp.Tips = fmt.Sprintf("为您找到相关结果约%d个", total)
	if total != display {
		resp.Tips += fmt.Sprintf(",篇幅有限只展示%d条", display)
	}
	for _, hit := range result.Hits {

		t := dto.SearchResultList{
			Id:        hit.Source.ID,
			Name:      hit.Source.Name,
			Icon:      hit.Source.IconPath,
			MainImg:   hit.Source.MainImg,
			Desc:      hit.Source.Description,
			Plaintext: hit.Source.Plaintext,
			Tags:      hit.Source.Tags,
			Version:   hit.Source.Version,
			Platform:  cast.ToInt(hit.Source.Platform),
			Maps:      hit.Source.Maps,
			Spell:     prettyHeroDesc(ctx, hit.Source.Description, req.Platform, req.Category),
		}

		resp.List = append(resp.List, &t)
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

func prettyHeroDesc(ctx *context.Context, desc, platform, category string) []*dto.HeroSpell {
	if desc == "" || category != "lol_heroes" {
		return nil
	}

	//pretty := ""
	heroDesc := make([]*dto.HeroDescription, 0, 5)
	err := json.Unmarshal([]byte(desc), &heroDesc)
	if err != nil {
		log.Logger.Error(ctx, err)
		return nil
	}

	mDesc := make(map[int]*dto.HeroDescription)
	for _, d := range heroDesc {
		mDesc[d.Sort] = d
	}

	l := len(mDesc)
	sDesc := make([]*dto.HeroSpell, 0, l)

	for i := 0; i < l; i++ {
		if _, ok := mDesc[i]; !ok {
			log.Logger.Warn(ctx, "hero desc error", mDesc)
			return nil
		}
		sk := mDesc[i].SpellKey
		if i > 0 {
			// 第0项全部都是passive,从第1项往后:手游是1,2,3,4 端游是q,w,e,r
			if platform == "0" {
				sk = mDesc[i].SpellKey // q,w,e,r
			} else if platform == "1" {
				sk = strconv.Itoa(i) // 1,2,3,4
			}
		}

		spell := make([]string, 0, 4)
		err := json.Unmarshal([]byte(mDesc[i].Detail), &spell)
		if err != nil {
			log.Logger.Error(ctx, "Unmarshal mDesc[i].Detail fail:", err)
		}

		sp := ""
		if len(spell) > 0 {
			sp = "<p>" + strings.Join(spell, "<p></p>") + "</p>"
		}

		sDesc = append(sDesc, &dto.HeroSpell{
			Icon:       mDesc[i].AbilityIconPath,
			Name:       mDesc[i].Name,
			Sort:       sk,
			Desc:       mDesc[i].Description,
			LevelSpell: sp,
		})
	}

	return sDesc
}
