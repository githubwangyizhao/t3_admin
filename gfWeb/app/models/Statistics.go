package models

import (
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

// func RepairAllGameNodeDailyStatistics() {
// 	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
// 	DoUpdateAllGameNodeDailyStatistics(todayZeroTimestamp - 86400)
// 	DoUpdateAllGameNodeDailyStatistics(todayZeroTimestamp - 86400*2)
// 	DoUpdateAllGameNodeDailyStatistics(todayZeroTimestamp - 86400*3)
// 	DoUpdateAllGameNodeDailyStatistics(todayZeroTimestamp - 86400*4)
// 	DoUpdateAllGameNodeDailyStatistics(todayZeroTimestamp - 86400*5)
// }

//更新所有游戏节点  DailyStatistics
func UpdateAllGameNodeDailyStatistics() {
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()

	DoUpdateAllGameNodeDailyStatistics(todayZeroTimestamp - 86400)
	// DoUpdateAllGameNodeDailyStatistics(1630080000 - 86400)
}

// 补未执行的dailyStatistic
func PlusUpdateAllGameNodeDailyStatistics(startTime, endTime int) {
	for i := startTime; i <= endTime; i += 86400 {
		DoUpdateAllGameNodeDailyStatistics(i - 86400)
	}
}

//更新所有游戏节点  DailyStatistics
func UpdateAllGameNodeDailyLTV() {
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	DoUpdateAllGameNodeDailyLTV(todayZeroTimestamp - 86400)
}

//更新所有游戏节点  ReaminCharge
func UpdateAllGameNodeChargeRemain() {
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	DoUpdateAllGameNodeReaminCharge(todayZeroTimestamp - 86400)
}

//func Repire() {
//	//1532016000
//	//1536163200 退1536422400
//
//	for i := 1532016000; i <= 1538236800; i += 86400 {
//		DoUpdateAllGameNodeDailyLTV(i)
//	}
//}

func Repire() {
	//1532016000
	//1536163200 退1536422400

	for i := 1537200000; i <= 1538236800; i += 86400 {
		DoUpdateAllGameNodeReaminCharge(i)
	}
}

func DoUpdateAllGameNodeTenMinuteStatistics(timestamp int) {
	g.Log().Infof("更新每10分钟统计:%v", gtime.Timestamp())
	gameServerList, _ := GetAllGameServerDirty()

	for _, gameServer := range gameServerList {
		platformId := gameServer.PlatformId
		serverId := gameServer.Sid
		channelList := GetChannelListByPlatformId(platformId)
		if len(channelList) == 0 {
			g.Log().Errorf("渠道未配置:%v %+v", platformId, channelList)
		}

		if timestamp < utils.GetThatZeroTimestamp(int64(gameServer.OpenTime)) {
			continue
		}
		//for _, channel := range channelList {
		err := UpdateTenMinuteStatistics(platformId, serverId, channelList, timestamp)
		utils.CheckError(err)
		//}
	}
	g.Log().Info("更新每10分钟完毕.")
}

func DoUpdateAllGameNodeDailyStatistics(timestamp int) {
	//todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	g.Log().Infof("更新每日统计:%v", timestamp)
	gameServerList, _ := GetAllGameServerDirty()

	for _, gameServer := range gameServerList {
		platformId := gameServer.PlatformId
		serverId := gameServer.Sid
		channelList := GetChannelListByPlatformId(platformId)
		if len(channelList) == 0 {
			g.Log().Errorf("渠道未配置:%v %+v", platformId, channelList)
		}

		if timestamp < utils.GetThatZeroTimestamp(int64(gameServer.OpenTime)) {
			continue
		}
		err := UpdateDailyStatistics(platformId, serverId, channelList, timestamp)
		utils.CheckError(err)
	}
	g.Log().Info("更新每日统计完毕.")
}

func DoUpdateAllGameNodeDailyLTV(timestamp int) {
	//todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	g.Log().Info("更新每日LTV:%v", timestamp)
	gameServerList, _ := GetAllGameServerDirty()
	for _, gameServer := range gameServerList {
		//err := UpdateDailyStatistics(serverNode.Node, todayZeroTimestamp - 86400)
		//utils.CheckError(err)
		platformId := gameServer.PlatformId
		serverId := gameServer.Sid
		channelList := GetChannelListByPlatformId(platformId)
		if len(channelList) == 0 {
			g.Log().Error("渠道未配置:%v %+v", platformId, channelList)
		}

		if timestamp < utils.GetThatZeroTimestamp(int64(gameServer.OpenTime)) {
			continue
		}
		//for _, channel := range channelList {
		err := UpdateDailyLTV(platformId, serverId, channelList, timestamp)
		utils.CheckError(err)
		//}
	}
	g.Log().Info("更新每日LTV完毕.")
}

func DoUpdateAllGameNodeReaminCharge(timestamp int) {
	//todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	g.Log().Info("更新付费留存:%v", timestamp)
	gameServerList, _ := GetAllGameServerDirty()
	for _, gameServer := range gameServerList {
		//err := UpdateDailyStatistics(serverNode.Node, todayZeroTimestamp - 86400)
		//utils.CheckError(err)
		platformId := gameServer.PlatformId
		//if platformId == "af" || platformId == "djs" {
		serverId := gameServer.Sid
		channelList := GetChannelListByPlatformId(platformId)
		if len(channelList) == 0 {
			g.Log().Error("渠道未配置:%v %+v", platformId, channelList)
		}

		if timestamp < utils.GetThatZeroTimestamp(int64(gameServer.OpenTime)) {
			continue
		}
		//for _, channel := range channelList {
		err := UpdateRemainCharge(platformId, serverId, channelList, timestamp)
		utils.CheckError(err)
		//}
		//}
	}
	g.Log().Info("更新付费留存完毕.")
}

////更新所有游戏节点  DailyStatistics
//func UpdateAllGameNodeDailyStatistics() {
//	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
//	g.Log().Info("UpdateAllGameNodeDailyStatistics:%v", todayZeroTimestamp)
//	gameServerNodeList := models.GetAllGameServerNode()
//	for _, serverNode := range gameServerNodeList {
//		//err := models.UpdateDailyStatistics(serverNode.Node, todayZeroTimestamp - 86400)
//		//utils.CheckError(err)
//		err := models.UpdateDailyChargeStatistics(serverNode.Node, todayZeroTimestamp-86400)
//		utils.CheckError(err)
//		err = models.UpdateDailyOnlineStatistics(serverNode.Node, todayZeroTimestamp-86400)
//		utils.CheckError(err)
//		err = models.UpdateDailyRegisterStatistics(serverNode.Node, todayZeroTimestamp-86400)
//		utils.CheckError(err)
//		err = models.UpdateDailyActiveStatistics(serverNode.Node, todayZeroTimestamp-86400)
//		utils.CheckError(err)
//	}
//}

//更新所有游戏节点  RemainTotal
func UpdateAllGameNodeRemainTotal() {
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400)
}

func DoUpdateAllGameNodeRemainTotal(timestamp int) {
	now := utils.GetTimestamp()
	//todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	g.Log().Info("更新所有总体留存:%v", timestamp)
	gameServerList, _ := GetAllGameServerDirty()

	for _, gameServer := range gameServerList {
		serverNode, err := GetServerNode(gameServer.Node)
		utils.CheckError(err)
		platformId := gameServer.PlatformId
		serverId := gameServer.Sid
		channelList := GetChannelListByPlatformId(platformId)
		if len(channelList) == 0 {
			g.Log().Error("渠道未配置:%v %+v", platformId, channelList)
		}
		if now >= serverNode.OpenTime {
			//for _, channel := range channelList {
			err := UpdateRemainTotal(platformId, serverId, channelList, timestamp)
			utils.CheckError(err)
			//}
		}
	}
	g.Log().Info("更新所有总体留存完毕。")
}

//func UpdateAllGameNodeLTV() {
//	now := utils.GetTimestamp()
//	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
//	g.Log().Info("更新所有LTV:%v", todayZeroTimestamp)
//	gameServerList, _ := models.GetAllGameServer()
//
//	for _, gameServer := range gameServerList {
//		serverNode, err := models.GetServerNode(gameServer.Node)
//		utils.CheckError(err)
//		platformId := gameServer.PlatformId
//		serverId := gameServer.Sid
//		channelList := models.GetChannelListByPlatformId(platformId)
//		if len(channelList) == 0 {
//			g.Log().Error("渠道未配置:%v %+v", platformId, channelList)
//		}
//		if now >= serverNode.OpenTime {
//			for _, channel := range channelList {
//				err := models.UpdateRemainTotal(platformId, serverId, channel.Channel, todayZeroTimestamp-86400)
//				utils.CheckError(err)
//			}
//		}
//	}
//}

//func UpdateAllGameNodeRemainTotal() {
//	now := utils.GetTimestamp()
//	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
//	g.Log().Info("更新所有总体留存:%v", todayZeroTimestamp)
//	gameServerNodeList := models.GetAllGameServerNode()
//	for _, serverNode := range gameServerNodeList {
//		if now >= serverNode.OpenTime {
//			err := models.UpdateRemainTotal(serverNode.Node, todayZeroTimestamp-86400)
//			utils.CheckError(err)
//		}
//
//	}
//}

func RepairAllGameNodeRemainActive() {
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	for i := 0; i < 10; i++ {
		DoUpdateAllGameNodeRemainActive(todayZeroTimestamp - 86400*i)
	}
}

func UpdateAllGameNodeRemainActive() {
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	g.Log().Info("更新所有活跃留存:%v", todayZeroTimestamp)
	DoUpdateAllGameNodeRemainActive(todayZeroTimestamp)
}

//更新所有游戏节点  RemainActive
func DoUpdateAllGameNodeRemainActive(todayZeroTimestamp int) {
	now := utils.GetTimestamp()
	gameServerList, _ := GetAllGameServerDirty()

	for _, gameServer := range gameServerList {
		serverNode, err := GetServerNode(gameServer.Node)
		utils.CheckError(err)
		platformId := gameServer.PlatformId
		serverId := gameServer.Sid
		channelList := GetChannelListByPlatformId(platformId)
		if len(channelList) == 0 {
			g.Log().Error("渠道未配置:%v %+v", platformId, channelList)
		}
		if now >= serverNode.OpenTime {
			for _, channel := range channelList {
				err := UpdateRemainActive(platformId, serverId, channel.Channel, todayZeroTimestamp-86400)
				utils.CheckError(err)
			}
		}
	}
}

//func TmpUpdateAllGameNodeRemainTotal(time int) {
//	now := utils.GetTimestamp()
//	g.Log().Info("更新所有总体留存:%v", time)
//	gameServerNodeList := models.GetAllGameServerNode()
//	for _, serverNode := range gameServerNodeList {
//		if now >= serverNode.OpenTime {
//			err := models.UpdateRemainTotal(serverNode.Node, time)
//			utils.CheckError(err)
//		}
//
//	}
//}
//
//func TmpUpdateAllGameNodeRemainActive(time int) {
//	now := utils.GetTimestamp()
//	g.Log().Info("更新所有活跃留存:%v", time)
//	gameServerNodeList := models.GetAllGameServerNode()
//	for _, serverNode := range gameServerNodeList {
//		if now >= serverNode.OpenTime {
//			err := models.UpdateRemainActive(serverNode.Node, time-86400)
//			utils.CheckError(err)
//		}
//
//	}
//}
