package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var GlobalConfig *Config

type Config struct {
	App      AppCfg      `yaml:"app"`
	Database DatabaseCfg `yaml:"database"`
	Redis    RedisCfg    `yaml:"redis"`
	MQ       MQCfg       `yaml:"mq"`
	ES       ESCfg       `yaml:"es"`
	Log      LogCfg      `yaml:"log"`
	Lol      LolCfg      `yaml:"lol"`
	LolM     LolmCfg     `yaml:"lolm"`
	Search   SearchCfg   `yaml:"search"`
}

type AppCfg struct {
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
}

type DatabaseCfg struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	DB       string `yaml:"db"`
}
type RedisCfg struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type MQCfg struct {
	Schema   string `yaml:"schema"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ESCfg struct {
	Host       string   `yaml:"host"`
	Port       string   `yaml:"port"`
	BuildIndex []string `yaml:"buildIndex"`
}

type LogCfg struct {
	LogLevel int    `yaml:"logLevel"`
	Path     string `yaml:"path"`
	SqlLog   string `yaml:"sqlLog"`
	EsLog    string `yaml:"esLog"`
}

type LolCfg struct {
	Equipment string `yaml:"equipment"`
	Heroes    string `yaml:"heroes"`
	Hero      string `yaml:"hero"`
	Rune      string `yaml:"rune"`
	Skill     string `yaml:"skill"`
}

type LolmCfg struct {
	Equipment       string `yaml:"equipment"`
	Heroes          string `yaml:"heroes"`
	Hero            string `yaml:"hero"`
	Rune            string `yaml:"rune"`
	RuneType        string `yaml:"runeType"`
	Skill           string `yaml:"skill"`
	RecommendHeroes string `yaml:"recommendHeroes"`
}
type SearchCfg struct {
	MapsLOL  []string `yaml:"mapsLOL"`
	MapsLOLM []string `yaml:"mapsLOLM"`
}

func Init() {
	path := "./configs/app.dev.yaml"
	//path := "/Users/limingze/GolandProjects/whisper/configs/app.dev.yaml"
	viper.SetConfigFile(path) // 指定配置文件路径
	//viper.SetConfigName("config")         // 配置文件名称(无扩展名)
	//viper.SetConfigType("yaml")           // 如果配置文件的名称中没有扩展名，则需要配置此项
	//viper.AddConfigPath("/etc/appname/")  // 查找配置文件所在的路径
	//viper.AddConfigPath("$HOME/.appname") // 多次调用以添加多个搜索路径
	//viper.AddConfigPath(".")              // 还可以在工作目录中查找配置
	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(&GlobalConfig)
	if err != nil {
		panic(fmt.Errorf("Failed to unmarshal config: %s \n", err))
		return
	}
}
