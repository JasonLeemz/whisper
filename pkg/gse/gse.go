package gse

import (
	"github.com/go-ego/gse"
	"github.com/go-ego/gse/hmm/idf"
)

func Analyzer(text string) idf.Segments {
	x, _ := gse.New()
	var te idf.TagExtracter
	te.WithGse(x)
	_ = te.LoadIdf()
	defer x.Empty()

	return te.ExtractTags(text, 10)
}
