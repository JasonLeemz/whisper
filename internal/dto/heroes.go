package dto

type LOLHeroes struct {
	Hero     []Hero `json:"hero,omitempty"`
	Version  string `json:"version,omitempty"`
	FileName string `json:"fileName,omitempty"`
	FileTime string `json:"fileTime,omitempty"`
}

type Hero struct {
	HeroId              string   `json:"heroId,omitempty"`
	Name                string   `json:"name,omitempty"`
	Alias               string   `json:"alias,omitempty"`
	Title               string   `json:"title,omitempty"`
	Roles               []string `json:"roles,omitempty"`
	IsWeekFree          string   `json:"isWeekFree,omitempty"`
	Attack              string   `json:"attack,omitempty"`
	Defense             string   `json:"defense,omitempty"`
	Magic               string   `json:"magic,omitempty"`
	Difficulty          string   `json:"difficulty,omitempty"`
	SelectAudio         string   `json:"selectAudio,omitempty"`
	BanAudio            string   `json:"banAudio,omitempty"`
	IsARAMWeekFree      string   `json:"isARAMweekfree,omitempty"`
	IsPermanentWeekFree string   `json:"ispermanentweekfree,omitempty"`
	ChangeLabel         string   `json:"changeLabel,omitempty"`
	GoldPrice           string   `json:"goldPrice,omitempty"`
	CouponPrice         string   `json:"couponPrice,omitempty"`
	Camp                string   `json:"camp,omitempty"`
	CampId              string   `json:"campId,omitempty"`
	Keywords            string   `json:"keywords,omitempty"`
	InstanceId          string   `json:"instance_id,omitempty"`
}

// --------------------------------------------------------

type LOLMHeroes struct {
	HeroList map[string]HeroInfo `json:"heroList,omitempty"`
	Version  string              `json:"version,omitempty"`
	FileName string              `json:"fileName,omitempty"`
	FileTime string              `json:"fileTime,omitempty"`
}
type HeroInfo struct {
	HeroId         string   `json:"heroId,omitempty"`
	Name           string   `json:"name,omitempty"`
	Title          string   `json:"title,omitempty"`
	Roles          []string `json:"roles,omitempty"`
	Intro          string   `json:"intro,omitempty"`
	Avatar         string   `json:"avatar,omitempty"`
	Card           string   `json:"card,omitempty"`
	Poster         string   `json:"poster,omitempty"`
	HighlightPrice string   `json:"highlightprice,omitempty"`
	CouponPrice    string   `json:"couponprice,omitempty"`
	Alias          string   `json:"alias,omitempty"`
	Lane           string   `json:"lane,omitempty"`
	Tags           string   `json:"tags,omitempty"`
	SearchKey      string   `json:"searchkey,omitempty"`
	IsWeekFree     string   `json:"isWeekFree,omitempty"`
	DifficultyL    string   `json:"difficultyL,omitempty"`
	Damage         string   `json:"damage,omitempty"`
	SurviveL       string   `json:"surviveL,omitempty"`
	AssistL        string   `json:"assistL,omitempty"`
}

// --------------------------------------------------------

type HeroAttribute struct {
	Hero     HeroBaseInfo `json:"hero,omitempty"`
	Skins    []Skins      `json:"skins,omitempty"`
	Spells   []Spells     `json:"spells,omitempty"`
	Version  string       `json:"version,omitempty"`
	FileName string       `json:"fileName,omitempty"`
	FileTime string       `json:"fileTime,omitempty"`
}

type HeroBaseInfo struct {
	HeroId string   `json:"heroId,omitempty"`
	Title  string   `json:"title,omitempty"` // LOL 布隆, LOLM 弗雷尔卓德之心
	Name   string   `json:"name,omitempty"`  // LOL 弗雷尔卓德之心, LOLM 布隆
	Alias  string   `json:"alias,omitempty"`
	Roles  []string `json:"roles,omitempty"`

	Defense string `json:"defense,omitempty"` // 防守 仅LOL有值
	Magic   string `json:"magic,omitempty"`

	Difficulty  string `json:"difficulty,omitempty"` // 仅LOL有值
	DifficultyL string `json:"difficultyL,omitempty"`

	Attack              string `json:"attack,omitempty"`
	AttackRange         string `json:"attackrange,omitempty"`  // 仅LOL有值
	AttackDamage        string `json:"attackdamage,omitempty"` // 仅LOL有值
	Attackspeed         string `json:"attackspeed,omitempty"`
	Attackspeedperlevel string `json:"attackspeedperlevel,omitempty"`
	Hp                  string `json:"hp,omitempty"`
	Hpperlevel          string `json:"hpperlevel,omitempty"`
	Mp                  string `json:"mp,omitempty"`
	Mpperlevel          string `json:"mpperlevel,omitempty"`
	Movespeed           string `json:"movespeed,omitempty"`
	Armor               string `json:"armor,omitempty"`
	Armorperlevel       string `json:"armorperlevel,omitempty"`
	Spellblock          string `json:"spellblock,omitempty"`
	Spellblockperlevel  string `json:"spellblockperlevel,omitempty"`
	Hpregen             string `json:"hpregen,omitempty"`
	Hpregenperlevel     string `json:"hpregenperlevel,omitempty"`
	Mpregen             string `json:"mpregen,omitempty"`
	Mpregenperlevel     string `json:"mpregenperlevel,omitempty"`
	Crit                string `json:"crit,omitempty"`

	Damage     string `json:"damage,omitempty"`
	Durability string `json:"durability,omitempty"`

	Mobility string `json:"mobility,omitempty"`

	Avatar         string `json:"avatar,omitempty"`         // 仅LOLM有值,LOL可以自行拼接或者取Skins[0] https://game.gtimg.cn/images/lol/act/img/skin/small201000.jpg
	Highlightprice string `json:"highlightprice,omitempty"` // 仅LOLM有值
	GoldPrice      string `json:"goldPrice,omitempty"`      // 仅LOL有值
	Couponprice    string `json:"couponprice,omitempty"`

	//Card             string `json:"card,omitempty"`
	//Poster           string `json:"poster,omitempty"`

	IsWeekFree string `json:"isWeekFree,omitempty"`
}

type Skins struct {
	SkinId          string `json:"skinId,omitempty"`
	HeroId          string `json:"heroId,omitempty"`
	HeroName        string `json:"heroName,omitempty"`
	HeroTitle       string `json:"heroTitle,omitempty"`
	Name            string `json:"name,omitempty"`
	Chromas         string `json:"chromas,omitempty"`
	ChromasBelongId string `json:"chromasBelongId,omitempty"`
	IsBase          string `json:"isBase,omitempty"`
	EmblemsName     string `json:"emblemsName,omitempty"`
	Description     string `json:"description,omitempty"`
	MainImg         string `json:"mainImg,omitempty"`
	IconImg         string `json:"iconImg,omitempty"`
	LoadingImg      string `json:"loadingImg,omitempty"`
	VideoImg        string `json:"videoImg,omitempty"`
	SourceImg       string `json:"sourceImg,omitempty"`
	VedioPath       string `json:"vedioPath,omitempty"`
	SuitType        string `json:"suitType,omitempty"`
	PublishTime     string `json:"publishTime,omitempty"`
	ChromaImg       string `json:"chromaImg,omitempty"`
}

type Spells struct {
	HeroId          string `json:"heroId,omitempty"`
	SpellKey        string `json:"spellKey,omitempty"` // LOL:passive q w e r, LOLM:passive active active active active
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"` // LOLM:通过普通攻击和寒冬之咬来对目标施加标记。在施加第一层标记后，我方英雄的普攻也会叠加此标记。  当目标身上的此标记达到4层，就会陷入晕眩1秒并承受32魔法伤害
	AbilityIconPath string `json:"abilityIconPath,omitempty"`
	//AbilityVideoPath string `json:"abilityVideoPath,omitempty"`

	Detail []string `json:"detail,omitempty"` // 仅LOLM有值 通过普通攻击和寒冬之咬来对目标施加标记，持续4秒。在施加第一层标记后，我方英雄的普攻命中也会对目标叠加此标记。  当目标身上的此标记达到4层，就会陷入晕眩1秒并承受32魔法伤害（*自身等级）。  敌人在晕眩后的8秒内不再获得标记层数，但被布隆的普攻命中时会受到额外7魔法伤害（*自身等级）。
	//Cdtime    string   `json:"cdtime,omitempty"` // 仅LOLM有值
	//Costtype  string   `json:"costtype,omitempty"`
	//Costvalue string   `json:"costvalue,omitempty"`
	//SpellId string `json:"spellId,omitempty"`
}
