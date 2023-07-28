package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type LOLHeroesDAO struct {
	db *gorm.DB
}

func (dao *LOLHeroesDAO) Add(heroes []*model.LOLHeroes) (int64, error) {
	result := dao.db.Create(heroes)
	return result.RowsAffected, result.Error
}

func NewLOLHeroesDAO() *LOLHeroesDAO {
	return &LOLHeroesDAO{
		db: mysql.DB,
	}
}

type LOLHeroes interface {
	Add([]*model.LOLHeroes) (int64, error)
}

// --------------------------------------------------------

type LOLMHeroesDAO struct {
	db *gorm.DB
}

func (dao *LOLMHeroesDAO) Add(heroes []*model.LOLMHeroes) (int64, error) {
	result := dao.db.Create(heroes)
	return result.RowsAffected, result.Error
}

func NewLOLMHeroesDAO() *LOLMHeroesDAO {
	return &LOLMHeroesDAO{
		db: mysql.DB,
	}
}

type LOLMHeroes interface {
	Add([]*model.LOLMHeroes) (int64, error)
}
