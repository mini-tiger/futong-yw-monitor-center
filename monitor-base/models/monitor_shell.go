package models

import "time"

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  agent
 * @Version: 1.0.0
 * @Date: 2021/11/25 上午11:20
 */
type MonitorShells struct {
	ID             uint `gorm:"primary_key" json:"id"`
	UpdatedAt      time.Time
	MetricName     string           `gorm:"column:metricname; UNIQUE_INDEX" json:"metricname" validate:"required"` //定义指标名称
	Desc           string           `gorm:"column:desc;" json:"desc"`
	Ident          string           `gorm:"column:ident; default:'shell';" json:"ident"` //
	Path           string           `gorm:"column:path;" json:"path"`                    //文件路径
	MonitorDevices []*MonitorDevice `gorm:"many2many:itm_yw_monitor_device_shell;"`
}

func (*MonitorShells) TableName() string {
	return "itm_yw_monitor_shells"
}
