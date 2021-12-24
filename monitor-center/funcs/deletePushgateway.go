package funcs

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"futong-yw-monitor-center/monitor-base/bg"
	bgmodels "futong-yw-monitor-center/monitor-base/models"
	bgutils "futong-yw-monitor-center/monitor-base/utils"
	"futong-yw-monitor-center/monitor-center/g"
	"futong-yw-monitor-center/monitor-center/models"
	mapset "github.com/deckarep/golang-set"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: authfunc
 * @File:  deletePushgateway
 * @Version: 1.0.0
 * @Date: 2021/11/26 下午3:14
 */

func DeletePushGateway() {
	for {
		delPushDataHandle()

		time.Sleep(time.Duration(g.GetConfig().DelDirtyDataMinute) * time.Minute)
	}
}
func delPushDataHandle() error {
	oldTime := time.Now().Add(-(time.Duration(g.GetConfig().DelPushGatewayMinute) * time.Minute)).Unix()
	//oldTime := time.Now().Unix()
	//fmt.Println(oldTime)
	var mds []bgmodels.MonitorDevice
	err := models.Db.Select("hostid,ip").Where("lastreqcfg < ?", oldTime).Or("lastreqcfg is null").Find(&mds).Error
	if err != nil {
		return err
	}
	if len(mds) == 0 {
		return nil
	}

	g.GetLog().Warn("需要清除oldData的monitorDevice共 [%d] 个: %+v\n", len(mds), mds)
	var metrics []bgmodels.MonitorMetrics
	err = models.Db.Not("ident", []string{"disk", "net"}).Select("ident").Find(&metrics).Error
	if err != nil {
		return err
	}
	// cpu mem load
	for _, md := range mds {
		if md.HostId != "" {
			for _, purl := range bg.PushCfgEntry.PushGateWay {
				for _, metric := range metrics {
					httpPath := fmt.Sprintf("metrics/job/%s/hostid/%s/instance/%s",
						metric.Ident, md.HostId, md.Ip)
					httpUrl := fmt.Sprintf("%s/%s", purl, httpPath)
					_, err := bg.RestyClient.R().
						//SetBody(wh).
						//SetResult(&rd). // or SetResult(AuthSuccess{}).
						//SetError(&AuthError{}).       // or SetError(AuthError{}).
						Delete(httpUrl)
					if err != nil {
						g.GetLog().Error("Delete pushgateway old data http:%s Err:%v\n", httpUrl, err)
					} else {
						g.GetLog().Debug("Delete pushgateway old data http:%s success\n", httpUrl)
					}
				}
			}

		}
	}

	// net disk
	for index, purl := range bg.PushCfgEntry.PushGateWay {
		file := fmt.Sprintf(fmt.Sprintf("pushMetrics%d.log", index))
		err := bgutils.DownLoadPushMetrics(purl, file)
		if err != nil {
			g.GetLog().Error("http DownLoad url:%s err:%v\n ", purl, err)
			continue
		}
		disk_metrics_set, nm, err := GetPushGateWayMetrics(file)
		if err != nil {
			g.GetLog().Error("GetPushWayMetrics Err:%v\n", err)
		}
		// disk
		for _, md := range mds {
			if md.HostId != "" {
				if dmslice, ok := disk_metrics_set[md.HostId]; ok {

					for value := range dmslice.Iterator().C {
						path := base64.RawURLEncoding.EncodeToString([]byte(value.(string)))
						httpPath := fmt.Sprintf("metrics/job/%s/hostid/%s/instance/%s/path@base64/%s",
							"disk", md.HostId, md.Ip, path)
						httpUrl := fmt.Sprintf("%s/%s", purl, httpPath)

						_, err := bg.RestyClient.R().
							//SetBody(wh).
							//SetResult(&rd). // or SetResult(AuthSuccess{}).
							//SetError(&AuthError{}).       // or SetError(AuthError{}).
							Delete(httpUrl)
						if err != nil {
							g.GetLog().Error("Delete pushgateway old data http:%s Err:%v\n", httpUrl, err)
						} else {
							g.GetLog().Debug("Delete pushgateway old data http:%s success\n", httpUrl)
						}
					}
				}

				// net
				if nmslice, ok := nm[md.HostId]; ok {

					for value := range nmslice.Iterator().C {
						name := base64.RawURLEncoding.EncodeToString([]byte(value.(string)))
						httpPath := fmt.Sprintf("metrics/job/%s/hostid/%s/instance/%s/name@base64/%s",
							"net", md.HostId, md.Ip, name)
						httpUrl := fmt.Sprintf("%s/%s", purl, httpPath)

						_, err := bg.RestyClient.R().
							//SetBody(wh).
							//SetResult(&rd). // or SetResult(AuthSuccess{}).
							//SetError(&AuthError{}).       // or SetError(AuthError{}).
							Delete(httpUrl)
						if err != nil {
							g.GetLog().Error("Delete pushgateway old data http:%s Err:%v\n", httpUrl, err)
						} else {
							g.GetLog().Debug("Delete pushgateway old data http:%s success\n", httpUrl)
						}
					}
				}
			}
		}

	}

	return nil
}

type PushMetrics struct {
	HostId   string `json:"hostid"`
	Instance string `json:"instance"`
	Job      string `json:"job"`
	Name     string `json:"name"`
	Path     string `json:"path"`
}

func GetPushGateWayMetrics(file string) (map[string]mapset.Set, map[string]mapset.Set, error) {
	dm, nm, err := GetFileStruct(file)
	//for key,value:=range nm{
	//	fmt.Println(key,value)
	//}
	//for key,value:=range dm{
	//	fmt.Println(key,value)
	//}}
	if err != nil {
		return dm, nm, err
	}

	return dm, nm, err
}

func GetFileStruct(file string) (map[string]mapset.Set, map[string]mapset.Set, error) {
	fi, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	defer fi.Close()

	pushMetricsNet := make(map[string]mapset.Set, 0)
	pushMetricsDisk := make(map[string]mapset.Set, 0)
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		//fmt.Println(string(a))
		row := string(a)
		if strings.Contains(row, "{") && strings.Contains(row, "}") {
			//fmt.Println(1111,row)
			rexjson := regexp.MustCompile(`\{([^)]+)\}`)
			jsonstr := rexjson.FindString(row)
			if strings.Contains(jsonstr, "hostid") &&
				strings.Contains(jsonstr, "disk") && (strings.Contains(jsonstr, "path") || strings.Contains(jsonstr, "Path")) {
				p := PushMetrics{}
				jsonstr = strings.Trim(jsonstr, "}")
				jsonstr = strings.Trim(jsonstr, "{")
				rowSlice := strings.Split(jsonstr, ",")

				//fmt.Println(rowSlice,len(rowSlice))
				hostid := GetSliceStr(rowSlice, "hostid=")
				if hostid != "" {
					p.HostId = strings.Trim(hostid, "\"")
				} else {
					continue
				}
				instance := GetSliceStr(rowSlice, "instance=")
				if instance != "" {
					p.Instance = strings.Trim(instance, "\"")
				}
				job := GetSliceStr(rowSlice, "job=")
				if job != "" {
					p.Job = strings.Trim(job, "\"")
				}
				path := GetSliceStr(rowSlice, "path=")
				if path != "" {
					p.Path = strings.Trim(path, "\"")
				} else {
					path = GetSliceStr(rowSlice, "Path=")
					if path != "" {
						p.Path = strings.Trim(path, "\"")
					} else {
						continue
					}
				}

				if _, ok := pushMetricsDisk[p.HostId]; ok {
					pushMetricsDisk[p.HostId].Add(p.Path)
				} else {
					pushMetricsDisk[p.HostId] = mapset.NewSet()
					pushMetricsDisk[p.HostId].Add(p.Path)
				}
			}
			if strings.Contains(jsonstr, "hostid") &&
				strings.Contains(jsonstr, "net") && strings.Contains(jsonstr, "name") {
				p := PushMetrics{}
				jsonstr = strings.Trim(jsonstr, "}")
				jsonstr = strings.Trim(jsonstr, "{")
				rowSlice := strings.Split(jsonstr, ",")

				//fmt.Println(rowSlice,len(rowSlice))
				hostid := GetSliceStr(rowSlice, "hostid=")
				if hostid != "" {
					p.HostId = strings.Trim(hostid, "\"")
				} else {
					continue
				}
				instance := GetSliceStr(rowSlice, "instance=")
				if instance != "" {
					p.Instance = strings.Trim(instance, "\"")
				}
				job := GetSliceStr(rowSlice, "job=")
				if job != "" {
					p.Job = strings.Trim(job, "\"")
				}
				name := GetSliceStr(rowSlice, "name=")
				if name != "" {
					p.Name = strings.Trim(name, "\"")
				} else {
					continue
				}
				//fmt.Printf("%+v\n",p)
				if _, ok := pushMetricsNet[p.HostId]; ok {
					pushMetricsNet[p.HostId].Add(p.Name)
				} else {
					pushMetricsNet[p.HostId] = mapset.NewSet()
					pushMetricsNet[p.HostId].Add(p.Name)
				}
			}
		}
	}
	return pushMetricsDisk, pushMetricsNet, nil
}

func GetSliceStr(slice []string, substr string) string {
	for _, v := range slice {
		if strings.Contains(v, substr) {
			tslice := strings.Split(v, substr)
			if len(tslice) > 1 {
				return tslice[1]
			}
		}
	}
	return ""
}
