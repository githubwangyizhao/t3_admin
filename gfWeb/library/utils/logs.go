package utils

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

func InitLog() {
	g.Log().SetFlags(glog.F_TIME_TIME | glog.F_FILE_SHORT) // 时间和 输出文件名:行号
	backDir := g.Cfg().GetString("database.backDir", "mysql_back")
	err := EnsureDir(backDir)
	CheckError(err, "创建数据备份文件夹失败")
	err = EnsureDir(GetShowFileDir())
	CheckError(err, "创建实时文件夹失败")
	//g.Log().New()
	//g.Log()
	//beego.BConfig.Log.AccessglogFormat = ""
	//level := beego.AppConfig.String("glog::level")
	//g.Log().SetLogger(g.Log().AdapterMultiFile, `{"filename":"glog/admin.log",
	//	"separate":["critical", "error", "warning", "info", "debug"],
	//	"level":`+ level+ `,
	//	"daily":true,
	//	"maxdays":10}`)
	//g.Log().Async() //异步
	////输出文件名和行号
	//g.Log().EnableFuncCallDepth(true)
	//g.Log().SetLogFuncCallDepth(3)
}
