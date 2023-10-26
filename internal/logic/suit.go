package logic

import (
	context2 "context"
	"encoding/json"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/redis"
)

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
