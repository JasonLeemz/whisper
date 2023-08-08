package dao

import (
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"whisper/internal/model"
	"whisper/pkg/context"
	"whisper/pkg/es"
)

type ESSkill interface {
	CreateIndex(ctx *context.Context) error
	Heroes2ES(ctx *context.Context, data []*model.ESSkill) error
}

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

func (dao *ESSkillDAO) Skill2ES(ctx *context.Context, data []*model.ESSkill) error {
	var esModel model.ESSkill
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

func NewESSkillDAO() *ESSkillDAO {
	return &ESSkillDAO{
		esClient: es.ESClient,
	}
}
