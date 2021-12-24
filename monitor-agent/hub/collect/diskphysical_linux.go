package collect

import (
	"os/exec"
	"strconv"
	"strings"
	"unicode"

	diskHelper "github.com/shirou/gopsutil/disk"
)

func CollectDiskPhysicalInfo() ([]map[string]interface{}, error) {
	str1 := "export LANG=en_US.UTF-8"
	// str2 := `ls -l /dev/ | grep " 8," | awk '$NF ~ /[[:alpha:]]$/{print $NF}'`
	// str3 := `export LANG=en_US.UTF-8 && fdisk -l | grep "^Disk /dev/$disk:" | awk '{print $7}'`
	str3 := `export LANG=en_US.UTF-8 && fdisk -l | grep "^Disk $disk:" | awk '{print $7}'`
	// str4 := `partx /dev/$disk | awk '{if (NR>1) print $4}' | awk 'BEGIN{sum=2048}{sum+=$1}END{print sum}'`
	str4 := `partx $disk | awk '{if (NR>1) print $4}' | awk 'BEGIN{sum=2048}{sum+=$1}END{print sum}'`
	var res []map[string]interface{}

	if _, err := exec.Command("/bin/bash", "-c", str1).Output(); err != nil {
		return res, nil
	}

	// data2, err := exec.Command("/bin/bash", "-c", str2).Output()
	// if err != nil {
	// 	return res, err
	// }
	// diskList := strings.Split(string(data2), "\n")
	partitions, err := diskHelper.Partitions(false)
	if err != nil {
		return res, err
	}
	diskMap := make(map[string]struct{}, 0)
	for _, v := range partitions {
		temp := ""
		for _, s := range v.Device {
			if !unicode.IsNumber(s) {
				temp += string(s)
			}
		}
		diskMap[temp] = struct{}{}
	}

	for disk, _ := range diskMap {
		if disk == "" {
			continue
		}
		str3 := strings.Replace(str3, "$disk", disk, 4)
		str4 := strings.Replace(str4, "$disk", disk, 4)
		data3, err := exec.Command("/bin/bash", "-c", str3).Output()
		if err != nil {
			continue
		}
		data4, err := exec.Command("/bin/bash", "-c", str4).Output()
		if err != nil {
			continue
		}

		totalSize, err := strconv.ParseInt(strings.TrimSpace(string(data3)), 10, 64)
		if err != nil {
			continue
		}

		phyUsed, err := strconv.ParseInt(strings.TrimSpace(string(data4)), 10, 64)
		if err != nil {
			continue
		}

		res = append(res, map[string]interface{}{
			"name":      disk,
			"totalSize": totalSize * 512,
			"phyUsed":   phyUsed * 512,
		})
	}

	return res, nil
}
