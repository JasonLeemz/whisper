package dao

import (
	"gorm.io/gorm"
	"sync"
	"whisper/internal/model"
	"whisper/pkg/mysql"
)

type HeroAttributeDAO struct {
	db *gorm.DB
}

func (dao *HeroAttributeDAO) FindWithExt(cond map[string]interface{}) ([]*model.HeroAttrWithExt, error) {
	sql := `SELECT
	attr.heroId,
	attr.title,
	attr.name,
	attr.alias,
	attr.shortBio,
	attr.defense,
	attr.magic,
	attr.difficulty,
	attr.difficultyL,
	attr.attack,
	attr.attackrange,
	attr.attackdamage,
	attr.attackspeed,
	attr.attackspeedperlevel,
	attr.hp,
	attr.hpperlevel,
	attr.mp,
	attr.mpperlevel,
	attr.movespeed,
	attr.armor,
	attr.armorperlevel,
	attr.spellblock,
	attr.spellblockperlevel,
	attr.hpregen,
	attr.hpregenperlevel,
	attr.mpregen,
	attr.mpregenperlevel,
	attr.crit,
	attr.damage,
	attr.durability,
	attr.mobility,
	attr.avatar,
	attr.highlightprice,
	attr.goldPrice,
	attr.couponprice,
	attr.isWeekFree,
	attr.platform,
	skin.skinId,
	skin.name AS skin_name,
	skin.description,
	skin.emblemsName,
	skin.mainImg,
	skin.iconImg,
	skin.loadingImg,
	skin.videoImg,
	skin.sourceImg
FROM
	hero_attribute attr
	LEFT JOIN hero_skin skin ON attr.heroId = skin.heroId;
`
	result := make([]*model.HeroAttrWithExt, 0)
	err := dao.db.Raw(sql).Scan(&result).Error

	return result, err
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

func (dao *HeroAttributeDAO) GetMaxVersion() ([]*model.HeroAttribute, error) {
	var result []*model.HeroAttribute
	sql := "select id,`version`,platform,MAX(ctime) as ctime,fileTime from hero_attribute group by platform order by platform asc"
	err := dao.db.Raw(sql).Scan(&result).Error
	return result, err
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

func (dao *HeroAttributeDAO) QueryAllHeroes(cond map[string]interface{}) ([]*model.HeroAttribute, error) {
	tx := dao.db.Model(&model.HeroAttribute{})
	// "DISTINCT(heroId)", "name", "title", "platform",
	var data []*model.HeroAttribute
	tx = tx.Select("heroId", "name", "title", "platform").
		Where(cond).
		Find(&data)
	return data, tx.Error
}

var (
	attrDao  *HeroAttributeDAO
	attrOnce sync.Once
)

func NewHeroAttributeDAO() *HeroAttributeDAO {
	attrOnce.Do(func() {
		attrDao = &HeroAttributeDAO{
			db: mysql.DB,
		}
	})
	return attrDao
}

type HeroAttribute interface {
	Find(query []string, cond map[string]interface{}) ([]*model.HeroAttribute, error)
	FindWithExt(cond map[string]interface{}) ([]*model.HeroAttrWithExt, error)
	Add([]*model.HeroAttribute) (int64, error)
	Delete(cond map[string]interface{}) (int64, error)
	DeleteAndInsert(delCond map[string]interface{}, addData []*model.HeroAttribute) error
	GetMaxVersion() ([]*model.HeroAttribute, error)
	QueryAllHeroes(cond map[string]interface{}) ([]*model.HeroAttribute, error)
}
