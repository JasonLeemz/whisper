package logic

import (
	errors2 "errors"
	"fmt"
	"github.com/spf13/cast"
	"math"
	"strings"
	"time"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/internal/service"
	"whisper/pkg/config"
	"whisper/pkg/context"
	"whisper/pkg/errors"
	"whisper/pkg/jieba"
	"whisper/pkg/log"
	"whisper/pkg/pinyin"
)

func QueryEquipments(ctx *context.Context, platform int) (any, *errors.Error) {

	if platform == common.PlatformForLOL {
		equip, err := service.QueryEquipmentsForLOL(ctx)
		if err != nil {
			log.Logger.Warn(ctx, err)
		}
		reloadEquipmentForLOL(ctx, equip)
		return equip, nil
	} else if platform == common.PlatformForLOLM {
		equip, err := service.QueryEquipmentsForLOLM(ctx)
		if err != nil {
			log.Logger.Warn(ctx, err)
		}
		reloadEquipmentForLOLM(ctx, equip)
		return equip, nil
	}

	return nil, errors.New(errors2.New("请指定游戏平台"), errors.ErrNoInvalidInput)
}

func reloadEquipmentForLOL(ctx *context.Context, equip *dto.LOLEquipment) {

	equipDao := dao.NewLOLEquipmentDAO()

	// 判断库中是否存在最新版本，如果存在就不更新
	result, err := equipDao.GetLOLEquipmentMaxVersion()
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
		return
	}

	if result != nil {
		log.Logger.Info(ctx,
			fmt.Sprintf("DB Version[%s] fileTime[%s],Data Version:[%s] fileTime[%s]", result.Version, result.FileTime, equip.Version, equip.FileTime),
		)
		x, err := common.CompareTime(result.FileTime, equip.FileTime)
		if err != nil {
			log.Logger.Error(ctx, errors.New(err))
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
	heroesSuit := make([]*model.HeroesSuit, 0)
	equipType := make([]*model.EquipType, 0)

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

	// 记录装备信息
	_, err = equipDao.Add(equips)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}

	log.Logger.Info(ctx, fmt.Sprintf("finish record LOL equipment data. Since:%fs", time.Since(startT).Seconds()))
}

func reloadEquipmentForLOLM(ctx *context.Context, equip *dto.LOLMEquipment) {
	equipDao := dao.NewLOLMEquipmentDAO()

	// 判断库中是否存在最新版本，如果存在就不更新
	result, err := equipDao.GetLOLMEquipmentMaxVersion()
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
		return
	}

	if result != nil {
		log.Logger.Info(ctx,
			fmt.Sprintf("DB Version[%s] fileTime[%s],Data Version:[%s] fileTime[%s]", result.Version, result.FileTime, equip.Version, equip.FileTime),
		)
		x, err := common.CompareTime(result.FileTime, equip.FileTime)
		if err != nil {
			log.Logger.Error(ctx, errors.New(err))
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

	for _, item := range equip.EquipList {

		namePY, nameF := pinyin.Trans(item.Name)
		searchKey := namePY + "," + nameF + "," + item.Type + "," + item.Level
		tmp := model.LOLMEquipment{
			EquipId:  item.EquipId,
			Name:     item.Name,
			IconPath: item.IconPath,
			//From:            item.From,
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
		}

		equips = append(equips, &tmp)
	}
	// 记录装备信息
	_, err = equipDao.Add(equips)
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
	}

	log.Logger.Info(ctx, fmt.Sprintf("finish record LOLM equipment data. Since:%fs", time.Since(startT).Seconds()))
}

func GetCurrentLOLVersion(ctx *context.Context) string {
	equipDao := dao.NewLOLEquipmentDAO()
	result, err := equipDao.GetLOLEquipmentMaxVersion()
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
		return ""
	}
	if result != nil {
		return result.Version
	}
	return ""
}

func GetCurrentLOLMVersion(ctx *context.Context) string {
	equipDao := dao.NewLOLMEquipmentDAO()
	result, err := equipDao.GetLOLMEquipmentMaxVersion()
	if err != nil {
		log.Logger.Error(ctx, errors.New(err))
		return ""
	}
	if result != nil {
		return result.Version
	}
	return ""
}

type equipIntro struct {
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	Price    string   `json:"price"`
	Keywords []string `json:"keywords"`
}

func ExtractKeyWords(ctx *context.Context, platform int) map[string]equipIntro {

	ed := dao.NewLOLEquipmentDAO()
	v, err := ed.GetLOLEquipmentMaxVersion()
	if err != nil {
		log.Logger.Error(ctx, err)
		return nil
	}
	equips, err := ed.GetLOLEquipment(v.Version)
	if err != nil {
		log.Logger.Error(ctx, err)
		return nil
	}

	result := make(map[string]equipIntro)

	for i, equip := range equips {
		if i > 30 {
			break
		}
		words, err := jieba.Analyzer(ctx, equip.Description, config.EquipDict.Extract.EquipWords, config.EquipDict.Stopwords)
		if err != nil {
			log.Logger.Error(ctx, err)
			return nil
		}
		result[equip.ItemId] = equipIntro{
			Name:     equip.Name,
			Desc:     equip.Description,
			Price:    equip.Total,
			Keywords: words,
		}
	}

	recordMongo(result)
	return result
}
func recordMongo(ctx map[string]equipIntro) {
}

func GetEquipTypes(ctx *context.Context) map[string][]string {
	return config.EquipDict.Extract.EquipWords
}
