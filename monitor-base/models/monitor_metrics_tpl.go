package models

import (
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  agent
 * @Version: 1.0.0
 * @Date: 2021/11/25 上午11:20
 */
type MonitorMetricsTpl struct {
	ID        uint `gorm:"primary_key" json:"id"`
	UpdatedAt time.Time
	Name      string `gorm:"column:name; " json:"name" validate:"required"`
	Desc      string `gorm:"column:desc;" json:"desc"`
	//同一个主机同一个指标 只能采集一次（shell）
	MonitorDeviceID uint64                   `gorm:"column:monitor_device_id;comment:'关联MetricsTpl'" json:"Metrics_tpl_id"`
	MonitorMetrics  []*MonitorMetricsDefault `gorm:"many2many:itm_yw_monitor_metrics_tpl;"`
}

func (*MonitorMetricsTpl) TableName() string {
	return "itm_yw_monitor_metrics_tpl"
}
