package logic

import (
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/internal/service"
	"whisper/pkg/context"
	"whisper/pkg/errors"
	"whisper/pkg/log"
)

func QueryHeroes(ctx *context.Context, platform int) (any, *errors.Error) {

	heroList, err := service.QueryHeroesForLOL(ctx)
	if err != nil {
		log.Logger.Warn(ctx, err)
	}

	// 入库更新
	heroes := make([]*model.Heroes, 0, len(heroList.Hero))
	heroRole := make([]*model.HeroRole, 0, len(heroList.Hero))

	for _, hero := range heroList.Hero {
		tmp := model.Heroes{
			Platform:            platform,
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
			IsARAMweekfree:      hero.IsARAMweekfree,
			Ispermanentweekfree: hero.Ispermanentweekfree,
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
			hsTmp := &model.HeroRole{
				HeroId: hero.HeroId,
				Role:   role,
			}
			heroRole = append(heroRole, hsTmp)
		}
	}

	// 记录装备信息
	heroesDao := dao.NewHeroesDAO()
	_, err = heroesDao.Add(heroes)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}

	hrDao := dao.NewHeroRoleDAO()
	_, err = hrDao.Add(heroRole)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}

	return heroList, errors.New(err)

}
