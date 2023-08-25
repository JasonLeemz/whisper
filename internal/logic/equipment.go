package logic

import (
	errors2 "errors"
	"fmt"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"sort"
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
	"whisper/pkg/log"
	"whisper/pkg/pinyin"
	"whisper/pkg/utils"
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

func ExtractKeyWords(ctx *context.Context, platform int) map[string]model.EquipIntro {
	result := extractEquipKeywords(ctx, platform)
	recordMongo(ctx, result, platform)
	return result
}

func extractEquipKeywords(ctx *context.Context, platform int) map[string]model.EquipIntro {
	_, dict := GetEquipTypes(ctx)
	re := utils.CompileKeywordsRegex(dict)

	result := make(map[string]model.EquipIntro)
	if platform == common.PlatformForLOL {
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

		for _, equip := range equips {
			words := utils.ExtractKeywords(equip.Description, re)
			result[equip.ItemId] = model.EquipIntro{
				ID:        equip.ItemId,
				Name:      equip.Name,
				Icon:      equip.IconPath,
				Desc:      utils.RemoveRepeatedBRTag(equip.Description),
				Plaintext: equip.Plaintext,
				Price:     cast.ToFloat64(equip.Total),
				Maps:      equip.Maps,
				Platform:  common.PlatformForLOL,
				Keywords:  words,
			}
		}
	} else {
		ed := dao.NewLOLMEquipmentDAO()
		v, err := ed.GetLOLMEquipmentMaxVersion()
		if err != nil {
			log.Logger.Error(ctx, err)
			return nil
		}
		equips, err := ed.GetLOLMEquipment(v.Version)
		if err != nil {
			log.Logger.Error(ctx, err)
			return nil
		}

		for _, equip := range equips {
			words := utils.ExtractKeywords(equip.Description, re)
			result[equip.EquipId] = model.EquipIntro{
				ID:        equip.EquipId,
				Name:      equip.Name,
				Icon:      equip.IconPath,
				Desc:      utils.RemoveRepeatedBRTag(equip.Description),
				Plaintext: "-",
				Price:     cast.ToFloat64(equip.Price),
				Maps:      "召唤师峡谷",
				Platform:  common.PlatformForLOLM,
				Keywords:  words,
			}
		}
	}
	return result
}

func GetEquipTypes(ctx *context.Context) ([]*dto.EquipType, []string) {
	equipTypes := make([]*dto.EquipType, 0)
	dict := make([]string, 0)

	// 为了保证输出有序
	// http://nacos.ybdx.xyz/nacos/#/configeditor?serverId=center&dataId=lol_equip_dict&group=dev&namespace=f320980d-d47e-4b63-896e-29879ea5a72e&edasAppName=&edasAppId=&searchDataId=&searchGroup=&pageSize=10&pageNo=1
	for _, cate := range config.EquipDict.Extract.EquipShow {
		if sub, ok := config.EquipDict.Extract.Equip[cate]; ok {
			equipType := &dto.EquipType{
				Cate: cate,
			}

			var sortKeys []string
			for key := range sub {
				sortKeys = append(sortKeys, key)
			}
			sort.Strings(sortKeys)

			subCateStr := make([]map[string]string, 0)
			for _, sk := range sortKeys {
				split := strings.Split(sk, ".")
				if len(split) < 2 {
					continue
				}
				equipType.SubCate = append(equipType.SubCate, dto.SubCate{
					Name:          split[1],
					KeywordsSlice: sub[sk],
					KeywordsStr:   strings.Join(sub[sk], ","),
				})
				subCateStr = append(subCateStr, map[string]string{
					split[1]: strings.Join(sub[sk], ","),
				})

				dict = append(dict, sub[sk]...)

			}

			equipTypes = append(equipTypes, equipType)
		}
	}

	return equipTypes, dict
}

func recordMongo(ctx *context.Context, data map[string]model.EquipIntro, platform int) {

	md := dao.NewMongoEquipmentDAO()
	equips := make([]*model.EquipIntro, 0, len(data))
	for _, intro := range data {
		introCopy := intro
		equips = append(equips, &introCopy)
	}

	cond := map[string]interface{}{
		"platform": platform,
	}
	err := md.Delete(ctx, cond)
	if err != nil {
		log.Logger.Error(ctx, err)
		return
	}
	err = md.Add(ctx, equips)
	if err != nil {
		log.Logger.Error(ctx, err)
	}
}

func FilterKeyWords(ctx *context.Context, keywords []string, platform int) ([]*model.EquipIntro, error) {
	log.Logger.Info(ctx, keywords)
	// FromMongo
	md := dao.NewMongoEquipmentDAO()

	kw := make([]bson.M, 0, len(keywords))
	for _, words := range keywords {
		in := strings.Split(words, ",")
		kw = append(kw, bson.M{
			"keywords": bson.M{
				"$in": in,
			},
		})
	}
	// 构建查询条件
	filter := bson.M{
		"platform": platform,
		"maps":     "召唤师峡谷",
		"$and":     kw,
	}

	result, err := md.Find(ctx, filter)
	return result, err
}
