package logic

import (
	"encoding/json"
	"strconv"
	"whisper/internal/logic/common"
	"whisper/internal/model"
	"whisper/internal/model/DAO"
	"whisper/pkg/es"

	"github.com/olivere/elastic/v7"
	"whisper/pkg/config"
	"whisper/pkg/context"
	"whisper/pkg/log"
)

type SearchParams struct {
	KeyWords string `json:"key_words"`
	Platform int    `json:"platform,omitempty"`
	Category int    `json:"category,omitempty"`
	Map      string `json:"map,omitempty"`
}

func EsSearch(ctx *context.Context, p *SearchParams) (*common.EsEquipHits, error) {

	query := elastic.NewBoolQuery()
	if p.Category == 0 {
		// 按名字搜索
		query = query.Must(elastic.NewMultiMatchQuery(p.KeyWords, "name", "keywords"))
	} else if p.Category == 1 {
		// 按功能介绍搜索
		query = query.Must(elastic.NewMultiMatchQuery(p.KeyWords, "description", "plaintext"))
	} else {
		// 未指定搜索范围，全字段搜索
		query = query.Must(elastic.NewMultiMatchQuery(p.KeyWords, "name", "description", "plaintext", "keywords"))

	}

	if p.Map != "" {
		query = query.Filter(elastic.NewTermQuery(p.Map, "maps"))
	} else {
		cond := make([]interface{}, 0)
		for _, m := range config.GlobalConfig.Search.MapsLOL {
			cond = append(cond, m)
		}
		query = query.Filter(elastic.NewTermsQuery("maps", cond...))
	}

	if p.Platform == 1 {
		query = query.Filter(elastic.NewTermQuery("1", "platform"))
	} else {
		query = query.Filter(elastic.NewTermQuery(strconv.Itoa(p.Platform), "platform"))
	}

	//query = query.Filter(elastic.NewTermQuery("version", "13.14")) // Filter 不会算分，Must会算分

	//query = query.Filter(elastic.NewTermsQuery("name",[]string{
	//	"xz","xie",
	//}))

	//query = query.Filter(elastic.NewRangeQuery("id").Gte(0))
	//query = query.Filter(elastic.NewRangeQuery("id").Lte(9999999))

	//es.ESClient.Search().Index("lol_equipment_13.14").Query(query).From(int(10)).Size(int(10)).Do(ctx)

	var esEquip model.ESEquipment
	res, err := es.ESClient.Search().Index(esEquip.GetIndexName()).Query(query).From(int(0)).Size(int(10)).Do(ctx)
	if err != nil {
		return nil, err
	}
	//res.Hits.TotalHits.Value
	resp := common.EsEquipHits{}
	data, _ := json.Marshal(res.Hits)
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// BuildIndex 重建索引
func BuildIndex(ctx *context.Context, index string) error {
	queue := make([]string, 0)
	// 如果没有指定index，就重建所有
	if index == "" {
		queue = config.GlobalConfig.ES.BuildIndex
	} else {
		queue = append(queue, index)
	}

	// 表数据直接建索引
	for _, tbl := range queue {
		err := mysql2es(ctx, tbl)
		if err != nil {
			log.Logger.Error(ctx, err)
			return err
		}
	}
	return nil
}

var sourceChan = make(chan []*model.ESEquipment, 100) // 最多起100个协程处理索引

func mysql2es(ctx *context.Context, tblName string) error {
	// 创建索引
	ed := dao.NewESEquipmentDAO()
	err := ed.CreateIndex(ctx)
	if err != nil {
		return err
	}
	go customer(ctx)

	var equipModel model.LOLEquipment
	var m_equipModel model.LOLMEquipment

	switch tblName {
	case m_equipModel.TableName():
		d := dao.NewLOLMEquipmentDAO()
		version, err := d.GetLOLMEquipmentMaxVersion()
		if err != nil {
			return err
		}

		equipment, err := d.GetLOLMEquipment(version)
		if err != nil {
			return err
		}

		go func() {
			for _, equip := range equipment {
				var esEquip []*model.ESEquipment
				tmp := equip
				esEquip = append(esEquip, &model.ESEquipment{
					ID:          tmp.EquipId + "_lolm",
					EquipId:     tmp.EquipId,
					Name:        tmp.Name,
					IconPath:    tmp.IconPath,
					Price:       tmp.Price,
					Description: tmp.Description,
					//Plaintext:    tmp.,
					//Sell:         tmp.,
					Total: tmp.Price,
					//SuitHeroId:   tmp.SuitHeroId,
					//SuitHeroName: tmp.SuitHeroId, // todo
					//SuitHeroIcon: tmp.SuitHeroId, // todo
					Keywords: tmp.SearchKey, // 数据为空
					Maps:     "LOLM",
					From:     tmp.From, // todo
					Into:     tmp.Into, // todo
					Types:    tmp.Type,
					Version:  tmp.Version,
					FileTime: tmp.FileTime,
					Platform: strconv.Itoa(common.PlatformForLOLM),
				})

				// 放入阻塞队列
				sourceChan <- esEquip
			}
		}()
	case equipModel.TableName():
		d := dao.NewLOLEquipmentDAO()
		version, err := d.GetLOLEquipmentMaxVersion()
		if err != nil {
			return err
		}

		equipment, err := d.GetLOLEquipment(version)
		if err != nil {
			return err
		}

		go func() {
			for _, equip := range equipment {
				var esEquip []*model.ESEquipment
				tmp := equip
				esEquip = append(esEquip, &model.ESEquipment{
					ID:           tmp.ItemId + "_" + tmp.Maps,
					EquipId:      tmp.ItemId,
					Name:         tmp.Name,
					IconPath:     tmp.IconPath,
					Price:        tmp.Price,
					Description:  tmp.Description,
					Plaintext:    tmp.Plaintext,
					Sell:         tmp.Sell,
					Total:        tmp.Total,
					SuitHeroId:   tmp.SuitHeroId,
					SuitHeroName: tmp.SuitHeroId, // todo
					SuitHeroIcon: tmp.SuitHeroId, // todo
					Keywords:     tmp.Keywords,
					Maps:         tmp.Maps,
					From:         tmp.From,  // todo
					Into:         tmp.Into,  // todo
					Types:        tmp.Types, // todo
					Version:      tmp.Version,
					FileTime:     tmp.FileTime,
					Platform:     strconv.Itoa(common.PlatformForLOL),
				})

				// 放入阻塞队列
				sourceChan <- esEquip

			}
		}()
	}
	return nil
}

func customer(ctx *context.Context) {
	var equips []*model.ESEquipment
	ed := dao.NewESEquipmentDAO()
	for {
		equips = <-sourceChan
		// LOLEquipment2ES 支持批量索引，这里避免占用内存过大每次只处理一行数据
		err := ed.Equipment2ES(ctx, equips)
		if err != nil {
			log.Logger.Error(ctx, err)
		}

		// TODO
	}
}
