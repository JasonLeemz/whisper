package controller

import (
	"whisper/internal/logic/position"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

type ReqHeroesPosition struct {
	Platform int `form:"platform" json:"platform" binding:"-"`
}

func HeroesPosition(ctx *context.Context) {
	req := &ReqHeroesPosition{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	data, err := position.NewPosition()(ctx, req.Platform).HeroesPosition()
	ctx.Reply(data, errors.New(err))
}
