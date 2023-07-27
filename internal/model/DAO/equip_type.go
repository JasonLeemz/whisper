package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type EquipTypeDAO struct {
	db *gorm.DB
}

func (dao *EquipTypeDAO) Add(equips []*model.EquipType) (int64, error) {
	result := mysql.DB.Create(equips)
	return result.RowsAffected, result.Error
}

func NewEquipTypeDAO() *EquipTypeDAO {
	return &EquipTypeDAO{
		db: mysql.DB,
	}
}

type EquipType interface {
	Add(equips []*model.EquipType) (int64, error)
}
