package dto

type HeroSuit struct {
	HeroID string                        `json:"hero_id"`
	Equips map[string]RecommendSuitEquip `json:"equips"` // map["top|bottom|support..."]RecommendSuitEquip
}

type RecommendSuitEquip struct {
	Out   [][]*SuitData `json:"out"`
	Shoe  [][]*SuitData `json:"shoe"`
	Core  [][]*SuitData `json:"core"`
	Other [][]*SuitData `json:"other"`
}

type SuitData struct {
	ID        int    `json:"ID"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Maps      string `json:"maps"`
	Level     string `json:"level"`
	Plaintext string `json:"plaintext"`
	Desc      string `json:"desc"`
	Price     int    `json:"price"`
	Sell      int    `json:"sell"`
	Version   string `json:"version"`

	Igamecnt int32 `json:"igamecnt,omitempty"`
	Wincnt   int32 `json:"wincnt,omitempty"`
	Winrate  int32 `json:"winrate,omitempty"`
	Allcnt   int32 `json:"allcnt,omitempty"`
	Showrate int32 `json:"showrate,omitempty"`

	Title       string `json:"title"`
	Author      string `json:"author"`
	AuthorIcon  string `json:"author_icon"`
	RecommendID string `json:"recommend_id"`
}
