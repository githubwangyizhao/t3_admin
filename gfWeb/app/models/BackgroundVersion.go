package models

import (
	"gfWeb/library/utils"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

// 后台通知信息模板
type BackgroundVersion struct {
	Version    string `json:"version"`
	RunVersion string `json:"runVersion"`
	ChangeTime int64  `json:"changeTime"`
}

// 获取后台版本数据
func GetBackgroundVersion() (*BackgroundVersion, error) {
	cmdResultStr, err :=
		utils.Cmd("ps", []string{
			//"ax -o pid= -o command= | grep 'gfAdmin' | grep -v grep",
			"ax",
			"-o",
			"pid=",
			"-o",
			"command=",
		})
	utils.CheckError(err, "获取后台运行版本")
	pidStr := gconv.String(gproc.Pid())
	RunVersion := "pid:" + pidStr
	for _, cmdStr := range strings.Split(cmdResultStr, "\n") {
		cmdStrL := strings.Split(cmdStr, " ")
		if pidStr == cmdStrL[0] {
			RunVersion += " \tfile:" + cmdStrL[1]
		}
	}
	data := &BackgroundVersion{}
	err = Db.First(&data).Error
	data.RunVersion = RunVersion
	return data, err
}

// 停用后台
func StopBackground() error {
	s := g.Server()
	return s.Shutdown()
}

// 更新后台版本
func UpdateBackgroundVersion() error {
	versionData, err := GetBackgroundVersion()
	currVersion := ""
	versionFileName := ""
	if err == nil {
		currVersionL := strings.Split(versionData.Version, "-")
		if len(currVersionL) == 3 {
			currVersion = versionData.Version
		}
	}
	for _, fileName := range utils.GetCurrDirOrAllFile(".") {
		fileNameL := strings.Split(fileName, "_")
		if len(fileNameL) != 2 {
			continue
		}
		fileNameStr := fileNameL[1]
		fileNameLL1 := strings.Split(fileNameStr, "-")
		if len(fileNameLL1) != 3 {
			continue
		}
		if currVersion != "" && fileNameStr <= currVersion {
			continue
		}
		currVersion = fileNameStr
		versionFileName = fileName

	}
	if versionFileName == "" {
		g.Log().Infof("未找到新版本")
		return gerror.New("未找到新版本")
	}
	g.Log().Infof("更新版本文件:%v pid:%v", versionFileName, gproc.Pid())
	err = ghttp.RestartAllServer(versionFileName)
	utils.CheckError(err)
	if err != nil {
		return gerror.Wrap(err, "重启后台文件失败")
	}
	g.Log().Infof("更新版本文件后:pid:%v", gproc.Pid())
	NewBackgroudVersionData := &BackgroundVersion{
		Version:    currVersion,
		ChangeTime: gtime.Timestamp(),
	}
	err = Db.Save(NewBackgroudVersionData).Error
	return err
}
