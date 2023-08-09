package logic

import (
	context2 "context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/internal/service"
	"whisper/pkg/context"
	"whisper/pkg/log"
)

// HeroAttribute
// 根据传过来的id获取数据
// 只有一个id时直接返回
// 多个id时放入队列: 等待全部执行结束，或者有报错时退出
func HeroAttribute(ctx *context.Context, heroID string, platform int) (*dto.HeroAttribute, error) {
	if heroID != "0" {
		attribute, err := heroAttribute(ctx, heroID, platform)
		if err != nil {
			return nil, err
		}
		err = recordHeroAttr(ctx, attribute, platform)
		return attribute, err
	}

	cancelCtx, cancelFunc := context2.WithCancel(ctx) // 用于结束正在运行的任务
	defer cancelFunc()

	var (
		taskErr        error                     // 记录任务报错
		taskChan       = make(chan struct{}, 10) // 控制同时并行的任务
		taskWaiter     = &sync.WaitGroup{}       // 用于等待所有任务完成
		successTaskNum = int32(0)                // 记录成功完成的任务数
		doTaskNum      = int32(0)                // 记录做过的任务数
	)

	var heroIDs []string
	// 获取所有英雄数据
	if platform == common.PlatformForLOL {
		heroesDao := dao.NewLOLHeroesDAO()
		v, err := heroesDao.GetLOLHeroesMaxVersion()
		if err != nil {
			return nil, err
		}
		heroes, err := heroesDao.GetLOLHeroes(v.Version)
		if err != nil {
			return nil, err
		}
		heroIDs = make([]string, 0, len(heroes))
		for i, _ := range heroes {
			heroIDs = append(heroIDs, heroes[i].HeroId)
		}
	} else {
		heroesDao := dao.NewLOLMHeroesDAO()
		v, err := heroesDao.GetLOLMHeroesMaxVersion()
		if err != nil {
			return nil, err
		}
		heroes, err := heroesDao.GetLOLMHeroes(v.Version)
		if err != nil {
			return nil, err
		}
		heroIDs = make([]string, 0, len(heroes))
		for i, _ := range heroes {
			heroIDs = append(heroIDs, heroes[i].HeroId)
		}
	}

	for _, id := range heroIDs {
		select {
		case <-cancelCtx.Done():
			break
		default:
			// 标记 chan 才能执行任务
			taskChan <- struct{}{}

			log.Logger.Info(ctx, "------>", id, "开始执行", "<--------")
			doTaskNum++
			taskWaiter.Add(1)

			var lk sync.Mutex
			// do task
			go func() {
				defer func() {
					taskWaiter.Done()

					// release chan
					<-taskChan
				}()

				select {
				case <-cancelCtx.Done():
					return
				default:
					attribute, err := heroAttribute(ctx, id, platform)
					if err == nil {
						atomic.AddInt32(&successTaskNum, 1) // 执行成功的task数量+1
						if err := recordHeroAttr(ctx, attribute, platform); err != nil {
							log.Logger.Error(ctx, err)
							lk.Lock()
							lk.Unlock()
							taskErr = err
							cancelFunc()
						}
						return
					} else {
						log.Logger.Error(ctx, err)
						lk.Lock()
						lk.Unlock()
						taskErr = err
						cancelFunc()
					}
				}

			}()

		}

	}

	// 等待任务结束
	taskWaiter.Wait()

	log.Logger.Info(ctx, fmt.Sprintf("处理了: %d 个任务", doTaskNum))
	log.Logger.Info(ctx, fmt.Sprintf("提前结束,执行出错: %d 个任务", doTaskNum-successTaskNum))
	log.Logger.Info(ctx, fmt.Sprintf("成功执行了: %d 个任务", successTaskNum))
	log.Logger.Info(ctx, fmt.Sprintf("剩余: %d 个任务待处理", int32(len(heroIDs))-doTaskNum))

	//go heroesAttribute(ctx, platform)

	return nil, taskErr
}

func QueryHeroes(ctx *context.Context, platform int) (any, error) {

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

	return nil, errors.New("请指定platform")

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
		log.Logger.Error(ctx, err)
	}

	//// TODO 记录英雄角色对应表
	//hrDao := dao.NewHeroRoleDAO()
	//_, err = hrDao.Add(heroRole)
	//if err != nil {
	//	log.Logger.Error(ctx,err)
	//}
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
		log.Logger.Error(ctx, err)
	}

	// TODO
	//hrDao := dao.NewHeroRoleDAO()
	//_, err = hrDao.Add(heroRole)
	//if err != nil {
	//	log.Logger.Error(ctx,err)
	//}
}

func heroAttribute(ctx *context.Context, heroID string, platform int) (*dto.HeroAttribute, error) {
	if platform == common.PlatformForLOL {
		return service.GetLOLHeroAttribute(ctx, heroID)
	} else {
		return service.GetLOLMHeroAttribute(ctx, heroID)
	}
}

var lock sync.Mutex

func recordHeroAttr(ctx *context.Context, data *dto.HeroAttribute, platform int) error {

	lock.Lock()
	defer lock.Unlock()

	// 记录HeroRole
	if err := recordHeroRole(ctx, data, platform); err != nil {
		return err
	}

	// 记录HeroSpell
	if err := recordHeroSpell(ctx, data, platform); err != nil {
		return err
	}

	// 记录HeroAttr
	avatar := data.Hero.Avatar
	if len(data.Skins) > 0 {
		avatar = data.Skins[0].IconImg
	}
	attr := &model.HeroAttribute{
		HeroId:              data.Hero.HeroId,
		Title:               data.Hero.Title,
		Name:                data.Hero.Name,
		Alias:               data.Hero.Alias,
		Defense:             data.Hero.Defense,
		Magic:               data.Hero.Magic,
		Difficulty:          data.Hero.Difficulty,
		DifficultyL:         data.Hero.DifficultyL,
		Attack:              data.Hero.Attack,
		Attackrange:         data.Hero.AttackRange,
		Attackdamage:        data.Hero.AttackDamage,
		Attackspeed:         data.Hero.Attackspeed,
		Attackspeedperlevel: data.Hero.Attackspeedperlevel,
		Hp:                  data.Hero.Hp,
		Hpperlevel:          data.Hero.Hpperlevel,
		Mp:                  data.Hero.Mp,
		Mpperlevel:          data.Hero.Mpperlevel,
		Movespeed:           data.Hero.Movespeed,
		Armor:               data.Hero.Armor,
		Armorperlevel:       data.Hero.Armorperlevel,
		Spellblock:          data.Hero.Spellblock,
		Spellblockperlevel:  data.Hero.Spellblockperlevel,
		Hpregen:             data.Hero.Hpregen,
		Hpregenperlevel:     data.Hero.Hpregenperlevel,
		Mpregen:             data.Hero.Mpregen,
		Mpregenperlevel:     data.Hero.Mpregenperlevel,
		Crit:                data.Hero.Crit,
		Damage:              data.Hero.Damage,
		Durability:          data.Hero.Durability,
		Mobility:            data.Hero.Mobility,
		Avatar:              avatar,
		Highlightprice:      data.Hero.Highlightprice,
		GoldPrice:           data.Hero.GoldPrice,
		Couponprice:         data.Hero.Couponprice,
		IsWeekFree:          data.Hero.IsWeekFree,
		Platform:            platform,
		Version:             data.Version,
		FileTime:            data.FileTime,
	}
	ha := dao.NewHeroAttributeDAO()
	_, err2 := ha.Delete(map[string]interface{}{
		"heroId": data.Hero.HeroId,
	})
	if err2 != nil {
		return err2
	}
	affected, err := ha.Add([]*model.HeroAttribute{attr})
	log.Logger.Info(ctx, "attr add record:", affected)
	return err
}
func recordHeroRole(ctx *context.Context, data *dto.HeroAttribute, platform int) error {
	hrs := make([]*model.HeroRole, 0, len(data.Hero.Roles))
	for i, _ := range data.Hero.Roles {
		hrs = append(hrs, &model.HeroRole{
			Platform: platform,
			HeroId:   data.Hero.HeroId,
			Role:     data.Hero.Roles[i],
		})
	}
	hrdao := dao.NewHeroRoleDAO()
	_, err3 := hrdao.Delete(map[string]interface{}{
		"hero_id": data.Hero.HeroId,
	})
	if err3 != nil {
		return err3
	}

	add, err2 := hrdao.Add(hrs)
	if err2 != nil {
		return err2
	}
	log.Logger.Info(ctx, "add hero role success:", add)
	return nil
}
func recordHeroSpell(ctx *context.Context, data *dto.HeroAttribute, platform int) error {
	hrs := make([]*model.HeroSpell, 0, len(data.Spells))
	var keyMap = map[string]int{
		"passive": 0,
		"q":       1,
		"w":       2,
		"e":       3,
		"r":       4,
	}
	apellMap := make(map[string]bool)
	for i, spell := range data.Spells {
		if _, ok := apellMap[spell.Name]; ok {
			continue
		} else {
			apellMap[spell.Name] = true
		}
		detail, _ := json.Marshal(spell.Detail)
		hs := &model.HeroSpell{
			HeroId:          data.Hero.HeroId,
			SpellKey:        spell.SpellKey,
			Sort:            keyMap[spell.SpellKey],
			Name:            spell.Name,
			Description:     spell.Description,
			AbilityIconPath: spell.AbilityIconPath,
			Detail:          string(detail),
			Platform:        platform,
			Version:         data.Version,
			FileTime:        data.FileTime,
		}
		if sort, ok := keyMap[spell.SpellKey]; ok {
			hs.Sort = sort
		} else {
			hs.Sort = -1
			if platform == common.PlatformForLOLM {
				hs.Sort = i
			}
		}
		hrs = append(hrs, hs)
	}
	hrdao := dao.NewHeroSpellDAO()
	_, err3 := hrdao.Delete(map[string]interface{}{
		"heroId": data.Hero.HeroId,
	})
	if err3 != nil {
		return err3
	}

	add, err2 := hrdao.Add(hrs)
	if err2 != nil {
		return err2
	}
	log.Logger.Info(ctx, "add hero spell success:", add)
	return nil
}
