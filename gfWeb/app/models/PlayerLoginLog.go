package models

import (
	"gfWeb/library/utils"
	//"github.com/zaaksam/dproxy/go/db"
	"github.com/jinzhu/gorm"
)

type PlayerLoginLog struct {
	Id         int    `json:"id"`
	PlayerId   int    `json:"playerId"`
	PlayerName string `json:"playerName" gorm:"-"`
	Ip         string `json:"ip"`
	Timestamp  int    `json:"time"`
}

type PlayerLoginLogQueryParam struct {
	BaseQueryParam
	PlatformId string
	ServerId   string `json:"serverId"`
	Ip         string
	PlayerId   int
	PlayerName string
	StartTime  int
	EndTime    int
}

func GetPlayerLoginLogList(params *PlayerLoginLogQueryParam) ([]*PlayerLoginLog, int64) {
	gameServer, err := GetGameServerOne(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	if err != nil {
		return nil, 0
	}
	node := gameServer.Node
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil, 0
	}
	defer gameDb.Close()
	data := make([]*PlayerLoginLog, 0)
	var count int64
	sortOrder := "id"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	//if params.Ip != "" {
	//	gameDb = gameDb.Where("ip = ?", params.Ip)
	//}
	//if params.PlayerId != 0 {
	//	gameDb = gameDb.Where("player_id = ?", params.PlayerId)
	//}
	//if params.StartTime != 0 {
	//	gameDb = gameDb.Where("timestamp >= ?", params.StartTime)
	//}
	//if params.EndTime != 0 {
	//	gameDb = gameDb.Where("timestamp <= ?", params.EndTime)
	//}
	f := func(db *gorm.DB) *gorm.DB {
		if params.StartTime > 0 {
			return db.Where("timestamp between ? and ?", params.StartTime, params.EndTime)
		}
		return db
	}
	f(gameDb.Model(&PlayerLoginLog{}).Where(&PlayerLoginLog{
		Ip:       params.Ip,
		PlayerId: params.PlayerId,
	})).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count)
	for _, e := range data {
		e.PlayerName = GetPlayerName(gameDb, e.PlayerId)
		//e.Ip = e.Ip + "(" + utils.GetIpLocation(e.Ip) + ")"
	}
	return data, count
}
