package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"sync"
	"whisper/internal/dto"
	"whisper/internal/model"
	"whisper/internal/model/common"
	"whisper/pkg/context"
	"whisper/pkg/es"
)

type ESHeroes interface {
	CreateIndex(ctx *context.Context) error
	DeleteIndex(ctx *context.Context) error
	Heroes2ES(ctx *context.Context, data []*model.ESHeroes) error
	Find(ctx *context.Context, cond *common.QueryCond) ([]*model.ESHeroes, error)
}

type ESHeroesDAO struct {
	esClient *elastic.Client
}

func (dao *ESHeroesDAO) CreateIndex(ctx *context.Context) error {

	// 索引是否存在
	var esModel model.ESHeroes
	idxName := esModel.GetIndexName()

	exists, err := es.ESClient.IndexExists(idxName).Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		// 创建索引
		createIndex, err := es.ESClient.CreateIndex(idxName).Body(esModel.GetMapping()).Do(ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return errors.New(fmt.Sprintf("expected IndicesCreateResult.Acknowledged %v; got %v", true, createIndex.Acknowledged))
		}
	}
	return nil
}

func (dao *ESHeroesDAO) DeleteIndex(ctx *context.Context) error {

	// 索引是否存在
	var esModel model.ESHeroes
	idxName := esModel.GetIndexName()

	exists, err := es.ESClient.IndexExists(idxName).Do(ctx)
	if err != nil {
		return err
	}
	if exists {
		// 删除索引
		deleteIndex, err := es.ESClient.DeleteIndex(idxName).Do(ctx)
		if err != nil {
			return err
		}
		if !deleteIndex.Acknowledged {
			return errors.New(fmt.Sprintf("expected IndicesDeleteResult.Acknowledged %v; got %v", true, deleteIndex.Acknowledged))
		}
	}
	return nil
}

func (dao *ESHeroesDAO) Heroes2ES(ctx *context.Context, data []*model.ESHeroes) error {
	var esModel model.ESHeroes
	idxName := esModel.GetIndexName()
	// 导入数据
	for _, e := range data {
		_, err := es.ESClient.Index().Index(idxName).BodyJson(e).Id(e.ID).Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dao *ESHeroesDAO) Find(ctx *context.Context, cond *common.QueryCond) ([]*model.ESHeroes, error) {
	var esModel model.ESHeroes
	idxName := esModel.GetIndexName()
	query := elastic.NewBoolQuery()
	if cond.MultiMatchQuery != nil {
		query = query.Must(elastic.NewMultiMatchQuery(cond.MultiMatchQuery.Text, cond.MultiMatchQuery.Fields...))
	}

	if cond.TermsQuery != nil {
		query = query.Must(elastic.NewTermsQuery(cond.TermsQuery.Name, cond.TermsQuery.Values...))
	}

	if cond.TermQuery != nil {
		for _, c := range cond.TermQuery {
			query = query.Must(elastic.NewTermQuery(c.Name, c.Value))
		}
	}

	res, err := es.ESClient.Search().
		Index(idxName).
		Query(query).
		From(0).Size(10000).
		Pretty(true).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	resp := dto.EsResultHits{}
	data, _ := json.Marshal(res.Hits)
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	var heroes []*model.ESHeroes
	for _, hit := range resp.Hits {
		sourceStr, _ := json.Marshal(hit.TmpSource)
		hitData := &model.ESHeroes{}
		err = json.Unmarshal(sourceStr, hitData)
		if err != nil {
			return nil, err
		}
		heroes = append(heroes, hitData)
	}
	return heroes, nil
}

var (
	esHDao  *ESHeroesDAO
	esHOnce sync.Once
)

func NewESHeroesDAO() *ESHeroesDAO {

	esHOnce.Do(func() {
		esHDao = &ESHeroesDAO{
			esClient: es.ESClient,
		}
	})
	return esHDao
}
