package controllers

import (
	"fmt"
	"futong-yw-monitor-center/monitor-base/bg"
	bgmodels "futong-yw-monitor-center/monitor-base/models"
	"futong-yw-monitor-center/monitor-center/g"
	"futong-yw-monitor-center/monitor-center/models"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: controllers
 * @File:  monitor
 * @Version: 1.0.0
 * @Date: 2021/11/25 下午4:24
 */

func UpdateFileMonitorShell(c *gin.Context) {

	file, err := c.FormFile("file")

	if err != nil {

		ResponseErrorMsg(c, fmt.Sprintf("file err:%v", err))
		return
	}
	monitorShell := &bgmodels.MonitorMetricsDefault{}
	metricName, b := c.GetPostForm("metricname")
	if !b {
		ResponseErrorMsg(c, fmt.Sprintf("Params Not Found metricname"))
		return
	}

	//save
	fileName := fmt.Sprintf("%d.sh", time.Now().Unix())
	shelldir := path.Join(g.CurrentDir, "alertmanager", "shellData")
	os.MkdirAll(shelldir, 0644)
	savefile := path.Join(shelldir, fileName)
	if err = c.SaveUploadedFile(file, savefile); err != nil {
		errstr := fmt.Sprintf("shell saveUploadFile err:%v\n", err)
		g.GetLog().Error(errstr)
		ResponseErrorMsg(c, errstr)
		return
	}
	monitorShell.MetricName = metricName
	monitorShell.Path = fmt.Sprintf("shell_data/%s", fileName)
	monitorShell.Ident = "shell"
	monitorShell.Desc = c.PostForm("desc")

	err = bg.Validate.Struct(monitorShell)
	if err != nil {
		g.GetLog().Error("http Path:%+v Validate Err:%v\n", c.Request.URL, err)
		ResponseError(c, err)
		return
	}

	if err = models.Db.Save(monitorShell).Error; err != nil {
		g.GetLog().Error("http path:%s Save Record err:%v\n", c.Request.URL, err)
		ResponseErrorMsg(c, err)
		return
	}
	ResponseSuccess(c, "ok")

}
