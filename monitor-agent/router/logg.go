package router

import (
	"bytes"
	"encoding/json"
	"futong-yw-monitor-center/monitor-agent/g"
	"github.com/gin-gonic/gin"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogMiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {
		//Start timer
		start := time.Now()
		//fmt.Println(start)
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		//raw := c.Request.URL.RawQuery

		// Process request
		// xxx c.Next以上是请求时的运行
		c.Next() // 注意 next()方法的作用是跳过该调用链去直接后面的中间件以及api路由
		// xxx c.Next以下是返回时的运行
		//end:=time.Now()
		//fmt.Println(end)
		//log.Printf("Request :%+v\n", c.Request.Header["X-Forwarded-For"])
		//path := c.Request.URL.Path
		g.GetLog().Warn("%15s | %20s | %15s | X-Forwarded-For: %s | %s | %d | %s\n",
			" RequestLog ",
			c.Request.URL.Path,
			c.ClientIP(),
			c.Request.Header["X-Forwarded-For"],
			c.Request.Method,
			c.Writer.Status(),
			time.Now().Sub(start))

		//fmt.Printf("%+v\n",c)
		//log.Printf("%s | %s | %s | %s | %d | %s",
		//	" RequestLog ",
		//	path,
		//	c.ClientIP(),
		//	c.Request.Method,
		//	c.Writer.Status(),
		//	time.Now().Sub(start))

		//statusCode := c.Writer.Status()
		var bodyMap map[string]interface{} = make(map[string]interface{}, 3)
		_ = json.Unmarshal(blw.body.Bytes(), &bodyMap)
		if code, ok := bodyMap["code"]; ok && code == "0" {
			g.GetLog().Error("%15s | %12s | %15s | X-Forwarded-For: %s | %s | Req: %+v | Resp: %+v\n",
				" ResponseLog ",
				c.Request.URL.Path,
				c.ClientIP(),
				c.Request.Header["X-Forwarded-For"],
				c.Request.Method,
				c.Request.Form,
				bodyMap)
			//c.Writer.Status(),
		}
		//if statusCode >= 400 {
		//	//ok this is an request with error, let's make a record for it
		//	// now print body (or log in your preferred way)
		//	fmt.Println("Response body: " + blw.body.String())
		//}
	}
}
