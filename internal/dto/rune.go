package dto

type LOLRune struct {
	Rune     map[string]Rune `json:"rune,omitempty"`
	Version  string          `json:"version,omitempty"`
	FileName string          `json:"fileName,omitempty"`
	FileTime string          `json:"fileTime,omitempty"`
}

type Rune struct {
	Name      string `json:"name,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Key       string `json:"key,omitempty"`
	Tooltip   string `json:"tooltip,omitempty"`
	ShortDesc string `json:"shortdesc,omitempty"`
	LongDesc  string `json:"longdesc,omitempty"`
	SlotLabel string `json:"slotLabel,omitempty"`
	StyleName string `json:"styleName,omitempty"`
}

// --------------------------------------------------------

type LOLMRune struct {
	RuneList []RuneList `json:"runeList,omitempty"`
	Version  string     `json:"version,omitempty"`
	FileName string     `json:"fileName,omitempty"`
	FileTime string     `json:"fileTime,omitempty"`
}

type RuneList struct {
	RuneId               string `json:"runeId,omitempty"`
	Name                 string `json:"name,omitempty"`
	Description          string `json:"description,omitempty"`
	DetailInfo           string `json:"detailInfo,omitempty"`
	AttrName             string `json:"attrName,omitempty"`
	Type                 string `json:"type,omitempty"`
	IconPath             string `json:"iconPath,omitempty"`
	SortOrder            string `json:"sortOrder,omitempty"`
	UnlockLv             string `json:"unlockLv,omitempty"`
	PrimarySlotIndex     string `json:"primarySlotIndex,omitempty"`
	PrimarySlotSortOrder string `json:"primarySlotSortOrder,omitempty"`
}

// --------------------------------------------------------

// LOLMRuneType 符文类别
type LOLMRuneType struct {
	MainPage   MainPage    `json:"main_page,omitempty"`
	MainTabs   MainTabs    `json:"main_tabs,omitempty"`
	MainTagKey string      `json:"main_tag_key,omitempty"`
	RuneTypes  []RuneTypes `json:"runeTypes,omitempty"`
}

type MainPage struct {
	IsNavRectEdge     int `json:"is_nav_rect_edge,omitempty"`
	JustNowUpdate     int `json:"just_now_update,omitempty"`
	MinNumShowTabView int `json:"min_num_show_tab_view,omitempty"`
	RefreshFlag       int `json:"refresh_flag,omitempty"`
	StatusBarStyle    int `json:"status_bar_style,omitempty"`
}

type MainTabs struct {
	TabUrl string `json:"tab_url,omitempty"`
}

type RuneTypes struct {
	Name    string `json:"name,omitempty"`
	SubType string `json:"sub_type,omitempty"`
	Type    string `json:"type,omitempty"`
}
