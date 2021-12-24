package controllers

import (
	"futong-yw-monitor-center/monitor-base/bg"
	bgmodels "futong-yw-monitor-center/monitor-base/models"
	"futong-yw-monitor-center/monitor-center/g"
	"futong-yw-monitor-center/monitor-center/models"
	"github.com/gin-gonic/gin"
)

/**
 * @Author: Tao Jun
 * @Description: controllers
 * @File:  monitor
 * @Version: 1.0.0
 * @Date: 2021/11/19 上午11:23
 */

func GetAgentVer(c *gin.Context) {

	var h bg.AgentVerReq
	err := c.BindJSON(&h)
	if err != nil {
		ResponseError(c, err)
		return
	}

	hostid := h.HostId
	agentVer := h.AgentVer

	//没有hostid 返回 默认配置
	if len([]byte(hostid)) <= 1 {
		g.GetLog().Error("HostId is nil\n")
		ResponseErrorMsg(c, "HostId is nil")
		return
	}

	agentversion := uint64(agentVer)
	//agentversion, err := strconv.ParseUint(uint64(agentVer), 10, 64)
	//if err != nil || agentversion == 0 {
	//	g.GetLog().Error("agentVer parse err:%v\n", err)
	//	ResponseError(c, err)
	//	return
	//}

	md := &bgmodels.MonitorDevice{}
	//err := models.Db.Where(map[string]interface{}{"hostid": hostid}).p.Take(md).Error

	err = models.Db.Model(md).Take(&md, map[string]interface{}{"hostid": hostid}).Error

	if err != nil {
		g.GetLog().Error("http Path:%+v Err:%v\n", c.Request.URL, err)
		ResponseError(c, err)
		return
	}

	if md.ID == 0 {
		g.GetLog().Error("http Path:%+v 0 record\n", c.Request.URL)
		ResponseErrorMsg(c, "record is 0")
		return
	}

	err = bg.Validate.Struct(md)
	if err != nil {
		g.GetLog().Error("http Path:%+v Validate Err:%v\n", c.Request.URL, err)
		ResponseError(c, err)
		return
	}

	agentVerRecord := &bgmodels.MonitorAgent{}
	models.Db.Model(md).Related(agentVerRecord, "agent_update_version_id") // Related 返回 一条数据

	var agentResp *bg.AgentVerResp = &bg.AgentVerResp{Update: false}

	if agentVerRecord.AgentVersion > agentversion {
		agentResp.Update = true
		agentResp.AgentUrl = agentVerRecord.DownLoadPath
		agentResp.AgentNewVer = bg.AgentVersion(agentVerRecord.AgentVersion)
		agentResp.Pattern = md.Pattern
		//agentResp.AgentUrl = utils.UrlFormat(g.GetConfig().AgentDownLoadUrl,
		//	fmt.Sprintf("%s/%s/%s", "linux", "amd64", bg.AgentName))
	}

	ResponseSuccess(c, agentResp)
	return
}

func UpdateAgentVer(c *gin.Context) {
	var agent bgmodels.MonitorAgent
	var err error

	err = c.BindJSON(&agent)
	if err != nil {
		ResponseError(c, err)
		return
	}
	err = bg.Validate.Struct(agent)
	if err != nil {
		g.GetLog().Error("http Path:%+v Validate Err:%v\n", c.Request.URL, err)
		ResponseError(c, err)
		return
	}
	// create
	if agent.ID == 0 {

		if err = models.Db.Create(&agent).Error; err != nil {
			g.GetLog().Error("http Path:%+v ,Create Record Err:%v\n", c.Request.URL, err)
			ResponseError(c, err)
			return
		} else {
			ResponseSuccess(c, "success")
			return
		}
	}

	// update
	//var db *gorm.DB
	//agentNew:= &bgmodels.MonitorAgent{}
	//if db=models.Db.Where(map[string]interface{}{"id":agent.ID}).Take(agentNew);db.Error!=nil {
	//	g.GetLog().Error("http Path:%+v Find Err:%v\n", c.Request.URL, err)
	//	ResponseError(c, err)
	//	return
	//}
	//
	//if db.RowsAffected == 0 {
	//	g.GetLog().Error("http Path:%+v not record\n", c.Request.URL)
	//	ResponseError(c, err)
	//	return
	//}
	//
	//err = bg.Validate.Struct(agent)
	//if err != nil {
	//	g.GetLog().Error("http Path:%+v Validate Err:%v\n", c.Request.URL, err)
	//	ResponseError(c, err)
	//	return
	//}

	if err := models.Db.Save(&agent).Error; err != nil {
		g.GetLog().Error("http Path:%+v ,Update Record Err:%v\n", c.Request.URL, err)
		ResponseError(c, err)
		return
	}

	ResponseSuccess(c, "success")
	return
}
