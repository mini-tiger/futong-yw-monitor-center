package models

import "time"

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  agent
 * @Version: 1.0.0
 * @Date: 2021/11/25 上午11:20
 */
type MonitorAlerts struct {
	ID        uint `gorm:"primary_key" json:"id"`
	UpdatedAt time.Time
	//MetricId      uint64           `gorm:"column:metric_id; size:64" json:"metric_id"` // 关联monitor_metrics  一对多,本表多
	Ident     string `gorm:"column:ident;" json:"ident" validate:"required"`      // 使用率
	Level     string `gorm:"column:level;" json:"level"`                          // warn
	LevelName string `json:"levelname;" json:"levelname" validate:"required"`     // 警告
	Term      string `json:"term" validate:"required"`                            // 计算符号  < >
	Value     int64  `gorm:"column:value;" json:"value" validate:"required,gt=0"` // 80
	Expr      string `gorm:"column:expr;comment:'expr'" json:"expr"`
	//MonitorDevice []*MonitorDevice `gorm:"many2many:itm_yw_monitor_device_alerts;"` //多对多
	MonitorDeviceID int64 `gorm:"column:monitor_device_id" json:"monitor_device_id" validate:"required"`
}

func (*MonitorAlerts) TableName() string {
	return "itm_yw_monitor_alerts"
}
