package daemon

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

var appName string
var appPath string

func init() {
	file, _ := filepath.Abs(os.Args[0])
	appPath = filepath.Dir(file)
	log.Println("app path: ", appPath)
	appName = filepath.Base(file)
	Daemon()

}
func Daemon() {
	//// daemon
	os.Chdir(appPath)
	fmt.Println(os.Args[0])
	fmt.Println(os.Args[1:])
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // 使子进程拥有自己的 pgid，等同于子进程的 pid
		//Credential: &syscall.Credential{
		//	Uid: uint32(uid),
		//	Gid: uint32(gid),
		//},
	}
	_, err := cmd.CombinedOutput() //不阻塞
	//fmt.Println(string(oo))
	if err != nil {
		log.Fatalf("cmd err:%v\n", err)
		return
	}
	fmt.Printf("%s [PID] %d running...\n", os.Args[0], cmd.Process.Pid)
	os.Exit(0)
}
