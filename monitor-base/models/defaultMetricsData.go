package models

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  defaultModelsData
 * @Version: 1.0.0
 * @Date: 2021/11/25 下午5:00
 */

func DefaultMonitorMetricsData() []*MonitorMetricsDefault {
	data := make([]*MonitorMetricsDefault, 0)
	data = append(data, &MonitorMetricsDefault{
		ID:         1,
		MetricName: "cpu",
		Desc:       "cpu使用率",
		Ident:      "cpu",
	})
	data = append(data, &MonitorMetricsDefault{
		ID:         2,
		MetricName: "内存",
		Desc:       "内存使用率",
		Ident:      "mem",
	})
	data = append(data, &MonitorMetricsDefault{
		ID:         3,
		MetricName: "磁盘",
		Desc:       "磁盘使用率",
		Ident:      "disk",
	})
	data = append(data, &MonitorMetricsDefault{
		ID:         4,
		MetricName: "网络",
		Desc:       "网络流量",
		Ident:      "net",
	})
	data = append(data, &MonitorMetricsDefault{
		ID:         5,
		MetricName: "硬盘速率",
		Desc:       "硬盘速率",
		Ident:      "diskrw",
	})
	return data
}
