package logic

import (
	errors2 "errors"
	"fmt"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/internal/service"
	"whisper/pkg/context"
	"whisper/pkg/errors"
	"whisper/pkg/log"
)

func QueryHeroes(ctx *context.Context, platform int) (any, *errors.Error) {

	if platform == common.PlatformForLOL {
		heroList, err := service.QueryHeroesForLOL(ctx)
		if err != nil {
			log.Logger.Warn(ctx, err)
		}
		reloadHeroesForLOL(ctx, heroList)
		return heroList, nil
	} else if platform == common.PlatformForLOLM {
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

	// 判断库中是否存在最新版本，如果存在就不更新
	heroesDao := dao.NewLOLHeroesDAO()
	result, err := heroesDao.Find([]string{
		"max(ctime) as ctime",
		"fileTime",
		"version",
	}, nil)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
		return
	}

	if len(result) > 0 {
		log.Logger.Info(ctx,
			fmt.Sprintf("DB Version[%s] fileTime[%s],Data Version:[%s] fileTime[%s]", result[0].Version, result[0].FileTime, heroList.Version, heroList.FileTime),
		)
		x, err := common.CompareTime(result[0].FileTime, heroList.FileTime)
		if err != nil {
			log.Logger.Error(ctx, errors.New(err))
			return
		}
		if x != "<" {
			// 如果原始数据版本和当前获取数据的版本相等，就不更新数据库
			log.Logger.Info(ctx, "原始数据版本和当前获取数据的版本相等,不更新")
			return
		}
	}

	// 入库更新
	heroes := make([]*model.LOLHeroes, 0, len(heroList.Hero))
	heroRole := make([]*model.HeroRole, 0, len(heroList.Hero))

	for _, hero := range heroList.Hero {
		tmp := model.LOLHeroes{
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
				Platform: common.PlatformForLOL,
				HeroId:   hero.HeroId,
				Role:     role,
			}
			heroRole = append(heroRole, hrTmp)
		}
	}

	// 记录英雄列表信息
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
	// 判断库中是否存在最新版本，如果存在就不更新
	heroesDao := dao.NewLOLMHeroesDAO()
	result, err := heroesDao.Find([]string{
		"max(ctime) as ctime",
		"fileTime",
		"version",
	}, nil)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
		return
	}

	if len(result) > 0 {
		log.Logger.Info(ctx,
			fmt.Sprintf("DB Version[%s] fileTime[%s],Data Version:[%s] fileTime[%s]", result[0].Version, result[0].FileTime, heroList.Version, heroList.FileTime),
		)
		x, err := common.CompareTime(result[0].FileTime, heroList.FileTime)
		if err != nil {
			log.Logger.Error(ctx, errors.New(err))
			return
		}
		if x != "<" {
			// 如果原始数据版本和当前获取数据的版本相等，就不更新数据库
			log.Logger.Info(ctx, "原始数据版本和当前获取数据的版本相等,不更新")
			return
		}
	}

	// 入库更新
	heroes := make([]*model.LOLMHeroes, 0, len(heroList.HeroList))
	heroRole := make([]*model.HeroRole, 0, len(heroList.HeroList))

	for _, hero := range heroList.HeroList {
		tmp := model.LOLMHeroes{
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
				Platform: common.PlatformForLOLM,
				HeroId:   hero.HeroId,
				Role:     role,
			}
			heroRole = append(heroRole, hsTmp)
		}
	}

	// 记录英雄列表信息
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
