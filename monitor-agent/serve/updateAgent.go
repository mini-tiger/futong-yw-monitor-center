package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"futong-yw-monitor-center/monitor-agent/g"
	"futong-yw-monitor-center/monitor-agent/utils"
	"futong-yw-monitor-center/monitor-base/bg"
	"github.com/goinggo/mapstructure"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sync"
	"syscall"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: serve
 * @File:  update
 * @Version: 1.0.0
 * @Date: 2021/11/20 下午11:33
 */

var mutex sync.Mutex

func UpdateAgentDaemon(td int) {

	for {

		if agentResp, err := GetHttpCheckAgent(); err == nil {
			mutex.Lock()

			if agentResp.Update {
				SelfUpdate(agentResp)
			}
			mutex.Unlock()
		}
		time.Sleep(time.Duration(td) * time.Second)

	}
}
func GetHttpCheckAgent() (*bg.AgentVerResp, error) {
	agentResp := &bg.AgentVerResp{}
	var err error

	defer func() {
		if errp := recover(); errp != nil {
			//g.GetLog().Error("http GetHttpCheckAgent panic err:%v\n", errp)
			err = errors.New(fmt.Sprintf("http GetHttpCheckAgent panic err:%v\n", errp))
		}
	}()
	g.GetLog().Info("[ AgentVer 更新检查 ] Current Version: %v\n", bg.AgentVer)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	agentReq := &bg.AgentVerReq{
		Os:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		AgentVer: bg.AgentVer,
		HostId:   g.HostID}

	reqbyte, err := json.Marshal(agentReq)
	if err != nil {
		//g.GetLog().Error("http Req Json err:%s\n", err.Error())
		return agentResp, err
	}
	req, err := http.NewRequest("POST", bg.PushCfgEntry.GetSelfUpdateUrl, bytes.NewBuffer(reqbyte))

	if err != nil {
		//g.GetLog().Error("http Req err:%s\n", err.Error())
		return agentResp, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		//g.GetLog().Error("http Resp err:%s\n", err.Error())
		return agentResp, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//g.GetLog().Error("http Resp body err:%s\n", err.Error())
		return agentResp, err
	}

	//fmt.Println("response Body:", string(body))

	var respMap map[string]interface{} = make(map[string]interface{})
	err = jsoniter.Unmarshal(body, &respMap)
	if err != nil {
		//g.GetLog().Error("http Resp body Json Unmarshal err:%s\n", err.Error())
		return agentResp, err
	}

	if code, ok := respMap["code"]; ok {
		if code.(string) == "1" {
			agentRespMap := respMap["data"].(map[string]interface{})
			//agentResp := &bg.AgentVerResp{}

			if err := mapstructure.Decode(agentRespMap, agentResp); err != nil {
				//g.GetLog().Error("http Resp body Json agentResp Unmarshal err:%s\n", err.Error())
				return agentResp, err
			}

			return agentResp, nil
		}
	}
	return agentResp, errors.New(fmt.Sprintf("resp err:%v", respMap))
}

func SelfUpdate(agentResp *bg.AgentVerResp) {

	g.GetLog().Info("[ 自更新进程 ] Current Dir: %s\n", g.CurrentPwd)
	//agent download
	ftAgentTmpName := path.Join(g.CurrentPwd, "ft-agent-new")
	agentFileName := path.Join(g.CurrentPwd, bg.AgentName)
	err := utils.HttpDownFile(agentResp.AgentUrl, ftAgentTmpName)

	if err != nil {
		g.GetLog().Error("selfUpdate DownFile err:%s\n", err.Error())
		return
	}
	g.GetLog().Debug("[ 自更新进程 ] DownFile 成功 \n")
	// agent remove
	err = os.Remove(agentFileName)
	if err != nil {
		g.GetLog().Error("selfUpdate os.Remove err:%s\n", err.Error())

	}
	// newagent rename
	err = os.Rename(ftAgentTmpName, agentFileName)
	if err != nil {
		g.GetLog().Error("selfUpdate os.Rename err:%s\n", err.Error())
		return
	}

	cmdstr := fmt.Sprintf("chmod 755 %s  ", agentFileName)
	cmd := exec.Command("/bin/bash", "-c", cmdstr)

	_, err = cmd.CombinedOutput()
	//fmt.Println(string(o))
	if err != nil {
		g.GetLog().Error("cmd err:%v\n", err)
		return
	}

	var startCmdStr string
	if startCmdStr = bg.GetPatternStartParams(agentResp.Pattern, bg.ConfigUrl); startCmdStr == "" {
		g.GetLog().Error("Pattern err:%v\n", agentResp.Pattern)
		return
	}

	if bg.Pattern == "web" {
		g.HttpSrv.Shutdown(context.TODO())
		time.Sleep(1 * time.Second)
	}
	g.GetLog().Debug("[ 自更新进程 ] 准备开启新版本程序\n")

	//cmdstr2 := fmt.Sprintf(bg.GetPatternStartParams(agentResp.Pattern), bg.AgentName)
	//fmt.Println(cmdstr2)
	//restartCmd := fmt.Sprintf("nohup %s &", startCmdStr)
	restartCmd := "systemctl restart ft-agent"
	cmd2 := exec.Command("/bin/bash", "-c", restartCmd)
	// xxx 拥有自己的进程组,子进程 独立于 父进程
	cmd2.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // 使子进程拥有自己的 pgid，等同于子进程的 pid
		//Credential: &syscall.Credential{
		//	Uid: uint32(uid),
		//	Gid: uint32(gid),
		//},
	}
	err = cmd2.Start() //xxx 新进程独立于 本进程
	//fmt.Println(string(oo))
	if err != nil {
		g.GetLog().Error("cmd2 err:%v\n", err)
		return
	}

	g.GetLog().Warn("[ 自更新进程 ]  end process,New %s Father PID: %+v ,bye bye!!!\n", bg.AgentName, cmd2.Process.Pid)
	os.Exit(0)
}
