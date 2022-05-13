package models

import (
	"github.com/gogf/gf/frame/g"
)

type PlayerIdPlatformDataParam struct {
	PlayerId   int    `json:"playerId"`
	PlatformId string `json:"platformId"`
	ServerId   string `json:"serverId"`
}

type PlayerIdPlatformQueryParam struct {
	BaseQueryParam
	PlayerIdStrList string
	PlayerIdList    []string
}

// 查询玩家Id平台数据
func GetPlayerIdPlatformList(params *PlayerIdPlatformQueryParam) ([]*PlayerIdPlatformData, int) {
	data := GetPlayerIdPlatformData(params.PlayerIdList)
	len := len(data)
	limit := params.BaseQueryParam.Limit
	start := params.BaseQueryParam.Offset
	if start >= len {
		return nil, len
	}
	if start+limit > len {
		limit = len - start
	}
	g.Log().Debug(len, start, limit)
	return data[start : start+limit], len
}
