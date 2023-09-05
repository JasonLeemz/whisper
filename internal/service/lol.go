package service

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
	"whisper/internal/dto"
	"whisper/pkg/config"
	"whisper/pkg/context"
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

// GetLOLHeroAttribute https://game.gtimg.cn/images/lol/act/img/js/hero/%d.js
func GetLOLHeroAttribute(ctx *context.Context, heroID string) (*dto.HeroAttribute, error) {
	heroAttrUrl := fmt.Sprintf(config.LOLConfig.Lol.Hero, heroID)
	url := fmt.Sprintf("%s?ts=%d", heroAttrUrl, time.Now().Unix()/600)
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	r := dto.HeroAttribute{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	return &r, err
}

// GetLOLMHeroAttribute https://game.gtimg.cn/images/lgamem/act/lrlib/js/hero/%d.js
func GetLOLMHeroAttribute(ctx *context.Context, heroID string) (*dto.HeroAttribute, error) {
	heroAttrUrl := fmt.Sprintf(config.LOLConfig.LolM.Hero, heroID)
	url := fmt.Sprintf("%s?ts=%d", heroAttrUrl, time.Now().Unix()/600)
	log.Logger.Info(ctx, "url="+url)

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

	header := http.Header{
		Key:   "Referer",
		Value: "https://101.qq.com/",
	}
	body, err := http.GetForm(ctx, url, header)
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

	header := http.Header{
		Key:   "Referer",
		Value: "https://101.qq.com/",
	}
	body, err := http.GetForm(ctx, url, header)
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

// HeroRankList 手游个位置英雄胜率
func HeroRankList(ctx *context.Context) (*dto.HeroRankList, error) {
	url := config.LOLConfig.LolM.HeroWinRate
	log.Logger.Info(ctx, "url="+url)

	// 发送 GetForm 请求
	championFightData := dto.HeroRankList{}

	header := http.Header{
		Key:   "Referer",
		Value: "https://lolm.qq.com/",
	}
	body, err := http.GetForm(ctx, url, header)
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
func HeroSuit(ctx *context.Context, heroID string) (*dto.HeroTech, *dto.EquipTech, error) {
	heroTechUrl := fmt.Sprintf(config.LOLConfig.LolM.HeroSuit, heroID)
	equipTechUrl := fmt.Sprintf(config.LOLConfig.LolM.HeroEquip, heroID)
	log.Logger.Info(ctx, "heroTechUrl="+heroTechUrl, "equipTechUrl="+equipTechUrl)

	// 发送 GetForm 请求
	heroTech := dto.HeroTech{}
	equipTech := dto.EquipTech{}

	header := http.Header{
		Key:   "Referer",
		Value: "https://101.qq.com/",
	}
	// -----------------------------
	body, err := http.GetForm(ctx, heroTechUrl, header)
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(body, &heroTech)
	if err != nil {
		return nil, nil, err
	}

	// -----------------------------
	body, err = http.GetForm(ctx, equipTechUrl, header)
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(body, &equipTech)
	if err != nil {
		return nil, nil, err
	}

	return &heroTech, &equipTech, nil
}
