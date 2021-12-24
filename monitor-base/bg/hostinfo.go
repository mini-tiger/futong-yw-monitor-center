package bg

import "github.com/shirou/gopsutil/host"

/**
 * @Author: Tao Jun
 * @Description: bg
 * @File:  hostinfo
 * @Version: 1.0.0
 * @Date: 2021/11/25 下午5:30
 */

type MonitorDeviceHostInfo struct {
	HostInfoFeature HostInfoFeature `json:"hostInfoFeature"`
	HostFeature     HostFeature     `json:"hostFeature"`
	//OutIp      string `json:"outIp"`
}
type HostInfoFeature struct {
	*host.InfoStat
	CpuBaseInfo
	//AgentId    string   `json:"agentId"`
	//ResourceID string   `json:"resourceId"`
	Ips      []string   `json:"ips"`
	Memory   uint64     `json:"memory"`
	DiskInfo []DiskInfo `json:"diskInfo"`
	IP       string     `json:"ip"`
}

type HostFeature struct {
	AgentVer AgentVersion `json:"agentVer"`
	Pattern  string       `json:"pattern"`
}

type CpuBaseInfo struct {
	ModelName string  `json:"modelName"`
	Cores     int     `json:"cores"`
	Threads   int     `json:"threads"`
	GHz       float64 `json:"GHz"`
}

type DiskInfo struct {
	Device      string  `json:"device"`
	MountPoint  string  `json:"mountpoint"`
	FsType      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}
