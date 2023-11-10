package attribute

import (
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/context"
	"whisper/pkg/log"
)

type Inner struct {
	ctx      *context.Context
	platform int
}

func NewInnerIns(ctx *context.Context) *Inner {
	return &Inner{ctx: ctx}
}

// GetAll 获取所有的英雄id，包括端游和手游
func (e *Inner) GetAll() []*model.HeroAttribute {
	hd := dao.NewHeroAttributeDAO()
	data, err := hd.Find([]string{
		"DISTINCT(heroId)", "name", "title", "platform",
	}, nil)
	log.Logger.Info(e.ctx, "Hero Attribute find:", len(data))
	if err != nil {
		log.Logger.Info(e.ctx, err.Error())
	}
	return data
}
