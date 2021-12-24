package funcs

import (
	"futong-yw-monitor-center/monitor-center/g"
	"github.com/olivere/elastic/v7"
)

/**
 * @Author: Tao Jun
 * @Description: authfunc
 * @File:  handler
 * @Version: 1.0.0
 * @Date: 2021/8/16 下午3:27
 */

var EsClient *elastic.Client

func InitES() {
	var err error
	EsClient, err = elastic.NewClient(
		elastic.SetURL(g.GetConfig().EsServer...),
		elastic.SetSniff(false), //docker es
	)

	if err != nil {
		g.GetLog().Error("ElasticSearch%v Fail !!!!\n", g.GetConfig().EsServer)
	}
}
