package model

type ESHeroes struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	IconPath    string `json:"iconPath,omitempty"`
	Price       string `json:"price,omitempty"`
	Description string `json:"description,omitempty"`
	Plaintext   string `json:"plaintext,omitempty"`
	Keywords    string `json:"keywords,omitempty"`
	Maps        string `json:"maps,omitempty"`
	Types       string `json:"types,omitempty"`
	Version     string `json:"version,omitempty"`
	FileTime    string `json:"fileTime,omitempty"`
	Platform    string `json:"platform"`
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
            "fileTime": {
                "type": "keyword"
            },
            "iconPath": {
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
	return "lol_heroes"
}
