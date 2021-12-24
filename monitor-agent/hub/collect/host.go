package collect

import (
	"errors"
	"fmt"
	"futong-yw-monitor-center/monitor-agent/g"
	"futong-yw-monitor-center/monitor-base/bg"
	"github.com/shirou/gopsutil/host"
)

type HostInfo struct {
	*host.InfoStat
	bg.CpuBaseInfo
	AgentId    string   `json:"agentId"`
	ResourceID string   `json:"resourceId"`
	Ips        []string `json:"ips"`
	Memory     uint64   `json:"memory"`
}

// 采集主机信息
// BootTime returns the system boot time expressed in seconds since the epoch.
func CollectHostInfo() (*HostInfo, error) {
	hostinfo := &HostInfo{}
	hostStat, err := host.Info()
	if err != nil {
		//g.GetLog().Error("collect hostStat err: %v\n", err)
		return hostinfo, errors.New(fmt.Sprintf("collect hostStat err: %v\n", err))
	}
	g.HostID = hostStat.HostID

	ips, err := CollectIPList()
	if err != nil {
		//g.GetLog().Error("collect ips err:%v\n", err)
		return hostinfo, errors.New(fmt.Sprintf("collect ips err:%v\n", err))
	}

	threads, err := CollectCpuThreads()
	if err != nil {
		//g.GetLog().Error("collect cpu threads err:%v\n", err)
		return hostinfo, errors.New(fmt.Sprintf("collect cpu threads err:%v\n", err))
	}
	cores, err := CollectCpuCors()
	if err != nil {
		//g.GetLog().Error("collect cpu cores err:%v\n", err)
		return hostinfo, errors.New(fmt.Sprintf("collect cpu cores err:%v\n", err))
	}

	cpuInfoStats, err := CollectCpuBaseInfo()
	if err != nil {
		//g.GetLog().Error("collect cpu infostats err:%v\n", err)
		return hostinfo, errors.New(fmt.Sprintf("collect cpu infostats err:%v\n", err))
	}

	memoryStat, err := CollectMemoryMetrics()
	if err != nil {
		//g.GetLog().Error()
		return hostinfo, errors.New(fmt.Sprintf("collect memory info err: %v\n", err))
	}

	//resourceId := conf.Config.MustValue("", "resourceId")

	hostinfo = &HostInfo{
		//AgentId:  agent_id.AgentId,
		//ResourceID: resourceId,
		//ResourceID: "",
		InfoStat: hostStat,
		Ips:      ips,
		Memory:   memoryStat.Total,
		CpuBaseInfo: bg.CpuBaseInfo{
			ModelName: cpuInfoStats[0].ModelName,
			Cores:     cores,
			Threads:   threads,
			GHz:       cpuInfoStats[0].Mhz / 1000,
		},
	}
	return hostinfo, nil
}

func CollectBasicHostInfo() {
	var err error
	g.HostBaseInfo, err = host.Info()
	if err != nil {
		g.GetLog().Error("collect hostStat err: %v\n", err)
	}
	g.HostID = g.HostBaseInfo.HostID

}

func CollectMonitorDeviceHostInfo() (*bg.MonitorDeviceHostInfo, error) {
	monitorDevHost := &bg.MonitorDeviceHostInfo{}
	ips, err := CollectIPList()
	if err != nil {
		//g.GetLog().Error()
		return monitorDevHost, errors.New(fmt.Sprintf("collect ips err:%v\n", err))
	}

	threads, err := CollectCpuThreads()
	if err != nil {
		//g.GetLog().Error("collect cpu threads err:%v\n", err)
		return monitorDevHost, errors.New(fmt.Sprintf("collect cpu threads err:%v\n", err))
	}
	cores, err := CollectCpuCors()
	if err != nil {
		//g.GetLog().Error("collect cpu cores err:%v\n", err)
		return monitorDevHost, errors.New(fmt.Sprintf("collect cpu cores err:%v\n", err))
	}

	cpuInfoStats, err := CollectCpuBaseInfo()
	if err != nil {
		//g.GetLog().Error("collect cpu infostats err:%v\n", err)
		return monitorDevHost, errors.New(fmt.Sprintf("collect cpu infostats err:%v\n", err))

	}

	memoryStat, err := CollectMemoryMetrics()
	if err != nil {
		//g.GetLog().Error("collect memory info err: %v\n", err)
		return monitorDevHost, errors.New(fmt.Sprintf("collect memory info err: %v\n", err))
	}

	distList, err := CollectDiskphysicalPartitionMetrics()
	if err != nil {
		//g.GetLog().Error("collect disk info err: %v\n", err)
		return monitorDevHost, errors.New(fmt.Sprintf("collect disk info err: %v\n", err))
	}

	if g.OutIP == "" {
		url := ""
		if bg.Pattern != "ssh" {
			url = g.ViperCfg.ConfWeb.GetConfUrl
		} else {
			url = bg.ConfigUrl
		}
		if err := GetOutBandIPHandle(url); err != nil {
			return monitorDevHost, err
		}

	}

	monitorDevHost = &bg.MonitorDeviceHostInfo{
		//AgentId:  agent_id.AgentId,
		//ResourceID: resourceId,
		//ResourceID: "",
		HostInfoFeature: bg.HostInfoFeature{
			InfoStat: g.HostBaseInfo,
			CpuBaseInfo: bg.CpuBaseInfo{
				ModelName: cpuInfoStats[0].ModelName,
				Cores:     cores,
				Threads:   threads,
				GHz:       cpuInfoStats[0].Mhz / 1000,
			},
			Ips:      ips,
			Memory:   memoryStat.Total,
			DiskInfo: distList,
			IP:       g.OutIP,
		},
		HostFeature: bg.HostFeature{
			AgentVer: bg.AgentVer,
			Pattern:  bg.Pattern,
		},
	}
	return monitorDevHost, nil
}
