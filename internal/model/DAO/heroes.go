package dao

import (
	"errors"
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type LOLHeroesDAO struct {
	db *gorm.DB
}

func (dao *LOLHeroesDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLHeroes, error) {
	tx := dao.db.Model(&model.LOLHeroes{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.LOLHeroes
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}

func (dao *LOLHeroesDAO) Add(heroes []*model.LOLHeroes) (int64, error) {
	result := dao.db.Create(heroes)
	return result.RowsAffected, result.Error
}

func (dao *LOLHeroesDAO) GetLOLHeroesMaxVersion() (*model.LOLHeroes, error) {
	tx := dao.db.Model(&model.LOLHeroes{})
	var result model.LOLHeroes
	tx = tx.Order("version desc").First(&result)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, tx.Error
}

func (dao *LOLHeroesDAO) GetLOLHeroes(version string) ([]*model.LOLHeroes, error) {
	// 查当前版本所有数据
	cond := map[string]interface{}{
		"version": version,
	}
	data, err := dao.Find(nil, cond)
	if err != nil {
		return nil, err
	}
	return data, err
}

func NewLOLHeroesDAO() *LOLHeroesDAO {
	return &LOLHeroesDAO{
		db: mysql.DB,
	}
}

type LOLHeroes interface {
	Add([]*model.LOLHeroes) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLHeroes, error)
	GetLOLHeroesMaxVersion() (*model.LOLHeroes, error)
	GetLOLHeroes(version string) ([]*model.LOLHeroes, error)
}

// --------------------------------------------------------

type LOLMHeroesDAO struct {
	db *gorm.DB
}

func (dao *LOLMHeroesDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLMHeroes, error) {
	tx := dao.db.Model(&model.LOLMHeroes{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.LOLMHeroes
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}

func (dao *LOLMHeroesDAO) Add(heroes []*model.LOLMHeroes) (int64, error) {
	result := dao.db.Create(heroes)
	return result.RowsAffected, result.Error
}

func (dao *LOLMHeroesDAO) GetLOLMHeroesMaxVersion() (*model.LOLMHeroes, error) {
	tx := dao.db.Model(&model.LOLMHeroes{})
	var result model.LOLMHeroes
	tx = tx.Order("version desc").First(&result)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, tx.Error
}

func (dao *LOLMHeroesDAO) GetLOLMHeroes(version string) ([]*model.LOLMHeroes, error) {
	// 查当前版本所有数据
	cond := map[string]interface{}{
		"version": version,
	}
	data, err := dao.Find(nil, cond)
	if err != nil {
		return nil, err
	}
	return data, err
}

func NewLOLMHeroesDAO() *LOLMHeroesDAO {
	return &LOLMHeroesDAO{
		db: mysql.DB,
	}
}

type LOLMHeroes interface {
	Add([]*model.LOLMHeroes) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLMHeroes, error)
	GetLOLMHeroesMaxVersion() (*model.LOLMHeroes, error)
	GetLOLMHeroes(version string) ([]*model.LOLMHeroes, error)
}
