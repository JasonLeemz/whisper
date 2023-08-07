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
	KeyWords string   `json:"key_words"`
	Platform string   `json:"platform,omitempty"`
	Category []string `json:"category,omitempty"`
	Map      []string `json:"map,omitempty"`
}

func EsSearch(ctx *context.Context, p *SearchParams) (*common.EsEquipHits, error) {

	query := elastic.NewBoolQuery()

	// 按名字介绍
	cate := make([]string, 0)
	for _, m := range p.Category {
		cate = append(cate, m)
	}
	query = query.Filter(elastic.NewMultiMatchQuery(p.KeyWords, cate...))

	// 按地图
	maps := make([]interface{}, 0)
	for _, m := range p.Map {
		maps = append(maps, m)
	}
	query = query.Filter(elastic.NewTermsQuery("maps", maps...))

	// 端游or手游
	query = query.Filter(elastic.NewTermQuery("platform", p.Platform))

	//query = query.Filter(elastic.NewRangeQuery("id").Gte(0))
	//query = query.Filter(elastic.NewRangeQuery("id").Lte(9999999))

	var esEquip model.ESEquipment
	res, err := es.ESClient.Search().Index(esEquip.GetIndexName()).Query(query).From(0).Size(10).Do(ctx)
	if err != nil {
		return nil, err
	}
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
					Maps:     "召唤师峡谷",
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
