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

// UserDynamic https://api.bilibili.com/x/space/dynamic/search?keyword=%E7%90%B4%E5%A5%B3&pn=1&ps=30&mid=424730226
type UserDynamic struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Ttl     int         `json:"ttl"`
	Data    DynamicData `json:"data"`
}

type DynamicData struct {
	Cards []Cards `json:"cards"`
	Total int     `json:"total"`
}

type Cards struct {
	Desc Desc   `json:"desc"`
	Card string `json:"card"` // DynamicCardsCard
	//ExtendJson string `json:"extend_json"`
	//Display    Display `json:"display"`
}

type DynamicCardsCard struct {
	Aid         int    `json:"aid"`
	Attribute   int    `json:"attribute"`
	AttributeV2 int    `json:"attribute_v2"`
	Cid         int    `json:"cid"`
	Copyright   int    `json:"copyright"`
	Ctime       int    `json:"ctime"`
	Desc        string `json:"desc"`
	//Dimension   Dimension `json:"dimension"`
	Duration   int    `json:"duration"`
	Dynamic    string `json:"dynamic"`
	FirstFrame string `json:"first_frame"`
	JumpUrl    string `json:"jump_url"`
	Owner      struct {
		Face string `json:"face"`
		Mid  int    `json:"mid"`
		Name string `json:"name"`
	} `json:"owner"`
	Pic         string `json:"pic"`
	PubLocation string `json:"pub_location"`
	Pubdate     int    `json:"pubdate"`
	//Rights      Rights `json:"rights"`
	SeasonId    int    `json:"season_id"`
	ShortLinkV2 string `json:"short_link_v2"`
	Stat        struct {
		Aid      int   `json:"aid"`
		Coin     int   `json:"coin"`
		Danmaku  int   `json:"danmaku"`
		Dislike  int   `json:"dislike"`
		Favorite int   `json:"favorite"`
		HisRank  int   `json:"his_rank"`
		Like     int   `json:"like"`
		NowRank  int   `json:"now_rank"`
		Reply    int   `json:"reply"`
		Share    int   `json:"share"`
		View     int64 `json:"view"`
		Vt       int   `json:"vt"`
		Vv       int   `json:"vv"`
	} `json:"stat"`
	State     int    `json:"state"`
	Tid       int    `json:"tid"`
	Title     string `json:"title"`
	Tname     string `json:"tname"`
	Videos    int    `json:"videos"`
	VtDisplay string `json:"vt_display"`
}

type Desc struct {
	Uid          int         `json:"uid"`
	Type         int         `json:"type"`
	Rid          int         `json:"rid"`
	Acl          int         `json:"acl"`
	View         int         `json:"view"`
	Repost       int         `json:"repost"`
	Like         int         `json:"like"`
	IsLiked      int         `json:"is_liked"`
	DynamicId    int         `json:"dynamic_id"`
	Timestamp    int64       `json:"timestamp"`
	PreDyId      int         `json:"pre_dy_id"`
	OrigDyId     int         `json:"orig_dy_id"`
	OrigType     int         `json:"orig_type"`
	UserProfile  UserProfile `json:"user_profile"`
	UidType      int         `json:"uid_type"`
	Stype        int         `json:"stype"`
	RType        int         `json:"r_type"`
	InnerId      int         `json:"inner_id"`
	Status       int         `json:"status"`
	DynamicIdStr string      `json:"dynamic_id_str"`
	PreDyIdStr   string      `json:"pre_dy_id_str"`
	OrigDyIdStr  string      `json:"orig_dy_id_str"`
	RidStr       string      `json:"rid_str"`
	Bvid         string      `json:"bvid"`
}

type UserProfile struct {
	Info Info `json:"info"`
	//Card    Card    `json:"card"`
	//Vip     Vip     `json:"vip"`
	//Pendant Pendant `json:"pendant"`
	//Rank    string  `json:"rank"`
	//Sign    string  `json:"sign"`
	//LevelInfo LevelInfo `json:"level_info"`
}

type Info struct {
	Uid     int    `json:"uid"`
	Uname   string `json:"uname"`
	Face    string `json:"face"`
	FaceNft int    `json:"face_nft"`
}

type Relation struct {
	Status     int `json:"status"`
	IsFollow   int `json:"is_follow"`
	IsFollowed int `json:"is_followed"`
}

type LiveInfo struct {
	LiveStatus int    `json:"live_status"`
	LiveUrl    string `json:"live_url"`
}
