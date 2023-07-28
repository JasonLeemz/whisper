package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type LOLRuneDAO struct {
	db *gorm.DB
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
}

// ---------------------------------------

type LOLMRuneDAO struct {
	db *gorm.DB
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
}
