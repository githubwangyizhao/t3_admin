package models

import (
	"fmt"
	"gfWeb/library/utils"
	"strconv"
	"strings"

	"github.com/gogf/gf/frame/g"
)

type PlayerGameLog struct {
	Id         int    `json:"id"`
	PlayerId   int    `json:"playerId"`
	SceneId    int    `json:"sceneId"`
	PlayerName string `json:"playerName" gorm:"-"`
	//CostList   string `json:"costList"`
	CostList  string
	Cost      []*Prop `json:"cost"`
	AwardList string
	Award     []*Prop `json:"award"`
	Time      int     `json:"time"`
	CostTime  int     `json:"costTime"`
	Times     int     `json:"times"`
}

type PlayerGameSceneLogQueryParam struct {
	// BaseQueryParam
	PlatformId     string
	ServerId       string `json:"serverId"`
	PlayerId       int
	SceneId        int `json:"sceneId"`
	PlayerName     string
	StartTime      int
	EndTime        int
	IsChargePlayer int
}

func GetPlayerGameSceneLogList(params *PlayerGameSceneLogQueryParam) []*PlayerGameLog {
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
	data := make([]*PlayerGameLog, 0)
	logData := make([]*PlayerGameLog, 0)
	// sortOrder := "id"
	// if params.Order == "descending" {
	// 	sortOrder = sortOrder + " desc"
	// } else if params.Order == "ascending" {
	// 	sortOrder = sortOrder + " asc"
	// } else {
	// 	sortOrder = sortOrder + " desc"
	// }

	// 重要
	// f := func(db *gorm.DB) *gorm.DB {
	// 	if params.StartTime > 0 {
	// 		return db.Where("time between ? and ?", params.StartTime, params.EndTime)
	// 	}
	// 	return db
	// }
	// f(gameDb.Model(&PlayerGameLog{}).Where(&PlayerGameLog{
	// 	SceneId:  params.SceneId,
	// 	PlayerId: params.PlayerId,
	// })).Find(&data)
	whereArray := make([]string, 0)
	joinStr := ""
	column := "scene_id, cost_list, award_list, cost_time"

	// 0 全部玩家  1 充值过的玩家   2 没充值的玩家
	if params.IsChargePlayer == 1 {
		joinStr = " inner join (select DISTINCT(player_id) as id from player_charge_shop) as player on player_game_log.player_id = player.id "
	}
	if params.IsChargePlayer == 2 {
		whereArray = append(whereArray, fmt.Sprintf(" player.id IS NULL "))
		joinStr = " left join (select DISTINCT(player_id) as id from player_charge_shop) as player on player_id = player.id "
	}

	if params.SceneId > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" scene_id = %d ", params.SceneId))
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
		`select %s from player_game_log %s  %s ; `,
		column,
		joinStr,
		whereParam,
	)

	g.Log().Debug("sql:", sql)

	err = gameDb.Raw(sql).Find(&data).Error

	// g.Log().Debug("data:", data)

	// f(gameDb.Model(&PlayerGameLog{}).Where(&PlayerGameLog{
	// 	SceneId:  params.SceneId,
	// 	PlayerId: params.PlayerId,
	// })).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0)
	for _, e := range data {
		isAdd := 1
		// e.PlayerName = GetPlayerName(gameDb, e.PlayerId)
		e.Award = GetPropStruct(e.AwardList)
		e.Cost = GetPropStruct(e.CostList)
		//e.Ip = e.Ip + "(" + utils.GetIpLocation(e.Ip) + ")"
		for _, log := range logData {
			if log.SceneId == e.SceneId {
				log.CostTime = log.CostTime + e.CostTime
				log.Award = MergeProp(log.Award, e.Award)
				log.Cost = MergeProp(log.Cost, e.Cost)
				log.Times = log.Times + 1
				isAdd = 0
				break
			}
		}
		if isAdd == 1 {
			e.AwardList = ""
			e.CostList = ""
			e.Times = 1
			logData = append(logData, e)
		}
	}
	return logData
}

// func GetPlayerGameSceneLogList(params *PlayerGameSceneLogQueryParam) ([]*PlayerGameLog, int64) {
// 	gameServer, err := GetGameServerOne(params.PlatformId, params.ServerId)
// 	utils.CheckError(err)
// 	if err != nil {
// 		return nil, 0
// 	}
// 	node := gameServer.Node
// 	gameDb, err := GetGameDbByNode(node)
// 	utils.CheckError(err)
// 	if err != nil {
// 		return nil, 0
// 	}
// 	defer gameDb.Close()
// 	data := make([]*PlayerGameLog, 0)
// 	var count int64
// 	sortOrder := "time"
// 	if params.Order == "descending" {
// 		sortOrder = sortOrder + " desc"
// 	} else if params.Order == "ascending" {
// 		sortOrder = sortOrder + " asc"
// 	} else {
// 		sortOrder = sortOrder + " desc"
// 	}

// 	f := func(db *gorm.DB) *gorm.DB {
// 		if params.StartTime > 0 {
// 			return db.Where("time between ? and ?", params.StartTime, params.EndTime)
// 		}
// 		return db
// 	}
// 	f(gameDb.Model(&PlayerGameLog{}).Where(&PlayerGameLog{
// 		SceneId:  params.SceneId,
// 		PlayerId: params.PlayerId,
// 	})).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count)
// 	for _, e := range data {
// 		e.PlayerName = GetPlayerName(gameDb, e.PlayerId)
// 		e.Award = GetPropStruct(e.AwardList)
// 		e.Cost = GetPropStruct(e.CostList)
// 		//e.Ip = e.Ip + "(" + utils.GetIpLocation(e.Ip) + ")"
// 	}
// 	return data, count
// }

// func GetAllGameSceneLogList(params *PlayerGameSceneLogQueryParam) ([]*PlayerGameLog, int64) {
// 	gameServer, err := GetGameServerOne(params.PlatformId, params.ServerId)
// 	utils.CheckError(err)
// 	if err != nil {
// 		return nil, 0
// 	}
// 	node := gameServer.Node
// 	gameDb, err := GetGameDbByNode(node)
// 	utils.CheckError(err)
// 	if err != nil {
// 		return nil, 0
// 	}
// 	defer gameDb.Close()
// 	data := make([]*PlayerGameLog, 0)
// 	var count int64
// 	sortOrder := "time"
// 	if params.Order == "descending" {
// 		sortOrder = sortOrder + " desc"
// 	} else if params.Order == "ascending" {
// 		sortOrder = sortOrder + " asc"
// 	} else {
// 		sortOrder = sortOrder + " desc"
// 	}

// 	f := func(db *gorm.DB) *gorm.DB {
// 		if params.StartTime > 0 {
// 			return db.Where("time between ? and ?", params.StartTime, params.EndTime)
// 		}
// 		return db
// 	}
// 	f(gameDb.Model(&PlayerGameLog{}).Where(&PlayerGameLog{
// 		SceneId:  params.SceneId,
// 		PlayerId: params.PlayerId,
// 	})).Find(&data).Count(&count)
// 	resultMap := make(map[int]interface{})
// 	for _, e := range data {

// 		e.PlayerName = GetPlayerName(gameDb, e.PlayerId)
// 		e.Award = GetPropStruct(e.AwardList)
// 		e.Cost = GetPropStruct(e.CostList)
// 		//e.Ip = e.Ip + "(" + utils.GetIpLocation(e.Ip) + ")"
// 	}
// 	return data, count
// }

// 合并道具 todo 因为t3版删除了类型 ，下面注释待检查
func MergeProp(prop []*Prop, addProp []*Prop) []*Prop {
	for _, addP := range addProp {
		isAdd := 1
		for _, p := range prop {
			// if p.PropType == addP.PropType && p.PropId == addP.PropId {
			if p.PropId == addP.PropId {
				p.PropNum = p.PropNum + addP.PropNum
				isAdd = 0
				break
			}
		}
		if isAdd == 1 {
			prop = append(prop, addP)
		}
	}
	return prop
}

// 获得道具结构体  根据字符串解析
func GetPropStruct(str string) []*Prop {
	var data []*Prop
	if str == "[]" {
		return data
	}
	if strings.HasPrefix(str, "[[") && strings.HasSuffix(str, "]]") {
		itemListStr1 := strings.Replace(strings.Replace(str, "[[", "", -1), "]]", "", -1)
		line := strings.Split(itemListStr1, "],[")
		for _, PropStr := range line {
			PropList := strings.Split(PropStr, ",")
			// PropType, err := strconv.Atoi(PropList[0])
			// utils.CheckError(err)
			PropId, err := strconv.Atoi(PropList[1])
			utils.CheckError(err)
			PropNum, err := strconv.Atoi(PropList[2])
			utils.CheckError(err)
			data = append(data, &Prop{
				// todo v3 del PropType: PropType,
				PropId:  PropId,
				PropNum: PropNum,
			})
		}
	}
	return data
}
