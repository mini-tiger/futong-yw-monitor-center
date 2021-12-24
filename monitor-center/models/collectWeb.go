package models

import (
	"fmt"
	"github.com/mini-tiger/tjtools/control"
	"sync"

	"futong-yw-monitor-center/monitor-base/bg"
	bgmodels "futong-yw-monitor-center/monitor-base/models"
	"futong-yw-monitor-center/monitor-center/g"
)

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  collectWeb
 * @Version: 1.0.0
 * @Date: 2021/11/28 下午12:21
 */

type CollectWebEntry struct {
	FusionWebHost []bgmodels.CollectHost
	sync.WaitGroup
	Semaphore *control.Semaphore
}

func (c *CollectWebEntry) CollectWeb(cycle int) {
	mds := make([]bgmodels.MonitorDevice, 0)
	err := Db.Preload("MonitorMetrics").
		//Where(map[string]interface{}{"pattern":"web"}).
		Where(map[string]interface{}{"status": 1, "monitorcycle": cycle, "pattern": "web"}).
		Where("hostid != ''").Find(&mds).Error

	if err != nil {
		g.GetLog().Error("collectWeb Err:%v\n", err)
		return
	}
	if len(mds) == 0 {
		return
	}

	fusionWebHost := make([]bgmodels.CollectHost, len(mds))

	for i, md := range mds {
		//fmt.Printf("2222 %+v\n", md)
		//webHost := make(map[string]interface{}, 3)
		webHost := bgmodels.CollectHost{}
		webHost.Ip = md.Ip
		webHost.HostId = md.HostId
		webHost.PushGateAddr = g.GetConfig().PushGateway
		webHost.ESAddr = g.GetConfig().EsServer
		webHost.Arch = md.Arch
		webHost.OS = md.OsType

		if len(md.MonitorMetrics) == 0 { // 没有指标跳过
			//webHost.Metrics = bg.PushCfgEntry.Metrics
			continue
		} else {
			// web 方式 向客户端传送指标
			webHost.Metrics = md.GetMetricsSlice("Ident")

			// shell
			shells, err := md.GetMetrics(Db, "shell")
			if err != nil {
				g.GetLog().Error("IP:%s HostId:%s get shell metrics Err:%v\n", md.Ip, md.HostId, err)
			}
			if len(shells) != 0 {
				webHost.Shells = shells
			}

		}
		fusionWebHost[i] = webHost
	}
	c.FusionWebHost = fusionWebHost
}

func (c *CollectWebEntry) CollectHandle() {
	c.Semaphore = SemaphorePool.Get().(*control.Semaphore)
	for _, wh := range c.FusionWebHost {
		c.Semaphore.Acquire()
		c.Add(1)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					g.GetLog().Error("Weh host:%s err:%v\n", wh.Ip, err)
				}
				c.Semaphore.Release()
				c.Done()
			}()
			// exec
			c.collectWebProcess(wh)
		}()
		c.Wait()
	}
	SemaphorePool.Put(c.Semaphore)

}

func (c *CollectWebEntry) collectWebProcess(wh bgmodels.CollectHost) {
	httpUrl := fmt.Sprintf("http://%s:%d/metrics", wh.Ip, bg.PushCfgEntry.System.WebPort)
	//resp, err := fetch.NewDefault().IsJSON().Post(httpUrl, fetch.NewReader(wh))
	var rd bgmodels.RespData
	resp, err := bg.RestyClient.R().
		SetBody(wh).
		SetResult(&rd). // or SetResult(AuthSuccess{}).
		//SetError(&AuthError{}).       // or SetError(AuthError{}).
		Post(httpUrl)

	if err != nil {
		g.GetLog().Error("web host %+v Request Fail:%v\n", wh, err)
		return
	}

	//body, err := ioutil.ReadAll(resp.Body())
	if resp.StatusCode() != 200 {
		//g.GetLog().Error("http Update HostInfo Resp StatusCode Err: %s\n", resp.StatusCode)
		g.GetLog().Error("web host %+v Resp StatusCode: %v\n", wh, resp.StatusCode())
		return
	}

	if rd.Code != "1" {
		g.GetLog().Error("web host %+v Resp:%v\n", wh, rd)
		return
	}
	wh.HostMetrics = rd.Data.(map[string]interface{})
	allpusherrs, allEserrs := wh.Push2EsPushGateWay()
	if len(allpusherrs) > 0 {
		for metric, e := range allpusherrs {
			g.GetLog().Error("ip:%s hostid:%s Push metrics:%s err:%v\n", wh.Ip, wh.HostId, metric, e)
		}
	} else {
		g.GetLog().Info("IP:%s HostId:%s Push metrics success\n", wh.Ip, wh.HostId)
	}

	if len(allEserrs) > 0 {
		for metric, e := range allEserrs {
			g.GetLog().Error("ip:%s hostid:%s Es metrics:%s err:%v\n", wh.Ip, wh.HostId, metric, e)
		}
	} else {
		g.GetLog().Info("IP:%s HostId:%s Es metrics success\n", wh.Ip, wh.HostId)
	}
}
