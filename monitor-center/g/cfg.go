package g

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mini-tiger/tjtools/file"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var cfg Config
var configLock = new(sync.RWMutex)

type Config struct {
	Port       int      `json:"port"`
	MysqlDsn   string   `json:"mysqlDsn"`
	Logfile    string   `json:"logfile"`
	LogMaxDays int      `json:"logMaxDays"`
	Level      string   `json:"level"`
	Stdout     bool     `json:"stdout"`
	EsServer   []string `json:"es_server"`
	//AgentDownLoadUrl     string   `json:"agentDownLoadUrl"`
	DelPushGatewayMinute int      `json:"delPushGatewayMinute"`
	DelDirtyDataMinute   int      `json:"delDirtyDataMinute"`
	GetConfUrl           string   `json:"getConfUrl"`
	PushGateway          []string `json:"pushGateway"`
	AlertManagerRulesDir string   `json:"alert_manager_rules_dir"`
	Prometheus           []string `json:"prometheus"`
}

func (c *Config) IsDebug() bool {
	return strings.ToLower(cfg.Level) == "debug"
}

func ParseConfig(cfg string) string {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)

	}

	//var c GlobalConfig

	lock.Lock()
	defer lock.Unlock()

	log.Println("read config file:", cfg, "successfully")
	return configContent
	//WLog(fmt.Sprintf("read config file: %s successfully",cfg))
}

func CheckConfig(fp string) (e error, conf string) {
	// 兼容开发与生产环境

	if file.IsExist(fp) {
		return nil, fp
	} else {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		if file.IsExist(filepath.Join(dir, fp)) {
			return nil, filepath.Join(dir, fp)
		} else {
			return errors.New(fmt.Sprintf("confile :%s Not Found", fp)), ""
		}
	}

}

func readconfig(cfgfile string) {
	cfgstr := ParseConfig(cfgfile)
	err := json.Unmarshal([]byte(cfgstr), &cfg)
	if err != nil {
		log.Fatalln("parse config file fail:", err)
	}

	// var
	//TokenTableName = cfg.TokenTableName
	//ClientTableName = cfg.ClientTableName

}

func LoadConfig(cfgPath string) (e error) {
	var confile string
	//var e error
	// 多级目录查找配置文件
	for _, basedir := range Basedirs {
		e, confile = CheckConfig(path.Join(basedir, cfgPath))
		if e == nil {
			CurrentDir = basedir
			log.Printf("Work Dir:%s\n", basedir)
			break
		}
	}

	if e == nil {
		readconfig(confile)
		//cfgStr, _ := json.MarshalIndent(cfg, "", "\t")
		//log.Printf("config file read success! data:%+v\n", string(cfgStr))
	} else {
		log.Fatalln("config file fail :", e)
	}
	return
}

func GetConfig() *Config {

	configLock.RLock()
	defer configLock.RUnlock()
	return &cfg
}
