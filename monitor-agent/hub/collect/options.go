package collect

import "github.com/shirou/gopsutil/net"

type Options struct {
	TcpCount          int     `json:"tcpCount"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}

// 采集tcp连接数
func CollectTcpCount() (int, error) {
	tcpCount, err := net.Connections("tcp")
	return len(tcpCount), err
}
