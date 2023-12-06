package suit

import (
	context2 "context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"whisper/internal/dto"
	"whisper/internal/logic/attribute"
	"whisper/internal/logic/common"
	"whisper/internal/logic/equipment"
	"whisper/internal/logic/position"
	rune2 "whisper/internal/logic/rune"
	"whisper/internal/logic/skill"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	lol "whisper/internal/service/lol"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/redis"
)

type task struct {
	total   int32
	success int32
	fail    int32
	done    int32
	wg      *sync.WaitGroup
	ch      chan struct{}
}

func newTask(total int32, wg *sync.WaitGroup, ch chan struct{}) *task {
	return &task{total: total, wg: wg, ch: ch}
}

type suit struct {
	ctx      *context.Context
	platform int
}

type Func func(ctx *context.Context, platform int) *suit

func NewSuit() Func {
	return func(ctx *context.Context, platform int) *suit {
		return &suit{ctx: ctx, platform: platform}
	}
}

func (suit *suit) BatchUpdateSuitEquip() {
	if suit.platform == common.PlatformForLOLM {
		// LOLM的position数据是单独接口获取的
		_, _ = position.NewPosition()(suit.ctx, common.PlatformForLOLM).HeroesPosition()
	}

	// 获取所有英雄ID
	ha := dao.NewHeroAttributeDAO()
	heroes, err := ha.QueryAllHeroes(nil)
	if err != nil {
		log.Logger.Error(suit.ctx, err)
		return
	}

	cancelCtx, cancelFunc := context2.WithCancel(suit.ctx)
	defer cancelFunc()

	t := newTask(int32(len(heroes)), &sync.WaitGroup{}, make(chan struct{}, 10))

	for _, hero := range heroes {
		select {
		case <-cancelCtx.Done():
			break
		default:
			t.ch <- struct{}{}
			t.wg.Add(1)

			go func(hero *model.HeroAttribute) {
				defer func() {
					<-t.ch
					t.wg.Done()
					atomic.AddInt32(&t.done, 1)
				}()

				suit.platform = hero.Platform
				_, err2 := suit.QuerySuitEquip(hero.HeroId)
				// 任务执行失败，这个地方可以使用锁，也可以使用原子操作，优先原子操作
				if err2 != nil {
					// todo 报警

					atomic.AddInt32(&t.fail, 1)
					cancelFunc()
					log.Logger.Error(suit.ctx, err2)
					return
				} else {
					atomic.AddInt32(&t.success, 1)
				}
			}(hero)
		}
	}

	t.wg.Wait()

	log.Logger.Info(suit.ctx, "批量更新英雄适配装备:")
	log.Logger.Info(suit.ctx, fmt.Sprintf("处理了: %d 个任务", t.done))
	log.Logger.Info(suit.ctx, fmt.Sprintf("提前结束,执行出错: %d 个任务", t.fail))
	log.Logger.Info(suit.ctx, fmt.Sprintf("成功执行了: %d 个任务", t.success))
	log.Logger.Info(suit.ctx, fmt.Sprintf("剩余: %d 个任务待处理", t.total-t.done))
}

func genPosData(heroId string, fightData *dto.ChampionFightData) []*model.HeroesPosition {
	posData := make([]*model.HeroesPosition, 0, 3)
	for pos, _ := range fightData.List.ChampionFight {
		posData = append(posData, &model.HeroesPosition{
			HeroId:   heroId,
			Pos:      pos,
			Platform: common.PlatformForLOL,
			Version:  fightData.GameVer,
			FileTime: fightData.Date,
		})
	}
	return posData
}
func (suit *suit) QuerySuitEquip(heroId string) (any, error) {

	if suit.platform == common.PlatformForLOL {
		fightData, err := suit.getFightData(heroId)
		if err != nil {
			return nil, errors.New("LOL英雄的rank数据/getFightData:" + err.Error())
		}

		// reload heroes_position 表
		err = position.NewPosition()(suit.ctx, suit.platform).
			UpdateHeroesPosition(map[string]interface{}{
				"heroId": heroId,
			}, genPosData(heroId, fightData),
			)
		if err != nil {
			return nil, errors.New("updateHeroesPosition:" + err.Error())
		}

		// reload heroes_suit 表
		err = updateLOLHeroesSuit(suit.ctx, heroId, fightData)
		if err != nil {
			return nil, errors.New("updateHeroesSuit:" + err.Error())
		}
		return fightData, nil
	} else {
		// common.PlatformForLOLM
		equips, err := lol.CreateLOLProduct(suit.platform)().QuerySuitEquip(suit.ctx, heroId)
		if err != nil {
			return nil, err
		}
		// reload heroes_suit 表
		err = updateLOLMHeroesSuit(suit.ctx, heroId, equips.(*dto.HeroTech))
		if err != nil {
			return nil, errors.New("updateLOLMHeroesSuit:" + err.Error())
		}

		return equips, nil
	}
}

// LOL英雄的rank数据
func (suit *suit) getFightData(heroId string) (*dto.ChampionFightData, error) {
	data, err := lol.CreateLOLProduct(suit.platform)().HeroRankData(suit.ctx, heroId)
	if err != nil {
		return nil, err
	}

	fightData := data.(*dto.ChampionFightData)
	for pos, posData := range fightData.List.ChampionLane {
		equipData := map[string]dto.Itemjson{}
		tmp := dto.ChampionLaneItem{}

		var err error

		if posData.Itemoutjson != "" {
			err = json.Unmarshal([]byte(posData.Itemoutjson), &equipData)
			if err != nil {
				log.Logger.Warn(suit.ctx, err, "heroid:", heroId)
			} else {
				tmp.Itemout = equipData
			}
		}

		if posData.Core3itemjson != "" {
			equipData = *new(map[string]dto.Itemjson)
			err = json.Unmarshal([]byte(posData.Core3itemjson), &equipData)
			if err != nil {
				log.Logger.Warn(suit.ctx, err, "heroid:", heroId)
			} else {
				tmp.Core3item = equipData
			}
		}

		if posData.Shoesjson != "" {
			equipData = *new(map[string]dto.Itemjson)
			err = json.Unmarshal([]byte(posData.Shoesjson), &equipData)
			if err != nil {
				log.Logger.Warn(suit.ctx, err, "heroid:", heroId)
			} else {
				tmp.Shoes = equipData
			}
		}

		if posData.Hold3 != "" {
			var suits []dto.Itemjson
			err = json.Unmarshal([]byte(posData.Hold3), &suits)
			if err != nil {
				log.Logger.Warn(suit.ctx, err, "heroid:", heroId)
			} else {
				tmp.Suits = suits
			}
		}

		fightData.List.ChampionLane[pos] = tmp
	}

	return fightData, nil
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
func updateLOLMHeroesSuit(ctx *context.Context, heroId string, heroTech *dto.HeroTech) error {
	platform := common.PlatformForLOLM
	now := time.Now().Format("2006-01-02 15:04:05")
	hpd := dao.NewHeroesSuitDAO()
	var m model.HeroesSuit

	// 构建入库数据
	hsdata := make([]*model.HeroesSuit, 0)
	for _, eqs := range heroTech.Data.AnchorRecommend.List {
		et := heroTech.EquipData[eqs.Head.Id]
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
func inArray(id string, ids []string) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}

// SuitData2Redis 批量更新英雄适配的装备、天赋、技能到redis
var (
	mEquip4LOL  map[string]*model.LOLEquipment
	mEquip4LOLM map[string]*model.LOLMEquipment

	mRune4LOL  map[string]*model.LOLRune
	mRune4LOLM map[string]*model.LOLMRune

	mSkill4LOL  map[string]*model.LOLSkill
	mSkill4LOLM map[string]*model.LOLMSkill
)

func SuitData2Redis(ctx *context.Context) error {
	attrs := attribute.NewInnerIns(ctx).GetAll()
	mEquip4LOL = equipment.NewInnerIns(ctx).GetAll(common.PlatformForLOL).(map[string]*model.LOLEquipment)
	mEquip4LOLM = equipment.NewInnerIns(ctx).GetAll(common.PlatformForLOLM).(map[string]*model.LOLMEquipment)

	mRune4LOL = rune2.NewInnerIns(ctx).GetAll(common.PlatformForLOL).(map[string]*model.LOLRune)
	mRune4LOLM = rune2.NewInnerIns(ctx).GetAll(common.PlatformForLOLM).(map[string]*model.LOLMRune)

	mSkill4LOL = skill.NewInnerIns(ctx).GetAll(common.PlatformForLOL).(map[string]*model.LOLSkill)
	mSkill4LOLM = skill.NewInnerIns(ctx).GetAll(common.PlatformForLOLM).(map[string]*model.LOLMSkill)

	log.Logger.Info(ctx,
		"attrs", len(attrs),
		"mEquip4LOL", len(mEquip4LOL),
		"mEquip4LOLM", len(mEquip4LOLM),
		"mRune4LOL", len(mRune4LOL),
		"mRune4LOLM", len(mRune4LOLM),
		"mSkill4LOL", len(mSkill4LOL),
		"mSkill4LOLM", len(mSkill4LOLM),
	)

	// --------------------------------
	sd := dao.NewHeroesSuitDAO()

	cancelCtx, cancelFunc := context2.WithCancel(ctx)
	defer cancelFunc()

	t := newTask(int32(len(attrs)), &sync.WaitGroup{}, make(chan struct{}, 10))

	for i, hero := range attrs {
		select {
		case <-cancelCtx.Done():
			break
		default:
			log.Logger.Info(ctx, ">>>>>>>>>>开始处理 hero:<<<<<<<<<<<", i, "/", hero.HeroId)
			t.ch <- struct{}{}
			t.wg.Add(1)

			go func(hero *model.HeroAttribute) {
				defer func() {
					<-t.ch
					t.wg.Done()
					atomic.AddInt32(&t.done, 1)
				}()

				equipForHero, err2 := sd.GetSuitForHero(hero.Platform, hero.HeroId)
				if err2 != nil {
					atomic.AddInt32(&t.fail, 1)
					cancelFunc()
					log.Logger.Error(ctx, "GetSuitForHero", err2)
					return
				} else {

					hsm := make(map[string][]*model.HeroesSuit)
					for idx, equip := range equipForHero {
						hsm[equip.Pos] = append(hsm[equip.Pos], equipForHero[idx])
					}
					mhs := model.HeroesSuit{}

					eqs := make(map[string]dto.RecommendSuitEquip)
					for pos, posdata := range hsm {
						suitForOut := make([][]*dto.SuitData, 0)   // 出门装
						suitForShoe := make([][]*dto.SuitData, 0)  // 鞋子
						suitForCore := make([][]*dto.SuitData, 0)  // 核心套件
						suitForOther := make([][]*dto.SuitData, 0) // 其他适配

						var (
							suitForRune  [][]*dto.SuitData // 符文
							suitForSkill [][]*dto.SuitData // 召唤师技能
						)

						for _, data := range posdata {
							// 符文 和 召唤师技能 在同一个posdata下是一样的，这里随机取一个有值的
							if suitForRune == nil && data.Runeids != "" {
								ids := strings.Split(data.Runeids, ",")
								suitData := make([]*dto.SuitData, 0)
								for _, id := range ids {
									suitData = append(suitData, genRuneData(data, id))
								}

								if len(suitData) > 0 {
									if suitData[0].Platform == common.PlatformForLOLM && len(suitData) == 5 {
										suitForRune = [][]*dto.SuitData{
											[]*dto.SuitData{suitData[0]},
											[]*dto.SuitData{suitData[1], suitData[2], suitData[3]},
											[]*dto.SuitData{suitData[4]},
										}
									} else {
										suitForRune = [][]*dto.SuitData{
											suitData,
										}
									}
								}
							}
							if suitForSkill == nil && data.Skillids != "" {
								ids := strings.Split(data.Skillids, ",")
								suitData := make([]*dto.SuitData, 0)
								for _, id := range ids {
									suitData = append(suitData, genSkillData(data, id))
								}
								if len(suitData) > 0 {
									suitForSkill = [][]*dto.SuitData{
										suitData,
									}
								}
							}
							switch data.Type {
							case mhs.TypeShoes():
								ids := strings.Split(data.Itemids, ",")
								suitData := make([]*dto.SuitData, 0)
								for _, id := range ids {
									if equipData := genEquipData(data, id); equipData != nil {
										suitData = append(suitData, equipData)
									}
								}

								if len(suitData) > 0 {
									suitForShoe = append(suitForShoe, suitData)
								}
							case mhs.TypeOther():
								ids := strings.Split(data.Itemids, ",")
								suitData := make([]*dto.SuitData, 0)
								for _, id := range ids {
									if equipData := genEquipData(data, id); equipData != nil {
										suitData = append(suitData, equipData)
									}
								}

								if len(suitData) > 0 {
									suitForOther = append(suitForOther, suitData)
								}

							case mhs.TypeOut():
								ids := strings.Split(data.Itemids, ",")
								suitData := make([]*dto.SuitData, 0)
								for _, id := range ids {
									if equipData := genEquipData(data, id); equipData != nil {
										suitData = append(suitData, equipData)
									}
								}

								if len(suitData) > 0 {
									suitForOut = append(suitForOut, suitData)
								}

							case mhs.TypeCore():
								ids := strings.Split(data.Itemids, ",")
								suitData := make([]*dto.SuitData, 0)
								for _, id := range ids {
									if equipData := genEquipData(data, id); equipData != nil {
										suitData = append(suitData, equipData)
									}
								}

								if len(suitData) > 0 {
									suitForCore = append(suitForCore, suitData)
								}

							}
						}

						eqs[pos] = dto.RecommendSuitEquip{
							Out:   suitForOut,
							Shoe:  suitForShoe,
							Core:  suitForCore,
							Other: suitForOther,
							Rune:  suitForRune,
							Skill: suitForSkill,
						}
					}

					jsonData, _ := json.Marshal(eqs)
					redis.RDB.HSet(ctx, redis.KeyCacheHeroEquip, hero.HeroId, jsonData)
					atomic.AddInt32(&t.success, 1)
				}
			}(hero)
		}
	}

	t.wg.Wait()

	log.Logger.Info(ctx, fmt.Sprintf("处理了: %d 个任务", t.done))
	log.Logger.Info(ctx, fmt.Sprintf("提前结束,执行出错: %d 个任务", t.fail))
	log.Logger.Info(ctx, fmt.Sprintf("成功执行了: %d 个任务", t.success))
	log.Logger.Info(ctx, fmt.Sprintf("剩余: %d 个任务待处理", t.total-t.done))

	return nil
}

func genEquipData(data *model.HeroesSuit, id string) *dto.SuitData {
	var suitData *dto.SuitData
	if data.Platform == common.PlatformForLOL {
		// 端游
		key := fmt.Sprintf(redis.KeyCacheEquip, "召唤师峡谷", strconv.Itoa(common.PlatformForLOL), id)
		if _, ok := mEquip4LOL[key]; !ok {
			return nil
		}
		suitData = &dto.SuitData{
			ID:        cast.ToInt(id),
			Name:      mEquip4LOL[key].Name,
			Icon:      mEquip4LOL[key].IconPath,
			Maps:      mEquip4LOL[key].Maps,
			Plaintext: mEquip4LOL[key].Plaintext,
			Desc:      mEquip4LOL[key].Description,
			Price:     cast.ToInt(mEquip4LOL[key].Total),
			Sell:      cast.ToInt(mEquip4LOL[key].Sell),
			Version:   mEquip4LOL[key].Version,

			Igamecnt: data.Igamecnt,
			Wincnt:   data.Wincnt,
			Winrate:  data.Winrate,
			Allcnt:   data.Allcnt,
			Showrate: data.Showrate,
		}
	} else if data.Platform == common.PlatformForLOLM {
		// 手游
		key := fmt.Sprintf(redis.KeyCacheEquip, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), id)
		if _, ok := mEquip4LOLM[key]; !ok {
			return nil
		}
		suitData = &dto.SuitData{
			ID:           cast.ToInt(id),
			Name:         mEquip4LOLM[key].Name,
			Icon:         mEquip4LOLM[key].IconPath,
			Desc:         mEquip4LOLM[key].Description,
			Price:        cast.ToInt(mEquip4LOLM[key].Price),
			Version:      mEquip4LOLM[key].Version,
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
		}
	}
	return suitData
}
func genRuneData(data *model.HeroesSuit, id string) *dto.SuitData {
	var suitData *dto.SuitData
	if data.Platform == common.PlatformForLOL {
		// 端游
		key := fmt.Sprintf(redis.KeyCacheRune, "召唤师峡谷", strconv.Itoa(common.PlatformForLOL), id)
		if _, ok := mRune4LOL[key]; ok {
			suitData = &dto.SuitData{
				ID:        cast.ToInt(id),
				Name:      mRune4LOL[key].Name,
				Icon:      mRune4LOL[key].Icon,
				Maps:      "召唤师峡谷",
				Plaintext: mRune4LOL[key].Tooltip,
				Desc:      fmt.Sprintf("<short>%s</short><long>%s</long>", mRune4LOL[key].Shortdesc, mRune4LOL[key].Longdesc),
				Version:   mRune4LOL[key].Version,
				RuneType:  mRune4LOL[key].StyleName,

				Igamecnt: data.Igamecnt,
				Wincnt:   data.Wincnt,
				Winrate:  data.Winrate,
				Allcnt:   data.Allcnt,
				Showrate: data.Showrate,

				Platform: data.Platform,
			}
		}
	} else {
		// 手游
		key := fmt.Sprintf(redis.KeyCacheRune, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), id)
		if _, ok := mRune4LOLM[key]; ok {
			suitData = &dto.SuitData{
				ID:           cast.ToInt(id),
				Name:         mRune4LOLM[key].Name,
				Icon:         mRune4LOLM[key].IconPath,
				Plaintext:    fmt.Sprintf("<short>%s</short><long>%s</long>", mRune4LOLM[key].AttrName, mRune4LOLM[key].Description),
				Desc:         mRune4LOLM[key].DetailInfo,
				Version:      mRune4LOLM[key].Version,
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
				RuneType:     mRune4LOLM[key].StyleName,

				Platform: data.Platform,
			}
		}
	}
	return suitData
}
func genSkillData(data *model.HeroesSuit, id string) *dto.SuitData {
	var suitData *dto.SuitData
	if data.Platform == common.PlatformForLOL {
		// 端游
		key := fmt.Sprintf(redis.KeyCacheSkill, "召唤师峡谷", strconv.Itoa(common.PlatformForLOL), id)
		if _, ok := mSkill4LOL[key]; ok {
			suitData = &dto.SuitData{
				ID:        cast.ToInt(id),
				Name:      mSkill4LOL[key].Name,
				Icon:      mSkill4LOL[key].Icon,
				Maps:      "召唤师峡谷",
				CD:        cast.ToInt(mSkill4LOL[key].Cooldown),
				Plaintext: mSkill4LOL[key].Description,
				Desc:      mSkill4LOL[key].Description,
				Version:   mSkill4LOL[key].Version,

				Igamecnt: data.Igamecnt,
				Wincnt:   data.Wincnt,
				Winrate:  data.Winrate,
				Allcnt:   data.Allcnt,
				Showrate: data.Showrate,
			}
		}
	} else {
		// 手游
		key := fmt.Sprintf(redis.KeyCacheSkill, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), id)
		if _, ok := mSkill4LOLM[key]; ok {
			suitData = &dto.SuitData{
				ID:           cast.ToInt(id),
				Name:         mSkill4LOLM[key].Name,
				Icon:         mSkill4LOLM[key].IconPath,
				Desc:         mSkill4LOLM[key].FuncDesc,
				CD:           cast.ToInt(mSkill4LOLM[key].Cd),
				Version:      mSkill4LOLM[key].Version,
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
			}
		}
	}

	return suitData
}
