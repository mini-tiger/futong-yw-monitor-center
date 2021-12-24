package collect

import (
	"fmt"
	"futong-yw-monitor-center/monitor-base/bg"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/disk"
)

// 采集disk physical 分区使用情况指标

func CollectDiskphysicalPartitionMetrics() ([]bg.DiskInfo, error) {
	partitionStats, err := disk.Partitions(false) //所有分区
	var diskUsageStatList []bg.DiskInfo
	if err != nil {
		return diskUsageStatList, nil
	}

	for _, partition := range partitionStats {
		// 采集各分区的使用情况
		us, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
			// return nil, err
		}
		UsedPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", us.UsedPercent), 64)
		diskUsageStatList = append(diskUsageStatList, bg.DiskInfo{
			Device:      partition.Device,
			MountPoint:  partition.Mountpoint,
			FsType:      partition.Fstype,
			Total:       us.Total,
			Used:        us.Used,
			UsedPercent: UsedPercent,
		})
	}
	return diskUsageStatList, nil

}

// 采集disk分区使用情况指标
func CollectDiskPartitionMetrics(all bool) ([]*disk.UsageStat, error) {
	partitionStats, err := collectDiskPartitions(all)
	if err != nil {
		return nil, err
	}

	var diskUsageStatList []*disk.UsageStat

	for _, partition := range partitionStats {
		// 采集各分区的使用情况
		us, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
			// return nil, err
		}

		us.UsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", us.UsedPercent), 64)
		us.InodesUsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", us.InodesUsedPercent), 64)

		diskUsageStatList = append(diskUsageStatList, us)
	}
	return diskUsageStatList, nil
}

// 采集disk分区 挂载点
func collectDiskPartitions(all bool) ([]disk.PartitionStat, error) {
	// Partitions returns disk partitions. If all is false, returns
	// physical devices only (e.g. hard disks, cd-rom drives, USB keys)
	// and ignore all others (e.g. memory partitions such as /dev/shm)
	return disk.Partitions(all)
}

type DiskRW struct {
	WriteBytesPs uint64 `json:"writeBytesPs"`
	ReadBytesPs  uint64 `json:"readBytesPs"`
	WriteCountPs uint64 `json:"writeCountPs"`
	ReadCountPs  uint64 `json:"readCountPs"`
}

// 磁盘读写速率
func CollectDiskRW() (*DiskRW, error) {
	var (
		firstReadBytesTotal  uint64
		firstWriteBytesTotal uint64
		firstReadCountTotal  uint64
		firstWriteCountTotal uint64

		secondReadBytesTotal  uint64
		secondWriteBytesTotal uint64
		secondReadCountTotal  uint64
		secondWriteCountTotal uint64
	)

	firstIOCounters, err := disk.IOCounters()
	if err != nil {
		return nil, err
	}

	for _, v := range firstIOCounters {
		firstReadBytesTotal += v.ReadBytes
		firstWriteBytesTotal += v.WriteBytes
		firstReadCountTotal += v.ReadCount + v.MergedReadCount
		firstWriteCountTotal += v.WriteCount + v.MergedWriteCount
	}

	time.Sleep(5 * time.Second)

	secondIOCounters, err := disk.IOCounters()
	if err != nil {
		return nil, err
	}

	for _, v := range secondIOCounters {
		secondReadBytesTotal += v.ReadBytes
		secondWriteBytesTotal += v.WriteBytes
		secondReadCountTotal += v.ReadCount + v.MergedReadCount
		secondWriteCountTotal += v.WriteCount + v.MergedWriteCount
	}

	return &DiskRW{
		WriteBytesPs: (secondWriteBytesTotal - firstWriteBytesTotal) / 5,
		ReadBytesPs:  (secondReadBytesTotal - firstReadBytesTotal) / 5,
		WriteCountPs: (secondReadCountTotal - firstReadCountTotal) / 5,
		ReadCountPs:  (secondWriteCountTotal - firstWriteCountTotal) / 5,
	}, nil
}
