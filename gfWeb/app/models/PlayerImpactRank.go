package models

import (
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
)

type PlayerImpactRank struct {
	Id         int    `json:"id"`
	Rank       int    `json:"rank"`
	PlayerId   int    `json:"playerId"`
	PlayerName string `json:"playerName" gorm:"-"`
	Value      int    `json:"value"`
	ChangeTime int    `json:"changeTime"`
}

type ImpactRankQueryParam struct {
	BaseQueryParam
	PlatformId string
	ServerId   string `json:"serverId"`
	PlayerId   int    `json:"playerId"`
	PlayerName string `json:"playerName" gorm:"-"`
	ActivityId int    `json:"activityId"`
}

// 查询冲榜排名
func GetImpactRankList(params *ImpactRankQueryParam) ([]*PlayerImpactRank, int64) {
	data := make([]*PlayerImpactRank, 0)
	gameDb, err := GetGameDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	if err != nil {
		//c.HttpResult(r, enums.CodeFail, "未找到服务区服", 0)
		return data, 0
	}
	defer gameDb.Close()

	var count int64
	sortOrder := "rank"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}

	PlayerImpactRankDb := &PlayerImpactRank{
		Id: params.ActivityId,
	}
	if params.PlayerId > 0 {
		PlayerImpactRankDb.PlayerId = params.PlayerId
	}
	err = gameDb.Model(&PlayerImpactRank{}).Where(PlayerImpactRankDb).Order(sortOrder).Offset(params.Offset).Limit(params.Limit).Find(&data).Offset(0).Count(&count).Error
	if err != nil {
		g.Log().Error("查询冲榜排名数据错误:%v", err)
		return data, 0
	}
	for _, e := range data {
		e.PlayerName = GetPlayerName(gameDb, e.PlayerId)
	}
	return data, count
}
