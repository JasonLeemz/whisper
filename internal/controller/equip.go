package controller

import (
	"fmt"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"whisper/internal/dto"
	"whisper/internal/logic"
	"whisper/internal/logic/common"
	"whisper/internal/logic/equipment"
	"whisper/internal/service/mq"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

type ReqEquipment struct {
	Platform int `form:"platform" json:"platform" binding:"-"`
	//Platform int `form:"platform" json:"platform" binding:"required"`
}

func Equipment(ctx *context.Context) {

	req := &ReqEquipment{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	equip, err := equipment.QueryEquipments(ctx, req.Platform)

	ctx.Reply(equip, errors.New(err))
}

type ReqGetEquipHeroSuit struct {
	Platform int    `form:"platform" json:"platform" binding:"-"`
	EquipId  int    `json:"id"`
	Version  string `json:"version"`
}

func GetEquipHeroSuit(ctx *context.Context) {
	req := &ReqGetEquipHeroSuit{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	suit, err := logic.GetEquipHeroSuit(ctx, req.Platform, cast.ToString(req.EquipId))
	ctx.Reply(suit, errors.New(err))
}

type ReqEquipExtract struct {
	Platform int `form:"platform" json:"platform" binding:"-"`
}

func EquipExtract(ctx *context.Context) {

	req := &ReqEquipExtract{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	words := equipment.NewInnerIns(ctx).WithPlatform(common.PlatformForLOL).ExtractKeyWords()

	ctx.Reply(words, nil)
}

type ReqEquipFilter struct {
	Platform string          `form:"platform" json:"platform" binding:"-"`
	Keywords map[string]bool `json:"keywords"`
}

func EquipFilter(ctx *context.Context) {

	req := &ReqEquipFilter{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	words := make([]string, 0)
	for kws, state := range req.Keywords {
		if state {
			words = append(words, kws)
		}
	}
	defer func() {
		for _, v := range words {
			mq.ProduceMessage(mq.Exchange, mq.RoutingKeyEquipBox, []byte(v))
		}
	}()

	platform, err := strconv.Atoi(req.Platform)
	if err != nil {
		ctx.Reply(err, errors.New(err, errors.WithMsg(errors.ErrNoInvalidInput)))
	}
	equips, err := equipment.FilterKeyWords(ctx, words, platform)

	resp := dto.SearchResult{}
	total := len(equips)
	resp.Tips = fmt.Sprintf("为您找到相关结果约%d个", total)

	for _, equip := range equips {
		tag := make([]string, 0)
		if equip.Plaintext != "" && !strings.EqualFold(equip.Plaintext, equip.Desc) {
			tag = append(tag, fmt.Sprintf("%s", equip.Plaintext))
		}
		price := int(equip.Price)
		tag = append(tag, fmt.Sprintf("价格:%d", price))
		tag = append(tag, fmt.Sprintf("Version:%s", equip.Version))
		tag = append(tag, fmt.Sprintf("%s", equip.Maps))

		t := dto.SearchResultList{
			Tags:      tag,
			Id:        equip.ID,
			Name:      equip.Name,
			Icon:      equip.Icon,
			Desc:      equip.Desc,
			Plaintext: equip.Plaintext,
			Price:     price,
			Maps:      equip.Maps,
			Platform:  int(equip.Platform),
			Version:   equip.Version,
			Keywords:  equip.Keywords,
		}

		resp.List = append(resp.List, &t)
	}

	ctx.Reply(resp, errors.New(err))
}

func QueryEquipTypes(ctx *context.Context) {

	types, _ := equipment.NewInnerIns(ctx).WithPlatform(common.PlatformForLOL).GetEquipTypes()

	ctx.Reply(map[string]interface{}{
		"types": types,
	}, nil)
}
