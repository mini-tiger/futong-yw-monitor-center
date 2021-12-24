package funcs

import (
	bgmodels "futong-yw-monitor-center/monitor-base/models"
	"futong-yw-monitor-center/monitor-center/g"
	"futong-yw-monitor-center/monitor-center/models"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: funcs
 * @File:  deleteNullDbData
 * @Version: 1.0.0
 * @Date: 2021/12/6 下午6:10
 */

func DeleteNullAndDirtyDataWithDb() {
	for {
		delNullDataWithDb()
		delDirtyDataWithDb()
		time.Sleep(time.Duration(g.GetConfig().DelDirtyDataMinute) * time.Minute)
	}
}

func delNullDataWithDb() {
	models.Db.Where("monitor_device_id is null").Or("monitor_device_id = ?", 0).Delete(&bgmodels.MonitorMetrics{})
	models.Db.Where("monitor_device_id is null").Or("monitor_device_id = ?", 0).Delete(&bgmodels.MonitorAlerts{})
}
func delDirtyDataWithDb() {
	//删除 itm_yw_monitor_alerts 没有关联 itm_yw_monitor_device
	models.Db.Exec("DELETE  from itm_yw_monitor_alerts where monitor_device_id not in  \n" +
		"(select a.monitor_device_id from itm_yw_monitor_alerts a , itm_yw_monitor_device d \n" +
		"	where a.monitor_device_id =d.id  GROUP BY a.monitor_device_id)")

	//删除 itm_yw_monitor_metrics 没有关联 itm_yw_monitor_device

	models.Db.Exec("delete from itm_yw_monitor_metrics where monitor_device_id not in \n" +
		" (select a.monitor_device_id from itm_yw_monitor_metrics a , itm_yw_monitor_device d \n" +
		"where a.monitor_device_id =d.id  GROUP BY a.monitor_device_id)")
}
