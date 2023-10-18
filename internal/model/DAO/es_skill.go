package dao

import (
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"sync"
	"whisper/internal/model"
	"whisper/pkg/context"
	"whisper/pkg/es"
)

type ESSkillDAO struct {
	esClient *elastic.Client
}

func (dao *ESSkillDAO) CreateIndex(ctx *context.Context) error {
	// 索引是否存在
	var esModel model.ESSkill
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

func (dao *ESSkillDAO) DeleteIndex(ctx *context.Context) error {
	// 索引是否存在
	var esModel model.ESSkill
	idxName := esModel.GetIndexName()

	exists, err := es.ESClient.IndexExists(idxName).Do(ctx)
	if err != nil {
		return err
	}
	if exists {
		// 创建索引
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

func (dao *ESSkillDAO) Data2ES(ctx *context.Context, data interface{}) error {
	var esModel model.ESSkill
	idxName := esModel.GetIndexName()
	// 导入数据
	for _, e := range data.([]*model.ESSkill) {
		_, err := es.ESClient.Index().Index(idxName).BodyJson(e).Id(e.ID).Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

var (
	esKDao     *ESSkillDAO
	onceEsKDao sync.Once
)

func NewESSkillDAO() EsDaoFunc {
	return func() ESIndex {
		onceEsKDao.Do(func() {
			esKDao = &ESSkillDAO{
				esClient: es.ESClient,
			}
		})
		return esKDao
	}
}
