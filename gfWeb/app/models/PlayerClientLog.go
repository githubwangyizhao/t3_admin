package models

import (
	"fmt"
	"gfWeb/library/utils"
	"strings"

	"github.com/gogf/gf/frame/g"
)

type PlayerClientLog struct {
	Id         int    `json:"id"`
	PlayerId   int    `json:"playerId"`
	LogId      int    `json:"logId"`
	PlayerName string `json:"playerName" gorm:"-"`
	Times      int    `json:"times"`
}

type OptLogjsonReq struct {
	JsonName string `json:"json_name"`
}

type PlayerClientLogQueryParam struct {
	// BaseQueryParam
	PlatformId     string
	ServerId       string `json:"serverId"`
	PlayerId       int
	LogId          int `json:"logId"`
	PlayerName     string
	StartTime      int
	EndTime        int
	IsChargePlayer int
}

func ClientLogList(params *PlayerClientLogQueryParam) []*PlayerClientLog {
	gameServer, err := GetGameServerOne(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	if err != nil {
		return nil
	}
	node := gameServer.Node
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil
	}
	defer gameDb.Close()
	data := make([]*PlayerClientLog, 0)
	whereArray := make([]string, 0)
	joinStr := ""
	column := "log_id , COUNT(log_id) as times"

	// 0 全部玩家  1 充值过的玩家   2 没充值的玩家
	if params.IsChargePlayer == 1 {
		joinStr = " inner join (select DISTINCT(player_id) as id from player_charge_shop) as player on player_game_log.player_id = player.id "
	}
	if params.IsChargePlayer == 2 {
		whereArray = append(whereArray, fmt.Sprintf(" player.id IS NULL "))
		joinStr = " left join (select DISTINCT(player_id) as id from player_charge_shop) as player on player_id = player.id "
	}

	if params.LogId > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" log_id = %d ", params.LogId))
	}

	if params.PlayerId > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" player_id = %d ", params.PlayerId))
	}

	if params.StartTime > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" time between %d and %d ", params.StartTime, params.EndTime))
	}

	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}

	sql := fmt.Sprintf(
		`select %s from player_client_log %s  %s group by log_id; `,
		column,
		joinStr,
		whereParam,
	)

	g.Log().Debug("sql:", sql)

	err = gameDb.Raw(sql).Find(&data).Error

	return data
}
