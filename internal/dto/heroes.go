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
	IsARAMweekfree      string   `json:"isARAMweekfree,omitempty"`
	Ispermanentweekfree string   `json:"ispermanentweekfree,omitempty"`
	ChangeLabel         string   `json:"changeLabel,omitempty"`
	GoldPrice           string   `json:"goldPrice,omitempty"`
	CouponPrice         string   `json:"couponPrice,omitempty"`
	Camp                string   `json:"camp,omitempty"`
	CampId              string   `json:"campId,omitempty"`
	Keywords            string   `json:"keywords,omitempty"`
	InstanceId          string   `json:"instance_id,omitempty"`
}
