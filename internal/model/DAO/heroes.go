package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type LOLHeroesDAO struct {
	db *gorm.DB
}

func (dao *LOLHeroesDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLHeroes, error) {
	query = append(query, "id")
	result := make([]*model.LOLHeroes, 0)
	tx := dao.db.Model(&model.LOLHeroes{}).Select(query).Where(cond).Find(&result)

	if result != nil && result[0].Id == 0 {
		result = nil
	}
	return result, tx.Error
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
	Find(query []string, cond map[string]interface{}) ([]*model.LOLHeroes, error)
}

// --------------------------------------------------------

type LOLMHeroesDAO struct {
	db *gorm.DB
}

func (dao *LOLMHeroesDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLMHeroes, error) {
	query = append(query, "id")
	result := make([]*model.LOLMHeroes, 0)
	tx := dao.db.Model(&model.LOLMHeroes{}).Select(query).Where(cond).Find(&result)

	if result != nil && result[0].Id == 0 {
		result = nil
	}
	return result, tx.Error
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
	Find(query []string, cond map[string]interface{}) ([]*model.LOLMHeroes, error)
}
