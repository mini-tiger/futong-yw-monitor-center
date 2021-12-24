package collect

import (
	psutil_net "github.com/shirou/gopsutil/net"
	"net"
	"strings"
	"time"
)

type IOCountersStat struct {
	Name          string `json:"name"`          // interface name
	BytesSentPs   int    `json:"bytesSentPs"`   // number of bytes sent ps
	BytesRecvPs   int    `json:"bytesRecvPs"`   // number of bytes received ps
	PacketsSentPs int    `json:"packetsSentPs"` // number of packets sent ps
	PacketsRecvPs int    `json:"packetsRecvPs"` // number of packets received ps
}

// 采集网络相关指标
func CollectNetMetrics() ([]IOCountersStat, error) {
	var res []IOCountersStat

	b, err := psutil_net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	time.Sleep(5 * time.Second)

	a, err := psutil_net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(a); i++ {
		iOCountersStat := IOCountersStat{}
		iOCountersStat.Name = a[i].Name
		iOCountersStat.BytesRecvPs = int(a[i].BytesRecv-b[i].BytesRecv) / 5
		iOCountersStat.BytesSentPs = int(a[i].BytesSent-b[i].BytesSent) / 5
		iOCountersStat.PacketsSentPs = int(a[i].PacketsSent-b[i].PacketsSent) / 5
		iOCountersStat.PacketsRecvPs = int(a[i].PacketsRecv-b[i].PacketsRecv) / 5
		res = append(res, iOCountersStat)
	}

	return res, nil
}

// 基于 https://github.com/toolkits/net
// 采集ipList
func CollectIPList() ([]string, error) {
	ips := make([]string, 0)

	ifaces, err := net.Interfaces()
	if err != nil {
		return ips, err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		// ignore docker and warden bridge
		if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "w-") {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return ips, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			ips = append(ips, ip.String())
		}
	}
	return ips, nil
}
