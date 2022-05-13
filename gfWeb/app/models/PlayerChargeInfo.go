package models

import (
	"fmt"
	"gfWeb/library/utils"
	"strconv"

	"github.com/gogf/gf/frame/g"
)

type PlayerChargeInfoRecord struct {
	PlayerId        int     `json:"playerId" gorm:"primary_key"`
	PlayerName      string  `json:"playerName" gorm:"-"`
	Account         string  `json:"account" gorm:"-"`
	PlatformId      string  `json:"platformId" gorm:"column:part_id"`
	ServerId        string  `json:"serverId"`
	TotalMoney      float32 `json:"totalMoney"`
	MaxMoney        float32 `json:"maxMoney"`
	MinMoney        float32 `json:"minMoney"`
	ChargeCount     int     `json:"chargeCount"`
	LastLoginTime   int     `json:"lastLoginTime" gorm:"-"`
	RegisterTime    int     `json:"registerTime" gorm:"-"`
	LastChargeTime  int     `json:"lastChargeTime" gorm:"column:last_time"`
	FirstChargeTime int     `json:"firstChargeTime" gorm:"column:first_time"`
}

type PlayerChargeDataQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId    string   `json:"serverId"`
	ChannelList []string `json:"channelList"`
	StartTime   int      `json:"startTime"`
	EndTime     int      `json:"endTime"`
}

func GetPlayerChargeDataOne(playerId int) (*PlayerChargeInfoRecord, error) {
	playerChargeInfo := &PlayerChargeInfoRecord{
		PlayerId: playerId,
	}
	err := DbCharge.FirstOrInit(&playerChargeInfo).Error
	return playerChargeInfo, err
}

func GetPlayerChargeDataList(params *PlayerChargeDataQueryParam) ([]*PlayerChargeInfoRecord, int64) {
	data := make([]*PlayerChargeInfoRecord, 0)
	var count int64
	sortOrder := "total_money desc"
	g.Log().Debug("params: ", params)
	//if params.Node == "" {
	DbCharge.Model(&PlayerChargeInfoRecord{}).Where(&PlayerChargeInfoRecord{
		PlatformId: params.PlatformId,
		ServerId:   params.ServerId,
	}).Where("charge_count > 0 ").Where("channel in(?)", params.ChannelList).Where(fmt.Sprintf("record_time BETWEEN %d AND %d", params.StartTime, params.EndTime)).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count)
	//} else {
	//	DbCharge.Model(&PlayerChargeInfoRecord{}).Where(&PlayerChargeInfoRecord{
	//		PlatformId: params.PlatformId,
	//	}).Where("charge_count > 0 ").Where("server_id in (?)", GetGameServerIdListByNode(params.Node)).Count(&count).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data)
	//}

	//exchangeRate := GetExchangeRate(params.PlatformId)
	// 转cny
	// cnyExchangeRate := GetExchangeRate("")

	for _, e := range data {
		gameDb, err := GetGameDbByPlatformIdAndSid(e.PlatformId, e.ServerId)
		utils.CheckError(err)
		if err != nil {
			continue
		}
		defer gameDb.Close()
		player, err := GetPlayerByDb(gameDb, e.PlayerId)
		utils.CheckError(err)
		if err != nil {
			continue
		}
		e.PlayerName = player.ServerId + "." + player.Nickname
		e.Account = player.AccId
		e.LastLoginTime = player.LastLoginTime
		e.RegisterTime = player.RegTime

		//TotalMoney, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", e.TotalMoney/exchangeRate), 64)
		TotalMoney, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", e.TotalMoney), 64)
		e.TotalMoney = float32(TotalMoney)

		//MaxMoney, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", e.MaxMoney/exchangeRate), 64)
		MaxMoney, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", e.MaxMoney), 64)
		e.MaxMoney = float32(MaxMoney)

		//MinMoney, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", e.MinMoney/exchangeRate), 64)
		MinMoney, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", e.MinMoney), 64)
		e.MinMoney = float32(MinMoney)
	}
	return data, count
}
