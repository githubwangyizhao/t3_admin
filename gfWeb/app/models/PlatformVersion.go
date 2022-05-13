package models

import (
	"fmt"
	"gfWeb/library/utils"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

// 版本路径
type BranchPath struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	UserId     int    `json:"userId"`
	ChangeTime int64  `json:"changeTime"`
	IsAdd      int    `json:"isAdd" gorm:"-"`
}

// 平台版本路径
type PlatformVersionPath struct {
	Id         int    `json:"id"`
	Type       int    `json:"type"` // 类型1:客户端 2:服务器端
	Name       string `json:"name"`
	BranchId   int    `json:"branchId"`
	ChangeType int    `json:"changeType"` // 操作类型1:打包 2:同步 3:更新
	State      int    `json:"state"`
	Path       string `json:"path"`
	UserId     int    `json:"userId"`
	ChangeTime int64  `json:"changeTime"`
	IsAdd      int    `json:"isAdd" gorm:"-"`
}

// 版本工具操作
type VersionToolChange struct {
	Id                int    `json:"id"`
	Type              int    `json:"type"` // 渠道操作类型
	Name              string `json:"name"`
	PlatformVersionId int    `json:"platformVersionId"`
	State             int    `json:"state"`
	ShPath            string `json:"shPath"`
	UserId            int    `json:"userId"`
	ChangeTime        int64  `json:"changeTime"`
	IsAdd             int    `json:"isAdd" gorm:"-"`
}

// 定时版本工具操作
type VersionToolChangeCron struct {
	Id           int    `json:"id"`
	State        int    `json:"state"`      // 1:可执行2:已完成3:关闭
	ChangeType   int    `json:"changeType"` // 操作类型1:打包 2:同步 3:更新
	RobotType    int    `json:"robotType"`  // 类型1:客户端 2:服务器端
	ChangeIdStr  string `json:"changeIdStr"`
	ChangeIdList []int  `json:"changeIdList" gorm:"-"`
	CronTime     int64  `json:"cronTime" gorm:"-"`
	CronTimeStr  string `json:"cronTimeStr"`
	SendTimes    int    `json:"sendTimes"`
	LastSendTime int64  `json:"lastSendTime"`
	UserId       int    `json:"userId"`
	ChangeTime   int64  `json:"changeTime"`
	IsAdd        int    `json:"isAdd" gorm:"-"`
	UserName     string `json:"userName" gorm:"-"`
}

// 版本操作层级
type PlatformVersionChildren struct {
	Id       int                       `json:"id"`
	Name     string                    `json:"name"`
	Children []PlatformVersionChildren `json:"children"`
}
type ParamVersionToolChangeCron struct {
	BaseQueryParam
	ChangeType int `json:"changeType"` // 操作类型0:打包 1:同步 2:更新
	RobotType  int `json:"robotType"`  // 机器类型:0=客户端 1=服务器端
}

func PlatformVersionPathTBName() string {
	return "platform_version_path"
}
func VersionToolChangeTBName() string {
	return "version_tool_change"
}

// 获得订时器名字
func (v VersionToolChangeCron) getCronName() string {
	return CRON_NAME_VERSION_TOOL_CHANGE + gconv.String(v.Id)
}

// 获得订时器名字
func getCronVersionToolName(id int) string {
	return CRON_NAME_VERSION_TOOL_CHANGE + gconv.String(id)
}

// changeType:操作类型1:打包 2:同步 3:更新
func GetVersionToolChangeInfo(changeType, robotType int) []*PlatformVersionChildren {
	changeDataList := make([]*VersionToolChange, 0)
	sql := fmt.Sprintf(`SELECT * FROM %s AS T0 where T0.state = 1 and platform_version_id in (SELECT id FROM %s AS T1 where T1.state = 1 and T1.change_type = ? and T1.type = ? )`, VersionToolChangeTBName(), PlatformVersionPathTBName())
	err := Db.Raw(sql, changeType, robotType).Find(&changeDataList).Error
	utils.CheckError(err)
	resultList := make([]*PlatformVersionChildren, 0)
	for _, e := range changeDataList {
		platformVersionPath, err := GetPlatformVersionPathOne(e.PlatformVersionId)
		versionName := ""
		if err == nil {
			versionName = platformVersionPath.Name + ">"
		}
		Children := &PlatformVersionChildren{
			Id:   e.Id,
			Name: versionName + e.Name,
		}
		resultList = append(resultList, Children)
	}
	return resultList
}

// robotType:0=客户端 1=服务器端
func GetPlatformVersionByChangeTypeOrVersionType(changeType, robotType int) []*PlatformVersionChildren {
	changeDataList := make([]*VersionToolChange, 0)
	//sql := fmt.Sprintf(`SELECT T0.* FROM %s AS T0 INNER JOIN %s AS T1 ON T0.platform_version_id = T1.id WHERE T0.state = T1.state and T0.state = 1 and T0.change_type = ? and T1.type = ?`, VersionToolChangeTBName(), PlatformVersionPathTBName())
	// 	select * from platform_version_change AS T0 where T0.state = 0 and platform_version_id in (select id from platform_version_path as T1 where T1.state = 0 and change_type = 1 and type = 0);
	sql := fmt.Sprintf(`SELECT * FROM %s AS T0 where T0.state = 1 and platform_version_id in (SELECT id FROM %s AS T1 where T1.state = 1 and T1.change_type = ? and T1.type = ? )`, VersionToolChangeTBName(), PlatformVersionPathTBName())
	err := Db.Raw(sql, changeType, robotType).Find(&changeDataList).Error
	//mapChildren := gmap.New()
	utils.CheckError(err)

	platformVersionList := make(map[int][]PlatformVersionChildren)
	for _, e := range changeDataList {
		Children := PlatformVersionChildren{
			Id:   e.Id,
			Name: e.Name,
		}
		tmpL, ok := platformVersionList[e.PlatformVersionId]
		if ok {
			platformVersionList[e.PlatformVersionId] = append(tmpL, Children)
		} else {
			platformVersionList[e.PlatformVersionId] = append([]PlatformVersionChildren{}, Children)
		}
		//handlePlatformVersionLoopChildren(mapChildren, e.PlatformVersionId, Children)
	}
	branchChildrenList := make(map[int][]PlatformVersionChildren)
	for platformVersionId, childrenPlatformVersionList := range platformVersionList {
		platformVersionPath, err := GetPlatformVersionPathOne(platformVersionId)
		if err != nil || platformVersionPath.State != 1 {
			//g.Log().Warningf("获取平台版本路径数据失败:%+v  err:%+v", mapChildrenKey, err)
			g.Log().Warningf("获取平台版本路径数据失败:%+v  err:%+v", platformVersionId, err)
			continue
		}
		Children := PlatformVersionChildren{
			Id:       platformVersionPath.Id,
			Name:     platformVersionPath.Name,
			Children: childrenPlatformVersionList,
		}
		tmpL, ok := branchChildrenList[platformVersionPath.BranchId]
		if ok {
			branchChildrenList[platformVersionPath.BranchId] = append(tmpL, Children)
		} else {
			branchChildrenList[platformVersionPath.BranchId] = append([]PlatformVersionChildren{}, Children)
		}
		//handlePlatformVersionLoopChildren(mapBranchChildren, platformVersionPath.BranchId, Children)
	}
	//var resultList  []PlatformVersionChildren
	resultList := make([]*PlatformVersionChildren, 0)
	for branchId, childrenBranchList := range branchChildrenList {
		branch, err := GetBranchPathOne(branchId)
		if err != nil {
			g.Log().Warningf("获取平台版本路径数据失败:%+v  err:%+v", branchId, err)
			continue
		}
		Children := &PlatformVersionChildren{
			Id:       branch.Id,
			Name:     branch.Name,
			Children: childrenBranchList,
		}
		resultList = append(resultList, Children)
	}
	return resultList
}

// 获得平台版本路径列表
func GetBranchPathOne(branchId int) (*BranchPath, error) {
	branch := &BranchPath{
		Id: branchId,
	}
	err := Db.Where(branch).First(branch).Error
	return branch, err
}

// 获得版本路径列表
func GetBranchPathList(params *BaseQueryParam) ([]*BranchPath, int64) {
	data := make([]*BranchPath, 0)
	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "id")
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	err := Db.Model(&BranchPath{}).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	return data, count
}

// 更新版本路径数据
func EditBranchPath(data *BranchPath) error {
	data.ChangeTime = gtime.Timestamp()
	err := Db.Save(data).Error
	return err
}

// 删除版本路径数据
func DeleteBranchPath(ids []int) error {
	err := Db.Where(ids).Delete(&BranchPath{}).Error
	return err
}

// 获得平台版本路径
func GetPlatformVersionPathOne(platformVersionId int) (*PlatformVersionPath, error) {
	platformVersionPath := &PlatformVersionPath{
		Id: platformVersionId,
	}
	err := Db.Where(platformVersionPath).First(platformVersionPath).Error
	return platformVersionPath, err
}

// 获得平台版本路径列表
func GetPlatformVersionPathList(params *BaseQueryParam) ([]*PlatformVersionPath, int64) {
	data := make([]*PlatformVersionPath, 0)
	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "id")
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	err := Db.Model(&PlatformVersionPath{}).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	return data, count
}

// 更新平台版本路径数据
func EditPlatformVersionPath(data *PlatformVersionPath) error {
	data.ChangeTime = gtime.Timestamp()
	err := Db.Save(data).Error
	return err
}

// 删除平台版本路径数据
func DeletePlatformVersionPath(ids []int) error {
	err := Db.Where(ids).Delete(&PlatformVersionPath{}).Error
	return err
}

// 获得版本工具操作数据
func GetVersionToolChangeOne(changeId int) (*VersionToolChange, error) {
	versionToolChange := &VersionToolChange{
		Id: changeId,
	}
	err := Db.Where(versionToolChange).First(versionToolChange).Error
	return versionToolChange, err
}

// 获得版本工具操作数据
func GetVersionToolChangeOnChangeName(changeId int) string {
	versionToolChange, err := GetVersionToolChangeOne(changeId)
	if err != nil {
		g.Log().Warningf("未找到版本工具数据:%+v", changeId)
		return ""
	}
	platformVersionPath, err := GetPlatformVersionPathOne(versionToolChange.PlatformVersionId)
	if err != nil {
		g.Log().Warningf("未找到版本路径数据:%+v", versionToolChange.PlatformVersionId)
		return "()" + versionToolChange.Name
	}
	return "(" + platformVersionPath.Name + ")" + versionToolChange.Name
}

// 获得版本工具操作列表
func GetVersionToolChangeList(params *BaseQueryParam) ([]*VersionToolChange, int64) {
	data := make([]*VersionToolChange, 0)
	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "id")
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	err := Db.Model(&VersionToolChange{}).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	return data, count
}

// 获得版本工具根据机器类型和操作类型
func GetVersionToolChangeListByRobotTypeAndChangeType(robotType, changeType int) ([]*VersionToolChange, error) {
	data := make([]*VersionToolChange, 0)
	sql := fmt.Sprintf(`SELECT * FROM %s AS T0 where T0.state = 1 and platform_version_id in (SELECT id FROM %s AS T1 where T1.state = 1 and T1.type = ? and T1.change_type = ? )`, VersionToolChangeTBName(), PlatformVersionPathTBName())
	err := Db.Raw(sql, robotType, changeType).Find(&data).Error
	//platformVersionPath := &PlatformVersionPath{
	//	Type: robotType,
	//	ChangeType: changeType,
	//}
	//err := Db.Where(platformVersionPath).First(platformVersionPath).Error
	return data, err
}

// 更新版本工具路径数据
func EditVersionToolChangePath(data *VersionToolChange) error {
	data.ChangeTime = gtime.Timestamp()
	err := Db.Save(data).Error
	return err
}

// 删除版本工具数据
func DeleteVersionToolChange(ids []int) error {
	err := Db.Where(ids).Delete(&VersionToolChange{}).Error
	return err
}

// 发送版本工具操作处理
func SendVersionToolChange(userId int, ids []int) error {
	for _, id := range ids {
		Key := GetSendVersionToolChangeKey(id)
		OldTime := utils.GetCacheInt64(Key)
		if OldTime != 0 {
			errStr := "正在处理中"
			return gerror.New(errStr)
		}
	}
	go handleSendVersionToolChange(userId, 1, ids)
	return nil
	//for _, id := range ids {
	//	Key := GetSendVersionToolChangeKey(id)
	//	var valueStr = ""
	//	utils.GetCache(Key, &valueStr)
	//	if valueStr != "" {
	//		errStr := "正在处理中"
	//		return errStr, gerror.New(errStr)
	//	}
	//	startTimeMilli := gtime.TimestampMilli()
	//	utils.SetCache(Key, Key, 0)
	//	defer utils.DelCache(Key)
	//	 changeData, err := GetVersionToolChangeOne(id)
	//	if err != nil {
	//		g.Log().Warningf("未找到操作数据:%d", id)
	//		continue
	//	}
	//	platformVersionData, err := GetPlatformVersionPathOne(changeData.PlatformVersionId)
	//	if err != nil {
	//		g.Log().Warningf("未找到平台版本路径:%d", changeData.PlatformVersionId)
	//		continue
	//	}
	//	branchData, err := GetBranchPathOne(platformVersionData.BranchId)
	//	if err != nil {
	//		g.Log().Warningf("未找到版本路径:%d", platformVersionData.BranchId)
	//		continue
	//	}
	//	dirStr := branchData.Path + platformVersionData.Path
	//	for _, shName := range strings.Split(changeData.ShPath, " "){
	//		if len(shName) < 1 {
	//			g.Log().Debugf("包含内容有问题%+v:%+v", shName, changeData.ShPath)
	//			continue
	//		}
	//		//out, err := utils.CmdShByDirOrParam(dirStr, []string{shName})
	//		err := utils.CmdShShowFileByDirOrParam("platformVersion_"+gconv.String(id)+"_"+shName+".log", dirStr, []string{shName})
	//		if err != nil {
	//			flagErr = err
	//			_ = SaveVersionToolChangeLog(userId, id, 1, updateType, shName, gconv.String(err), startTimeMilli)
	//			continue
	//		}
	//		_ = SaveVersionToolChangeLog(userId, id, 0, updateType, shName, "", startTimeMilli)
	//		successStr += platformVersionData.Name + " ==>> "
	//	}
	//	//shName :=  branchData.Path + platformVersionData.Path + changeData.ShPath
	//}
	//if flagErr == nil {
	//	return successStr, flagErr
	//}
	//return flagStr, flagErr
}

// 处理版本操作处理
func handleSendVersionToolChange(userId, updateType int, ids []int) {
	g.Log().Infof("%+v开始处理版本操作处理userId:%+v updateType:%+v ids:%+v", gtime.Datetime(), userId, updateType, ids)
	changeStr := ""
	for _, id := range ids {
		Key := GetSendVersionToolChangeKey(id)
		OldTime := utils.GetCacheInt64(Key)
		if OldTime != 0 {
			g.Log().Warningf("正在处理中:%+v", Key)
			continue
		}
		startTimeMilli := gtime.TimestampMilli()
		utils.SetCache(Key, startTimeMilli, 0)
		defer utils.DelCache(Key)
		changeData, err := GetVersionToolChangeOne(id)
		if err != nil {
			g.Log().Warningf("未找到操作数据:%d", id)
			continue
		}
		platformVersionData, err := GetPlatformVersionPathOne(changeData.PlatformVersionId)
		if err != nil {
			g.Log().Warningf("未找到平台版本路径数据:%d", changeData.PlatformVersionId)
			continue
		}
		branchData, err := GetBranchPathOne(platformVersionData.BranchId)
		if err != nil {
			g.Log().Warningf("未找到版本路径数据:%d", platformVersionData.BranchId)
			continue
		}
		dirStr := branchData.Path + platformVersionData.Path
		for _, shName := range strings.Split(changeData.ShPath, " ") {
			if len(shName) < 1 {
				g.Log().Debugf("包含内容有问题%+v:%+v", shName, changeData.ShPath)
				continue
			}
			//out, err := utils.CmdShByDirOrParam(dirStr, []string{shName})

			logData, err := SaveVersionToolChangeLogInit(userId, id, 1, updateType, shName)
			if err != nil {
				g.Log().Warningf("保存版本工具操作日志初始数据失败%+v:%+v", shName, changeData.ShPath)
				continue
			}
			err = utils.CmdShShowFileByDirOrParam("platformVersion_"+gconv.String(logData.Id)+"_"+shName+".log", dirStr, []string{shName})
			if err != nil {
				err = UpdateVersionToolChangeLog(logData, 5, gconv.String(err), startTimeMilli)
				utils.CheckError(err)
				continue
			}
			err = UpdateVersionToolChangeLog(logData, 9, "", startTimeMilli)
			utils.CheckError(err)
			if changeStr != "" {
				changeStr += " ; "
			}
			changeStr += platformVersionData.Name + ":" + changeData.Name
		}
		//shName :=  branchData.Path + platformVersionData.Path + changeData.ShPath
	}

	g.Log().Infof("版本操作处理结束:%+v %+v", changeStr, gtime.Datetime())
}

func GetSendVersionToolChangeKey(id int) string {
	return "SendVersionToolChange" + "_" + gconv.String(id)
}

// 获取单个订时版本工具操作
func GetVersionToolChangeCron(id int) (*VersionToolChangeCron, error) {
	data := &VersionToolChangeCron{
		Id: id,
	}
	err := Db.First(&data).Error
	data.ChangeIdList = gconv.Ints(strings.Split(data.ChangeIdStr, ","))
	return data, err
}

// 获取单订时版本工具操作列表
func GetVersionToolChangeCronList(params *ParamVersionToolChangeCron) ([]*VersionToolChangeCron, int64) {
	data := make([]*VersionToolChangeCron, 0)
	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "id")
	if params.Order == "ascending" {
		sortOrder = sortOrder + " asc"
	} else {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	versionToolChangeCron := &VersionToolChangeCron{
		ChangeType: params.ChangeType,
		RobotType:  params.RobotType,
	}
	err := Db.Model(versionToolChangeCron).Where(versionToolChangeCron).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		e.ChangeIdList = gconv.Ints(strings.Split(e.ChangeIdStr, ","))
		e.UserName = GetUserName(e.UserId)
	}
	return data, count
}

// 删除订时版本工具
func DeleteVersionToolChangeCron(ids []int) error {
	for _, id := range ids {
		gcron.Remove(getCronVersionToolName(id))
	}
	err := Db.Where(ids).Delete(&VersionToolChangeCron{}).Error
	return err
}

// 初始时订时版本工具操作
func InitVersionToolChangeCron(data *VersionToolChangeCron) {
	gcron.Remove(data.getCronName())

	if data.State >= 2 {
		g.Log().Warningf("初始时订时版本工具操作已结束:%+v", data.Id)
		return
	}
	if len(data.ChangeIdStr) == 0 {
		g.Log().Warningf("初始时订时版本工具操作内容为空:%+v", data.Id)
		return
	}
	g.Log().Infof("订时版本工具操作开始CronTimeStr：%+v name:%+v", data.CronTimeStr, data.getCronName())
	g.Log().Debugf("订时版本工具操作开始CronTimeStr:%T", data.CronTimeStr)
	cronFun := func() {
		StartVersionToolChangeCron(data.Id)
	}
	//if data.CronTimes > 0 {
	//	_, err := gcron.AddTimes(data.CronTimeStr, data.CronTimes, cronFun, data.getCronName())
	//	utils.CheckError(err)
	//	return
	//} else {
	_, err := gcron.AddOnce(data.CronTimeStr, cronFun, data.getCronName())
	utils.CheckError(err)
	return
	//}
}

// 开始订时版本工具操作
func StartVersionToolChangeCron(id int) {
	data, err := GetVersionToolChangeCron(id)
	if err != nil {
		g.Log().Warningf("不存在当前订时版本工具操作数据:%+v", id)
		return
	}
	if data.State >= 2 {
		g.Log().Warningf("订时版本工具操作已完成:%+v", data.Id)
		return
	}
	go handleSendVersionToolChange(data.UserId, 2, data.ChangeIdList)
	//go handleSendVersionToolChange(userId, 2, ids)
	data.LastSendTime = gtime.Timestamp()
	data.SendTimes++
	data.State = 2
	err = Db.Save(&data).Error
	utils.CheckError(err, "保存公告日志失败")

}
