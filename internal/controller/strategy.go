package controller

import (
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

type ReqStrategyHero struct {
	Keywords string `form:"keywords" json:"keywords" binding:"required"`
	Platform int    `form:"platform" json:"platform" binding:"-"`
}

func StrategyHero(ctx *context.Context) {
	req := &ReqStrategyHero{}
	if err := ctx.Bind(req); err != nil {
		return
	}
	list, err := strategy.NewStrategy(ctx).Hero(req.Keywords, req.Platform)
	ctx.Reply(list, errors.New(err))
}
