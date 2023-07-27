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
	//// Create a Resty Client
	//client := resty.New()
	//
	//_, err := client.R().
	//	SetHeader("Accept", "application/x-www-form-urlencoded").
	//	SetResult(&equip).
	//	ForceContentType("application/json").
	//	Get(url)

	return &equip, err
}
