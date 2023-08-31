package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type HeroesPositionDAO struct {
	db *gorm.DB
}

func (dao *HeroesPositionDAO) Add(hs []*model.HeroesPosition) (int64, error) {
	result := dao.db.Create(hs)
	return result.RowsAffected, result.Error
}
func (dao *HeroesPositionDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.HeroesPosition{}, cond)
	return tx.RowsAffected, tx.Error
}
func NewHeroesPositionDAO() *HeroesPositionDAO {
	return &HeroesPositionDAO{
		db: mysql.DB,
	}
}

type HeroesPosition interface {
	Add([]*model.HeroesPosition) (int64, error)
	Delete(cond map[string]interface{}) (int64, error)
}
