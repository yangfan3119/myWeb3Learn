package databases

import (
	"go01_4/logger"

	_ "github.com/mattn/go-sqlite3"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Mdb *gorm.DB
var mlog = logger.GetLogger()

func InitDB(sqlitName string) {
	mlog.Infoln(sqlitName, "数据库连接初始化")
	var err error
	Mdb, err = gorm.Open(sqlite.Open(sqlitName), &gorm.Config{})
	if err != nil {
		mlog.Errorln("failed to connect database.", err.Error())
		Mdb = nil
		return // 如果连接失败，直接返回
	}
	sqlx, _ := Mdb.DB()

	if err := sqlx.Ping(); err != nil {
		mlog.Errorln("测试连接失败！")
	} else {
		mlog.Println("数据库连接成功！")
	}
}

func GetDB() *gorm.DB {
	if Mdb == nil {
		mlog.Errorln("数据库连接未初始化！")
	}
	return Mdb
}
