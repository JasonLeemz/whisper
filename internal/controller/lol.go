package controller

import (
	"whisper/internal/logic"
	"whisper/internal/logic/common"
	"whisper/internal/logic/suit"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

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

type ReqSuitEquip struct {
	Platform int    `form:"platform" json:"platform" binding:"-"`
	HeroId   string `json:"hero_id"`
}

func SuitEquip(ctx *context.Context) {
	req := &ReqSuitEquip{}
	if err := ctx.Bind(req); err != nil {
		return
	}
	equips, err := suit.NewSuit()(ctx, req.Platform).QuerySuitEquip(req.HeroId)
	ctx.Reply(equips, errors.New(err))
}

func BatchUpdateSuitEquip(ctx *context.Context) {
	suit.NewSuit()(ctx, common.PlatformForLOL).BatchUpdateSuitEquip()
	suit.NewSuit()(ctx, common.PlatformForLOLM).BatchUpdateSuitEquip()
	ctx.Reply(nil, errors.New(nil))
}

func SuitData2Redis(ctx *context.Context) {
	err := suit.SuitData2Redis(ctx)
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
