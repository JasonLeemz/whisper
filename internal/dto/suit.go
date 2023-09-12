package dto

type HeroSuit struct {
	HeroID   string                        `json:"hero_id"`
	Platform int                           `json:"platform"`
	Equips   map[string]RecommendSuitEquip `json:"equips"` // map["top|bottom|support..."]RecommendSuitEquip
	ExtInfo  HeroSuitExtInfo               `json:"ext_info"`
}

type HeroSuitExtInfo struct {
	RecommendReason map[string]string `json:"recommend_reason"`
	AuthorInfo      map[string]struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	} `json:"author_info"`
}

type RecommendSuitEquip struct {
	Out   [][]*SuitData `json:"out"`
	Shoe  [][]*SuitData `json:"shoe"`
	Core  [][]*SuitData `json:"core"`
	Other [][]*SuitData `json:"other"`

	Rune  [][]*SuitData `json:"rune,omitempty"`
	Skill [][]*SuitData `json:"skill,omitempty"`
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

	CD int `json:"cd"`

	Igamecnt int32 `json:"igamecnt,omitempty"`
	Wincnt   int32 `json:"wincnt,omitempty"`
	Winrate  int32 `json:"winrate,omitempty"`
	Allcnt   int32 `json:"allcnt,omitempty"`
	Showrate int32 `json:"showrate,omitempty"`

	Title        string `json:"title"`
	Author       string `json:"author"`
	AuthorIcon   string `json:"author_icon"`
	RecommendID  string `json:"recommend_id"`
	ThinkingInfo string `json:"thinking_info"`

	Platform int `json:"platform"`
}
