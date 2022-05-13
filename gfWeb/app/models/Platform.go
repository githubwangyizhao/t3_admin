package models

import (
	"encoding/json"
	"fmt"
	"gfWeb/library/utils"
	"sort"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/util/gconv"
)

type Platform struct {
	Id                        string                       `json:"id"`
	Name                      string                       `json:"name"`
	IsAutoOpenServer          int                          `json:"isAutoOpenServer"`
	InventoryDatabaseId       int                          `json:"inventoryDatabaseId"`
	ZoneInventoryServerId     int                          `json:"zoneInventoryServerId"`
	CreateRoleLimit           int                          `json:"createRoleLimit"`
	Version                   string                       `json:"version"`
	PlatformInventorySeverRel []*PlatformInventorySeverRel `json:"-"`
	ChannelList               []*Channel                   `json:"channelList"`
	InventorySeverIds         []int                        `json:"inventorySeverIds" gorm:"-"`
	OpenServerTakeTime        int                          `json:"openServerTakeTime"`   // 定固定时间开服
	IntervalInitTime          int                          `json:"intervalInitTime"`     // 间隔开服初始时间
	IntervalDay               int                          `json:"intervalDay"`          // 间隔几天开服
	OpenServerTimeScope       string                       `json:"openServerTimeScope"`  // 开服时间段24小时,分割
	ServerAliasStr            string                       `json:"serverAliasStr"`       // 区服别名xx%dxx
	OpenServerChangeTime      int64                        `json:"openServerChangeTime"` // 开服管理操作时间',
	CronUpdateTime            int64                        `json:"cronUpdateTime"`       // 订时更新时间',
	CronUpdateType            int                          `json:"cronUpdateType"`       // 更新类型（0:热更;1:冷更,10:开启，11:关闭）',
	EnterStateType            int                          `json:"enterStateType"`       // 入口状态类型', 0:不操作 1:关闭 2:关闭更新完开启
	EnterState                int                          `json:"enterState"`           // 入口状态',
	UpdateState               int                          `json:"updateState"`          // 更新状态', 1:可更新 2:更新中 3:已更新 4:关闭 5:失败
	UpdateUserId              int                          `json:"updateUserId"gorm:"-"` // 更新操作者
	Time                      int64                        `json:"time"`
	TrackerToken              string                       `json:"trackerToken"gorm:"tracker_token"` // adjust的tracker_token
}

func (a *Platform) TableName() string {
	return PlatformDatabaseTBName()
}

func PlatformDatabaseTBName() string {
	return "platform"
}

// 获得订时器名字
func (p Platform) getCronName() string {
	return CRON_NAME_PLATFORM_UPDATE_VERSION + p.Id
}

// 获得订时开服订时器名字
func (p Platform) getCronNameOpenServer() string {
	return CRON_NAME_PLATFORM_OPEN_SERVER + p.Id
}

type PlatformParam struct {
	BaseQueryParam
}

//type Platform struct {
//	Id   int    `json:"id"`
//	Name string `json:"name"`
//	//User string `json:"user"`
//	//Port int    `json:"port"`
//	//Host string `json:"host"`
//	//AddTime           int    `json:"addTime"`
//	Time        int    `json:"time"`
//}

//获取平台列表
func GetPlatformListByPlatformIdList(platformIdList []string) []*Platform {
	data := make([]*Platform, 0)
	if len(platformIdList) == 0 {
		err := Db.Model(&Platform{}).Find(&data).Error
		utils.CheckError(err)
	} else {
		err := Db.Model(&Platform{}).Where("id in (?)", platformIdList).Find(&data).Error
		utils.CheckError(err)
	}

	for _, v := range data {
		v.InventorySeverIds = make([]int, 0)

		err := Db.Model(&v).Related(&v.PlatformInventorySeverRel).Error
		utils.CheckError(err)

		for _, e := range v.PlatformInventorySeverRel {
			v.InventorySeverIds = append(v.InventorySeverIds, e.InventoryServerId)
		}
		sort.Ints(v.InventorySeverIds)
		v.ChannelList = GetChannelListByPlatformId(v.Id)
	}
	return data
}

//获取平台列表
func GetPlatformList() []*Platform {
	return GetPlatformListByPlatformIdList([]string{})
}

//获取平台精简的数据列表(不能有数据写入)
func GetPlatformSimpleList() []*Platform {
	return GetPlatformSimpleListByPlatformIdList([]string{})
}
func GetPlatformSimpleListByPlatformIdList(platformIdList []string) []*Platform {
	data := make([]*Platform, 0)
	if len(platformIdList) == 0 {
		err := Db.Model(&Platform{}).Find(&data).Error
		utils.CheckError(err)
	} else {
		err := Db.Model(&Platform{}).Where("id in (?)", platformIdList).Find(&data).Error
		utils.CheckError(err)
	}
	//err := Db.Model(&Platform{}).Find(&data).Error
	//utils.CheckError(err)
	return data
}

//获取平台精简的数据(不能有数据写入)
func GetPlatformSimpleOne(platformId string) *Platform {
	p := &Platform{
		Id: platformId,
	}
	err := Db.First(p).Error
	utils.CheckError(err)
	return p
}

//获取平台列表
func GetPlatformListByUserId(userId int) []*Platform {
	var list []*Platform
	var channelList []*Channel
	user, err := GetUserOne(userId)
	utils.CheckError(err)
	if user.IsSuperUser() {
		list = GetPlatformList()
	} else {
		sql := fmt.Sprintf(`SELECT DISTINCT T2.*
		FROM %s AS T0
		INNER JOIN %s AS T1 ON T0.role_id = T1.role_id
		INNER JOIN %s AS T2 ON T2.id = T0.channel_id
		WHERE T1.user_id = ?`, RoleChannelRelTBName(), RoleUserRelTBName(), ChannelDatabaseTBName())
		rows, err := Db.Raw(sql, userId).Rows()
		defer rows.Close()
		utils.CheckError(err)
		for rows.Next() {
			var channel Channel
			Db.ScanRows(rows, &channel)
			//g.Log().Debug("channel:%+v", channel)
			channelList = append(channelList, &channel)
		}
		flag := make(map[string]bool)
		for _, v := range channelList {
			_, ok := flag[v.PlatformId]
			if ok {

			} else {
				flag[v.PlatformId] = true
				platform, err := GetPlatformOne(v.PlatformId)
				utils.CheckError(err)
				list = append(list, platform)
			}
		}
		for _, v := range list {
			for _, channel := range channelList {
				if channel.PlatformId == v.Id {
					v.ChannelList = append(v.ChannelList, channel)
				}
			}
		}
	}
	return list
}

//获取用户平台列表
func GetPlatformIdListByUserId(userId int) []string {
	platformIdList := make([]string, 0)
	platformList := GetPlatformListByUserId(userId)
	for _, e := range platformList {
		platformIdList = append(platformIdList, e.Id)
	}
	return platformIdList
}

////获取平台列表
//func GetPlatformListByUserId(userId int) []*Platform {
//	var list []*Platform
//	user, err := GetUserOne(userId)
//	utils.CheckError(err)
//	if user.IsSuper == 1 {
//		list = GetPlatformList()
//	} else {
//		sql := fmt.Sprintf(`SELECT DISTINCT T2.*
//		FROM %s AS T0
//		INNER JOIN %s AS T1 ON T0.role_id = T1.role_id
//		INNER JOIN %s AS T2 ON T2.id = T0.platform_id
//		WHERE T1.user_id = ?`, RoleChannelRelTBName(), RoleUserRelTBName(), PlatformDatabaseTBName())
//		rows, err := Db.Raw(sql, userId).Rows()
//		defer rows.Close()
//		utils.CheckError(err)
//		for rows.Next() {
//			var platform Platform
//			Db.ScanRows(rows, &platform)
//			list = append(list, &platform)
//		}
//	}
//	return list
//}

//获取单个平台
func GetPlatformOne(id string) (*Platform, error) {
	r := &Platform{
		Id: id,
	}
	err := Db.First(&r).Error
	r.InventorySeverIds = make([]int, 0)

	err = Db.Model(&r).Related(&r.PlatformInventorySeverRel).Error
	utils.CheckError(err)

	for _, e := range r.PlatformInventorySeverRel {
		r.InventorySeverIds = append(r.InventorySeverIds, e.InventoryServerId)
	}
	sort.Ints(r.InventorySeverIds)
	return r, err
}

// 删除平台列表
func DeletePlatform(ids []string) error {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := Db.Where(ids).Delete(&Platform{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	//删除平台渠道关系
	//if _, err := DeleteChannelRelByPlatformIdList(ids); err != nil {
	//	tx.Rollback()
	//	g.Log().Info("删除平台渠道关系")
	//	return err
	//}
	//删除角色平台关系
	//if _, err := DeleteRoleChannelRelByPlatformIdList(ids); err != nil {
	//	tx.Rollback()
	//	g.Log().Info("删除角色平台关系")
	//	return err
	//}
	//删除服务器平台关系
	if _, err := DeletePlatformInventorySeverRelByPlatformIdList(ids); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// 重置定时开服时间
func UpdatePlatformOpenServerTime(platformId string, openServerTime int) error {
	platform, err := GetPlatformOne(platformId)
	utils.CheckError(err, "获取平台失败:"+platformId)
	if err != nil {
		return err
	}
	platform.OpenServerTakeTime = openServerTime
	err = Db.Save(&platform).Error
	utils.CheckError(err, "重置开服时间失败")
	if err != nil {
		return err
	}
	err = InitCronPlatformOpenServerTime(platform)
	return err
}

// 保存开服服务管理
func SaveOpenServerManage(userId int, oldPlatform, newPlatform *Platform) error {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := SaveOpenServerManageLog(userId, oldPlatform, newPlatform); err != nil {
		utils.CheckError(err, "开服管理日志保存失败")
		return err
	}
	if err := Db.Save(newPlatform).Error; err != nil {
		utils.CheckError(err, "开服管理设置失败")
		return err
	}
	if err := InitCronPlatformOpenServerTime(newPlatform); err != nil {
		utils.CheckError(err, "初始更新订时器启动失败")
		return err
	}
	return tx.Commit().Error
}

// 初始更新订时器启动
func InitCronPlatformOpenServerTime(platform *Platform) error {
	gcron.Remove(platform.getCronNameOpenServer())
	if platform.OpenServerTakeTime <= 0 {
		return nil
	}
	cronFun := func() {
		OpenServerType(0, platform.Id, platform.OpenServerTakeTime, 0)
	}
	cronTimeStr := utils.TimestampToCronStr(gconv.Int64(platform.OpenServerTakeTime))
	_, err := gcron.AddOnce(cronTimeStr, cronFun, platform.getCronNameOpenServer())
	if err != nil {
		g.Log().Errorf("平台更新订时器启动失败:%+v  err:%+v", platform.Id, err)
		return err
	}
	return nil
}

//
//func GetPlatformList()[] *Platform {
//	filename := "views/static/json/Platform.json"
//	bytes, err := ioutil.ReadFile(filename)
//	list := make([] *Platform, 0)
//	if err != nil {
//		fmt.Println("ReadFile: ", err.Error())
//		g.Log().Error("GetPlatformList:%v, %v", filename, err)
//		return nil, err
//	}
//
//	if err := json.Unmarshal(bytes, &list); err != nil {
//		g.Log().Error("Unmarshal json:%v, %v", filename, err)
//		return nil, err
//	}
//	return list, nil
//}

func GetPlatFormPayTimes(AppId string, platformClientInfo PlatformClientInfo) map[string]interface{} {
	err := Db.Where("app_id = ?", AppId).First(&platformClientInfo).Error
	utils.CheckError(err)
	data := map[string]interface{}{}
	if platformClientInfo.Id > 0 {
		data["appId"] = platformClientInfo.AppId
		data["payTimes"] = platformClientInfo.PayTimes
	}
	return data
}

func AsyncNoticeCenter(platform string, channel string, trackerToken string) error {
	pool := utils.GetAsyncPool()

	err := pool.Add(func() {
		var request struct {
			PlatformId   string `json:"platform_id"`
			TrackerToken string `json:"tracker_token"`
			Channel      string `json:"channel"`
		}
		request.PlatformId = platform
		request.TrackerToken = trackerToken
		request.Channel = channel

		data, err := json.Marshal(request)
		utils.CheckError(err)
		if err != nil {
			g.Log().Errorf("request failure: %+v", err)
		} else {
			url := utils.GetCenterURL() + "/set_platform_tracker_token"
			resp, _ := utils.HttpRequest(url, string(data))
			g.Log().Infof("call url: %s response: %+v", url, resp)
		}
	})
	if err != nil {
		g.Log().Errorf("error: %+v", err)
	}
	return err
}
