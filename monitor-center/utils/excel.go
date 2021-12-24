package utils

import (
	"fmt"
	"futong-yw-monitor-center/monitor-base/bg"
	bgmodels "futong-yw-monitor-center/monitor-base/models"
	"github.com/go-playground/validator/v10"
	"github.com/szyhf/go-excel"
)

func MonitorExcelRead(exFile string) (map[int]bgmodels.MonitorDeviceExcelRow, map[int]interface{}, error) {

	exData := make(map[int]bgmodels.MonitorDeviceExcelRow, 0)
	errData := make(map[int]interface{}, 0)

	conn := excel.NewConnecter()
	err := conn.Open(exFile)
	if err != nil {
		return exData, errData, err
	}

	defer conn.Close()

	// Generate an new reader of a sheet
	// sheetNamer: if sheetNamer is string, will use sheet as sheet name.
	//             if sheetNamer is int, will i'th sheet in the workbook, be careful the hidden sheet is counted. i ∈ [1,+inf]
	//             if sheetNamer is a object implements `GetXLSXSheetName()string`, the return value will be used.
	//             otherwise, will use sheetNamer as struct and reflect for it's name.
	// 			   if sheetNamer is a slice, the type of element will be used to infer like before.

	cfg := &excel.Config{
		Sheet:         "monitor_devices",
		TitleRowIndex: 1, // xxx title 行的位置
		Skip:          0,
		Prefix:        "",
		Suffix:        "",
	}
	rd, err := conn.NewReaderByConfig(cfg)

	var index = 3
	for rd.Next() {
		//fmt.Printf("rd:%v\n",rd.GetTitles())
		var r bgmodels.MonitorDeviceExcelRow
		// Read a row into a struct.
		err := rd.Read(&r)
		if err != nil {
			if err.Error() != "EOF" {
				errData[index] = err.Error()
			}
			index++
			continue
		}
		//fmt.Println("============", index)

		//xxx 检验数据
		err = bg.Validate.Struct(&r)
		if err != nil {
			errslice, ok := err.(validator.ValidationErrors)
			if len(errslice) > 0 && ok {
				errData[index] = fmt.Sprintf("行: %d 列: %s 值: %s 不符合:%s 规范", index, (errslice[0]).Field(), errslice[0].Value(), errslice[0].Field())
				index++
				continue
			}
		}
		//fmt.Printf("%+v\n", r)
		exData[index] = r
		index++
	}
	return exData, errData, nil
}
