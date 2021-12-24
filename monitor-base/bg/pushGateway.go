package bg

import (
	"net"
	"net/http"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: g
 * @File:  pushGateway
 * @Version: 1.0.0
 * @Date: 2021/11/19 下午4:22
 */
var PushGatewayClient *http.Client

func init() {
	// 创建PushGateway连接池
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 连接超时
			KeepAlive: 30 * time.Second, // 探活时间
		}).DialContext,
		MaxIdleConns:          1000,             // 最大空闲连接
		IdleConnTimeout:       90 * time.Second, // 空闲超时时间
		TLSHandshakeTimeout:   10 * time.Second, // tls握手超时时间
		ExpectContinueTimeout: 1 * time.Second,  // 100-continue状态码超时时间
	}
	// 创建客户端
	PushGatewayClient = &http.Client{
		Timeout:   time.Second * 30, // 请求超时时间
		Transport: transport,
	}
}
