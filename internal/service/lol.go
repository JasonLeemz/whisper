package service

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"
	"time"
	"whisper/internal/dto"
	header "whisper/internal/service/common"
	"whisper/pkg/config"
	"whisper/pkg/context"
	"whisper/pkg/errors"
	"whisper/pkg/http"
	"whisper/pkg/log"
)

// QueryEquipmentsForLOL 通过 https://101.qq.com/#/equipment 查询端游的所有装备列表
func QueryEquipmentsForLOL(ctx *context.Context) (*dto.LOLEquipment, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.Lol.Equipment, time.Now().Unix()/600)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	equip := dto.LOLEquipment{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &equip)
	return &equip, err
}

// QueryHeroesForLOL 通过 https://101.qq.com/#/hero 查询端游的所有英雄
func QueryHeroesForLOL(ctx *context.Context) (*dto.LOLHeroes, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.Lol.Heroes, time.Now().Unix()/600)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	heroes := dto.LOLHeroes{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &heroes)
	return &heroes, err
}

// QueryHeroesForLOLM 通过 https://game.gtimg.cn/images/lgamem/act/lrlib/js/heroList/hero_list.js 查询端游的所有英雄
func QueryHeroesForLOLM(ctx *context.Context) (*dto.LOLMHeroes, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.LolM.Heroes, time.Now().Unix()/600)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	heroes := dto.LOLMHeroes{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &heroes)
	return &heroes, err
}

// QueryEquipmentsForLOLM 通过 https://game.gtimg.cn/images/lgamem/act/lrlib/js/equip/equip.js 查询手游的所有装备列表
func QueryEquipmentsForLOLM(ctx *context.Context) (*dto.LOLMEquipment, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.LolM.Equipment, time.Now().Unix()/600)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	equip := dto.LOLMEquipment{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &equip)
	return &equip, err
}

// QueryRuneForLOL https://game.gtimg.cn/images/lol/act/img/js/runeList/rune_list.js
func QueryRuneForLOL(ctx *context.Context) (*dto.LOLRune, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.Lol.Rune, time.Now().Unix()/600)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	r := dto.LOLRune{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	return &r, err
}

// QueryRuneForLOLM https://game.gtimg.cn/images/lgamem/act/lrlib/js/rune/rune.js
func QueryRuneForLOLM(ctx *context.Context) (*dto.LOLMRune, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.LolM.Rune, time.Now().Unix()/600)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	r := dto.LOLMRune{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	return &r, err
}

// QuerySkillForLOL https://game.gtimg.cn/images/lol/act/img/js/summonerskillList/summonerskill_list.js
func QuerySkillForLOL(ctx *context.Context) (*dto.LOLSkill, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.Lol.Skill, time.Now().Unix()/600)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	r := dto.LOLSkill{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	return &r, err
}

// QuerySkillForLOLM https://game.gtimg.cn/images/lgamem/act/lrlib/js/skill/skill.js
func QuerySkillForLOLM(ctx *context.Context) (*dto.LOLMSkill, error) {
	url := fmt.Sprintf("%s?ts=%d", config.LOLConfig.LolM.Skill, time.Now().Unix()/600)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	r := dto.LOLMSkill{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	return &r, err
}

// GetLOLHeroAttribute https://xxx/%d.js
func GetLOLHeroAttribute(ctx *context.Context, heroID string) (*dto.HeroAttribute, error) {
	heroAttrUrl := fmt.Sprintf(config.LOLConfig.Lol.Hero, heroID)
	url := fmt.Sprintf("%s?ts=%d", heroAttrUrl, time.Now().Unix()/600)
	log.Logger.Debug(ctx, "url="+url)

	// 发送 GetForm 请求
	r := dto.HeroAttribute{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	return &r, err
}

// GetLOLMHeroAttribute https://xx/%d.js
func GetLOLMHeroAttribute(ctx *context.Context, heroID string) (*dto.HeroAttribute, error) {
	heroAttrUrl := fmt.Sprintf(config.LOLConfig.LolM.Hero, heroID)
	url := fmt.Sprintf("%s?ts=%d", heroAttrUrl, time.Now().Unix()/600)
	//log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	r := dto.HeroAttribute{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	return &r, err
}

// QueryRuneTypeForLOLM https://mlol.qt.qq.com/go/zone/views_layout?key=lr_rune_type
func QueryRuneTypeForLOLM(ctx *context.Context) (*dto.LOLMRuneType, error) {
	url := fmt.Sprintf("%s", config.LOLConfig.LolM.RuneType)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	runeType := dto.LOLMRuneType{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &runeType)
	return &runeType, err
}

func QuerySuitEquipForLOL(ctx *context.Context, heroId string) (*dto.JDataDataResult, error) {
	dtstatdate := time.Now().AddDate(0, 0, -1).Format("20060102")
	url := fmt.Sprintf(config.LOLConfig.Lol.SuitEquip, dtstatdate, heroId)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	suitEquip := dto.HeroSuitEquip{}

	body, err := http.GetForm(ctx, url, header.Referer...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &suitEquip)
	if err != nil {
		return nil, err
	}

	result := dto.JDataDataResult{}
	err = json.Unmarshal([]byte(suitEquip.JData.Data.Result), &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

// ChampionFightData 英雄对战数据详情 LOL
func ChampionFightData(ctx *context.Context, heroID string) (*dto.ChampionFightData, error) {
	//jsonpResponse := `var CHAMPION_DETAIL_17={"gameVer":"13.16","date":"2023-08-30 16:15:26"};/*  |xGv00|b214aa8b2b62d14489dce9170b96cdee */`
	champDetailUrl := fmt.Sprintf(config.LOLConfig.Lol.ChampDetail, heroID)
	url := fmt.Sprintf("%s?ts=%d", champDetailUrl, time.Now().Unix()/600)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	championFightData := dto.ChampionFightData{}

	body, err := http.GetForm(ctx, url, header.Referer...)
	if err != nil {
		return nil, err
	}

	// 使用正则表达式提取 JSON 数据部分
	re := regexp.MustCompile(`{.*}`)
	match := re.FindString(string(body))

	err = json.Unmarshal([]byte(match), &championFightData)
	if err != nil {
		return nil, err
	}

	return &championFightData, nil
}

// HeroRankList 手游各位置英雄胜率
func HeroRankList(ctx *context.Context) (*dto.HeroRankList, error) {
	url := config.LOLConfig.LolM.HeroWinRate
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	championFightData := dto.HeroRankList{}

	body, err := http.GetForm(ctx, url, header.Referer...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &championFightData)
	if err != nil {
		return nil, err
	}

	return &championFightData, nil
}

// HeroSuit 手游英雄推荐出装
func HeroSuit(ctx *context.Context, heroID string) (*dto.HeroTech, map[string]*dto.EquipTech, error) {
	heroTechUrl := fmt.Sprintf(config.LOLConfig.LolM.HeroSuit, heroID)
	log.Logger.Info(ctx, "heroTechUrl="+heroTechUrl)

	// 发送 GetForm 请求
	heroTech := dto.HeroTech{}
	// -----------------------------
	body, err := http.GetForm(ctx, heroTechUrl, header.Referer...)
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(body, &heroTech)
	if err != nil {
		return nil, nil, err
	}

	// -----------------------------
	wg := sync.WaitGroup{}

	syncMap := sync.Map{}
	for _, eqs := range heroTech.Data.AnchorRecommend.List {
		wg.Add(1)
		go func(eqs dto.AnchorRecommendList) {
			defer wg.Done()

			equipTechUrl := fmt.Sprintf(config.LOLConfig.LolM.HeroEquip, eqs.Head.Id)
			log.Logger.Info(ctx, "equipTechUrl="+equipTechUrl)
			body, err = http.GetForm(ctx, equipTechUrl, header.Referer...)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}

			et := dto.EquipTech{}
			err = json.Unmarshal(body, &et)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}

			syncMap.Store(eqs.Head.Id, &et)
		}(eqs)
	}
	wg.Wait()

	met := make(map[string]*dto.EquipTech)
	syncMap.Range(func(key, value interface{}) bool {
		met[key.(string)] = value.(*dto.EquipTech)
		return true
	})

	return &heroTech, met, nil
}

// LOLMVersionList 手游版本列表
func LOLMVersionList(ctx *context.Context) (*dto.LOLMVersionList, map[string]*dto.LOLMVersionDetail, error) {
	versionListUrl := config.LOLConfig.LolM.VersionList
	log.Logger.Info(ctx, "versionListUrl="+versionListUrl)

	// 发送 GetForm 请求
	versionList := dto.LOLMVersionList{}
	// -----------------------------
	body, err := http.GetForm(ctx, versionListUrl, header.CommonHeaders()...)
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(body, &versionList)
	if err != nil {
		return nil, nil, err
	}

	if versionList.Result != 0 {
		return nil, nil, errors.New(versionList.Msg, versionList.ErrMsg)
	}

	vd, err := VersionDetail(ctx, 1, versionList.Data[0].Vkey)
	if err != nil {
		return nil, nil, err
	}

	return &versionList, vd, nil
}

// VersionDetail 版本更新详情
func VersionDetail(ctx *context.Context, platform int, version string) (map[string]*dto.LOLMVersionDetail, error) {
	// -----------------------------
	wg := sync.WaitGroup{}

	keys := []string{
		fmt.Sprintf("lgame_%s_hero", version),
		fmt.Sprintf("lgame_%s_prop", version),
		fmt.Sprintf("lgame_%s_system", version),
		fmt.Sprintf("lgame_%s_rune", version),
	}
	syncMap := sync.Map{}
	for _, k := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()

			// https://mlol.qt.qq.com/go/database/versiondetail?key=%s
			versionDetailUrl := config.LOLConfig.LolM.VersionDetail
			detailUrl := fmt.Sprintf(versionDetailUrl, k)
			log.Logger.Info(ctx, "detailUrl="+detailUrl)
			body, err := http.GetForm(ctx, detailUrl, header.CommonHeaders()...)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}

			detail := dto.LOLMVersionDetail{}
			err = json.Unmarshal(body, &detail)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}

			syncMap.Store(k, &detail)
		}(k)
	}
	wg.Wait()

	vd := make(map[string]*dto.LOLMVersionDetail)
	syncMap.Range(func(key, value interface{}) bool {
		vd[key.(string)] = value.(*dto.LOLMVersionDetail)
		return true
	})

	return vd, nil

}
