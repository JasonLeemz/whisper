package dto

type EsResultHits struct {
	Total    Total   `json:"total,omitempty"`
	MaxScore float64 `json:"max_score,omitempty"`
	Hits     []Hits  `json:"hits,omitempty"`
}

type Total struct {
	Value    int    `json:"value,omitempty"`
	Relation string `json:"relation,omitempty"`
}

type Hits struct {
	Score     float64             `json:"_score,omitempty"`
	Index     string              `json:"_index,omitempty"`
	Type      string              `json:"_type,omitempty"`
	Id        string              `json:"_id,omitempty"`
	TmpSource interface{}         `json:"_source,omitempty"`
	Source    Source              `json:"source,omitempty"` // 冗余的结构，用于存放_source字符串解析出的内容
	Highlight map[string][]string `json:"highlight,omitempty"`
}

type Source struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	IconPath    string   `json:"iconPath"`
	MainImg     string   `json:"mainImg"`
	ItemId      string   `json:"itemId"`
	Description string   `json:"description"`
	Plaintext   string   `json:"plaintext"`
	Version     string   `json:"version"`
	Platform    string   `json:"platform"`
	Maps        string   `json:"maps"`
	Tags        []string `json:"tags"`
}
