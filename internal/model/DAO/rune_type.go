package dao

import (
	"gorm.io/gorm"
	"sync"
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

var (
	rtDao  *RuneTypeDAO
	rtOnce sync.Once
)

func NewRuneTypeDAO() *RuneTypeDAO {
	rtOnce.Do(func() {
		rtDao = &RuneTypeDAO{
			db: mysql.DB,
		}
	})
	return rtDao
}

type RuneType interface {
	Add([]*model.RuneType) (int64, error)
	DeleteAll(map[string]interface{}) (int64, error)
}
