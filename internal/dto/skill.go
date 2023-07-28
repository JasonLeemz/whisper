package dto

type LOLSkill struct {
	SummonerSkill map[string]SummonerSkill `json:"summonerskill,omitempty"`
	Version       string                   `json:"version,omitempty"`
	FileName      string                   `json:"fileName,omitempty"`
	FileTime      string                   `json:"fileTime,omitempty"`
}

type SummonerSkill struct {
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	SummonerLevel string `json:"summonerlevel,omitempty"`
	CoolDown      string `json:"cooldown,omitempty"`
	GameMode      string `json:"gamemode,omitempty"`
	Icon          string `json:"icon,omitempty"`
}

// --------------------------------------------------------

type LOLMSkill struct {
	SkillList []SkillList `json:"skillList,omitempty"`
	Version   string      `json:"version,omitempty"`
	FileName  string      `json:"fileName,omitempty"`
	FileTime  string      `json:"fileTime,omitempty"`
}

type SkillList struct {
	SkillId   string `json:"skillId,omitempty"`
	Name      string `json:"name,omitempty"`
	IconPath  string `json:"iconPath,omitempty"`
	FuncDesc  string `json:"funcDesc,omitempty"`
	Cd        string `json:"cd,omitempty"`
	Video     string `json:"video,omitempty"`
	UnlockLv  string `json:"unlocklv,omitempty"`
	Mode      string `json:"mode,omitempty"`
	SortOrder string `json:"sortOrder,omitempty"`
}
