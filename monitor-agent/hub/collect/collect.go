package collect

import (
	"futong-yw-monitor-center/monitor-agent/g"
)

type HostMetrics struct {
	//AgentId      string      `json:"agentId"`
	//AgentVersion string      `json:"agentVersion"`
	Cpu          interface{} `json:"cpu"`
	Memory       interface{} `json:"memory"`
	Load         interface{} `json:"load"`
	Process      interface{} `json:"process"`
	Net          interface{} `json:"net"`
	Disk         interface{} `json:"disk"`
	DiskPhysical interface{} `json:"diskPhysical"`
	DiskRW       interface{} `json:"diskRW"`
	Options      *Options    `json:"options"`
}

// 采集监控信息
func CollectHostMetrics() *HostMetrics {
	var inodesUsedPercent float64

	// windows无此指标
	loadAvg, err := CollectAvgLoadMetrics()
	if err != nil {
		//logger.Sugar.Errorf("collect loadAvg err:", err)
	}

	cpuInfo, err := CollectCpuMetrics()
	if err != nil {
		g.GetLog().Error("collect cpu err:%v\n", err)
	}

	diskUsage, err := CollectDiskPartitionMetrics(true)
	if err != nil {
		// logger.Sugar.Errorf("collect disk err:", err)

	} else {
		inodesUsedPercent = diskUsage[0].InodesUsedPercent
	}
	// only linux
	diskPhysicalInfo, err := CollectDiskPhysicalInfo()
	if err != nil {
		g.GetLog().Error("collect disk physical partitions err:%v\n", err)
	}

	tcpCount, err := CollectTcpCount()
	if err != nil {
		g.GetLog().Error("collect tcpCount err:%v\n", err)
	}

	netStat, err := CollectNetMetrics()
	if err != nil {
		g.GetLog().Error("collect net err: %v\n", err)
	}

	memStat, err := CollectMemoryMetrics()
	if err != nil {
		g.GetLog().Error("collect memory err:%v\n", err)
	}

	processMetrics, err := CollectProcessMetrics()
	if err != nil {
		g.GetLog().Error("collect process err:%v\n", err)
	}

	diskRW, err := CollectDiskRW()
	if err != nil {
		g.GetLog().Error("collect diskRW err:%v\n", err)
	}

	return &HostMetrics{
		Cpu:          cpuInfo,
		Memory:       memStat,
		Load:         loadAvg,
		Process:      processMetrics,
		Net:          netStat,
		Disk:         diskUsage,
		DiskPhysical: diskPhysicalInfo,
		DiskRW:       diskRW,
		Options:      &Options{TcpCount: tcpCount, InodesUsedPercent: inodesUsedPercent},
	}
}
