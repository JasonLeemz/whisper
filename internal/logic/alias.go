package logic

import (
	"whisper/internal/logic/common"
	"whisper/internal/model"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/context"
	"whisper/pkg/pinyin"
)

func AliasHeroes(ctx *context.Context) (any, error) {

	// 查询所有Heroes
	hdao := dao.NewLOLHeroesDAO() // LOL 英雄覆盖LOLM，这里只处理LOL
	v, err := hdao.GetLOLHeroesMaxVersion()
	if err != nil {
		return nil, err
	}
	heroes, err := hdao.GetLOLHeroes(v.Version)
	if err != nil {
		return nil, err
	}

	// insert or update hero_alias
	hadao := dao.NewHeroAliasDAO()
	for _, hero := range heroes {
		name := hero.Name + " " + hero.Title

		py, first := pinyin.Trans(hero.Name)
		keywordsPy := py + "," + first
		py, first = pinyin.Trans(hero.Title)
		keywordsPy += "," + py + "," + first
		alias := &model.HeroAlias{
			KeywordsPy: keywordsPy,
		}

		cond := map[string]interface{}{
			"name": name,
		}
		exists, err := hadao.Exists(cond)
		if err != nil {
			return nil, err
		}
		if !exists {
			// 不存在就插入
			alias.Name = name
			_, err = hadao.Add([]*model.HeroAlias{
				alias,
			})
		} else {
			// 存在就更新
			_, err = hadao.Update(alias, cond)
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, err
}

func AliasEquip(ctx *context.Context) (any, error) {

	// 查询所有Equipments
	equipDao := dao.NewLOLEquipmentDAO()
	v, err := equipDao.GetLOLEquipmentMaxVersion()
	if err != nil {
		return nil, err
	}
	equips, err := equipDao.GetLOLEquipment(v.Version)
	if err != nil {
		return nil, err
	}

	m_equipDao := dao.NewLOLMEquipmentDAO()
	m_v, err := m_equipDao.GetLOLMEquipmentMaxVersion()
	if err != nil {
		return nil, err
	}
	m_equips, err := m_equipDao.GetLOLMEquipment(m_v.Version)
	if err != nil {
		return nil, err
	}

	// insert or update hero_alias
	eadao := dao.NewEquipAliasDAO()
	for _, equip := range equips {
		name := equip.Name

		py, first := pinyin.Trans(equip.Name)
		keywordsPy := py + "," + first
		alias := &model.EquipAlias{
			KeywordsPy: keywordsPy,
			Platform:   common.PlatformForLOL,
		}

		cond := map[string]interface{}{
			"name": name,
		}
		exists, err := eadao.Exists(cond)
		if err != nil {
			return nil, err
		}
		if !exists {
			// 不存在就插入
			alias.Name = name
			_, err = eadao.Add([]*model.EquipAlias{
				alias,
			})
		} else {
			// 存在就更新
			_, err = eadao.Update(alias, cond)
		}
		if err != nil {
			return nil, err
		}
	}

	for _, equip := range m_equips {
		name := equip.Name

		py, first := pinyin.Trans(equip.Name)
		keywordsPy := py + "," + first
		alias := &model.EquipAlias{
			KeywordsPy: keywordsPy,
			Platform:   common.PlatformForLOLM,
		}

		cond := map[string]interface{}{
			"name": name,
		}
		exists, err := eadao.Exists(cond)
		if err != nil {
			return nil, err
		}
		if !exists {
			// 不存在就插入
			alias.Name = name
			_, err = eadao.Add([]*model.EquipAlias{
				alias,
			})
		} else {
			// 存在就更新
			_, err = eadao.Update(alias, cond)
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, err
}
