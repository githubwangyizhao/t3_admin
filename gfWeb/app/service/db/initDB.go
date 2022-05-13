package db

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

var (
	Db       gdb.DB
	DbCenter gdb.DB
	DbCharge gdb.DB
	//DbLoginServer *gorm.DB
)

// 现在是否在开服
//var IsNowOpenServer = false
//var IsNowOpenServerMap map[string]bool
var IsNowOpenServerMap = make(map[string]bool, 0)

//初始化
func init() {
	g.Log().Info("准备连接数据库")
	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	return defaultTableName
	//}
	initDb()
	initCenter()
	initCharge()
	g.Log().Info("完成连接数据库")
}

func initDb() {
	g.Log().Info(".............连接后台数据库")
	Db = g.DB()
	_, err := Db.Query("select * from db_version")
	if err != nil {
		g.Log().Error("连接后台失败:", g.Cfg().Get("database.default.link"))
		return
	}
}

func initCenter() {
	g.Log().Info(".............连接中心服数据库")
	DbCenter = g.DB("center")
	_, err := DbCenter.Query("select * from db_version")
	if err != nil {
		g.Log().Error("连接中心服失败:", g.Cfg().Get("database.center.link"))
		return
	}
}
func initCharge() {
	g.Log().Info(".............连接充值服数据库")
	DbCharge = g.DB("charge")
	_, err := DbCharge.Query("select * from db_version")
	if err != nil {
		g.Log().Error("连接充值服数据库失败:", g.Cfg().Get("database.charge.link"))
		return
	}
}

func PingDb(db gdb.DB) {
	sql := `show databases`
	_, err := db.Exec(sql)
	if err != nil {
		g.Log().Error("ping 数据库失败:%v", err)
	}
}
