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

	// 发送 PostJson 请求
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

	// 发送 PostJson 请求
	heroes := dto.LOLHeroes{}

	body, err := http.GetForm(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &heroes)
	return &heroes, err
}
