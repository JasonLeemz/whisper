package controller

import (
	"whisper/internal/logic/attribute"
	"whisper/internal/logic/common"
	"whisper/internal/logic/equipment"
	rune2 "whisper/internal/logic/rune"
	"whisper/internal/logic/spider"
	"whisper/internal/logic/strategy"
	"whisper/internal/model"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

func GrabStrategy(ctx *context.Context) {
	sp := spider.NewSpider(ctx)
	go sp.BilibiliGrab()
	ctx.Reply("ok", nil)
}

type ReqStrategy struct {
	ID       interface{} `form:"id" json:"id"`
	Keywords string      `form:"keywords" json:"keywords"`
	Platform int         `form:"platform" json:"platform" binding:"-"`
}

func StrategyHero(ctx *context.Context) {
	req := &ReqStrategy{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	if req.Keywords == "" && req.ID != "" {
		// 查询英雄名字
		attr, err := attribute.NewInnerIns(ctx).GetOne(req.ID)
		if err != nil {
			ctx.Reply(nil, errors.New(err))
			return
		}

		if req.Platform == common.PlatformForLOL {
			req.Keywords = attr.Title
		} else {
			req.Keywords = attr.Name
		}
	}

	list, err := strategy.NewStrategy(ctx).List(req.Keywords, req.Platform)
	ctx.Reply(list, errors.New(err))
}

func StrategyEquip(ctx *context.Context) {
	req := &ReqStrategy{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	if req.Keywords == "" && req.ID != "" {
		// 查询装备名字
		if req.Platform == common.PlatformForLOL {
			r := equipment.NewInnerIns(ctx).GetOne(req.Platform, req.ID).(*model.LOLEquipment)
			req.Keywords = r.Name
		} else {
			r := equipment.NewInnerIns(ctx).GetOne(req.Platform, req.ID).(*model.LOLMEquipment)
			req.Keywords = r.Name
		}
	}

	list, err := strategy.NewStrategy(ctx).List(req.Keywords, req.Platform)
	ctx.Reply(list, errors.New(err))
}

func StrategyRune(ctx *context.Context) {
	req := &ReqStrategy{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	if req.Keywords == "" && req.ID != "" {
		// 查询符文名字
		if req.Platform == common.PlatformForLOL {
			r := rune2.NewInnerIns(ctx).GetOne(req.Platform, req.ID).(*model.LOLRune)
			req.Keywords = r.Name
		} else {
			r := rune2.NewInnerIns(ctx).GetOne(req.Platform, req.ID).(*model.LOLMRune)
			req.Keywords = r.Name
		}
	}

	list, err := strategy.NewStrategy(ctx).List(req.Keywords, req.Platform)
	ctx.Reply(list, errors.New(err))
}
