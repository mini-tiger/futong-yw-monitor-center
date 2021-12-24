package main

import (
	"flag"
	"fmt"
	"futong-yw-monitor-center/monitor-base/bg"
	funcs "futong-yw-monitor-center/monitor-center/funcs"
	"futong-yw-monitor-center/monitor-center/g"
	"futong-yw-monitor-center/monitor-center/middleware"
	"futong-yw-monitor-center/monitor-center/models"

	"futong-yw-monitor-center/monitor-center/routers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
)

// xxx gin doc https://www.kancloud.cn/shuangdeyu/gin_book/949420

const ConfigJson = "config.json"

func SetupCfg() {
	_ = g.LoadConfig(ConfigJson)

	// CurrentDir 需要在 LoadConfig 设置
	_ = os.Chdir(g.CurrentDir)

	// 初始化 日志
	g.InitLog()
	g.PrintConf()
}
func SetupServer() (r *gin.Engine) {
	// 默认已经连接了 Logger and Recovery 中间件
	//r := gin.Default()
	// xxx创建一个默认的没有任何中间件的路由
	r = gin.New()

	// windows 无法显示日志颜色
	if runtime.GOOS == "windows" {
		gin.DisableConsoleColor()
	} else {
		gin.ForceConsoleColor()
	}
	gin.SetMode(g.GetConfig().Level)
	//gin.SetMode(gin.ReleaseMode)

	// xxx 全局中间件

	// todo ip 白名单
	//r.Use(middleware.IPWhiteList())

	// Logger 中间件将写日志到 gin.DefaultWriter 即使你设置 GIN_MODE=release.
	// 默认设置 gin.DefaultWriter = os.Stdout
	// r.Use(gin.Logger())

	// 自定义日志中间件,和django一样,中间件 往返都要执行
	r.Use(middleware.LogMiddleWare())
	// 需要将 r.Use(middlewares.Cors()) 在使用路由前进行设置，否则会导致不生效
	r.Use(middleware.Cors())
	// Recovery 中间件从任何 panic 恢复，如果出现 panic，它会写一个 500 错误。
	//r.Use(gin.Recovery())
	r.Use(middleware.Recovery())
	// xxx 加载路由
	routers.LoadRouterMonitor(r)
	//fmt.Println(path.Join(g.CurrentDir, "alertmanager", "shellData"))
	r.StaticFS("/packages", http.Dir(path.Join(g.CurrentDir, "packages")))
	r.StaticFS("/shell_data", http.Dir(path.Join(g.CurrentDir, "alertmanager", "shellData")))
	return
}

func SetupPlugins() {
	// xxx 初始化mysql conn
	models.InitDB()
	err := models.SqlDB.Ping()
	if err != nil {
		g.GetLog().Fatalf("sql err:%s\n", err.Error())
	}
	// xxx 初始化ES
	//funcs.InitES()
}

var (
	GoVersion = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	//GoVersion   string
	Branch    string
	Commit    string
	BuildTime string
	lowercase string // 小写也可以
)

func printInfo() {
	log.Printf("Go Build version: %s\n", GoVersion)
	log.Printf("Branch: %s\n", Branch)
	log.Printf("Commit: %s\n", Commit)
	log.Printf("BuildTime: %s\n", BuildTime)
	log.Printf("lowercase: %s\n", lowercase)
}

func main() {
	/*
		go build -ldflags=" \
		-X main.Branch=`git rev-parse --abbrev-ref HEAD` \
		-X main.Commit=`git rev-parse HEAD` \
		-X main.BuildTime=`date '+%Y-%m-%d_%H:%M:%S'`" \
		-v -o main main.go
	*/

	versionFlag := flag.Bool("version", false, "print the version")
	flag.Parse()

	if *versionFlag {
		printInfo()
		os.Exit(0)
	}

	//printInfo()

	SetupCfg()

	r := SetupServer()

	// init es mysql
	SetupPlugins()

	// default pushconfig
	err := bg.InitPushConfig()
	if err != nil {
		g.GetLog().Fatalf("Init PushConfig Err:%s\n", err.Error())
		log.Fatalf("Init PushConfig Err:%s\n", err.Error())
	}
	//删除没更新pushgateway 旧数据
	go funcs.DeletePushGateway()

	// 删除alert,metrics表中没关联tm_yw_monitor_device 和 空数据
	go funcs.DeleteNullAndDirtyDataWithDb()

	if g.Product {
		//call web monitor
		funcs.CollectWebBegin()

		//call ssh monitor
		funcs.CollectSSHBegin()
	}

	os.MkdirAll(g.GetConfig().AlertManagerRulesDir, 0644)

	err = r.Run(":" + strconv.Itoa(g.GetConfig().Port))
	if err != nil {
		log.Fatalf("Run web Port Err:%s\n", err.Error())
	}
}
