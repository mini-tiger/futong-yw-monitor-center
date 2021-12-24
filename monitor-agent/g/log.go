package g

import (
	"encoding/json"
	"fmt"
	"futong-yw-monitor-center/monitor-base/bg"
	nxlog "github.com/ccpaging/nxlog4go"
	logDiy "github.com/mini-tiger/tjtools/logDiyNew"
	"log"
	"os"
	"path"
	"sync"
)

var (
	lock  = new(sync.RWMutex)
	logge *nxlog.Logger
)

func InitLog() *nxlog.Logger {
	// 初始化 日志
	logdir := path.Join("logs")

	err := os.MkdirAll(logdir, os.ModePerm)
	if err != nil {
		log.Println(fmt.Sprintf("Log Dir %s Create Err: %v", path.Join(logdir, string(os.PathSeparator), "logs"), err))
	}
	//fmt.Println(CurrentDir)
	//fmt.Println(logdir)
	logge = logDiy.InitLog1(path.Join(logdir, ViperCfg.Log.Logfile),
		ViperCfg.Log.LogMaxDays,
		true,
		ViperCfg.Log.Level,
		ViperCfg.Log.Stdout)

	return logge

}
func InitSSHLog() *nxlog.Logger {
	logdir := path.Join("logs")

	err := os.MkdirAll(logdir, os.ModePerm)
	if err != nil {
		log.Println(fmt.Sprintf("Log Dir %s Create Err: %v", path.Join(logdir, string(os.PathSeparator), "logs"), err))
	}
	// 初始化 日志
	logge = logDiy.InitLog1(path.Join(logdir, "run.log"),
		7, true, "debug", true)
	return logge
}
func GetLog() *nxlog.Logger {
	lock.RLock()
	defer lock.RUnlock()
	return logge
}

func PrintConf() {

	cfgStr, _ := json.MarshalIndent(ViperCfg, "", "\t")
	logge.Debug("config file data:%+v\n", string(cfgStr))
	logge.Debug("Current Agent Version:%d\n", bg.AgentVer)

}
