package common

import (
	"errors"
	"strings"
)

type QueryCond struct {
	indexName       string
	MultiMatchQuery *MultiMatchQuery
	TermsQuery      *TermsQuery
	TermQuery       []*TermQuery
	FieldSort       *FieldSort
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

// CondBuilder 建造者模式构建查询条件
type CondBuilder struct {
	indexName       string
	multiMatchQuery *MultiMatchQuery
	termsQuery      *TermsQuery
	termQuery       []*TermQuery
	fieldSort       *FieldSort
}

func (cb *CondBuilder) SetMultiMatchQuery(multiMatchQuery *MultiMatchQuery) {
	cb.multiMatchQuery = multiMatchQuery
}

func (cb *CondBuilder) SetTermsQuery(termsQuery *TermsQuery) {
	cb.termsQuery = termsQuery
}

func (cb *CondBuilder) SetTermQuery(termQuery []*TermQuery) {
	cb.termQuery = termQuery
}

func (cb *CondBuilder) SetFieldSort(fieldSort *FieldSort) {
	// 统一转换为小写
	fieldSort.Direction = strings.ToLower(fieldSort.Direction)
	cb.fieldSort = fieldSort
}

func NewCondBuilder(indexName string) *CondBuilder {
	return &CondBuilder{indexName: indexName}
}
func (c *QueryCond) Builder(indexName string) *CondBuilder {
	// 这里可以更换为任意Builder
	return NewCondBuilder(indexName)
}
func (cb *CondBuilder) Build() (*QueryCond, error) {
	cond := &QueryCond{}

	if cb.multiMatchQuery != nil {
		if len(cb.multiMatchQuery.Fields) == 0 {
			return nil, errors.New("multiMatchQuery.Fields[] must have at least 1 item")
		}
		cond.MultiMatchQuery = cb.multiMatchQuery
	}

	if cb.termsQuery != nil {
		if len(cb.termsQuery.Values) == 0 {
			return nil, errors.New("termsQuery.Values[] must have at least 1 item")
		}
		cond.TermsQuery = cb.termsQuery
	}

	if cb.termQuery != nil {
		cond.TermQuery = cb.termQuery
	}

	if cb.fieldSort == nil {
		// 默认按照_score降序
		cond.FieldSort = &FieldSort{
			Field:     "_score",
			Direction: "desc",
		}
	} else {
		if cb.fieldSort.Field == "" {
			return nil, errors.New("fieldSort.Field must be a useful value")
		}
		if cb.fieldSort.Direction != "asc" && cb.fieldSort.Direction != "desc" {
			return nil, errors.New("fieldSort.Direction must be one of 'asc' or 'desc'")
		}
		cond.FieldSort = cb.fieldSort
	}

	return cond, nil
}
