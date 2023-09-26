package logic

import (
	context2 "context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	redis2 "github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/internal/service"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/redis"
	"whisper/pkg/utils"
)

func BatchUpdateSuitEquip(ctx *context.Context) error {

	// 获取所有英雄ID
	ha := dao.NewHeroAttributeDAO()
	heroes, err := ha.Find([]string{
		"DISTINCT(heroId)", "name", "title", "platform",
	}, nil)
	if err != nil {
		return err
	}

	cancelCtx, cancelFunc := context2.WithCancel(ctx)
	defer cancelFunc()

	var (
		taskAll  = int32(len(heroes))
		taskSucc = int32(0)
		taskFail = int32(0)
		taskDone = int32(0)
		wg       = &sync.WaitGroup{}
		ch       = make(chan struct{}, 20)
	)

	for i, hero := range heroes {
		select {
		case <-cancelCtx.Done():
			break
		default:
			log.Logger.Info(ctx, ">>>>>>>>>>开始处理 hero:<<<<<<<<<<<", i, "/", hero.HeroId)
			ch <- struct{}{}
			wg.Add(1)

			go func(hero *model.HeroAttribute) {
				defer func() {
					<-ch
					wg.Done()
					atomic.AddInt32(&taskDone, 1)
				}()

				_, err2 := QuerySuitEquip(ctx, hero.Platform, hero.HeroId)
				// 任务执行失败，这个地方可以使用锁，也可以使用原子操作，优先原子操作
				if err2 != nil {
					atomic.AddInt32(&taskFail, 1)
					cancelFunc()
					log.Logger.Error(ctx, err2)
					return
				} else {
					atomic.AddInt32(&taskSucc, 1)
				}
			}(hero)

		}
	}

	wg.Wait()

	log.Logger.Info(ctx, fmt.Sprintf("处理了: %d 个任务", taskDone))
	log.Logger.Info(ctx, fmt.Sprintf("提前结束,执行出错: %d 个任务", taskFail))
	log.Logger.Info(ctx, fmt.Sprintf("成功执行了: %d 个任务", taskSucc))
	log.Logger.Info(ctx, fmt.Sprintf("剩余: %d 个任务待处理", taskAll-taskDone))

	return nil
}
func QuerySuitEquip(ctx *context.Context, platform int, heroId string) (any, error) {

	if platform == common.PlatformForLOL {
		fightData, err := getFightData(ctx, heroId)
		if err != nil {
			return nil, errors.New("getFightData:" + err.Error())
		}

		// reload heroes_position 表
		err = updateHeroesPosition(ctx, platform, heroId, fightData)
		if err != nil {
			return nil, errors.New("updateHeroesPosition:" + err.Error())
		}

		// reload heroes_suit 表
		err = updateLOLHeroesSuit(ctx, heroId, fightData)
		if err != nil {
			return nil, errors.New("updateHeroesSuit:" + err.Error())
		}
		return fightData, nil
	} else {
		// common.PlatformForLOLM
		heroTech, equipTechs, err := service.HeroSuit(ctx, heroId)
		if err != nil {
			return nil, errors.New("service.HeroSuit:" + err.Error())
		}
		// reload heroes_suit 表
		err = updateLOLMHeroesSuit(ctx, heroId, heroTech, equipTechs)
		if err != nil {
			return nil, errors.New("updateLOLMHeroesSuit:" + err.Error())
		}

		return []any{
			heroTech, equipTechs,
		}, nil
	}
}

// LOL英雄的rank数据
func getFightData(ctx *context.Context, heroId string) (*dto.ChampionFightData, error) {
	fightData, err := service.ChampionFightData(ctx, heroId)
	if err != nil {
		return nil, err
	}
	for pos, posData := range fightData.List.ChampionLane {
		equipData := map[string]dto.Itemjson{}
		tmp := dto.ChampionLaneItem{}

		var err error
		err = json.Unmarshal([]byte(posData.Itemoutjson), &equipData)
		if err != nil {
			log.Logger.Warn(ctx, err, "heroid:", heroId)
		} else {
			tmp.Itemout = equipData
		}

		equipData = *new(map[string]dto.Itemjson)
		err = json.Unmarshal([]byte(posData.Core3itemjson), &equipData)
		if err != nil {
			log.Logger.Warn(ctx, err, "heroid:", heroId)
		} else {
			tmp.Core3item = equipData
		}

		equipData = *new(map[string]dto.Itemjson)
		err = json.Unmarshal([]byte(posData.Shoesjson), &equipData)
		if err != nil {
			log.Logger.Warn(ctx, err, "heroid:", heroId)
		} else {
			tmp.Shoes = equipData
		}

		var suits []dto.Itemjson
		err = json.Unmarshal([]byte(posData.Hold3), &suits)
		if err != nil {
			log.Logger.Warn(ctx, err, "heroid:", heroId)
		} else {
			tmp.Suits = suits
		}

		fightData.List.ChampionLane[pos] = tmp
	}

	return fightData, nil
}
func updateHeroesPosition(ctx *context.Context, platform int, heroId string, fightData *dto.ChampionFightData) error {
	hpd := dao.NewHeroesPositionDAO()
	posData := make([]*model.HeroesPosition, 0, 3)
	for pos, _ := range fightData.List.ChampionFight {
		posData = append(posData, &model.HeroesPosition{
			HeroId:   heroId,
			Pos:      pos,
			Platform: platform,
			Version:  fightData.GameVer,
			FileTime: fightData.Date,
		})
	}
	if len(posData) == 0 {
		log.Logger.Warn(ctx, "posData is nil", "heroId:", heroId)
		return nil
	}

	err := hpd.DeleteAndInsert(map[string]interface{}{
		"heroId": heroId,
	}, posData)
	if err != nil {
		return errors.New("Add LOL HeroesPosition " + err.Error() + ",heroId:" + heroId)
	}
	log.Logger.Info(ctx, "Add HeroesPosition heroId:", heroId)

	return nil
}
func updateLOLHeroesSuit(ctx *context.Context, heroId string, fightData *dto.ChampionFightData) error {
	platform := common.PlatformForLOL
	hpd := dao.NewHeroesSuitDAO()

	posData := make([]*model.HeroesSuit, 0)
	var m model.HeroesSuit
	for pos, pds := range fightData.List.ChampionLane {
		posCopy := pos
		for _, pdsd := range pds.Itemout {
			itemidArr := strings.Split(pdsd.Itemid, "&")
			itemids := strings.Join(itemidArr, ",")
			posData = append(posData, &model.HeroesSuit{
				HeroId:   heroId,
				Pos:      posCopy,
				Itemids:  itemids,
				Igamecnt: pdsd.Igamecnt,
				Wincnt:   pdsd.Wincnt,
				Winrate:  pdsd.Winrate,
				Allcnt:   pdsd.Allcnt,
				Showrate: pdsd.Showrate,
				Type:     m.TypeOut(),
				Platform: platform,
				Version:  fightData.GameVer,
				FileTime: fightData.Date,
			})
		}

		for _, pdsd := range pds.Core3item {
			itemidArr := strings.Split(pdsd.Itemid, "&")
			itemids := strings.Join(itemidArr, ",")
			posData = append(posData, &model.HeroesSuit{
				HeroId:   heroId,
				Pos:      posCopy,
				Itemids:  itemids,
				Igamecnt: pdsd.Igamecnt,
				Wincnt:   pdsd.Wincnt,
				Winrate:  pdsd.Winrate,
				Allcnt:   pdsd.Allcnt,
				Showrate: pdsd.Showrate,
				Type:     m.TypeCore(),
				Platform: platform,
				Version:  fightData.GameVer,
				FileTime: fightData.Date,
			})
		}

		for _, pdsd := range pds.Shoes {
			itemidArr := strings.Split(pdsd.Itemid, "&")
			itemids := strings.Join(itemidArr, ",")
			posData = append(posData, &model.HeroesSuit{
				HeroId:   heroId,
				Pos:      posCopy,
				Itemids:  itemids,
				Igamecnt: pdsd.Igamecnt,
				Wincnt:   pdsd.Wincnt,
				Winrate:  pdsd.Winrate,
				Allcnt:   pdsd.Allcnt,
				Showrate: pdsd.Showrate,
				Type:     m.TypeShoes(),
				Platform: platform,
				Version:  fightData.GameVer,
				FileTime: fightData.Date,
			})
		}

		for _, pdsd := range pds.Suits {
			itemidArr := strings.Split(pdsd.Itemid, "&")
			itemids := strings.Join(itemidArr, ",")
			posData = append(posData, &model.HeroesSuit{
				HeroId:   heroId,
				Pos:      posCopy,
				Itemids:  itemids,
				Igamecnt: pdsd.Igamecnt,
				Wincnt:   pdsd.Wincnt,
				Winrate:  pdsd.Winrate,
				Allcnt:   pdsd.Allcnt,
				Showrate: pdsd.Showrate,
				Type:     m.TypeOther(),
				Platform: platform,
				Version:  fightData.GameVer,
				FileTime: fightData.Date,
			})
		}
	}
	if len(posData) == 0 {
		log.Logger.Warn(ctx, "posData is nil", "heroId:", heroId)
		return nil
	}

	err := hpd.DeleteAndInsert(map[string]interface{}{
		"heroId": heroId,
	}, posData)
	if err != nil {
		return errors.New("Add LOLM HeroesSuit " + err.Error() + ",heroId:" + heroId)
	}
	log.Logger.Info(ctx, "Add LOLM HeroesSuit:", heroId)
	return nil
}
func updateLOLMHeroesSuit(ctx *context.Context, heroId string, heroTech *dto.HeroTech, equipTech map[string]*dto.EquipTech) error {
	platform := common.PlatformForLOLM
	now := time.Now().Format("2006-01-02 15:04:05")
	hpd := dao.NewHeroesSuitDAO()
	var m model.HeroesSuit

	// 构建入库数据
	hsdata := make([]*model.HeroesSuit, 0)
	for _, eqs := range heroTech.Data.AnchorRecommend.List {
		et := equipTech[eqs.Head.Id]
		desc := make([]string, 0)
		desc = append(desc, fmt.Sprintf("<h4>%s</h4>", eqs.Body.Desc.Title))
		desc = append(desc, fmt.Sprintf("<p>%s</p>", eqs.Body.Desc.Content))
		desc = append(desc, fmt.Sprintf("<h5>%s</h5>", et.Data.ThinkingInfo.Title))
		for _, item := range et.Data.ThinkingInfo.List {
			desc = append(desc, fmt.Sprintf("<h6>%s</h6>", item.Name))
			desc = append(desc, fmt.Sprintf("<p>%s</p>", item.Content))
		}

		skillids := make([]string, 0)
		for _, item := range et.Data.SkillInfo.List {
			skillids = append(skillids, item.Id)
		}
		runeids := make([]string, 0)
		for _, nl := range et.Data.RuneInfo.NewList {
			for _, item := range nl.Items {
				runeids = append(runeids, item.Id)
			}
		}

		// 备战推荐 => 认作是LOL中的Other
		itemids := make([]string, 0)
		for _, equip := range et.Data.EquipInfo.List {
			itemids = append(itemids, equip.Id)
		}
		// 还有几个可选装备，也放到这里（和备战推荐可能有重复，需要去重）
		for _, l := range et.Data.EquipList {
			if strings.Contains(l.Title, "可选装备") {
				for _, eq := range l.List {
					if !inArray(eq.Id, itemids) {
						itemids = append(itemids, eq.Id)
					}
				}
			}
		}
		hsdata = append(hsdata, &model.HeroesSuit{
			HeroId:      eqs.Head.HeroId,
			Title:       et.Data.TopInfo.Title,
			RecommendId: eqs.Head.Id,
			Runeids:     strings.Join(runeids, ","),
			Skillids:    strings.Join(skillids, ","),
			Desc:        strings.Join(desc, ""),
			Author:      et.Data.TopInfo.Author,
			AuthorIcon:  et.Data.TopInfo.AuthorIcon,
			//Pos:         common.PositionNameEN[0],
			Pos:      et.Data.TopInfo.Title,
			Itemids:  strings.Join(itemids, ","),
			Type:     m.TypeOther(), // 契合装备
			Platform: platform,
			Version:  now,
			FileTime: now,
		})

		for _, l := range et.Data.EquipList {
			itemids = itemids[:0]
			if strings.Contains(l.Title, "鞋子推荐") || strings.Contains(l.Title, "附魔推荐") {
				// 鞋子+附魔 => 认作是LOL中的shoe
				for _, eq := range l.List {
					itemids = append(itemids, eq.Id)
				}

				hsdata = append(hsdata, &model.HeroesSuit{
					HeroId:      eqs.Head.HeroId,
					Title:       et.Data.TopInfo.Title,
					RecommendId: eqs.Head.Id,
					Runeids:     strings.Join(runeids, ","),
					Skillids:    strings.Join(skillids, ","),
					Desc:        strings.Join(desc, ""),
					Author:      et.Data.TopInfo.Author,
					AuthorIcon:  et.Data.TopInfo.AuthorIcon,
					//Pos:         common.PositionNameEN[0],
					Pos:      et.Data.TopInfo.Title,
					Itemids:  strings.Join(itemids, ","),
					Type:     m.TypeShoes(), // 鞋子装备
					Platform: platform,
					Version:  now,
					FileTime: now,
				})
			}
		}

		// 核心出装中有概率出现重复的，会导致唯一索引报错,这里用map去重
		itemidsMap := make(map[string]*model.HeroesSuit)
		for _, l := range et.Data.EquipList {
			itemids = itemids[:0]
			if strings.Contains(l.Title, "核心出装") {
				// 核心出装[1|2|3...] => 认作是LOL中的core
				for _, eq := range l.List {
					itemids = append(itemids, eq.Id)
				}
				key := strings.Join(itemids, ",")
				itemidsMap[key] = &model.HeroesSuit{
					HeroId:      eqs.Head.HeroId,
					Title:       et.Data.TopInfo.Title,
					RecommendId: eqs.Head.Id,
					Runeids:     strings.Join(runeids, ","),
					Skillids:    strings.Join(skillids, ","),
					Desc:        strings.Join(desc, ""),
					Author:      et.Data.TopInfo.Author,
					AuthorIcon:  et.Data.TopInfo.AuthorIcon,
					//Pos:         common.PositionNameEN[0],
					Pos:      et.Data.TopInfo.Title,
					Itemids:  strings.Join(itemids, ","),
					Type:     m.TypeCore(), // 核心装备
					Platform: platform,
					Version:  now,
					FileTime: now,
				}
			}
		}
		for _, equip := range itemidsMap {
			hsdata = append(hsdata, equip)
		}
	}

	if len(hsdata) == 0 {
		log.Logger.Warn(ctx, "hsdata is nil", "heroId:", heroId)
		return nil
	}

	err := hpd.DeleteAndInsert(map[string]interface{}{
		"heroId": heroId,
	}, hsdata)
	if err != nil {
		return errors.New("Add LOLM HeroesSuit " + err.Error() + ",heroId:" + heroId)
	}
	log.Logger.Info(ctx, "Add LOLM HeroesSuit:", heroId)

	return nil
}

func SuitData2Redis(ctx *context.Context) error {
	err := heroesSuits2Redis(ctx)
	if err != nil {
		return err
	}

	return nil
}

func heroesSuits2Redis(ctx *context.Context) error {
	hd := dao.NewHeroAttributeDAO()
	heroes, err := hd.Find([]string{
		"DISTINCT(heroId)", "name", "title", "platform",
	}, nil)
	if err != nil {
		return err
	}

	// LOL
	// 获取全部装备
	ed := dao.NewLOLEquipmentDAO()
	eVersion, err := ed.GetLOLEquipmentMaxVersion()
	if err != nil {
		return err
	}
	equips, err := ed.GetLOLEquipment(eVersion.Version)
	if err != nil {
		return err
	}

	mequip := make(map[string]*model.LOLEquipment)
	for _, equip := range equips {
		key := fmt.Sprintf(redis.KeyCacheEquip, equip.Maps, strconv.Itoa(common.PlatformForLOL), equip.ItemId)
		value, _ := json.Marshal(equip)
		mequip[key] = equip
		redis.RDB.Set(ctx, key, value, redis2.KeepTTL)
	}

	// 获取全部符文
	rd := dao.NewLOLRuneDAO()
	runes, err := rd.Find([]string{
		"*",
	}, map[string]interface{}{
		"status": 0,
	})
	if err != nil {
		return err
	}
	mrune := make(map[string]*model.LOLRune)
	for _, lolRune := range runes {
		key := fmt.Sprintf(redis.KeyCacheRune, "召唤师峡谷", strconv.Itoa(common.PlatformForLOL), lolRune.RuneID)
		value, _ := json.Marshal(lolRune)
		mrune[key] = lolRune
		redis.RDB.Set(ctx, key, value, redis2.KeepTTL)
	}

	// 获取全部召唤师技能
	sk := dao.NewLOLSkillDAO()
	skills, err := sk.Find([]string{
		"*",
	}, map[string]interface{}{
		"status": 0,
	})
	if err != nil {
		return err
	}
	mskill := make(map[string]*model.LOLSkill)
	for _, lolskill := range skills {
		key := fmt.Sprintf(redis.KeyCacheSkill, "召唤师峡谷", strconv.Itoa(common.PlatformForLOL), lolskill.SkillID)
		value, _ := json.Marshal(lolskill)
		mskill[key] = lolskill
		redis.RDB.Set(ctx, key, value, redis2.KeepTTL)
	}

	// LOLM
	// 获取全部装备
	med := dao.NewLOLMEquipmentDAO()
	meVersion, err := med.GetLOLMEquipmentMaxVersion()
	if err != nil {
		return err
	}
	mequips, err := med.GetLOLMEquipment(meVersion.Version)
	if err != nil {
		return err
	}

	mmequip := make(map[string]*model.LOLMEquipment)
	for _, equip := range mequips {
		key := fmt.Sprintf(redis.KeyCacheEquip, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), equip.EquipId) // todo
		value, _ := json.Marshal(equip)
		mmequip[key] = equip
		redis.RDB.Set(ctx, key, value, redis2.KeepTTL)
	}

	// 获取全部符文
	mrd := dao.NewLOLMRuneDAO()
	version, err := mrd.GetLOLMRuneMaxVersion()
	if err != nil {
		return err
	}
	mrunes, err := mrd.GetLOLMRune(version.Version)
	//mrunes, err := mrd.Find([]string{
	//	"*",
	//}, map[string]interface{}{
	//	"status": 0,
	//})
	if err != nil {
		return err
	}
	mmrune := make(map[string]*model.LOLMRune)
	for _, lolmRune := range mrunes {
		key := fmt.Sprintf(redis.KeyCacheRune, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), lolmRune.RuneId)
		value, _ := json.Marshal(lolmRune)
		mmrune[key] = lolmRune
		redis.RDB.Set(ctx, key, value, redis2.KeepTTL)
	}

	// 获取全部召唤师技能
	msk := dao.NewLOLMSkillDAO()
	mskills, err := msk.Find([]string{
		"*",
	}, map[string]interface{}{
		"status": 0,
	})
	if err != nil {
		return err
	}
	mmskill := make(map[string]*model.LOLMSkill)
	for _, lolmskill := range mskills {
		// todo maps "召唤师峡谷"
		key := fmt.Sprintf(redis.KeyCacheSkill, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), lolmskill.SkillID)
		value, _ := json.Marshal(lolmskill)
		mmskill[key] = lolmskill
		redis.RDB.Set(ctx, key, value, redis2.KeepTTL)
	}

	// --------------------------------
	sd := dao.NewHeroesSuitDAO()

	cancelCtx, cancelFunc := context2.WithCancel(ctx)
	defer cancelFunc()

	var (
		taskAll  = int32(len(heroes))
		taskSucc = int32(0)
		taskFail = int32(0)
		taskDone = int32(0)
		wg       = &sync.WaitGroup{}
		ch       = make(chan struct{}, 10)
	)

	for i, hero := range heroes {
		select {
		case <-cancelCtx.Done():
			break
		default:
			log.Logger.Info(ctx, ">>>>>>>>>>开始处理 hero:<<<<<<<<<<<", i, "/", hero.HeroId)
			ch <- struct{}{}
			wg.Add(1)

			go func(hero *model.HeroAttribute) {
				defer func() {
					<-ch
					wg.Done()
					atomic.AddInt32(&taskDone, 1)
				}()

				equipForHero, err2 := sd.GetSuitForHero(hero.Platform, hero.HeroId)
				if err2 != nil {
					atomic.AddInt32(&taskFail, 1)
					cancelFunc()
					log.Logger.Error(ctx, err2)
					return
				} else {

					hsm := make(map[string][]*model.HeroesSuit)
					for idx, equip := range equipForHero {
						hsm[equip.Pos] = append(hsm[equip.Pos], equipForHero[idx])
					}
					mhs := model.HeroesSuit{}

					eqs := make(map[string]dto.RecommendSuitEquip)
					for pos, posdata := range hsm {
						out := make([][]*dto.SuitData, 0)   // 出门装
						shoe := make([][]*dto.SuitData, 0)  // 鞋子
						core := make([][]*dto.SuitData, 0)  // 核心套件
						other := make([][]*dto.SuitData, 0) // 其他适配

						var (
							rrune [][]*dto.SuitData // 符文
							skill [][]*dto.SuitData // 召唤师技能
						)

						for _, data := range posdata {
							// 符文 和 召唤师技能 在同一个posdata下是一样的，这里随机取一个有值的
							if rrune == nil {
								ids := strings.Split(data.Runeids, ",")
								rune2 := make([]*dto.SuitData, 0)
								for _, id := range ids {
									if data.Platform == common.PlatformForLOL {
										// 端游
										key := fmt.Sprintf(redis.KeyCacheRune, "召唤师峡谷", strconv.Itoa(common.PlatformForLOL), id)
										if _, ok := mrune[key]; ok {
											rune2 = append(rune2, &dto.SuitData{
												ID:        cast.ToInt(id),
												Name:      mrune[key].Name,
												Icon:      mrune[key].Icon,
												Maps:      "召唤师峡谷",
												Plaintext: mrune[key].Tooltip,
												Desc:      fmt.Sprintf("<short>%s</short><long>%s</long>", mrune[key].Shortdesc, mrune[key].Longdesc),
												Version:   mrune[key].Version,
												RuneType:  mrune[key].StyleName,

												Igamecnt: data.Igamecnt,
												Wincnt:   data.Wincnt,
												Winrate:  data.Winrate,
												Allcnt:   data.Allcnt,
												Showrate: data.Showrate,

												Platform: data.Platform,
											})
										}
									} else {
										// 手游
										key := fmt.Sprintf(redis.KeyCacheRune, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), id)
										if _, ok := mmrune[key]; ok {
											rune2 = append(rune2, &dto.SuitData{
												ID:           cast.ToInt(id),
												Name:         mmrune[key].Name,
												Icon:         mmrune[key].IconPath,
												Plaintext:    fmt.Sprintf("<short>%s</short><long>%s</long>", mmrune[key].AttrName, mmrune[key].Description),
												Desc:         mmrune[key].DetailInfo,
												Version:      mmrune[key].Version,
												Igamecnt:     data.Igamecnt,
												Wincnt:       data.Wincnt,
												Winrate:      data.Winrate,
												Allcnt:       data.Allcnt,
												Showrate:     data.Showrate,
												Title:        data.Title,
												Author:       data.Author,
												AuthorIcon:   data.AuthorIcon,
												RecommendID:  data.RecommendId,
												ThinkingInfo: data.Desc,
												RuneType:     mmrune[key].StyleName,

												Platform: data.Platform,
											})
										}
									}
								}

								if len(rune2) > 0 {
									if rune2[0].Platform == common.PlatformForLOLM && len(rune2) == 5 {
										rrune = append(rrune, []*dto.SuitData{rune2[0]})
										rrune = append(rrune, []*dto.SuitData{rune2[1], rune2[2], rune2[3]})
										rrune = append(rrune, []*dto.SuitData{rune2[4]})
									} else {
										rrune = append(rrune, rune2)
									}
								}
							}
							if skill == nil {
								ids := strings.Split(data.Skillids, ",")
								skill2 := make([]*dto.SuitData, 0)
								for _, id := range ids {
									if data.Platform == common.PlatformForLOL {
										// 端游
										key := fmt.Sprintf(redis.KeyCacheSkill, "召唤师峡谷", strconv.Itoa(common.PlatformForLOL), id)
										if _, ok := mskill[key]; ok {
											skill2 = append(skill2, &dto.SuitData{
												ID:        cast.ToInt(id),
												Name:      mskill[key].Name,
												Icon:      mskill[key].Icon,
												Maps:      "召唤师峡谷",
												CD:        cast.ToInt(mskill[key].Cooldown),
												Plaintext: mskill[key].Description,
												Desc:      mskill[key].Description,
												Version:   mskill[key].Version,

												Igamecnt: data.Igamecnt,
												Wincnt:   data.Wincnt,
												Winrate:  data.Winrate,
												Allcnt:   data.Allcnt,
												Showrate: data.Showrate,
											})
										}
									} else {
										// 手游
										key := fmt.Sprintf(redis.KeyCacheSkill, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), id)
										if _, ok := mmskill[key]; ok {
											skill2 = append(skill2, &dto.SuitData{
												ID:           cast.ToInt(id),
												Name:         mmskill[key].Name,
												Icon:         mmskill[key].IconPath,
												Desc:         mmskill[key].FuncDesc,
												CD:           cast.ToInt(mmskill[key].Cd),
												Version:      mmskill[key].Version,
												Igamecnt:     data.Igamecnt,
												Wincnt:       data.Wincnt,
												Winrate:      data.Winrate,
												Allcnt:       data.Allcnt,
												Showrate:     data.Showrate,
												Title:        data.Title,
												Author:       data.Author,
												AuthorIcon:   data.AuthorIcon,
												RecommendID:  data.RecommendId,
												ThinkingInfo: data.Desc,
											})
										}
									}
								}
								if len(skill2) > 0 {
									skill = append(skill, skill2)
								}
							}
							switch data.Type {
							case mhs.TypeShoes():
								ids := strings.Split(data.Itemids, ",")
								shoe2 := make([]*dto.SuitData, 0)
								for _, id := range ids {
									if data.Platform == common.PlatformForLOL {
										// 端游
										key := fmt.Sprintf(redis.KeyCacheEquip, "召唤师峡谷", strconv.Itoa(common.PlatformForLOL), id)
										if _, ok := mequip[key]; !ok {
											continue
										}
										shoe2 = append(shoe2, &dto.SuitData{
											ID:        cast.ToInt(id),
											Name:      mequip[key].Name,
											Icon:      mequip[key].IconPath,
											Maps:      mequip[key].Maps,
											Plaintext: mequip[key].Plaintext,
											Desc:      mequip[key].Description,
											Price:     cast.ToInt(mequip[key].Total),
											Sell:      cast.ToInt(mequip[key].Sell),
											Version:   mequip[key].Version,

											Igamecnt: data.Igamecnt,
											Wincnt:   data.Wincnt,
											Winrate:  data.Winrate,
											Allcnt:   data.Allcnt,
											Showrate: data.Showrate,
										})
									} else if data.Platform == common.PlatformForLOLM {
										// 手游
										key := fmt.Sprintf(redis.KeyCacheEquip, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), id)
										if _, ok := mmequip[key]; !ok {
											continue
										}
										shoe2 = append(shoe2, &dto.SuitData{
											ID:           cast.ToInt(id),
											Name:         mmequip[key].Name,
											Icon:         mmequip[key].IconPath,
											Desc:         mmequip[key].Description,
											Price:        cast.ToInt(mmequip[key].Price),
											Version:      mmequip[key].Version,
											Igamecnt:     data.Igamecnt,
											Wincnt:       data.Wincnt,
											Winrate:      data.Winrate,
											Allcnt:       data.Allcnt,
											Showrate:     data.Showrate,
											Title:        data.Title,
											Author:       data.Author,
											AuthorIcon:   data.AuthorIcon,
											RecommendID:  data.RecommendId,
											ThinkingInfo: data.Desc,
										})
									} else {
										continue
									}
								}

								if len(shoe2) > 0 {
									shoe = append(shoe, shoe2)
								}
							case mhs.TypeOther():
								ids := strings.Split(data.Itemids, ",")
								other2 := make([]*dto.SuitData, 0)
								for _, id := range ids {
									if data.Platform == common.PlatformForLOL {
										// 端游
										key := fmt.Sprintf(redis.KeyCacheEquip, "召唤师峡谷", strconv.Itoa(common.PlatformForLOL), id)
										if _, ok := mequip[key]; !ok {
											continue
										}
										other2 = append(other2, &dto.SuitData{
											ID:        cast.ToInt(id),
											Name:      mequip[key].Name,
											Icon:      mequip[key].IconPath,
											Maps:      mequip[key].Maps,
											Plaintext: mequip[key].Plaintext,
											Desc:      mequip[key].Description,
											Price:     cast.ToInt(mequip[key].Total),
											Sell:      cast.ToInt(mequip[key].Sell),
											Version:   mequip[key].Version,

											Igamecnt: data.Igamecnt,
											Wincnt:   data.Wincnt,
											Winrate:  data.Winrate,
											Allcnt:   data.Allcnt,
											Showrate: data.Showrate,
										})
									} else if data.Platform == common.PlatformForLOLM {
										// 手游
										key := fmt.Sprintf(redis.KeyCacheEquip, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), id)
										if _, ok := mmequip[key]; !ok {
											continue
										}
										other2 = append(other2, &dto.SuitData{
											ID:      cast.ToInt(id),
											Name:    mmequip[key].Name,
											Icon:    mmequip[key].IconPath,
											Desc:    mmequip[key].Description,
											Price:   cast.ToInt(mmequip[key].Price),
											Version: mmequip[key].Version,

											Igamecnt: data.Igamecnt,
											Wincnt:   data.Wincnt,
											Winrate:  data.Winrate,
											Allcnt:   data.Allcnt,
											Showrate: data.Showrate,

											Title:        data.Title,
											Author:       data.Author,
											AuthorIcon:   data.AuthorIcon,
											RecommendID:  data.RecommendId,
											ThinkingInfo: data.Desc,
										})
									} else {
										continue
									}
								}

								if len(other2) > 0 {
									other = append(other, other2)
								}

							case mhs.TypeOut():
								ids := strings.Split(data.Itemids, ",")
								out2 := make([]*dto.SuitData, 0)
								for _, id := range ids {
									if data.Platform == common.PlatformForLOL {
										// 端游
										key := fmt.Sprintf(redis.KeyCacheEquip, "召唤师峡谷", strconv.Itoa(common.PlatformForLOL), id)
										if _, ok := mequip[key]; !ok {
											continue
										}
										out2 = append(out2, &dto.SuitData{
											ID:        cast.ToInt(id),
											Name:      mequip[key].Name,
											Icon:      mequip[key].IconPath,
											Maps:      mequip[key].Maps,
											Plaintext: mequip[key].Plaintext,
											Desc:      mequip[key].Description,
											Price:     cast.ToInt(mequip[key].Total),
											Sell:      cast.ToInt(mequip[key].Sell),
											Version:   mequip[key].Version,

											Igamecnt: data.Igamecnt,
											Wincnt:   data.Wincnt,
											Winrate:  data.Winrate,
											Allcnt:   data.Allcnt,
											Showrate: data.Showrate,
										})
									} else if data.Platform == common.PlatformForLOLM {
										// 手游
										key := fmt.Sprintf(redis.KeyCacheEquip, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), id)
										if _, ok := mmequip[key]; !ok {
											continue
										}

										out2 = append(out2, &dto.SuitData{
											ID:      cast.ToInt(id),
											Name:    mmequip[key].Name,
											Icon:    mmequip[key].IconPath,
											Desc:    mmequip[key].Description,
											Price:   cast.ToInt(mmequip[key].Price),
											Version: mmequip[key].Version,

											Igamecnt: data.Igamecnt,
											Wincnt:   data.Wincnt,
											Winrate:  data.Winrate,
											Allcnt:   data.Allcnt,
											Showrate: data.Showrate,

											Title:        data.Title,
											Author:       data.Author,
											AuthorIcon:   data.AuthorIcon,
											RecommendID:  data.RecommendId,
											ThinkingInfo: data.Desc,
										})
									} else {
										continue
									}
								}

								if len(out2) > 0 {
									out = append(out, out2)
								}

							case mhs.TypeCore():
								ids := strings.Split(data.Itemids, ",")
								core2 := make([]*dto.SuitData, 0)
								for _, id := range ids {
									if data.Platform == common.PlatformForLOL {
										// 端游
										key := fmt.Sprintf(redis.KeyCacheEquip, "召唤师峡谷", strconv.Itoa(common.PlatformForLOL), id)
										if _, ok := mequip[key]; !ok {
											continue
										}
										core2 = append(core2, &dto.SuitData{
											ID:        cast.ToInt(id),
											Name:      mequip[key].Name,
											Icon:      mequip[key].IconPath,
											Maps:      mequip[key].Maps,
											Plaintext: mequip[key].Plaintext,
											Desc:      mequip[key].Description,
											Price:     cast.ToInt(mequip[key].Total),
											Sell:      cast.ToInt(mequip[key].Sell),
											Version:   mequip[key].Version,

											Igamecnt: data.Igamecnt,
											Wincnt:   data.Wincnt,
											Winrate:  data.Winrate,
											Allcnt:   data.Allcnt,
											Showrate: data.Showrate,
										})
									} else if data.Platform == common.PlatformForLOLM {
										// 手游
										key := fmt.Sprintf(redis.KeyCacheEquip, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), id)
										if _, ok := mmequip[key]; !ok {
											continue
										}
										core2 = append(core2, &dto.SuitData{
											ID:      cast.ToInt(id),
											Name:    mmequip[key].Name,
											Icon:    mmequip[key].IconPath,
											Desc:    mmequip[key].Description,
											Price:   cast.ToInt(mmequip[key].Price),
											Version: mmequip[key].Version,

											Igamecnt: data.Igamecnt,
											Wincnt:   data.Wincnt,
											Winrate:  data.Winrate,
											Allcnt:   data.Allcnt,
											Showrate: data.Showrate,

											Title:        data.Title,
											Author:       data.Author,
											AuthorIcon:   data.AuthorIcon,
											RecommendID:  data.RecommendId,
											ThinkingInfo: data.Desc,
										})
									} else {
										continue
									}
								}

								if len(core2) > 0 {
									core = append(core, core2)
								}

							}
						}

						eqs[pos] = dto.RecommendSuitEquip{
							Out:   out,
							Shoe:  shoe,
							Core:  core,
							Other: other,
							Rune:  rrune,
							Skill: skill,
						}
					}

					jsonData, _ := json.Marshal(eqs)
					redis.RDB.HSet(ctx, redis.KeyCacheHeroEquip, hero.HeroId, jsonData)
					atomic.AddInt32(&taskSucc, 1)
				}
			}(hero)
		}
	}

	wg.Wait()

	log.Logger.Info(ctx, fmt.Sprintf("处理了: %d 个任务", taskDone))
	log.Logger.Info(ctx, fmt.Sprintf("提前结束,执行出错: %d 个任务", taskFail))
	log.Logger.Info(ctx, fmt.Sprintf("成功执行了: %d 个任务", taskSucc))
	log.Logger.Info(ctx, fmt.Sprintf("剩余: %d 个任务待处理", taskAll-taskDone))

	return nil
}

func GetHeroSuit(ctx *context.Context, heroID string) (dto.HeroSuit, error) {
	d := redis.RDB.HGet(ctx, redis.KeyCacheHeroEquip, heroID)
	hs := dto.HeroSuit{
		HeroID: heroID,
		ExtInfo: dto.HeroSuitExtInfo{
			RecommendReason: make(map[string]string),
			AuthorInfo: make(map[string]struct {
				Name string `json:"name"`
				Icon string `json:"icon"`
			}),
		},
	}

	hid := cast.ToInt(heroID)
	if hid < common.MinHeroIDForLOLM {
		hs.Platform = common.PlatformForLOL
	} else {
		hs.Platform = common.PlatformForLOLM
	}

	var rs map[string]dto.RecommendSuitEquip
	err := json.Unmarshal([]byte(d.Val()), &rs)
	hs.Equips = rs

	for title, data := range rs {
		var mTypeEquips map[string][][]*dto.SuitData
		marshal, _ := json.Marshal(data)
		_ = json.Unmarshal(marshal, &mTypeEquips)

		for _, equips := range mTypeEquips {
			if len(equips) == 0 {
				continue
			}

			for _, suitData := range equips {
				for _, equip := range suitData {
					hs.ExtInfo.RecommendReason[title] = equip.ThinkingInfo
					hs.ExtInfo.AuthorInfo[title] = struct {
						Name string `json:"name"`
						Icon string `json:"icon"`
					}{
						Name: equip.Author,
						Icon: equip.AuthorIcon,
					}
					break // 只执行一次
				}
				break // 只执行一次
			}

			break // 只执行一次
		}
	}

	return hs, err
}

func HeroesPosition(ctx *context.Context, platform int) (*dto.HeroRankList, error) {
	rankList, err := service.HeroRankList(ctx)
	if err != nil {
		log.Logger.Error(ctx, err)
		return nil, err
	}

	hpd := dao.NewHeroesPositionDAO()
	// 删除旧数据
	cond := map[string]interface{}{
		"platform": common.PlatformForLOLM,
	}

	hp := make([]*model.HeroesPosition, 0)
	// 只取钻石以上分段
	if levData, ok := rankList.Data[common.LevelDiamond]; ok {
		for pos, heroes := range levData {
			posName := common.PositionNameEN[pos]
			for _, data := range heroes {

				hp = append(hp, &model.HeroesPosition{
					HeroId:   data.HeroId,
					Pos:      posName,
					ShowRate: utils.Str2Int(data.AppearRate),
					WinRate:  utils.Str2Int(data.WinRate),
					Platform: common.PlatformForLOLM,
					Version:  data.Dtstatdate,
					FileTime: data.Dtstatdate,
				})
			}
		}
	}

	err = hpd.DeleteAndInsert(cond, hp)
	if err != nil {
		log.Logger.Error(ctx, err)
		return nil, err
	}
	log.Logger.Info(ctx, "add LOLM position success")
	return rankList, nil
}

func inArray(id string, ids []string) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}

func SuitHeroData2Redis(ctx *context.Context) error {
	err := suitHero2Redis(ctx)
	if err != nil {
		return err
	}

	return nil
}

func suitHero2Redis(ctx *context.Context) error {
	hsd := dao.NewHeroesSuitDAO()
	list, err := hsd.FindHighRateEquip([]string{
		"heroId", "itemids", "skillids", "runeids", "winrate", "showrate", "platform", "author", "version", "fileTime",
	}, nil)
	if err != nil {
		return err
	}

	var (
		ch = make(chan struct{}, 5000)
		wg = sync.WaitGroup{}
	)

	for _, hero := range list {
		wg.Add(1)
		ch <- struct{}{}

		go func(hero *model.HeroesSuit) {
			defer func() {
				wg.Done()
				<-ch
			}()

			// equip
			if hero.Itemids != "" {
				ids := strings.Split(hero.Itemids, ",")
				for _, id := range ids {
					key := fmt.Sprintf(redis.KeyCacheEquipHeroSuit, hero.Platform, id)
					redis.RDB.ZIncrBy(ctx, key, 1, hero.HeroId)
				}
			}

			// rune
			if hero.Runeids != "" {
				ids := strings.Split(hero.Runeids, ",")
				for _, id := range ids {
					key := fmt.Sprintf(redis.KeyCacheRuneHeroSuit, hero.Platform, id)
					redis.RDB.ZIncrBy(ctx, key, 1, hero.HeroId)
				}
			}

			// skill
			if hero.Skillids != "" {
				ids := strings.Split(hero.Skillids, ",")
				for _, id := range ids {
					key := fmt.Sprintf(redis.KeyCacheSkillHeroSuit, hero.Platform, id)
					redis.RDB.ZIncrBy(ctx, key, 1, hero.HeroId)
				}
			}
		}(hero)

	}

	wg.Wait()
	log.Logger.Info(ctx, "suitHero2Redis ok")
	return nil
}

func GetEquipHeroSuit(ctx *context.Context, platform int, equipID string) ([]*dto.SearchResultList, error) {
	// 获取英雄适配数据
	suitHeroes := make([]*dto.SearchResultList, 0)

	min := "-inf"
	max := "+inf"
	// ZREVRANGE my_rankings 0 2 WITHSCORES
	key := fmt.Sprintf(redis.KeyCacheEquipHeroSuit, platform, equipID)
	score := redis.RDB.ZRevRangeByScoreWithScores(ctx, key, &redis2.ZRangeBy{
		Min: min,
		Max: max,
		//Offset: 0,
		//Count:  100,
	})

	var heroesID []string
	for _, k := range score.Val() {
		heroesID = append(heroesID, k.Member.(string))
	}

	if len(heroesID) == 0 {
		return suitHeroes, nil
	}

	had := dao.NewHeroAttributeDAO()
	heroes, err := had.Find([]string{
		"heroId", "name", "title", "avatar", "platform", "version",
	}, map[string]interface{}{
		"heroId": heroesID,
	})
	if err != nil {
		return nil, err
	}

	for _, hero := range heroes {
		name := ""
		if platform == common.PlatformForLOL {
			name = hero.Name + " " + hero.Title
		} else {
			name = hero.Title + " " + hero.Name
		}
		suitHeroes = append(suitHeroes, &dto.SearchResultList{
			Id:       hero.HeroId,
			Name:     name,
			Icon:     hero.Avatar,
			Platform: hero.Platform,
			Version:  hero.Version,
		})
	}

	return suitHeroes, nil
}

func GetRuneHeroSuit(ctx *context.Context, platform int, runeID string) ([]*dto.SearchResultList, error) {
	// 获取英雄适配数据
	suitHeroes := make([]*dto.SearchResultList, 0)

	min := "-inf"
	max := "+inf"
	// ZREVRANGE my_rankings 0 2 WITHSCORES
	key := fmt.Sprintf(redis.KeyCacheRuneHeroSuit, platform, runeID)
	score := redis.RDB.ZRevRangeByScoreWithScores(ctx, key, &redis2.ZRangeBy{
		Min: min,
		Max: max,
		//Offset: 0,
		//Count:  100,
	})

	var heroesID []string
	for _, k := range score.Val() {
		heroesID = append(heroesID, k.Member.(string))
	}

	if len(heroesID) == 0 {
		return suitHeroes, nil
	}

	had := dao.NewHeroAttributeDAO()
	heroes, err := had.Find([]string{
		"heroId", "name", "title", "avatar", "platform", "version",
	}, map[string]interface{}{
		"heroId": heroesID,
	})
	if err != nil {
		return nil, err
	}

	for _, hero := range heroes {
		name := ""
		if platform == common.PlatformForLOL {
			name = hero.Name + " " + hero.Title
		} else {
			name = hero.Title + " " + hero.Name
		}
		suitHeroes = append(suitHeroes, &dto.SearchResultList{
			Id:       hero.HeroId,
			Name:     name,
			Icon:     hero.Avatar,
			Platform: hero.Platform,
			Version:  hero.Version,
		})
	}

	return suitHeroes, nil
}

func GetSkillHeroSuit(ctx *context.Context, platform int, skillID string) ([]*dto.SearchResultList, error) {
	// 获取英雄适配数据
	suitHeroes := make([]*dto.SearchResultList, 0)

	min := "-inf"
	max := "+inf"
	// ZREVRANGE my_rankings 0 2 WITHSCORES
	key := fmt.Sprintf(redis.KeyCacheSkillHeroSuit, platform, skillID)
	score := redis.RDB.ZRevRangeByScoreWithScores(ctx, key, &redis2.ZRangeBy{
		Min: min,
		Max: max,
		//Offset: 0,
		//Count:  100,
	})

	var heroesID []string
	for _, k := range score.Val() {
		heroesID = append(heroesID, k.Member.(string))
	}

	if len(heroesID) == 0 {
		return suitHeroes, nil
	}

	had := dao.NewHeroAttributeDAO()
	heroes, err := had.Find([]string{
		"heroId", "name", "title", "avatar", "platform", "version",
	}, map[string]interface{}{
		"heroId": heroesID,
	})
	if err != nil {
		return nil, err
	}

	for _, hero := range heroes {
		name := ""
		if platform == common.PlatformForLOL {
			name = hero.Name + " " + hero.Title
		} else {
			name = hero.Title + " " + hero.Name
		}
		suitHeroes = append(suitHeroes, &dto.SearchResultList{
			Id:       hero.HeroId,
			Name:     name,
			Icon:     hero.Avatar,
			Platform: hero.Platform,
			Version:  hero.Version,
		})
	}

	return suitHeroes, nil
}
