package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type EquipmentDAO struct {
	db *gorm.DB
}

func (dao *EquipmentDAO) Add(equips []*model.Equipment) (int64, error) {
	result := mysql.DB.Create(equips)
	return result.RowsAffected, result.Error
}

func NewEquipmentDAO() *EquipmentDAO {
	return &EquipmentDAO{
		db: mysql.DB,
	}
}

type Equipment interface {
	Add(equips []*model.Equipment) (int64, error)
}
