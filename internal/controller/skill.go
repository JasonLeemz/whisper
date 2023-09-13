package controller

import (
	"whisper/internal/logic"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

type ReqGetSkillHeroSuit struct {
	Platform int    `form:"platform" json:"platform" binding:"-"`
	RuneId   string `json:"id"`
	Version  string `json:"version"`
}

func GetSkillHeroSuit(ctx *context.Context) {
	req := &ReqGetSkillHeroSuit{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	suit, err := logic.GetSkillHeroSuit(ctx, req.Platform, req.RuneId)
	ctx.Reply(suit, errors.New(err))
}
