package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type LOLSkillDAO struct {
	db *gorm.DB
}

func (dao *LOLSkillDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLSkill, error) {
	query = append(query, "id")
	result := make([]*model.LOLSkill, 0)
	tx := dao.db.Model(&model.LOLSkill{}).Select(query).Where(cond).Find(&result)
	if result != nil && result[0].Id == 0 {
		result = nil
	}
	return result, tx.Error
}

func (dao *LOLSkillDAO) Add(r []*model.LOLSkill) (int64, error) {
	result := dao.db.Create(r)
	return result.RowsAffected, result.Error
}

func NewLOLSkillDAO() *LOLSkillDAO {
	return &LOLSkillDAO{
		db: mysql.DB,
	}
}

type LOLSkill interface {
	Add([]*model.LOLSkill) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLSkill, error)
}

// -----------------------------------------

type LOLMSkillDAO struct {
	db *gorm.DB
}

func (dao *LOLMSkillDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLMSkill, error) {
	query = append(query, "id")
	result := make([]*model.LOLMSkill, 0)
	tx := dao.db.Model(&model.LOLMSkill{}).Select(query).Where(cond).Find(&result)

	if result != nil && result[0].Id == 0 {
		result = nil
	}
	return result, tx.Error
}

func (dao *LOLMSkillDAO) Add(r []*model.LOLMSkill) (int64, error) {
	result := dao.db.Create(r)
	return result.RowsAffected, result.Error
}

func NewLOLMSkillDAO() *LOLMSkillDAO {
	return &LOLMSkillDAO{
		db: mysql.DB,
	}
}

type LOLMSkill interface {
	Add([]*model.LOLMSkill) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLMSkill, error)
}
