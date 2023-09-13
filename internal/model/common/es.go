package common

type QueryCond struct {
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
