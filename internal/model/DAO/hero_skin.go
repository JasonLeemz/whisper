package dao

import (
	"gorm.io/gorm"
	"sync"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type HeroSkinDAO struct {
	db *gorm.DB
}

func (dao *HeroSkinDAO) Find(query []string, cond map[string]interface{}) ([]*model.HeroSkin, error) {
	tx := dao.db.Model(&model.HeroSkin{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.HeroSkin
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}

func (dao *HeroSkinDAO) Add(hr []*model.HeroSkin) (int64, error) {
	result := dao.db.Create(hr)
	return result.RowsAffected, result.Error
}

func (dao *HeroSkinDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.HeroSkin{}, cond)
	return tx.RowsAffected, tx.Error
}

func (dao *HeroSkinDAO) DeleteAndInsert(delCond map[string]interface{}, addData []*model.HeroSkin) error {
	tx := dao.db.Begin()
	tx.Delete(&model.HeroSkin{}, delCond)
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
	hskDao  *HeroSkinDAO
	hskOnce sync.Once
)

func NewHeroSkinDAO() *HeroSkinDAO {
	hskOnce.Do(func() {
		hskDao = &HeroSkinDAO{
			db: mysql.DB,
		}
	})
	return hskDao
}

type HeroSkin interface {
	Find(query []string, cond map[string]interface{}) ([]*model.HeroSkin, error)
	Add(hr []*model.HeroSkin) (int64, error)
	Delete(cond map[string]interface{}) (int64, error)
	DeleteAndInsert(delCond map[string]interface{}, addData []*model.HeroSkin) error
}
