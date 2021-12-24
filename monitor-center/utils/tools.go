package utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/url"
	"path"
	"strings"
)

/**
 * @Author: Tao Jun
 * @Description: utils
 * @File:  tools
 * @Version: 1.0.0
 * @Date: 2021/8/16 下午5:45
 */

func ParamsVerify(c *gin.Context, s string) (string, error) {
	var p string
	var ok bool

	if p, ok = c.GetPostForm(s); ok {
		pLower := strings.ToLower(p)
		if pLower == "undefined" || pLower == "none" || pLower == "nil" {
			return p, errors.New(fmt.Sprintf("request Params [%s] invalid", s))
		} else {
			return p, nil
		}

	} else {
		return p, errors.New(fmt.Sprintf("request Params Miss [%s]", s))
	}

}

func Md5V3(str string) string {
	w := md5.New()
	io.WriteString(w, str)

	return strings.ToUpper(fmt.Sprintf("%x", w.Sum(nil)))
}

func FixReg(c *gin.Context) (regexTpl string) {
	regexTpl = "*%s*"

	if regexFix, o := c.GetPostForm("regexFix"); o {
		regexFix = strings.ToLower(regexFix)
		switch true {
		case regexFix == "left":
			regexTpl = strings.TrimRight(regexTpl, "*")
			break
		case regexFix == "right":
			regexTpl = strings.TrimLeft(regexTpl, "*")
			break
		}
	}
	return
}

func UrlFormat(prefix, pathfile string) string {
	u, _ := url.Parse(prefix)

	u.Path = path.Join(u.Path, pathfile)
	return u.String()
}
