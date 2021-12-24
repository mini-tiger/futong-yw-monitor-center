package models

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  httpmodels
 * @Version: 1.0.0
 * @Date: 2021/12/3 下午3:55
 */

type MonitorExcelResp struct {
	SuccessDB  interface{}
	ExcelParse interface{}
	InsertErrs interface{}
	UniqErrs   interface{}
	Err        error
}
