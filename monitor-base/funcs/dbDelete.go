package funcs

import (
	"github.com/jinzhu/gorm"
)

/**
 * @Author: Tao Jun
 * @Description: funcs
 * @File:  dbData
 * @Version: 1.0.0
 * @Date: 2021/12/5 上午8:43
 */

func DeleteMany(db *gorm.DB, deldata interface{}) error {
	return db.Debug().Unscoped().Delete(deldata).Error

}
