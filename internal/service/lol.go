package service

import (
	"encoding/json"
	"fmt"
	"time"
	"whisper/internal/dto"
	"whisper/pkg/config"
	"whisper/pkg/context"
	"whisper/pkg/http"
	"whisper/pkg/log"
)

// QueryEquipmentsForLOL 通过 https://101.qq.com/#/equipment 查询端游的所有装备列表
func QueryEquipmentsForLOL(ctx *context.Context) (*dto.LOLEquipment, error) {
	url := fmt.Sprintf("%s?ts=%d", config.GlobalConfig.Lol.Equipment, time.Now().Unix())
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
	url := fmt.Sprintf("%s?ts=%d", config.GlobalConfig.Lol.Heroes, time.Now().Unix())
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
	url := fmt.Sprintf("%s?ts=%d", config.GlobalConfig.LolM.Heroes, time.Now().Unix())
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
	url := fmt.Sprintf("%s?ts=%d", config.GlobalConfig.LolM.Equipment, time.Now().Unix())
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
	url := fmt.Sprintf("%s?ts=%d", config.GlobalConfig.Lol.Rune, time.Now().Unix())
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
	url := fmt.Sprintf("%s?ts=%d", config.GlobalConfig.LolM.Rune, time.Now().Unix())
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
	url := fmt.Sprintf("%s?ts=%d", config.GlobalConfig.Lol.Skill, time.Now().Unix())
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
	url := fmt.Sprintf("%s?ts=%d", config.GlobalConfig.LolM.Skill, time.Now().Unix())
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
