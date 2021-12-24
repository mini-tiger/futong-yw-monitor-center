package main

import (
	"fmt"
	//_ "futong-yw-monitor-center/monitor-agent/daemon"
	"futong-yw-monitor-center/monitor-agent/g"
	"futong-yw-monitor-center/monitor-agent/hub"
	"futong-yw-monitor-center/monitor-agent/hub/collect"
	"futong-yw-monitor-center/monitor-agent/router"
	"futong-yw-monitor-center/monitor-agent/serve"
	"futong-yw-monitor-center/monitor-base/bg"
	flag "github.com/spf13/pflag" //替换原生的flag，并兼容
	"os"
)

var indep bool      // 独立
var web bool        // 是否web 方式
var once bool       // 是否无代理模式
var config string   //无代理模式的配置
var getVersion bool // 打印 版本
//var goDaemon bool   //后台 服务  -d=true

func main() {
	flag.BoolVarP(&web, "web", "w", false, "web方式 ")
	flag.BoolVarP(&once, "once", "o", false, "无代理模式")
	flag.BoolVarP(&getVersion, "version", "v", false, "打印版本")
	flag.StringVarP(&config, "config", "c", "http://172.16.71.20:8001/api/conf", "打印get config url")
	//flag.BoolVarP(&goDaemon, "daemon", "d", false, "run app as a daemon with -d=true.")
	flag.BoolVarP(&indep, "indep", "i", true, "agent 独立方式 ") // xxx default
	flag.Parse()

	if getVersion {
		bg.GetAgentVer()
		os.Exit(0)
	}
	fmt.Printf("web:%v once:%v indep:%v\n", web, once, indep)

	//获取 主机信息hostid
	serve.GetHostBaseInfo()

	if once {
		//fmt.Println(config)
		bg.Pattern = "ssh"
		bg.ConfigUrl = config

		g.InitSSHLog()
		g.GetLog().Debug("bg config:%v;config:%v\n", bg.ConfigUrl, config)

		// 获取outip
		if err := collect.GetOutBandIPHandle(bg.ConfigUrl); err != nil {
			g.GetLog().Error("Get OutIP Err:%v\n", err)
			os.Exit(0)
		}

		// read config
		err := serve.GetHttpConfAndCheck(config)
		if err != nil {
			g.GetLog().Error("get configUrl:%s Fail:%v\n", config, err)
			os.Exit(0)
		}

		// upload hostinfo  需要配置文件中的uploadhost url
		if err := serve.UpdateHostInfoHandle(); err != nil {
			g.GetLog().Error("[update host] err %v\n", err)
		} else {
			g.GetLog().Info("[update host] success\n")
		}

		//XXX check agentUpdate once
		if agentResp, err := serve.GetHttpCheckAgent(); err == nil {
			if agentResp.Update {
				serve.SelfUpdate(agentResp)
			}
		} else {
			g.GetLog().Error("agent update err:%v\n", err)
		}

		new(hub.GetData).Run()

		os.Exit(0)
	}

	if web || indep {

		g.InitConfig()

		g.InitLog()

		if web {
			bg.Pattern = "web"
		}
		if indep {
			bg.Pattern = "agent"
		}

		// 获取outip, 获取 配置文件 需要 outip
		if err := collect.GetOutBandIPHandle(g.ViperCfg.ConfWeb.GetConfUrl); err != nil {
			g.GetLog().Error("Get OutIP Err:%v\n", err)
			os.Exit(0)
		}

		g.SwConfig.Add(1)

		// 每 S秒 获取配置,读取到配置后 往下执行
		serve.GetRemoteConfDaemon(30)

		g.SwConfig.Wait()

		// XXX agent update daemon
		go serve.UpdateAgentDaemon(30)

		g.PrintConf()
		g.GetLog().Warn("Push config %v\n", bg.PushCfgEntry.GetConf())

	}

	if web {
		router.InitWeb()

		//获得pushconfig之后 ,上传hostinfo 启动一次,
		go serve.UpdateHostInfo()
	}

	// 修改逻辑 先获取配置，配置不到 请求本地配置，只能配置正常在 往后走逻辑
	if indep {

		//if err := bg.InitPushConfig(); err != nil { //先从本地查找PushConfig.json
		//	g.GetLog().Error(err.Error())
		//}
		g.GetLog().Debug("初始化加载 PushConfig 成功 %+v \n", bg.PushCfgEntry.GetConf())
		//获得pushconfig之后 ,上传hostinfo 启动一次,
		go serve.UpdateHostInfo()

		//本地没有pushconfig.json 远程也无法获取配置，则无法采集上传
		go serve.MonitorDaemon()
		//go serve.GetOutboundIP() //通过连接远程 HBS接口 获取本机出网IP
	}

	fmt.Println("Current Pattern: ", bg.Pattern)
	select {}

}
