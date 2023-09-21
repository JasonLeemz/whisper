package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type LOLHeroes interface {
	Add([]*model.LOLHeroes) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLHeroes, error)
	Update(data *model.LOLHeroes, cond map[string]interface{}) (int64, error)
	GetLOLHeroesMaxVersion() (*model.LOLHeroes, error)
	GetLOLHeroes(version string) ([]*model.LOLHeroes, error)
	GetLOLHeroesWithExt(version string) ([]*model.LOLHeroesEXT, error)
}

type LOLHeroesDAO struct {
	db *gorm.DB
}

func (dao *LOLHeroesDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLHeroes, error) {
	tx := dao.db.Model(&model.LOLHeroes{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.LOLHeroes
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}

func (dao *LOLHeroesDAO) Add(heroes []*model.LOLHeroes) (int64, error) {
	result := dao.db.Create(heroes)
	return result.RowsAffected, result.Error
}

func (dao *LOLHeroesDAO) Update(data *model.LOLHeroes, cond map[string]interface{}) (int64, error) {
	result := dao.db.Model(model.LOLHeroes{}).Where(cond).Updates(data)
	return result.RowsAffected, result.Error
}

func (dao *LOLHeroesDAO) GetLOLHeroesMaxVersion() (*model.LOLHeroes, error) {
	tx := dao.db.Model(&model.LOLHeroes{})
	var result model.LOLHeroes
	tx = tx.Where("status = 0").Order("version desc").First(&result)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, tx.Error
}

func (dao *LOLHeroesDAO) GetLOLHeroes(version string) ([]*model.LOLHeroes, error) {
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

func (dao *LOLHeroesDAO) GetLOLHeroesWithExt(version string) ([]*model.LOLHeroesEXT, error) {
	// 查当前版本所有数据
	result := make([]*model.LOLHeroesEXT, 0)
	sql := `
SELECT
	GROUP_CONCAT(role.role SEPARATOR ',') AS roles,
	hero.heroId,
	hero.name,
	hero.alias,
	hero.title,
	attr.avatar,
	attr.mainImg,
	hero.goldPrice,
	hero.couponPrice,
	CONCAT_WS(',', hero.keywords, alias.keywords, alias.keywords_py) AS keywords,
	attr.difficulty,
	attr.difficultyL,
	hero.goldPrice,
	hero.couponPrice,
	hero.version,
	hero.fileTime
FROM
	lol_heroes hero
	LEFT JOIN hero_attribute attr ON hero.heroId = attr.heroId
	LEFT JOIN hero_role role ON hero.heroId = role.hero_id
	LEFT JOIN hero_alias alias ON CONCAT_WS(' ', hero.name, hero.title) = alias.name
WHERE
	hero.version = '%s'
	AND hero.status = 0
	AND attr.version = '%s'
	AND attr.platform = 0
	AND alias.platform = 0
GROUP BY
	hero.heroId,
	hero.version
`
	sql = fmt.Sprintf(sql, version, version)
	err := dao.db.Raw(sql).Scan(&result).Error

	return result, err
}

var (
	lolHDao  *LOLHeroesDAO
	lolHOnce sync.Once
)

func NewLOLHeroesDAO() *LOLHeroesDAO {
	lolHOnce.Do(func() {
		lolHDao = &LOLHeroesDAO{
			db: mysql.DB,
		}
	})
	return lolHDao
}

// --------------------------------------------------------

type LOLMHeroes interface {
	Add([]*model.LOLMHeroes) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLMHeroes, error)
	Update(data *model.LOLMHeroes, cond map[string]interface{}) (int64, error)
	GetLOLMHeroesMaxVersion() (*model.LOLMHeroes, error)
	GetLOLMHeroes(version string) ([]*model.LOLMHeroes, error)
	GetLOLMHeroesWithExt(version string) ([]*model.LOLMHeroesEXT, error)
}

type LOLMHeroesDAO struct {
	db *gorm.DB
}

func (dao *LOLMHeroesDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLMHeroes, error) {
	tx := dao.db.Model(&model.LOLMHeroes{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.LOLMHeroes
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}

func (dao *LOLMHeroesDAO) Add(heroes []*model.LOLMHeroes) (int64, error) {
	result := dao.db.Create(heroes)
	return result.RowsAffected, result.Error
}

func (dao *LOLMHeroesDAO) Update(data *model.LOLMHeroes, cond map[string]interface{}) (int64, error) {
	result := dao.db.Model(model.LOLMHeroes{}).Where(cond).Updates(data)
	return result.RowsAffected, result.Error
}

func (dao *LOLMHeroesDAO) GetLOLMHeroesMaxVersion() (*model.LOLMHeroes, error) {
	tx := dao.db.Model(&model.LOLMHeroes{})
	var result model.LOLMHeroes
	tx = tx.Where("status = 0").Order("version desc").First(&result)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, tx.Error
}

func (dao *LOLMHeroesDAO) GetLOLMHeroes(version string) ([]*model.LOLMHeroes, error) {
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

func (dao *LOLMHeroesDAO) GetLOLMHeroesWithExt(version string) ([]*model.LOLMHeroesEXT, error) {
	// 查当前版本所有数据
	result := make([]*model.LOLMHeroesEXT, 0)
	sql := `
SELECT
	GROUP_CONCAT(role.role SEPARATOR ',') AS roles,
	hero.heroId,
	hero.title,
	hero.name,
	hero.alias,
	attr.avatar,
	hero.card,
	hero.poster,
	hero.intro,
	hero.highlightprice,
	hero.couponprice,
	CONCAT_WS(',', hero.alias, hero.title, alias.keywords, alias.keywords_py) AS searchkey,
	hero.lane,
	hero.tags,
	attr.difficulty,
	attr.difficultyL,
	hero.version,
	hero.fileTime
FROM
	lolm_heroes hero
	LEFT JOIN hero_attribute attr ON hero.heroId = attr.heroId
	LEFT JOIN hero_alias alias ON CONCAT_WS(' ', hero.title, hero.name) = alias.name
	LEFT JOIN hero_role role ON hero.heroId = role.hero_id
WHERE
	hero.version = '%s'
	AND hero.status = 0
	AND attr.version = '%s'
	AND attr.platform = 1
	-- AND alias.platform = 1
GROUP BY
	hero.heroId,
	hero.version;
`
	sql = fmt.Sprintf(sql, version, version)
	err := dao.db.Raw(sql).Scan(&result).Error

	return result, err
}

var (
	lolmHDao  *LOLMHeroesDAO
	lolmHOnce sync.Once
)

func NewLOLMHeroesDAO() *LOLMHeroesDAO {
	lolmHOnce.Do(func() {
		lolmHDao = &LOLMHeroesDAO{
			db: mysql.DB,
		}
	})
	return lolmHDao
}
