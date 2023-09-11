package dao

import (
	"gorm.io/gorm"
	"sync"
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

func (dao *HeroSpellDAO) GetSpells(heroID string) ([]*model.HeroSpell, error) {
	var result []*model.HeroSpell
	err := dao.db.Where("heroId", heroID).Find(&result).Order("sort desc").Error
	return result, err
}

func (dao *HeroSpellDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.HeroSpell{}, cond)
	return tx.RowsAffected, tx.Error
}

func (dao *HeroSpellDAO) DeleteAndInsert(delCond map[string]interface{}, addData []*model.HeroSpell) error {
	tx := dao.db.Begin()
	tx.Delete(&model.HeroSpell{}, delCond)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Create(addData)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()

	return nil
}

var (
	hsDao  *HeroSpellDAO
	hsOnce *sync.Once
)

func NewHeroSpellDAO() *HeroSpellDAO {
	hsOnce.Do(func() {
		hsDao = &HeroSpellDAO{
			db: mysql.DB,
		}
	})
	return hsDao
}

type HeroSpell interface {
	Add(hr []*model.HeroSpell) (int64, error)
	Delete(cond map[string]interface{}) (int64, error)
	DeleteAndInsert(delCond map[string]interface{}, addData []*model.HeroSpell) error
	GetSpells(heroID string) ([]*model.HeroSpell, error)
}
