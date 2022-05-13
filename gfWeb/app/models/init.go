package models

import (
	"fmt"
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
)

var (
	Db       *gorm.DB
	DbCenter *gorm.DB
	DbCharge *gorm.DB
	//DbLoginServer *gorm.DB
)

//初始化
func init() {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}
	initDb()
	initCenter()
	initCharge()
}

func initDb() {
	//数据库名称
	dbName := g.Cfg().GetString("database.default.name")
	//数据库用户名
	dbUser := g.Cfg().GetString("database.default.user")
	//数据库密码
	dbPwd := g.Cfg().GetString("database.default.pass")
	//数据库IP
	dbHost := g.Cfg().GetString("database.default.host")
	//数据库端口
	dbPort := g.Cfg().GetString("database.default.port")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPwd, dbHost, dbPort, dbName)
	//dsn := "root:game1234@tcp(192.168.31.100:3306)/center?charset=utf8&parseTime=True&loc=Local"
	var err error

	glog.Infof("initDb dsn:%v", dsn)
	Db, err = gorm.Open("mysql", dsn)
	utils.CheckError(err, "连接后台数据库失败")
	//Db.LogMode(true)
	Db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	Db.SingularTable(true)
	Db.DB().SetMaxIdleConns(50)
}

func initCenter() {
	//数据库名称
	dbName := g.Cfg().GetString("database.center.name")
	//数据库用户名
	dbUser := g.Cfg().GetString("database.center.user")
	//数据库密码
	dbPwd := g.Cfg().GetString("database.center.pass")
	//数据库IP
	dbHost := g.Cfg().GetString("database.center.host")
	//数据库端口
	dbPort := g.Cfg().GetString("database.center.port")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPwd, dbHost, dbPort, dbName)
	var err error
	DbCenter, err = gorm.Open("mysql", dsn)
	DbCenter.SingularTable(true)
	utils.CheckError(err, "连接中心服失败")
}

func PingDb(db *gorm.DB) {
	err := db.DB().Ping()
	if err != nil {
		g.Log().Errorf("ping 数据库失败:%v", err)
	}
	//sql := `show databases`
	//err := db.Raw(sql).Error
	//if err != nil {
	//	g.Log().Errorf("ping 数据库失败:%v", err)
	//}
}

func initCharge() {
	//数据库名称
	dbName := g.Cfg().GetString("database.charge.name")
	//数据库用户名
	dbUser := g.Cfg().GetString("database.charge.user", g.Cfg().GetString("database.center.user"))
	//数据库密码
	dbPwd := g.Cfg().GetString("database.charge.pass", g.Cfg().GetString("database.center.pass"))
	//数据库IP
	dbHost := g.Cfg().GetString("database.charge.host", g.Cfg().GetString("database.center.host"))
	//数据库端口
	dbPort := g.Cfg().GetString("database.charge.port", g.Cfg().GetString("database.center.port"))
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPwd, dbHost, dbPort, dbName)
	var err error

	glog.Infof("initCharge dsn:%v", dsn)
	DbCharge, err = gorm.Open("mysql", dsn)
	utils.CheckError(err, "连接充值数据库失败")
	//Db.LogMode(true)
	DbCharge.SetLogger(log.New(os.Stdout, "\r\n", 0))
	DbCharge.SingularTable(true)
}

func TableName(name string) string {
	prefix := g.Cfg().GetString("database.default.prefix")
	return prefix + name
}
