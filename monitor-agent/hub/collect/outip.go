package collect

import (
	"errors"
	"futong-yw-monitor-center/monitor-agent/g"
	"futong-yw-monitor-center/monitor-agent/utils"
	"futong-yw-monitor-center/monitor-base/bg"
	"regexp"
)

/**
 * @Author: Tao Jun
 * @Description: collect
 * @File:  hostinfo
 * @Version: 1.0.0
 * @Date: 2021/11/26 上午10:53
 */

func GetOutBandIPHandle(url string) error {
	for i := 0; i < 2; i++ {
		//解析正则表达式，如果成功返回解释器
		reg1 := regexp.MustCompile(`(http|https)://(.*?):(.*?)/(.*)`)
		if reg1 == nil { //解释失败，返回nil
			//fmt.Println("regexp err")
			return errors.New("regex err")
		}
		//根据规则提取关键信息
		result := reg1.FindAllStringSubmatch(url, -1)
		//fmt.Println("result1 = ", result)
		if len(result) > 0 {
			if len(result[0]) > 3 {
				tip := utils.GetOutboundIP(result[0][2], result[0][3])
				errs := bg.Validate.Var(tip, "required,ip")
				if errs != nil {
					return errs
				} else {
					g.OutIP = tip
					g.GetLog().Debug("Get OutIP :%s success\n", tip)
					return nil
				}
			}
		}
		continue
	}
	return errors.New("Get OutBandIP result 0")
}
