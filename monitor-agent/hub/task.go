package hub

import (
	"encoding/json"
	"fmt"
	"futong-yw-monitor-center/monitor-agent/hub/collect"
	"futong-yw-monitor-center/monitor-center/g"
)

//var hostAddr string = conf.Config.MustValue("", "serverAddr", "")

type MonitorMetrics struct {
	HostInfo    *collect.HostInfo    `json:"hostInfo"`
	HostMetrics *collect.HostMetrics `json:"hostMetrics"`
}

func sendHostMetrics() {
	hostinfo, err := collect.CollectHostInfo()
	if err != nil {
		g.GetLog().Error(err)
	}
	monitorMetrics := MonitorMetrics{
		HostInfo:    hostinfo,
		HostMetrics: collect.CollectHostMetrics(),
	}

	data, _ := json.Marshal(monitorMetrics)
	fmt.Println(string(data))

	//var (
	//	collectTime int
	//	ticker      *time.Ticker
	//)
	//collectTime = beego.AppConfig.DefaultInt("collect::interval", 15)
	//collectTime = conf.Config.MustInt("collect", "interval", 10)

	//for {
	//	ticker = time.NewTicker(time.Duration(collectTime) * time.Second)
	//	for  range ticker.C {
	//		hostMetrics, _ := json.Marshal(collect.CollectHostMetrics())
	//		fmt.Println(string(hostMetrics))
	//
	//		//postMetricsUrl := hostAddr + "/api/v1/yw/agent/monitor/"
	//
	//		//body, err := HttpPost(postMetricsUrl, hostMetrics)
	//		//if err != nil {
	//		//	logger.Sugar.Error("task sendHostMetrics err:", err)
	//		//	continue
	//		//}
	//
	//		//interval, err := strconv.Atoi(string(body))
	//		//if interval != collectTime && err == nil && interval >= utils.CollectTimeAllowedMin {
	//		//	collectTime = interval
	//		//	ticker.Stop()
	//		//	logger.Sugar.Infof("采集频率变为:%ds", interval)
	//		//	break
	//		//}
	//	}
	//}
}
