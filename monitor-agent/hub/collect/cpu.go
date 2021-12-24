package collect

import (
	"fmt"
	bgmodels "futong-yw-monitor-center/monitor-base/models"
	"github.com/shirou/gopsutil/cpu"
	"strconv"
	"time"
)

// 采集cpu线程数
func CollectCpuCors() (int, error) {
	return cpu.Counts(false)
}

// 采集cpu核数
func CollectCpuThreads() (int, error) {
	return cpu.Counts(true)
}

func CollectCpuBaseInfo() ([]cpu.InfoStat, error) {
	return cpu.Info()
}

// 采集cpu指标
func CollectCpuMetrics() (*bgmodels.CpuInfo, error) {
	var cpuInfo *bgmodels.CpuInfo

	cpuTimesStat, err := collectCpuTimes()
	if err != nil {
		return cpuInfo, err
	}

	cpuPercent, err := collectCpuPercent()
	if err != nil {
		return cpuInfo, err
	}

	cpuInfo = &bgmodels.CpuInfo{&cpuTimesStat[0], cpuPercent}

	return cpuInfo, nil
}

// 采集cpu使用率
func collectCpuPercent() (float64, error) {
	_, _ = cpu.Percent(0, false)
	time.Sleep(time.Second)
	res, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(fmt.Sprintf("%.2f", res[0]), 64)
}

// 采集cpu时间
func collectCpuTimes() ([]cpu.TimesStat, error) {
	return cpu.Times(false)
}
