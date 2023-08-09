package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type HeroSpellDAO struct {
	db *gorm.DB
}

func (dao *HeroSpellDAO) Add(hr []*model.HeroSpell) (int64, error) {
	result := dao.db.Create(hr)
	return result.RowsAffected, result.Error
}

func (dao *HeroSpellDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.HeroSpell{}, cond)
	return tx.RowsAffected, tx.Error
}

func NewHeroSpellDAO() *HeroSpellDAO {
	return &HeroSpellDAO{
		db: mysql.DB,
	}
}

type HeroSpell interface {
	Add([]*model.HeroSpell) (int64, error)
	Delete(map[string]interface{}) (int64, error)
}
