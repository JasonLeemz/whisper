package dto

type Version struct {
	LOL  VersionDetail `json:"LOL"`
	LOLM VersionDetail `json:"LOLM"`
}

type VersionDetail struct {
	Version    string `json:"version"`
	UpdateTime string `json:"update_time"`
}
