package router

/**
 * @Author: Tao Jun
 * @Description: router
 * @File:  service
 * @Version: 1.0.0
 * @Date: 2021/11/26 下午3:48
 */
import (
	"fmt"
	"futong-yw-monitor-center/monitor-agent/g"
	"futong-yw-monitor-center/monitor-agent/hub"
	"futong-yw-monitor-center/monitor-base/bg"
	bgmodels "futong-yw-monitor-center/monitor-base/models"
	bgutils "futong-yw-monitor-center/monitor-base/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func InitWeb() {
	g.GinRouter = gin.Default()
	g.GinRouter.Use(LogMiddleWare())

	g.GinRouter.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome ft-agent Server")
	})

	LoadMetrics()

	//fmt.Println(bg.PushCfgEntry.System.WebPort)
	g.HttpSrv = &http.Server{
		Addr:    fmt.Sprintf(":%s", strconv.Itoa(bg.PushCfgEntry.System.WebPort)),
		Handler: g.GinRouter,
	}

	go func() {
		// 服务连接
		if err := g.HttpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			g.GetLog().Error("listen: %s\n", err)
		} else {
			g.GetLog().Warn("listen: %d Release !!!  \n", bg.PushCfgEntry.System.WebPort)
		}
	}()

	//time.Sleep(10*time.Second)
	//err:=srv.Shutdown(context.Background())
	//if err!=nil{
	//	g.GetLog().Error("http shutdown err:%v\n",err)
	//}
	//g.GetLog().Debug("http shutdown success\n")
}
func LoadMetrics() {
	g.GinRouter.POST("/metrics", func(c *gin.Context) {

		reqMetricsData := &bgmodels.ReqMetricsData{}
		err := c.BindJSON(reqMetricsData)
		if err != nil {
			c.JSON(http.StatusOK, &bgmodels.RespData{
				Code: "0",
				Msg:  fmt.Sprintf("Bind Json err:%v", err),
				Data: nil,
			})
			return
		}

		if reqMetricsData.HostId == "" {
			c.JSON(http.StatusOK, &bgmodels.RespData{
				Code: "0",
				Msg:  "hostid is null",
				Data: nil,
			})
			return
		}

		if reqMetricsData.HostId != g.HostID {
			c.JSON(http.StatusOK, &bgmodels.RespData{
				Code: "0",
				Msg: fmt.Sprintf("hostid Not Match %s  ne %s",
					reqMetricsData.HostId, g.HostID),
				Data: nil,
			})
			return
		}

		if len(reqMetricsData.Metrics) == 0 {
			c.JSON(http.StatusOK, &bgmodels.RespData{
				Code: "0",
				Msg:  "metrics len 0",
				Data: nil,
			})
			return
		}

		// 基础采集
		getMetricsData := new(hub.GetData)
		getMetricsData.CollectData()
		// 采集所有 过滤需要的
		data := getMetricsData.FormatDataFilterMetrics(reqMetricsData.Metrics)

		monitorShells := reqMetricsData.Shells

		//fmt.Println(monitorShells)
		monitorCenterHost := bgutils.HttpHostSplit(g.ViperCfg.ConfWeb.GetConfUrl)
		if monitorCenterHost == "" {
			g.GetLog().Error("monitorCenterHost is null skip shell exec\n")
		} else {
			getMetricsData.MonitorShells = monitorShells
			data["shell"] = getMetricsData.CollectShellData(fmt.Sprintf("http://%s", monitorCenterHost))
		}

		g.GetLog().Debug("Metrics %+v\n", reqMetricsData.Metrics)
		c.JSON(http.StatusOK, &bgmodels.RespData{
			Code: "1",
			Msg:  "ok",
			Data: data,
		})
		return
	})
}
