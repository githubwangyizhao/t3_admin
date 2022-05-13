package models

import (
	"github.com/jinzhu/gorm"
)

type BackgroundChargeLog struct {
	Id          int    `json:"id"`
	PlatformId  string `json:"platformId"`
	ServerId    string `json:"serverId"`
	PlayerName  string `json:"playerName" gorm:"-"`
	PlayerId    int    `json:"playerId"`
	Time        int    `json:"time"`
	ChargeType  string `json:"chargeType"`
	ChargeValue int    `json:"chargeValue"`
	ItemId      int    `json:"itemId"`
	UserId      int    `json:"userId"`
	UserName    string `json:"userName" gorm:"-"`
}

type BackgroundChargeLogQueryParam struct {
	BaseQueryParam
	PlatformId string
	ServerId   string `json:"serverId"`
	PlayerName string
	PlayerId   int
	StartTime  int
	EndTime    int
	UserId     int
}

func GetBackgroundChargeLogList(params *BackgroundChargeLogQueryParam) ([]*BackgroundChargeLog, int64) {
	data := make([]*BackgroundChargeLog, 0)
	var count int64
	sortOrder := "time"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	f := func(db *gorm.DB) *gorm.DB {
		if params.StartTime > 0 {
			return db.Where("time between ? and ?", params.StartTime, params.EndTime)
		}
		return db
	}

	//if params.ServerId == "" {
	f(Db.Model(&BackgroundChargeLog{}).Where(&BackgroundChargeLog{
		PlatformId: params.PlatformId,
		PlayerId:   params.PlayerId,
		ServerId:   params.ServerId,
	})).Order(sortOrder).Offset(params.Offset).Limit(params.Limit).Find(&data).Offset(0).Count(&count)
	//} else {
	//	serverIdList := GetGameServerIdListByNode(params.Node)
	//	f(Db.Model(&BackgroundChargeLog{}).Where(&BackgroundChargeLog{
	//		PlatformId: params.PlatformId,
	//		PlayerId:   params.PlayerId,
	//	}).Where("server_id in (?)", serverIdList)).Order(sortOrder).Count(&count).Offset(params.Offset).Limit(params.Limit).Find(&data)
	//}

	for _, e := range data {
		u, err := GetUserOne(e.UserId)
		if err == nil {
			e.UserName = u.Name
		}
		e.PlayerName = GetPlayerName_2(e.PlatformId, e.ServerId, e.PlayerId)
	}
	return data, count
}
