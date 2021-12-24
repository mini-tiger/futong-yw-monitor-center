package bg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mini-tiger/tjtools/file"
	"io/ioutil"
	"sync"
)

/**
 * @Author: Tao Jun
 * @Description: g
 * @File:  pushconf
 * @Version: 1.0.0
 * @Date: 2021/11/19 上午11:33
 */

const DefaultPushConfigFile = "PushConfig.json"

func NewRespPushConfig() *RespPushConfigData {
	return &RespPushConfigData{}
}

type RespPushConfigData struct {
	Code string
	Msg  string
	Data *pushConfig
}

type pushConfig struct {
	System      System   `mapstructure:"system" json:"system" yaml:"system" validate:"required"`
	Es          Es       `mapstructure:"es" json:"es" yaml:"es" validate:"required"`
	PushGateWay []string `mapstructure:"pushGateway" json:"pushGateway" yaml:"pushGateway" validate:"required"`
	Metrics     []string `mapstructure:"metrics" json:"metrics" yaml:"metrics" validate:"required"`
	//Version    uint64     `mapstructure:"version" json:"version" yaml:"version" validate:"gte=1,required"`
	GetSelfUpdateUrl string      `mapstructure:"getSelfUpdateUrl" json:"getSelfUpdateUrl" yaml:"getSelfUpdateUrl" validate:"url,required"`
	AddHostInfoUrl   string      `mapstructure:"addHostInfoUrl" json:"addHostInfoUrl" yaml:"addHostInfoUrl" validate:"url,required"`
	Shells           interface{} `mapstructure:"shells" json:"shells" yaml:"shells" validate:"required"`
}

func NewPushConfig() *pushConfig {
	return &pushConfig{}
}

type System struct {
	Interval int `mapstructure:"interval" json:"interval" yaml:"interval" validate:"gte=1,required"`
	WebPort  int `mapstructure:"webPort" json:"webPort" yaml:"webPort" validate:"gte=1"`
}

type Es struct {
	Host string `mapstructure:"host" json:"host" yaml:"host" validate:"url,required"`
}

//type PushGateWay struct {
//	PushGateWayAddressFirst  string `mapstructure:"pushGateWayAddressFirst" json:"pushGateWayAddressFirst" yaml:"pushGateWayAddressFirst" validate:"required,url"`
//	PushGateWayAddressSecond string `mapstructure:"pushGateWayAddressSecond" json:"pushGateWayAddressSecond" yaml:"pushGateWayAddressSecond" validate:"required,url"`
//	PushGateWayAddressThird  string `mapstructure:"pushGateWayAddressThird" json:"pushGateWayAddressThird" yaml:"pushGateWayAddressThird" validate:"required,url"`
//}

// 服务 端默认配置只能有一个
var PushCfgEntry *pushConfigEntry = &pushConfigEntry{pushConfig: &pushConfig{}}

type pushConfigEntry struct {
	sync.RWMutex
	*pushConfig
}

func InitPushConfig() error {
	if !file.IsExist(DefaultPushConfigFile) {
		return errors.New(fmt.Sprintf("初始化加载 PushConfig 失败 %s 文件不存在\n", DefaultPushConfigFile))
	}
	bs, err := ioutil.ReadFile(DefaultPushConfigFile)
	if err != nil {
		return errors.New(fmt.Sprintf("初始化加载 PushConfig 读取失败 %s \n", err.Error()))

	}
	// 可能重复写入文件
	err = PushCfgEntry.SetConf(bs)
	if err != nil {

		return errors.New(fmt.Sprintf("初始化加载 PushConfig SetConf失败 %s \n", err.Error()))
	}
	return nil
}

func (p *pushConfigEntry) SetConf(bs []byte) (err error) {
	p.Lock()
	defer p.Unlock()
	tp := &pushConfig{}
	err = json.Unmarshal(bs, tp)
	if err != nil {
		return err
	}
	return p.SetStruct(tp)

}

func (p *pushConfigEntry) SetStruct(tp *pushConfig) (err error) {
	err = Validate.Struct(tp)
	if err != nil {
		return err
	}

	p.pushConfig = tp
	go p.WriteConf()

	return nil
}
func (p *pushConfigEntry) WriteConf() error {
	p.RLock()
	defer p.RUnlock()
	content, _ := json.MarshalIndent(p.pushConfig, "", "\t")
	err := ioutil.WriteFile(DefaultPushConfigFile, content, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("Json 写入失败 :%s\n", err.Error()))
	} else {
		return errors.New("Json 写入成功\n")
	}

}

func (p *pushConfigEntry) GetConf() *pushConfig {
	p.RLock()
	p.RUnlock()
	return p.pushConfig

}
