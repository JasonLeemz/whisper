package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type RuneTypeDAO struct {
	db *gorm.DB
}

func (dao *RuneTypeDAO) Add(hr []*model.RuneType) (int64, error) {
	result := dao.db.Create(hr)
	return result.RowsAffected, result.Error
}

func (dao *RuneTypeDAO) DeleteAll(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.RuneType{}, cond)
	return tx.RowsAffected, tx.Error
}

func NewRuneTypeDAO() *RuneTypeDAO {
	return &RuneTypeDAO{
		db: mysql.DB,
	}
}

type RuneType interface {
	Add([]*model.RuneType) (int64, error)
	DeleteAll(map[string]interface{}) (int64, error)
}
