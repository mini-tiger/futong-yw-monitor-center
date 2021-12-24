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

func CollectSSHBegin() {
	go func() {
		for {
			cc := new(models.CollectSSHEntry)
			cc.CollectSSH(1)
			cc.AgentStart()
			time.Sleep(1 * time.Minute)
		}
	}()
	go func() {
		for {
			cc := new(models.CollectSSHEntry)
			cc.CollectSSH(5)
			cc.AgentStart()
			time.Sleep(1 * time.Minute)
		}
	}()

}
