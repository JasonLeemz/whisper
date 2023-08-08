package common

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
	Score     float64     `json:"_score,omitempty"`
	Index     string      `json:"_index,omitempty"`
	Type      string      `json:"_type,omitempty"`
	Id        string      `json:"_id,omitempty"`
	TmpSource interface{} `json:"_source,omitempty"`
	Source    Source      `json:"source,omitempty"`
}

type Source struct {
	Name        string   `json:"name"`
	IconPath    string   `json:"iconPath"`
	Description string   `json:"description"`
	Plaintext   string   `json:"plaintext"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
}
