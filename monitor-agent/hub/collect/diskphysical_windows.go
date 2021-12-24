package collect

//   "github.com/StackExchange/wmi"
import (
	"github.com/StackExchange/wmi"
	"regexp"
	"runtime"
	"strings"
)

var rex, _ = regexp.Compile(`#\d`)

type Abc struct {
	Name string
	Size int
}

type Abc2 struct {
	Antecedent string
	Dependent  string
}

type Abc3 struct {
	Name string
	Size int
}

func CollectDiskPhysicalInfo() ([]map[string]interface{}, error) {
	var res = make([]map[string]interface{}, 0)

	if runtime.GOOS != "windows" {
		return res, nil
	}

	str1 := "Select * from Win32_DiskDrive"
	str2 := "Select * from Win32_DiskDriveToDiskPartition"
	str3 := "Select * from Win32_DiskPartition"

	//fmt.Println("********")
	//fmt.Println("语句1的输出：")
	var dst1 []Abc
	err := wmi.Query(str1, &dst1)
	if err != nil {
		return res, err
	}
	//for _, v := range dst1 {
	//	fmt.Println(v)
	//}

	//fmt.Println("********")
	//fmt.Println("语句2的输出：")
	var dst2 []Abc2
	err = wmi.Query(str2, &dst2)
	if err != nil {
		return res, err
	}
	//for _, v := range dst2 {
	//	fmt.Println(v)
	//}

	//fmt.Println("********")
	//fmt.Println("语句3的输出：")
	var dst3 []Abc3
	err = wmi.Query(str3, &dst3)
	if err != nil {
		return res, err
	}
	//for _, v := range dst3 {
	//	fmt.Println(v)
	//}

	//fmt.Println(">>>>>>>>>>>>")
	//fmt.Println("汇总输出：")
	for _, itm1 := range dst1 {
		name := itm1.Name
		name1 := name[10:]
		totalSize := itm1.Size
		phyUsed := 0
		for _, itm2 := range dst2 {
			ant := itm2.Antecedent
			bnt := itm2.Dependent
			for _, itm3 := range dst3 {
				//if strings.HasPrefix(itm3.Name, "磁盘") {
				//	itm3.Name = strings.Replace(itm3.Name, "磁盘", "Disk",8 )
				//	itm3.Name = strings.Replace(itm3.Name, "分区", "Partition",8 )
				//}
				regBnt := rex.FindAllString(bnt, 2)
				regName := rex.FindAllString(itm3.Name, 2)
				if len(regBnt) != 2 && len(regName) != 2 {
					continue
				}

				if strings.Contains(ant, name1) && regBnt[0] == regName[0] && regBnt[1] == regName[1] {
					phyUsed += itm3.Size
				}

				//if strings.Contains(bnt, itm3.Name) && strings.Contains(ant, name1) {
				//	phyUsed += itm3.Size
				//}
			}
		}
		res = append(res, map[string]interface{}{
			"name":      name,
			"totalSize": totalSize,
			"phyUsed":   phyUsed,
		})
	}
	return res, nil
}
