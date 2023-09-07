package controller

import (
	"fmt"
	"strconv"
	"strings"
	"whisper/internal/dto"
	"whisper/internal/logic"
	"whisper/internal/service/mq"
	"whisper/pkg/errors"

	"whisper/pkg/context"
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

	equip, err := logic.QueryEquipments(ctx, req.Platform)

	ctx.Reply(equip, errors.New(err))
}

type ReqHeroes struct {
	Platform int `form:"platform" json:"platform" binding:"-"`
}

func Heroes(ctx *context.Context) {

	req := &ReqHeroes{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	equip, err := logic.QueryHeroes(ctx, req.Platform)

	ctx.Reply(equip, errors.New(err))
}

type ReqHeroesAttribute struct {
	Platform int    `form:"platform" json:"platform" binding:"-"`
	HeroID   string `form:"hero_id" json:"hero_id" binding:"required"`
}

func HeroesAttribute(ctx *context.Context) {

	req := &ReqHeroesAttribute{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	attr, err := logic.HeroAttribute(ctx, req.HeroID, req.Platform)

	ctx.Reply(attr, errors.New(err))
}

type ReqRune struct {
	Platform int `form:"platform" json:"platform" binding:"-"`
}

func Rune(ctx *context.Context) {

	req := &ReqRune{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	equip, err := logic.QueryRune(ctx, req.Platform)

	ctx.Reply(equip, errors.New(err))
}

type ReqRuneType struct {
	Platform int `form:"platform" json:"platform" binding:"-"`
}

func RuneType(ctx *context.Context) {

	req := &ReqRuneType{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	equip, err := logic.QueryRuneType(ctx, req.Platform)

	ctx.Reply(equip, errors.New(err))
}

type ReqSkill struct {
	Platform int `form:"platform" json:"platform" binding:"-"`
}

func Skill(ctx *context.Context) {

	req := &ReqSkill{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	equip, err := logic.QuerySkill(ctx, req.Platform)

	ctx.Reply(equip, errors.New(err))
}

type ReqEquipExtract struct {
	Platform int `form:"platform" json:"platform" binding:"-"`
}

func EquipExtract(ctx *context.Context) {

	req := &ReqEquipExtract{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	words := logic.ExtractKeyWords(ctx, req.Platform)

	ctx.Reply(words, errors.New(nil))
}

type ReqEquipFilter struct {
	Platform string   `form:"platform" json:"platform" binding:"-"`
	Keywords []string `json:"keywords"`
}

func EquipFilter(ctx *context.Context) {

	req := &ReqEquipFilter{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	defer func() {
		for _, kws := range req.Keywords {
			mq.ProduceMessage(mq.Exchange, mq.RoutingKeyEquipBox, []byte(kws))
		}
	}()

	platform, err := strconv.Atoi(req.Platform)
	if err != nil {
		ctx.Reply(nil, errors.New(err, errors.ErrNoInvalidInput))
	}
	equips, err := logic.FilterKeyWords(ctx, req.Keywords, platform)

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

type ReqSuitEquip struct {
	Platform int    `form:"platform" json:"platform" binding:"-"`
	HeroId   string `json:"hero_id"`
}

func SuitEquip(ctx *context.Context) {
	req := &ReqSuitEquip{}
	if err := ctx.Bind(req); err != nil {
		return
	}
	equips, err := logic.QuerySuitEquip(ctx, req.Platform, req.HeroId)
	ctx.Reply(equips, errors.New(err))
}

func BatchUpdateSuitEquip(ctx *context.Context) {
	err := logic.BatchUpdateSuitEquip(ctx)
	ctx.Reply(nil, errors.New(err))
}

func SuitData2Redis(ctx *context.Context) {
	err := logic.SuitData2Redis(ctx)
	ctx.Reply(nil, errors.New(err))
}

func SuitHeroData2Redis(ctx *context.Context) {
	err := logic.SuitHeroData2Redis(ctx)
	ctx.Reply(nil, errors.New(err))
}

type ReqGetHeroSuit struct {
	Platform int    `form:"platform" json:"platform" binding:"-"`
	HeroId   string `json:"hero_id"`
}

func GetHeroSuit(ctx *context.Context) {
	req := &ReqGetHeroSuit{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	suit, err := logic.GetHeroSuit(ctx, req.HeroId)
	ctx.Reply(suit, errors.New(err))
}

type ReqHeroesPosition struct {
	Platform int `form:"platform" json:"platform" binding:"-"`
}

func HeroesPosition(ctx *context.Context) {
	req := &ReqHeroesPosition{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	data, err := logic.HeroesPosition(ctx, req.Platform)
	ctx.Reply(data, errors.New(err))
}
