package serve

import (
	"futong-yw-monitor-center/monitor-agent/g"
	"futong-yw-monitor-center/monitor-agent/hub"
	"futong-yw-monitor-center/monitor-base/bg"
	"github.com/robfig/cron/v3"
	"time"
)

var C *cron.Cron = cron.New()

func RunServerCron() {
	g.GetLog().Debug("begin work cron 间隔时间 : %d 分钟\n", bg.PushCfgEntry.System.Interval)

	for _, e := range C.Entries() {
		C.Remove(e.ID)
	}

	t := cron.Every(time.Duration(bg.PushCfgEntry.System.Interval) * time.Minute) // 每执行
	//t := cron.Every(5 * time.Second) // 每执行
	C.Schedule(t, &hub.GetData{}) // 调用Run方法
	C.Start()

}

func MonitorDaemon() {
	for {
		select {
		case <-g.BeginChan:

			if bg.PushCfgEntry.System.Interval <= 0 {
				g.GetLog().Error("Remote PushConfig Fail Not Cron,Interval:%d \n", bg.PushCfgEntry.System.Interval)
				return
			}
			if g.CurrentCycle == bg.PushCfgEntry.System.Interval {
				g.GetLog().Warn("Remote PushConfig Interval: %d minute [eq],Local old Interval: %d\n",
					bg.PushCfgEntry.System.Interval, g.CurrentCycle)
				continue
			} else {
				g.GetLog().Warn("Remote PushConfig Interval: %d Minute [change],Local old Interval: %d\n",
					bg.PushCfgEntry.System.Interval, g.CurrentCycle)

				g.CurrentCycle = bg.PushCfgEntry.System.Interval
			}
			//fmt.Printf("%+v\n",bg.PushCfgEntry.Metrics)
			g.GetLog().Info("[ 重启任务 ] 间隔时间 %d 分钟,指标集:%v\n",
				bg.PushCfgEntry.System.Interval, bg.PushCfgEntry.Metrics)
			C.Stop()
			RunServerCron()
		}
	}
}
