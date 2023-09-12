package dao

import (
	"gorm.io/gorm"
	"sync"
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

func (dao *HeroRoleDAO) DeleteAndInsert(delCond map[string]interface{}, addData []*model.HeroRole) error {
	tx := dao.db.Begin()
	tx.Delete(&model.HeroRole{}, delCond)
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
func (dao *HeroRoleDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.HeroRole{}, cond)
	return tx.RowsAffected, tx.Error
}

var (
	hrDao  *HeroRoleDAO
	hrOnce sync.Once
)

func NewHeroRoleDAO() *HeroRoleDAO {
	hrOnce.Do(func() {
		hrDao = &HeroRoleDAO{
			db: mysql.DB,
		}
	})
	return hrDao
}

type HeroRole interface {
	Add([]*model.HeroRole) (int64, error)
	Delete(map[string]interface{}) (int64, error)
	DeleteAndInsert(delCond map[string]interface{}, addData []*model.HeroRole) error
}
