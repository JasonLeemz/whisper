package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type HeroRoleDAO struct {
	db *gorm.DB
}

func (dao *HeroRoleDAO) Add(hr []*model.HeroRole) (int64, error) {
	result := mysql.DB.Create(hr)
	return result.RowsAffected, result.Error
}

func NewHeroRoleDAO() *HeroRoleDAO {
	return &HeroRoleDAO{
		db: mysql.DB,
	}
}

type HeroRole interface {
	Add([]*model.HeroRole) (int64, error)
}
