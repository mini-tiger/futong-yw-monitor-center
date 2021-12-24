package controllers

import (
	"fmt"
	"futong-yw-monitor-center/monitor-base/bg"
	bgfuncs "futong-yw-monitor-center/monitor-base/funcs"
	bgmodels "futong-yw-monitor-center/monitor-base/models"
	"futong-yw-monitor-center/monitor-center/g"
	"futong-yw-monitor-center/monitor-center/models"
	"github.com/gin-gonic/gin"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: controllers
 * @File:  monitor
 * @Version: 1.0.0
 * @Date: 2021/11/19 上午11:23
 */

func GetConf(c *gin.Context) {
	hostid := c.Query("hostid")
	//pattern := c.Query("pattern")
	ip := c.Query("ip")
	//没有hostid 返回 默认配置
	if len([]byte(hostid)) <= 1 {
		g.GetLog().Error("HostId is nil, Response defaultConfig\n")
		ResponseSuccess(c, bg.PushCfgEntry.GetConf())
		return
	}

	if len([]byte(ip)) <= 1 {
		g.GetLog().Error("IP is nil, IP:%v Response defaultConfig\n", c.Query(ip))
		ResponseSuccess(c, bg.PushCfgEntry.GetConf())
		return
	}

	md := &bgmodels.MonitorDevice{}
	//err := models.Db.Where(map[string]interface{}{"hostid": hostid}).p.Take(md).Error

	err := models.Db.Preload("MonitorMetrics").
		Where("hostid = ?", hostid).Or("ip = ?", ip).
		Take(&md).Error

	if err != nil {
		ResponseSuccess(c, bg.PushCfgEntry.GetConf())
		g.GetLog().Error("http Path:%+v Err:%v Resp default config\n", c.Request.URL, err)
		return
	}

	err = bg.Validate.Struct(md)
	if err != nil {
		ResponseSuccess(c, bg.PushCfgEntry.GetConf())
		g.GetLog().Error("http Path:%+v Validate Err:%v\n", c.Request.URL, err)
		return
	}

	//没有hostid 返回 默认配置
	if md.ID == 0 {
		g.GetLog().Debug("HostId:%s Not Found, Response defaultConfig\n", hostid)
		ResponseSuccess(c, bg.PushCfgEntry.GetConf())
		return
	}

	// 更新上传时间字段
	go func() {
		updateDbManyRecord := &bgfuncs.UpdateDbManyRecordEntry{}
		updateDbManyRecord.QueryParams = &bgmodels.MonitorDevice{ID: md.ID}
		updateDbManyRecord.UpdateParams = map[string]interface{}{"LastReqCfg": uint64(time.Now().Unix())}
		if err = updateDbManyRecord.RetryUpdateCol(models.Db); err != nil {
			g.GetLog().Error("IP:%s HostId:%s Update LastReqCfg Err:%v\n", md.Ip, md.HostId, err)
		}

	}()

	//复制 default pushconfig
	pcfg := bg.NewPushConfig()
	*pcfg = *(bg.PushCfgEntry.GetConf())

	//关联指标个数
	metrics := md.GetMetricsSlice("Ident")
	// 当前主机没有关联指标，使用默认指标集
	if len(metrics) != 0 {
		pcfg.Metrics = metrics
	}

	shells, err := md.GetMetrics(models.Db, "shell")
	if err != nil {
		g.GetLog().Error("IP:%s HostId:%s get shell metrics Err:%v\n", md.Ip, md.HostId, err)
	}

	if len(shells) != 0 {
		pcfg.Shells = shells
	}

	pcfg.System.Interval = int(md.MonitorCycle)

	g.GetLog().Debug("hostid:%s, cfg:%v\n", hostid, pcfg)
	ResponseSuccess(c, pcfg)
	return
}

func PutConf(c *gin.Context) {
	var err error
	jsonstr, err := c.GetRawData()
	if err != nil {
		//ResponseBasic(c, 0, fmt.Sprintf("Request json err:%s",err.Error()))
		ResponseErrorMsg(c, fmt.Sprintf("Request json err:%s", err.Error()))
		return
	}

	//pc := g.PushCfgEntry.GetConf()
	//fmt.Printf("%+v\n", pc)

	err = bg.PushCfgEntry.SetConf(jsonstr)
	if err != nil {
		ResponseErrorMsg(c, fmt.Sprintf("json Set解析 err:%s", err.Error()))
		return
	}

	ResponseSuccess(c, bg.PushCfgEntry.GetConf())
	return
}
