package routers

import (
	"futong-yw-monitor-center/monitor-center/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadRoute(router *gin.Engine) {

	// Hello World
	router.GET("/", func(c *gin.Context) {
		//c.String(http.StatusOK, "Hello, World,TaoJun First Deploy for DataCenter")
		c.JSON(http.StatusOK, gin.H{
			//"status":        http.StatusOK,
			"statusText": "Hello, World,TaoJun First Deploy for DataCenter",
		})
	})

	//xxx gin 自带 json parse方式 https://cloud.tencent.com/developer/article/1689928
}

var Api *gin.RouterGroup

func LoadRouterMonitor(router *gin.Engine) {
	// Hello World
	router.GET("/", func(c *gin.Context) {
		//c.String(http.StatusOK, "Hello, World,TaoJun First Deploy for DataCenter")
		c.JSON(http.StatusOK, gin.H{
			//"status":        http.StatusOK,
			"statusText": "Hello, World,TaoJun First Deploy for DataCenter",
		})
	})

	// work api
	Api = router.Group("/api")

	// xxx 配置文件
	InitConfigViews()

	// xxx MonitorDevice excelupload , uploadhostinfo
	InitHostInfoViews()

	// xxx agent version
	InitAgentViews()

	//
	InitModelMonitorDevice()

	InitModelMonitorShells()

	InitModelMonitorAlerts()
}
func InitAgentViews() {
	routerAgentSub := Api.Group("/agent")
	routerAgentSub.POST("", controllers.GetAgentVer)
	routerAgentSub.PUT("/updateVer", controllers.UpdateAgentVer) // update or add
}

func InitConfigViews() {
	routerConfSub := Api.Group("/conf")
	routerConfSub.GET("", controllers.GetConf) //通过 hostid Query
	//routerSub.POST("", controller.GetMonitorAuth)   // get

}

func InitHostInfoViews() {
	monitorRouter := Api.Group("/monitorDevice")
	monitorRouter.PUT("", controllers.FirstInitMonitorDevice)   // add update
	monitorRouter.POST("/excelupload", controllers.ExcelUpload) // add update
	monitorRouter.PUT("/metrics", controllers.MetricsBind)      // add update (shell,cpu,disk,mem...)
}

func InitModelMonitorDevice() {
	Router := Api.Group("/model/monitorDevice")

	Router.DELETE("", controllers.DeleteDevice) // add update
}

func InitModelMonitorShells() {
	Router := Api.Group("/model/monitorShells")

	// 创建default表记录 不与monitorDevice关联
	Router.PUT("/uploadShell", controllers.UpdateFileMonitorShell) // 上传文件 新增记录
}

func InitModelMonitorAlerts() {
	Router := Api.Group("/model/monitorAlerts")

	//前端
	/*
			1.选择一台主机 添加报警项metricname
		    2.批量添加 shell ,先筛选已经 添加 此监控指标 的主机
	*/
	Router.PUT("", controllers.AlertsBind) // add update

	//只能 在默认采集表中筛选metricname
	Router.PUT("/shell", controllers.AlertsAdd) //  默认报警表
}
