package models

import (
	"gfWeb/library/enums"
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/jinzhu/gorm"
	"github.com/mssola/user_agent"
)

// 请求日志
type RequestAdminLog struct {
	Id         int     `json:"id"`
	UserId     int     `json:"userId"`
	UseTime    float64 `json:"useTime"`
	Url        string  `json:"url"`
	ParamStr   string  `json:"paramStr"`
	ChangeTime int64   `json:"changeTime"`
	LoginType  int     `json:"loginType"` // 登录类型1:登录2:未登录',
	State      int     `json:"state"`
	FailMsg    string  `json:"failMsg"`
	UserName   string  `json:"userName" gorm:"-"`
}

// 登录日志
type LoginAdminLog struct {
	Id         int    `json:"id"`
	Account    string `json:"account"`
	UserName   string `json:"userName"`
	Ip         string `json:"ip"`
	LoginType  int    `json:"loginType"` // '登录类型1:登录2:刷新',
	Explorer   string `json:"explorer"`
	Os         string `json:"os"`
	State      int    `json:"state"`
	FailMsg    string `json:"failMsg"`
	ChangeTime int64  `json:"changeTime"`
	IpName     string `json:"ipName" gorm:"-"`
}

// 开服管理日志
type OpenServerManageLog struct {
	Id         int    `json:"id"`
	PlatformId string `json:"platformId"`
	UserId     int    `json:"userId"`
	ChangeStr  string `json:"changeStr"`
	ChangeTime int64  `json:"changeTime"`
	UserName   string `json:"userName" gorm:"-"`
}

// 平台版本更新日志
type UpdatePlatformVersionLog struct {
	Id         int     `json:"id"`
	PlatformId string  `json:"platformId"`
	UserId     int     `json:"userId"`
	State      int     `json:"state"`
	Type       int     `json:"type"` // 更新类型0:手动1:订时
	UseTime    float64 `json:"useTime"`
	FailMsg    string  `json:"failMsg"`
	ChangeTime int64   `json:"changeTime"`
	UserName   string  `json:"userName" gorm:"-"`
}

// 版本工具操作日志
type VersionToolChangeLog struct {
	Id         int     `json:"id"`
	UserId     int     `json:"userId"`
	State      int     `json:"state"`
	UpdateType int     `json:"updateType"` // 更新类型1:手动,2:订时
	ChangeId   int     `json:"changeId"`
	ShName     string  `json:"shName"`
	UseTime    float64 `json:"useTime"`
	FailMsg    string  `json:"failMsg"`
	ChangeTime int64   `json:"changeTime"`
	UserName   string  `json:"userName" gorm:"-"`
	ChangeName string  `json:"changeName" gorm:"-"`
}

// ----------请求参数-----------
// 参数请求日志
type ParamsRequestAdmin struct {
	BaseQueryParam
	UserId    int    `json:"userId"`
	Url       string `json:"url"`
	LoginType int    `json:"loginType"`
	StartTime int    `json:"startTime"`
	EndTime   int    `json:"endTime"`
}

// 参数登录日志
type ParamsLoginAdmin struct {
	BaseQueryParam
	Account string `json:"account"`
	Ip      string `json:"ip"`
	State   int    `json:"state"`
}

// 参数登录日志
type ParamsOpenServerManage struct {
	BaseQueryParam
	PlatformIdList []string `json:"platformIdList"`
	PlatformId     string   `json:"platformId"`
	UserId         int      `json:"userId"`
}

// 参数平台版本更新日志
type ParamsUpdatePlatformVersion struct {
	BaseQueryParam
	PlatformIdList []string `json:"platformIdList"`
	PlatformId     string   `json:"platformId"`
	UserId         int      `json:"userId"`
	State          int      `json:"state"`
}

// 参数版本工具操作日志
type ParamsVersionTool struct {
	BaseQueryParam
	RobotType  int `json:"robotType"`
	ChangeType int `json:"changeType"`
	State      int `json:"state"`
}

// 获取请求日志列表
func GetRequestAdminLogList(params *ParamsRequestAdmin) ([]*RequestAdminLog, int64) {
	data := make([]*RequestAdminLog, 0)
	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "id")
	if params.Order == "ascending" {
		sortOrder = sortOrder + " asc"
	} else {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	tempF := func(db *gorm.DB) *gorm.DB {
		if params.UserId > 0 {
			db = db.Where(&RequestAdminLog{UserId: params.UserId})
		}
		if len(params.Url) > 0 {
			db = db.Where(&RequestAdminLog{Url: params.Url})
		}
		if params.LoginType > 0 {
			db = db.Where(&RequestAdminLog{LoginType: params.LoginType})
		}
		if params.StartTime > 0 && params.EndTime > 0 {
			db = db.Where("? <= change_time AND change_time <= ?", params.StartTime, params.EndTime)
		}
		return db
	}
	err := tempF(Db.Model(&RequestAdminLog{})).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		if e.UserId == 0 {
			e.UserName = "不存在：" + gconv.String(e.UserId)
			continue
		}
		e.UserName = GetUserName(e.UserId)
	}
	return data, count
}

// 保存请求日志
func SaveRequestAdminLog(userId, loginType int, code enums.ResultCode, msg string, r *ghttp.Request) {
	url := r.URL.Path
	//g.Log().Debugf("保存请求日志:%+v  %+v", u, url)
	//if url == "/log/request_log" { // 防止和分页冲突一直刷新
	//	return
	//}
	if url == "/tool/show_file" { // 无意义请求
		return
	}
	if code == enums.CodeSuccess {
		msg = ""
	}
	requestAdminLog :=
		&RequestAdminLog{
			UserId:     userId,
			UseTime:    gconv.Float64(gtime.TimestampMilli()-r.EnterTime) / 1000,
			Url:        url,
			ParamStr:   r.GetBodyString(),
			ChangeTime: gtime.Timestamp(),
			LoginType:  loginType,
			State:      gconv.Int(code),
			FailMsg:    msg,
		}
	err := Db.Save(requestAdminLog).Error
	if err != nil {
		g.Log().Errorf("保存请求日志失败:%+v", err)
	}
}

// 获取登录日志列表
func GetLoginAdminLogList(params *ParamsLoginAdmin) ([]*LoginAdminLog, int64) {
	data := make([]*LoginAdminLog, 0)
	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "id")
	if params.Order == "ascending" {
		sortOrder = sortOrder + " asc"
	} else {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	tempF := func(db *gorm.DB) *gorm.DB {
		if len(params.Account) > 0 {
			db = db.Where(&LoginAdminLog{Account: params.Account})
		}
		if len(params.Ip) > 0 {
			db = db.Where(&LoginAdminLog{Ip: params.Ip})
		}
		if params.State > 0 {
			db = db.Where(&LoginAdminLog{State: params.State})
		}
		return db
	}
	err := tempF(Db.Model(&LoginAdminLog{})).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		e.IpName = e.Ip
	}
	return data, count
}

// 保存登录日志
func SaveLoginAdminLog(u User, userAgent string, r *ghttp.Request, code enums.ResultCode, msg string, data interface{}) {
	url := r.URL.Path
	if url != "/login" && url != "/dashboard" {
		return
	}
	loginType := 1
	if url == "/dashboard" {
		loginType = 2
	}
	if code == enums.CodeSuccess {
		msg = ""
	} else {
		msg += ":" + gconv.String(data)
	}
	ua := user_agent.New(userAgent)
	explorer, _ := ua.Browser()
	//g.Log().Debugf("expVersion : %+v: %+v", explorer, expVersion)
	requestAdminLog :=
		&LoginAdminLog{
			Account:    u.Account,
			UserName:   u.Name,
			Ip:         r.GetClientIp(),
			LoginType:  loginType,
			Explorer:   explorer,
			Os:         ua.OS(),
			State:      gconv.Int(code),
			FailMsg:    msg,
			ChangeTime: gtime.Timestamp(),
		}
	err := Db.Save(requestAdminLog).Error
	if err != nil {
		g.Log().Errorf("保存登录日志失败:%+v", err)
	}
}

// 获取开服日志列表
func GetOpenServerManageLogList(params *ParamsOpenServerManage) ([]*OpenServerManageLog, int64) {
	data := make([]*OpenServerManageLog, 0)
	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "id")
	if params.Order == "ascending" {
		sortOrder = sortOrder + " asc"
	} else {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	tempF := func(db *gorm.DB) *gorm.DB {
		if len(params.PlatformId) > 0 {
			db = db.Where(&OpenServerManageLog{PlatformId: params.PlatformId})
		} else if len(params.PlatformIdList) > 0 {
			db = db.Where("platform_id in (?)", params.PlatformIdList)
		}
		if params.UserId > 0 {
			db = db.Where(&OpenServerManageLog{UserId: params.UserId})
		}
		return db
	}
	err := tempF(Db.Model(&OpenServerManageLog{})).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		e.UserName = GetUserName(e.UserId)
	}
	return data, count
}

// 保存开服管理日志
func SaveOpenServerManageLog(userId int, oldPlatform, newPlatform *Platform) error {
	changeStr := ""
	F := func(changeStr1, name, a1, a2 string) string {
		if changeStr1 != "" {
			changeStr1 += ";"
		}
		changeStr1 += name + ":" + a1 + "=>" + a2
		return changeStr1
	}
	if oldPlatform.CreateRoleLimit != newPlatform.CreateRoleLimit {
		g.Log().Debugf("CreateRoleLimit:%+v:%+v", oldPlatform.CreateRoleLimit, newPlatform.CreateRoleLimit)
		changeStr = F(changeStr, "创角人数开服", gconv.String(oldPlatform.CreateRoleLimit), gconv.String(newPlatform.CreateRoleLimit))
		//changeStr += gconv.String(oldPlatform.CreateRoleLimit) + "=>" + gconv.String(newPlatform.CreateRoleLimit)
	}
	if oldPlatform.OpenServerTakeTime != newPlatform.OpenServerTakeTime {
		changeStr = F(changeStr, "定时开服时间", gconv.String(oldPlatform.OpenServerTakeTime), gconv.String(newPlatform.OpenServerTakeTime))
		//changeStr += gconv.String(oldPlatform.OpenServerTakeTime) + "=>" + gconv.String(newPlatform.OpenServerTakeTime)
	}
	if oldPlatform.IntervalInitTime != newPlatform.IntervalInitTime {
		changeStr = F(changeStr, "间隔初始时间", gconv.String(oldPlatform.IntervalInitTime), gconv.String(newPlatform.IntervalInitTime))
		//changeStr += gconv.String(oldPlatform.IntervalInitTime) + "=>" + gconv.String(newPlatform.IntervalInitTime)
	}
	if oldPlatform.IntervalDay != newPlatform.IntervalDay {
		changeStr = F(changeStr, "间隔开服天", gconv.String(oldPlatform.IntervalDay), gconv.String(newPlatform.IntervalDay))
		//changeStr += gconv.String(oldPlatform.IntervalDay) + "=>" + gconv.String(newPlatform.IntervalDay)
	}
	if oldPlatform.OpenServerTimeScope != newPlatform.OpenServerTimeScope {
		changeStr = F(changeStr, "开服时间段", gconv.String(oldPlatform.OpenServerTimeScope), gconv.String(newPlatform.OpenServerTimeScope))
		//changeStr += gconv.String(oldPlatform.OpenServerTimeScope) + "=>" + gconv.String(newPlatform.OpenServerTimeScope)
	}
	if oldPlatform.ServerAliasStr != newPlatform.ServerAliasStr {
		changeStr = F(changeStr, "区服别名", gconv.String(oldPlatform.ServerAliasStr), gconv.String(newPlatform.ServerAliasStr))
		//changeStr += gconv.String(oldPlatform.ServerAliasStr) + "=>" + gconv.String(newPlatform.ServerAliasStr)
	}
	if len(changeStr) == 0 {
		return nil
	}
	err := SaveOpenServerManageLogData(userId, newPlatform.Id, changeStr)
	return err
}

// 保存开服管理日志数据
func SaveOpenServerManageLogData(userId int, PlatformId, changeStr string) error {
	logData :=
		&OpenServerManageLog{
			UserId:     userId,
			PlatformId: PlatformId,
			ChangeStr:  changeStr,
			ChangeTime: gtime.Timestamp(),
		}
	err := Db.Save(logData).Error
	if err != nil {
		g.Log().Errorf("保存开服日志日志失败:%+v", err)
		return err
	}
	return err
}

// 获取平台版本更新日志列表
func GetUpdatePlatformVersionLogList(params *ParamsUpdatePlatformVersion) ([]*UpdatePlatformVersionLog, int64) {
	data := make([]*UpdatePlatformVersionLog, 0)
	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "id")
	if params.Order == "ascending" {
		sortOrder = sortOrder + " asc"
	} else {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	tempF := func(db *gorm.DB) *gorm.DB {
		if params.UserId > 0 {
			db = db.Where(&UpdatePlatformVersionLog{UserId: params.UserId})
		}
		if params.State > 0 {
			db = db.Where(&UpdatePlatformVersionLog{State: params.State})
		}
		if len(params.PlatformId) > 0 {
			db = db.Where(&UpdatePlatformVersionLog{PlatformId: params.PlatformId})
		} else if len(params.PlatformIdList) > 0 {
			db = db.Where("platform_id in (?)", params.PlatformIdList)
		}
		return db
	}
	err := tempF(Db.Model(&UpdatePlatformVersionLog{})).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		e.UserName = GetUserName(e.UserId)
	}
	return data, count
}

// 保存平台版本更新日志数据
func SaveUpdatePlatformVersionLog(userId, updateType, changeType, state int, PlatformId, failMsg string, startTimeMilli int64) error {
	logData :=
		&UpdatePlatformVersionLog{
			UserId:     userId,
			PlatformId: PlatformId,
			State:      state,
			Type:       changeType,
			UseTime:    gconv.Float64(gtime.TimestampMilli()-startTimeMilli) / 1000,
			FailMsg:    failMsg,
			ChangeTime: gtime.Timestamp(),
		}
	err := Db.Save(logData).Error
	if err != nil {
		g.Log().Errorf("保存平台版本更新日志数据失败:%+v", err)
		return err
	}
	updateStr := "热更"
	if updateType == 1 {
		updateStr = "冷更"
	} else if updateType == 10 {
		updateStr = "开启"
	} else if updateType == 11 {
		updateStr = "关闭"
	}
	if len(failMsg) == 0 {
		failMsg = "操作成功：" + updateStr
	} else {
		failMsg = "操作失败: " + updateStr + "！！详情内容：" + failMsg
	}
	SendBackgroundMsgTemplateHandleByUser(userId, MSG_TEMPLATE_UPDATE_PLATFORM_VERSION, PlatformId, "版本更新:"+PlatformId, failMsg)
	return err
}

// 获取参数版本工具操作日志列表
func GetVersionToolLogList(params *ParamsVersionTool) ([]*VersionToolChangeLog, int64) {
	data := make([]*VersionToolChangeLog, 0)
	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "id")
	if params.Order == "ascending" {
		sortOrder = sortOrder + " asc"
	} else {
		sortOrder = sortOrder + " desc"
	}
	var count int64

	versionChangeList, err := GetVersionToolChangeListByRobotTypeAndChangeType(params.RobotType, params.ChangeType)
	if err != nil {
		g.Log().Errorf("获取版本工具失败:%+v", err)
		return data, 0
	}
	versionChangeIdList := make([]int, 0)
	for _, versionChangeData := range versionChangeList {
		versionChangeIdList = append(versionChangeIdList, versionChangeData.Id)
	}

	//mapChildren := gmap.New()
	utils.CheckError(err)
	tempF := func(db *gorm.DB) *gorm.DB {
		if params.State > 0 {
			db = db.Where(&VersionToolChangeLog{State: params.State})
		}
		return db
	}
	err = tempF(Db.Model(&VersionToolChangeLog{})).Where("change_id in ( ? )", versionChangeIdList).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		e.UserName = GetUserName(e.UserId)
		e.ChangeName = GetVersionToolChangeOnChangeName(e.ChangeId)
	}
	return data, count
}

// 保存版本工具操作日志初始数据
func SaveVersionToolChangeLogInit(userId, changeId, state, updateType int, shName string) (*VersionToolChangeLog, error) {
	logData :=
		&VersionToolChangeLog{
			UserId:     userId,
			State:      state,
			UpdateType: updateType,
			ChangeId:   changeId,
			ShName:     shName,
			ChangeTime: gtime.Timestamp(),
		}
	err := Db.Save(logData).Error
	if err != nil {
		g.Log().Errorf("保存版本工具操作日志初始数据失败:%+v", err)
		return logData, err
	}
	err = Db.First(logData).Error
	return logData, err
}

// 保存版本工具操作日志数据
func SaveVersionToolChangeLog(userId, changeId, state, updateType int, shName, failMsg string, startTimeMilli int64) error {
	logData :=
		&VersionToolChangeLog{
			UserId:     userId,
			State:      state,
			UpdateType: updateType,
			ChangeId:   changeId,
			ShName:     shName,
			UseTime:    gconv.Float64(gtime.TimestampMilli()-startTimeMilli) / 1000,
			FailMsg:    failMsg,
			ChangeTime: gtime.Timestamp(),
		}
	err := Db.Save(logData).Error
	if err != nil {
		g.Log().Errorf("保存版本工具操作日志数据失败:%+v", err)
		return err
	}
	return err
}

// 更新版本工具操作日志数据
func UpdateVersionToolChangeLog(logData *VersionToolChangeLog, state int, failMsg string, startTimeMilli int64) error {
	logData.State = state
	logData.FailMsg = failMsg
	logData.UseTime = gconv.Float64(gtime.TimestampMilli()-startTimeMilli) / 1000
	err := Db.Save(logData).Error
	if err != nil {
		g.Log().Errorf("更新版本工具操作日志数据失败:%+v", err)
		return err
	}
	return err
}
