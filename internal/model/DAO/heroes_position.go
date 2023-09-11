package dao

import (
	"gorm.io/gorm"
	"sync"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type HeroesPositionDAO struct {
	db *gorm.DB
}

func (dao *HeroesPositionDAO) DeleteAndInsert(delCond map[string]interface{}, addData []*model.HeroesPosition) error {
	tx := dao.db.Begin()
	tx.Delete(&model.HeroesPosition{}, delCond)
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
func (dao *HeroesPositionDAO) Add(hs []*model.HeroesPosition) (int64, error) {
	result := dao.db.Create(hs)
	return result.RowsAffected, result.Error
}

func (dao *HeroesPositionDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.HeroesPosition{}, cond)
	return tx.RowsAffected, tx.Error
}

var (
	hpDao  *HeroesPositionDAO
	hpOnce *sync.Once
)

func NewHeroesPositionDAO() *HeroesPositionDAO {
	hpOnce.Do(func() {
		hpDao = &HeroesPositionDAO{
			db: mysql.DB,
		}
	})
	return hpDao
}

type HeroesPosition interface {
	Add([]*model.HeroesPosition) (int64, error)
	Delete(cond map[string]interface{}) (int64, error)
	DeleteAndInsert(delCond map[string]interface{}, addData []*model.HeroesPosition) error
}
