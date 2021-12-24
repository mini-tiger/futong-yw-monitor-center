package g

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @Author: Tao Jun
 * @Description: g
 * @File:  http
 * @Version: 1.0.0
 * @Date: 2021/11/26 下午4:01
 */

var GinRouter *gin.Engine

var HttpSrv = &http.Server{
	Addr:    "0",
	Handler: GinRouter,
}
