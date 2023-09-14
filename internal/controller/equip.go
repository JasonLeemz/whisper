package controller

import (
	"github.com/spf13/cast"
	"whisper/internal/logic"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

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
