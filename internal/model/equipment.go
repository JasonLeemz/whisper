package model

import (
	"time"
)

type LOLEquipment struct {
	Id          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	ItemId      string    `gorm:"column:itemId;default:;NOT NULL"`
	Name        string    `gorm:"column:name;default:;NOT NULL"`
	IconPath    string    `gorm:"column:iconPath;default:;NOT NULL"`
	Price       string    `gorm:"column:price;default:;NOT NULL"`
	Description string    `gorm:"column:description;default:;NOT NULL"`
	Plaintext   string    `gorm:"column:plaintext;default:;NOT NULL"`
	Sell        string    `gorm:"column:sell;default:;NOT NULL"`
	Total       string    `gorm:"column:total;default:;NOT NULL"`
	SuitHeroId  string    `gorm:"column:suitHeroId;default:;NOT NULL"`
	Tag         string    `gorm:"column:tag;default:;NOT NULL"`
	Keywords    string    `gorm:"column:keywords;default:;NOT NULL"`
	Maps        string    `gorm:"column:maps;default:;NOT NULL"`
	From        string    `gorm:"column:from;default:;NOT NULL;comment:'合成自'"`
	Into        string    `gorm:"column:into;default:;NOT NULL;comment:'由谁合成'"`
	Types       string    `gorm:"column:types;default:;NOT NULL"`
	Version     string    `gorm:"column:version;default:;NOT NULL"`
	FileTime    string    `gorm:"column:fileTime;default:;NOT NULL"`
	Ctime       time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime       time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
	Status      uint8     `gorm:"column:status;default:0;NOT NULL"`
}

func (e *LOLEquipment) TableName() string {
	return "lol_equipment"
}

// --------------------------------------------------------

type LOLMEquipment struct {
	Id              uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	EquipId         string    `gorm:"column:equipId;default:;NOT NULL"`
	Name            string    `gorm:"column:name;default:;NOT NULL"`
	IconPath        string    `gorm:"column:iconPath;default:;NOT NULL"`
	From            string    `gorm:"column:from;default:;NOT NULL"`
	Type            string    `gorm:"column:type;default:;NOT NULL"`
	Level           string    `gorm:"column:level;default:;NOT NULL"`
	Price           string    `gorm:"column:price;default:;NOT NULL"`
	Description     string    `gorm:"column:description;default:;NOT NULL"`
	Hp              string    `gorm:"column:hp;default:;NOT NULL"`
	HpRegen         string    `gorm:"column:hpRegen;default:;NOT NULL"`
	HpRegenRate     string    `gorm:"column:hpRegenRate;default:;NOT NULL"`
	Armor           string    `gorm:"column:armor;default:;NOT NULL"`
	ArmorPene       string    `gorm:"column:armorPene;default:;NOT NULL"`
	ArmorPeneRate   string    `gorm:"column:armorPeneRate;default:;NOT NULL"`
	CritRate        string    `gorm:"column:critRate;default:;NOT NULL"`
	CritDamage      string    `gorm:"column:critDamage;default:;NOT NULL"`
	AttackSpeed     string    `gorm:"column:attackSpeed;default:;NOT NULL"`
	HealthPerAttack string    `gorm:"column:healthPerAttack;default:;NOT NULL"`
	MagicAttack     string    `gorm:"column:magicAttack;default:;NOT NULL"`
	Mp              string    `gorm:"column:mp;default:;NOT NULL"`
	MpRegen         string    `gorm:"column:mpRegen;default:;NOT NULL"`
	MagicBlock      string    `gorm:"column:magicBlock;default:;NOT NULL"`
	MagicPene       string    `gorm:"column:magicPene;default:;NOT NULL"`
	MagicPeneRate   string    `gorm:"column:magicPeneRate;default:;NOT NULL"`
	HealthPerMagic  string    `gorm:"column:healthPerMagic;default:;NOT NULL"`
	Cd              string    `gorm:"column:cd;default:;NOT NULL"`
	DuctRate        string    `gorm:"column:ductRate;default:;NOT NULL"`
	MoveSpeed       string    `gorm:"column:moveSpeed;default:;NOT NULL"`
	MoveRate        string    `gorm:"column:moveRate;default:;NOT NULL"`
	ComposeLevel    string    `gorm:"column:composeLevel;default:;NOT NULL"`
	Ad              string    `gorm:"column:ad;default:;NOT NULL"`
	Into            string    `gorm:"column:into;default:;NOT NULL"`
	Tags            string    `gorm:"column:tags;default:;NOT NULL"`
	UnName          string    `gorm:"column:unName;default:;NOT NULL"`
	SearchKey       string    `gorm:"column:searchKey;default:;NOT NULL"`
	Version         string    `gorm:"column:version;default:;NOT NULL"`
	FileTime        string    `gorm:"column:fileTime;default:;NOT NULL"`
	Ctime           time.Time `gorm:"column:ctime;default:current_timestamp();NOT NULL"`
	Utime           time.Time `gorm:"column:utime;default:current_timestamp();NOT NULL"`
	Status          uint8     `gorm:"column:status;default:0;NOT NULL"`
}

func (e *LOLMEquipment) TableName() string {
	return "lolm_equipment"
}
