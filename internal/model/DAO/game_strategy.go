package dao

import (
	"gorm.io/gorm"
	"sync"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type GameStrategyDAO struct {
	db *gorm.DB
}

func (dao *GameStrategyDAO) InsertORIgnore(strategy *model.GameStrategy) (int64, error) {
	tx := dao.db.Begin()
	var exists bool
	err := tx.Select("count(*) > 0").
		Where(map[string]interface{}{
			"bvid": strategy.Bvid,
		}).
		Find(&exists).Error
	if exists {
		tx.Rollback()
		return 0, err
	}

	// 插入
	tx.Create(strategy)
	if tx.Error != nil {
		tx.Rollback()
		return 0, tx.Error
	}
	tx.Commit()

	return tx.RowsAffected, nil
}

func (dao *GameStrategyDAO) Exists(cond map[string]interface{}) (bool, error) {
	var exists bool
	err := dao.db.Model(&model.GameStrategy{}).Select("count(*) > 0").Where(cond).Find(&exists).Error

	return exists, err
}

func (dao *GameStrategyDAO) Update(hr *model.GameStrategy, cond map[string]interface{}) (int64, error) {
	result := dao.db.Where(cond).Updates(hr)
	return result.RowsAffected, result.Error
}

func (dao *GameStrategyDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.GameStrategy{}, cond)
	return tx.RowsAffected, tx.Error
}

func (dao *GameStrategyDAO) Find(query []string, cond map[string]interface{}) ([]*model.GameStrategy, error) {
	//select * from `game_strategy`
	//where `hero`='安妮' and `platform` = 1 and status = 0
	//order by `public_date` desc,`played` desc;
	tx := dao.db.Model(&model.GameStrategy{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.GameStrategy
	tx = tx.Where(cond).Order("public_date desc, played desc").Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}

var (
	gameStrategyDao  *GameStrategyDAO
	gameStrategyOnce sync.Once
)

func NewGameStrategyDAO() *GameStrategyDAO {
	gameStrategyOnce.Do(func() {
		gameStrategyDao = &GameStrategyDAO{
			db: mysql.DB,
		}
	})
	return gameStrategyDao
}

type GameStrategy interface {
	InsertORIgnore(*model.GameStrategy) (int64, error)
	Delete(map[string]interface{}) (int64, error)
	Update(hr *model.GameStrategy, cond map[string]interface{}) (int64, error)
	Exists(cond map[string]interface{}) (bool, error)
	Find(query []string, cond map[string]interface{}) ([]*model.GameStrategy, error)
}
