package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"sync"
	"whisper/internal/model"
	"whisper/pkg/config"
	"whisper/pkg/mysql"
)

type LOLEquipment interface {
	Add(data []*model.LOLEquipment) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLEquipment, error)
	Update(data *model.LOLEquipment, cond map[string]interface{}) (int64, error)
	GetLOLEquipmentMaxVersion() (*model.LOLEquipment, error)
	GetLOLEquipment(version string) ([]*model.LOLEquipment, error)
	GetLOLEquipmentWithExt(version string) ([]*model.LOLEquipment, error)
	GetRoadmap(version string, id string) (map[string][]*model.LOLMEquipment, error)
}

type LOLEquipmentDAO struct {
	db *gorm.DB
}

func (dao *LOLEquipmentDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLEquipment, error) {
	tx := dao.db.Model(&model.LOLEquipment{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.LOLEquipment
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}

func (dao *LOLEquipmentDAO) Add(equips []*model.LOLEquipment) (int64, error) {
	result := dao.db.Create(equips)
	return result.RowsAffected, result.Error
}

func (dao *LOLEquipmentDAO) Update(data *model.LOLEquipment, cond map[string]interface{}) (int64, error) {
	result := dao.db.Model(model.LOLEquipment{}).Where(cond).Updates(data)
	return result.RowsAffected, result.Error
}

func (dao *LOLEquipmentDAO) GetLOLEquipmentMaxVersion() (*model.LOLEquipment, error) {
	tx := dao.db.Model(&model.LOLEquipment{})
	var result model.LOLEquipment
	tx = tx.Where("status = 0").Order("version desc").First(&result)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, tx.Error
}

func (dao *LOLEquipmentDAO) GetLOLEquipment(version string) ([]*model.LOLEquipment, error) {
	// 查当前版本所有数据
	cond := map[string]interface{}{
		"version": version,
		"status":  0,
	}

	not := map[string]interface{}{
		"itemId": config.EquipDict.Exclude,
	}

	var equip []*model.LOLEquipment
	tx := dao.db.Model(model.LOLEquipment{}).Not(not).Find(&equip, cond)
	//equip, err := dao.Find(nil, cond)
	//if err != nil {
	//	return nil, err
	//}
	return equip, tx.Error
}

func (dao *LOLEquipmentDAO) GetRoadmap(version string, id string) (map[string][]*model.LOLEquipment, error) {
	result := make(map[string][]*model.LOLEquipment)

	current, err := dao.Find(nil, map[string]interface{}{
		"status":  0,
		"version": version,
		"itemId":  id,
	})
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, errors.New("can not find id " + id)
	}
	result["current"] = current

	from, err := dao.Find(nil, map[string]interface{}{
		"status":  0,
		"version": version,
		"itemId":  strings.Split(current[0].From, ","),
	})
	if err != nil {
		return nil, err
	}
	result["from"] = from

	into, err := dao.Find(nil, map[string]interface{}{
		"status":  0,
		"version": version,
		"itemId":  strings.Split(current[0].Into, ","),
	})
	if err != nil {
		return nil, err
	}
	result["into"] = into

	return result, nil
}

func (dao *LOLEquipmentDAO) GetLOLEquipmentWithExt(version string) ([]*model.LOLEquipment, error) {

	notin := strings.Join(config.EquipDict.Exclude, ",")
	// 查当前版本所有数据
	result := make([]*model.LOLEquipment, 0)
	sql := `
SELECT
	equip.itemId,
	equip.name,
	CONCAT_WS(',', equip.keywords, alias.keywords, alias.keywords_py) AS keywords,
	equip.description,
	equip.plaintext,
	equip.iconPath,
	equip.price,
	equip.sell,
	equip.total,
	equip.suitHeroId,
	equip.maps,
	equip.from,
	equip.into,
	equip.types,
	equip.version,
	equip.fileTime
FROM
	lol_equipment equip
	LEFT JOIN equip_alias alias ON equip.name = alias.name
WHERE
	equip.version = '%s' and equip.status = 0 and total <> 0 and equip.maps <> '' and alias.platform = 0
	and equip.itemId not in (%s)
`
	sql = fmt.Sprintf(sql, version, notin)
	err := dao.db.Exec(sql).Find(&result).Error
	return result, err
}

var (
	lolEDao  *LOLEquipmentDAO
	lolEOnce *sync.Once
)

func NewLOLEquipmentDAO() *LOLEquipmentDAO {
	lolEOnce.Do(func() {
		lolEDao = &LOLEquipmentDAO{
			db: mysql.DB,
		}
	})
	return lolEDao
}

// --------------------------------------------------------

type LOLMEquipment interface {
	Add(equips []*model.LOLMEquipment) (int64, error)
	Find(query []string, cond map[string]interface{}) ([]*model.LOLMEquipment, error)
	Update(data *model.LOLMEquipment, cond map[string]interface{}) (int64, error)
	UpdatesInto(fileTime, version string, ups map[string][]string) (int64, error)
	GetLOLMEquipmentMaxVersion() (*model.LOLMEquipment, error)
	GetLOLMEquipment(version string) ([]*model.LOLMEquipment, error)
	GetLOLMEquipmentWithExt(version string) ([]*model.LOLMEquipment, error)
	GetRoadmap(version string, id string) (map[string][]*model.LOLMEquipment, error)
}

type LOLMEquipmentDAO struct {
	db *gorm.DB
}

func (dao *LOLMEquipmentDAO) Find(query []string, cond map[string]interface{}) ([]*model.LOLMEquipment, error) {
	tx := dao.db.Model(&model.LOLMEquipment{})
	if query != nil {
		query = append(query, "id")
		tx = tx.Select(query)
	}
	var result []*model.LOLMEquipment
	tx = tx.Where(cond).Find(&result)
	if tx.RowsAffected > 0 && result[0].Id == 0 {
		return nil, nil
	}
	return result, tx.Error
}
func (dao *LOLMEquipmentDAO) Add(equips []*model.LOLMEquipment) (int64, error) {
	result := dao.db.Create(equips)
	return result.RowsAffected, result.Error
}

func (dao *LOLMEquipmentDAO) Update(data *model.LOLMEquipment, cond map[string]interface{}) (int64, error) {
	result := dao.db.Model(model.LOLMEquipment{}).Where(cond).Updates(data)
	return result.RowsAffected, result.Error
}

func (dao *LOLMEquipmentDAO) UpdatesInto(fileTime, version string, ups map[string][]string) (int64, error) {
	sql := `
UPDATE
	lolm_equipment
SET
	%s
	ELSE
		''
	END
WHERE %s
`

	caseWhen := "`into` = CASE"
	for equipID, ins := range ups {
		caseWhen += fmt.Sprintf(" WHEN equipId = %s THEN '%s' ", equipID, strings.Join(ins, ","))
	}

	where := "status = 0 and fileTime = '%s' and version = '%s'"
	where = fmt.Sprintf(where, fileTime, version)

	sql = fmt.Sprintf(sql, caseWhen, where)
	tx := dao.db.Exec(sql)
	return tx.RowsAffected, tx.Error
}

func (dao *LOLMEquipmentDAO) GetLOLMEquipmentMaxVersion() (*model.LOLMEquipment, error) {
	tx := dao.db.Model(&model.LOLMEquipment{})
	var result model.LOLMEquipment
	tx = tx.Where("status = 0").Order("version desc").First(&result)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, tx.Error
}
func (dao *LOLMEquipmentDAO) GetLOLMEquipment(version string) ([]*model.LOLMEquipment, error) {
	// 查当前版本所有数据
	cond := map[string]interface{}{
		"version": version,
		"status":  0,
	}
	not := map[string]interface{}{
		"equipId": config.EquipDict.Exclude,
	}

	var equip []*model.LOLMEquipment
	tx := dao.db.Model(model.LOLMEquipment{}).Not(not).Find(&equip, cond)
	//equip, err := dao.db.Model(model.LOLMEquipment{}).Not(not).Find(nil, cond)
	//if tx.Error != nil {
	//	return nil, tx.Error
	//}
	return equip, tx.Error
}

func (dao *LOLMEquipmentDAO) GetRoadmap(version string, id string) (map[string][]*model.LOLMEquipment, error) {
	result := make(map[string][]*model.LOLMEquipment)

	current, err := dao.Find(nil, map[string]interface{}{
		"status":  0,
		"version": version,
		"equipId": id,
	})
	if err != nil {
		return nil, err
	}
	if len(current) == 0 {
		return nil, errors.New("can not find id " + id)
	}
	result["current"] = current

	from, err := dao.Find(nil, map[string]interface{}{
		"status":  0,
		"version": version,
		"equipId": strings.Split(current[0].From, ","),
	})
	if err != nil {
		return nil, err
	}
	result["from"] = from

	into, err := dao.Find(nil, map[string]interface{}{
		"status":  0,
		"version": version,
		"equipId": strings.Split(current[0].Into, ","),
	})
	if err != nil {
		return nil, err
	}
	result["into"] = into

	return result, nil
}

func (dao *LOLMEquipmentDAO) GetLOLMEquipmentWithExt(version string) ([]*model.LOLMEquipment, error) {
	notin := strings.Join(config.EquipDict.Exclude, ",")

	// 查当前版本所有数据
	result := make([]*model.LOLMEquipment, 0)
	sql := `
SELECT
	equip.equipId,
	equip.name,
	equip.iconPath,
	equip.price,
	equip.description,
	CONCAT_WS(',', equip.searchKey, alias.keywords, alias.keywords_py) AS keywords,
	equip.from,
	equip.into,
	equip.type,
	equip.version,
	equip.fileTime
FROM
	lolm_equipment equip
	LEFT JOIN equip_alias alias ON equip.name = alias.name
WHERE
	equip.version = '%s' and equip.status = 0 and alias.platform = 0
	and equip.equipId not in (%s)
`
	sql = fmt.Sprintf(sql, version, notin)
	err := dao.db.Exec(sql).Find(&result).Error
	return result, err
}

var (
	lolmEDao  *LOLMEquipmentDAO
	lolmEOnce *sync.Once
)

func NewLOLMEquipmentDAO() *LOLMEquipmentDAO {
	lolmEOnce.Do(func() {
		lolmEDao = &LOLMEquipmentDAO{
			db: mysql.DB,
		}

	})
	return lolmEDao
}
