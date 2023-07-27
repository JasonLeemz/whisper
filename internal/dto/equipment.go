package dto

type LOLEquipment struct {
	Tree     []Tree  `json:"tree,omitempty"`
	Items    []Items `json:"items,omitempty"`
	Version  string  `json:"version,omitempty"`
	FileName string  `json:"fileName,omitempty"`
	FileTime string  `json:"fileTime,omitempty"`
}

type Tree struct {
	Header string   `json:"header,omitempty"`
	Tags   []string `json:"tags,omitempty"`
}

type Items struct {
	ItemId      string   `json:"itemId,omitempty"`
	Name        string   `json:"name,omitempty"`
	IconPath    string   `json:"iconPath,omitempty"`
	Price       string   `json:"price,omitempty"`
	Description string   `json:"description,omitempty"`
	Plaintext   string   `json:"plaintext,omitempty"`
	Sell        string   `json:"sell,omitempty"`
	Total       string   `json:"total,omitempty"`
	Tag         string   `json:"tag,omitempty"`
	Keywords    string   `json:"keywords,omitempty"`
	Maps        []string `json:"maps,omitempty"`
	Types       []string `json:"types,omitempty"`
	From        any      `json:"from,omitempty"`
	SuitHeroId  any      `json:"suitHeroId,omitempty"`
	Into        any      `json:"into,omitempty"`
}
