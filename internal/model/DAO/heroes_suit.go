package dao

import (
	"fmt"
	"gorm.io/gorm"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type HeroesSuitDAO struct {
	db *gorm.DB
}

func (dao *HeroesSuitDAO) DeleteAndInsert(delCond map[string]interface{}, addData []*model.HeroesSuit) error {
	tx := dao.db.Begin()
	tx.Delete(&model.HeroesSuit{}, delCond)
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
func (dao *HeroesSuitDAO) Add(hs []*model.HeroesSuit) (int64, error) {
	result := dao.db.Create(hs)
	return result.RowsAffected, result.Error
}
func (dao *HeroesSuitDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.HeroesSuit{}, cond)
	return tx.RowsAffected, tx.Error
}
func (dao *HeroesSuitDAO) GetSuitForHero(platform int, heroID string) ([]*model.HeroesSuit, error) {
	sql := `
SELECT
	suit.heroId,
	suit.pos,
	suit.title,
	suit.author,
	suit.author_icon,
	suit.recommend_id,
	suit.runeids,
	suit.skillids,
	suit.itemids,
	suit.igamecnt,
	suit.wincnt,
	suit.winrate,
	suit.allcnt,
	suit.showrate,
	suit.desc,
	suit.type,
	suit.platform
FROM
	heroes_suit suit
	%s
WHERE
	suit.heroId = %s
ORDER BY
	winrate desc
`

	rightJoin := "RIGHT JOIN heroes_position pos ON suit.pos = pos.pos AND suit.heroId = pos.heroId"
	if platform == 0 {
		sql = fmt.Sprintf(sql, rightJoin, heroID)
	} else {
		sql = fmt.Sprintf(sql, "", heroID)
	}

	result := make([]*model.HeroesSuit, 0)
	err := dao.db.Raw(sql).Scan(&result).Error

	return result, err
}
func NewHeroesSuitDAO() *HeroesSuitDAO {
	return &HeroesSuitDAO{
		db: mysql.DB,
	}
}

type HeroesSuit interface {
	Add([]*model.HeroesSuit) (int64, error)
	Delete(cond map[string]interface{}) (int64, error)
	DeleteAndInsert(delCond map[string]interface{}, addData []*model.HeroesSuit) error
	GetSuitForHero(platform int, heroID string) ([]*model.HeroesSuit, error)
}
