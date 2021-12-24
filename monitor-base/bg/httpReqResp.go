package bg

import (
	"github.com/go-resty/resty/v2"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: bg
 * @File:  httpReqResp
 * @Version: 1.0.0
 * @Date: 2021/11/25 下午6:26
 */

var RestyClient *resty.Client

func init() {
	RestyClient = resty.New()

	// Retries are configured per client
	RestyClient.
		//SetLogger(nil). // xxx 修改 默认日志
		SetDebug(false).
		// xxx timeout
		SetTimeout(80 * time.Second).
		// xxx Setting a Proxy URL and Port
		//SetProxy("http://proxyserver:8888")

		// xxx Want to remove proxy setting
		//	client.RemoveProxy().
		// Set retry count to non zero to enable retries
		SetRetryCount(1).
		// You can override initial retry wait time.
		// Default is 100 milliseconds.
		SetRetryWaitTime(100 * time.Millisecond)
	// MaxWaitTime can be overridden as well.
	// Default is 2 seconds.
	//SetRetryMaxWaitTime(20 * time.Second).
	// SetRetryAfter sets callback to calculate wait time between retries.
	// Default (nil) implies exponential backoff with jitter
	//SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
	//	return 0, errors.New("quota exceeded")
	//}).
	//OnError(func(req *resty.Request, err error) { // xxx 请求错误hook
	//	if v, ok := err.(*resty.ResponseError); ok {
	//		// v.Response contains the last response from the server
	//		// v.Err contains the original error
	//		fmt.Printf("111111 %+v\n", v)
	//	}
	//	fmt.Printf("22222 req:%+v err:%+v\n", req, err)
	//	// Log the error, increment a metric, etc...
	//})
}
