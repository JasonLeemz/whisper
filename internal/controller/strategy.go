package controller

import (
	"whisper/internal/logic/spider"
	"whisper/pkg/context"
)

func GrabStrategy(ctx *context.Context) {
	sp := spider.NewSpider(ctx)
	sp.BilibiliGrab()
	ctx.Reply("ok", nil)
}
