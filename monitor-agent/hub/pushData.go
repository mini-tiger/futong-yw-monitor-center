package hub

import (
	"fmt"
	"futong-yw-monitor-center/monitor-agent/g"
	"futong-yw-monitor-center/monitor-agent/hub/collect"
	"futong-yw-monitor-center/monitor-agent/utils"
	"futong-yw-monitor-center/monitor-base/bg"
	bgmodels "futong-yw-monitor-center/monitor-base/models"
	bgutils "futong-yw-monitor-center/monitor-base/utils"
	"github.com/fatih/structs"
	"github.com/goinggo/mapstructure"
	"github.com/shirou/gopsutil/disk"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

/**
 * @Author: Tao Jun
 * @Description: hub
 * @File:  pushGateway
 * @Version: 1.0.0
 * @Date: 2021/11/20 下午10:27
 */

type GetData struct {
	MonitorMetrics *MonitorMetrics
	MonitorShells  []*bgmodels.MonitorMetricsDefault
	Metrics        []string `json:"metrics"` //采集 的指标
	HostId         string
	IP             string
	Labels         map[string]interface{}
}

func (d *GetData) CollectData() {
	hostinfo, err := collect.CollectHostInfo()
	if err != nil {
		g.GetLog().Error(err)
	}

	d.MonitorMetrics = &MonitorMetrics{
		HostInfo:    hostinfo,
		HostMetrics: collect.CollectHostMetrics(),
	}

	d.HostId = d.MonitorMetrics.HostInfo.HostID

	if g.OutIP == "" {
		if len(d.MonitorMetrics.HostInfo.Ips) > 0 {
			d.IP = d.MonitorMetrics.HostInfo.Ips[0]
		}
	} else {
		d.IP = g.OutIP
	}
}

func (d *GetData) CovertShellMetrics(data interface{}) {
	dataSlice, ok := data.([]interface{})
	if !ok {
		d.MonitorShells = nil
		return
	}
	for _, value := range dataSlice {
		dataMap, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		var mmd bgmodels.MonitorMetricsDefault
		//将 map 转换为指定的结构体
		if err := mapstructure.Decode(dataMap, &mmd); err != nil {
			continue
		}
		d.MonitorShells = append(d.MonitorShells, &mmd)
	}
}

func (d *GetData) CollectShellData(baseUrl string) map[string]interface{} {
	shellMap := make(map[string]interface{}, 0)
	if len(d.MonitorShells) == 0 {
		return shellMap
	}
	baseShellPath := path.Join(g.CurrentPwd, "tmpShell")

	os.MkdirAll(baseShellPath, 0644)
	//todo 重复下载shell   too many open files
	for _, shellData := range d.MonitorShells {
		u, err := url.Parse(baseUrl)
		if err != nil {
			continue
		}
		u.Path = path.Join(u.Path, shellData.Path)
		shellUrl := u.String()
		_, fileName := path.Split(shellData.Path)
		shellFile := path.Join(baseShellPath, fileName)

		if err = utils.HttpDownFile(shellUrl, shellFile); err != nil {
			g.GetLog().Error("http down %s err:%v\n", shellUrl, err)
			continue
		}
		g.GetLog().Info("http down %s success\n", shellUrl)
		errout, out := bgutils.RunCommand(shellFile)
		if errout != "" {
			g.GetLog().Error("shell %s err:%v\n", shellData.Path, err)
			continue
		}
		var num int
		if num, err = strconv.Atoi(bgutils.TrimEnterSpace(out)); err == nil {
			shellMap[shellData.MetricName] = num
			g.GetLog().Debug("shell desc:%s, path:%s out:%v strconv success\n", shellData.Desc, shellData.Path, num)
		} else {
			g.GetLog().Error("shell desc:%s, path:%s out:%v strconv err:%v\n", shellData.Desc, shellData.Path, out, err)
		}

	}
	return shellMap
}

func (d *GetData) Run() {
	defer func() {
		if err := recover(); err != nil {
			g.GetLog().Error("Err Task Cron Err : %v\n", err)
		}
	}()

	d.CollectData()

	mc := new(bgmodels.CollectHost)
	mc.Ip = d.IP
	mc.HostId = d.HostId
	mc.Labels = d.Labels
	mc.OS = g.HostBaseInfo.OS
	mc.Arch = g.HostBaseInfo.KernelArch
	mc.PushGateAddr = bg.PushCfgEntry.PushGateWay
	mc.ESAddr = []string{bg.PushCfgEntry.Es.Host}
	//mc.Metrics = bg.PushCfgEntry.Metrics
	d.Metrics = bg.PushCfgEntry.Metrics
	//allMetricsData := d.FormatData()
	mc.HostMetrics = d.FormatDataFilterMetrics(d.Metrics)

	// shell
	g.GetLog().Debug("metrics shells %+v\n", bg.PushCfgEntry.Shells)

	d.MonitorShells = nil //重新
	d.CovertShellMetrics(bg.PushCfgEntry.Shells)
	var confurl string

	if bg.Pattern == "ssh" {
		confurl = bg.ConfigUrl
	} else {
		confurl = g.ViperCfg.ConfWeb.GetConfUrl
	}
	monitorCenterHost := bgutils.HttpHostSplit(confurl)

	if monitorCenterHost == "" {
		g.GetLog().Error("getConfUrl Host is null skip shell exec\n")
	} else {
		shellmap := d.CollectShellData(fmt.Sprintf("http://%s", monitorCenterHost))
		if len(shellmap) > 0 {
			mc.HostMetrics["shell"] = shellmap
		}
	}

	g.GetLog().Debug("hostid:%s,ip:%s metrics:%v\n", mc.HostId, mc.Ip,
		strings.Join(d.Metrics, ","))

	allpusherrs, allEserrs := mc.Push2EsPushGateWay()
	if len(allpusherrs) > 0 {
		for metric, e := range allpusherrs {
			g.GetLog().Error("ip:%s hostid:%s Push metrics:%s err:%v\n", d.IP, d.HostId, metric, e)
		}
	} else {
		g.GetLog().Info("IP:%s HostId:%s Push metrics success\n", d.IP, d.HostId)
	}

	if len(allEserrs) > 0 {
		for metric, e := range allEserrs {
			g.GetLog().Error("ip:%s hostid:%s Es metrics:%s err:%v\n", d.IP, d.HostId, metric, e)
		}
	} else {
		g.GetLog().Info("IP:%s HostId:%s Es metrics success\n", d.IP, d.HostId)
	}

}

func (d *GetData) FormatDataFilterMetrics(metrics []string) map[string]interface{} {
	formatData := make(map[string]interface{}, 0)
	for _, value := range metrics {
		metricName := value
		switch metricName {
		case "cpu":
			//d.PushCpuToPushGateway("cpu")

			cpuinfo := d.MonitorMetrics.HostMetrics.Cpu.(*bgmodels.CpuInfo)
			m := structs.Map(cpuinfo.TimesStat)
			m["usedPercent"] = cpuinfo.UsedPercent
			formatData[metricName] = m

			continue
		case "disk":
			diskinfo := d.MonitorMetrics.HostMetrics.Disk.([]*disk.UsageStat)
			alldata := make([]interface{}, 0)
			for _, di := range diskinfo {
				fstype := di.Fstype
				path := di.Path
				if di.Total > 0 && !strings.Contains(fstype, "tmpfs") &&
					!strings.Contains(path, "overlay2") {
					dm := bgutils.Struct2Map(*di, true)
					dm["usedPercent"] = dm["usedpercent"]
					delete(dm, "usedpercent")
					alldata = append(alldata, dm)
				}
			}
			formatData[metricName] = alldata
			continue
		case "diskrw":
			diskrw := d.MonitorMetrics.HostMetrics.DiskRW.(*collect.DiskRW)
			m := bgutils.Struct2Map(*diskrw, true)
			formatData[metricName] = m

			continue
		case "net":
			//	d.PushNetToPushGateway("net")
			alldata := make([]interface{}, 0)
			netinfo := d.MonitorMetrics.HostMetrics.Net.([]collect.IOCountersStat)
			for _, ni := range netinfo {
				//byte to Byte
				ni.BytesRecvPs = ni.BytesRecvPs * 8
				ni.BytesSentPs = ni.BytesSentPs * 8
				nm := bgutils.Struct2Map(ni, true)
				alldata = append(alldata, nm)
			}
			formatData[metricName] = alldata
			continue
		case "mem":
			meminfo := d.MonitorMetrics.HostMetrics.Memory.(*collect.MemoryStat)

			m := bgutils.Struct2Map(*meminfo, true)
			m["usedPercent"] = m["usedpercent"]
			delete(m, "usedpercent")
			formatData[metricName] = m
			//d.PushMemToPushGateway("mem")
			continue
		//case "load":
		//	if strings.ToLower(runtime.GOOS) == "linux" {
		//		d.PushLoadToPushGateway("load")
		//	}
		//	continue
		default:
			continue
		}
	}

	//d.PushHealthToPushGateway("up")
	return formatData
}
