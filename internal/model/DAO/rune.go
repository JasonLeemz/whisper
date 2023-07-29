package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type LOLRuneDAO struct {
	db *gorm.DB
}

func (dao *LOLRuneDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLRune, error) {
	query = append(query, "id")
	result := make([]*model.LOLRune, 0)
	tx := dao.db.Model(&model.LOLRune{}).Select(query).Where(cond).Find(&result)

	if result != nil && result[0].Id == 0 {
		result = nil
	}
	return result, tx.Error
}

func (dao *LOLRuneDAO) Add(r []*model.LOLRune) (int64, error) {
	result := dao.db.Create(r)
	return result.RowsAffected, result.Error
}

func NewLOLRuneDAO() *LOLRuneDAO {
	return &LOLRuneDAO{
		db: mysql.DB,
	}
}

type LOLRune interface {
	Add([]*model.LOLRune) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLRune, error)
}

// ---------------------------------------

type LOLMRuneDAO struct {
	db *gorm.DB
}

func (dao *LOLMRuneDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLMRune, error) {
	query = append(query, "id")
	result := make([]*model.LOLMRune, 0)
	tx := dao.db.Model(&model.LOLMRune{}).Select(query).Where(cond).Find(&result)

	if result != nil && result[0].Id == 0 {
		result = nil
	}
	return result, tx.Error
}

func (dao *LOLMRuneDAO) Add(r []*model.LOLMRune) (int64, error) {
	result := dao.db.Create(r)
	return result.RowsAffected, result.Error
}

func NewLOLMRuneDAO() *LOLMRuneDAO {
	return &LOLMRuneDAO{
		db: mysql.DB,
	}
}

type LOLMRune interface {
	Add([]*model.LOLMRune) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLMRune, error)
}
