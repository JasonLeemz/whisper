package es

import (
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	"strings"
	"whisper/internal/dto"
	"whisper/pkg/context"
	"whisper/pkg/es"
)

type Instance struct {
	indexName       string
	multiMatchQuery *MultiMatchQuery
	termsQuery      *TermsQuery
	termQuery       []*TermQuery
	fieldSort       *FieldSort
}

type MultiMatchQuery struct {
	Text   interface{}
	Fields []string
}

type TermsQuery struct {
	Name   string
	Values []interface{}
}

type TermQuery struct {
	Name  string
	Value interface{}
}

type FieldSort struct {
	Field     string
	Direction string
}

// InsBuilder 建造者模式构建查询条件
type InsBuilder struct {
	indexName       string
	multiMatchQuery *MultiMatchQuery
	termsQuery      *TermsQuery
	termQuery       []*TermQuery
	fieldSort       *FieldSort
}

func (i *InsBuilder) SetMultiMatchQuery(multiMatchQuery *MultiMatchQuery) *InsBuilder {
	i.multiMatchQuery = multiMatchQuery
	return i
}

func (i *InsBuilder) SetTermsQuery(termsQuery *TermsQuery) *InsBuilder {
	i.termsQuery = termsQuery
	return i
}

func (i *InsBuilder) SetTermQuery(termQuery []*TermQuery) *InsBuilder {
	i.termQuery = termQuery
	return i
}

func (i *InsBuilder) SetFieldSort(fieldSort *FieldSort) *InsBuilder {
	// 统一转换为小写
	fieldSort.Direction = strings.ToLower(fieldSort.Direction)
	i.fieldSort = fieldSort
	return i
}

func (i *InsBuilder) Build() (*Instance, error) {
	ins := &Instance{indexName: i.indexName}

	if i.multiMatchQuery != nil {
		if len(i.multiMatchQuery.Fields) == 0 {
			return nil, errors.New("multiMatchQuery.Fields[] must have at least 1 item")
		}
		ins.multiMatchQuery = i.multiMatchQuery
	}

	if i.termsQuery != nil {
		if len(i.termsQuery.Values) == 0 {
			return nil, errors.New("termsQuery.Values[] must have at least 1 item")
		}
		ins.termsQuery = i.termsQuery
	}

	if i.termQuery != nil {
		ins.termQuery = i.termQuery
	}

	if i.fieldSort == nil {
		// 默认按照_score降序
		ins.fieldSort = &FieldSort{
			Field:     "_score",
			Direction: "desc",
		}
	} else {
		if i.fieldSort.Field == "" {
			return nil, errors.New("fieldSort.Field must be a useful value")
		}
		if i.fieldSort.Direction != "asc" && i.fieldSort.Direction != "desc" {
			return nil, errors.New("fieldSort.Direction must be one of 'asc' or 'desc'")
		}
		ins.fieldSort = i.fieldSort
	}

	return ins, nil
}

func NewInsBuilder(indexName string) *InsBuilder {
	return &InsBuilder{indexName: indexName}
}

func (i *Instance) Builder(indexName string) *InsBuilder {
	// 这里可以更换为任意Builder
	return NewInsBuilder(indexName)
}

func (i *Instance) Query(ctx *context.Context) (*dto.EsResultHits, error) {
	query := elastic.NewBoolQuery()

	if i.multiMatchQuery != nil {
		query = query.Must(elastic.NewMultiMatchQuery(i.multiMatchQuery.Text, i.multiMatchQuery.Fields...))
	}

	if i.termsQuery != nil {
		query = query.Must(elastic.NewTermsQuery(i.termsQuery.Name, i.termsQuery.Values...))
	}

	if i.termQuery != nil {
		for _, c := range i.termQuery {
			query = query.Must(elastic.NewTermQuery(c.Name, c.Value))
		}
	}

	sortByScore := &elastic.FieldSort{}
	if i.fieldSort != nil {
		if i.fieldSort.Direction == "asc" {
			sortByScore = elastic.NewFieldSort(i.fieldSort.Field).Asc()
		} else {
			sortByScore = elastic.NewFieldSort(i.fieldSort.Field).Desc()
		}
	}

	res, err := es.ESClient.Search().
		Index(i.indexName).
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
