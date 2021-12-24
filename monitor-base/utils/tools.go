package utils

import (
	"errors"
	"futong-yw-monitor-center/monitor-base/bg"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"
)

/**
 * @Author: Tao Jun
 * @Description: utils
 * @File:  tools
 * @Version: 1.0.0
 * @Date: 2021/11/29 上午9:41
 */

func ConvertNum(v interface{}) (float64, error) {
	switch v.(type) {
	case float64:
		return v.(float64), nil
	case float32:
		return float64(v.(float32)), nil
	case int:
		return float64(v.(int)), nil
	case int32:
		return float64(v.(int32)), nil
	case uint32:
		return float64(v.(uint32)), nil
	case uint:
		return float64(v.(uint)), nil
	case uint64:
		return float64(v.(uint64)), nil
	}
	return 0, errors.New("err type")
}

func GetRandAddress(rn int, all []string) string {
	//pushGateWayAddress := ""
	//switch collectPeriod {
	//case 1:
	//	pushGateWayAddress = bg.PushCfgEntry.Prometheus.PushGateWayAddressFirst
	//case 2:
	//	pushGateWayAddress = bg.PushCfgEntry.Prometheus.PushGateWayAddressSecond
	//case 3:
	//	pushGateWayAddress = bg.PushCfgEntry.Prometheus.PushGateWayAddressThird
	//default:
	//	pushGateWayAddress = bg.PushCfgEntry.Prometheus.PushGateWayAddressFirst
	//}
	//return bg.PushCfgEntry.Prometheus[collectPeriod]
	return all[rn]
}

func HttpDownFile(url string, file string) error {
	os.Remove(file)
	//imgUrl := "http://172.16.71.31:19091/metrics"

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//os.Chdir("/home/go/src/godev")
	err = ioutil.WriteFile(file, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func DownLoadPushMetrics(purl string, file string) error {
	os.Remove(file)
	u, _ := url.Parse(purl)
	u.Path = path.Join(u.Path, "metrics")
	urlpath := u.String()
	_, err := bg.RestyClient.R().
		SetOutput(file).
		ForceContentType("text/plain").
		Get(urlpath)
	//fmt.Println(string(resp.Body()))

	if err != nil {
		return err
	}
	return nil
}

func Struct2Map(d interface{}, lower bool) map[string]interface{} {
	m := make(map[string]interface{}, 0)
	t := reflect.TypeOf(d)
	v := reflect.ValueOf(d)
	for k := 0; k < t.NumField(); k++ {
		if lower {
			m[strings.ToLower(t.Field(k).Name)] = v.Field(k).Interface()
		} else {
			m[(t.Field(k).Name)] = v.Field(k).Interface()
		}
		//fmt.Println("name:", fmt.Sprintf("%+v", t.Field(k).Name),
		//	", value:", fmt.Sprintf("%v", v.Field(k).Interface()),
		//	", yaml:", t.Field(k).Tag.Get("yaml"))
	}
	return m
}

func SliceIsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func HttpHostSplit(urlstr string) string {
	u, err := url.Parse(urlstr)
	if err != nil {
		return ""
	}

	return u.Host
}

func TrimEnterSpace(s string) string {
	return strings.Trim(strings.TrimSpace(s), "\n")
}
