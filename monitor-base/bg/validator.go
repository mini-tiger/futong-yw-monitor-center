package bg

import (
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
)

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  models_vaild
 * @Version: 1.0.0
 * @Date: 2021/11/25 下午12:41
 */
var Validate *validator.Validate

func init() {
	Validate = validator.New()
	Validate.RegisterValidation("OsValidation", OsValidationFunc)
	Validate.RegisterValidation("PatternValidation", PatternValidationFunc)
	//Validate.RegisterValidation("PortValidation", PortValidationFunc)
}

func OsValidationFunc(f1 validator.FieldLevel) bool {
	// f1 包含了字段相关信息
	// f1.Field() 获取当前字段信息
	// f1.Param() 获取tag对应的参数
	// f1.FieldName() 获取字段名称
	//fmt.Printf("%+v\n",f1)
	value := strings.ToLower(f1.Field().String())
	return value == "linux" || value == "windows"
}

func PortValidationFunc(f1 validator.FieldLevel) bool {
	// f1 包含了字段相关信息
	// f1.Field() 获取当前字段信息
	// f1.Param() 获取tag对应的参数
	// f1.FieldName() 获取字段名称
	//fmt.Printf("%+v\n",f1)
	value, err := strconv.Atoi(f1.Field().String())
	if err != nil {
		return false
	}
	return value >= 0 && value < 65535
}

func PatternValidationFunc(f1 validator.FieldLevel) bool {
	//fmt.Printf("%+v\n",f1)
	v := strings.ToLower(f1.Field().String())

	return v == "ssh" || v == "web" || v == "agent"
}
