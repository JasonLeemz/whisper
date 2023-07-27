package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type HeroesDAO struct {
	db *gorm.DB
}

func (dao *HeroesDAO) Add(heroes []*model.Heroes) (int64, error) {
	result := mysql.DB.Create(heroes)
	return result.RowsAffected, result.Error
}

func NewHeroesDAO() *HeroesDAO {
	return &HeroesDAO{
		db: mysql.DB,
	}
}

type Heroes interface {
	Add([]*model.Heroes) (int64, error)
}
