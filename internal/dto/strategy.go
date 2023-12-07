package dto

type StrategyHero struct {
	Keywords   string `json:"keywords"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	PublicDate string `json:"public_date"`
	Author     string `json:"author"`
	MainImg    string `json:"main_img"`
	JumpURL    string `json:"jump_url"`
	VideoID    string `json:"video_id"`
	Length     string `json:"length"`
	Source     int32  `json:"source"`
	Played     int64  `json:"played"`
	Platform   int    `json:"platform"`
}
