package models

import (
	"fmt"
	"gfWeb/library/utils"
	"strings"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
)

type ConsumeStatistics struct {
	PlayerId int `json:"player_id"`
	PropId   int `json:"prop_id"`
	Type     int `json:"type"`
	logType  int `json:"log_type"`
	Value    int `json:"value"`
	SceneId  int `json:"scene_id"`
}

type PropConsumeStatistics struct {
	OpType int     `json:"opType"`
	Count  int     `json:"count"`
	Rate   float32 `json:"rate"`
}
type PropConsumeStatisticsQueryParam struct {
	PlayerName  string
	PlayerId    int
	PlatformId  string
	ServerId    string   `json:"serverId"`
	ChannelList []string `json:"channelList"`
	PromoteId   int
	StartTime   int
	EndTime     int
	// PropType    int
	PropId int
	Type   int
}

func GetPropConsumeStatistics(param *PropConsumeStatisticsQueryParam) ([]*PropConsumeStatistics, error) {
	// if param.PropType == 0 || param.PropId == 0 {
	if param.PropId == 0 {
		return nil, gerror.New("请选择道具")
	}
	gameServer, err := GetGameServerOne(param.PlatformId, param.ServerId)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	node := gameServer.Node
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	defer gameDb.Close()
	list := make([]*PropConsumeStatistics, 0)

	whereArray := make([]string, 0)
	whereArray = append(whereArray, fmt.Sprintf("type = %d", param.Type))
	// whereArray = append(whereArray, fmt.Sprintf("prop_type = %d", param.PropType))
	whereArray = append(whereArray, fmt.Sprintf("prop_id = %d", param.PropId))
	//whereArray = append(whereArray, fmt.Sprintf("server_id = %s", param.ServerId))
	if param.PromoteId == 0 {
		whereArray = append(whereArray, fmt.Sprintf("player_id in (select id from player where server_id = '%s' and channel in(%s))", param.ServerId, GetSQLWhereParam(param.ChannelList)))
	} else {
		promoteData, err := GetPromoteDataOne(param.PromoteId)
		g.Log().Infof("查看数据:%+v", promoteData)
		utils.CheckError(err)
		globalAccountList := make([]*GlobalAccount, 0)
		serverStr := ""
		for serverByte := range []byte(param.ServerId) {
			if serverStr == "" {
				serverStr = fmt.Sprintf("%v", serverByte)
			} else {
				serverStr = serverStr + "," + fmt.Sprintf("%v", serverByte)
			}
		}
		serverStr = "[" + serverStr + "]"
		thisSql := fmt.Sprintf(`select account from global_account where platform_id = '%s' and promote = '%s' and recent_server_list LIKE '%%%s%%'; `, param.PlatformId, promoteData.Promote, serverStr)
		g.Log().Infof("查看thissql:%+v", thisSql)
		err = DbCenter.Raw(thisSql).Scan(&globalAccountList).Error
		g.Log().Infof("查看数据:%+v", globalAccountList)
		utils.CheckError(err)
		accountList := make([]string, 0)
		for _, a := range globalAccountList {
			accountList = append(accountList, a.Account)
		}
		whereArray = append(whereArray, fmt.Sprintf("player_id in (select id from player where server_id = '%s' and channel in(%s) and acc_id in(%s))", param.ServerId, GetSQLWhereParam(param.ChannelList), GetSQLWhereParam(accountList)))
	}
	if param.PlayerId > 0 {
		whereArray = append(whereArray, fmt.Sprintf("player_id = %d", param.PlayerId))
	}

	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}

	order := ""

	if param.Type == 0 {
		order = " order by count desc"
	} else {
		order = "order by count asc"
	}

	sql := fmt.Sprintf(
		` select log_type as op_type, sum(value) as count from consume_statistics  %s group by log_type %s; `, whereParam, order)

	g.Log().Infof("查看sql:%+v", sql)

	err = gameDb.Raw(sql).Scan(&list).Error

	g.Log().Infof("查看list:%+v", list)

	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	var sum = 0
	for _, e := range list {
		sum += e.Count
	}
	for _, e := range list {
		e.Rate = float32(e.Count) / float32(sum) * 100
	}
	return list, nil
}

//
//func GetPropConsumeStatistics(param *PropConsumeStatisticsQueryParam) ([]*PropConsumeStatistics, error) {
//	if param.PropType == 0 || param.PropId == 0 {
//		return nil, errors.New("请选择道具")
//	}
//	gameServer, err := GetGameServerOne(param.PlatformId, param.ServerId)
//	utils.CheckError(err)
//	if err != nil {
//		return nil, err
//	}
//	node := gameServer.Node
//	gameDb, err := GetGameDbByNode(node)
//	utils.CheckError(err)
//	if err != nil {
//		return nil, err
//	}
//	defer gameDb.Close()
//	list := make([]*PropConsumeStatistics, 0)
//
//	var selectPlayer string
//	if param.PlayerId > 0 {
//		selectPlayer = fmt.Sprintf("and player_id = %d", param.PlayerId)
//	}
//
//	order := ""
//
//	if param.Type == 0 {
//		order = "order by count desc"
//	} else {
//		order = "order by count asc"
//	}
//
//	sql := fmt.Sprintf(
//		` select log_type as op_type, sum(value) as count from consume_statistics where type = %d and prop_type = %d and prop_id = %d %s group by log_type %s; `, param.Type, param.PropType, param.PropId, selectPlayer, order)
//	err = gameDb.Raw(sql).Scan(&list).Error
//	utils.CheckError(err)
//	if err != nil {
//		return nil, err
//	}
//	var sum = 0
//	for _, e := range list {
//		sum += e.Count
//	}
//	for _, e := range list {
//		e.Rate = float32(e.Count) / float32(sum) * 100
//	}
//	return list, nil
//}

//func GetPropConsumeStatistics(param *PropConsumeStatisticsQueryParam) ([]*PropConsumeStatistics, error) {
//	if param.PropType == 0 || param.PropId == 0 {
//		return nil, errors.New("请选择道具")
//	}
//	gameDb, err := GetGameDbByNode(param.Node)
//	utils.CheckError(err)
//	if err != nil {
//		return nil, err
//	}
//	defer gameDb.Close()
//	list := make([]*PropConsumeStatistics, 0)
//	var changeValue string
//	if param.Type == 0 {
//		changeValue = "change_value < 0"
//	} else {
//		changeValue = "change_value > 0"
//	}
//	var timeRange string
//	if param.StartTime > 0 {
//		timeRange = fmt.Sprintf("and op_time between %d and %d", param.StartTime, param.EndTime)
//	}
//
//	var selectPlayer string
//	if param.PlayerId > 0 {
//		timeRange = fmt.Sprintf("and player_id = %d", param.PlayerId)
//	}
//
//	sql := fmt.Sprintf(
//		` select op_type, sum(change_value) as count from player_prop_log where %s and prop_type = ? and prop_id = ? %s %s group by op_type; `, changeValue, timeRange, selectPlayer)
//	err = gameDb.Raw(sql, param.PropType, param.PropId).Scan(&list).Error
//	utils.CheckError(err)
//	if err != nil {
//		return nil, err
//	}
//	var sum = 0
//	for _, e := range list {
//		sum += e.Count
//	}
//	for _, e := range list {
//		e.Rate = float32(e.Count) / float32(sum) * 100
//	}
//	return list, nil
//}
