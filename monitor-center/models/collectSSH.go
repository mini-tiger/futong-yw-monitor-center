package models

import (
	bgmodels "futong-yw-monitor-center/monitor-base/models"
	"futong-yw-monitor-center/monitor-center/g"
	"github.com/mini-tiger/tjtools/control"
	"sync"
)

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  collectWeb
 * @Version: 1.0.0
 * @Date: 2021/11/28 下午12:21
 */

var SemaphorePool = &sync.Pool{
	New: func() interface{} {
		return control.NewSemaphore(10)
	},
}

type CollectSSHEntry struct {
	Fusion_SSH_Host []SSH_Host
	sync.WaitGroup
	Semaphore *control.Semaphore
}

func (c *CollectSSHEntry) CollectSSH(cycle int) {
	mds := make([]bgmodels.MonitorDevice, 0)
	err := Db.
		Preload("MonitorAuth").Preload("MonitorMetrics").
		//Where(map[string]interface{}{"pattern":"web"}).
		Where(map[string]interface{}{"status": 1, "monitorcycle": cycle, "pattern": "ssh"}).
		Find(&mds).Error

	if err != nil {
		g.GetLog().Error("collectSSH  Err:%v\n", err)
		return
	}
	if len(mds) == 0 {
		return
	}
	sshHosts := make([]SSH_Host, 0)
	for _, md := range mds {
		//ssh 第一次上传没有hostid 没有关联指标
		if len(md.MonitorMetrics) == 0 && md.HostId != "" {
			g.GetLog().Debug("SSH HostID:%s,IP:%s Metrics is 0 skip\n", md.HostId, md.Ip)
			continue
		}

		if len(md.MonitorAuth) > 0 {
			sshHost := SSH_Host{
				SSHClient: nil,
				Username:  md.MonitorAuth[0].Username,
				Password:  md.MonitorAuth[0].AuthStr,
				Port:      md.Port,
				Ip:        md.Ip,
				HostId:    md.HostId,
				OsType:    md.OsType,
				AuthType:  md.MonitorAuth[0].AuthType,
			}
			sshHosts = append(sshHosts, sshHost)
		} else {
			g.GetLog().Debug("SSH HostID:%s,IP:%s Auth nil skip\n", md.HostId, md.Ip)
		}

	}
	c.Fusion_SSH_Host = sshHosts
}

func (c *CollectSSHEntry) AgentStart() {
	c.Semaphore = SemaphorePool.Get().(*control.Semaphore)
	for _, value := range c.Fusion_SSH_Host {
		c.Semaphore.Acquire()
		c.Add(1)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					g.GetLog().Error("SSH host:%s err:%v\n", value.Ip, err)
				}
				c.Semaphore.Release()
				c.Done()
			}()
			// exec
			value.CmdExec()
		}()
		c.Wait()
	}
	SemaphorePool.Put(c.Semaphore)
}
