package dto

type RespHeroSkins struct {
	HeroId     string `json:"heroId,omitempty"`
	Title      string `json:"title,omitempty"`
	Name       string `json:"name,omitempty"`
	ShortBio   string `json:"ShortBio,omitempty"`
	MainImg    string `json:"mainImg,omitempty"`
	IconImg    string `json:"iconImg,omitempty"`
	LoadingImg string `json:"loadingImg,omitempty"`
	VideoImg   string `json:"videoImg,omitempty"`
	SourceImg  string `json:"sourceImg,omitempty"`
	Platform   int    `json:"platform,omitempty"`
	Version    string `json:"version,omitempty"`
}
