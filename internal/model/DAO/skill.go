package dao

import (
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type LOLSkillDAO struct {
	db *gorm.DB
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
}

// -----------------------------------------

type LOLMSkillDAO struct {
	db *gorm.DB
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
}
