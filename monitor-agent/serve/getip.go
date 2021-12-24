package serve

import (
	"futong-yw-monitor-center/monitor-agent/g"
	"futong-yw-monitor-center/monitor-agent/hub/collect"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: serve
 * @File:  getip
 * @Version: 1.0.0
 * @Date: 2021/11/19 下午6:36
 */

func GetOutboundIP() {
	for {
		if err := collect.GetOutBandIPHandle(g.ViperCfg.ConfWeb.GetConfUrl); err != nil {
			g.GetLog().Error(err)
		} else {
			break
		}

		time.Sleep(60 * time.Second)

	}

}
