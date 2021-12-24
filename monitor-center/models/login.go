package models

/**
 * @Author: Tao Jun
 * @Description: modules
 * @File:  login
 * @Version: 1.0.0
 * @Date: 2021/4/15 下午1:06
 */

type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
