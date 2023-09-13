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
		esd := dao.NewESEquipmentDAO()
		esEquipments, err := esd.Find(ctx, cond)
		if err != nil {
			ctx.Reply(keywords, errors.New(err))
		}

		for _, item := range esEquipments {
			if item.Maps == "召唤师峡谷" {
				name := strings.TrimSpace(item.Name)
				keywords = append(keywords, name)
			}
		}
	case model.NewModelESHeroes().GetIndexName():
		esd := dao.NewESHeroesDAO()
		esHeroes, err := esd.Find(ctx, cond)
		if err != nil {
			ctx.Reply(keywords, errors.New(err))
		}

		for _, item := range esHeroes {
			name := strings.TrimSpace(item.Name)
			keywords = append(keywords, name)
		}
	case model.NewModelESRune().GetIndexName():
		esd := dao.NewESRuneDAO()
		esRune, err := esd.Find(ctx, cond)
		if err != nil {
			ctx.Reply(keywords, errors.New(err))
		}

		for _, item := range esRune {
			name := strings.TrimSpace(item.Name)
			keywords = append(keywords, name)
		}
	case model.NewModelESSkill().GetIndexName():
		esd := dao.NewESSkillDAO()
		esSkill, err := esd.Find(ctx, cond)
		if err != nil {
			ctx.Reply(keywords, errors.New(err))
		}

		for _, item := range esSkill {
			name := strings.TrimSpace(item.Name)
			keywords = append(keywords, name)
		}
	}

	ctx.Reply(keywords, nil)
}
