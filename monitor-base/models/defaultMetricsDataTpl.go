package models

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  defaultModelsData
 * @Version: 1.0.0
 * @Date: 2021/11/25 下午5:00
 */

func DefaultMonitorMetricsTplData() *MonitorMetricsTpl {
	tpl := new(MonitorMetricsTpl)
	tpl.ID = 1
	tpl.Desc = "default monitorTpl"
	tpl.Name = "默认指标模板"
	tpl.MonitorMetrics = DefaultMonitorMetricsData()
	return tpl
}
