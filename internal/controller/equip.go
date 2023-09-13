package controller

import (
	"whisper/internal/logic"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

type ReqGetEquipHeroSuit struct {
	Platform int    `form:"platform" json:"platform" binding:"-"`
	SkillId  string `json:"id"`
	Version  string `json:"version"`
}

func GetEquipHeroSuit(ctx *context.Context) {
	req := &ReqGetEquipHeroSuit{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	suit, err := logic.GetHeroSuit(ctx, req.SkillId)
	ctx.Reply(suit, errors.New(err))
}
