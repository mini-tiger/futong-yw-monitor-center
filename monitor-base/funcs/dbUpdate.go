package funcs

import (
	"github.com/jinzhu/gorm"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: funcs
 * @File:  dbData
 * @Version: 1.0.0
 * @Date: 2021/12/5 上午8:43
 */

type UpdateDbManyRecordEntry struct {
	QueryParams  interface{}            // where sql
	UpdateParams map[string]interface{} //
}

func (u *UpdateDbManyRecordEntry) RetryUpdateCol(db *gorm.DB) error {
	var err error
	for i := 0; i < 2; i++ {
		if err = u.UpdateManyCol(db); err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	return err
}

func (u *UpdateDbManyRecordEntry) UpdateManyCol(db *gorm.DB) error {
	return db.Model(u.QueryParams).Updates(u.UpdateParams).Error

}
