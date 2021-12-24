package models

import "github.com/shirou/gopsutil/cpu"

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  pushgateway
 * @Version: 1.0.0
 * @Date: 2021/11/29 上午10:20
 */

type CpuInfo struct {
	TimesStat   *cpu.TimesStat `json:"timesStat"`
	UsedPercent float64        `json:"usedPercent"`
}
