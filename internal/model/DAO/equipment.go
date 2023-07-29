package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type LOLEquipment interface {
	Add(data []*model.LOLEquipment) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLEquipment, error)
}

type LOLEquipmentDAO struct {
	db *gorm.DB
}

func (dao *LOLEquipmentDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLEquipment, error) {
	query = append(query, "id")
	result := make([]*model.LOLEquipment, 0)
	tx := dao.db.Model(&model.LOLEquipment{}).Select(query).Where(cond).Find(&result)

	if result != nil && result[0].Id == 0 {
		result = nil
	}
	return result, tx.Error
}

func (dao *LOLEquipmentDAO) Add(equips []*model.LOLEquipment) (int64, error) {
	result := dao.db.Create(equips)
	return result.RowsAffected, result.Error
}

func NewLOLEquipmentDAO() *LOLEquipmentDAO {
	return &LOLEquipmentDAO{
		db: mysql.DB,
	}
}

// --------------------------------------------------------

type LOLMEquipment interface {
	Add([]*model.LOLMEquipment) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLMEquipment, error)
}

type LOLMEquipmentDAO struct {
	db *gorm.DB
}

func (dao *LOLMEquipmentDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLMEquipment, error) {
	query = append(query, "id")
	result := make([]*model.LOLMEquipment, 0)
	tx := dao.db.Model(&model.LOLMEquipment{}).Select(query).Where(cond).Find(&result)

	if result != nil && result[0].Id == 0 {
		result = nil
	}
	return result, tx.Error
}
func (dao *LOLMEquipmentDAO) Add(equips []*model.LOLMEquipment) (int64, error) {
	result := dao.db.Create(equips)
	return result.RowsAffected, result.Error
}

func NewLOLMEquipmentDAO() *LOLMEquipmentDAO {

	d := &LOLMEquipmentDAO{
		db: mysql.DB,
	}
	return d
}
