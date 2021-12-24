package models

import "github.com/jinzhu/gorm"

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  agent
 * @Version: 1.0.0
 * @Date: 2021/11/25 上午11:20
 */
type MonitorAlertsDefault struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	Ident     string `gorm:"column:ident;" json:"ident" validate:"required"`      // 使用率
	Level     string `gorm:"column:level;" json:"level" validate:"required"`      // warn
	LevelName string `json:"levelname"`                                           // 警告
	Term      string `json:"term" validate:"required"`                            // 计算符号  < >
	Value     int64  `gorm:"column:value;" json:"value" validate:"required,gt=0"` // 80
	Expr      string `gorm:"column:expr;comment:'expr'" json:"expr" validate:"required"`
}

func (*MonitorAlertsDefault) TableName() string {
	return "itm_yw_monitor_alerts_default"
}

type MonitorMetricsDefault struct {
	ID             uint   `gorm:"primary_key" json:"id"`
	MetricName     string `gorm:"column:metricname; size:64" json:"metricname"`
	Desc           string `gorm:"column:desc;" json:"desc"`
	Ident          string `gorm:"column:ident;" json:"ident"`
	Path           string `gorm:"column:path;" json:"path"` //文件路径
	MonitorMetrics []*MonitorMetrics
}

func (*MonitorMetricsDefault) TableName() string {
	return "itm_yw_monitor_metrics_default"
}

func CopyDefaultAlertData(db *gorm.DB) []*MonitorAlerts {
	var defaultMonitorAlerts []*MonitorAlertsDefault = DefaultMonitorMetricsAlerts()
	//db.Find(&defaultMonitorAlerts)

	MonitorMetricsSlice := make([]*MonitorAlerts, len(defaultMonitorAlerts))

	for index, defaultmonitorAlert := range defaultMonitorAlerts {
		MonitorMetricsSlice[index] = &MonitorAlerts{
			//ID:              0,
			//UpdatedAt:       time.Time{},
			Ident:     defaultmonitorAlert.Ident,
			Level:     defaultmonitorAlert.Level,
			LevelName: defaultmonitorAlert.LevelName,
			Term:      defaultmonitorAlert.Term,
			Value:     defaultmonitorAlert.Value,
			Expr:      defaultmonitorAlert.Expr,
			//MonitorDeviceID: 0,
		}
	}
	return MonitorMetricsSlice
}

func CopyDefaultMetricsData(db *gorm.DB) []*MonitorMetrics {

	var defaultMonitorMetrics []*MonitorMetricsDefault = DefaultMonitorMetricsData()
	//db.Find(&defaultMonitorMetrics)

	MonitorMetricsSlice := make([]*MonitorMetrics, len(defaultMonitorMetrics))
	for index, monitorMetric := range defaultMonitorMetrics {
		MonitorMetricsSlice[index] = &MonitorMetrics{
			//ID:              0,
			//UpdatedAt:       time.Time{},
			MetricName: monitorMetric.MetricName,
			Desc:       monitorMetric.Desc,
			Ident:      monitorMetric.Ident,
			//MonitorDeviceID: 0,
		}
	}
	return MonitorMetricsSlice
}
