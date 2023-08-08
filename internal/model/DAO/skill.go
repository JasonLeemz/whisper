package dao

import (
	"errors"
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type LOLSkillDAO struct {
	db *gorm.DB
}

func (dao *LOLSkillDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLSkill, error) {
	tx := dao.db.Model(&model.LOLSkill{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.LOLSkill
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}

func (dao *LOLSkillDAO) Add(r []*model.LOLSkill) (int64, error) {
	result := dao.db.Create(r)
	return result.RowsAffected, result.Error
}

func (dao *LOLSkillDAO) GetLOLSkillMaxVersion() (*model.LOLSkill, error) {
	tx := dao.db.Model(&model.LOLSkill{})
	var result model.LOLSkill
	tx = tx.Order("version desc").First(&result)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, tx.Error
}

func (dao *LOLSkillDAO) GetLOLSkill(version string) ([]*model.LOLSkill, error) {
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

func NewLOLSkillDAO() *LOLSkillDAO {
	return &LOLSkillDAO{
		db: mysql.DB,
	}
}

type LOLSkill interface {
	Add([]*model.LOLSkill) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLSkill, error)
	GetLOLSkillMaxVersion() (*model.LOLSkill, error)
	GetLOLSkill(version string) ([]*model.LOLSkill, error)
}

// -----------------------------------------

type LOLMSkillDAO struct {
	db *gorm.DB
}

func (dao *LOLMSkillDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLMSkill, error) {
	tx := dao.db.Model(&model.LOLMSkill{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.LOLMSkill
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}

func (dao *LOLMSkillDAO) Add(r []*model.LOLMSkill) (int64, error) {
	result := dao.db.Create(r)
	return result.RowsAffected, result.Error
}

func (dao *LOLMSkillDAO) GetLOLMSkillMaxVersion() (*model.LOLMSkill, error) {
	tx := dao.db.Model(&model.LOLMSkill{})
	var result model.LOLMSkill
	tx = tx.Order("version desc").First(&result)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, tx.Error
}

func (dao *LOLMSkillDAO) GetLOLMSkill(version string) ([]*model.LOLMSkill, error) {
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

func NewLOLMSkillDAO() *LOLMSkillDAO {
	return &LOLMSkillDAO{
		db: mysql.DB,
	}
}

type LOLMSkill interface {
	Add([]*model.LOLMSkill) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLMSkill, error)
	GetLOLMSkillMaxVersion() (*model.LOLMSkill, error)
	GetLOLMSkill(version string) ([]*model.LOLMSkill, error)
}
