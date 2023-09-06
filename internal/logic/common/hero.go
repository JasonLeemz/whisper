package common

const (
	MaxHeroIDForLOL  = 950   // select max(CAST(heroId AS SIGNED)) from lol_heroes;
	MinHeroIDForLOLM = 10001 // select min(CAST(heroId AS SIGNED)) from lolm_heroes;
)
