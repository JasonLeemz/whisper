package dto

// --------------------------------------------------------

type SearchResult struct {
	Tips string              `json:"tips"`
	List []*SearchResultList `json:"list"`
}

type SearchResultList struct {
	Tags      []string `json:"tags"`
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Icon      string   `json:"icon"`
	Desc      string   `json:"desc"`
	Plaintext string   `json:"plaintext"`
	Price     int      `json:"price"`
	Maps      string   `json:"maps"`
	Platform  int      `json:"platform"`
	Version   string   `json:"version"`
	Keywords  []string `json:"keywords"`

	Spell []*HeroSpell `json:"spell"`
}

type HeroSpell struct {
	Icon       string `json:"icon"`
	Name       string `json:"name"`
	Sort       string `json:"sort"`
	Desc       string `json:"desc"`
	LevelSpell string `json:"levelSpell"`
}
