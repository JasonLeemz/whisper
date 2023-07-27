package logic

import (
	"github.com/spf13/cast"
	"math"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/internal/service"
	"whisper/pkg/context"
	"whisper/pkg/errors"
	"whisper/pkg/log"
)

func QueryEquipments(ctx *context.Context, platform int) (any, *errors.Error) {

	equip, err := service.QueryEquipmentsForLOL(ctx)
	if err != nil {
		log.Logger.Warn(ctx, err)
	}

	// 入库更新
	equips := make([]*model.Equipment, 0, len(equip.Items)+int(math.Floor(float64(len(equip.Items)/3))))
	heroesSuit := make([]*model.HeroesSuit, 0)
	equipType := make([]*model.EquipType, 0)

	for _, item := range equip.Items {
		tmp := model.Equipment{
			Platform:    platform,
			ItemId:      item.ItemId,
			Name:        item.Name,
			IconPath:    item.IconPath,
			Price:       item.Price,
			Description: item.Description,
			Plaintext:   item.Plaintext,
			Sell:        item.Sell,
			Total:       item.Total,
			Tag:         item.Tag,
			Keywords:    item.Keywords,
			Version:     equip.Version,
			FileTime:    equip.FileTime,

			//From:       item.From,
			//Into:       item.Into,
		}

		for _, m := range item.Maps {
			eqModel := tmp
			eqModel.Maps = m

			equips = append(equips, &eqModel)
		}

		// 记录英雄适配装备表
		switch item.SuitHeroId.(type) {
		case []interface{}:
			for _, heroID := range item.SuitHeroId.([]interface{}) {
				id := cast.ToString(heroID)
				hsTmp := &model.HeroesSuit{
					HeroId: id,
					ItemId: item.ItemId,
				}
				heroesSuit = append(heroesSuit, hsTmp)
			}
		}

		// 记录装备所属类型表
		for _, t := range item.Types {
			etTmp := model.EquipType{
				Types:  t,
				ItemId: item.ItemId,
			}
			equipType = append(equipType, &etTmp)
		}

	}

	// 记录英雄适配装备表
	hsDao := dao.NewHeroesSuitDAO()
	_, err = hsDao.Add(heroesSuit)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}

	// 记录装备所属类型表
	etDao := dao.NewEquipTypeDAO()
	_, err = etDao.Add(equipType)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}

	// 记录装备信息
	equipDao := dao.NewEquipmentDAO()
	_, err = equipDao.Add(equips)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}

	return equip, errors.New(err)

}
