package logic

import (
	errors2 "errors"
	"whisper/internal/dto"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/internal/service"
	"whisper/pkg/context"
	"whisper/pkg/errors"
	"whisper/pkg/log"
)

func QueryHeroes(ctx *context.Context, platform int) (any, *errors.Error) {

	if platform == platformForLOL {
		heroList, err := service.QueryHeroesForLOL(ctx)
		if err != nil {
			log.Logger.Warn(ctx, err)
		}
		reloadHeroesForLOL(ctx, heroList)
		return heroList, nil
	} else if platform == platformForLOLM {
		heroList, err := service.QueryHeroesForLOLM(ctx)
		if err != nil {
			log.Logger.Warn(ctx, err)
		}
		reloadHeroesForLOLM(ctx, heroList)
		return heroList, nil
	}

	return nil, errors.New(errors2.New("请指定游戏平台"), errors.ErrNoInvalidInput)

}

func reloadHeroesForLOL(ctx *context.Context, heroList *dto.LOLHeroes) {
	// 入库更新
	heroes := make([]*model.LOLHeroes, 0, len(heroList.Hero))
	heroRole := make([]*model.HeroRole, 0, len(heroList.Hero))
	var err error

	for _, hero := range heroList.Hero {
		tmp := model.LOLHeroes{
			Platform:            platformForLOL,
			HeroId:              hero.HeroId,
			Name:                hero.Name,
			Alias:               hero.Alias,
			Title:               hero.Title,
			IsWeekFree:          hero.IsWeekFree,
			Attack:              hero.Attack,
			Defense:             hero.Defense,
			Magic:               hero.Magic,
			Difficulty:          hero.Difficulty,
			SelectAudio:         hero.SelectAudio,
			BanAudio:            hero.BanAudio,
			IsARAMWeekFree:      hero.IsARAMWeekFree,
			IsPermanentWeekFree: hero.IsPermanentWeekFree,
			ChangeLabel:         hero.ChangeLabel,
			GoldPrice:           hero.GoldPrice,
			CouponPrice:         hero.CouponPrice,
			Camp:                hero.Camp,
			CampId:              hero.CampId,
			Keywords:            hero.Keywords,
			InstanceId:          hero.InstanceId,
			Version:             heroList.Version,
			FileTime:            heroList.FileTime,
		}
		heroes = append(heroes, &tmp)

		// 记录英雄角色表
		for _, role := range hero.Roles {
			hrTmp := &model.HeroRole{
				Platform: platformForLOL,
				HeroId:   hero.HeroId,
				Role:     role,
			}
			heroRole = append(heroRole, hrTmp)
		}
	}

	// 记录英雄列表信息
	heroesDao := dao.NewLOLHeroesDAO()
	_, err = heroesDao.Add(heroes)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}

	// 记录英雄角色对应表
	hrDao := dao.NewHeroRoleDAO()
	_, err = hrDao.Add(heroRole)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}
}
func reloadHeroesForLOLM(ctx *context.Context, heroList *dto.LOLMHeroes) {

	// 入库更新
	heroes := make([]*model.LOLMHeroes, 0, len(heroList.HeroList))
	heroRole := make([]*model.HeroRole, 0, len(heroList.HeroList))
	var err error

	for _, hero := range heroList.HeroList {
		tmp := model.LOLMHeroes{
			Platform:       platformForLOLM,
			HeroId:         hero.HeroId,
			Name:           hero.Name,
			Title:          hero.Title,
			Intro:          hero.Intro,
			Avatar:         hero.Avatar,
			Card:           hero.Card,
			Poster:         hero.Poster,
			Highlightprice: hero.HighlightPrice,
			Couponprice:    hero.CouponPrice,
			Alias:          hero.Alias,
			Lane:           hero.Lane,
			Tags:           hero.Tags,
			Searchkey:      hero.SearchKey,
			IsWeekFree:     hero.IsWeekFree,
			DifficultyL:    hero.DifficultyL,
			Damage:         hero.Damage,
			SurviveL:       hero.SurviveL,
			AssistL:        hero.AssistL,
			Version:        heroList.Version,
			FileTime:       heroList.FileTime,
		}
		heroes = append(heroes, &tmp)

		// 记录英雄角色表
		for _, role := range hero.Roles {
			hsTmp := &model.HeroRole{
				Platform: platformForLOLM,
				HeroId:   hero.HeroId,
				Role:     role,
			}
			heroRole = append(heroRole, hsTmp)
		}
	}

	// 记录英雄列表信息
	heroesDao := dao.NewLOLMHeroesDAO()
	_, err = heroesDao.Add(heroes)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}

	hrDao := dao.NewHeroRoleDAO()
	_, err = hrDao.Add(heroRole)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}
}
