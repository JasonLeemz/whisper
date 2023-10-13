package controller

import (
	"strings"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/internal/model/common"
	"whisper/pkg/context"
	"whisper/pkg/errors"
)

type ReqAutoComplete struct {
	Platform string `form:"platform" json:"platform" binding:"-"`
	KeyWords string `form:"key_words" json:"key_words" binding:"-"`
	Category string `form:"category" json:"category" binding:"-"`
}

func AutoComplete(ctx *context.Context) {
	req := &ReqAutoComplete{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	var (
		keywords []string
		mk       = make(map[string]struct{})
	)

	cond := &common.QueryCond{
		MultiMatchQuery: &common.MultiMatchQuery{
			Text: req.KeyWords,
			Fields: []string{
				"name", "keywords",
			},
		},
		//TermsQuery: &common.TermsQuery{
		//	Name:   "",
		//	Values: nil,
		//},
		TermQuery: []*common.TermQuery{
			&common.TermQuery{
				Name:  "platform",
				Value: req.Platform,
			},
		},
		//FieldSort: &common.FieldSort{
		//	Field:     "",
		//	Direction: "",
		//},
	}

	if req.KeyWords == "" {
		cond.MultiMatchQuery = nil
	}

	switch req.Category {
	case model.NewModelESEquipment().GetIndexName():
		data, err := dao.CreateEsDao(model.NewModelESEquipment().GetIndexName()).(*dao.ESEquipmentDAO).Find(ctx, cond)
		if err != nil {
			ctx.Reply(keywords, errors.New(err))
		}

		for _, item := range data {
			if item.Maps == "召唤师峡谷" {
				name := strings.TrimSpace(item.Name)
				mk[name] = struct{}{}
			}
		}
	case model.NewModelESHeroes().GetIndexName():
		data, err := dao.CreateEsDao(model.NewModelESHeroes().GetIndexName()).(*dao.ESHeroesDAO).Find(ctx, cond)
		if err != nil {
			ctx.Reply(keywords, errors.New(err))
		}

		for _, item := range data {
			name := strings.TrimSpace(item.Name)
			mk[name] = struct{}{}
		}
	case model.NewModelESRune().GetIndexName():
		data, err := dao.CreateEsDao(model.NewModelESRune().GetIndexName()).(*dao.ESRuneDAO).Find(ctx, cond)
		if err != nil {
			ctx.Reply(keywords, errors.New(err))
		}

		for _, item := range data {
			name := strings.TrimSpace(item.Name)
			mk[name] = struct{}{}
		}
	case model.NewModelESSkill().GetIndexName():
		data, err := dao.CreateEsDao(model.NewModelESSkill().GetIndexName()).(*dao.ESSkillDAO).Find(ctx, cond)
		if err != nil {
			ctx.Reply(keywords, errors.New(err))
		}

		for _, item := range data {
			name := strings.TrimSpace(item.Name)
			mk[name] = struct{}{}
		}
	}

	// 去重
	for name, _ := range mk {
		keywords = append(keywords, name)
	}

	ctx.Reply(keywords, nil)
}
