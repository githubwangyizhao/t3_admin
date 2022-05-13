package models

import (
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
	//"fmt"
)

type RemainTotal struct {
	//Node         string `json:"node" gorm:"primary_key"`
	PlatformId   string `json:"platformId" gorm:"primary_key"`
	ServerId     string `json:"serverId" gorm:"primary_key"`
	Channel      string `json:"channel" gorm:"primary_key"`
	Time         int    `json:"time" gorm:"primary_key"`
	RegisterRole int    `json:"registerRole"`
	CreateRole   int    `json:"createRole"`
	Remain2      int    `json:"remain2"`
	Remain3      int    `json:"remain3"`
	Remain4      int    `json:"remain4"`
	//Remain5      int    `json:"remain5"`
	//Remain6      int    `json:"remain6"`
	Remain7 int `json:"remain7"`
	//Remain8      int    `json:"remain8"`
	//Remain9      int    `json:"remain9"`
	//Remain10     int    `json:"remain10"`
	//Remain11     int    `json:"remain11"`
	//Remain12     int    `json:"remain12"`
	//Remain13     int    `json:"remain13"`
	Remain14 int `json:"remain14"`
	//Remain15     int    `json:"remain15"`
	//Remain16     int    `json:"remain16"`
	//Remain17     int    `json:"remain17"`
	//Remain18     int    `json:"remain18"`
	//Remain19     int    `json:"remain19"`
	//Remain20     int    `json:"remain20"`
	//Remain21     int    `json:"remain21"`
	//Remain22     int    `json:"remain22"`
	//Remain23     int    `json:"remain23"`
	//Remain24     int    `json:"remain24"`
	//Remain25     int    `json:"remain25"`
	//Remain26     int    `json:"remain26"`
	//Remain27     int    `json:"remain27"`
	//Remain28     int    `json:"remain28"`
	//Remain29     int    `json:"remain29"`
	Remain30 int `json:"remain30"`
}

type TotalRemainQueryParam struct {
	BaseQueryParam
	PlatformId string
	//ServerId   string
	ServerId    string   `json:"serverId"`
	ChannelList []string `json:"channelList"`
	Channel     string
	StartTime   int
	EndTime     int
}

// 获取总体留存
func GetRemainTotalList(params *TotalRemainQueryParam) ([]*RemainTotal, int64) {
	if params.EndTime < params.StartTime || params.StartTime == 0 {
		g.Log().Error("开始结束时间错误")
		return nil, 0
	}
	data := make([]*RemainTotal, 0, (params.EndTime-params.StartTime)/86400)
	var count int64
	//f := func(db *gorm.DB) *gorm.DB {
	//	if params.StartTime > 0 {
	//		return db.Where("time between ? and ?", params.StartTime, params.EndTime)
	//	}
	//	return db
	//}
	//err := Db.Model(&RemainTotal{}).Where(&RemainTotal{PlatformId: params.PlatformId, ServerId: params.ServerId, Channel: params.Channel}).Where(" channel in (?)", params.ChannelList).Count(&count).Offset(params.Offset).Limit(params.Limit).Find(&data).Error
	//utils.CheckError(err)

	channelLen := len(params.ChannelList)
	for i := params.StartTime; i <= params.EndTime; i = i + 86400 {
		tmpData := make([]*RemainTotal, 0, channelLen)
		err := Db.Model(&RemainTotal{}).Where(&RemainTotal{PlatformId: params.PlatformId, ServerId: params.ServerId, Time: i}).Where(" channel in (?)", params.ChannelList).Find(&tmpData).Error
		//err := Db.Model(&RemainCharge{}).Where(&RemainCharge{PlatformId: params.PlatformId, ServerId: params.ServerId, Time: i}).Where("channel in (?)", params.ChannelList).Find(&tmpData).Error
		utils.CheckError(err)
		if len(tmpData) > 0 {
			tmpE := &RemainTotal{
				PlatformId: params.PlatformId,
				ServerId:   params.ServerId,
				Time:       i,
			}
			for _, e := range tmpData {
				tmpDailyData := make([]*DailyStatistics, 0, channelLen)
				//fmt.Println("go", params.PlatformId, params.ServerId, i, e.Channel)
				//err := Db.Model(&DailyStatistics{}).Where(&DailyStatistics{PlatformId: params.PlatformId, ServerId: e.ServerId, Time: i, Channel: e.Channel}).Find(&tmpDailyData).Error
				err := Db.Debug().Model(&DailyStatistics{}).Where(&DailyStatistics{PlatformId: params.PlatformId, ServerId: e.ServerId, Time: i, Channel: e.Channel}).Group("time").Find(&tmpDailyData).Error
				utils.CheckError(err)
				if err != nil {
					tmpE.CreateRole += e.CreateRole
					tmpE.RegisterRole += e.RegisterRole
				} else {
					//fmt.Println("len", len(tmpDailyData))
					for _, e1 := range tmpDailyData {
						tmpE.CreateRole += e1.CreateRoleCount
						tmpE.RegisterRole += e1.RegisterCount
						g.Log().Debugf("create_role_count: %d create_count: %d, time: %d tmpE: %d %d",
							e1.CreateRoleCount, e1.RegisterCount, e1.Time, tmpE.CreateRole, tmpE.RegisterRole)
					}
				}

				tmpE.Remain2 += e.Remain2
				tmpE.Remain3 += e.Remain3
				tmpE.Remain4 += e.Remain4
				tmpE.Remain7 += e.Remain7
				tmpE.Remain14 += e.Remain14
				tmpE.Remain30 += e.Remain30
			}
			data = append(data, tmpE)
		}
	}
	for _, e := range data {
		if e.CreateRole > 0 {
			e.Remain2 = int(float32(e.Remain2) / float32(e.CreateRole) * 10000)
		} else {
			e.Remain2 = 0
		}
		if e.CreateRole > 0 {
			e.Remain3 = int(float32(e.Remain3) / float32(e.CreateRole) * 10000)
		} else {
			e.Remain3 = 0
		}
		if e.CreateRole > 0 {
			e.Remain4 = int(float32(e.Remain4) / float32(e.CreateRole) * 10000)
		} else {
			e.Remain4 = 0
		}
		if e.CreateRole > 0 {
			e.Remain7 = int(float32(e.Remain7) / float32(e.CreateRole) * 10000)
		} else {
			e.Remain7 = 0
		}
		if e.CreateRole > 0 {
			e.Remain14 = int(float32(e.Remain14) / float32(e.CreateRole) * 10000)
		} else {
			e.Remain14 = 0
		}
		if e.CreateRole > 0 {
			e.Remain30 = int(float32(e.Remain30) / float32(e.CreateRole) * 10000)
		} else {
			e.Remain30 = 0
		}
	}
	//for _, e := range data {
	//	dailyStatistics, err := GetDailyStatisticsOne(e.PlatformId, e.ServerId, e.Channel, e.Time)
	//	//g.Log().Info("dailyStatistics:%+v", dailyStatistics)
	//	utils.CheckError(err)
	//	//e.CreateRole = dailyStatistics.CreateRoleCount
	//	e.RegisterRole = dailyStatistics.RegisterCount
	//}
	return data, count
}

//更新 总体留存
func UpdateRemainTotal(platformId string, serverId string, channelList []*Channel, timestamp int) error {
	g.Log().Info("更新总体留存:%v, %v, %v, %v", platformId, serverId, len(channelList), timestamp)
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

	for _, e := range channelList {
		channel := e.Channel
		//g.Log().Info("%s", channel)
		for i := 1; i < 30; i++ {
			if i == 1 || i == 2 || i == 3 || i == 6 || i == 13 || i == 29 {

			} else {
				continue
			}
			thatDayZeroTimestamp := timestamp - i*86400
			if openDayZeroTimestamp > thatDayZeroTimestamp {
				continue
			}
			createRolePlayerIdList := GetThatDayCreateRolePlayerIdList(gameDb, serverId, channel, thatDayZeroTimestamp)
			registerCount := GetRegisterRoleCount(gameDb, serverId, channel, thatDayZeroTimestamp, thatDayZeroTimestamp+86400)
			createRoleNum := len(createRolePlayerIdList)
			rate := 0
			if createRoleNum > 0 {
				loginNum := 0
				for _, playerId := range createRolePlayerIdList {
					if IsThatDayPlayerLogin(gameDb, timestamp, playerId) {
						loginNum += 1
					}
				}
				rate = loginNum
				//rate = int(float32(loginNum) / float32(createRoleNum) * 10000)
				m := &RemainTotal{
					//Node:       node,
					PlatformId: platformId,
					ServerId:   serverId,
					Channel:    channel,
					CreateRole: createRoleNum,
					Time:       thatDayZeroTimestamp,
					Remain2:    -1,
					Remain3:    -1,
					Remain4:    -1,
					Remain7:    -1,
					Remain14:   -1,
					Remain30:   -1,
				}
				err = Db.FirstOrCreate(&m).Error
				if err != nil {
					return err
				}
				err = Db.Model(&m).Update("create_role", createRoleNum).Error
				utils.CheckError(err)
				err = Db.Model(&m).Update("register_role", registerCount).Error
				utils.CheckError(err)
				switch i {
				case 1:
					err = Db.Model(&m).Update("Remain2", rate).Error
				case 2:
					err = Db.Model(&m).Update("Remain3", rate).Error
				case 3:
					err = Db.Model(&m).Update("Remain4", rate).Error
				case 6:
					err = Db.Model(&m).Update("Remain7", rate).Error
				case 13:
					err = Db.Model(&m).Update("Remain14", rate).Error
				case 29:
					err = Db.Model(&m).Update("Remain30", rate).Error
				}
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

//func RepireRemainTotal() {
//	g.Log().Info("修复 remain_total")
//	data := make([]*RemainTotal, 0)
//	err := Db.Model(&RemainTotal{}).Find(&data).Error
//	utils.CheckError(err)
//	for _, e := range data {
//
//		nodeName := strings.Split(e.Node, "@")[0]
//
//		platform := "wx"
//		serverId := ""
//		if strings.Contains(nodeName, "_") {
//			platform = strings.Split(nodeName, "_")[0]
//			serverId = strings.Split(nodeName, "_")[1]
//		} else {
//			serverId = nodeName
//		}
//		//e.PlatformId = platform
//		//e.ServerId = serverId
//		//err = Db.Save(&e).Error
//		sql := fmt.Sprintf("update remain_total set platform_id = '%s', server_id = '%s', channel = '%s' where node = '%s';", platform, serverId, platform, e.Node)
//		g.Log().Info("sql:%s\n", sql)
//		err = Db.Debug().Exec(sql).Error
//		//err =Db.Model(&e).Updates(map[string]interface{}{"platform_id": platform, "server_id": serverId, "channel": platform}).Error
//		utils.CheckError(err)
//	}
//	g.Log().Info("修复remain_total完毕")
//}

//func RepireRemainActive() {
//	g.Log().Info("修复 remain_active")
//	data := make([]*RemainActive, 0)
//	err := Db.Model(&RemainActive{}).Find(&data).Error
//	utils.CheckError(err)
//	for _, e := range data {
//
//		nodeName := strings.Split(e.Node, "@")[0]
//
//		platform := "wx"
//		serverId := ""
//		if strings.Contains(nodeName, "_") {
//			platform = strings.Split(nodeName, "_")[0]
//			serverId = strings.Split(nodeName, "_")[1]
//		} else {
//			serverId = nodeName
//		}
//		//e.PlatformId = platform
//		//e.ServerId = serverId
//		//err = Db.Save(&e).Error
//		sql := fmt.Sprintf("update remain_active set platform_id = '%s', server_id = '%s', channel = '%s' where node = '%s';", platform, serverId, platform, e.Node)
//		g.Log().Info("sql:%s\n", sql)
//		err = Db.Debug().Exec(sql).Error
//		//err =Db.Model(&e).Updates(map[string]interface{}{"platform_id": platform, "server_id": serverId, "channel": platform}).Error
//		utils.CheckError(err)
//	}
//	g.Log().Info("修复remain_active完毕")
//}
