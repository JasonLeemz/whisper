package logic

import (
	"fmt"
	"time"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	lol "whisper/internal/service/lol"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/pinyin"
)

func QueryHeroes(ctx *context.Context, platform int) (any, error) {
	heroes, err := lol.CreateLOLProduct(platform)().QueryHeroes(ctx)
	if err != nil {
		log.Logger.Error(ctx, err)
		return nil, err
	}

	if platform == common.PlatformForLOL {
		reloadHeroesForLOL(ctx, heroes.(*dto.LOLHeroes))
	} else if platform == common.PlatformForLOLM {
		reloadHeroesForLOLM(ctx, heroes.(*dto.LOLMHeroes))
	}

	return heroes, nil

}

func reloadHeroesForLOL(ctx *context.Context, heroList *dto.LOLHeroes) {

	// 判断库中是否存在最新版本，如果存在就不更新
	heroesDao := dao.NewLOLHeroesDAO()
	result, err := heroesDao.GetLOLHeroesMaxVersion()
	if err != nil {
		log.Logger.Error(ctx, err)
		return
	}

	if result != nil {
		log.Logger.Info(ctx,
			fmt.Sprintf("DB Version[%s] fileTime[%s],Data Version:[%s] fileTime[%s]", result.Version, result.FileTime, heroList.Version, heroList.FileTime),
		)
		x, err := common.CompareTime(result.FileTime, heroList.FileTime)
		if err != nil {
			log.Logger.Error(ctx, err)
			return
		}
		if x != "<" {
			// 如果原始数据版本和当前获取数据的版本相等，就不更新数据库
			log.Logger.Info(ctx, "原始数据版本和当前获取数据的版本相等,不更新")
			return
		}

		// 有可能日期更新了，但是版本号没变
		if result.Version == heroList.Version {
			// 将db中该版本的数据status置为1
			cond := map[string]interface{}{
				"version": result.Version,
				"status":  0,
			}
			data := model.LOLHeroes{
				Status: 1,
			}
			up, err := heroesDao.Update(&data, cond)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}
			log.Logger.Info(ctx, "当前版本数据不是最新,已经软删除,生效行数:", up)
		}

	}

	log.Logger.Info(ctx, "running record LOL heroes data...")
	startT := time.Now()
	// 入库更新
	heroes := make([]*model.LOLHeroes, 0, len(heroList.Hero))

	for _, hero := range heroList.Hero {
		namePY, nameF := pinyin.Trans(hero.Name)
		titlePY, titleF := pinyin.Trans(hero.Title)
		searchKey := hero.Keywords + "," + namePY + "," + nameF + "," + titlePY + "," + titleF

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
			Keywords:            searchKey,
			InstanceId:          hero.InstanceId,
			Version:             heroList.Version,
			FileTime:            heroList.FileTime,
		}
		heroes = append(heroes, &tmp)

	}

	// 记录英雄列表信息
	_, err = heroesDao.Add(heroes)
	if err != nil {
		log.Logger.Error(ctx, err)
	}

	log.Logger.Info(ctx, fmt.Sprintf("finish record LOL heroes data. Since:%fs", time.Since(startT).Seconds()))
}
func reloadHeroesForLOLM(ctx *context.Context, heroList *dto.LOLMHeroes) {
	// 判断库中是否存在最新版本，如果存在就不更新
	heroesDao := dao.NewLOLMHeroesDAO()
	result, err := heroesDao.GetLOLMHeroesMaxVersion()
	if err != nil {
		log.Logger.Error(ctx, err)
		return
	}

	if result != nil {
		log.Logger.Info(ctx,
			fmt.Sprintf("DB Version[%s] fileTime[%s],Data Version:[%s] fileTime[%s]", result.Version, result.FileTime, heroList.Version, heroList.FileTime),
		)
		x, err := common.CompareTime(result.FileTime, heroList.FileTime)
		if err != nil {
			log.Logger.Error(ctx, err)
			return
		}
		if x != "<" {
			// 如果原始数据版本和当前获取数据的版本相等，就不更新数据库
			log.Logger.Info(ctx, "原始数据版本和当前获取数据的版本相等,不更新")
			return
		}
		// 有可能日期更新了，但是版本号没变
		if result.Version == heroList.Version {
			// 将db中该版本的数据status置为1
			cond := map[string]interface{}{
				"version": result.Version,
				"status":  0,
			}
			data := model.LOLMHeroes{
				Status: 1,
			}
			up, err := heroesDao.Update(&data, cond)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}
			log.Logger.Info(ctx, "当前版本数据不是最新,已经软删除,生效行数:", up)
		}
	}

	log.Logger.Info(ctx, "running record LOLM heroes data...")
	startT := time.Now()
	// 入库更新
	heroes := make([]*model.LOLMHeroes, 0, len(heroList.HeroList))
	heroRole := make([]*model.HeroRole, 0, len(heroList.HeroList))

	for _, hero := range heroList.HeroList {
		namePY, nameF := pinyin.Trans(hero.Name)
		titlePY, titleF := pinyin.Trans(hero.Title)
		searchKey := hero.SearchKey + "," + namePY + "," + nameF + "," + titlePY + "," + titleF
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
			Searchkey:      searchKey,
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
		log.Logger.Error(ctx, err)
	}

	log.Logger.Info(ctx, fmt.Sprintf("finish record LOLM heroes data. Since:%fs", time.Since(startT).Seconds()))
}
