package dao

import (
	"gorm.io/gorm"
	"sync"
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

var (
	lolETDao  *EquipTypeDAO
	lolETOnce sync.Once
)

func NewEquipTypeDAO() *EquipTypeDAO {
	lolETOnce.Do(func() {
		lolETDao = &EquipTypeDAO{
			db: mysql.DB,
		}
	})
	return lolETDao
}

type EquipType interface {
	Add([]*model.EquipType) (int64, error)
}
