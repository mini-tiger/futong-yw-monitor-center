package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  agent
 * @Version: 1.0.0
 * @Date: 2021/11/25 上午11:20
 */
type MonitorMetrics struct {
	ID         uint `gorm:"primary_key" json:"id"`
	UpdatedAt  time.Time
	MetricName string `gorm:"column:metricname; size:64; UNIQUE_INDEX:uniqe_monitordevice_metric_MetricName" json:"metricname" validate:"required"`
	Desc       string `gorm:"column:desc;" json:"desc"`
	Ident      string `gorm:"column:ident;" json:"ident"  validate:"required"`
	//同一个主机同一个指标 只能采集一次（shell）
	MonitorDeviceID         int64 `gorm:"column:monitor_device_id; UNIQUE_INDEX:uniqe_monitordevice_metric_MetricName" json:"monitor_device_id"`                   // 关联monitor_device  一对多,本表多
	MonitorMetricsDefaultID int64 `gorm:"column:monitor_metrics_default_id; UNIQUE_INDEX:uniqe_monitordevice_metric_MetricName" json:"monitor_metrics_default_id"` // 关联MonitorMetricsDefault 一对多,本表多
	//MonitorDevice         []*MonitorDevice `gorm:"many2many:itm_yw_monitor_device_metrics;"`
	//MonitorAlerts         []*MonitorAlerts  // 关联monitor_alerts  一对多,本表一
}

func (*MonitorMetrics) TableName() string {
	return "itm_yw_monitor_metrics"
}

func (m *MonitorMetrics) SaveMany(tx *gorm.DB) error {
	return tx.Create(m).Error
}
