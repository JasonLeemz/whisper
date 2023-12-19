package attribute

import (
	"errors"
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
		log.Logger.Error(e.ctx, err.Error())
	}
	return data
}

func (e *Inner) GetOne(heroID interface{}) (*model.HeroAttribute, error) {
	result, err := dao.NewHeroAttributeDAO().Find([]string{"*"}, map[string]interface{}{
		"heroId": heroID,
	})
	if err != nil {
		log.Logger.Error(e.ctx, err.Error())
		return nil, err
	}

	if len(result) == 0 {
		log.Logger.Warn(e.ctx, err.Error())
		return nil, errors.New("heroid not found")
	}

	return result[0], err
}
