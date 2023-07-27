package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type HeroesSuitDAO struct {
	db *gorm.DB
}

func (dao *HeroesSuitDAO) Add(equips []*model.HeroesSuit) (int64, error) {
	result := mysql.DB.Create(equips)
	return result.RowsAffected, result.Error
}

func NewHeroesSuitDAO() *HeroesSuitDAO {
	return &HeroesSuitDAO{
		db: mysql.DB,
	}
}

type HeroesSuit interface {
	Add(equips []*model.HeroesSuit) (int64, error)
}
