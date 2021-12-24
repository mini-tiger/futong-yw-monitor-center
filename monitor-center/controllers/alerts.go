package controllers

import (
	"fmt"
	"futong-yw-monitor-center/monitor-base/bg"
	bgmodels "futong-yw-monitor-center/monitor-base/models"
	"futong-yw-monitor-center/monitor-center/funcs"
	"futong-yw-monitor-center/monitor-center/g"
	"futong-yw-monitor-center/monitor-center/models"
	"github.com/gin-gonic/gin"
)

/**
 * @Author: Tao Jun
 * @Description: controllers
 * @File:  monitor
 * @Version: 1.0.0
 * @Date: 2021/11/25 下午4:24
 */

type AlertsBindReqParams struct {
	DefaultAlerts    *bgmodels.MonitorAlertsDefault `json:"defaultmetrics"`
	MonitorDeviceIDS []int64                        `json:"monitor_device_ids"`
}

func AlertsBind(c *gin.Context) {
	req := new(AlertsBindReqParams)
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

	defaultalerts := req.DefaultAlerts
	//models.Db.Take(&defaultMetric)
	//fmt.Println(defaultMetric)

	//metric.MonitorMetricsDefault = defaultMetric
	alerts := bgmodels.MonitorAlerts{}
	alerts.Expr = defaultalerts.Expr
	alerts.Ident = defaultalerts.Ident
	alerts.Value = defaultalerts.Value
	alerts.Term = defaultalerts.Term
	alerts.Level = defaultalerts.Level
	alerts.LevelName = defaultalerts.LevelName

	tx := models.Db.Begin()
	for _, monitorDeviceid := range req.MonitorDeviceIDS {
		atemp := alerts
		atemp.MonitorDeviceID = monitorDeviceid
		fmt.Printf("%+v\n", atemp)
		//同一个hostid ident相同 只能有一个，没检测
		err = tx.Create(&atemp).Error
		if err != nil {
			g.GetLog().Error("tx Create Record err:%v\n", err)
			continue
		}
		go func() {
			monitorDevice := new(bgmodels.MonitorDevice)
			models.Db.Where("id = ?", monitorDeviceid).Take(monitorDevice)
			err = monitorDevice.GenerateAlertRulesFile(models.Db)
			if err != nil {
				g.GetLog().Error("ip:%s GenerateAlertRulesFile err:%v\n", monitorDevice.Ip, err)
				return
			}
		}()
	}
	if err := tx.Commit().Error; err != nil {
		g.GetLog().Error("http path:%s Create Record err:%v\n", c.Request.URL, err)
		ResponseErrorMsg(c, err)
		return
	}

	funcs.PrometheusReload()

	ResponseSuccess(c, "ok")
	return
}

func AlertsAdd(c *gin.Context) {

	monitorAlert := new(bgmodels.MonitorAlertsDefault)
	err := c.BindJSON(monitorAlert)
	if err != nil {
		//errStr = fmt.Sprintf("Pattern %s ; err:%v", h.HostFeature.Pattern, err)
		g.GetLog().Error("http path:%s BindJson err:%v\n", c.Request.URL, err)
		ResponseErrorMsg(c, err)
		return
	}
	if err = bg.Validate.Struct(monitorAlert); err != nil {
		g.GetLog().Error("http path:%s Validate.Struct err:%v\n", c.Request.URL, err)
		ResponseErrorMsg(c, err)
		return
	}

	if monitorAlert.ID == 0 {
		// new
		expr := monitorAlert.Expr
		expr = fmt.Sprintf("%s{job=\"%s\",hostid=\"%%s\"}", expr, "shell")
		monitorAlert.Expr = expr
		fmt.Printf("%+v\n", monitorAlert)
		if err = models.Db.Create(monitorAlert).Error; err != nil {
			g.GetLog().Error("http path:%s Create record err:%v\n", c.Request.URL, err)
			ResponseErrorMsg(c, err)
			return
		}
	}

	ResponseSuccess(c, "ok")
	return
	//fmt.Printf("%+v\n",dd.MonitorMetricsDefault)

}
