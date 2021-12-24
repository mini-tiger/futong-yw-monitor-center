package g

import (
	"fmt"
	"futong-yw-monitor-center/monitor-base/bg"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

const defaultConfigFile = "config.yaml"

var ViperCfg LocalConfig

const Interval = 5 // 默认5分钟

func InitConfig() {
	v := viper.New()
	v.SetConfigFile(defaultConfigFile)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		logge.Warn("config file changed: %s", e.Name)
		PrintConf()
		//tmpVipCfg:=Config{}
		//if err := v.Unmarshal(&tmpVipCfg); err != nil {
		//	//log.Println("config Unmarshal err:",err)
		//	//log.Println("config use default Interval cron :",Interval)
		//	defaultConfLoad()
		//}
		//BeginChan<- struct{}{}

	})
	if err := v.Unmarshal(&ViperCfg); err != nil {
		log.Fatal("config fail err :", err)
	}

}
func CheckConfig() error {
	return bg.Validate.Struct(&ViperCfg)
}

type LocalConfig struct {
	Log     *Log     `mapstructure:"log" json:"log" yaml:"log" validate:"required"`
	ConfWeb *ConfWeb `mapstructure:"confWeb" json:"confWeb" yaml:"confWeb" validate:"required"`
}
type ConfWeb struct {
	GetConfUrl string `mapstructure:"getConfUrl" json:"getConfUrl" yaml:"getConfUrl" validate:"url,required"`
}

type Log struct {
	Stdout     bool   `mapstructure:"stdout" yaml:"stdout" json:"stdout"`
	LogMaxDays int    `mapstructure:"LogMaxDays" yaml:"logMaxDays" json:"logMaxDays"`
	Logfile    string `mapstructure:"logfile" yaml:"logfile" json:"logfile"`
	Level      string `mapstructure:"level" yaml:"level" json:"level"`
}
