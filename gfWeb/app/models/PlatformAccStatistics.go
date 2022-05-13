package models

import (
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

// 平台订阅数据
type PlatformDingYue struct {
	PlatformId  string `json:"platformId" gorm:"PRIMARY_KEY"`
	DingYueTime int    `json:"dingYueTime" gorm:"PRIMARY_KEY"`
	DingYueNum  int    `json:"dingYueNum"`
}

func platformDingYueTBName() string {
	//return TableName("user")
	return "platform_ding_yue"
}

// 平台账号统计数据
type PlatformAccStatistics struct {
	PlatformId  string `json:"platformId" gorm:"PRIMARY_KEY"`
	AccId       string `json:"accId" gorm:"PRIMARY_KEY"`
	DingYueTime int    `json:"dingYueTime"`
}

// 获得合服数据列表
func GetPlatformDingYue(platformId string, baseQueryParam BaseQueryParam) ([]*PlatformDingYue, int64, int64) {
	var count int64
	var dingYueSum struct {
		DingYueCount int64
	}
	data := make([]*PlatformDingYue, 0)
	platformDingYue := &PlatformDingYue{PlatformId: platformId}
	Db.Model(&PlatformDingYue{}).Where(platformDingYue).Order("ding_yue_time desc").Offset(baseQueryParam.Offset).Limit(baseQueryParam.Limit).Find(&data).Count(&count)
	Db.Table(platformDingYueTBName()).Select("sum(ding_yue_num) as ding_yue_count").Where(platformDingYue).Scan(&dingYueSum)
	//Db.Model(&PlatformDingYue{PlatformId:platformId}).Select("sum(ding_yue_num) as ding_yue_count").Scan(&dingYueSum)
	return data, dingYueSum.DingYueCount, count
}

// 更新订阅统计
func UpdateDingYueStatistics(platformId string, timestamp int) {
	g.Log().Infof("更新订阅统计:%s 时间:%s 处理时间:%s", platformId, utils.TimeIntFormDefault(timestamp), gtime.Datetime())

	channelList := GetChannelListByPlatformId(platformId)
	if len(channelList) == 0 {
		g.Log().Errorf("渠道未配置:%v %+v", platformId, channelList)
		return
	}
	gameServerList, _ := GetAllGameServerDirtyByPlatformId(platformId)
	dingYueTime := utils.GetThatZeroTimestamp(int64(timestamp))
	g.Log().Debug("---更新订阅统计::%s 时间:%s 处理时间:%s", platformId, utils.TimeIntFormDefault(dingYueTime), utils.TimeIntFormDefault(dingYueTime+86400-1))
	var AccIdList []string
	var nodeList []string
	for _, gameServer := range gameServerList {
		serverId := gameServer.Sid

		if timestamp < utils.GetThatZeroTimestamp(int64(gameServer.OpenTime)) {
			continue
		}
		//for _, channel := range channelList {
		//g.Log().Debug("getGameDingYueData", platformId, " serverId:",serverId)
		AccIdList1, nodeList1 := getGameDingYueData(platformId, serverId, dingYueTime, dingYueTime+86400-1, nodeList)
		nodeList = nodeList1
		for _, AccId := range AccIdList1 {
			IsHave := false
			for _, CheckAccId := range AccIdList {
				if CheckAccId == AccId {
					IsHave = true
					break
				}
			}
			if IsHave == false {
				platformAccStatistics := getPlatformAccStatistics(platformId, AccId)
				if platformAccStatistics.DingYueTime == 0 {
					platformAccStatistics.DingYueTime = dingYueTime
					err := Db.Save(platformAccStatistics).Error
					g.Log().Infof("保存账号订阅数据:%s ;AccId:%s ; err:%+v", platformId, AccId, err)
					AccIdList = append(AccIdList, AccId)
				} else if platformAccStatistics.DingYueTime == dingYueTime {
					AccIdList = append(AccIdList, AccId)
				}
			}
		}
	}
	platformDingYue, err := getPlatformDingYue(platformId, dingYueTime)
	//if err != nil {
	//	g.Log().Error("未获得平台账号统计数据:", AccIdList)
	//	return
	//}
	AccIdLen := len(AccIdList)

	g.Log().Debug("查询平台订阅数据:", AccIdLen, ":", AccIdList)
	if AccIdLen >= platformDingYue.DingYueNum {
		platformDingYue.DingYueNum = AccIdLen
		err = Db.Save(platformDingYue).Error
		if err != nil {
			g.Log().Error("保存平台订阅数据失败:", AccIdList)
			return
		}
		g.Log().Infof("更新平台订阅数据:%s ;DingYueNum:%s", platformId, AccIdLen)
	}

}

// 获得游戏服时间范围内订阅玩家账号
func getGameDingYueData(platformId string, serverId string, startTime int, endTime int, nodeList []string) (AccIdList []string, newNodeList []string) {
	gameServer, err := GetGameServerOne(platformId, serverId)
	if err != nil {
		return AccIdList, nodeList
	}
	node := gameServer.Node
	for _, checkNode := range nodeList {
		if node == checkNode {
			return AccIdList, nodeList
		}
	}
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return AccIdList, nodeList
	}
	defer gameDb.Close()

	err = gameDb.Table("player").Select("player.acc_id").Joins("left join player_game_data on player_game_data.player_id = player.id").
		Where("player_game_data.data_id = 33 and ? <= player_game_data.int_data and player_game_data.int_data <= ?", startTime, endTime).Pluck("player.acc_id", &AccIdList).Error
	//sql := fmt.Sprintf(
	//	`SELECT player_id FROM player_game_data  WHERE data_id = 33 and %d <= int_data and int_data <= %d `, startTime, endTime)
	//err := gameDb.Raw(sql).Scan(&playerStruct).Error
	utils.CheckError(err)
	nodeList = append(nodeList, node)
	return AccIdList, nodeList
}

// 获得平台账号统计数据
func getPlatformAccStatistics(platformId string, accId string) *PlatformAccStatistics {
	platformAccStatistics := &PlatformAccStatistics{
		PlatformId: platformId,
		AccId:      accId,
	}
	Db.First(platformAccStatistics)
	return platformAccStatistics
}

// 获得平台账号统计数据
func getPlatformDingYue(platformId string, dingYueTime int) (*PlatformDingYue, error) {
	platformDingYue := &PlatformDingYue{
		PlatformId:  platformId,
		DingYueTime: dingYueTime,
	}
	err := Db.First(platformDingYue).Error
	return platformDingYue, err
}
