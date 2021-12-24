package models

import (
	"database/sql"
	"futong-yw-monitor-center/monitor-base/models"
	"futong-yw-monitor-center/monitor-center/g"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	tjfile "github.com/mini-tiger/tjtools/file"

	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var Db *gorm.DB
var SqlDB *sql.DB

func InitDB() {

	var err error
	admin := g.GetConfig().MysqlDsn
	//dsn := admin.Username + ":" + admin.Password + "@tcp(" + admin.Path + ")/" + admin.Dbname + "?" + admin.Config
	//if MysqlClient, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
	//	log.Fatal("mysql init err:", err)
	//}

	//if Db, err = gorm.Open(mysql.New(mysql.Config{
	//	//DSN: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name
	//	DSN: dsn,
	//	//DefaultStringSize: 64, // default size for string fields
	//	DisableDatetimePrecision: true, // disable datetime precision, which not supported before MySQL 5.6
	//	DontSupportRenameIndex: true, // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
	//	DontSupportRenameColumn: true, // `change` when rename column, rename column not supported before MySQL 8, MariaDB
	//	SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	//}), &gorm.Config{}); err != nil {
	//	log.Fatal("mysql init err:", err)
	//}

	Db, err = gorm.Open("mysql", admin)
	if err != nil {
		log.Fatalf("Mysql Conn Err:%s\n", err.Error())
	}

	SqlDB = Db.DB()
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	SqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	SqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	SqlDB.SetConnMaxLifetime(time.Minute)

	if g.Product {
		Db.LogMode(false)
	}

	// 迁移数据库表
	DBTables()
}

func GetDbInfo() sql.DBStats {
	return SqlDB.Stats()
}

// 注册数据库表专用
func DBTables() {
	var err error
	if !Db.HasTable(&models.MonitorMetricsTpl{}) {
		err = Db.CreateTable(&models.MonitorMetricsTpl{}).Error
		if err != nil {
			log.Printf("DB Tables:%s Create Err:%s\n", new(models.MonitorMetrics).TableName(), err.Error())
		}
	}

	if !Db.HasTable(&models.MonitorDevice{}) {
		err = Db.CreateTable(&models.MonitorDevice{}).Error
		if err != nil {
			log.Printf("DB Tables:%s Create Err:%s\n", new(models.MonitorDevice).TableName(), err.Error())
		}
	}

	if !Db.HasTable(&models.MonitorAuth{}) {
		Db.CreateTable(&models.MonitorAuth{})
	}
	if !Db.HasTable(&models.MonitorAgent{}) {
		Db.CreateTable(&models.MonitorAgent{})
	}

	if !Db.HasTable(&models.MonitorAlerts{}) {
		Db.CreateTable(&models.MonitorAlerts{})
	}

	if !Db.HasTable(&models.MonitorMetricsDefault{}) {
		Db.CreateTable(&models.MonitorMetricsDefault{})
	}
	if !Db.HasTable(&models.MonitorAlertsDefault{}) {
		Db.CreateTable(&models.MonitorAlertsDefault{})
	}

	InitData()
}

func InitData() {
	if Db.HasTable(&models.MonitorMetricsDefault{}) {
		//saveData("monitor_metrics.sql")
		for _, value := range models.DefaultMonitorMetricsData() {
			Db.FirstOrCreate(value)
		}
	}

	if Db.HasTable(&models.MonitorMetricsTpl{}) {

		Db.FirstOrCreate(models.DefaultMonitorMetricsTplData())
	}

	if Db.HasTable(&models.MonitorAlertsDefault{}) {
		//saveData("monitor_metrics.sql")
		for _, value := range models.DefaultMonitorMetricsAlerts() {
			Db.FirstOrCreate(value)
		}
	}
	//models.CopyDefaultAlertData(Db)
}

func saveData(file string) {
	baseDir := filepath.Dir(g.CurrentDir)
	modelDir := path.Join(baseDir, "monitor-base", "models")

	SqlFile := path.Join(modelDir, file)

	if !tjfile.IsExist(SqlFile) {
		log.Printf("Sql File:%s Not Exist\n", SqlFile)
		return
	}

	sqls, _ := ioutil.ReadFile(SqlFile)
	sqlArr := strings.Split(string(sqls), ";")
	tx := Db.Begin()
	for _, sql := range sqlArr {
		if sql == "" {
			continue
		}
		//fmt.Printf("%+v\n", sql)
		tx.Exec(sql)
	}
	err := tx.Commit().Error
	if err != nil {
		log.Printf("Sql File %s Err:%v\n", file, err)
	}

}
