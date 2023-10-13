package model

import (
	"sync"
	"whisper/internal/model/common"
)

type ESHeroes struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	IconPath    string `json:"iconPath,omitempty"`
	MainImg     string `json:"mainImg,omitempty"`
	Roles       string `json:"roles,omitempty"`
	Price       string `json:"price,omitempty"`
	Description string `json:"description,omitempty"`
	Plaintext   string `json:"plaintext,omitempty"`
	Spells      string `json:"spells,omitempty"`
	Keywords    string `json:"keywords,omitempty"`
	Maps        string `json:"maps,omitempty"`
	Types       string `json:"types,omitempty"`
	Version     string `json:"version,omitempty"`
	FileTime    string `json:"fileTime,omitempty"`
	Platform    string `json:"platform"`
}

var (
	modelHeroes *ESHeroes
	onceHeroes  sync.Once
)

func NewModelESHeroes() *ESHeroes {
	onceHeroes.Do(func() {
		modelHeroes = new(ESHeroes)
	})

	return modelHeroes
}

func (e *ESHeroes) GetMapping() string {
	return `
{
    "mappings": {
        "properties": {
            "name": {
                "type": "text",
                "analyzer": "ik_smart"
            },
            "keywords": {
                "type": "text",
                "analyzer": "ik_smart"
            },
            "description": {
                "type": "text",
                "analyzer": "ik_smart"
            },
			"spells": {
                "type": "text",
                "analyzer": "ik_smart"
            },
			"roles": {
                "type": "keyword"
            },
            "fileTime": {
                "type": "keyword"
            },
            "iconPath": {
                "type": "keyword"
            },
			"mainImg": {
                "type": "keyword"
            },
            "itemId": {
                "type": "keyword"
            },
            "maps": {
                "type": "keyword"
            },
            "plaintext": {
                "type": "text",
                "analyzer": "ik_smart"
            },
            "platform": {
                "type": "keyword"
            },
            "price": {
                "type": "keyword"
            },
            "sell": {
                "type": "keyword"
            },
            "total": {
                "type": "keyword"
            },
            "version": {
                "type": "keyword"
            }
        }
    }
}
`
}

func (e *ESHeroes) GetIndexName() string {
	return common.IndexHeroes
}
