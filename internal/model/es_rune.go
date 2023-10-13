package model

import (
	"sync"
	. "whisper/internal/model/common"
)

type ESRune struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Keywords    string `json:"keywords,omitempty"`
	Description string `json:"description,omitempty"`
	Plaintext   string `json:"plaintext,omitempty"`
	IconPath    string `json:"iconPath,omitempty"`
	Tooltip     string `json:"tooltip,omitempty"`
	SlotLabel   string `json:"slotLabel,omitempty"`
	StyleName   string `json:"styleName,omitempty"`
	Maps        string `json:"maps,omitempty"`
	Types       string `json:"types,omitempty"`
	Version     string `json:"version,omitempty"`
	FileTime    string `json:"fileTime,omitempty"`
	Platform    string `json:"platform"`
}

var (
	modelRune *ESRune
	onceRune  sync.Once
)

func NewModelESRune() *ESRune {
	onceRune.Do(func() {
		modelRune = new(ESRune)
	})

	return modelRune
}

func (e *ESRune) GetMapping() string {
	return `
{
    "settings": {
        "analysis": {
            "analyzer": {
                "my_analyzer": {
                    "tokenizer": "ik_max_word",
                    "filter": "py"
                }
            },
            "filter": {
                "py": {
                    "type": "pinyin",
                    "keep_full_pinyin": false,
                    "keep_joined_full_pinyin": true,
                    "keep_original": true,
                    "limit_first_letter_length": 16,
                    "remove_duplicated_term": true,
                    "none_chinese_pinyin_tokenize": false
                }
            }
        }
    },
    "mappings": {
        "properties": {
            "name": {
                "type": "text",
                "analyzer": "my_analyzer",
                "search_analyzer": "ik_smart"
            },
            "keywords": {
                "type": "text",
                "analyzer": "my_analyzer",
                "search_analyzer": "ik_smart"
            },
            "plaintext": {
                "type": "text",
                "analyzer": "ik_smart"
            },
            "description": {
                "type": "text",
                "analyzer": "ik_smart"
            },
            "tooltip": {
                "type": "text",
                "analyzer": "ik_smart"
            },
            "styleName": {
                "type": "text",
                "analyzer": "ik_smart"
            },
            "types": {
                "type": "text",
                "analyzer": "ik_smart"
            },
            "slotLabel": {
                "type": "keyword"
            },
            "iconPath": {
                "type": "keyword"
            },
            "maps": {
                "type": "keyword"
            },
            "fileTime": {
                "type": "keyword"
            },
            "platform": {
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

func (e *ESRune) GetIndexName() string {
	return IndexRune
}
