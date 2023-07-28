package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type EquipTypeDAO struct {
	db *gorm.DB
}

func (dao *EquipTypeDAO) Add(et []*model.EquipType) (int64, error) {
	result := dao.db.Create(et)
	return result.RowsAffected, result.Error
}

func NewEquipTypeDAO() *EquipTypeDAO {
	return &EquipTypeDAO{
		db: mysql.DB,
	}
}

type EquipType interface {
	Add([]*model.EquipType) (int64, error)
}
