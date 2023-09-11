package dao

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type EquipAliasDAO struct {
	db *gorm.DB
}

func (dao *EquipAliasDAO) Add(hr []*model.EquipAlias) (int64, error) {
	result := dao.db.Create(hr)
	return result.RowsAffected, result.Error
}

func (dao *EquipAliasDAO) Exists(cond map[string]interface{}) (bool, error) {
	var exists bool
	err := dao.db.Model(&model.EquipAlias{}).Select("count(*) > 0").Where(cond).Find(&exists).Error

	return exists, err
}

func (dao *EquipAliasDAO) Update(hr *model.EquipAlias, cond map[string]interface{}) (int64, error) {
	result := dao.db.Where(cond).Updates(hr)
	return result.RowsAffected, result.Error
}

func (dao *EquipAliasDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.EquipAlias{}, cond)
	return tx.RowsAffected, tx.Error
}

func (dao *EquipAliasDAO) Truncate() error {
	tx := dao.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", new(model.EquipAlias).TableName()))
	return tx.Error
}

var (
	lolEADao  *EquipAliasDAO
	lolEAOnce *sync.Once
)

func NewEquipAliasDAO() *EquipAliasDAO {
	lolEAOnce.Do(func() {
		lolEADao = &EquipAliasDAO{
			db: mysql.DB,
		}
	})
	return lolEADao
}

type EquipAlias interface {
	Add([]*model.EquipAlias) (int64, error)
	Delete(map[string]interface{}) (int64, error)
	Truncate() error
	Update(hr *model.EquipAlias, cond map[string]interface{}) (int64, error)
	Exists(cond map[string]interface{}) (bool, error)
}
