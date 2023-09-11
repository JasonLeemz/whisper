package dao

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type HeroAliasDAO struct {
	db *gorm.DB
}

func (dao *HeroAliasDAO) Add(hr []*model.HeroAlias) (int64, error) {
	result := dao.db.Create(hr)
	return result.RowsAffected, result.Error
}

func (dao *HeroAliasDAO) Exists(cond map[string]interface{}) (bool, error) {
	var exists bool
	err := dao.db.Model(&model.HeroAlias{}).Select("count(*) > 0").Where(cond).Find(&exists).Error

	return exists, err
}

func (dao *HeroAliasDAO) Update(hr *model.HeroAlias, cond map[string]interface{}) (int64, error) {
	result := dao.db.Where(cond).Updates(hr)
	return result.RowsAffected, result.Error
}

func (dao *HeroAliasDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.HeroAlias{}, cond)
	return tx.RowsAffected, tx.Error
}

func (dao *HeroAliasDAO) Truncate() error {
	tx := dao.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", new(model.HeroAlias).TableName()))
	return tx.Error
}

var (
	haDao  *HeroAliasDAO
	haOnce *sync.Once
)

func NewHeroAliasDAO() *HeroAliasDAO {
	haOnce.Do(func() {
		haDao = &HeroAliasDAO{
			db: mysql.DB,
		}
	})
	return haDao
}

type HeroAlias interface {
	Add([]*model.HeroAlias) (int64, error)
	Delete(map[string]interface{}) (int64, error)
	Truncate() error
	Update(hr *model.HeroAlias, cond map[string]interface{}) (int64, error)
	Exists(cond map[string]interface{}) (bool, error)
}
