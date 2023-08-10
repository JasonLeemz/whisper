package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type LOLEquipment interface {
	Add(data []*model.LOLEquipment) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLEquipment, error)
	GetLOLEquipmentMaxVersion() (*model.LOLEquipment, error)
	GetLOLEquipment(version string) ([]*model.LOLEquipment, error)
	GetLOLEquipmentWithExt(version string) ([]*model.LOLEquipment, error)
}

type LOLEquipmentDAO struct {
	db *gorm.DB
}

func (dao *LOLEquipmentDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLEquipment, error) {
	tx := dao.db.Model(&model.LOLEquipment{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.LOLEquipment
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}

func (dao *LOLEquipmentDAO) Add(equips []*model.LOLEquipment) (int64, error) {
	result := dao.db.Create(equips)
	return result.RowsAffected, result.Error
}

func (dao *LOLEquipmentDAO) GetLOLEquipmentMaxVersion() (*model.LOLEquipment, error) {
	tx := dao.db.Model(&model.LOLEquipment{})
	var result model.LOLEquipment
	tx = tx.Order("version desc").First(&result)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, tx.Error
}

func (dao *LOLEquipmentDAO) GetLOLEquipment(version string) ([]*model.LOLEquipment, error) {
	// 查当前版本所有数据
	cond := map[string]interface{}{
		"version": version,
	}
	equip, err := dao.Find(nil, cond)
	if err != nil {
		return nil, err
	}
	return equip, err
}

func (dao *LOLEquipmentDAO) GetLOLEquipmentWithExt(version string) ([]*model.LOLEquipment, error) {
	// 查当前版本所有数据
	result := make([]*model.LOLEquipment, 0)
	sql := `
SELECT
	equip.itemId,
	equip.name,
	CONCAT_WS(',', equip.keywords, alias.keywords, alias.keywords_py) AS keywords,
	equip.description,
	equip.plaintext,
	equip.iconPath,
	equip.price,
	equip.sell,
	equip.total,
	equip.suitHeroId,
	equip.maps,
	equip.from,
	equip.into,
	equip.types,
	equip.version,
	equip.fileTime
FROM
	lol_equipment equip
	LEFT JOIN equip_alias alias ON equip.name = alias.name
WHERE
	equip.version = '%s' and alias.platform = 0
`
	sql = fmt.Sprintf(sql, version)
	err := dao.db.Exec(sql).Find(&result).Error
	return result, err
}

func NewLOLEquipmentDAO() *LOLEquipmentDAO {
	return &LOLEquipmentDAO{
		db: mysql.DB,
	}
}

// --------------------------------------------------------

type LOLMEquipment interface {
	Add(equips []*model.LOLMEquipment) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLMEquipment, error)
	GetLOLMEquipmentMaxVersion() (*model.LOLMEquipment, error)
	GetLOLMEquipment(version string) ([]*model.LOLMEquipment, error)
	GetLOLMEquipmentWithExt(version string) ([]*model.LOLMEquipment, error)
}

type LOLMEquipmentDAO struct {
	db *gorm.DB
}

func (dao *LOLMEquipmentDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLMEquipment, error) {
	tx := dao.db.Model(&model.LOLMEquipment{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.LOLMEquipment
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}
func (dao *LOLMEquipmentDAO) Add(equips []*model.LOLMEquipment) (int64, error) {
	result := dao.db.Create(equips)
	return result.RowsAffected, result.Error
}
func (dao *LOLMEquipmentDAO) GetLOLMEquipmentMaxVersion() (*model.LOLMEquipment, error) {
	tx := dao.db.Model(&model.LOLMEquipment{})
	var result model.LOLMEquipment
	tx = tx.Order("version desc").First(&result)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, tx.Error
}
func (dao *LOLMEquipmentDAO) GetLOLMEquipment(version string) ([]*model.LOLMEquipment, error) {
	// 查当前版本所有数据
	cond := map[string]interface{}{
		"version": version,
	}
	equip, err := dao.Find(nil, cond)
	if err != nil {
		return nil, err
	}
	return equip, err
}

func (dao *LOLMEquipmentDAO) GetLOLMEquipmentWithExt(version string) ([]*model.LOLMEquipment, error) {
	// 查当前版本所有数据
	result := make([]*model.LOLMEquipment, 0)
	sql := `
SELECT
	equip.equipId,
	equip.name,
	equip.iconPath,
	equip.price,
	equip.description,
	CONCAT_WS(',', equip.searchKey, alias.keywords, alias.keywords_py) AS keywords,
	equip.from,
	equip.into,
	equip.type,
	equip.version,
	equip.fileTime
FROM
	lolm_equipment equip
	LEFT JOIN equip_alias alias ON equip.name = alias.name
WHERE
	equip.version = '%s' and alias.platform = 0
`
	sql = fmt.Sprintf(sql, version)
	err := dao.db.Exec(sql).Find(&result).Error
	return result, err
}

func NewLOLMEquipmentDAO() *LOLMEquipmentDAO {

	d := &LOLMEquipmentDAO{
		db: mysql.DB,
	}
	return d
}
