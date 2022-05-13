package models

import (
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
)

type TenMinuteStatistics struct {
	//Node          string `json:"node" gorm:"primary_key"`
	PlatformId        string  `json:"platformId" gorm:"primary_key"`
	ServerId          string  `json:"serverId" gorm:"primary_key"`
	Channel           string  `json:"channel" gorm:"primary_key"`
	Time              int     `json:"time" gorm:"primary_key"`
	OnlineCount       int     `json:"onlineNum"`
	RegisterCount     int     `json:"registerCount"`
	CreateRoleCount   int     `json:"createRoleCount"`
	ChargeCount       float32 `json:"chargeCount"`
	ChargePlayerCount int     `json:"chargePlayerCount"`
}

func UpdateTenMinuteStatistics(platformId string, serverId string, channelList []*Channel, timestamp int) error {
	//g.Log().Info("UpdateTenMinuteStatistics:%v, %v, %v, %v", platformId, serverId, channel, timestamp)
	serverNode, err := GetGameServerOne(platformId, serverId)
	if err != nil {
		return err
	}
	//node := serverNode.Node
	gameDb, err := GetGameDbByNode(serverNode.Node)
	if err != nil {
		return err
	}
	defer gameDb.Close()

	zeroTime := utils.GetThatZeroTimestamp(int64(timestamp))

	if timestamp == zeroTime {
		zeroTime = zeroTime - 86400
	}
	for _, e := range channelList {
		channel := e.Channel
		onlineCount := GetNowOnlineCount2(gameDb, serverId, []string{channel})
		registerCount := GetRegisterRoleCount(gameDb, serverId, channel, timestamp-600, timestamp-1)
		chargePlayerCount := GetTotalChargePlayerCount(platformId, serverId, channel, zeroTime, timestamp)
		// if onlineCount > 0 || registerCount > 0 || chargePlayerCount > 0 {
		m := &TenMinuteStatistics{
			//Node:          serverNode.Node,
			PlatformId:        platformId,
			ServerId:          serverId,
			Channel:           channel,
			Time:              timestamp,
			OnlineCount:       onlineCount,
			RegisterCount:     registerCount,
			CreateRoleCount:   GetCreateRoleCount(gameDb, serverId, channel, timestamp-600, timestamp-1),
			ChargeCount:       GetTotalChargeMoney(platformId, serverId, channel, timestamp-600, timestamp-1, 0),
			ChargePlayerCount: chargePlayerCount,
		}
		err = Db.Save(&m).Error
		if err != nil {
			return err
		}
		// }
	}
	return nil
}

func UpdateTenMinuteStatistics2(platformId string, serverId string, channelList []*Channel, timestamp int) error {
	g.Log().Infof("UpdateTenMinuteStatistics2:%v, %v, %v, %v", platformId, serverId, len(channelList), timestamp)
	serverNode, err := GetGameServerOne(platformId, serverId)
	utils.CheckError(err)
	if err != nil {
		return err
	}
	g.Log().Info("0")
	//node := serverNode.Node
	gameDb, err := GetGameDbByNode(serverNode.Node)
	utils.CheckError(err)
	if err != nil {
		return err
	}
	g.Log().Info("1")
	defer gameDb.Close()

	zeroTime := utils.GetThatZeroTimestamp(int64(timestamp))

	if timestamp == zeroTime {
		zeroTime = zeroTime - 86400
	}
	g.Log().Info("2")
	for _, e := range channelList {
		channel := e.Channel
		registerCount := GetRegisterRoleCount(gameDb, serverId, channel, timestamp-600, timestamp-1)
		chargePlayerCount := GetTotalChargePlayerCount(platformId, serverId, channel, zeroTime, timestamp)
		//g.Log().Info("go:%s, %d, %d", channel, registerCount, chargePlayerCount)
		if registerCount > 0 || chargePlayerCount > 0 {
			m := &TenMinuteStatistics{
				//Node:          serverNode.Node,
				PlatformId: platformId,
				ServerId:   serverId,
				Channel:    channel,
				Time:       timestamp,
			}
			//err = Db.Save(&m).Error
			//if timestamp <= 1543708800 + 600{
			m.RegisterCount = registerCount
			m.CreateRoleCount = GetCreateRoleCount(gameDb, serverId, channel, timestamp-600, timestamp-1)
			m.ChargeCount = GetTotalChargeMoney(platformId, serverId, channel, timestamp-600, timestamp-1, 0)
			m.ChargePlayerCount = chargePlayerCount
			err = Db.Save(&m).Error
			//} else {
			//	err = Db.Debug().Model(&m).Updates(TenMinuteStatistics{
			//		RegisterCount: registerCount,
			//		CreateRoleCount: GetCreateRoleCount(gameDb, serverId, channel, timestamp-600, timestamp-1),
			//		ChargeCount: GetTotalChargeMoney(platformId, serverId, channel, timestamp-600, timestamp-1),
			//		ChargePlayerCount:chargePlayerCount,
			//	}).Error
			//}

			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 这个方法没有使用
func RepireTenMinuteStatistics() {
	g.Log().Info("开始修复每10分钟统计")
	gameServerList, _ := GetAllGameServerDirty()
	for i := 1543766400; i <= 1543807800; i += 600 {
		for _, gameServer := range gameServerList {
			platformId := gameServer.PlatformId
			serverId := gameServer.Sid
			channelList := GetChannelListByPlatformId(platformId)
			if len(channelList) == 0 {
				g.Log().Error("渠道未配置:%v %+v", platformId, channelList)
			}

			if i < utils.GetThatZeroTimestamp(int64(gameServer.OpenTime)) {
				continue
			}
			//for _, channel := range channelList {
			err := UpdateTenMinuteStatistics2(platformId, serverId, channelList, i)
			utils.CheckError(err)
			//}
		}
	}
	g.Log().Info("修复每10分钟统计完毕")
}
func DoRepireTenMinuteStatistics(platformId string, serverId string, channel string, timestamp int) error {
	m := &TenMinuteStatistics{
		PlatformId: platformId,
		ServerId:   serverId,
		Channel:    channel,
		Time:       timestamp,
	}
	zeroTime := utils.GetThatZeroTimestamp(int64(timestamp))
	if timestamp == zeroTime {
		zeroTime = zeroTime - 86400
	}
	c := GetTotalChargePlayerCount(platformId, serverId, channel, zeroTime, timestamp)
	err := Db.Model(&m).Update("charge_player_count", c).Error
	return err
}

//func RepireTenMinuteStatistics() {
//	g.Log().Info("开始修复每10分钟统计")
//	gameServerList, _ := GetAllGameServerDirty()
//	for i := 1538755200; i <= 1538978400; i += 600 {
//		for _, gameServer := range gameServerList {
//			//err := models.UpdateDailyStatistics(serverNode.Node, todayZeroTimestamp - 86400)
//			//utils.CheckError(err)
//			platformId := gameServer.PlatformId
//			//if platformId == "af" || platformId == "djs" {
//			serverId := gameServer.Sid
//			channelList := GetChannelListByPlatformId(platformId)
//			if len(channelList) == 0 {
//				g.Log().Error("渠道未配置:%v %+v", platformId, channelList)
//			}
//
//			if i < utils.GetThatZeroTimestamp(int64(gameServer.OpenTime)) {
//				continue
//			}
//			for _, channel := range channelList {
//				err := DoRepireTenMinuteStatistics(platformId, serverId, channel.Channel, i)
//				utils.CheckError(err)
//			}
//			//}
//		}
//	}
//	g.Log().Info("修复每10分钟统计完毕")
//}
//func DoRepireTenMinuteStatistics(platformId string, serverId string, channel string, timestamp int) error {
//	m := &TenMinuteStatistics{
//		PlatformId: platformId,
//		ServerId:   serverId,
//		Channel:    channel,
//		Time:       timestamp,
//	}
//	c := GetTotalChargeMoney(platformId, serverId, channel, timestamp-600, timestamp-1)
//	err := Db.Model(&m).Update("charge_count", c).Error
//	return err
//}
