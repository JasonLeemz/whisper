package equipment

import (
	"fmt"
	"github.com/spf13/cast"
	"sort"
	"strconv"
	"strings"
	"whisper/internal/dto"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/internal/service/bloomfilter"
	"whisper/pkg/config"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/redis"
	"whisper/pkg/utils"
)

type Inner struct {
	ctx      *context.Context
	platform int
}

func NewInnerIns(ctx *context.Context) *Inner {
	return &Inner{ctx: ctx}
}

func (e *Inner) WithPlatform(platform int) *Inner {
	e.platform = platform
	return e
}

func (e *Inner) SetBit() {
	data := e.GetAll(common.PlatformForLOL).(map[string]*model.LOLEquipment)
	datam := e.GetAll(common.PlatformForLOLM).(map[string]*model.LOLMEquipment)

	bf := bloomfilter.NewEquipBloomFilter(redis.RDB, 3)

	for _, equip := range data {
		d := equip.ItemId
		bf.Add(d, redis.KeyBitSetEquip)
	}

	for _, equip := range datam {
		d := equip.EquipId
		bf.Add(d, redis.KeyBitSetEquip)
	}

}

func (e *Inner) Contain(data string) bool {
	bf := bloomfilter.NewEquipBloomFilter(redis.RDB, 3)
	return bf.Contains(data, redis.KeyBitSetEquip)
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

// GetAll return map[string]*model.LOLEquipment|map[string]*model.LOLMEquipment
func (e *Inner) GetAll(platform int) interface{} {
	// 获取全部装备
	if platform == common.PlatformForLOL {
		d := dao.NewLOLEquipmentDAO()
		eVersion, _ := d.GetLOLEquipmentMaxVersion()
		data, _ := d.GetLOLEquipment(eVersion.Version)
		mequip := make(map[string]*model.LOLEquipment)
		for _, equip := range data {
			key := fmt.Sprintf(redis.KeyCacheEquip, equip.Maps, strconv.Itoa(common.PlatformForLOL), equip.ItemId)
			mequip[key] = equip
		}
		return mequip
	} else {
		d := dao.NewLOLMEquipmentDAO()
		v, _ := d.GetLOLMEquipmentMaxVersion()
		data, _ := d.GetLOLMEquipment(v.Version)
		mequip := make(map[string]*model.LOLMEquipment)
		for _, equip := range data {
			key := fmt.Sprintf(redis.KeyCacheEquip, "召唤师峡谷", strconv.Itoa(common.PlatformForLOLM), equip.EquipId) // todo
			mequip[key] = equip
		}
		return mequip
	}
}

// GetOne 获取符文
func (e *Inner) GetOne(platform int, id interface{}) interface{} {
	// 获取全部装备
	if platform == common.PlatformForLOL {
		d := dao.NewLOLEquipmentDAO()
		data, _ := d.Find([]string{"*"}, map[string]interface{}{
			"itemId": id,
		})
		if len(data) == 0 {
			return nil
		} else {
			return data[0]
		}
	} else {
		d := dao.NewLOLMEquipmentDAO()
		data, _ := d.Find([]string{"*"}, map[string]interface{}{
			"equipId": id,
		})
		if len(data) == 0 {
			return nil
		} else {
			return data[0]
		}
	}

}
