package models

import (
	"fmt"
	bgfuncs "futong-yw-monitor-center/monitor-base/funcs"
	"futong-yw-monitor-center/monitor-center/g"
	"github.com/ghodss/yaml"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"path/filepath"
	"reflect"
	_ "strconv"
	"strings"
	"time"
)

type MonitorDevice struct {
	ID                 int    `gorm:"primaryKey" json:"id"`
	OsType             string `gorm:"column:ostype" json:"ostype" validate:"required,OsValidation"`
	Arch               string `gorm:"column:arch" json:"arch" validate:"required"`
	HostName           string `gorm:"column:name; size:128;" json:"Hostname"`
	Ip                 string `gorm:"column:ip; size:64; UNIQUE_INDEX:uniqe_hostid_ip" json:"ip" validate:"required,ipv4"`
	Port               int    `gorm:"column:port; default:22;" json:"port" validate:"required,gt=0,lte=65535"`
	Pattern            string `gorm:"column:pattern; default:'ssh'; size:64; not null; comment:'采集方式agent web ssh'" json:"pattern" validate:"required,PatternValidation" `
	Status             int    `gorm:"column:status; default:1; comment:'正常=1暂停=2失败=3'" json:"status" validate:"required,gte=0,lte=3"`
	HostId             string `gorm:"column:hostid; comment:'uuid'; UNIQUE_INDEX:uniqe_hostid_ip" json:"hostid"`
	UpdatedAt          time.Time
	MonitorAuth        []*MonitorAuth    `gorm:"many2many:itm_yw_monitor_device_auth;"`                                                // 关联monitor_auth
	AgentVerCurrent    uint64            `gorm:"column:agent_version_current" json:"agent_version_current"`                            //当前agent版本
	AgentUpDateVersion uint64            `gorm:"column:agent_update_version_id" json:"agent_update_version_id"`                        // monitor_agent 关联   待更新agent版本??                                  // 关联agentVersion 不是id 关联
	MonitorCycle       uint              `gorm:"column:monitorcycle; default:1" json:"monitorcycle" validate:"required,gte=0,lte=120"` // 采集间隔
	MonitorMetrics     []*MonitorMetrics //`gorm:"many2many:itm_yw_monitor_device_metrics;"`                                             //指标 项对多
	MonitorAlerts      []*MonitorAlerts  //`gorm:"many2many:itm_yw_monitor_device_alerts;"`                                              //报警 项对多
	HostInfo           string            `gorm:"column:hostinfo;comment:'主机配置';type:longtext" json:"hostinfo"`
	LastReqCfg         uint64            `gorm:"column:lastreqcfg;comment:'最新获取配置时间'" json:"lastreqcfg"`
}

func (*MonitorDevice) TableName() string {
	return "itm_yw_monitor_device"
}

func (m *MonitorDevice) GetMetricsSlice(filed string) []string {
	metrics := make([]string, len(m.MonitorMetrics))
	if len(m.MonitorMetrics) != 0 {
		for index, value := range m.MonitorMetrics {
			v := reflect.ValueOf(*value)
			metrics[index] = (v.FieldByName(filed)).String()
		}
	}
	return metrics
}

func (m *MonitorDevice) DelMetricsAndAlertsData(db *gorm.DB) error {
	tx := db.Begin()
	for _, deldata := range m.MonitorMetrics {
		tx.Debug().Unscoped().Delete(deldata)
	}
	for _, deldata := range m.MonitorAlerts {
		tx.Debug().Unscoped().Delete(deldata)
	}
	tx.Commit()
	return tx.Error
}

func (m *MonitorDevice) GetMetrics(db *gorm.DB, metric string) ([]*MonitorMetricsDefault, error) {
	metrics := make([]*MonitorMetricsDefault, 0)
	if len(m.MonitorMetrics) == 0 {
		return metrics, nil
	}

	for _, mm := range m.MonitorMetrics {
		if mm.Ident == metric {
			//metrics = append(metrics, mm)
			dm := new(MonitorMetricsDefault)
			if err := db.Model(mm).Related(dm, "monitor_metrics_default_id").Error; err != nil {
				g.GetLog().Debug(err)
				continue
			}
			metrics = append(metrics, dm)

		}
	}
	return metrics, nil

}

func (m *MonitorDevice) GenerateAlertRulesFile(db *gorm.DB) error {
	monitorDevice := &MonitorDevice{}
	mustOneRecord := bgfuncs.GetDbOneRecord{}
	mustOneRecord.Result = monitorDevice
	mustOneRecord.Preload = []string{"MonitorMetrics", "MonitorAlerts"}
	mustOneRecord.Params = map[string]interface{}{"hostid": m.HostId}
	err := mustOneRecord.MustOneRecord(db)
	if err != nil {
		return err
	}

	metrics := make([]Metric, 0)
	for _, alert := range monitorDevice.MonitorAlerts {
		//fmt.Println(alert.Expr)
		//fmt.Println(strings.Contains(alert.Expr, "shell"))
		//fmt.Printf("%+v\n",alert)
		var annotations Annotations
		switch true {
		case strings.Contains(alert.Ident, "硬盘"):
			annotations = Annotations{
				Summary:     fmt.Sprintf("IP:{{$labels.instance}} HOSTID:{{$labels.hostid}}  Path: {{$labels.path}} %s %s %d", alert.Ident, alert.Term, alert.Value),
				Description: "{{$labels.instance}}: {{$labels.job}}  (current value is:{{ $value }})",
			}
			break
		case strings.Contains(alert.Ident, "网络"):
			annotations = Annotations{
				Summary:     fmt.Sprintf("IP:{{$labels.instance}} HOSTID:{{$labels.hostid}}  Net: {{$labels.name}} %s %s %d", alert.Ident, alert.Term, alert.Value),
				Description: "{{$labels.instance}}: {{$labels.job}}  (current value is:{{ $value }})",
			}
			break
		case strings.Contains(alert.Expr, "shell"):
			annotations = Annotations{
				Summary:     fmt.Sprintf("IP:{{$labels.instance}} HOSTID:{{$labels.hostid}}  Net: {{$labels.name}} %s %s %d", alert.Ident, alert.Term, alert.Value),
				Description: "{{$labels.instance}}: {{$labels.job}}  (current value is:{{ $value }})",
			}
			break
		default:
			annotations = Annotations{
				Summary:     fmt.Sprintf("IP:{{$labels.instance}} HOSTID:{{$labels.hostid}}: %s %s %d", alert.Ident, alert.Term, alert.Value),
				Description: "{{$labels.instance}}: {{$labels.job}}  (current value is:{{ $value }})",
			}
		}
		metric := Metric{
			Expr:        fmt.Sprintf(" %s %s %d", fmt.Sprintf(alert.Expr, m.HostId), alert.Term, alert.Value),
			Alert:       fmt.Sprintf("%s %s", alert.Ident, alert.LevelName),
			For:         "2m",
			Labels:      map[string]interface{}{"term": "node", "severity": alert.Level},
			Annotations: annotations,
		}

		metrics = append(metrics, metric)

	}
	ruleGroup := RuleGroup{
		Groups: []SubGroup{
			{
				Name:  m.HostId,
				Rules: metrics,
			},
		},
	}
	y, err := yaml.Marshal(ruleGroup)
	if err != nil {
		//fmt.Printf("err: %v\n", err)
		return err
	}
	fileName := fmt.Sprintf("%s_%s.rules", m.Ip, m.HostId)

	return ioutil.WriteFile(filepath.Join(g.GetConfig().AlertManagerRulesDir, fileName), y, 0644)
}

// defined a struct
type MonitorDeviceExcelRow struct {
	// xxx column 要与title 名字一样
	// use field name as default column name
	//Phone int `xlsx:"column(phone)"`
	// column means to map the column name
	Name string `xlsx:"column(Name)"`
	// you can map a column into more than one field
	Os string `xlsx:"column(OS)" validate:"required,OsValidation"`
	// omit `column` if only want to map to column name, it's equal to `column(AgeOf)`
	Port int `xlsx:"column(port)" validate:"required,gt=0,lte=65535"`
	//Mail string `xlsx:"column(mail);default(abc@mail.com)" validate:"required,email"`
	IP string `xlsx:"column(IP)" validate:"required,ipv4"`

	// split means to split the string into slice by the `|`

	//Slice   []int `xlsx:"split(|);req(Slice);"` // xxx req： 没有title 返回错误

	// *Temp implement the `encoding.BinaryUnmarshaler`
	//Temp    *Temp `xlsx:"column(UnmarshalString)"`
	// support default encoding of json
	//TempEncoding *TempEncoding `xlsx:"column(UnmarshalString);encoding(json)"`
	// use '-' to ignore.
	//Ignored string `xlsx:"-"`
}

func (m *MonitorDevice) Excel2Model(db *gorm.DB, monitorDeviceData map[int]MonitorDeviceExcelRow,
	auth *MonitorAuth) (map[int]MonitorDeviceExcelRow, map[int]interface{}, map[int]MonitorDeviceExcelRow, error) {

	successData := make(map[int]MonitorDeviceExcelRow, 0)
	InsertErrs := make(map[int]interface{}, 0)
	UniqueErrs := make(map[int]MonitorDeviceExcelRow, 0)

	tx := db.Begin()
	for index, value := range monitorDeviceData {
		var count int
		db.Model(&MonitorDevice{}).Where("ip = ?", value.IP).Count(&count)
		if count > 0 {
			UniqueErrs[index] = value
			continue
		}

		// 上传不关联 metrics alerts  调用接口updatehostinfo 创建
		if err := tx.Create(&MonitorDevice{Ip: value.IP, Port: value.Port,
			HostName: value.Name, OsType: value.Os, MonitorAuth: []*MonitorAuth{auth},
			//MonitorMetrics: CopyDefaultMetricsData(db),
			//MonitorAlerts:  CopyDefaultAlertData(db),
		},
		).Error; err != nil {
			InsertErrs[index] = map[string]interface{}{"err": err, "record": value}
		} else {
			successData[index] = value
		}
	}

	return successData, InsertErrs, UniqueErrs, tx.Commit().Error

}
