package controller

import (
	"whisper/internal/dto"
	"whisper/internal/logic"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

type ReqGetSkillHeroSuit struct {
	Platform int    `form:"platform" json:"platform" binding:"-"`
	SkillID  string `json:"id"`
	Version  string `json:"version"`
}

func GetSkillHeroSuit(ctx *context.Context) {
	req := &ReqGetSkillHeroSuit{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	suit, err := logic.GetSkillHeroSuit(ctx, req.Platform, req.SkillID)
	resp := dto.SthHeroSuit{
		ID:       req.SkillID,
		Platform: req.Platform,
	}
	for _, s := range suit {
		resp.Data = append(resp.Data, s)
	}
	ctx.Reply(resp, errors.New(err))
}
