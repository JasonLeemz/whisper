package logic

import (
	"encoding/json"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"strings"
	"sync"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/redis"
)

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
