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
	result := dao.db.Create(hr)
	return result.RowsAffected, result.Error
}

func (dao *HeroRoleDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.HeroRole{}, cond)
	return tx.RowsAffected, tx.Error
}

func NewHeroRoleDAO() *HeroRoleDAO {
	return &HeroRoleDAO{
		db: mysql.DB,
	}
}

type HeroRole interface {
	Add([]*model.HeroRole) (int64, error)
	Delete(map[string]interface{}) (int64, error)
}
