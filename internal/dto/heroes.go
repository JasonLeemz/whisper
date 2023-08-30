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

// --------------------------------------------------------

// HeroDescription 英雄技能介绍
type HeroDescription struct {
	SpellKey        string `json:"spellKey"`
	Sort            int    `json:"sort"` // 用于展示数据的时候排序
	Name            string `json:"name"`
	Description     string `json:"description"`
	AbilityIconPath string `json:"abilityIconPath"`
	Detail          string `json:"detail"`
	Version         string `json:"version"`
}

// --------------------------------------------------------

type HeroSuitEquip struct {
	Ret   int    `json:"ret,omitempty"`
	Msg   string `json:"msg,omitempty"`
	IRet  int    `json:"iRet,omitempty"`
	SMsg  string `json:"sMsg,omitempty"`
	JData JData  `json:"jData,omitempty"`
}

type JDataData struct {
	Result string `json:"result,omitempty"`
}

type JData struct {
	Code      int       `json:"code,omitempty"`
	Msg       string    `json:"msg,omitempty"`
	Data      JDataData `json:"data,omitempty"`
	AmsSerial string    `json:"ams_serial,omitempty"`
}

type JDataDataResult struct {
	Timesectionwinrate    string `json:"timesectionwinrate,omitempty"`
	Styledetails          string `json:"styledetails,omitempty"`
	Spellid               string `json:"spellid,omitempty"`
	Itemover              string `json:"itemover,omitempty"`
	Itemout               string `json:"itemout,omitempty"`
	Itemcore              string `json:"itemcore,omitempty"`
	Itemshoes             string `json:"itemshoes,omitempty"`
	Doublechampiondetails string `json:"doublechampiondetails,omitempty"`
	Skilllist             string `json:"skilllist,omitempty"`
}

// --------------------------------------------------------

type ChampionFightData struct {
	Status  string `json:"status"`
	GameVer string `json:"gameVer"`
	Date    string `json:"date"`

	List FightDataList `json:"list"`
}

type FightDataList struct {
	ChampionLane map[string]ChampionLaneItem `json:"championLane"`
	//ChampionTrend interface{}            `json:"championTrend"`
	ChampionFight map[string]interface{} `json:"championFight"`
}

type Itemjson struct {
	Itemid   string `json:"itemid,omitempty"`
	Igamecnt int    `json:"igamecnt,omitempty"`
	Wincnt   int    `json:"wincnt,omitempty"`
	Winrate  int    `json:"winrate,omitempty"`
	Allcnt   int    `json:"allcnt,omitempty"`
	Showrate int    `json:"showrate,omitempty"`
}

type FightDataDetail struct {
	Itemout   map[string]Itemjson `json:"Itemout"`
	Core3item map[string]Itemjson `json:"Core3item"`
	Shoes     map[string]Itemjson `json:"Shoes"`
	Suits     []Itemjson          `json:"Suits"`
}
type ChampionLaneItem struct {

	//Itemoutjson   map[string]Itemjson `json:"itemoutjson,omitempty"`
	//Core3itemjson map[string]Itemjson `json:"core3itemjson,omitempty"`
	//Shoesjson     map[string]Itemjson `json:"shoesjson,omitempty"`
	//Hold3         []Itemjson          `json:"hold3,omitempty"`
	Itemout   map[string]Itemjson `json:"Itemout"`
	Core3item map[string]Itemjson `json:"Core3item"`
	Shoes     map[string]Itemjson `json:"Shoes"`
	Suits     []Itemjson          `json:"Suits"`

	Itemoutjson   string `json:"itemoutjson,omitempty"`
	Core3itemjson string `json:"core3itemjson,omitempty"`
	Shoesjson     string `json:"shoesjson,omitempty"`
	Hold3         string `json:"hold3,omitempty"`

	//Dtstatdate          string `json:"dtstatdate,omitempty"`
	//Championid          string `json:"championid,omitempty"`
	//Gameversion         string `json:"gameversion,omitempty"`
	//Lane                string `json:"lane,omitempty"`
	//Wincnt              string `json:"wincnt,omitempty"`
	//Igamecnt            string `json:"igamecnt,omitempty"`
	//Lanrate             string `json:"lanrate,omitempty"`
	//Lanewinrate         string `json:"lanewinrate,omitempty"`
	//Lanshowrate         string `json:"lanshowrate,omitempty"`
	//Champlanorder       string `json:"champlanorder,omitempty"`
	//Mainviceperk        string `json:"mainviceperk,omitempty"`
	//Perkdetail          string `json:"perkdetail,omitempty"`
	//Spellidjson         string `json:"spellidjson,omitempty"`
	//Skilljson           string `json:"skilljson,omitempty"`
	//WinrateFlowPlaytime string `json:"winrate_flow_playtime,omitempty"`
	//Deaths              string `json:"deaths,omitempty"`
	//Kills               string `json:"kills,omitempty"`
	//Assists             string `json:"assists,omitempty"`
	//Kda                 string `json:"kda,omitempty"`
	//Timeplayed          string `json:"timeplayed,omitempty"`
	//Goldearned          string `json:"goldearned,omitempty"`
	//Goldearnedpergame   string `json:"goldearnedpergame,omitempty"`
	//Minionskilled       string `json:"minionskilled,omitempty"`
	//Damagerate          string `json:"damagerate,omitempty"`
	//Damagetochamprate   string `json:"damagetochamprate,omitempty"`
	//Killsrate           string `json:"killsrate,omitempty"`
	//Hold1               string `json:"hold1,omitempty"`
	//Hold2               string `json:"hold2,omitempty"`
	//Hold4               string `json:"hold4,omitempty"`
	//Hold5               string `json:"hold5,omitempty"`
	//Hold6               string `json:"hold6,omitempty"`
	//Hold7               string `json:"hold7,omitempty"`
	//Hold8               string `json:"hold8,omitempty"`
	//Hold9               string `json:"hold9,omitempty"`
	//Hold10              string `json:"hold10,omitempty"`
	//EtlStamp            string `json:"etl_stamp,omitempty"`
}
type ChampionFight struct {
	Bottom  interface{} `json:"bottom"`
	Mid     interface{} `json:"mid"`
	Support interface{} `json:"support"`
	Top     interface{} `json:"top"`
	Jungle  interface{} `json:"jungle"`
}
