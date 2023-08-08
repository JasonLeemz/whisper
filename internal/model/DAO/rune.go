package dao

import (
	"errors"
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type LOLRuneDAO struct {
	db *gorm.DB
}

func (dao *LOLRuneDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLRune, error) {
	tx := dao.db.Model(&model.LOLRune{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.LOLRune
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}

func (dao *LOLRuneDAO) Add(r []*model.LOLRune) (int64, error) {
	result := dao.db.Create(r)
	return result.RowsAffected, result.Error
}

func (dao *LOLRuneDAO) GetLOLRuneMaxVersion() (*model.LOLRune, error) {
	tx := dao.db.Model(&model.LOLRune{})
	var result model.LOLRune
	tx = tx.Order("version desc").First(&result)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, tx.Error
}

func (dao *LOLRuneDAO) GetLOLRune(version string) ([]*model.LOLRune, error) {
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
func NewLOLRuneDAO() *LOLRuneDAO {
	return &LOLRuneDAO{
		db: mysql.DB,
	}
}

type LOLRune interface {
	Add([]*model.LOLRune) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLRune, error)
	GetLOLRuneMaxVersion() (*model.LOLRune, error)
	GetLOLRune(version string) ([]*model.LOLRune, error)
}

// ---------------------------------------

type LOLMRuneDAO struct {
	db *gorm.DB
}

func (dao *LOLMRuneDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLMRune, error) {
	tx := dao.db.Model(&model.LOLMRune{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.LOLMRune
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}

func (dao *LOLMRuneDAO) Add(r []*model.LOLMRune) (int64, error) {
	result := dao.db.Create(r)
	return result.RowsAffected, result.Error
}

func (dao *LOLMRuneDAO) GetLOLMRuneMaxVersion() (*model.LOLMRune, error) {
	tx := dao.db.Model(&model.LOLMRune{})
	var result model.LOLMRune
	tx = tx.Order("version desc").First(&result)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, tx.Error
}

func (dao *LOLMRuneDAO) GetLOLMRune(version string) ([]*model.LOLMRune, error) {
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

func NewLOLMRuneDAO() *LOLMRuneDAO {
	return &LOLMRuneDAO{
		db: mysql.DB,
	}
}

type LOLMRune interface {
	Add([]*model.LOLMRune) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLMRune, error)
	GetLOLMRuneMaxVersion() (*model.LOLMRune, error)
	GetLOLMRune(version string) ([]*model.LOLMRune, error)
}
