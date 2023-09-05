package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type HeroAttributeDAO struct {
	db *gorm.DB
}

func (dao *HeroAttributeDAO) Find(query []string, cond map[string]interface{}) ([]*model.HeroAttribute, error) {
	tx := dao.db.Model(&model.HeroAttribute{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.HeroAttribute
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}

func (dao *HeroAttributeDAO) Add(et []*model.HeroAttribute) (int64, error) {
	result := dao.db.Create(et)
	return result.RowsAffected, result.Error
}

func (dao *HeroAttributeDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.HeroAttribute{}, cond)
	return tx.RowsAffected, tx.Error
}

func (dao *HeroAttributeDAO) DeleteAndInsert(delCond map[string]interface{}, addData []*model.HeroAttribute) error {
	tx := dao.db.Begin()
	tx.Delete(&model.HeroAttribute{}, delCond)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Create(addData)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()

	return nil
}

func NewHeroAttributeDAO() *HeroAttributeDAO {
	return &HeroAttributeDAO{
		db: mysql.DB,
	}
}

type HeroAttribute interface {
	Find(query []string, cond map[string]interface{}) ([]*model.HeroAttribute, error)
	Add([]*model.HeroAttribute) (int64, error)
	Delete(cond map[string]interface{}) (int64, error)
	DeleteAndInsert(delCond map[string]interface{}, addData []*model.HeroAttribute) error
}
