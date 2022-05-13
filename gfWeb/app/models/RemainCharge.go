package models

import (
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
)

type RemainCharge struct {
	//Node       string `json:"node" gorm:"primary_key"`
	PlatformId string `json:"platformId" gorm:"primary_key"`
	ServerId   string `json:"serverId" gorm:"primary_key"`
	Channel    string `json:"channel" gorm:"primary_key"`
	Time       int    `json:"time" gorm:"primary_key"`
	ChargeNum  int    `json:"chargeNum"`
	Remain2    int    `json:"remain2"`
	Remain3    int    `json:"remain3"`
	Remain4    int    `json:"remain4"`
	Remain7    int    `json:"remain7"`
	Remain14   int    `json:"remain14"`
	Remain30   int    `json:"remain30"`
}

type RemainChargeQueryParam struct {
	BaseQueryParam
	PlatformId string
	//ServerId   string
	ServerId    string   `json:"serverId"`
	ChannelList []string `json:"channelList"`
	Channel     string
	StartTime   int
	EndTime     int
}

// 获取付费留存
func GetRemainChargeList(params *RemainChargeQueryParam) []*RemainCharge {
	if params.EndTime < params.StartTime || params.StartTime == 0 {
		g.Log().Error("开始结束时间错误")
		return nil
	}
	data := make([]*RemainCharge, 0, (params.EndTime-params.StartTime)/86400)
	channelLen := len(params.ChannelList)
	for i := params.StartTime; i <= params.EndTime; i = i + 86400 {
		tmpData := make([]*RemainCharge, 0, channelLen)
		err := Db.Model(&RemainCharge{}).Where(&RemainCharge{PlatformId: params.PlatformId, ServerId: params.ServerId, Time: i}).Where("channel in (?)", params.ChannelList).Find(&tmpData).Error
		utils.CheckError(err)
		if len(tmpData) > 0 {
			tmpE := &RemainCharge{
				PlatformId: params.PlatformId,
				ServerId:   params.ServerId,
				Time:       i,
			}
			for _, e := range tmpData {
				tmpE.ChargeNum += e.ChargeNum
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
		if e.ChargeNum > 0 {
			e.Remain2 = int(float32(e.Remain2) / float32(e.ChargeNum) * 10000)
		} else {
			e.Remain2 = 0
		}
		if e.ChargeNum > 0 {
			e.Remain3 = int(float32(e.Remain3) / float32(e.ChargeNum) * 10000)
		} else {
			e.Remain3 = 0
		}
		if e.ChargeNum > 0 {
			e.Remain4 = int(float32(e.Remain4) / float32(e.ChargeNum) * 10000)
		} else {
			e.Remain4 = 0
		}
		if e.ChargeNum > 0 {
			e.Remain7 = int(float32(e.Remain7) / float32(e.ChargeNum) * 10000)
		} else {
			e.Remain7 = 0
		}
		if e.ChargeNum > 0 {
			e.Remain14 = int(float32(e.Remain14) / float32(e.ChargeNum) * 10000)
		} else {
			e.Remain14 = 0
		}
		if e.ChargeNum > 0 {
			e.Remain30 = int(float32(e.Remain30) / float32(e.ChargeNum) * 10000)
		} else {
			e.Remain30 = 0
		}

	}
	return data
	//tmpData := make([]*RemainCharge, 0)
	//data := make([]*RemainCharge, 0)
	//var count int64
	//f := func(db *gorm.DB) *gorm.DB {
	//	if params.StartTime > 0 {
	//		return db.Where("time between ? and ?", params.StartTime, params.EndTime)
	//	}
	//	return db
	//}
	//err := f(Db.Model(&RemainCharge{}).Where(&RemainCharge{PlatformId: params.PlatformId, ServerId: params.ServerId, Channel: params.Channel})).Where(" channel in (?)", params.ChannelList).Count(&count).Offset(params.Offset).Limit(params.Limit).Find(&data).Error
	//utils.CheckError(err)
	//
	//if len(tmpData) > 0 {
	//	tmpE := &DailyStatistics{
	//		PlatformId: params.PlatformId,
	//		ServerId:   params.ServerId,
	//		Time:       i,
	//	}
	//	for _, e := range tmpData {
	//		tmpE.ChargeMoney += e.ChargeMoney
	//	}
	//	data = append(data, tmpE)
	//}
	//
	//for _, e := range tmpData {
	//	chargePlayerIdList := GetChargePlayerIdList(e.PlatformId, e.ServerId, e.Channel, e.Time, e.Time+86400)
	//	chargeNum := len(chargePlayerIdList)
	//	//g.Log().Info("dailyStatistics:%+v", dailyStatistics)
	//	utils.CheckError(err)
	//	e.ChargeNum = chargeNum
	//}
	//return tmpData, count
}

//更新 付费留存
func UpdateRemainCharge(platformId string, serverId string, channelList []*Channel, timestamp int) error {
	g.Log().Info("付费留存:%v, %v, %v, %v", platformId, serverId, len(channelList), timestamp)
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
		for i := 1; i < 30; i++ {
			//if i == 1 || i == 2 || i == 3 || i == 6 || i == 13 || i == 29 {
			if i == 1 {

			} else {
				continue
			}
			thatDayZeroTimestamp := timestamp - i*86400
			if openDayZeroTimestamp > thatDayZeroTimestamp {
				continue
			}
			chargePlayerIdList := GetChargePlayerIdList(platformId, serverId, channel, thatDayZeroTimestamp, thatDayZeroTimestamp+86400)
			chargeNum := len(chargePlayerIdList)
			rate := 0
			loginNum := 0
			if chargeNum > 0 {
				for _, playerId := range chargePlayerIdList {
					if IsThatDayPlayerLogin(gameDb, timestamp, playerId) {
						loginNum += 1
					}
				}
				//rate = int(float32(loginNum) / float32(chargeNum) * 10000)
				rate = loginNum
				m := &RemainCharge{
					//Node:       node,
					PlatformId: platformId,
					ServerId:   serverId,
					Channel:    channel,
					ChargeNum:  chargeNum,
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
