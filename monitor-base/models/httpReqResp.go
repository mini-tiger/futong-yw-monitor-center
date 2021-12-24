package models

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  httpReqResp
 * @Version: 1.0.0
 * @Date: 2021/12/8 下午5:59
 */

type RespData struct {
	Code string
	Msg  string
	Data interface{}
}

type ReqMetricsData struct {
	HostId  string                   `json:"hostId" validate:"required"`
	Metrics []string                 `mapstructure:"metrics" json:"metrics" yaml:"metrics" validate:"required"`
	Shells  []*MonitorMetricsDefault `mapstructure:"shells" json:"shells" yaml:"shells" validate:"required"`
}
