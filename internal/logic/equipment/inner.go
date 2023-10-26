package equipment

import (
	"github.com/spf13/cast"
	"sort"
	"strings"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/config"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/utils"
)

type Inner struct {
	ctx      *context.Context
	platform int
}

func NewInnerEquip(ctx *context.Context, platform int) *Inner {
	return &Inner{ctx: ctx, platform: platform}
}

func (e *Inner) ExtractKeyWords() map[string]model.EquipIntro {
	result := e.extractEquipKeywords()
	e.recordMongo(result)
	return result
}

func (e *Inner) recordMongo(data map[string]model.EquipIntro) {

	md := dao.NewMongoEquipmentDAO()
	equips := make([]*model.EquipIntro, 0, len(data))
	for _, intro := range data {
		introCopy := intro
		equips = append(equips, &introCopy)
	}

	cond := map[string]interface{}{
		"platform": e.platform,
	}
	err := md.Delete(e.ctx, cond)
	if err != nil {
		log.Logger.Error(e.ctx, err)
		return
	}
	err = md.Add(e.ctx, equips)
	if err != nil {
		log.Logger.Error(e.ctx, err)
	}
}

func (e *Inner) extractEquipKeywords() map[string]model.EquipIntro {
	_, dict := e.GetEquipTypes()
	re := utils.CompileKeywordsRegex(dict)

	result := make(map[string]model.EquipIntro)
	if e.platform == common.PlatformForLOL {
		ed := dao.NewLOLEquipmentDAO()
		v, err := ed.GetLOLEquipmentMaxVersion()
		if err != nil {
			log.Logger.Error(e.ctx, err)
			return nil
		}
		equips, err := ed.GetLOLEquipment(v.Version)
		if err != nil {
			log.Logger.Error(e.ctx, err)
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
				Version:   equip.Version,
				Keywords:  words,
			}
		}
	} else {
		ed := dao.NewLOLMEquipmentDAO()
		v, err := ed.GetLOLMEquipmentMaxVersion()
		if err != nil {
			log.Logger.Error(e.ctx, err)
			return nil
		}
		equips, err := ed.GetLOLMEquipment(v.Version)
		if err != nil {
			log.Logger.Error(e.ctx, err)
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
				Version:   equip.Version,
				Keywords:  words,
			}
		}
	}
	return result
}

func (e *Inner) GetEquipTypes() ([]*dto.EquipType, []string) {
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

func (e *Inner) GetAll(platform int) (interface{}, error) {
	// 获取全部装备
	if platform == common.PlatformForLOL {
		d := dao.NewLOLEquipmentDAO()
		eVersion, err := d.GetLOLEquipmentMaxVersion()
		if err != nil {
			return nil, err
		}
		data, err := d.GetLOLEquipment(eVersion.Version)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		d := dao.NewLOLMEquipmentDAO()
		v, err := d.GetLOLMEquipmentMaxVersion()
		if err != nil {
			return nil, err
		}
		data, err := d.GetLOLMEquipment(v.Version)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

}
