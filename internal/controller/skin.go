package controller

import (
	"math/rand"
	"whisper/internal/dto"
	"whisper/internal/logic"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

type ReqGetHeroSkins struct {
	Platform int    `form:"platform" json:"platform" binding:"-"`
	HeroId   string `json:"id"`
}

func GetHeroSkins(ctx *context.Context) {
	req := &ReqGetHeroSkins{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	heroID := req.HeroId
	if heroID == "" {
		attrs, err := logic.GetAllHeroesFromAttr(ctx, []int{0})
		if err == nil {
			idx := rand.Intn(len(attrs)) // [0,n)
			heroID = attrs[idx].HeroId
		}
	}
	attr, err := logic.GetAttribute(ctx, req.Platform, heroID)
	if err != nil {
		ctx.Reply(nil, errors.New(err))
	}

	skins, err := logic.GetHeroSkins(ctx, req.Platform, heroID)
	var resp []*dto.RespHeroSkins
	for _, skin := range skins {
		resp = append(resp, &dto.RespHeroSkins{
			HeroId:     skin.HeroId,
			Title:      attr.Title,
			Name:       attr.Name,
			ShortBio:   attr.ShortBio,
			MainImg:    skin.MainImg,
			IconImg:    skin.IconImg,
			LoadingImg: skin.LoadingImg,
			VideoImg:   skin.VideoImg,
			SourceImg:  skin.SourceImg,
			Platform:   skin.Platform,
			Version:    skin.Version,
		})
	}
	ctx.Reply(resp, errors.New(err))
}
