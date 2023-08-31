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

func (dao *HeroesSuitDAO) Add(hs []*model.HeroesSuit) (int64, error) {
	result := dao.db.Create(hs)
	return result.RowsAffected, result.Error
}
func (dao *HeroesSuitDAO) Delete(cond map[string]interface{}) (int64, error) {
	tx := dao.db.Delete(&model.HeroesSuit{}, cond)
	return tx.RowsAffected, tx.Error
}
func (dao *HeroesSuitDAO) GetSuitForHero(heroID string) ([]*model.HeroesSuit, error) {
	sql := `
SELECT
	suit.heroId,
	suit.pos,
	suit.itemids,
	suit.igamecnt,
	suit.wincnt,
	suit.winrate,
	suit.allcnt,
	suit.showrate,
	suit.type,
	suit.platform
FROM
	heroes_suit suit
	RIGHT JOIN heroes_position pos ON suit.pos = pos.pos
		AND suit.heroId = pos.heroId
WHERE
	suit.heroId = %s
ORDER BY
	showrate desc
`
	sql = fmt.Sprintf(sql, heroID)

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
	GetSuitForHero(heroID string) ([]*model.HeroesSuit, error)
}
