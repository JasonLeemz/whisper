package dto

type SearchKeywords struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Ttl     int    `json:"ttl,omitempty"`
	Data    Data   `json:"data,omitempty"`
}

type BList struct {
	Vlist []Vlist       `json:"vlist,omitempty"`
	Slist []interface{} `json:"slist,omitempty"`
}
type Page struct {
	Pn    int `json:"pn,omitempty"`
	Ps    int `json:"ps,omitempty"`
	Count int `json:"count,omitempty"`
}
type Data struct {
	List        BList `json:"list,omitempty"`
	Page        Page  `json:"page,omitempty"`
	IsRisk      bool  `json:"is_risk,omitempty"`
	GaiaResType int   `json:"gaia_res_type,omitempty"`
}

type Vlist struct {
	Comment          int    `json:"comment,omitempty"`
	Typeid           int    `json:"typeid,omitempty"`
	Play             int64  `json:"play,omitempty"`
	Pic              string `json:"pic,omitempty"`
	Subtitle         string `json:"subtitle,omitempty"`
	Description      string `json:"description,omitempty"`
	Copyright        string `json:"copyright,omitempty"`
	Title            string `json:"title,omitempty"`
	Review           int    `json:"review,omitempty"`
	Author           string `json:"author,omitempty"`
	Mid              int    `json:"mid,omitempty"`
	Created          int64  `json:"created,omitempty"`
	Length           string `json:"length,omitempty"`
	VideoReview      int    `json:"video_review,omitempty"`
	Aid              int    `json:"aid,omitempty"`
	Bvid             string `json:"bvid,omitempty"`
	HideClick        bool   `json:"hide_click,omitempty"`
	IsPay            int    `json:"is_pay,omitempty"`
	IsUnionVideo     int    `json:"is_union_video,omitempty"`
	IsSteinsGate     int    `json:"is_steins_gate,omitempty"`
	IsLivePlayback   int    `json:"is_live_playback,omitempty"`
	IsLessonVideo    int    `json:"is_lesson_video,omitempty"`
	IsLessonFinished int    `json:"is_lesson_finished,omitempty"`
	LessonUpdateInfo string `json:"lesson_update_info,omitempty"`
	JumpUrl          string `json:"jump_url,omitempty"`
	Meta             Meta   `json:"meta,omitempty"`
	IsAvoided        int    `json:"is_avoided,omitempty"`
	SeasonId         int    `json:"season_id,omitempty"`
	Attribute        int    `json:"attribute,omitempty"`
	IsChargingArc    bool   `json:"is_charging_arc,omitempty"`
	Vt               int    `json:"vt,omitempty"`
	EnableVt         int    `json:"enable_vt,omitempty"`
	VtDisplay        string `json:"vt_display,omitempty"`
	PlaybackPosition int    `json:"playback_position,omitempty"`
}

type Meta struct {
	Id        int    `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	Cover     string `json:"cover,omitempty"`
	Mid       int    `json:"mid,omitempty"`
	Intro     string `json:"intro,omitempty"`
	SignState int    `json:"sign_state,omitempty"`
	Attribute int    `json:"attribute,omitempty"`
	Stat      Stat   `json:"stat,omitempty"`
	EpCount   int    `json:"ep_count,omitempty"`
	FirstAid  int    `json:"first_aid,omitempty"`
	Ptime     int    `json:"ptime,omitempty"`
	EpNum     int    `json:"ep_num,omitempty"`
}

type Stat struct {
	SeasonId int `json:"season_id,omitempty"`
	View     int `json:"view,omitempty"`
	Danmaku  int `json:"danmaku,omitempty"`
	Reply    int `json:"reply,omitempty"`
	Favorite int `json:"favorite,omitempty"`
	Coin     int `json:"coin,omitempty"`
	Share    int `json:"share,omitempty"`
	Like     int `json:"like,omitempty"`
	Mtime    int `json:"mtime,omitempty"`
	Vt       int `json:"vt,omitempty"`
	Vv       int `json:"vv,omitempty"`
}
