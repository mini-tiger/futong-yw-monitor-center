package models

import (
	"github.com/jinzhu/gorm"
	_ "strconv"
	"time"
)

type MonitorAuth struct {
	ID             uint `gorm:"primary_key" json:"id"`
	UpdatedAt      time.Time
	AuthType       uint8            `gorm:"column:authtype; comment:'采集方式1=pass 2=pub'" json:"authType" validate:"gte=0,lte=3"`
	Name           string           `gorm:"column:name; size:64" json:"name"`
	Username       string           `gorm:"column:username; size:64" json:"username"`
	AuthStr        string           `gorm:"column:authstr;comment:'私钥或密码';type:longtext" json:"authStr"`
	MonitorDevices []*MonitorDevice `gorm:"many2many:itm_yw_monitor_device_auth;"`
}

func (*MonitorAuth) TableName() string {
	return "itm_yw_monitor_auth"
}

func (m *MonitorAuth) MustRecordMap(db *gorm.DB, mr map[string]interface{}) error {
	if err := db.Where(mr).Take(m).Error; err != nil {
		//g.GetLog().Error("authid find err\n")
		//ResponseErrorMsg(c, fmt.Sprintf("authid :%d find err:%v",authid,err))
		return err
	}

	return nil
}
