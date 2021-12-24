package funcs

import (
	"fmt"
	"futong-yw-monitor-center/monitor-base/bg"
	"futong-yw-monitor-center/monitor-center/g"
)

/**
 * @Author: Tao Jun
 * @Description: funcs
 * @File:  reload
 * @Version: 1.0.0
 * @Date: 2021/12/7 上午10:00
 */

func PrometheusReload() {
	for _, url := range g.GetConfig().Prometheus {
		reloadUrl := fmt.Sprintf("%s/-/reload", url)
		resp, err := bg.RestyClient.R().
			//SetBody(wh).
			//SetResult(&rd). // or SetResult(AuthSuccess{}).
			//SetError(&AuthError{}).       // or SetError(AuthError{}).
			Post(reloadUrl)
		if err != nil {
			g.GetLog().Error("Prometheus %s err:%v\n", reloadUrl, err)
			continue
		}

		if resp.StatusCode() != 200 {
			g.GetLog().Error("Prometheus %s reload statucode:%d\n", reloadUrl, resp.StatusCode())
			continue
		}
		g.GetLog().Debug("Prometheus %s Success\n", reloadUrl)
	}

}
