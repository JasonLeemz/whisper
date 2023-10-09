package dto

type Version struct {
	LOL  LocalVersionDetail `json:"LOL"`
	LOLM LocalVersionDetail `json:"LOLM"`
}

type LocalVersionDetail struct {
	Version    string `json:"version"`
	UpdateTime string `json:"update_time"`
}

type PageVersionList struct {
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Title        string `json:"title,omitempty"`
	Vkey         string `json:"vkey,omitempty"`
	Introduction string `json:"introduction,omitempty"`
	Isnew        string `json:"isnew,omitempty"`
	Image        string `json:"image,omitempty"`
	PublicDate   string `json:"public_date,omitempty"`
	Platform     int    `json:"platform"`
}

type VersionList struct {
	Baton   string            `json:"baton,omitempty"`
	Code    int               `json:"code,omitempty"`
	Data    []VersionListData `json:"data,omitempty"`
	ErrMsg  string            `json:"errMsg,omitempty"`
	Err_Msg string            `json:"err_msg,omitempty"`
	Msg     string            `json:"msg,omitempty"`
	Result  int               `json:"result,omitempty"`
}

type VersionListData struct {
	Id           string   `json:"id,omitempty"`
	Name         string   `json:"name,omitempty"`
	Vkey         string   `json:"vkey,omitempty"`
	Title        string   `json:"title,omitempty"`
	Introduction string   `json:"introduction,omitempty"`
	Isnew        string   `json:"isnew,omitempty"`
	Image        string   `json:"image,omitempty"`
	PublicDate   string   `json:"public_date,omitempty"`
	Versions     []string `json:"versions,omitempty"`
	Intent       string   `json:"intent,omitempty"`
	//Name2  string `json:"Name,omitempty"`   // LOL
	//Change string `json:"change,omitempty"` // LOL
}

type VersionDetail struct {
	Baton   string            `json:"baton,omitempty"`
	Code    int               `json:"code,omitempty"`
	Data    VersionDetailData `json:"data,omitempty"`
	ErrMsg  string            `json:"errMsg,omitempty"`
	Err_Msg string            `json:"err_msg,omitempty"`
	Msg     string            `json:"msg,omitempty"`
	Result  int               `json:"result,omitempty"`
}

type BuyInfo struct {
	Text         string `json:"text,omitempty"`
	Intent       string `json:"intent,omitempty"`
	TextColor    string `json:"text_color,omitempty"`
	TextBgColor1 string `json:"text_bg_color1,omitempty"`
	TextBgColor2 string `json:"text_bg_color2,omitempty"`
}

type List struct {
	Title         string       `json:"title,omitempty"`
	Content       string       `json:"content,omitempty"`
	Descirbe      string       `json:"descirbe,omitempty"`
	MarkBgColor   string       `json:"mark_bg_color,omitempty"`
	ImgUrl        string       `json:"img_url,omitempty"`
	AttachContent string       `json:"attach_content,omitempty"`
	HeadUrl       string       `json:"head_url,omitempty"`
	VideoPlayUrl  string       `json:"video_play_url,omitempty"`
	HeroId        string       `json:"heroId,omitempty"`
	IsExpend      int          `json:"is_expend,omitempty"`
	IsEnableShare bool         `json:"is_enable_share,omitempty"`
	IsColorType   bool         `json:"is_color_type,omitempty"`
	BuyInfo       BuyInfo      `json:"buy_info,omitempty"`
	HeroTag       string       `json:"hero_tag,omitempty"`
	ShowIndex     int          `json:"showIndex,omitempty"`
	List          []DetailList `json:"list,omitempty"`

	ItemId int `json:"item_id,omitempty"` // LOL
}

type DetailList struct {
	Icon    string `json:"icon,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

type VersionDetailData struct {
	ChangeCount string `json:"change_count,omitempty"`
	List        []List `json:"list,omitempty"`
	VersionKey  string `json:"version_key,omitempty"`
}

type VersionInfo struct {
	Code    int             `json:"code,omitempty"`
	Data    VersionInfoData `json:"data,omitempty"`
	ErrMsg  string          `json:"errMsg,omitempty"`
	Err_Msg string          `json:"err_msg,omitempty"`
	Msg     string          `json:"msg,omitempty"`
	Result  int             `json:"result,omitempty"`
}

type Tabs struct {
	Title        string `json:"title,omitempty"`
	SchemeUrl    string `json:"schemeUrl,omitempty"`
	IsDefaultTab int    `json:"is_default_tab,omitempty"`
	Key          string `json:"key,omitempty"`
}

type VersionInfoData struct {
	ImgUrl  string `json:"imgUrl,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Tabs    []Tabs `json:"tabs,omitempty"`
}
