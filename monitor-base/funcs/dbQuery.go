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
type GetDbOneRecord struct {
	Params  map[string]interface{} // where sql
	Preload []string               // 预加载的字段
	Result  interface{}            // 返回记录
}

func (g *GetDbOneRecord) MustOneRecord(db *gorm.DB) error {
	exec_db := db.Where(g.Params)
	for _, preload := range g.Preload {
		exec_db = exec_db.Preload(preload)
	}

	return exec_db.Take(g.Result).Error
}
