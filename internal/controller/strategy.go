package controller

import (
	"github.com/spf13/cast"
	"whisper/internal/logic"
	"whisper/internal/logic/common"
	"whisper/internal/logic/spider"
	"whisper/internal/logic/strategy"
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
		attr, err := logic.QueryHeroAttribute(ctx, cast.ToString(req.ID), req.Platform)
		if err != nil {
			ctx.Reply(nil, errors.New(err))
			return
		}

		if req.Platform == common.PlatformForLOL {
			req.Keywords = attr.Hero.Title
		} else {
			req.Keywords = attr.Hero.Name
		}
	}

	list, err := strategy.NewStrategy(ctx).Hero(req.Keywords, req.Platform)
	ctx.Reply(list, errors.New(err))
}

func StrategyEquip(ctx *context.Context) {
	req := &ReqStrategy{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	if req.Keywords == "" && req.ID != "" {
		// 查询装备名字
		attr, err := logic.QueryHeroAttribute(ctx, cast.ToString(req.ID), req.Platform)
		if err != nil {
			ctx.Reply(nil, errors.New(err))
			return
		}

		if req.Platform == common.PlatformForLOL {
			req.Keywords = attr.Hero.Title
		} else {
			req.Keywords = attr.Hero.Name
		}
	}

	list, err := strategy.NewStrategy(ctx).Hero(req.Keywords, req.Platform)
	ctx.Reply(list, errors.New(err))
}

func StrategyRune(ctx *context.Context) {
	req := &ReqStrategy{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	if req.Keywords == "" && req.ID != "" {
		// 查询英雄名字
		attr, err := logic.QueryHeroAttribute(ctx, cast.ToString(req.ID), req.Platform)
		if err != nil {
			ctx.Reply(nil, errors.New(err))
			return
		}

		if req.Platform == common.PlatformForLOL {
			req.Keywords = attr.Hero.Title
		} else {
			req.Keywords = attr.Hero.Name
		}
	}

	list, err := strategy.NewStrategy(ctx).Hero(req.Keywords, req.Platform)
	ctx.Reply(list, errors.New(err))
}
