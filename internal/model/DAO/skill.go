package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
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
	tx = tx.Where("status = 0").Order("version desc").First(&result)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, tx.Error
}

func (dao *LOLSkillDAO) GetLOLSkill(version string) ([]*model.LOLSkill, error) {
	// 查当前版本所有数据
	//cond := map[string]interface{}{
	//	"version": version,
	//}
	var result []*model.LOLSkill
	cond := fmt.Sprintf("version = '%s' and status = 0 and gamemode <> ''", version)
	err := dao.db.Where(cond).Find(&result).Error
	//data, err := dao.Find(nil, cond)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (dao *LOLSkillDAO) Update(data *model.LOLSkill, cond map[string]interface{}) (int64, error) {
	result := dao.db.Model(model.LOLSkill{}).Where(cond).Updates(data)
	return result.RowsAffected, result.Error
}

var (
	lolSKDao  *LOLSkillDAO
	lolSKOnce *sync.Once
)

func NewLOLSkillDAO() *LOLSkillDAO {
	lolSKOnce.Do(func() {
		lolSKDao = &LOLSkillDAO{
			db: mysql.DB,
		}
	})
	return lolSKDao
}

type LOLSkill interface {
	Add([]*model.LOLSkill) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLSkill, error)
	Update(data *model.LOLSkill, cond map[string]interface{}) (int64, error)
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

func (dao *LOLMSkillDAO) Update(data *model.LOLMSkill, cond map[string]interface{}) (int64, error) {
	result := dao.db.Model(model.LOLMSkill{}).Where(cond).Updates(data)
	return result.RowsAffected, result.Error
}

func (dao *LOLMSkillDAO) GetLOLMSkillMaxVersion() (*model.LOLMSkill, error) {
	tx := dao.db.Model(&model.LOLMSkill{})
	var result model.LOLMSkill
	tx = tx.Where("status = 0").Order("version desc").First(&result)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, tx.Error
}

func (dao *LOLMSkillDAO) GetLOLMSkill(version string) ([]*model.LOLMSkill, error) {
	// 查当前版本所有数据
	cond := map[string]interface{}{
		"version": version,
		"status":  0,
	}
	data, err := dao.Find(nil, cond)
	if err != nil {
		return nil, err
	}
	return data, err
}

var (
	lolmSKDao  *LOLMSkillDAO
	lolmSKOnce *sync.Once
)

func NewLOLMSkillDAO() *LOLMSkillDAO {
	lolmSKOnce.Do(func() {
		lolmSKDao = &LOLMSkillDAO{
			db: mysql.DB,
		}
	})
	return lolmSKDao
}

type LOLMSkill interface {
	Add([]*model.LOLMSkill) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLMSkill, error)
	Update(data *model.LOLMSkill, cond map[string]interface{}) (int64, error)
	GetLOLMSkillMaxVersion() (*model.LOLMSkill, error)
	GetLOLMSkill(version string) ([]*model.LOLMSkill, error)
}
