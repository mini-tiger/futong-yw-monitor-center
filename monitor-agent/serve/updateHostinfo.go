package serve

import (
	"errors"
	"fmt"
	"futong-yw-monitor-center/monitor-agent/g"
	"futong-yw-monitor-center/monitor-agent/hub/collect"
	"futong-yw-monitor-center/monitor-base/bg"
	bgmodels "futong-yw-monitor-center/monitor-base/models"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: serve
 * @File:  updateHostinfo
 * @Version: 1.0.0
 * @Date: 2021/11/25 下午5:48
 */

func GetHostBaseInfo() {
	collect.CollectBasicHostInfo()
}

func UpdateHostInfo() {
	for {
		err := UpdateHostInfoHandle()
		if err != nil {
			g.GetLog().Error("[ 更新主机信息 ] Err:%v \n", err)
		} else {
			g.GetLog().Info("[ 更新主机信息 ] Success\n")
			break
		}
		time.Sleep(10 * time.Second)

	}

}

func UpdateHostInfoHandle() error {
	hostinfo, err := collect.CollectMonitorDeviceHostInfo()
	if err != nil {
		return err
	}
	//_, _ = json.MarshalIndent(hostinfo, "", "\t")
	//fmt.Println(string(b))
	//fmt.Println(g.ViperCfg.ConfWeb.AddHostInfoUrl)

	//resp, err := fetch.NewDefault().IsJSON().Put(bg.PushCfgEntry.AddHostInfoUrl, fetch.NewReader(hostinfo))
	var rd bgmodels.RespData

	//var rd bg.RespData
	resp, err := bg.RestyClient.R().
		SetBody(*hostinfo).
		SetResult(&rd). //
		//SetError(&AuthError{}).
		Put(bg.PushCfgEntry.AddHostInfoUrl)

	if err != nil {

		return errors.New(fmt.Sprintf("resp err:%v", err))
	}

	//body, err := ioutil.ReadAll(resp.Body)
	//g.GetLog().Debug("http Response body:%v\n", string(body))

	if resp.StatusCode() != 200 {
		//g.GetLog().Error("http Update HostInfo Resp StatusCode Err: %s\n", resp.StatusCode)
		return errors.New(fmt.Sprintf("Resp StatusCode Err: %v", resp))
	}
	//err = json.Unmarshal(body, &rd)
	//if err != nil {
	//	//g.GetLog().Error("http Update HostInfo json.Unmarshal Err: %s\n", err)
	//	return errors.New(fmt.Sprintf("json.Unmarshal Err: %v", err))
	//}

	//fmt.Println(rd)

	if rd.Code != "1" {
		//g.GetLog().Error("http Update resp code:%d Err: %s\n", rd.Code, rd.Msg)
		return errors.New(fmt.Sprintf("resp err:%+v", rd))
	}
	return nil
}
