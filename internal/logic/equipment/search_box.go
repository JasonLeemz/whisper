package equipment

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	lol "whisper/internal/service/lol"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/pinyin"
)

func QueryEquipments(ctx *context.Context, platform int) (any, error) {
	equipments, err := lol.CreateLOLProduct(platform)().QueryEquipments(ctx)
	if err != nil {
		log.Logger.Error(ctx, err)
		return nil, err
	}

	if platform == common.PlatformForLOL {
		reloadEquipmentForLOL(ctx, equipments.(*dto.LOLEquipment))
	} else {
		reloadEquipmentForLOLM(ctx, equipments.(*dto.LOLMEquipment))
	}

	return equipments, nil
}

func reloadEquipmentForLOL(ctx *context.Context, equip *dto.LOLEquipment) {

	equipDao := dao.NewLOLEquipmentDAO()

	// 判断库中是否存在最新版本，如果存在就不更新
	result, err := equipDao.GetLOLEquipmentMaxVersion()
	if err != nil {
		log.Logger.Error(ctx, err)
		return
	}

	if result != nil {
		log.Logger.Info(ctx,
			fmt.Sprintf("DB Version[%s] fileTime[%s],Data Version:[%s] fileTime[%s]", result.Version, result.FileTime, equip.Version, equip.FileTime),
		)
		x, err := common.CompareTime(result.FileTime, equip.FileTime)
		if err != nil {
			log.Logger.Error(ctx, err)
			return
		}
		if x != "<" {
			// 如果原始数据版本和当前获取数据的版本相等，就不更新数据库
			log.Logger.Info(ctx, "原始数据版本和当前获取数据的版本相等,不更新")
			return
		}

		// 有可能日期更新了，但是版本号没变
		if result.Version == equip.Version {
			// 将db中该版本的数据status置为1
			cond := map[string]interface{}{
				"version": result.Version,
				"status":  0,
			}
			data := model.LOLEquipment{
				Status: 1,
			}
			up, err := equipDao.Update(&data, cond)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}
			log.Logger.Info(ctx, "当前版本数据不是最新,已经软删除,生效行数:", up)
		}
	}

	log.Logger.Info(ctx, "running record LOL equipment data...")
	startT := time.Now()
	// 入库更新
	equips := make([]*model.LOLEquipment, 0, len(equip.Items)+int(math.Floor(float64(len(equip.Items)/3))))

	for _, item := range equip.Items {
		namePY, nameF := pinyin.Trans(item.Name)
		searchKey := namePY + "," + nameF

		tmp := model.LOLEquipment{
			ItemId:      item.ItemId,
			Name:        item.Name,
			IconPath:    item.IconPath,
			Price:       item.Price,
			Description: item.Description,
			Plaintext:   item.Plaintext,
			Sell:        item.Sell,
			Total:       item.Total,
			Tag:         item.Tag,
			Keywords:    item.Keywords + "," + searchKey,
			Version:     equip.Version,
			FileTime:    equip.FileTime,

			From: convRoadmapData4LOL(item.From),
			Into: convRoadmapData4LOL(item.Into),
		}

		for _, m := range item.Maps {
			eqModel := tmp
			eqModel.Maps = m

			equips = append(equips, &eqModel)
		}
	}

	// 记录装备信息
	_, err = equipDao.Add(equips)
	if err != nil {
		log.Logger.Error(ctx, err)
	}

	log.Logger.Info(ctx, fmt.Sprintf("finish record LOL equipment data. Since:%fs", time.Since(startT).Seconds()))
}

func reloadEquipmentForLOLM(ctx *context.Context, equip *dto.LOLMEquipment) {
	equipDao := dao.NewLOLMEquipmentDAO()

	// 判断库中是否存在最新版本，如果存在就不更新
	result, err := equipDao.GetLOLMEquipmentMaxVersion()
	if err != nil {
		log.Logger.Error(ctx, err)
		return
	}

	if result != nil {
		log.Logger.Info(ctx,
			fmt.Sprintf("DB Version[%s] fileTime[%s],Data Version:[%s] fileTime[%s]", result.Version, result.FileTime, equip.Version, equip.FileTime),
		)
		x, err := common.CompareTime(result.FileTime, equip.FileTime)
		if err != nil {
			log.Logger.Error(ctx, err)
			return
		}
		if x != "<" {
			// 如果原始数据版本和当前获取数据的版本相等，就不更新数据库
			log.Logger.Info(ctx, "原始数据版本和当前获取数据的版本相等,不更新")
			return
		}

		// 有可能日期更新了，但是版本号没变
		if result.Version == equip.Version {
			// 将db中该版本的数据status置为1
			cond := map[string]interface{}{
				"version": result.Version,
				"status":  0,
			}
			data := model.LOLMEquipment{
				Status: 1,
			}
			up, err := equipDao.Update(&data, cond)
			if err != nil {
				log.Logger.Error(ctx, err)
				return
			}
			log.Logger.Info(ctx, "当前版本数据不是最新,已经软删除,生效行数:", up)
		}
	}

	log.Logger.Info(ctx, "running record LOLM equipment data...")
	startT := time.Now()
	// 入库更新
	equips := make([]*model.LOLMEquipment, 0, len(equip.EquipList))

	into := make(map[string][]string)

	for _, item := range equip.EquipList {
		fs, fa := convRoadmapData4LOLM(item.From)
		for _, fid := range fa {
			into[fid] = append(into[fid], item.EquipId)
		}

		namePY, nameF := pinyin.Trans(item.Name)
		searchKey := namePY + "," + nameF + "," + item.Type + "," + item.Level
		tmp := model.LOLMEquipment{
			EquipId:         item.EquipId,
			Name:            item.Name,
			IconPath:        item.IconPath,
			Type:            item.Type,
			Level:           item.Level,
			Price:           item.Price,
			Hp:              item.Hp,
			HpRegen:         item.HpRegen,
			HpRegenRate:     item.HpRegenRate,
			Armor:           item.Armor,
			ArmorPene:       item.ArmorPene,
			ArmorPeneRate:   item.ArmorPeneRate,
			CritRate:        item.CritRate,
			CritDamage:      item.CritDamage,
			AttackSpeed:     item.AttackSpeed,
			HealthPerAttack: item.HealthPerAttack,
			MagicAttack:     item.MagicAttack,
			Mp:              item.Mp,
			MpRegen:         item.MpRegen,
			MagicBlock:      item.MagicBlock,
			MagicPene:       item.MagicPene,
			MagicPeneRate:   item.MagicPeneRate,
			HealthPerMagic:  item.HealthPerMagic,
			Cd:              item.Cd,
			DuctRate:        item.DuctRate,
			MoveSpeed:       item.MoveSpeed,
			MoveRate:        item.MoveRate,
			ComposeLevel:    item.ComposeLevel,
			Ad:              item.Ad,
			Into:            item.Into,
			Tags:            item.Tags,
			UnName:          item.UnName,
			SearchKey:       searchKey,
			Version:         equip.Version,
			FileTime:        equip.FileTime,

			Description: strings.Join(item.Description, ""),
			//Description: strings.Join(item.Description, "<br>"),

			From: fs,
		}

		equips = append(equips, &tmp)
	}
	// 记录装备信息
	_, err = equipDao.Add(equips)
	if err != nil {
		log.Logger.Error(ctx, err)
	}

	updatesInto, err := equipDao.UpdatesInto(equip.FileTime, equip.Version, into)
	log.Logger.Info(ctx, fmt.Sprintf("Update LOLM equipment into. Rows:%d,err:%v", updatesInto, err))

	log.Logger.Info(ctx, fmt.Sprintf("finish record LOLM equipment data. Since:%fs", time.Since(startT).Seconds()))
}

func convRoadmapData4LOL(data any) string {
	var t []string
	b, _ := json.Marshal(data)
	err := json.Unmarshal(b, &t)
	if err != nil {
		return ""
	}

	return strings.Join(t, ",")
}

func convRoadmapData4LOLM(data any) (string, []string) {
	var t []string
	b, _ := json.Marshal(data)
	err := json.Unmarshal(b, &t)
	if err != nil {
		return "", nil
	}

	return strings.Join(t, ","), t
}
