package models

import (
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
	"github.com/jinzhu/gorm"
)

type RemainActive struct {
	Node         string `json:"node" gorm:"primary_key"`
	PlatformId   string `json:"platformId" gorm:"primary_key"`
	ServerId     string `json:"serverId" gorm:"primary_key"`
	Channel      string `json:"channel" gorm:"primary_key"`
	Time         int    `json:"time" gorm:"primary_key"`
	RegisterRole int    `json:"registerRole" gorm:"-"`
	CreateRole   int    `json:"createRole" gorm:"-"`
	Remain2      int    `json:"remain2"`
	Remain3      int    `json:"remain3"`
	Remain4      int    `json:"remain4"`
	Remain5      int    `json:"remain5"`
	Remain6      int    `json:"remain6"`
	Remain7      int    `json:"remain7"`
	Remain8      int    `json:"remain8"`
	Remain9      int    `json:"remain9"`
	Remain10     int    `json:"remain10"`
	Remain11     int    `json:"remain11"`
	Remain12     int    `json:"remain12"`
	Remain13     int    `json:"remain13"`
	Remain14     int    `json:"remain14"`
	Remain15     int    `json:"remain15"`
	Remain16     int    `json:"remain16"`
	Remain17     int    `json:"remain17"`
	Remain18     int    `json:"remain18"`
	Remain19     int    `json:"remain19"`
	Remain20     int    `json:"remain20"`
	Remain21     int    `json:"remain21"`
	Remain22     int    `json:"remain22"`
	Remain23     int    `json:"remain23"`
	Remain24     int    `json:"remain24"`
	Remain25     int    `json:"remain25"`
	Remain26     int    `json:"remain26"`
	Remain27     int    `json:"remain27"`
	Remain28     int    `json:"remain28"`
	Remain29     int    `json:"remain29"`
	Remain30     int    `json:"remain30"`
}

type ActiveRemainQueryParam struct {
	BaseQueryParam
	PlatformId string
	//ServerId   string
	ServerId    string   `json:"serverId"`
	ChannelList []string `json:"channelList"`
	StartTime   int
	EndTime     int
}

// 获取活跃留存
func GetRemainActiveList(params *ActiveRemainQueryParam) ([]*RemainActive, int64) {
	data := make([]*RemainActive, 0)
	var count int64
	f := func(db *gorm.DB) *gorm.DB {
		if params.StartTime > 0 {
			return db.Where("time between ? and ?", params.StartTime, params.EndTime)
		}
		return db
	}
	err := f(Db.Model(&RemainActive{}).Where(&RemainActive{PlatformId: params.PlatformId, ServerId: params.ServerId})).Where(" channel in (?)", params.ChannelList).Offset(params.Offset).Limit(params.Limit).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		dailyStatistics, err := GetDailyStatisticsOne(e.PlatformId, e.ServerId, []string{e.Channel}, e.Time)
		utils.CheckError(err)
		e.CreateRole = dailyStatistics.CreateRoleCount
		e.RegisterRole = dailyStatistics.RegisterCount
	}
	return data, count
}

//更新 活跃留存
func UpdateRemainActive(platformId string, serverId string, channel string, timestamp int) error {
	g.Log().Info("更新 活跃留存:%v, %v, %v, %v", platformId, serverId, channel, timestamp)
	gameServer, err := GetGameServerOne(platformId, serverId)
	if err != nil {
		return err
	}
	node := gameServer.Node
	serverNode, err := GetServerNode(node)
	if err != nil {
		return err
	}
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return err
	}
	defer gameDb.Close()
	openDayZeroTimestamp := utils.GetThatZeroTimestamp(int64(serverNode.OpenTime))

	for i := 1; i < 30; i++ {
		//if i == 1 || i == 2 || i == 3 || i == 6 || i == 13 || i == 29  {
		//
		//} else {
		//	continue
		//}
		thatDayZeroTimestamp := timestamp - i*86400
		if openDayZeroTimestamp > thatDayZeroTimestamp {
			continue
		}
		createRolePlayerIdList := GetThatDayCreateRolePlayerIdList(gameDb, serverId, channel, thatDayZeroTimestamp)
		createRoleNum := len(createRolePlayerIdList)
		rate := 0
		if createRoleNum > 0 {
			lastDayLoginNum := 0
			LastDayLoginPlayerIdList := make([]int, 0)
			for _, playerId := range createRolePlayerIdList {
				if IsThatDayPlayerLogin(gameDb, timestamp-86400, playerId) {
					lastDayLoginNum += 1
					LastDayLoginPlayerIdList = append(LastDayLoginPlayerIdList, playerId)
				}
			}
			if lastDayLoginNum > 0 {
				loginNum := 0
				for _, playerId := range LastDayLoginPlayerIdList {
					if IsThatDayPlayerLogin(gameDb, timestamp, playerId) {
						loginNum += 1
					}
				}
				rate = int(float32(loginNum) / float32(lastDayLoginNum) * 10000)
			}
			//loginNum := 0
			//for _, playerId := range createRolePlayerIdList {
			//	if IsThatDayPlayerLogin(gameDb, timestamp, playerId) {
			//		loginNum += 1
			//	}
			//}
			//if loginNum > 0 {
			//	LastDayLoginNum := 0
			//	for _, playerId := range createRolePlayerIdList {
			//		if IsThatDayPlayerLogin(gameDb, timestamp-86400, playerId) {
			//			LastDayLoginNum += 1
			//		}
			//	}
			//	if LastDayLoginNum == 0 {
			//		g.Log().Error("更新活跃留存遇到 分母为0:%v, %v, %v, %v", node, timestamp, loginNum, LastDayLoginNum)
			//		LastDayLoginNum = 1
			//	}
			//	rate = int(float32(loginNum) / float32(LastDayLoginNum) * 10000)
			//}

		}
		m := &RemainActive{
			Node:       node,
			PlatformId: platformId,
			ServerId:   serverId,
			Channel:    channel,
			Time:       thatDayZeroTimestamp,
			Remain2:    -1,
			Remain3:    -1,
			Remain4:    -1,
			Remain5:    -1,
			Remain6:    -1,
			Remain7:    -1,
			Remain8:    -1,
			Remain9:    -1,
			Remain10:   -1,
			Remain11:   -1,
			Remain12:   -1,
			Remain13:   -1,
			Remain14:   -1,
			Remain15:   -1,
			Remain16:   -1,
			Remain17:   -1,
			Remain18:   -1,
			Remain19:   -1,
			Remain20:   -1,
			Remain21:   -1,
			Remain22:   -1,
			Remain23:   -1,
			Remain24:   -1,
			Remain25:   -1,
			Remain26:   -1,
			Remain27:   -1,
			Remain28:   -1,
			Remain29:   -1,
			Remain30:   -1,
		}
		err = Db.FirstOrCreate(&m).Error
		if err != nil {
			return err
		}
		switch i {
		case 1:
			err = Db.Model(&m).Update("Remain2", rate).Error
		case 2:
			err = Db.Model(&m).Update("Remain3", rate).Error
		case 3:
			err = Db.Model(&m).Update("Remain4", rate).Error
		case 4:
			err = Db.Model(&m).Update("Remain5", rate).Error
		case 5:
			err = Db.Model(&m).Update("Remain6", rate).Error
		case 6:
			err = Db.Model(&m).Update("Remain7", rate).Error

		case 7:
			err = Db.Model(&m).Update("Remain8", rate).Error
		case 8:
			err = Db.Model(&m).Update("Remain9", rate).Error
		case 9:
			err = Db.Model(&m).Update("Remain10", rate).Error
		case 10:
			err = Db.Model(&m).Update("Remain11", rate).Error
		case 11:
			err = Db.Model(&m).Update("Remain12", rate).Error
		case 12:
			err = Db.Model(&m).Update("Remain13", rate).Error

		case 13:
			err = Db.Model(&m).Update("Remain14", rate).Error
		case 14:
			err = Db.Model(&m).Update("Remain15", rate).Error
		case 15:
			err = Db.Model(&m).Update("Remain16", rate).Error
		case 16:
			err = Db.Model(&m).Update("Remain17", rate).Error
		case 17:
			err = Db.Model(&m).Update("Remain18", rate).Error
		case 18:
			err = Db.Model(&m).Update("Remain19", rate).Error

		case 19:
			err = Db.Model(&m).Update("Remain20", rate).Error
		case 20:
			err = Db.Model(&m).Update("Remain21", rate).Error
		case 21:
			err = Db.Model(&m).Update("Remain22", rate).Error
		case 22:
			err = Db.Model(&m).Update("Remain23", rate).Error
		case 23:
			err = Db.Model(&m).Update("Remain24", rate).Error
		case 24:
			err = Db.Model(&m).Update("Remain25", rate).Error

		case 25:
			err = Db.Model(&m).Update("Remain26", rate).Error
		case 26:
			err = Db.Model(&m).Update("Remain27", rate).Error
		case 27:
			err = Db.Model(&m).Update("Remain28", rate).Error
		case 28:
			err = Db.Model(&m).Update("Remain29", rate).Error
		case 29:
			err = Db.Model(&m).Update("Remain30", rate).Error
		}

	}

	return err
}
