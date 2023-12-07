package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

var (
	NacosConfig  *Config
	ConfigClient config_client.IConfigClient
)

type Config struct {
	Nacos  NacosCfg  `yaml:"nacos"`
	LOL    LOLCfg    `yaml:"lol"`
	Equip  EquipCfg  `yaml:"equip"`
	Spider SpiderCfg `yaml:"spider"`
}

type SpiderCfg struct {
	NameSpace string `yaml:"nameSpace"`
	DataID    string `yaml:"dataID"`
	Group     string `yaml:"group"`
}

type EquipCfg struct {
	NameSpace string `yaml:"nameSpace"`
	DataID    string `yaml:"dataID"`
	Group     string `yaml:"group"`
}

type LOLCfg struct {
	NameSpace string `yaml:"nameSpace"`
	DataID    string `yaml:"dataID"`
	Group     string `yaml:"group"`
}

type NacosCfg struct {
	IP                  string `yaml:"ip"`
	Port                uint64 `yaml:"port"`
	TimeoutMs           uint64 `yaml:"timeoutMs"`
	NotLoadCacheAtStart bool   `yaml:"notLoadCacheAtStart"`
	LogDir              string `yaml:"logDir"`
	CacheDir            string `yaml:"cacheDir"`
	LogLevel            string `yaml:"logLevel"`
	NameSpace           string `yaml:"nameSpace"`
	DataID              string `yaml:"dataID"`
	Group               string `yaml:"group"`
}

func Init() {
	path := "./configs/app.dev.yaml"
	viper.SetConfigFile(path)   // 指定配置文件路径
	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(&NacosConfig)
	if err != nil {
		panic(fmt.Errorf("Failed to unmarshal config: %s \n", err))
	}

	//create clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         NacosConfig.Nacos.NameSpace,
		TimeoutMs:           NacosConfig.Nacos.TimeoutMs,
		NotLoadCacheAtStart: NacosConfig.Nacos.NotLoadCacheAtStart,
		LogDir:              NacosConfig.Nacos.LogDir,
		CacheDir:            NacosConfig.Nacos.CacheDir,
		LogLevel:            NacosConfig.Nacos.LogLevel,
	}

	// At least one ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      NacosConfig.Nacos.IP,
			ContextPath: "/nacos",
			Port:        NacosConfig.Nacos.Port,
			Scheme:      "http",
		},
	}

	ConfigClient, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}
}
