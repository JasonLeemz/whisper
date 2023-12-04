package dao

import (
	"gorm.io/gorm"
	"sync"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type AuthorSpaceDAO struct {
	db *gorm.DB
}

func (dao *AuthorSpaceDAO) Add(hr []*model.AuthorSpace) (int64, error) {
	result := dao.db.Create(hr)
	return result.RowsAffected, result.Error
}

func (dao *AuthorSpaceDAO) Exists(cond map[string]interface{}) (bool, error) {
	var exists bool
	err := dao.db.Model(&model.AuthorSpace{}).Select("count(*) > 0").Where(cond).Find(&exists).Error

	return exists, err
}

func (dao *AuthorSpaceDAO) Update(hr *model.AuthorSpace, cond map[string]interface{}) (int64, error) {
	result := dao.db.Where(cond).Updates(hr)
	return result.RowsAffected, result.Error
}

func (dao *AuthorSpaceDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.AuthorSpace{}, cond)
	return tx.RowsAffected, tx.Error
}

func (dao *AuthorSpaceDAO) Find(query []string, cond map[string]interface{}) ([]*model.AuthorSpace, error) {
	tx := dao.db.Model(&model.AuthorSpace{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.AuthorSpace
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}

var (
	authorSpaceDao  *AuthorSpaceDAO
	authorSpaceOnce sync.Once
)

func NewAuthorSpaceDAO() *AuthorSpaceDAO {
	authorSpaceOnce.Do(func() {
		authorSpaceDao = &AuthorSpaceDAO{
			db: mysql.DB,
		}
	})
	return authorSpaceDao
}

type AuthorSpace interface {
	Add([]*model.AuthorSpace) (int64, error)
	Delete(map[string]interface{}) (int64, error)
	Update(hr *model.AuthorSpace, cond map[string]interface{}) (int64, error)
	Exists(cond map[string]interface{}) (bool, error)
	Find(query []string, cond map[string]interface{}) ([]*model.AuthorSpace, error)
}
