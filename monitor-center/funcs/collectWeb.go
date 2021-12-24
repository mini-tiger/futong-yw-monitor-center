package funcs

import (
	"futong-yw-monitor-center/monitor-center/models"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: funcs
 * @File:  collectWeb
 * @Version: 1.0.0
 * @Date: 2021/11/26 下午5:29
 */

func CollectWebBegin() {
	go func() {
		for {
			cc := new(models.CollectWebEntry)
			cc.CollectWeb(1)
			cc.CollectHandle()
			time.Sleep(1 * time.Minute)
		}
	}()
	go func() {
		for {
			cc := new(models.CollectWebEntry)
			cc.CollectWeb(5)
			cc.CollectHandle()
			time.Sleep(5 * time.Minute)
		}
	}()

}
