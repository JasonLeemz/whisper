package dao

import (
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"whisper/internal/model"
	"whisper/pkg/context"
	"whisper/pkg/es"
)

type ESHeroes interface {
	CreateIndex(ctx *context.Context) error
	DeleteIndex(ctx *context.Context) error
	Heroes2ES(ctx *context.Context, data []*model.ESHeroes) error
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

func NewESHeroesDAO() *ESHeroesDAO {
	return &ESHeroesDAO{
		esClient: es.ESClient,
	}
}
