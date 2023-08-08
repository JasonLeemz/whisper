package model

type ESSkill struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Keywords    string `json:"keywords,omitempty"`
	Description string `json:"description,omitempty"`
	Plaintext   string `json:"plaintext,omitempty"`
	IconPath    string `json:"iconPath,omitempty"`
	Maps        string `json:"maps,omitempty"`
	CoolDown    string `json:"cooldown,omitempty"`
	Version     string `json:"version,omitempty"`
	FileTime    string `json:"fileTime,omitempty"`
	Platform    string `json:"platform"`
}

func (e *ESSkill) GetMapping() string {
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
            "plaintext": {
                "type": "text",
                "analyzer": "ik_smart"
            },
            "description": {
                "type": "text",
                "analyzer": "ik_smart"
            },
            "iconPath": {
                "type": "keyword"
            },
            "maps": {
                "type": "keyword"
            },
            "cooldown": {
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

func (e *ESSkill) GetIndexName() string {
	return "lol_skill"
}
