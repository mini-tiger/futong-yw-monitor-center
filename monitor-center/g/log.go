package g

import (
	"encoding/json"
	nxlog "github.com/ccpaging/nxlog4go"
	logDiy "github.com/mini-tiger/tjtools/logDiyNew"
	"path"
	"sync"
)

var (
	lock  = new(sync.RWMutex)
	logge *nxlog.Logger
)

func InitLog() *nxlog.Logger {
	// 初始化 日志

	logge = logDiy.InitLog1(path.Join(CurrentDir, cfg.Logfile), cfg.LogMaxDays,
		true, GetConfig().Level, !Product)
	return logge

}

func GetLog() *nxlog.Logger {
	lock.RLock()
	defer lock.RUnlock()
	return logge
}

func PrintConf() {
	if cfg.IsDebug() {
		cfgStr, _ := json.MarshalIndent(cfg, "", "\t")
		logge.Debug("config file data:%+v\n", string(cfgStr))
	}
	logge.Warn("Current Dir:%v\n", CurrentDir)

}
