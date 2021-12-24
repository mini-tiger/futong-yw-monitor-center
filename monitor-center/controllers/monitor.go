package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"futong-yw-monitor-center/monitor-base/bg"
	bgfuncs "futong-yw-monitor-center/monitor-base/funcs"
	bgmodels "futong-yw-monitor-center/monitor-base/models"
	"futong-yw-monitor-center/monitor-center/funcs"
	"futong-yw-monitor-center/monitor-center/g"
	"futong-yw-monitor-center/monitor-center/models"
	"futong-yw-monitor-center/monitor-center/utils"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"strconv"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: controllers
 * @File:  monitor
 * @Version: 1.0.0
 * @Date: 2021/11/25 下午4:24
 */

func ExcelUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	resp := new(models.MonitorExcelResp)
	if err != nil {
		resp.Err = err
		ResponseErrorMsg(c, resp)
		return
	}

	var authid int
	auth, b := c.GetPostForm("auth")
	if b {
		if authid, err = strconv.Atoi(auth); err != nil {
			g.GetLog().Error("authid is format err\n")
			resp.Err = errors.New(fmt.Sprintf("authid is format err"))
			ResponseErrorMsg(c, resp)
			return
		}
	}
	//check auth
	monitorAuth := bgmodels.MonitorAuth{}

	if err = monitorAuth.MustRecordMap(models.Db, map[string]interface{}{"id": authid}); err != nil {
		g.GetLog().Error("Find authid:%d err:%v\n", authid, err)
		resp.Err = errors.New(fmt.Sprintf(fmt.Sprintf("Find authid:%d err:%v", authid, err)))
		ResponseErrorMsg(c, resp)
		return
	}

	//save excel
	fileName := fmt.Sprintf("%d.xlsx", time.Now().Unix())
	tmpdir := path.Join(g.CurrentDir, "tmp")
	os.MkdirAll(tmpdir, 0644)
	savefile := path.Join(tmpdir, fileName)
	if err := c.SaveUploadedFile(file, savefile); err != nil {
		g.GetLog().Error("saveUploadFile err:%v\n", err)
		resp.Err = err
		ResponseErrorMsg(c, resp)

		return
	}

	// format excel
	exData, errData, err := utils.MonitorExcelRead(savefile)
	fmt.Printf("%+v\n", exData)
	fmt.Println(errData)
	fmt.Println(err)
	go func() {
		os.Remove(savefile)
	}()
	if err != nil {
		g.GetLog().Error("ExcelRead err:%v\n", err)
		resp.ExcelParse = errors.New(fmt.Sprintf("ExcelRead err:%v\n", err))
		ResponseErrorMsg(c, resp)
		return
	}

	success, inserterr, uniquerr, err := new(bgmodels.MonitorDevice).Excel2Model(models.Db, exData, &monitorAuth)

	resp.SuccessDB = success
	resp.Err = err
	resp.InsertErrs = inserterr
	resp.UniqErrs = uniquerr

	if err != nil {
		ResponseSuccess(c, resp)
	}
	ResponseSuccess(c, resp)
}

func DeleteDevice(c *gin.Context) {
	var monitorDevice *bgmodels.MonitorDevice = &bgmodels.MonitorDevice{}
	var errStr string
	err := c.BindJSON(monitorDevice)
	if err != nil {
		errStr = fmt.Sprintf("BindJson err:%v", err)
		g.GetLog().Error("http path:%s err:%v\n", c.Request.URL, errStr)
		ResponseErrorMsg(c, errStr)
		return
	}
	if monitorDevice.ID == 0 {
		errStr = fmt.Sprintf("ID is 0")
		g.GetLog().Error("http path:%s err:%v\n", c.Request.URL, errStr)
		ResponseErrorMsg(c, errStr)
		return
	}

	mustOneRecord := bgfuncs.GetDbOneRecord{}
	mustOneRecord.Result = monitorDevice
	mustOneRecord.Preload = []string{"MonitorMetrics", "MonitorAlerts"}
	mustOneRecord.Params = map[string]interface{}{"id": monitorDevice.ID}

	err = mustOneRecord.MustOneRecord(models.Db)

	if err != nil {
		errStr = fmt.Sprintf("%v", err)
		g.GetLog().Error("http path:%s err:%v\n", c.Request.URL, errStr)
		ResponseError(c, err)
		return
	}
	//fmt.Printf("%+v\n",monitorDevice)
	//models.Db.Debug().Unscoped().Delete(monitorDevice.MonitorAlerts)
	//models.Db.Debug().Unscoped().Delete(monitorDevice.MonitorMetrics)
	go func() {
		monitorDevice.DelMetricsAndAlertsData(models.Db)
	}()

	if err = models.Db.Delete(monitorDevice).Error; err != nil {
		errStr = fmt.Sprintf("id:%d delete err:%v", monitorDevice.ID, err)
		g.GetLog().Error("http path:%s err:%v\n", c.Request.URL, errStr)
		ResponseError(c, err)
	} else {
		ResponseSuccess(c, fmt.Sprintf("id:%d delete success", monitorDevice.ID))
	}
}

type MetricsBindReqParams struct {
	DefaultMetrics   *bgmodels.MonitorMetricsDefault `json:"defaultmetrics"`
	MonitorDeviceIDS []int64                         `json:"monitor_device_ids"`
}

func MetricsBind(c *gin.Context) {

	req := new(MetricsBindReqParams)
	err := c.BindJSON(req)
	if err != nil {
		//errStr = fmt.Sprintf("Pattern %s ; err:%v", h.HostFeature.Pattern, err)
		g.GetLog().Error("http path:%s BindJson err:%v\n", c.Request.URL, err)
		ResponseErrorMsg(c, err)
		return
	}
	if len(req.MonitorDeviceIDS) == 0 {
		g.GetLog().Error("http path:%s MonitorDeviceID len is 0\n", c.Request.URL)
		ResponseErrorMsg(c, "MonitorDeviceID len is 0")
		return
	}

	defaultMetric := req.DefaultMetrics
	//models.Db.Take(&defaultMetric)
	//fmt.Println(defaultMetric)

	//metric.MonitorMetricsDefault = defaultMetric
	metric := bgmodels.MonitorMetrics{}
	metric.MonitorMetricsDefaultID = int64(defaultMetric.ID)
	metric.MetricName = defaultMetric.MetricName
	metric.Ident = defaultMetric.Ident
	metric.Desc = defaultMetric.Desc

	tx := models.Db.Begin()
	for _, monitorDeviceid := range req.MonitorDeviceIDS {
		mtemp := metric
		mtemp.MonitorDeviceID = monitorDeviceid
		err = mtemp.SaveMany(tx)
		if err != nil {
			g.GetLog().Error("tx Create Record err:%v\n", err)
			continue
		}
	}
	if err := tx.Commit().Error; err != nil {
		g.GetLog().Error("http path:%s Create Record err:%v\n", c.Request.URL, err)
		ResponseErrorMsg(c, err)
		return
	}

	ResponseSuccess(c, "ok")
	return
	//fmt.Printf("%+v\n",dd.MonitorMetricsDefault)

}

func FirstInitMonitorDevice(c *gin.Context) {
	var h bg.MonitorDeviceHostInfo
	var errStr string
	err := c.BindJSON(&h)
	if err != nil {
		errStr = fmt.Sprintf("Pattern %s ; err:%v", h.HostFeature.Pattern, err)
		g.GetLog().Error("http path:%s BindJson err:%v\n", c.Request.URL, errStr)
		ResponseErrorMsg(c, errStr)
		return
	}
	//fmt.Printf("%+v\n",h.DiskInfo)
	monitorDevice := &bgmodels.MonitorDevice{}

	//每次 update  or  Insert
	// 不能确定agent 硬件信息什么时候更新，agent启动上传一次硬件

	mustOneRecord := bgfuncs.GetDbOneRecord{}
	mustOneRecord.Result = monitorDevice
	mustOneRecord.Preload = []string{"MonitorMetrics", "MonitorAlerts"}
	mustOneRecord.Params = map[string]interface{}{"hostid": h.HostInfoFeature.HostID}
	//err = monitorDevice.GetMap(models.Db, map[string]interface{}{"hostid": h.HostInfoFeature.HostID})
	//if err != nil {
	//	g.GetLog().Error("http path:%s err:%v\n", c.Request.URL, err)
	//	ResponseError(c, err)
	//	return
	//}
	err = mustOneRecord.MustOneRecord(models.Db)
	if err != nil && err.Error() != "record not found" {
		errStr = fmt.Sprintf("Pattern %s,err:%v", h.HostFeature.Pattern, err)
		g.GetLog().Error("http path:%s err:%v\n", c.Request.URL, errStr)
		ResponseError(c, err)
		return
	}

	var sshFirstUpload bool = false
	// ssh 方式 第一次上传, ssh第一次上传查不到hostid
	//web 和 agent方式 不需要提前在表中创建记录
	if monitorDevice.ID == 0 && h.HostFeature.Pattern == "ssh" {
		//monitorDevice = &bgmodels.MonitorDevice{}
		//mustOneRecord = bgfuncs.GetDbOneRecord{}
		//mustOneRecord.Result = monitorDevice
		//mustOneRecord.Params = map[string]interface{}{"hostid": h.HostInfoFeature.HostID}
		//xxx 只改变 查询参数
		mustOneRecord.Params = map[string]interface{}{"ip": h.HostInfoFeature.IP}

		err = mustOneRecord.MustOneRecord(models.Db)
		if err != nil {
			errStr = fmt.Sprintf("Pattern %s,err:%v", h.HostFeature.Pattern, err)
			g.GetLog().Error("http path:%s err:%v\n", c.Request.URL, errStr)
			ResponseErrorMsg(c, errStr)
			return
		}
		// 查找后确认ssh第一次上传
		if len(monitorDevice.HostId) <= 1 {
			sshFirstUpload = true
		}
	}

	monitorDevice.Status = 1
	monitorDevice.HostName = h.HostInfoFeature.Hostname
	monitorDevice.OsType = h.HostInfoFeature.OS
	monitorDevice.Arch = h.HostInfoFeature.KernelArch
	monitorDevice.HostId = h.HostInfoFeature.HostID
	if h.HostInfoFeature.IP != "" {
		monitorDevice.Ip = h.HostInfoFeature.IP
	} else {
		if len(h.HostInfoFeature.Ips) > 0 {
			monitorDevice.Ip = h.HostInfoFeature.Ips[0]
		}
	}
	monitorDevice.AgentVerCurrent = uint64(h.HostFeature.AgentVer)
	monitorDevice.Pattern = h.HostFeature.Pattern

	b, _ := json.Marshal(&h.HostInfoFeature)
	monitorDevice.HostInfo = string(b)

	if monitorDevice.ID == 0 {
		monitorDevice.LastReqCfg = uint64(time.Now().Unix())
		//添加默认监控和报警项
		//monitorDevice.MonitorAlerts = bgmodels.CopyDefaultAlertData(models.Db)
		//monitorDevice.MonitorMetrics = bgmodels.CopyDefaultMetricsData(models.Db)

		if err = models.Db.Create(monitorDevice).Error; err != nil {
			errStr = fmt.Sprintf("Pattern %s,err:%v", h.HostFeature.Pattern, fmt.Sprintf("DB Create MonitorDevice err:%v,data:%v\n",
				err, monitorDevice))
			g.GetLog().Error("http path:%s err:%v\n", c.Request.URL, errStr)
			ResponseErrorMsg(c, errStr)
			return
		}
		models.Db.Model(monitorDevice).Association("MonitorAlerts").
			Replace(bgmodels.CopyDefaultAlertData(models.Db))

		models.Db.Model(monitorDevice).Association("MonitorMetrics").
			Replace(bgmodels.CopyDefaultMetricsData(models.Db))

		go func() {
			err = monitorDevice.GenerateAlertRulesFile(models.Db)
			if err != nil {
				g.GetLog().Error("ip:%s GenerateAlertRulesFile err:%v\n", monitorDevice.Ip, err)
				return
			}
			funcs.PrometheusReload()
		}()
	} else {
		// 可能是存在pattern web or agent的记录 也可能是ssh 第一次上传

		if monitorDevice.Pattern == "ssh" && sshFirstUpload {
			//添加默认监控和报警项
			//monitorDevice.MonitorAlerts = bgmodels.CopyDefaultAlertData(models.Db)
			//monitorDevice.MonitorMetrics = bgmodels.CopyDefaultMetricsData(models.Db)
			if err = models.Db.Save(monitorDevice).Error; err != nil {
				errStr = fmt.Sprintf("Pattern %s,err:%v", h.HostFeature.Pattern,
					fmt.Sprintf("DB Save MonitorDevice Association err:%v,data:%v",
						err, monitorDevice))
				g.GetLog().Error("http path:%s err:%v\n", c.Request.URL, errStr)
				ResponseErrorMsg(c, errStr)
				return
			}
			models.Db.Model(monitorDevice).Association("MonitorAlerts").
				Replace(bgmodels.CopyDefaultAlertData(models.Db))

			models.Db.Model(monitorDevice).Association("MonitorMetrics").
				Replace(bgmodels.CopyDefaultMetricsData(models.Db))
			go func() {
				err = monitorDevice.GenerateAlertRulesFile(models.Db)
				if err != nil {
					g.GetLog().Error("ip:%s GenerateAlertRulesFile err:%v\n", monitorDevice.Ip, err)
					return
				}
				funcs.PrometheusReload()
			}()
		} else {
			// hostinfo 更新
			if err = models.Db.Save(monitorDevice).Error; err != nil {
				errStr = fmt.Sprintf("Pattern %s,err:%v", h.HostFeature.Pattern,
					fmt.Sprintf("DB Save MonitorDevice err:%v,data:%v",
						err, monitorDevice))
				g.GetLog().Error("http path:%s err:%v\n", c.Request.URL, errStr)
				ResponseErrorMsg(c, errStr)

				return
			}
		}
	}

	ResponseSuccess(c, &h)
}
