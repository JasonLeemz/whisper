package dto

type Version struct {
	LOL  VersionDetail `json:"LOL"`
	LOLM VersionDetail `json:"LOLM"`
}

type VersionDetail struct {
	Version    string `json:"version"`
	UpdateTime string `json:"update_time"`
}

type LOLMVersionList struct {
	Baton   string                `json:"baton,omitempty"`
	Code    int                   `json:"code,omitempty"`
	Data    []LOLMVersionListData `json:"data,omitempty"`
	ErrMsg  string                `json:"errMsg,omitempty"`
	Err_Msg string                `json:"err_msg,omitempty"`
	Msg     string                `json:"msg,omitempty"`
	Result  int                   `json:"result,omitempty"`
}

type LOLMVersionListData struct {
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
}

type LOLMVersionDetail struct {
	Baton   string                `json:"baton,omitempty"`
	Code    int                   `json:"code,omitempty"`
	Data    LOLMVersionDetailData `json:"data,omitempty"`
	ErrMsg  string                `json:"errMsg,omitempty"`
	Err_Msg string                `json:"err_msg,omitempty"`
	Msg     string                `json:"msg,omitempty"`
	Result  int                   `json:"result,omitempty"`
}

type BuyInfo struct {
	Text         string `json:"text,omitempty"`
	Intent       string `json:"intent,omitempty"`
	TextColor    string `json:"text_color,omitempty"`
	TextBgColor1 string `json:"text_bg_color1,omitempty"`
	TextBgColor2 string `json:"text_bg_color2,omitempty"`
}

type List struct {
	Title         string  `json:"title,omitempty"`
	Content       string  `json:"content,omitempty"`
	Descirbe      string  `json:"descirbe,omitempty"`
	MarkBgColor   string  `json:"mark_bg_color,omitempty"`
	ImgUrl        string  `json:"img_url,omitempty"`
	AttachContent string  `json:"attach_content,omitempty"`
	HeadUrl       string  `json:"head_url,omitempty"`
	VideoPlayUrl  string  `json:"video_play_url,omitempty"`
	HeroId        string  `json:"heroId,omitempty"`
	IsExpend      int     `json:"is_expend,omitempty"`
	IsEnableShare bool    `json:"is_enable_share,omitempty"`
	IsColorType   bool    `json:"is_color_type,omitempty"`
	BuyInfo       BuyInfo `json:"buy_info,omitempty"`
	HeroTag       string  `json:"hero_tag,omitempty"`
	ShowIndex     int     `json:"showIndex,omitempty"`
}

type LOLMVersionDetailData struct {
	ChangeCount string `json:"change_count,omitempty"`
	List        []List `json:"list,omitempty"`
	VersionKey  string `json:"version_key,omitempty"`
}
