package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @Author: Tao Jun
 * @Description: controllers
 * @File:  base
 * @Version: 1.0.0
 * @Date: 2021/8/24 上午11:48
 */

func ResponseSuccess(c *gin.Context, msg interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": "1",
		"msg":  "成功",
		"data": msg,
	})
}

func ResponseError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "失败",
		"data": err.Error(),
	})
}

func ResponseErrorMsg(c *gin.Context, msg interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "失败",
		"data": msg,
	})
}

func ResponseAuthError(c *gin.Context, err error) {
	//c.JSON(code, gin.H{
	//	"err": err.Error(),
	//})
	ResponseBasic(c, http.StatusUnauthorized, gin.H{"err": err.Error(), "status": http.StatusUnauthorized})
}

func ResponseBasic(c *gin.Context, code int, msg interface{}) {
	//c.JSON(http.StatusOK, msg)
	//fmt.Println(msg)
	c.JSON(code, msg)
}
