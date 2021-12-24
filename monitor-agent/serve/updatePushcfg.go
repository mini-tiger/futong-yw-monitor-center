package serve

import (
	"errors"
	"fmt"
	"futong-yw-monitor-center/monitor-agent/g"
	"futong-yw-monitor-center/monitor-base/bg"
	"strconv"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: serve
 * @File:  updatePushcfg
 * @Version: 1.0.0
 * @Date: 2021/12/2 下午1:50
 */

func GetRemoteConfDaemon(td int) {
	if err := GetHttpConfAndCheck(g.ViperCfg.ConfWeb.GetConfUrl); err == nil {

		beginIndep()
		g.SyncOnce.Do(func() {
			g.SwConfig.Done()
		})
	} else {
		g.GetLog().Error(err)
	}

	go func() {
		for {
			if err := GetHttpConfAndCheck(g.ViperCfg.ConfWeb.GetConfUrl); err == nil {
				beginIndep()
				g.SyncOnce.Do(func() {
					g.SwConfig.Done()
				})
			} else {
				g.GetLog().Error("[ PushCfg 更新检查 ] SyncFile Update 失败 %+v\n", bg.PushCfgEntry.GetConf())
			}

			time.Sleep(time.Duration(td) * time.Second)
		}
	}()
}

func beginIndep() {
	if bg.Pattern == "agent" {
		g.BeginChan <- struct{}{}
	}
}

func GetHttpConfAndCheck(confurl string) error {
	reqparams := map[string]string{
		"hostid":  g.HostID,
		"pattern": bg.Pattern,
		"ip":      g.OutIP,
	}
	tpdata := bg.NewRespPushConfig()

	g.GetLog().Info("[ PushCfg 更新检查 ] Req params:%v\n", reqparams)
	resp, err := bg.RestyClient.R().
		SetQueryParams(reqparams).
		SetResult(tpdata).
		ForceContentType("application/json").
		Get(confurl)
	if err != nil {
		g.GetLog().Error("Get config URL:%s err:%v\n", confurl, err)
	}

	if resp.StatusCode() != 200 {
		return errors.New(fmt.Sprintf("http Resp StatusCode Err: %v\n", resp.StatusCode()))
	}

	if code, err := strconv.Atoi(tpdata.Code); err != nil || code != 1 {
		//g.GetLog().Error("Resp Code Err %d\n", code)
		return errors.New(fmt.Sprintf("Resp Code Err %d\n", code))
	}

	err = bg.PushCfgEntry.SetStruct(tpdata.Data)
	if err != nil {

		return errors.New(fmt.Sprintf("PushConfig SetConf Err: %s", err.Error()))
	}

	g.GetLog().Info("[ PushCfg 更新检查 ] SyncFile Update 成功 %+v\n", bg.PushCfgEntry.GetConf())

	return nil
}
