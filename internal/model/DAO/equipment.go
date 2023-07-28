package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type LOLEquipment interface {
	Add([]*model.LOLEquipment) (int64, error)
}

type LOLEquipmentDAO struct {
	db *gorm.DB
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
}

type LOLMEquipmentDAO struct {
	db *gorm.DB
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
