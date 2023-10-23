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

type InnerEquip struct {
	ctx      *context.Context
	platform int
}

func NewInnerEquip(ctx *context.Context, platform int) *InnerEquip {
	return &InnerEquip{ctx: ctx, platform: platform}
}

func (e *InnerEquip) ExtractKeyWords() map[string]model.EquipIntro {
	result := extractEquipKeywords(e.ctx, e.platform)
	recordMongo(e.ctx, result, e.platform)
	return result
}

type InnerEquipCommand struct {
	*InnerEquip
}

func (e *InnerEquip) NewExtractKeyWordsCmd() *InnerEquipCommand {
	return &InnerEquipCommand{
		e,
	}
}

func (cmd InnerEquipCommand) Exec() error {
	cmd.ExtractKeyWords()
	return nil
}

func ExtractKeyWords(ctx *context.Context, platform int) map[string]model.EquipIntro {
	result := extractEquipKeywords(ctx, platform)
	recordMongo(ctx, result, platform)
	return result
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
				Version:   equip.Version,
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
				Version:   equip.Version,
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
