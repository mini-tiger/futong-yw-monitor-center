package middleware

import (
	"fmt"
	"futong-yw-monitor-center/monitor-center/controllers"
	"futong-yw-monitor-center/monitor-center/g"
	"github.com/gin-gonic/gin"
	"github.com/mini-tiger/tjtools/iprange"
	"net"
	"sync"
)

/**
 * @Author: Tao Jun
 * @Description: middleware
 * @File:  IPWhite
 * @Version: 1.0.0
 * @Date: 2021/10/13 下午4:45
 */

var ippool iprange.Pool
var once sync.Once

func initIPWhileList() {
	// xxx ip从mysql 获取 或者redis  定期 数据库sync
	var while string = "192.168.0.0/24 172.0.0.0/8"
	if ippool.Size().String() == "0" {
		var err error
		ippool, err = iprange.ParseRanges(while)
		if err != nil {
			g.GetLog().Error("ip 白名单 Init Err:%v\n", err)
			return
		}
		g.GetLog().Info("初始化 ip 白名单 完成,长度:%d,subnet:%v\n", ippool.Size(), ippool.String())

	}
}

func IPWhiteList() gin.HandlerFunc {
	once.Do(initIPWhileList)

	return func(c *gin.Context) {
		//fmt.Println("1111",c.ClientIP())
		//fmt.Println(ippool.Contains(net.ParseIP(c.ClientIP())))
		if !ippool.Contains(net.ParseIP(c.ClientIP())) {
			c.Abort() // 中止后续的函数处理
			g.GetLog().Error("Access Client IP :%s Not In IP whileList\n", c.ClientIP())
			controllers.ResponseError(c, fmt.Errorf("Access Client IP :%s Not In IP whileList\n", c.ClientIP()))
			return
		} else {
			c.Next()
		}

	}
}
