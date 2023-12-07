package strategy

import (
	"whisper/internal/dto"
	dao2 "whisper/internal/model/DAO"
	"whisper/pkg/context"
)

type Strategy struct {
	ctx *context.Context
}

func NewStrategy(ctx *context.Context) *Strategy {
	return &Strategy{ctx: ctx}
}

func (s *Strategy) List(name string, platform int) ([]*dto.StrategyHero, error) {
	dao := dao2.NewGameStrategyDAO()
	result, err := dao.Find([]string{"*"}, map[string]interface{}{
		"keywords": name,
		"platform": platform,
		"status":   0,
	})
	if err != nil {
		return nil, err
	}

	data := make([]*dto.StrategyHero, 0, len(result))
	for _, d := range result {
		data = append(data, &dto.StrategyHero{
			Keywords:   d.Keywords,
			Title:      d.Title,
			Subtitle:   d.Subtitle,
			PublicDate: d.PublicDate.Format("06/01/02 15:04"),
			Author:     d.Author,
			MainImg:    d.MainImage,
			JumpURL:    d.LinkUrl,
			VideoID:    d.Bvid,
			Length:     d.Length,
			Source:     d.Source,
			Played:     d.Played,
			Platform:   d.Platform,
		})
	}
	return data, nil
}
