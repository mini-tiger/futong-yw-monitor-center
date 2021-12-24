package models

import (
	"context"
	"errors"
	"fmt"
	"futong-yw-monitor-center/monitor-base/bg"
	bgutils "futong-yw-monitor-center/monitor-base/utils"
	"github.com/imdario/mergo"
	jsoniter "github.com/json-iterator/go"
	"github.com/olivere/elastic/v7"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"math/rand"
	"strings"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  collect_push
 * @Version: 1.0.0
 * @Date: 2021/12/1 上午10:04
 */

type CollectHost struct {
	Ip           string                   `json:"ip"`
	HostId       string                   `json:"hostId"`
	Metrics      []string                 `json:"metrics"` //采集 的指标
	Shells       []*MonitorMetricsDefault `json:"shells"`  //shell采集 的指标
	Labels       map[string]interface{}
	ESAddr       []string `json:"esAddr"`
	PushGateAddr []string `json:"pushGateAddr"`
	OS           string   `json:"os"`
	Arch         string   `json:"arch"`
	ClientEs     *elastic.Client
	HostMetrics  map[string]interface{}
}

func (ch *CollectHost) Push2EsPushGateWay() (map[string][]error, map[string][]error) {
	metrics := ch.HostMetrics
	allPushErrs := make(map[string][]error, 0)
	allEsErrs := make(map[string][]error, 0)

	var esCheck bool = true
	if err := ch.InitEs(); err != nil {
		esCheck = false
		allEsErrs["all"] = []error{errors.New("Es init Fail")}
	}

	for key, value := range metrics {
		metricName := key

		switch metricName {
		case "cpu":
			err := ch.PushCpuToPushGateway("cpu", value)
			if err != nil {
				allPushErrs[metricName] = []error{err}
			}
			if esCheck {
				err := ch.PushCpuToEs(metricName, value)
				if err != nil {
					allEsErrs[metricName] = []error{err}
				}
			}

			continue
		case "disk":

			if esCheck {
				errs := ch.PushDiskToEs(metricName, value)
				if len(errs) > 0 {
					allEsErrs[metricName] = errs
				}
			}
			// pushgateway 会删除path fstype ，所以后面执行
			errs := ch.PushDiskToPushGateway("disk", value)
			if len(errs) > 0 {
				allPushErrs[metricName] = errs
			}

			continue
		case "diskrw":
			err := ch.PushDiskRWToPushGateway(metricName, value)
			if err != nil {
				allPushErrs[metricName] = []error{err}
			}
			if esCheck {
				err := ch.PushDiskRwToEs(metricName, value)
				if err != nil {
					allEsErrs[metricName] = []error{err}
				}
			}

			continue
		case "net":
			if esCheck {
				errs := ch.PushNetToEs(metricName, value)
				if len(errs) > 0 {
					allEsErrs[metricName] = errs
				}
			}
			errs := ch.PushNetToPushGateway("net", value)
			if len(errs) > 0 {
				allPushErrs[metricName] = errs
			}

			continue
		case "mem":
			err := ch.PushMemToPushGateway(metricName, value)
			if err != nil {
				allPushErrs[metricName] = []error{err}
			}
			if esCheck {
				err := ch.PushMemToEs(metricName, value)
				if err != nil {
					allEsErrs[metricName] = []error{err}
				}
			}

			continue
		case "shell":
			err := ch.PushMemToPushGateway(metricName, value)
			if err != nil {
				allPushErrs[metricName] = []error{err}
			}
			if esCheck {
				err := ch.PushMemToEs(metricName, value)
				if err != nil {
					allEsErrs[metricName] = []error{err}
				}
			}
			continue
		default:
			continue
		}
	}
	return allPushErrs, allEsErrs
	//d.PushHealthToPushGateway("up")
}

func (ch *CollectHost) PushNetToPushGateway(metricName string, data interface{}) []error {
	errs := make([]error, 0)
	netinfos := data.([]interface{})
	for _, ni := range netinfos {
		//nm := structs.Map(ni)
		nm := ni.(map[string]interface{})
		ethName := (nm["name"]).(string)
		if !strings.Contains(ethName, "docker") &&
			!strings.Contains(ethName, "veth") && !strings.Contains(ethName, "lo") {
			dylabels := make(map[string]string, 1)

			dylabels["name"] = ethName
			err := ch.pushToPushGateway(metricName, nm, nil, dylabels)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func (ch *CollectHost) PushNetToEs(metricName string, data interface{}) []error {
	errs := make([]error, 0)
	netinfos := data.([]interface{})
	bulkRequest := ch.ClientEs.Bulk()
	for _, ni := range netinfos {
		//nm := structs.Map(ni)
		nm := ni.(map[string]interface{})
		ethName := (nm["name"]).(string)
		if !strings.Contains(ethName, "docker") &&
			!strings.Contains(ethName, "veth") && !strings.Contains(ethName, "lo") {
			//dylabels := make(map[string]string, 1)

			//dylabels["name"] = ethName

			jstr := ch.FillBaseDataCoverJson(nm)
			doc := elastic.NewBulkIndexRequest().Index(metricName).Doc(jstr)
			bulkRequest = bulkRequest.Add(doc)

		}
	}
	response, err := bulkRequest.Do(context.TODO())
	if err != nil {
		errs = []error{err}
		return errs
	}
	failed := response.Failed()
	l := len(failed)
	if l > 0 {
		for _, e := range failed {
			errs = []error{errors.New(fmt.Sprintf("%v", e.Error.Reason))}
			//errs = []error{errors.New(fmt.Sprintf("Error(%d),%v", l, response.Errors))}
		}
	}
	return errs
}

//
func (ch *CollectHost) PushDiskToPushGateway(metricName string, data interface{}) []error {
	ds := data.([]interface{})
	errs := make([]error, 0)

	for _, d := range ds {
		di := d.(map[string]interface{})
		fstype := (di["fstype"]).(string)
		path := (di["path"]).(string)
		total, err := bgutils.ConvertNum(di["total"])
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if total > 0 && !strings.Contains(fstype, "tmpfs") &&
			!strings.Contains(path, "overlay2") && !strings.Contains(fstype, "fuseblk") {
			//fmt.Println(di)
			//dm := structs.Map(di)
			//dm:=make(map[string]interface{},0)
			constlabels := make(map[string]string, 1)
			dyLabels := make(map[string]string, 1)

			constlabels["fstype"] = fstype
			dyLabels["path"] = path
			delete(di, "fstype")
			delete(di, "path")
			err := ch.pushToPushGateway(metricName, di, constlabels, dyLabels)
			if err != nil {
				errs = append(errs, err)
			}
		}
		//fmt.Println(di)
	}
	return errs
}

func (ch *CollectHost) PushDiskToEs(metricName string, data interface{}) []error {
	ds := data.([]interface{})
	errs := make([]error, 0)
	defer func() {
		if err := recover(); err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("metrics: %s To Es Err:%v", metricName, err)))
			return
		}
	}()

	bulkRequest := ch.ClientEs.Bulk()
	for _, d := range ds {
		di := d.(map[string]interface{})
		fstype := (di["fstype"]).(string)
		//path := (di["path"]).(string)
		total, err := bgutils.ConvertNum(di["total"])
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if total > 0 && !strings.Contains(fstype, "tmpfs") &&
			len(fstype) > 1 && !strings.Contains(fstype, "cgroupfs") && !strings.Contains(fstype, "fuseblk") {
			//constlabels := make(map[string]string, 1)
			//dyLabels := make(map[string]string, 1)

			jstr := ch.FillBaseDataCoverJson(di)
			doc := elastic.NewBulkIndexRequest().Index(metricName).Doc(jstr)
			bulkRequest = bulkRequest.Add(doc)
			//g.GetLog().Warn("%v\n", jstr)
		}
	}
	response, err := bulkRequest.Do(context.TODO())
	if err != nil {
		errs = []error{err}
		return errs
	}
	failed := response.Failed()
	l := len(failed)
	if l > 0 {
		for _, e := range failed {
			errs = []error{errors.New(fmt.Sprintf("%v", e.Error.Reason))}
			//errs = []error{errors.New(fmt.Sprintf("Error(%d),%v", l, response.Errors))}
		}
	}
	return errs
}

func (ch *CollectHost) PushMemToPushGateway(metricName string, data interface{}) error {
	meminfo := data.(map[string]interface{})

	//m := structs.Map(meminfo)

	//m["usedPercent"] = cpuinfo.UsedPercent

	return ch.pushToPushGateway(metricName, meminfo, nil, nil)
}

func (ch *CollectHost) PushMemToEs(metricName string, data interface{}) error {
	meminfo := data.(map[string]interface{})
	var jstr string
	if jstr = ch.FillBaseDataCoverJson(meminfo); jstr == "" {
		return errors.New(fmt.Sprintf("%s CovertJson is nil", metricName))
	}
	//fmt.Printf("print mString:%s\n", mString)
	_, err := ch.ClientEs.Index().
		Index(metricName).
		//Type(typeName).
		//Id(strconv.Itoa(subject.ID)).
		BodyJson(jstr).
		//Refresh("wait_for").
		Do(context.TODO())

	if err != nil {
		return err
	}
	return nil
}

func (ch *CollectHost) PushDiskRwToEs(metricName string, data interface{}) error {
	diskRWinfo := data.(map[string]interface{})
	var jstr string
	if jstr = ch.FillBaseDataCoverJson(diskRWinfo); jstr == "" {
		return errors.New(fmt.Sprintf("%s CovertJson is nil", metricName))
	}
	//fmt.Printf("print mString:%s\n", mString)
	_, err := ch.ClientEs.Index().
		Index(metricName).
		//Type(typeName).
		//Id(strconv.Itoa(subject.ID)).
		BodyJson(jstr).
		//Refresh("wait_for").
		Do(context.TODO())

	if err != nil {
		return err
	}
	return nil
}
func (ch *CollectHost) PushDiskRWToPushGateway(metricName string, data interface{}) error {
	diskRWinfo := data.(map[string]interface{})

	//m := structs.Map(meminfo)

	//m["usedPercent"] = cpuinfo.UsedPercent

	return ch.pushToPushGateway(metricName, diskRWinfo, nil, nil)
}
func (ch *CollectHost) PushCpuToPushGateway(metricName string, data interface{}) error {
	dm := data.(map[string]interface{})
	m := make(map[string]interface{}, 0)
	//m := structs.Map(cpuinfo["cpu"])
	if _, ok := dm["timeStat"]; ok {

		//m=dm["timesStat"].(map[string]interface{})
		if err := mergo.Map(&m, dm["timesStat"].(map[string]interface{})); err != nil {
			return err
		}
		m["usedPercent"] = dm["UsedPercent"]
	} else {
		if err := mergo.Map(&m, dm); err != nil {
			return err
		}
	}

	return ch.pushToPushGateway(metricName, m, nil, nil)
}

func (ch *CollectHost) PushCpuToEs(metricName string, data interface{}) error {
	dm := data.(map[string]interface{})
	m := make(map[string]interface{}, 0)
	//m := structs.Map(cpuinfo["cpu"])
	if _, ok := dm["timeStat"]; ok {

		//m=dm["timesStat"].(map[string]interface{})
		if err := mergo.Map(&m, dm["timesStat"].(map[string]interface{})); err != nil {
			return err
		}
		m["usedPercent"] = dm["UsedPercent"]
	} else {
		if err := mergo.Map(&m, dm); err != nil {
			return err
		}
	}

	var jstr string
	if jstr = ch.FillBaseDataCoverJson(m); jstr == "" {
		return errors.New(fmt.Sprintf("%s CovertJson is nil", metricName))
	}

	_, err := ch.ClientEs.Index().
		Index(metricName).
		//Type(typeName).
		//Id(strconv.Itoa(subject.ID)).
		BodyJson(jstr).
		//Refresh("wait_for").
		Do(context.TODO())

	if err != nil {
		return err
	}
	return nil

}
func (ch *CollectHost) FillBaseDataCoverJson(m map[string]interface{}) string {
	m["ip"] = ch.Ip
	m["hostid"] = ch.HostId
	m["os"] = ch.OS
	m["arch"] = ch.Arch
	m["collectTime"] = time.Now().Unix()

	mjson, err := jsoniter.Marshal(m)
	if err != nil {
		return ""
	}
	return string(mjson)
}

func (ch *CollectHost) pushToPushGateway(metricName string,
	m map[string]interface{}, lables map[string]string, dylables map[string]string) error {

	registry := prometheus.NewRegistry()
	//fmt.Printf("%+v\n",m)

	if err := mergo.Map(&lables, ch.Labels); err != nil {
		return errors.New(fmt.Sprintf("mergo map err:%v", err))
	}
	for k, v := range m {
		metric := prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        k,
			Help:        "this is monitor-web",
			ConstLabels: lables,
		})

		if value, err := bgutils.ConvertNum(v); err != nil {
			//g.GetLog().Debug("metricName : %s registry: %s CovertNum: %v  err:%v\n",metricName,k,v,err)
			continue
		} else {
			metric.Set(value)
		}

		if err := registry.Register(metric); err != nil {
			//g.GetLog().Error("metricName : %s registry: %s err:%v\n", metricName, k, err)

			continue
		}
	}

	pushGateWayAddress := bgutils.GetRandAddress(rand.Intn(len(ch.PushGateAddr)), ch.PushGateAddr)
	if pushGateWayAddress == "" {
		return errors.New(fmt.Sprintf("pushgateway Addr is nil"))
	}

	pushClient := push.New(pushGateWayAddress, metricName). // job
								Gatherer(registry).
								Grouping("instance", ch.Ip). // label
								Grouping("hostid", ch.HostId).
								Client(bg.PushGatewayClient)

	for lable, value := range dylables {

		pushClient.Grouping(lable, value)

	}

	err := pushClient.
		Add()
	if err != nil {
		//g.GetLog().Error("PushHostMetricsToPushGateWay job=%s instance=%s  err: %s\n", metricName, ch.Ip, err.Error())
		return errors.New(fmt.Sprintf("PushHostMetricsToPushGateWay job=%s instance=%s  err: %s\n", metricName, ch.Ip, err.Error()))
	} else {
		//g.GetLog().Info("PushHostMetricsToPushGateWay  job=%s instance=%s  data:%v success\n", metricName, ch.Ip, m)
	}
	return nil
}

func (ch *CollectHost) InitEs() error {
	var err error
	ch.ClientEs, err = elastic.NewClient(
		elastic.SetURL(ch.ESAddr...),
		elastic.SetSniff(false))

	return err

}
