package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type HeroesSuitDAO struct {
	db *gorm.DB
}

func (dao *HeroesSuitDAO) Add(hs []*model.HeroesSuit) (int64, error) {
	result := dao.db.Create(hs)
	return result.RowsAffected, result.Error
}

func NewHeroesSuitDAO() *HeroesSuitDAO {
	return &HeroesSuitDAO{
		db: mysql.DB,
	}
}

type HeroesSuit interface {
	Add([]*model.HeroesSuit) (int64, error)
}
