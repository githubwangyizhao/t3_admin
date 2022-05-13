package models

import (
	"fmt"
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"path/filepath"
)

// 备份数据库基础用法
func BackDatabaseBase(user, host, pass, dbName, fileName string) (string, error) {
	//数据库名称
	backDir := g.Cfg().GetString("database.backDir", "mysql_back")
	if fileName == "" {
		fileName = GetBackDatabaseName(dbName)
	}
	if pass == "" {
		pass = g.Cfg().GetString("database.game.pass")
	}
	backDir, err := filepath.Abs(backDir)
	if err != nil {
		return "", err
	}
	pathFileName := backDir + "/" + fileName
	g.Log().Info("备份数据库...")
	err = utils.EnsureDir(backDir)
	utils.CheckError(err)
	if err != nil {
		return pathFileName, err
	}
	//cmd := fmt.Sprintf("mysqldump -u root -h %s -p%s  %s > mysql_back/%s_%d.sql", dbConfig.TargetDb.DbHost, dbConfig.TargetDb.DbPwd, dbConfig.TargetDb.DbName, dbConfig.TargetDb.DbName, now)
	cmd := fmt.Sprintf("mysqldump -u %s -h %s -p%s  %s > %s", user, host, pass, dbName, pathFileName)
	err = utils.GfExecShellRun(cmd)
	utils.CheckError(err, "备份数据库失败:"+cmd)
	if err != nil {
		return pathFileName, err
	}
	g.Log().Infof("备份数据库完成:%+v", pathFileName)
	return pathFileName, err
}

// 备份数据库基础用法并压缩zip(带删除原文件)
func BackDatabaseBaseAndZip(user, host, pass, dbName string) (string, error) {
	pathFileName, err := BackDatabaseBase(user, host, pass, dbName, "")
	if err != nil {
		return "", err
	}
	zipFileName, err := utils.CompressorZip(pathFileName)
	if err != nil {
		return zipFileName, err
	}
	err = gfile.Remove(pathFileName)
	if err != nil {
		g.Log().Errorf("移除文件失败：%s  err:%+v", pathFileName, err)
		return zipFileName, err
	}
	return zipFileName, nil
}

// 备份的数据库名
func GetBackDatabaseName(dbName string) string {
	return dbName + "_" + utils.TimeFormFileName() + ".sql"
}
