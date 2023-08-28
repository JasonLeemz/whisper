package model

type EquipIntro struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Icon      string   `json:"icon"`
	Desc      string   `json:"desc"`
	Plaintext string   `json:"plaintext"`
	Price     float64  `json:"price"`
	Maps      string   `json:"maps"`
	Platform  uint8    `json:"platform"`
	Version   string   `json:"version"`
	Keywords  []string `json:"keywords"`
}

func (e *EquipIntro) CollectionName() string {
	return "equipment"
}
