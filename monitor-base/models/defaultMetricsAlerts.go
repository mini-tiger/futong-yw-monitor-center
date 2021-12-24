package models

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  defaultModelsData
 * @Version: 1.0.0
 * @Date: 2021/11/25 下午5:00
 */

func DefaultMonitorMetricsAlerts() []*MonitorAlertsDefault {
	data := make([]*MonitorAlertsDefault, 0)
	data = append(data, &MonitorAlertsDefault{
		ID:        1,
		Ident:     "cpu使用率",
		Level:     "warn",
		LevelName: "警告",
		Term:      ">",
		Value:     70,
		Expr:      "usedPercent{job=\"cpu\",hostid=\"%s\"} ",
		//MetricId:  1, // 与Metrics默认数据关联
	})

	data = append(data, &MonitorAlertsDefault{
		ID:        2,
		Ident:     "cpu使用率",
		Level:     "critical",
		LevelName: "严重",
		Term:      ">",
		Value:     85,
		Expr:      "usedPercent{job=\"cpu\",hostid=\"%s\"}",
		//MetricId:  1, // 与Metrics默认数据关联
	})

	data = append(data, &MonitorAlertsDefault{
		ID:        3,
		Ident:     "内存使用率",
		Level:     "warn",
		LevelName: "警告",
		Term:      ">",
		Value:     70,
		Expr:      "usedPercent{job=\"mem\",hostid=\"%s\"}",
		//MetricId:  2, // 与Metrics默认数据关联
	})

	data = append(data, &MonitorAlertsDefault{
		ID:        4,
		Ident:     "内存使用率",
		Level:     "critical",
		LevelName: "严重",
		Term:      ">",
		Value:     85,
		Expr:      "usedPercent{job=\"mem\",hostid=\"%s\"}",
		//MetricId:  2, // 与Metrics默认数据关联
	})

	data = append(data, &MonitorAlertsDefault{
		ID:        5,
		Ident:     "硬盘使用率",
		Level:     "warn",
		LevelName: "警告",
		Term:      ">",
		Value:     70,
		Expr:      "usedPercent{job=\"disk\",hostid=\"%s\"}",
		//MetricId:  3, // 与Metrics默认数据关联
	})

	data = append(data, &MonitorAlertsDefault{
		ID:        6,
		Ident:     "硬盘使用率",
		Level:     "critical",
		LevelName: "严重",
		Term:      ">",
		Value:     85,
		Expr:      "usedPercent{job=\"disk\",hostid=\"%s\"}",
		//MetricId:  3, // 与Metrics默认数据关联
	})

	data = append(data, &MonitorAlertsDefault{
		ID:        7,
		Ident:     "网络出流量",
		Level:     "warn",
		LevelName: "警告",
		Term:      ">",
		Value:     25 * 1024 * 1024,
		Expr:      "bytessentps{job=\"net\",hostid=\"%s\"}",
		//MetricId:  4, // 与Metrics默认数据关联
	})

	data = append(data, &MonitorAlertsDefault{
		ID:        8,
		Ident:     "网络出流量",
		Level:     "critical",
		LevelName: "严重",
		Term:      ">",
		Value:     35 * 1024 * 1024,
		Expr:      "bytessentps{job=\"net\",hostid=\"%s\"}",
		//MetricId:  4, // 与Metrics默认数据关联
	})

	data = append(data, &MonitorAlertsDefault{
		ID:        9,
		Ident:     "网络入流量",
		Level:     "warn",
		LevelName: "警告",
		Term:      ">",
		Value:     25 * 1024 * 1024,
		Expr:      "bytesrecvps{job=\"net\",hostid=\"%s\"}",
		//MetricId:  4, // 与Metrics默认数据关联
	})

	data = append(data, &MonitorAlertsDefault{
		ID:        10,
		Ident:     "网络入流量",
		Level:     "critical",
		LevelName: "严重",
		Term:      ">",
		Value:     35 * 1024 * 1024,
		Expr:      "bytesrecvps{job=\"net\",hostid=\"%s\"}",
		//MetricId:  4, // 与Metrics默认数据关联
	})
	return data
}
