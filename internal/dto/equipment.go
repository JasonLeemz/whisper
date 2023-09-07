package dto

type LOLEquipment struct {
	Tree     []Tree  `json:"tree,omitempty"`
	Items    []Items `json:"items,omitempty"`
	Version  string  `json:"version,omitempty"`
	FileName string  `json:"fileName,omitempty"`
	FileTime string  `json:"fileTime,omitempty"`
}

type Tree struct {
	Header string   `json:"header,omitempty"`
	Tags   []string `json:"tags,omitempty"`
}

type Items struct {
	ItemId      string   `json:"itemId,omitempty"`
	Name        string   `json:"name,omitempty"`
	IconPath    string   `json:"iconPath,omitempty"`
	Price       string   `json:"price,omitempty"`
	Description string   `json:"description,omitempty"`
	Plaintext   string   `json:"plaintext,omitempty"`
	Sell        string   `json:"sell,omitempty"`
	Total       string   `json:"total,omitempty"`
	Tag         string   `json:"tag,omitempty"`
	Keywords    string   `json:"keywords,omitempty"`
	Maps        []string `json:"maps,omitempty"`
	Types       []string `json:"types,omitempty"`
	From        any      `json:"from,omitempty"`
	SuitHeroId  any      `json:"suitHeroId,omitempty"`
	Into        any      `json:"into,omitempty"`
}

// --------------------------------------------------------

type LOLMEquipment struct {
	EquipList []EquipList `json:"equipList,omitempty"`
	Version   string      `json:"version,omitempty"`
	FileName  string      `json:"fileName,omitempty"`
	FileTime  string      `json:"fileTime,omitempty"`
}

type EquipList struct {
	EquipId         string        `json:"equipId,omitempty"`
	Name            string        `json:"name,omitempty"`
	IconPath        string        `json:"iconPath,omitempty"`
	From            []interface{} `json:"from,omitempty"`
	Type            string        `json:"type,omitempty"`
	Level           string        `json:"level,omitempty"`
	Price           string        `json:"price,omitempty"`
	Description     []string      `json:"description,omitempty"`
	Hp              string        `json:"hp,omitempty"`
	HpRegen         string        `json:"hpRegen,omitempty"`
	HpRegenRate     string        `json:"hpRegenRate,omitempty"`
	Armor           string        `json:"armor,omitempty"`
	ArmorPene       string        `json:"armorPene,omitempty"`
	ArmorPeneRate   string        `json:"armorPeneRate,omitempty"`
	CritRate        string        `json:"critRate,omitempty"`
	CritDamage      string        `json:"critDamage,omitempty"`
	AttackSpeed     string        `json:"attackSpeed,omitempty"`
	HealthPerAttack string        `json:"healthPerAttack,omitempty"`
	MagicAttack     string        `json:"magicAttack,omitempty"`
	Mp              string        `json:"mp,omitempty"`
	MpRegen         string        `json:"mpRegen,omitempty"`
	MagicBlock      string        `json:"magicBlock,omitempty"`
	MagicPene       string        `json:"magicPene,omitempty"`
	MagicPeneRate   string        `json:"magicPeneRate,omitempty"`
	HealthPerMagic  string        `json:"healthPerMagic,omitempty"`
	Cd              string        `json:"cd,omitempty"`
	DuctRate        string        `json:"ductRate,omitempty"`
	MoveSpeed       string        `json:"moveSpeed,omitempty"`
	MoveRate        string        `json:"moveRate,omitempty"`
	ComposeLevel    string        `json:"composeLevel,omitempty"`
	Ad              string        `json:"ad,omitempty"`
	Into            string        `json:"into,omitempty"`
	Tags            string        `json:"tags,omitempty"`
	UnName          string        `json:"unName,omitempty"`
	SearchKey       string        `json:"searchKey,omitempty"`
}

// --------------------------------------------------------

type EquipType struct {
	Cate    string    `json:"cate"`
	SubCate []SubCate `json:"sub_cate"`
}

type SubCate struct {
	Name          string   `json:"name"`
	KeywordsSlice []string `json:"keywordsSlice"`
	KeywordsStr   string   `json:"keywordsStr"`
}

// --------------------------------------------------------

type RespRoadmap struct {
	Current      Roadmap            `json:"current"`
	From         []Roadmap          `json:"from"`
	Into         []Roadmap          `json:"into"`
	GapPriceFrom int                `json:"gapPriceFrom"`
	SuitHeroes   []SearchResultList `json:"suit_heroes"`
}

type Roadmap struct {
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
}
