package dto

type RespHeroSkins struct {
	HeroId      string `json:"heroId,omitempty"`
	SkinName    string `json:"skinName,omitempty"`
	HeroName    string `json:"heroName,omitempty"`
	HeroTitle   string `json:"heroTitle,omitempty"`
	Description string `json:"description,omitempty"`
	MainImg     string `json:"mainImg,omitempty"`
	IconImg     string `json:"iconImg,omitempty"`
	LoadingImg  string `json:"loadingImg,omitempty"`
	VideoImg    string `json:"videoImg,omitempty"`
	SourceImg   string `json:"sourceImg,omitempty"`
	Platform    int    `json:"platform,omitempty"`
	Version     string `json:"version,omitempty"`
}
