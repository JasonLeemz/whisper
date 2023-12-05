package dto

type StrategyHero struct {
	Desc       string `json:"desc"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	PublicDate string `json:"public_date"`
	Author     string `json:"author"`
	MainImg    string `json:"main_img"`
	JumpURL    string `json:"jump_url"`
	Source     int32  `json:"source"`
	Played     int64  `json:"played"`
	Platform   int    `json:"platform"`
}
