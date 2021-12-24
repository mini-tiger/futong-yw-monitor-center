package g

import (
	"github.com/shirou/gopsutil/host"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	BeginChan    = make(chan struct{}, 1)
	CurrentDir   string
	CurrentPwd   string
	OutIP        string
	HostID       string
	CurrentCycle int
	SwConfig     sync.WaitGroup
	SyncOnce     sync.Once
	HostBaseInfo *host.InfoStat
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	CurrentDir = filepath.Dir(filepath.Dir(file))

	file, _ = filepath.Abs(os.Args[0])
	CurrentPwd = filepath.Dir(file)
	os.Chdir(CurrentDir)

	//log.Println("Current WorkDir :", CurrentDir)
}
