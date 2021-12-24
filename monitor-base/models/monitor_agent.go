package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  agent
 * @Version: 1.0.0
 * @Date: 2021/11/25 上午11:20
 */
type MonitorAgent struct {
	ID        uint `gorm:"primary_key" json:"id"`
	UpdatedAt time.Time
	//Os             string `gorm:"column:os; size:32" json:"os"`
	//Arch           string `gorm:"column:arch; size:32" json:"arch"`
	DownLoadPath   string `gorm:"column:downloadpath; size:256" json:"downLoadPath" validate:"required,url"`
	AgentVersion   uint64 `gorm:"column:agentVersion;" json:"agentVersion" validate:"required"`
	MonitorDevices []MonitorDevice
}

func (*MonitorAgent) TableName() string {
	return "itm_yw_monitor_agent"
}

func (m *MonitorAgent) GetMap(db *gorm.DB, mp map[string]interface{}) error {
	err := db.Where(mp).Take(m).Error
	if err != nil && err.Error() != "record not found" {
		return err
	}
	return nil
}
