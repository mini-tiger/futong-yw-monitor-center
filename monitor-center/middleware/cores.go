package middleware

import (
	"futong-yw-monitor-center/monitor-center/g"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @Author: Tao Jun
 * @Description: middleware
 * @File:  cores
 * @Version: 1.0.0
 * @Date: 2021/8/18 上午9:23
 */

// xxx 需要将 r.Use(middlewares.Cors()) 在使用路由前进行设置，否则会导致不生效
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		//fmt.Printf("%+v\n",c.Request)
		//fmt.Println(origin)
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			//c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			c.Header("Access-Control-Allow-Headers", "*")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		c.Next()

	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {

		defer func() {
			if err := recover(); err != nil {
				g.GetLog().Error(" [ %s ] 触发错误:%v\n", c.Request.URL.Path, err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": "0",
					"msg":  "失败",
					"data": err,
				})
			}
		}()

		c.Next()

	}
}
