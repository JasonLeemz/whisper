package dao

import (
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"whisper/internal/dto"
	. "whisper/internal/model/common"
	"whisper/pkg/context"
	"whisper/pkg/es"
)

const (
	ESIndexEquipment = IndexEquipment
	ESIndexHeroes    = IndexHeroes
	ESIndexRune      = IndexRune
	ESIndexSkill     = IndexSkill
)

type EsDaoFunc func() ESIndex

func CreateEsDao(t string) EsDaoFunc {
	switch t {
	case ESIndexEquipment:
		return NewESEquipmentDAO()
	case ESIndexHeroes:
		return NewESHeroesDAO()
	case ESIndexRune:
		return NewESRuneDAO()
	case ESIndexSkill:
		return NewESSkillDAO()

	}
	return nil
}

type ESIndex interface {
	CreateIndex(ctx *context.Context) error
	DeleteIndex(ctx *context.Context) error
	Data2ES(ctx *context.Context, data interface{}) error
	Find(ctx *context.Context, cond *QueryCond) ([]map[string]interface{}, error)
}

func ESQuery(ctx *context.Context, idxName string, cond *QueryCond) (*dto.EsResultHits, error) {
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

	sortByScore := elastic.NewFieldSort("_score").Desc()
	if cond.FieldSort != nil {
		if cond.FieldSort.Direction == "asc" {
			sortByScore = elastic.NewFieldSort(cond.FieldSort.Field).Asc()
		} else {
			sortByScore = elastic.NewFieldSort(cond.FieldSort.Field).Desc()
		}
	}

	res, err := es.ESClient.Search().
		Index(idxName).
		Query(query).
		SortBy(sortByScore).
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
	return &resp, nil
}
