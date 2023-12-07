package controller

import (
	"whisper/internal/dto"
	"whisper/internal/logic"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

type ReqGetRuneHeroSuit struct {
	Platform int    `form:"platform" json:"platform" binding:"-"`
	RuneId   string `json:"id"`
	Version  string `json:"version"`
}

func GetRuneHeroSuit(ctx *context.Context) {
	req := &ReqGetRuneHeroSuit{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	suit, err := logic.GetRuneHeroSuit(ctx, req.Platform, req.RuneId)

	resp := dto.SthHeroSuit{
		ID:       req.RuneId,
		Platform: req.Platform,
	}
	for _, s := range suit {
		resp.Data = append(resp.Data, s)
	}
	ctx.Reply(resp, errors.New(err))
}
