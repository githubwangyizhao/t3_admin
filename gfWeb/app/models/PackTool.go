package models

import (
	"fmt"
	"gfWeb/library/utils"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

type PackToolQueryParam struct {
	PackServer string `json:"packServer"`
	PackType   string `json:"packType"`
}

type FtpToolQueryParam struct {
	Ftp      int `json:"packServer"`
	PackType int `json:"packType"`
}

var IsPack = false

// 打包工具
func PackTool(PackServer, PackFile string) (string, error) {
	if IsPack {
		ErrStr := "已经存在打包"
		g.Log().Error(ErrStr)
		return ErrStr, gerror.New(ErrStr)
	}
	IsPack = true
	toolDir := fmt.Sprintf("/opt/branch/%s/filelist/", PackServer)
	//out, err := utils.CmdAndChangeDir(toolDir, "sh", []string{PackFile + ".sh"})
	out, err := utils.GfCmdByDir(toolDir, PackFile+".sh")
	IsPack = false
	return out, err
}

// 同步工具
func SyncTool(SyncDir, SyncPlatform, ShFile string) (string, error) {
	SyncDirStr := SyncDir + SyncPlatform
	if string([]byte(SyncDir)[:1]) != "/" { // 客户端同步
		SyncDirStr = fmt.Sprintf("/opt/branch/%s/filelist/%s", SyncDir, SyncPlatform)
	} else {
		SyncDirStr = fmt.Sprintf("/opt/admin_tool/%s", SyncPlatform)
	}
	//out, err := utils.CmdAndChangeDir(SyncDirStr, "sh", []string{ShFile + ".sh"})
	out, err := utils.GfCmdByDir(SyncDirStr, ShFile+".sh")
	return out, err
}

// 订时更新平台版本
func UpdatePlatformVersionCron(userId, UpdateType, updateEnterStateType, cronUpdateTime1, removeUpdate int, PlatformIdList []string) error {
	list := GetPlatformListByPlatformIdList(PlatformIdList)
	errStr := "没找到平台"
	Len := len(list)
	if Len == 0 {
		return gerror.New(errStr)
	}
	updatePlatformList := make([]Platform, 0)
	currTime := gtime.Timestamp() + 1
	cronUpdateTime := gconv.Int64(cronUpdateTime1)
	for _, platformData := range list {
		platformId := platformData.Id
		Key := GetUpdatePlatformVersionCacheKey(platformId)
		platformValue := utils.GetCacheInt64(Key)
		if platformValue == 0 {
			platformData.UpdateUserId = userId
			if removeUpdate > 0 {
				platformData.UpdateState = 4
				err := Db.Save(platformData).Error
				utils.CheckError(err)
				gcron.Remove(platformData.getCronName())
				continue
			}
			platformData.CronUpdateTime = cronUpdateTime
			platformData.CronUpdateType = UpdateType
			platformData.EnterStateType = updateEnterStateType
			platformData.UpdateState = 1
			//if (Len - 1) == Index {
			//	updatePlatformToolHandle(platformData.Id, platformData.Version, UpdateType)
			//} else {
			if currTime >= cronUpdateTime {
				if updateEnterStateType == 1 || updateEnterStateType == 2 {
					err := BatchPlatformDataUpdateState(platformData, 1)
					utils.CheckError(err, fmt.Sprintf("更新平台版本-更新维护状态失败:%v %v", platformId, err))
					if err != nil {
						continue
					}
					//platformData.EnterState = 1
					//err = Db.Save(platformData).Error
					//utils.CheckError(err)
				} else {
					err := Db.Save(platformData).Error
					utils.CheckError(err, "订时更新平台保存数据失败:"+platformId)
				}
				updatePlatformList = append(updatePlatformList, *platformData)
				continue
			}
			err := Db.Save(platformData).Error
			utils.CheckError(err, "订时更新平台保存数据失败:"+platformId)
			if err != nil {
				continue
			}
			updatePlatformVersionStartCron(platformData)
		}
	}
	if len(updatePlatformList) == 0 {
		if currTime >= cronUpdateTime {
			errStr = "平台正在更新中"
			return gerror.New(errStr)
		} else {
			return nil
		}
	}
	err := RefreshGameServer()
	utils.CheckError(err, "更新平台版本-刷新入口失败")
	if err != nil {
		return gerror.Newf("更新平台版本-刷新入口失败:%+v", err)
	}

	for _, platformData := range updatePlatformList {
		go updatePlatformToolHandle(platformData.Id, 0)
	}
	return err
}

// 获得订时更新版本key
func GetUpdatePlatformVersionCacheKey(platformId string) string {
	return "UpdatePlatformTool" + "_" + platformId
}

// 订时器处理平台更新
func InitPlatformVersionCron() {
	currTime := gtime.Timestamp() + 1
	currTimeInt := gconv.Int(currTime)
	platformList := GetPlatformSimpleList()
	for _, platform := range platformList {
		if platform.OpenServerTakeTime > currTimeInt {
			InitCronPlatformOpenServerTime(platform)
		}
		if platform.UpdateState != 1 {
			continue
		}
		if platform.CronUpdateTime > currTime {
			updatePlatformVersionStartCron(platform)
		}
	}
}

// 平台更新订时器启动
func updatePlatformVersionStartCron(platform *Platform) {
	g.Log().Infof("平台更新订时器启动:%+v", platform.Id)
	gcron.Remove(platform.getCronName())
	cronFun := func() {
		updatePlatformVersionHandle(platform.Id)
	}
	cronTimeStr := utils.TimestampToCronStr(platform.CronUpdateTime)
	_, err := gcron.AddOnce(cronTimeStr, cronFun, platform.getCronName())
	utils.CheckError(err)
	if err != nil {
		g.Log().Error("平台更新订时器启动失败:%+v  err:%+v", platform.Id, err)
	}
}

// 平台更新订时器执行
func updatePlatformVersionHandle(platformId string) {
	platform, err := GetPlatformOne(platformId)
	if err != nil {
		g.Log().Error("平台更新订时器执行更新不存在:%+v err:%+v", platformId, err)
		return
	}
	updateEnterStateType := platform.EnterStateType
	if platform.UpdateState != 1 {
		g.Log().Error("平台更新订时器执行更新不是可更新状态:%+v UpdateState:%+v", platformId, platform.UpdateState)
		return
	}
	if updateEnterStateType == 1 || updateEnterStateType == 2 {
		err := BatchPlatformDataUpdateState(platform, 1)
		utils.CheckError(err, fmt.Sprintf("平台更新订时器执行-更新维护状态失败:%v %v", platformId, err))
		//platform.EnterState = 1
		err = Db.Save(platform).Error
		utils.CheckError(err)
	}
	updatePlatformToolHandle(platformId, 1)
}

//// 更新平台版本
//func UpdatePlatformVersion(userId, UpdateType, updateEnterStateType int, PlatformIdList []string) (string, error) {
//	list := GetPlatformListByPlatformIdList(PlatformIdList)
//	errStr := "没找到平台"
//	Len := len(list)
//	if Len == 0 {
//		return errStr, gerror.New(errStr)
//	}
//	infoStr := ""
//	updatePlatformList := make([]Platform, 0)
//
//	toolDir := utils.GetToolDir()
//	shName := toolDir + "platform_hot_reload.sh"
//	if UpdateType == 1 {
//		shName = toolDir + "platform_cold_reload.sh"
//	}
//	for _, platformData := range list {
//		platformId := platformData.Id
//		Key := getUpdatePlatformVersionCacheKey(platformId, platformData.Version)
//		var valueStr = ""
//		utils.GetCache(Key, &valueStr)
//		if valueStr == "" {
//			//if (Len - 1) == Index {
//			//	updatePlatformToolHandle(platformData.Id, platformData.Version, UpdateType)
//			//} else {
//			if updateEnterStateType == 1 || updateEnterStateType == 2 {
//				err := BatchUpdateState(platformId, 1)
//				utils.CheckError(err, fmt.Sprintf("更新平台版本-更新维护状态失败:%v %v", platformId, err))
//			}
//			updatePlatformList = append(updatePlatformList, *platformData)
//			//go updatePlatformToolHandle(platformId, platformData.Version, shName, userId, updateEnterStateType)
//			//}
//		}
//	}
//	if len(updatePlatformList) == 0 {
//		errStr = "平台正在更新中"
//		return errStr, gerror.New(errStr)
//	}
//	err := RefreshGameServer()
//	utils.CheckError(err, fmt.Sprintf("更新平台版本-刷新入口失败:%v", err))
//	if err != nil {
//		return "", err
//	}
//
//	for _, platformData := range updatePlatformList {
//		go updatePlatformToolHandle(platformData.Id, platformData.Version, shName, userId, 0, updateEnterStateType)
//		infoStr += platformData.Id + "=》" + platformData.Version + "; "
//	}
//	return infoStr, err
//}

// 更新平台工具实现内容
func updatePlatformToolHandle(platformId string, changeType int) {
	g.Log().Infof("%v更新平台工具开始", platformId)
	//Key := "UpdatePlatformTool" + "_" + platformId + "_" + Version

	platform, err := GetPlatformOne(platformId)
	if err != nil {
		g.Log().Error("平台不存在:%+v err:%+v", platformId, err)
		return
	}
	updateType := platform.CronUpdateType
	userId := platform.UpdateUserId
	updateEnterStateType := platform.EnterStateType
	version := platform.Version

	Key := GetUpdatePlatformVersionCacheKey(platformId)
	platformValue := utils.GetCacheInt64(Key)
	if platformValue != 0 {
		g.Log().Warningf("已在更新中key:%+v ：%+v", Key, utils.TimeInt64FormDefault(platformValue/1000))
		return
	}
	toolDir := utils.GetToolDir()
	shName := "platform_hot_reload.sh"
	if updateType == 1 {
		shName = "platform_cold_reload.sh"
	} else if updateType == 10 {
		shName = "platform_tool.sh"
		version = "start"
	} else if updateType == 11 {
		shName = "platform_tool.sh"
		version = "stop"
	}
	startTimeMilli := gtime.TimestampMilli()
	failMsg := ""
	updateState := 5
	utils.SetCache(Key, startTimeMilli, 0)
	defer func() {
		utils.DelCache(Key)
		SaveUpdatePlatformVersionLog(userId, updateType, changeType, updateState, platformId, failMsg, startTimeMilli)
		platform.UpdateState = updateState
		err = Db.Save(platform).Error
		utils.CheckError(err, "保存数据失败")
	}()

	platform.UpdateState = 2
	Db.Save(platform)

	//shPath := path.Dir(shName)
	commandArgs := []string{shName, platformId, version}
	err = utils.CmdAndChangeDirToFile("version_"+platformId+".log", toolDir, "sh", commandArgs)
	utils.CheckError(err, fmt.Sprintf("更新平台工具实现内容失败:%v", platformId))
	if err != nil {
		failMsg = gconv.String(err)
		return
	}
	if updateEnterStateType == 2 {
		platform, _ = GetPlatformOne(platformId)
		err := BatchPlatformDataUpdateState(platform, 3)
		if err != nil {
			return
		}
		utils.CheckError(err, fmt.Sprintf("更新平台工具实现内容-更新火爆状态:%v %v", platformId, err))
		err = RefreshGameServer()
		utils.CheckError(err, fmt.Sprintf("更新平台工具实现内容-同步平台区服状态:%v %v", platformId, err))
		//platform.EnterState = 3
	}
	updateState = 3
	g.Log().Infof("更新完毕===>> %s", platformId)
}

// 平台操作
func ChangePlatformTool(changeType string, platformIdList []string) (string, error) {
	list := GetPlatformListByPlatformIdList(platformIdList)
	errStr := "没找到平台"
	Len := len(list)
	if Len == 0 {
		return errStr, gerror.New(errStr)
	}
	infoStr := ""
	for Index, platformData := range list {
		Key := "ChangePlatformTool" + "_" + platformData.Id
		OldTime := utils.GetCacheInt64(Key)
		if OldTime == 0 {
			if (Len - 1) == Index {
				changePlatformToolHandle(platformData.Id, changeType)
			} else {
				go changePlatformToolHandle(platformData.Id, changeType)
			}
			infoStr += platformData.Id + "=》" + changeType + "; "
		}
	}
	if infoStr == "" {
		errStr = "平台正在更新中"
		return errStr, gerror.New(errStr)
	}
	return infoStr, nil
}

// 平台操作工具实现内容
func changePlatformToolHandle(platformId string, changeType string) {
	g.Log().Info("%v平台操作工具实现内容开始:%d", platformId, changeType)
	Key := "ChangePlatformTool" + "_" + platformId
	utils.SetCache(Key, gtime.Timestamp(), 0)
	defer utils.DelCache(Key)
	toolDir := utils.GetToolDir()
	shName := "platform_tool.sh"
	commandArgs := []string{shName, platformId, changeType}
	err := utils.CmdAndChangeDirToFile("tool_"+platformId+".log", toolDir, "sh", commandArgs)
	utils.CheckError(err, fmt.Sprintf("平台操作工具实现内容失败:%v", platformId))
	g.Log().Info("平台操作完毕===>> %s", platformId)
}
