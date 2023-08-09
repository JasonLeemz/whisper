package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type HeroAttributeDAO struct {
	db *gorm.DB
}

func (dao *HeroAttributeDAO) Add(et []*model.HeroAttribute) (int64, error) {
	result := dao.db.Create(et)
	return result.RowsAffected, result.Error
}

func (dao *HeroAttributeDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.HeroAttribute{}, cond)
	return tx.RowsAffected, tx.Error
}

func NewHeroAttributeDAO() *HeroAttributeDAO {
	return &HeroAttributeDAO{
		db: mysql.DB,
	}
}

type HeroAttribute interface {
	Add([]*model.HeroAttribute) (int64, error)
	Delete(cond map[string]interface{}) (int64, error)
}
