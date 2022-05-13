package models

import (
	"fmt"
	"gfWeb/library/utils"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/jinzhu/gorm"
)

//获取某日创角的玩家Id列表
func GetThatDayCreateRolePlayerIdList(db *gorm.DB, serverId string, channel string, zeroTimestamp int) []int {
	var data []struct {
		Id int
	}

	sql := fmt.Sprintf(
		`SELECT id FROM player WHERE reg_time between %d and %d and channel = '%s' and server_id = '%s'`, zeroTimestamp, zeroTimestamp+86400, channel, serverId)
	err := db.Raw(sql).Find(&data).Error
	utils.CheckError(err)
	idList := make([]int, 0)
	for _, e := range data {
		idList = append(idList, e.Id)
	}
	return idList
}

// 是否该玩家某天登录过
func IsThatDayPlayerLogin(db *gorm.DB, zeroTimestamp int, playerId int) bool {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM player_login_log WHERE player_id = %d and timestamp between %d and %d`, playerId, zeroTimestamp, zeroTimestamp+86400)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	if data.Count == 0 {
		return false
	}
	return true
}

// 是否该玩家连续登录过
func IsPlayerContinueLogin(db *gorm.DB, zeroTimestamp int, continueDay int, playerId int) bool {
	for i := 0; i <= continueDay; i++ {
		if IsThatDayPlayerLogin(db, zeroTimestamp+86400*i, playerId) == false {
			return false
		}
	}
	return true
}

func GetExchangeRate(platform string) float32 {
	var exchangeRate float32
	switch platform {
	case "indonesia":
		exchangeRate = 14410
	default:
		exchangeRate = 6
	}
	return exchangeRate
}

func GetSqlSelectParam(args []string) string {
	return strings.Join(args, ", ")
}

func GetSQLWhereParam(args []string) string {
	return "'" + strings.Join(args, "','") + "'"
}

// 查询指定日期的24小时内，登陆过游戏的玩家总数
func getSpecifyDayLoginPlayerCount(db *gorm.DB, serverId string, channel []string, timestamp int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT count(DISTINCT(player_id)) AS count FROM player_login_log WHERE timestamp BETWEEN %d and %d and player_id in (select id from player where channel in (%s) and server_id = '%s')`, timestamp, timestamp+86400, GetSQLWhereParam(channel), serverId)
	g.Log().Infof("getSpecifyDayLoginPlayerCount sql: %s", sql)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

// 获取该天登录次数
func GetThatDayLoginTimes(db *gorm.DB, serverId string, channel string, zeroTimestamp int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM player_login_log WHERE timestamp between %d and %d and player_id in (select id from player where channel = '%s' and server_id = '%s')`, zeroTimestamp, zeroTimestamp+86400, channel, serverId)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

// 获取该天登录的玩家数量
func GetThatDayLoginPlayerCount(db *gorm.DB, serverId string, channel string, zeroTimestamp int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM player WHERE last_login_time between %d and %d and channel = '%s' and server_id = '%s'`, zeroTimestamp, zeroTimestamp+86400, channel, serverId)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

func GetActivePlayersBySpecifiedPeriod(db *gorm.DB, serverId, channel string, startTimestamp int, endTimestamp int) map[string]int {
	endTimestampIn64 := gconv.Int64(endTimestamp)
	zeroEndTimestamp := utils.GetThatZeroTimestamp(endTimestampIn64)
	dayBeforeEndZeroTimestamp := zeroEndTimestamp - 86400

	startTimestampIn64 := gconv.Int64(startTimestamp)
	zeroStartTimestamp := utils.GetThatZeroTimestamp(startTimestampIn64)

	var (
		playerSelect           = "count(p.id) AS login_times, p.nickname AS nickname, from_unixtime(p.reg_time, '%Y-%m-%d %H:%i:%s') AS reg_date"
		playerLoginLogSelect   = " FROM_UNIXTIME(pll.timestamp, '%Y-%m-%d') AS player_login_time"
		playerDataSelect       = " pd.level AS level, pd.vip_level AS vip_level"
		selectFrom             = "player AS p"
		joinPlayerLoginLog     = " LEFT JOIN player_login_log AS pll ON p.id = pll.player_id"
		joinPlayerData         = " LEFT JOIN player_data AS pd ON p.id = pd.player_id"
		channelServerCondition = fmt.Sprintf(" p.channel = '%s' AND p.server_id = '%s' ", channel, serverId)
		regTimeCondition       = fmt.Sprintf(` AND p.reg_time <= '%d'`, dayBeforeEndZeroTimestamp)
		loginTimeCondition     = fmt.Sprintf(` AND pll.timestamp BETWEEN '%d' AND '%d'`, zeroStartTimestamp, zeroEndTimestamp)
		groupBy                = fmt.Sprintf(` GROUP BY player_login_time, pll.player_id`)
		orderBy                = fmt.Sprintf(` ORDER BY pll.timestamp DESC, p.id DESC`)
	)
	if channel == "" {
		channelServerCondition = fmt.Sprintf(" p.server_id = '%s' ", serverId)
	}

	data := make([]*ActivePlayerData, 0)
	sql := fmt.Sprintf(`SELECT %s, %s, %s FROM %s %s %s WHERE %s %s %s %s %s`,
		playerSelect, playerLoginLogSelect, playerDataSelect, selectFrom, joinPlayerLoginLog, joinPlayerData,
		channelServerCondition, regTimeCondition, loginTimeCondition, groupBy, orderBy)
	err := db.Debug().Raw(sql).Scan(&data).Error
	utils.CheckError(err)

	returnData := make(map[string]int, 0)

	for _, item := range data {
		if _, ok := returnData[item.PlayerLoginTime]; ok {
			returnData[item.PlayerLoginTime] += 1
		} else {
			returnData[item.PlayerLoginTime] = 1
		}
	}

	return returnData
}

// GetActivePlayersBySpecifiedDay 查询指定游戏服下所有区服
func GetActivePlayersBySpecifiedDay(db *gorm.DB, serverId, channel string, timestamp int) []*ActivePlayerData {
	timestampIn64 := gconv.Int64(timestamp)
	zeroTimestamp := utils.GetThatZeroTimestamp(timestampIn64)
	dayBeforeZeroTimestamp := zeroTimestamp - 86400

	var (
		playerSelect           = "p.id AS player_id, p.nickname AS nickname, from_unixtime(p.reg_time, '%Y-%m-%d %H:%i:%s') AS reg_date"
		playerLoginLogSelect   = " FROM_UNIXTIME(pll.timestamp, '%Y-%m-%d %H:%i:%s') AS player_login_time"
		playerDataSelect       = " pd.level AS level, pd.vip_level AS vip_level"
		selectFrom             = "player AS p"
		joinPlayerLoginLog     = " LEFT JOIN player_login_log AS pll ON p.id = pll.player_id"
		joinPlayerData         = " LEFT JOIN player_data AS pd ON p.id = pd.player_id"
		channelServerCondition = fmt.Sprintf(" p.channel = '%s' AND p.server_id = '%s' ", channel, serverId)
		regTimeCondition       = fmt.Sprintf(` AND p.reg_time <= '%d'`, dayBeforeZeroTimestamp)
		loginTimeCondition     = fmt.Sprintf(` AND pll.timestamp BETWEEN '%d' AND '%d'`, dayBeforeZeroTimestamp, zeroTimestamp)
		orderBy                = fmt.Sprintf(` ORDER BY p.id`)
	)
	if channel == "" {
		channelServerCondition = fmt.Sprintf(" p.server_id = '%s' ", serverId)
	}

	data := make([]*ActivePlayerData, 0)
	sql := fmt.Sprintf(`SELECT %s, %s, %s FROM %s %s %s WHERE %s %s %s %s`,
		playerSelect, playerLoginLogSelect, playerDataSelect, selectFrom, joinPlayerLoginLog, joinPlayerData,
		channelServerCondition, regTimeCondition, loginTimeCondition, orderBy)
	err := db.Debug().Raw(sql).Scan(&data).Error
	utils.CheckError(err)

	returnData := make([]*ActivePlayerData, 0)
	matchMap := make(map[int]int, 0)
	i := 0

	for _, item := range data {
		if _, ok := matchMap[item.PlayerId]; ok {
			returnData[matchMap[item.PlayerId]].LoginTimes += 1
			returnData[matchMap[item.PlayerId]].LoginTime = item.LoginTime
		} else {
			item.LoginTimes = 1
			returnData = append(returnData, item)
			matchMap[item.PlayerId] = i
			i += 1
		}
	}

	return returnData
}

// GetThatDayActivePlayerCount 获取该天活跃玩家数量
func GetThatDayActivePlayerCount(db *gorm.DB, serverId, channel string, zeroTimestamp int) int {
	count := 0
	data := make([]*Player, 0)
	sql := fmt.Sprintf(
		`SELECT * FROM player WHERE last_login_time between %d and %d and channel = '%s' and server_id = '%s'`, zeroTimestamp, zeroTimestamp+86400, channel, serverId)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	for _, e := range data {
		if e.LoginTimes >= (zeroTimestamp+86400-e.RegTime)/86400 {
			count++
		}
	}
	return count
}

//获取当前在线人数
func GetNowOnlineCountByNode(node string) int {
	gameDb, err := GetGameDbByNode(node)

	utils.CheckError(err)
	if err != nil {
		return -1
	}
	defer gameDb.Close()
	return GetNowOnlineCount(gameDb)
}

//获取当前在线人数
func GetNowOnlineCount(db *gorm.DB) int {
	var count int
	db.Model(&Player{}).Where(&Player{IsOnline: 1}).Count(&count)
	return count
}

//获取当前在线人数
func GetNowOnlineCount2(db *gorm.DB, serverId string, channelList []string) int {
	var count int
	db.Model(&Player{}).Where(&Player{ServerId: serverId, IsOnline: 1}).Where(" channel in(?)", channelList).Count(&count)
	return count
}

//获取当前在线ip数
func GetNowOnlineIpCount(db *gorm.DB, serverId string, channelList []string) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT COUNT( DISTINCT last_login_ip ) as count FROM player where is_online = 1 and server_id = '%s' and channel in(%s);`, serverId, GetSQLWhereParam(channelList))
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取最高在线人数
//func GetMaxOnlineCount(node string) int {
//	var data struct {
//		Count int
//	}
//	sql := fmt.Sprintf(
//		`SELECT max(online_num) as count FROM c_ten_minute_statics WHERE node = '%s' `, node)
//	err := DbCenter.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.Count
//}

//获取时间段内最高在线人数
func GetThatDayMaxOnlineCount(platformId string, serverId string, channelList []string, startTime int, endTime int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT max(online_count) as count FROM ten_minute_statistics WHERE platform_id = '%s' and server_id = '%s' and channel in(%s) and time between %d and %d`, platformId, serverId, GetSQLWhereParam(channelList), startTime, endTime)
	err := Db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取时间段内最低在线人数
func GetThatDayMinOnlineCount(platformId string, serverId string, channelList []string, startTime int, endTime int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT min(online_count) as count FROM ten_minute_statistics WHERE platform_id = '%s' and server_id = '%s' and channel in(%s) and time between %d and %d`, platformId, serverId, GetSQLWhereParam(channelList), startTime, endTime)
	err := Db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取时间段内平均在线人数
func GetThatDayAvgOnlineCount(platformId string, serverId string, channelList []string, startTime int, endTime int) int {
	var data struct {
		Count float32
	}
	sql := fmt.Sprintf(
		`SELECT avg(online_count) as count FROM ten_minute_statistics WHERE platform_id = '%s' and server_id = '%s' and channel in(%s) and time between %d and %d`, platformId, serverId, GetSQLWhereParam(channelList), startTime, endTime)
	err := Db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return int(data.Count)
}

//// 获取当天 平均在线人数
//func GetThatDayAverageOnlineCountByChannel(node string, channel string, zeroTimestamp int) float32 {
//	var data struct {
//		Count float32
//	}
//	sql := fmt.Sprintf(
//		`SELECT avg(online_num)  as count FROM c_ten_minute_statics where node = '%s' and time between %d and %d and channel = '%s'`, node, zeroTimestamp, zeroTimestamp+86400, channel)
//	err := DbCenter.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.Count
//}

//// 获取当天 平均在线人数
//func GetThatDayAverageOnlineCount(node string, zeroTimestamp int) float32 {
//	var data struct {
//		Count float32
//	}
//	sql := fmt.Sprintf(
//		`SELECT avg(online_num)  as count FROM c_ten_minute_statics where node = '%s' and time between %d and %d `, node, zeroTimestamp, zeroTimestamp+86400)
//	err := DbCenter.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.Count
//}

//// 获取当天 最高在线人数
//func GetThatDayMaxOnlineCount(node string, channel string, zeroTimestamp int) int {
//	var data struct {
//		Count int
//	}
//	sql := fmt.Sprintf(
//		`SELECT max(online_num)  as count FROM c_ten_minute_statics where node = '%s' and time between %d and %d and channel = '%s'`, node, zeroTimestamp, zeroTimestamp+86400, channel)
//	err := DbCenter.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.Count
//}

//// 获取当天 最低在线人数
//func GetThatDayMinOnlineCount(node string, channel string, zeroTimestamp int) int {
//	var data struct {
//		Count int
//	}
//	sql := fmt.Sprintf(
//		`SELECT min(online_num)  as count FROM c_ten_minute_statics where node = '%s' and time between %d and %d and channel = '%s'`, node, zeroTimestamp, zeroTimestamp+86400, channel)
//	err := DbCenter.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.Count
//}

//获取该服最高等级
func GetMaxPlayerLevel(db *gorm.DB, serverId string, channelList []string) int {
	var data struct {
		MaxLevel int
	}
	sql := fmt.Sprintf(
		`SELECT max(level) as max_level FROM player_data where player_id in (select id from player where server_id = '%s' and channel in (%s))`, serverId, GetSQLWhereParam(channelList))
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.MaxLevel
}

//获取那天在线时长
func GetOnlineTime(node string, serverId string, channel string, zeroTimestamp int) int {
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return 0
	}
	defer gameDb.Close()
	var data struct {
		Time int
	}
	sql := fmt.Sprintf(
		`SELECT sum(online_time) as time FROM player_online_log where login_time between %d and %d and player_id in (select id from player where channel = '%s' and server_id = '%s')`, zeroTimestamp, zeroTimestamp+86400, channel, serverId)
	err = gameDb.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return int(data.Time)
}

////获取那天平均在线时长
//func GetAvgOnlineTime(node string, channel string, zeroTimestamp int) int {
//	gameDb, err := GetGameDbByNode(node)
//	utils.CheckError(err)
//	if err != nil {
//		return 0
//	}
//	defer gameDb.Close()
//	var data struct {
//		Time float32
//	}
//	sql := fmt.Sprintf(
//		`SELECT avg(online_time) as time FROM player_online_log where login_time between %d and %d and player_id in (select id from player where channel = '%s')`, zeroTimestamp, zeroTimestamp+86400, channel)
//	err = gameDb.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return int(data.Time)
//}

type RemainTask struct {
	TaskId     int     `json:"taskId"`
	Count      int     `json:"count"`
	LeaveCount int     `json:"leaveCount"`
	Rate       float32 `json:"rate"`
}

// 获取任务分布
func GetRemainTask(platformId string, serverId string, channelList []string, isChargePlayer int) []*RemainTask {
	gameServer, err := GetGameServerOne(platformId, serverId)
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

	type Element struct {
		TaskId   int
		PlayerId int
		IsOnline int
		Status   int
	}
	mapData := make(map[int]*RemainTask, 0)
	data := make([]*RemainTask, 0)
	elementList := make([]*Element, 0)

	joinStr := ""
	whereArray := make([]string, 0)

	// 0 全部玩家  1 充值过的玩家   2 没充值的玩家
	if isChargePlayer == 1 {
		joinStr = " inner join (select DISTINCT(player_charge_shop.player_id) as charge_player_id from player_charge_shop) as charge_player on player.id = charge_player.charge_player_id "
	}
	if isChargePlayer == 2 {
		whereArray = append(whereArray, fmt.Sprintf(" charge_player.charge_player_id IS NULL"))
		joinStr = " left join (select DISTINCT(player_charge_shop.player_id) as charge_player_id from player_charge_shop) as charge_player on player.id = charge_player.charge_player_id "
	}

	whereArray = append(whereArray, fmt.Sprintf(" player.channel in (%s) ", GetSQLWhereParam(channelList)))

	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}

	sql := fmt.Sprintf(
		`SELECT player.id, player.is_online, player_task.task_id, player_task.status FROM (player left join player_task on player.id = player_task.player_id) %s %s `,
		joinStr,
		whereParam)
	g.Log().Debug("sql: ", sql)
	err = gameDb.Raw(sql).Find(&elementList).Error
	utils.CheckError(err)
	if err != nil {
		return nil
	}
	for _, e := range elementList {
		if e.Status == 2 {
			e.TaskId = 10000
		}
		if r, ok := mapData[e.TaskId]; ok == true {
			if e.IsOnline == 1 {
				r.Count++
			} else {
				r.Count++
				r.LeaveCount++
			}
		} else {
			if e.IsOnline == 1 {
				mapData[e.TaskId] = &RemainTask{
					TaskId:     e.TaskId,
					Count:      1,
					LeaveCount: 0,
				}
			} else {
				mapData[e.TaskId] = &RemainTask{
					TaskId:     e.TaskId,
					Count:      1,
					LeaveCount: 1,
				}
			}
		}
	}

	var keys []int
	totalCreateRole := GetTotalCreateRoleCountByChannelList(gameDb, serverId, channelList, 0)
	for key, _ := range mapData {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, id := range keys {
		e := mapData[id]
		if totalCreateRole > 0 {
			e.Rate = float32(e.LeaveCount) / float32(totalCreateRole) * 100
		}
		data = append(data, e)
	}
	return data
}

//// 获取任务分布
//func GetRemainTask(platformId int, serverId string) [] *RemainTask{
//	gameDb, err:= GetGameDbByPlatformIdAndSid(platformId, serverId)
//	utils.CheckError(err)
//	defer gameDb.Close()
//	data := make([] *RemainTask, 0)
//	sql := fmt.Sprintf(
//		`SELECT task_id, count(*) as count FROM player_task group by task_id `)
//	err = gameDb.Raw(sql).Find(&data).Error
//
//	totalCreateRole := GetTotalCreateRoleCount(gameDb)
//	for _,e:= range data {
//		e.Rate = float32(e.Count) / float32(totalCreateRole) * 100
//	}
//	return data
//}

type RemainLevel struct {
	Level      int     `json:"level"`
	Count      int     `json:"count"`
	LeaveCount int     `json:"leaveCount"`
	Rate       float32 `json:"rate"`
}

// 获取等级分布
func GetRemainLevel(platformId string, serverId string, channelList []string, startTime int, endTime int, isChargePlayer int) []*RemainLevel {
	gameServer, err := GetGameServerOne(platformId, serverId)
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

	type Element struct {
		Level    int
		PlayerId int
		IsOnline int
	}
	mapData := make(map[int]*RemainLevel, 0)
	data := make([]*RemainLevel, 0)
	elementList := make([]*Element, 0)
	sql := ""
	joinStr := ""
	whereArray := make([]string, 0)

	// 0 全部玩家  1 充值过的玩家   2 没充值的玩家
	if isChargePlayer == 1 {
		joinStr = " inner join (select DISTINCT(player_id) as charge_player_id from player_charge_shop) as charge_player on player_id = charge_player.charge_player_id "
	}
	if isChargePlayer == 2 {
		whereArray = append(whereArray, fmt.Sprintf(" charge_player.charge_player_id IS NULL"))
		joinStr = " left join (select DISTINCT(player_id) as charge_player_id from player_charge_shop) as charge_player on player_id = charge_player.charge_player_id "
	}

	whereArray = append(whereArray, fmt.Sprintf(" player.channel in (%s) ", GetSQLWhereParam(channelList)))

	if startTime > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" reg_time BETWEEN %d AND %d ", startTime, endTime))
	}

	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}

	sql = fmt.Sprintf(
		`SELECT player.id, player.is_online, player_data.level FROM (player left join player_data on player.id = player_data.player_id) %s %s `, joinStr, whereParam)
	g.Log().Debug("sql: ", sql)
	err = gameDb.Raw(sql).Find(&elementList).Error
	utils.CheckError(err)
	if err != nil {
		return nil
	}
	for _, e := range elementList {
		if r, ok := mapData[e.Level]; ok == true {
			if e.IsOnline == 1 {
				r.Count++
			} else {
				r.Count++
				r.LeaveCount++
			}
		} else {
			if e.IsOnline == 1 {
				mapData[e.Level] = &RemainLevel{
					Level:      e.Level,
					Count:      1,
					LeaveCount: 0,
				}
			} else {
				mapData[e.Level] = &RemainLevel{
					Level:      e.Level,
					Count:      1,
					LeaveCount: 1,
				}
			}
		}

	}

	var keys []int
	totalCreateRole := GetTotalCreateRoleCountByChannelList(gameDb, serverId, channelList, 0)
	for key, _ := range mapData {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, id := range keys {
		e := mapData[id]
		e.Rate = float32(e.LeaveCount) / float32(totalCreateRole) * 100
		data = append(data, e)
	}
	return data
	//gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	//utils.CheckError(err)
	//defer gameDb.Close()
	//data := make([] *RemainLevel, 0)
	//sql := fmt.Sprintf(
	//	`SELECT level, count(*) as count FROM player_data group by level `)
	//err = gameDb.Raw(sql).Find(&data).Error
	//utils.CheckError(err)
	//
	//totalCreateRole := GetTotalCreateRoleCount(gameDb)
	//for _, e := range data {
	//	e.Rate = float32(e.Count) / float32(totalCreateRole) * 100
	//}
	//return data
}

type RemainTime struct {
	StartTime  int     `json:"-"`
	EndTime    int64   `json:"-"`
	TimeString string  `json:"timeString"`
	Count      int     `json:"count"`
	LeaveCount int     `json:"leaveCount"`
	Rate       float32 `json:"rate"`
}

// 获取时长分布
func GetRemainTime(platformId string, serverId string, channelList []string, startTime int, endTime int, isChargePlayer int) []*RemainTime {
	gameServer, err := GetGameServerOne(platformId, serverId)
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
	var data = []*RemainTime{
		&RemainTime{
			StartTime:  0,
			EndTime:    60,
			TimeString: "小于1分钟",
		},
		&RemainTime{
			StartTime:  60,
			EndTime:    300,
			TimeString: "1~5分钟",
		},
		&RemainTime{
			StartTime:  300,
			EndTime:    600,
			TimeString: "5~10分钟",
		},
		&RemainTime{
			StartTime:  600,
			EndTime:    1800,
			TimeString: "10~30分钟",
		},
		&RemainTime{
			StartTime:  1800,
			EndTime:    3600,
			TimeString: "30~60分钟",
		},
		&RemainTime{
			StartTime:  3600,
			EndTime:    3600 * 2,
			TimeString: "1~2小时",
		},
		&RemainTime{
			StartTime:  3600 * 2,
			EndTime:    3600 * 3,
			TimeString: "2~3小时",
		},
		&RemainTime{
			StartTime:  3600 * 3,
			EndTime:    3600 * 4,
			TimeString: "3~4小时",
		},
		&RemainTime{
			StartTime:  3600 * 4,
			EndTime:    3600 * 5,
			TimeString: "4~5小时",
		},
		&RemainTime{
			StartTime:  3600 * 5,
			EndTime:    3600 * 6,
			TimeString: "5~6小时",
		},
		&RemainTime{
			StartTime:  3600 * 6,
			EndTime:    3600 * 9,
			TimeString: "6~9小时",
		},
		&RemainTime{
			StartTime:  3600 * 9,
			EndTime:    3600 * 12,
			TimeString: "9~12小时",
		},
		&RemainTime{
			StartTime:  3600 * 12,
			EndTime:    3600 * 24,
			TimeString: "12~24小时",
		},
		&RemainTime{
			StartTime:  3600 * 24,
			EndTime:    3600 * 48,
			TimeString: "1~2天",
		},
		&RemainTime{
			StartTime:  3600 * 48,
			EndTime:    3600 * 72,
			TimeString: "2~3天",
		},
		&RemainTime{
			StartTime:  3600 * 72,
			EndTime:    3600 * 999999,
			TimeString: ">3天",
		},
	}
	type Element struct {
		OnlineTime int
		//PlayerId int
		IsOnline int
	}
	//elementList := make([] *Element, 0)
	//sql := fmt.Sprintf(
	//	`SELECT is_online, total_online_time FROM player`)
	//err = gameDb.Raw(sql).Find(&elementList).Error
	joinStr := ""
	whereArray := make([]string, 0)

	// 0 全部玩家  1 充值过的玩家   2 没充值的玩家
	if isChargePlayer == 1 {
		joinStr = " inner join (select DISTINCT(player_id) as charge_player_id from player_charge_shop) as charge_player on player.id = charge_player.charge_player_id "
	}
	if isChargePlayer == 2 {
		whereArray = append(whereArray, fmt.Sprintf(" charge_player.charge_player_id IS NULL"))
		joinStr = " left join (select DISTINCT(player_id) as charge_player_id from player_charge_shop) as charge_player on player.id = charge_player.charge_player_id "
	}

	whereArray = append(whereArray, fmt.Sprintf(" channel in (%s) ", GetSQLWhereParam(channelList)))

	if startTime > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" reg_time BETWEEN %d AND %d ", startTime, endTime))
	}

	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}
	totalCreateRole := GetTotalCreateRoleCountByChannelList(gameDb, serverId, channelList, 0)
	for _, e := range data {
		elementList := make([]*Element, 0)

		sql := fmt.Sprintf(
			`SELECT is_online, total_online_time FROM player %s %s and total_online_time >= %d and total_online_time < %d `,
			joinStr,
			whereParam,
			e.StartTime,
			e.EndTime)
		g.Log().Debug("sql: ", sql)
		err = gameDb.Raw(sql).Find(&elementList).Error
		utils.CheckError(err)
		if err != nil {
			return nil
		}
		for _, ee := range elementList {
			e.Count++
			if ee.IsOnline == 0 {
				e.LeaveCount++
			}
		}
		if totalCreateRole > 0 {
			e.Rate = float32(e.LeaveCount) / float32(totalCreateRole) * 100
		}

	}
	return data
	//totalCreateRole := GetTotalCreateRoleCount(gameDb)
	//for _, e := range data {
	//	sql := fmt.Sprintf(
	//		`SELECT count(*) as count FROM player where total_online_time >= ? and total_online_time < ? `)
	//	err = gameDb.Raw(sql, e.StartTime, e.EndTime).Find(&e).Error
	//	utils.CheckError(err)
	//	e.Rate = float32(e.Count) / float32(totalCreateRole) * 100
	//}
	//return data
}

func get24hoursOnlineCount(platformId string, serverId string, channelList []string, zeroTimestamp int) ([]string, int) {
	onlineCountList := make([]string, 0, 144)
	now := utils.GetTimestamp()
	//gameServer, _ := GetGameServerOne(platformId, serverId)
	nowOnline := 0

	whereArray := make([]string, 0)
	//whereArray = append(whereArray, fmt.Sprintf("time = %d", i))
	whereArray = append(whereArray, fmt.Sprintf("platform_id = '%s'", platformId))
	whereArray = append(whereArray, fmt.Sprintf("online_count > 0"))
	if serverId != "" {
		whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", serverId))
	}
	whereArray = append(whereArray, fmt.Sprintf("channel in(%s)", GetSQLWhereParam(channelList)))
	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}

	for i := zeroTimestamp; i < zeroTimestamp+86400; i = i + 10*60 {
		if i < now {
			var data struct {
				Sum int
			}
			sql := fmt.Sprintf(
				`SELECT sum(online_count) as sum from ten_minute_statistics %s and time = %d`, whereParam, i)
			err := Db.Raw(sql).Scan(&data).Error
			if err == nil {
				nowOnline = data.Sum
				onlineCountList = append(onlineCountList, strconv.Itoa(data.Sum))
			} else {
				onlineCountList = append(onlineCountList, "null")
			}
		} else {
			onlineCountList = append(onlineCountList, "null")
		}
	}
	return onlineCountList, nowOnline
}

func get24hoursRegisterCount(platformId string, serverId string, channelList []string, zeroTimestamp int) ([]string, int) {
	onlineCountList := make([]string, 0, 144)
	now := utils.GetTimestamp()
	totalCount := 0

	whereArray := make([]string, 0)
	//whereArray = append(whereArray, fmt.Sprintf("time = %d", i))
	whereArray = append(whereArray, fmt.Sprintf("platform_id = '%s'", platformId))
	whereArray = append(whereArray, fmt.Sprintf("register_count > 0"))
	if serverId != "" {
		whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", serverId))
	}
	whereArray = append(whereArray, fmt.Sprintf("channel in(%s)", GetSQLWhereParam(channelList)))
	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}
	for i := zeroTimestamp; i < zeroTimestamp+86400; i = i + 10*60 {
		if i < now {
			var data struct {
				Sum int
			}

			sql := fmt.Sprintf(
				`SELECT sum(register_count) as sum from ten_minute_statistics %s and time = %d`, whereParam, i)
			err := Db.Raw(sql).Scan(&data).Error
			utils.CheckError(err)
			if err == nil {
				totalCount += data.Sum
				onlineCountList = append(onlineCountList, strconv.Itoa(data.Sum))
			} else {
				onlineCountList = append(onlineCountList, "null")
			}
		} else {
			onlineCountList = append(onlineCountList, "null")
		}
	}
	return onlineCountList, totalCount
}

func get24hoursChargeCount(platformId string, serverId string, channelList []string, zeroTimestamp int) ([]string, int) {
	chargeCountList := make([]string, 0, 144)
	now := utils.GetTimestamp()
	totalCount := 0
	//gameServer, _ := GetGameServerOne(platformId, serverId)
	for i := zeroTimestamp + 600; i <= zeroTimestamp+86400; i = i + 10*60 {
		if i < now {
			var data struct {
				Sum int
			}
			whereArray := make([]string, 0)
			whereArray = append(whereArray, fmt.Sprintf("time = %d", i))
			whereArray = append(whereArray, fmt.Sprintf("platform_id = '%s'", platformId))
			if serverId != "" {
				whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", serverId))
			}
			whereArray = append(whereArray, fmt.Sprintf("channel in(%s)", GetSQLWhereParam(channelList)))
			whereParam := strings.Join(whereArray, " and ")
			if whereParam != "" {
				whereParam = " where " + whereParam
			}
			sql := fmt.Sprintf(
				`SELECT sum(charge_count) as sum from ten_minute_statistics %s`, whereParam)
			err := Db.Raw(sql).Scan(&data).Error
			utils.CheckError(err)
			if err == nil {
				totalCount += data.Sum
				chargeCountList = append(chargeCountList, strconv.Itoa(totalCount))
			} else {
				chargeCountList = append(chargeCountList, "null")
			}
		} else {
			chargeCountList = append(chargeCountList, "null")
		}
	}
	//g.Log().Info("%+v", len(onlineCountList))
	return chargeCountList, totalCount
}

func get24hoursChargePlayerCount(platformId string, serverId string, channelList []string, zeroTimestamp int) ([]string, int) {
	chargePlayerCountList := make([]string, 0, 144)
	now := utils.GetTimestamp()
	totalCount := 0
	//gameServer, _ := GetGameServerOne(platformId, serverId)
	for i := zeroTimestamp + 600; i <= zeroTimestamp+86400; i = i + 10*60 {
		if i < now {
			var data struct {
				Sum int
			}
			whereArray := make([]string, 0)
			whereArray = append(whereArray, fmt.Sprintf("time = %d", i))
			whereArray = append(whereArray, fmt.Sprintf("platform_id = '%s'", platformId))
			if serverId != "" {
				whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", serverId))
			}
			whereArray = append(whereArray, fmt.Sprintf("channel in(%s)", GetSQLWhereParam(channelList)))
			whereParam := strings.Join(whereArray, " and ")
			if whereParam != "" {
				whereParam = " where " + whereParam
			}
			sql := fmt.Sprintf(
				`SELECT sum(charge_player_count) as sum from ten_minute_statistics %s`, whereParam)
			err := Db.Raw(sql).Scan(&data).Error
			utils.CheckError(err)
			if err == nil {
				if data.Sum > totalCount {
					totalCount = data.Sum
				}
				chargePlayerCountList = append(chargePlayerCountList, strconv.Itoa(totalCount))
			} else {
				chargePlayerCountList = append(chargePlayerCountList, "null")
			}
		} else {
			chargePlayerCountList = append(chargePlayerCountList, "null")
		}
	}
	//g.Log().Info("%+v", len(onlineCountList))
	return chargePlayerCountList, totalCount
}

//获取玩家名字
func GetPlayerName(db *gorm.DB, playerId int) string {
	var data struct {
		ServerId string
		Name     string
	}

	sql := fmt.Sprintf(
		`SELECT server_id, nickname as name FROM player where id = %d `, playerId)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.ServerId + "." + data.Name
}

//获取玩家最近登录时间
func GetPlayerLastLoginTime(platformId string, serverId string, playerId int) int {
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	if err != nil {
		return 0
	}
	defer gameDb.Close()
	var data struct {
		Time int
	}
	sql := fmt.Sprintf(
		`SELECT last_login_time as time FROM player where id = %d `, playerId)
	err = gameDb.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Time
}

//获取玩家名字
func GetPlayerName_2(platformId string, serverId string, playerId int) string {
	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
	utils.CheckError(err)
	if err != nil {
		return ""
	}
	defer gameDb.Close()
	var data struct {
		ServerId string
		Name     string
	}
	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)

	sql := fmt.Sprintf(
		`SELECT server_id, nickname as name FROM player where id = %d `, playerId)
	//g.Log().Info("GetPlayerName:%v", playerId)
	err = gameDb.Raw(sql).Scan(&data).Error
	//utils.CheckError(err)
	if err != nil {
		return fmt.Sprintf("角色不存在:%s_%d", serverId, playerId)
	}
	//g.Log().Info("ppp:%v,%v", data.MaxLevel)
	return data.ServerId + "." + data.Name
}

//获取区服付费人数
func GetServerChargePlayerCount(platformId string, serverId string, channelList []string, endTime int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where record_time < %d and part_id = '%s' and server_id = '%s' and channel in(%s) and charge_type = 99;`, endTime, platformId, serverId, GetSQLWhereParam(channelList))
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服付费人数
func GetThatDayServerChargePlayerCount(platformId, serverId, channel string, time, source int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and source = %d and ( record_time between %d and %d);`, platformId, serverId, channel, source, time, time+86400)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服付费人数
func GetThatDayChargePlayerCount(platformId string, serverId string, channel string, time, source int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and source = %d and record_time < %d ;`, platformId, serverId, channel, source, time)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//// arpu
//func CaclARPU(totalChargeValueNum int, totalChargePlayerCount int) float32 {
//	if totalChargePlayerCount == 0 {
//		return 0
//	}
//	if totalChargePlayerCount > 0 {
//		return float32(totalChargeValueNum) / float32(totalChargePlayerCount)
//	}
//	return float32(0)
//}

//率
func CaclRate(totalChargePlayerCount int, totalRoleCount int) float32 {
	if totalRoleCount == 0 {
		return 0
	}
	if totalRoleCount > 0 {
		return float32(totalChargePlayerCount) / float32(totalRoleCount)
	}
	return float32(0)
}

func GetNodeListByServerIdList(platformId string, serverIdList []string) []string {
	nodeList := make([]string, 0)
	for _, serverId := range serverIdList {
		gameServer, err := GetGameServerOne(platformId, serverId)
		if err != nil {
			g.Log().Error("获取节点列表失败!!!! %+v, %+v, %+v", platformId, serverId, err)
			return nodeList
		}
		isContain := false
		for _, node := range nodeList {
			if node == gameServer.Node {
				isContain = true
			}
		}
		if isContain == false {
			nodeList = append(nodeList, gameServer.Node)
		}
	}
	return nodeList
}

////二次付费率
//func CaclChargeRate(secondChargePlayerCount int, totalChargePlayerCount int) float32 {
//	if totalChargePlayerCount == 0 {
//		return 0
//	}
//	if totalChargePlayerCount > 0 {
//		return float32(secondChargePlayerCount) / float32(totalChargePlayerCount)
//	}
//	return float32(0)
//}

//获取区服二次付费人数
//func GetServerSecondChargePlayerCount(node string) int {
//	var data struct {
//		Count int
//	}
//	sql := fmt.Sprintf(
//		`select count(DISTINCT player_id) as count from charge_info_record where server_id in (%s) and is_first = 0 and charge_type = 99;`, GetGameServerIdListStringByNode(node))
//	err := DbCharge.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.Count
//}
func GetServerSecondChargePlayerCount(platformId string, serverId string, channelList []string, endTime int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where record_time < %d and part_id = '%s' and server_id = '%s' and channel in(%s) and is_first = 0 and charge_type = 99;`, endTime, platformId, serverId, GetSQLWhereParam(channelList))
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服充值次数列表
func GetServerChargeCountList(platformId string, serverId string, channelList []string, endTime int) []int {
	var data []struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(player_id) as count from charge_info_record where record_time < %d and part_id = '%s' and server_id = '%s' and channel in (%s)  and charge_type = 99 group by player_id;`, endTime, platformId, serverId, GetSQLWhereParam(channelList))
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	r := make([]int, 0)
	//g.Log().Debug("data:%+v", data)
	for _, e := range data {
		r = append(r, e.Count)
	}
	//g.Log().Debug("r:%+v", r)
	return r
}

//获取区服首次付费人数
func GetThadDayServerFirstChargePlayerCount(platformId, serverId, channel string, time, source int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and is_first = 1 and charge_type = 99 and source = %d and (record_time between %d and %d);`, platformId, serverId, channel, source, time, time+86400)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服首次金额
func GetThadDayServerFirstChargeTotalMoney(platformId, serverId, channel string, time, source int) float32 {
	var data struct {
		Count float32
	}
	sql := fmt.Sprintf(
		`select sum(money) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and is_first = 1 and charge_type = 99 and source = %d and (record_time between %d and %d);`, platformId, serverId, channel, source, time, time+86400)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服新增付费人数
func GetThadDayNewChargePlayerCount(platformId, serverId, channel string, time, source int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and is_first = 1 and charge_type = 99 and source = %d and (record_time between %d and %d) and (reg_time between %d and %d);`, platformId, serverId, channel, source, time, time+86400, time, time+86400)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取时间内充值的玩家id列表
func GetChargePlayerIdList(platformId, serverId, channel string, startTime int, endTime int) []int {
	var data []struct {
		Id int
	}
	sql := fmt.Sprintf(
		`select DISTINCT player_id as id from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s'  and charge_type = 99 and (record_time between %d and %d) ;`, platformId, serverId, channel, startTime, endTime)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	idList := make([]int, 0)
	for _, e := range data {
		idList = append(idList, e.Id)
	}
	return idList
}

//获取区服总充值元宝
func GetServerTotalChargeIngot(platformId string, serverId string, channelList []string) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select sum(ingot) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel in(%s) and charge_type = 99;`, platformId, serverId, GetSQLWhereParam(channelList))
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取该时间内区服总充值人民币
func GetTotalChargeMoney(platformId, serverId, channel string, startTime, endTime, source int) float32 {
	var data struct {
		Count float32
	}
	sql := fmt.Sprintf(
		`select sum(money) as count from charge_info_record where status = 1 and  part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and source = %d and record_time between %d and %d;`, platformId, serverId, channel, source, startTime, endTime)

	if source == 0 {
		sql = fmt.Sprintf(
			`select sum(money) as count from charge_info_record where status = 1 and  part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and record_time between %d and %d;`, platformId, serverId, channel, startTime, endTime)
	}
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

func GetChargePlayerCount(platformId string, serverId string, channel string, startTime int, endTime int) int {
	var data struct {
		Count int
	}
	whereArr := make([]string, 0)
	whereArr = append(whereArr, fmt.Sprintf(`status = 1 and charge_type = 99`))
	if platformId != "" {
		whereArr = append(whereArr, fmt.Sprintf(`part_id = '%s'`, platformId))
	}
	if serverId != "" {
		whereArr = append(whereArr, fmt.Sprintf(`server_id = '%s'`, serverId))
	}
	if channel != "" {
		whereArr = append(whereArr, fmt.Sprintf(`channel in (%s)`, channel))
	}
	if startTime != 0 && endTime > startTime {
		whereArr = append(whereArr, fmt.Sprintf(`record_time between '%d' and '%d'`, startTime, endTime))
	}

	whereStr := ""
	if len(whereArr) >= 0 {
		whereStr = " WHERE " + strings.Join(whereArr, " AND ")
	}
	sql := fmt.Sprintf(
		`SELECT COUNT(DISTINCT player_id) AS count FROM charge_info_record %s`, whereStr)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取该时间内区服充值人数
func GetTotalChargePlayerCount(platformId string, serverId string, channel string, startTime int, endTime int) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`select count(DISTINCT player_id) as count from charge_info_record where status = 1 and part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and record_time between %d and %d;`, platformId, serverId, channel, startTime, endTime)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取某天注册的玩家 该时间区间内总充值人民币
func GetTotalChargeMoneyByRegisterTime(platformId string, serverId string, channel string, startTime int, endTime int, registerTime int) float32 {
	var data struct {
		Count float32
	}
	sql := fmt.Sprintf(
		`select sum(money) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and record_time between %d and %d and reg_time between %d and %d;`, platformId, serverId, channel, startTime, endTime, registerTime, registerTime+86400)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服总充值人民币
func GetServerTotalChargeMoneyByChannelList(platformId string, serverId string, channelList []string, endTime int) float32 {
	var data struct {
		Count float32
	}
	sql := fmt.Sprintf(
		`select sum(money) as count from charge_info_record where record_time < %d and part_id = '%s' and server_id = '%s' and channel in(%s) and charge_type = 99;`, endTime, platformId, serverId, GetSQLWhereParam(channelList))
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

//获取区服总充值人民币
//func GetTotalChargeMoney(platformId string, serverId string, channel string, time int) int {
//	var data struct {
//		Count float32
//	}
//	sql := fmt.Sprintf(
//		`select sum(money) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and record_time < %d;`, platformId, serverId, channel, time)
//	err := DbCharge.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return int(data.Count)
//}

func GetThatDayNewChargeMoney(platformId string, serverId string, channel string, time, source int) float32 {
	var data struct {
		Count float32
	}
	sql := fmt.Sprintf(
		`select sum(money) as count from charge_info_record where part_id = '%s' and server_id = '%s' and channel = '%s' and charge_type = 99 and source = %d and (record_time between %d and %d) and (reg_time between %d and %d);`, platformId, serverId, channel, source, time, time+86400, time, time+86400)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}

func SplitPlayerName(playerName string) (string, string, error) {
	array := strings.Split(playerName, ".")
	if len(array) == 2 {
		return array[0], array[1], nil
	}
	return "", "", gerror.New(fmt.Sprintf("解析玩家名字失败:%s", playerName))
}

type ChargeTaskDistribution struct {
	TaskId int     `json:"taskId"`
	Count  int     `json:"count"`
	Rate   float32 `json:"rate"`
}
type ChargeTaskDistributionQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId    string   `json:"serverId"`
	ChannelList []string `json:"channelList"`
	StartTime   int
	EndTime     int
	IsFirst     int
	Promote     string `json:"promote"`
}

// 获取充值任务分布
func GetChargeTaskDistribution(params ChargeTaskDistributionQueryParam) []*ChargeTaskDistribution {
	data := make([]*ChargeTaskDistribution, 0)
	whereArray := make([]string, 0)
	whereArray = append(whereArray, fmt.Sprintf("charge_type = 99"))
	whereArray = append(whereArray, fmt.Sprintf(" part_id = '%s'", params.PlatformId))
	whereArray = append(whereArray, fmt.Sprintf(" channel in (%s) ", GetSQLWhereParam(params.ChannelList)))

	g.Log().Debug("params: ", params)

	if params.Promote != "" {
		AccIdList, AccErr := GetGlobalAccountByPromote(params.Promote)
		utils.CheckError(AccErr)
		g.Log().Debug("AccIdList:", params.Promote, AccIdList)

		if len(AccIdList) > 0 {
			whereArray = append(whereArray, fmt.Sprintf(" acc_id in (%s) ", GetSQLWhereParam(AccIdList)))
		}
	}

	if params.IsFirst == 1 {
		whereArray = append(whereArray, fmt.Sprintf("is_first = 1"))
	}
	if params.ServerId != "" {
		whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", params.ServerId))
	}
	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}

	sql := fmt.Sprintf(
		`SELECT curr_task_id as task_id, count(*) as count FROM charge_info_record %s group by task_id `, whereParam)
	g.Log().Debug("sql: ", sql)
	err := DbCharge.Raw(sql).Find(&data).Error
	utils.CheckError(err)

	//if params.Node == "" {
	//	sql := fmt.Sprintf(
	//		`SELECT curr_task_id as task_id, count(*) as count FROM charge_info_record where charge_type = 99 group by task_id `)
	//	err := DbCharge.Raw(sql).Find(&data).Error
	//	utils.CheckError(err)
	//} else {
	//	serverIdList := GetGameServerIdListStringByNode(params.Node)
	//	sql := fmt.Sprintf(
	//		`SELECT curr_task_id as task_id, count(*) as count FROM charge_info_record where server_id in (%s) and charge_type = 99 group by task_id `, serverIdList)
	//	err := DbCharge.Raw(sql).Find(&data).Error
	//	utils.CheckError(err)
	//}

	sum := 0
	for _, e := range data {
		sum += e.Count
	}
	if sum > 0 {
		for _, e := range data {
			e.Rate = float32(e.Count) / float32(sum) * 100
		}
	}

	return data
}

type ChargeActivityDistribution struct {
	ChargeItemId int     `json:"chargeItemId"`
	Count        int     `json:"count"`
	Rate         float32 `json:"rate"`
	Money        float32 `json:"money"`
	MoneyRate    float32 `json:"moneyRate"`
}
type ChargeActivityDistributionQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId    string   `json:"serverId"`
	ChannelList []string `json:"channelList"`
	StartTime   int
	EndTime     int
	IsFirst     int
	Promote     string `json:"promote"`
}

// 获取充值任务分布
func GetChargeActivityDistribution(params ChargeActivityDistributionQueryParam) []*ChargeActivityDistribution {
	data := make([]*ChargeActivityDistribution, 0)
	whereArray := make([]string, 0)
	whereArray = append(whereArray, fmt.Sprintf("charge_type = 99"))
	whereArray = append(whereArray, fmt.Sprintf(" part_id = '%s'", params.PlatformId))
	whereArray = append(whereArray, fmt.Sprintf(" channel in (%s) ", GetSQLWhereParam(params.ChannelList)))
	if params.IsFirst == 1 {
		whereArray = append(whereArray, fmt.Sprintf("is_first = 1"))
	}
	if params.StartTime > 0 {
		whereArray = append(whereArray, fmt.Sprintf("record_time between %d and %d", params.StartTime, params.EndTime))
	}
	if params.ServerId != "" {
		whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", params.ServerId))
	}
	if params.Promote != "" {
		AccIdList, AccErr := GetGlobalAccountByPromote(params.Promote)
		utils.CheckError(AccErr)
		g.Log().Debug("AccIdList:", params.Promote, AccIdList)

		if len(AccIdList) > 0 {
			whereArray = append(whereArray, fmt.Sprintf(" acc_id in (%s) ", GetSQLWhereParam(AccIdList)))
		} else {
			return data
		}
	}
	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}

	sql := fmt.Sprintf(
		`SELECT charge_item_id, count(*) as count, sum(money) as money FROM charge_info_record %s group by charge_item_id `, whereParam)
	g.Log().Debug("sql: ", sql)
	err := DbCharge.Raw(sql).Find(&data).Error
	utils.CheckError(err)

	//if params.Node == "" {
	//	sql := fmt.Sprintf(
	//		`SELECT curr_task_id as task_id, count(*) as count FROM charge_info_record where charge_type = 99 group by task_id `)
	//	err := DbCharge.Raw(sql).Find(&data).Error
	//	utils.CheckError(err)
	//} else {
	//	serverIdList := GetGameServerIdListStringByNode(params.Node)
	//	sql := fmt.Sprintf(
	//		`SELECT curr_task_id as task_id, count(*) as count FROM charge_info_record where server_id in (%s) and charge_type = 99 group by task_id `, serverIdList)
	//	err := DbCharge.Raw(sql).Find(&data).Error
	//	utils.CheckError(err)
	//}

	sum := 0
	var moneySum float32
	for _, e := range data {
		sum += e.Count
		moneySum += e.Money
	}
	if sum > 0 {
		for _, e := range data {
			e.Rate = float32(e.Count) / float32(sum) * 100
		}
	}
	if moneySum > 0 {
		for _, e := range data {
			e.MoneyRate = float32(e.Money) / float32(moneySum) * 100
		}
	}
	return data
}

type ChargeMoneyDistribution struct {
	ValueString string  `json:"valueString"`
	Count       int     `json:"count"`
	Rate        float32 `json:"rate"`
	Min         int     `json:"-"`
	Max         int64   `json:"-"`
}
type ChargeMoneyDistributionV1 struct {
	Count  int `json:"count"`
	Level0 int
	Level2 int
}
type ChargeMoneyDistributionDataInDb struct {
	Count   int
	Level0  int
	Level1  int
	Level2  int
	Level3  int
	Level4  int
	Level5  int
	Level6  int
	Level7  int
	Level8  int
	Level9  int
	Level10 int
	Level11 int
	Level12 int
	Level13 int
}
type ChargeMoneyDistributionQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId    string   `json:"serverId"`
	ChannelList []string `json:"channelList"`
	StartTime   int
	EndTime     int
	Promote     string `json:"promote"`
}

func GetChargeMoneyDistributionV1(params ChargeMoneyDistributionQueryParam) []*ChargeMoneyDistribution {
	var data = &ChargeLevelDistributionDataInDb{
		Count:   0,
		Level0:  0,
		Level1:  0,
		Level2:  0,
		Level3:  0,
		Level4:  0,
		Level5:  0,
		Level6:  0,
		Level7:  0,
		Level8:  0,
		Level9:  0,
		Level10: 0,
		Level11: 0,
		Level12: 0,
		Level13: 0,
	}
	g.Log().Debug("data，", data)
	var matchData = []*ChargeMoneyDistribution{
		&ChargeMoneyDistribution{
			Min:         0,
			Max:         2,
			ValueString: "2",
		},
		&ChargeMoneyDistribution{
			Min:         3,
			Max:         8,
			ValueString: "3~8",
		},
		&ChargeMoneyDistribution{
			Min:         8,
			Max:         20,
			ValueString: "8~20",
		},
		&ChargeMoneyDistribution{
			Min:         20,
			Max:         50,
			ValueString: "20~50",
		},
		&ChargeMoneyDistribution{
			Min:         50,
			Max:         100,
			ValueString: "50~100",
		},
		&ChargeMoneyDistribution{
			Min:         100,
			Max:         200,
			ValueString: "100~200",
		},
		&ChargeMoneyDistribution{
			Min:         200,
			Max:         500,
			ValueString: "200~500",
		},
		&ChargeMoneyDistribution{
			Min:         500,
			Max:         1000,
			ValueString: "500~1000",
		},
		&ChargeMoneyDistribution{
			Min:         1000,
			Max:         2000,
			ValueString: "1000~2000",
		},
		&ChargeMoneyDistribution{
			Min:         2000,
			Max:         5000,
			ValueString: "2000~5000",
		},
		&ChargeMoneyDistribution{
			Min:         5000,
			Max:         10000,
			ValueString: "5000~10000",
		},
		&ChargeMoneyDistribution{
			Min:         10000,
			Max:         20000,
			ValueString: "1万~2万",
		},
		&ChargeMoneyDistribution{
			Min:         20001,
			Max:         100000000000,
			ValueString: "大于2万",
		},
	}
	selectArray := make([]string, 0)
	for matchKey, matchVal := range matchData {
		selectArray = append(selectArray, fmt.Sprintf("IFNULL(SUM(if(total_money > %d and total_money <= %d, 1, 0)), 0) as level%d", matchVal.Min, matchVal.Max, matchKey))
	}
	selectString := GetSqlSelectParam(selectArray)
	whereArray := make([]string, 0)
	if params.Promote != "" {
		AccIdList, AccErr := GetGlobalAccountByPromote(params.Promote)
		utils.CheckError(AccErr)
		g.Log().Debug("AccIdList:", params.Promote, AccIdList)

		if len(AccIdList) > 0 {
			whereArray = append(whereArray, fmt.Sprintf(" acc_id in (%s) ", GetSQLWhereParam(AccIdList)))
		} else {
			return matchData
		}
	}

	if params.ServerId != "" {
		whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", params.ServerId))
	}
	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}
	sql := fmt.Sprintf(`SELECT count(*) as count, %s FROM player_charge_info_record %s`, selectString, whereParam)
	g.Log().Debug("sql: ", sql)
	err := DbCharge.Raw(sql).Find(&data).Error
	utils.CheckError(err)

	//exchangeRate := GetExchangeRate(params.PlatformId)
	//exchangeRate = exchangeRate / GetExchangeRate("")

	Reflect := reflect.ValueOf(data).Elem()
	for matchKey, _ := range matchData {
		name := fmt.Sprintf("Level%d", matchKey)
		if Reflect.FieldByName(name).IsValid() {
			count := Reflect.FieldByName(name).Int()
			matchData[matchKey].Count = int(float32(count)) // / exchangeRate)
		} else {
			matchData[matchKey].Count = 0
		}

		if data.Count > 9 {
			matchData[matchKey].Rate = float32(matchData[matchKey].Count) / float32(data.Count) * 100
		}
	}

	return matchData
}

// 获取充值金额分布
func GetChargeMoneyDistribution(params ChargeMoneyDistributionQueryParam) []*ChargeMoneyDistribution {
	var data = []*ChargeMoneyDistribution{
		&ChargeMoneyDistribution{
			Min:         0,
			Max:         1,
			ValueString: "1元",
		},
		&ChargeMoneyDistribution{
			Min:         1,
			Max:         5,
			ValueString: "1~5元",
		},
		&ChargeMoneyDistribution{
			Min:         5,
			Max:         10,
			ValueString: "5~10元",
		},
		&ChargeMoneyDistribution{
			Min:         10,
			Max:         20,
			ValueString: "10~20元",
		},
		&ChargeMoneyDistribution{
			Min:         20,
			Max:         50,
			ValueString: "20~50元",
		},
		&ChargeMoneyDistribution{
			Min:         50,
			Max:         100,
			ValueString: "50~100",
		},
		&ChargeMoneyDistribution{
			Min:         100,
			Max:         200,
			ValueString: "100~200元",
		},
		&ChargeMoneyDistribution{
			Min:         200,
			Max:         500,
			ValueString: "200~500元",
		},
		&ChargeMoneyDistribution{
			Min:         500,
			Max:         1000,
			ValueString: "500~1000元",
		},
		&ChargeMoneyDistribution{
			Min:         1000,
			Max:         2000,
			ValueString: "1000~2000元",
		},
		&ChargeMoneyDistribution{
			Min:         2000,
			Max:         5000,
			ValueString: "2000~5000元",
		},
		&ChargeMoneyDistribution{
			Min:         5000,
			Max:         10000,
			ValueString: "5000~10000元",
		},
		&ChargeMoneyDistribution{
			Min:         10000,
			Max:         20000,
			ValueString: "1万~2万元",
		},
		&ChargeMoneyDistribution{
			Min:         20000,
			Max:         50000,
			ValueString: "2万~5万元",
		},
		&ChargeMoneyDistribution{
			Min:         50001,
			Max:         100000,
			ValueString: "5万~10万元",
		},
		&ChargeMoneyDistribution{
			Min:         100001,
			Max:         100000000000,
			ValueString: "大于10万元",
		},
	}
	maxCount := 0
	for _, e := range data {
		whereArray := make([]string, 0)
		whereArray = append(whereArray, fmt.Sprintf(" part_id = '%s'", params.PlatformId))
		whereArray = append(whereArray, fmt.Sprintf(" total_money > %d", e.Min))
		whereArray = append(whereArray, fmt.Sprintf(" total_money <= %d", e.Max))
		whereArray = append(whereArray, fmt.Sprintf(" channel in (%s) ", GetSQLWhereParam(params.ChannelList)))
		if params.ServerId != "" {
			whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", params.ServerId))
		}
		whereParam := strings.Join(whereArray, " and ")
		if whereParam != "" {
			whereParam = " where " + whereParam
		}

		sql := fmt.Sprintf(`SELECT count(*) as count FROM player_charge_info_record %s`, whereParam)
		err := DbCharge.Raw(sql).Find(&e).Error
		utils.CheckError(err)
	}
	for _, e := range data {
		maxCount += e.Count
	}
	if maxCount > 0 {
		for _, e := range data {
			e.Rate = float32(e.Count) / float32(maxCount) * 100
		}
	}

	return data
}

type ChargeLevelDistribution struct {
	LevelString string  `json:"levelString"`
	Count       int     `json:"count"`
	Rate        float32 `json:"rate"`
	Min         int     `json:"-"`
	Max         int     `json:"-"`
}
type ChargeLevelDistributionDataInDb struct {
	Count   int
	Level0  int
	Level1  int
	Level2  int
	Level3  int
	Level4  int
	Level5  int
	Level6  int
	Level7  int
	Level8  int
	Level9  int
	Level10 int
	Level11 int
	Level12 int
	Level13 int
	Level14 int
	Level15 int
	Level16 int
	Level17 int
	Level18 int
	Level19 int
	Level20 int
	Level21 int
	Level22 int
	Level23 int
	Level24 int
	Level25 int
}
type ChargeLevelDistributionV1 struct {
	LevelString string  `json:"levelString"`
	Count       int     `json:"count"`
	Rate        float32 `json:"rate"`
	Min         int     `json:"min"`
	Max         int     `json:"max"`
}
type ChargeLevelDistributionQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId    string   `json:"serverId"`
	ChannelList []string `json:"channelList"`
	StartTime   int
	EndTime     int
	IsFirst     int
	Promote     string `json:"promote"`
}

func GetChargeLevelDistributionV1(params ChargeLevelDistributionQueryParam) []*ChargeLevelDistributionV1 {
	var data = &ChargeLevelDistributionDataInDb{
		Count:   0,
		Level0:  0,
		Level1:  0,
		Level2:  0,
		Level3:  0,
		Level4:  0,
		Level5:  0,
		Level6:  0,
		Level7:  0,
		Level8:  0,
		Level9:  0,
		Level10: 0,
		Level11: 0,
		Level12: 0,
		Level13: 0,
		Level14: 0,
		Level15: 0,
		Level16: 0,
		Level17: 0,
		Level18: 0,
		Level19: 0,
		Level20: 0,
		Level21: 0,
		Level22: 0,
		Level23: 0,
		Level24: 0,
		Level25: 0,
	}
	var matchData = []*ChargeLevelDistributionV1{
		&ChargeLevelDistributionV1{
			Min:         1,
			Max:         5,
			LevelString: "1~5级",
		},
		&ChargeLevelDistributionV1{
			Min:         6,
			Max:         10,
			LevelString: "6~10级",
		},
		&ChargeLevelDistributionV1{
			Min:         11,
			Max:         20,
			LevelString: "11~20级",
		},
		&ChargeLevelDistributionV1{
			Min:         21,
			Max:         30,
			LevelString: "21~30级",
		},
		&ChargeLevelDistributionV1{
			Min:         31,
			Max:         40,
			LevelString: "31~40级",
		},
		&ChargeLevelDistributionV1{
			Min:         41,
			Max:         50,
			LevelString: "41~50级",
		},
		&ChargeLevelDistributionV1{
			Min:         51,
			Max:         60,
			LevelString: "51~60级",
		},
		&ChargeLevelDistributionV1{
			Min:         61,
			Max:         70,
			LevelString: "61~70级",
		},
		&ChargeLevelDistributionV1{
			Min:         71,
			Max:         80,
			LevelString: "71~80级",
		},
		&ChargeLevelDistributionV1{
			Min:         81,
			Max:         90,
			LevelString: "81~90级",
		},
		&ChargeLevelDistributionV1{
			Min:         91,
			Max:         100,
			LevelString: "91~100级",
		},
		&ChargeLevelDistributionV1{
			Min:         101,
			Max:         110,
			LevelString: "101~110级",
		},
		&ChargeLevelDistributionV1{
			Min:         111,
			Max:         120,
			LevelString: "111~120级",
		}, &ChargeLevelDistributionV1{
			Min:         121,
			Max:         130,
			LevelString: "121~130级",
		},
		&ChargeLevelDistributionV1{
			Min:         131,
			Max:         140,
			LevelString: "131~140级",
		},
		&ChargeLevelDistributionV1{
			Min:         141,
			Max:         150,
			LevelString: "141~150级",
		},
		&ChargeLevelDistributionV1{
			Min:         151,
			Max:         160,
			LevelString: "151~160级",
		},
		&ChargeLevelDistributionV1{
			Min:         161,
			Max:         170,
			LevelString: "161~170级",
		},
		&ChargeLevelDistributionV1{
			Min:         171,
			Max:         180,
			LevelString: "171~180级",
		},
		&ChargeLevelDistributionV1{
			Min:         181,
			Max:         190,
			LevelString: "181~190级",
		},
		&ChargeLevelDistributionV1{
			Min:         191,
			Max:         200,
			LevelString: "191~200级",
		},
		&ChargeLevelDistributionV1{
			Min:         201,
			Max:         250,
			LevelString: "201~250级",
		},
		&ChargeLevelDistributionV1{
			Min:         251,
			Max:         300,
			LevelString: "251~300级",
		},
		&ChargeLevelDistributionV1{
			Min:         301,
			Max:         400,
			LevelString: "301~400级",
		},
		&ChargeLevelDistributionV1{
			Min:         401,
			Max:         500,
			LevelString: "401~500级",
		},
		&ChargeLevelDistributionV1{
			Min:         501,
			Max:         10000,
			LevelString: "大于501级",
		},
	}
	selectArray := make([]string, 0)
	for matchKey, matchVal := range matchData {
		selectArray = append(selectArray, fmt.Sprintf("IFNULL(SUM(if(curr_level >= %d and curr_level <= %d, 1, 0)), 0) as level%d", matchVal.Min, matchVal.Max, matchKey))
	}
	selectString := GetSqlSelectParam(selectArray)

	whereArray := make([]string, 0)
	if params.Promote != "" {
		AccIdList, AccErr := GetGlobalAccountByPromote(params.Promote)
		utils.CheckError(AccErr)
		g.Log().Debug("AccIdList:", params.Promote, AccIdList)

		if len(AccIdList) > 0 {
			whereArray = append(whereArray, fmt.Sprintf(" acc_id in (%s) ", GetSQLWhereParam(AccIdList)))
		} else {
			return matchData
		}
	}

	if params.IsFirst == 1 {
		whereArray = append(whereArray, fmt.Sprintf("is_first = 1"))
	}
	if params.ServerId != "" {
		whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", params.ServerId))
	}
	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}
	sql := fmt.Sprintf(`SELECT count(*) as count, %s FROM charge_info_record %s`, selectString, whereParam)
	g.Log().Debug("sql: ", sql)
	err := DbCharge.Raw(sql).Scan(&data).Error
	utils.CheckError(err)

	Reflect := reflect.ValueOf(data).Elem()
	for matchKey, _ := range matchData {
		name := fmt.Sprintf("Level%d", matchKey)
		if Reflect.FieldByName(name).IsValid() {
			count := Reflect.FieldByName(name).Int()
			matchData[matchKey].Count = int(count)
		} else {
			matchData[matchKey].Count = 0
		}
	}

	return matchData
}

// 获取充值等级分布
func GetChargeLevelDistribution(params ChargeLevelDistributionQueryParam) []*ChargeLevelDistribution {
	var data = []*ChargeLevelDistribution{
		&ChargeLevelDistribution{
			Min:         1,
			Max:         5,
			LevelString: "1~5级",
		},
		&ChargeLevelDistribution{
			Min:         6,
			Max:         10,
			LevelString: "6~10级",
		},
		&ChargeLevelDistribution{
			Min:         11,
			Max:         20,
			LevelString: "11~20级",
		},
		&ChargeLevelDistribution{
			Min:         21,
			Max:         30,
			LevelString: "21~30级",
		},
		&ChargeLevelDistribution{
			Min:         31,
			Max:         40,
			LevelString: "31~40级",
		},
		&ChargeLevelDistribution{
			Min:         41,
			Max:         50,
			LevelString: "41~50级",
		},
		&ChargeLevelDistribution{
			Min:         51,
			Max:         60,
			LevelString: "51~60级",
		},
		&ChargeLevelDistribution{
			Min:         61,
			Max:         70,
			LevelString: "61~70级",
		},
		&ChargeLevelDistribution{
			Min:         71,
			Max:         80,
			LevelString: "71~80级",
		},
		&ChargeLevelDistribution{
			Min:         81,
			Max:         90,
			LevelString: "81~90级",
		},
		&ChargeLevelDistribution{
			Min:         91,
			Max:         100,
			LevelString: "91~100级",
		},
		&ChargeLevelDistribution{
			Min:         101,
			Max:         110,
			LevelString: "101~110级",
		},
		&ChargeLevelDistribution{
			Min:         111,
			Max:         120,
			LevelString: "111~120级",
		}, &ChargeLevelDistribution{
			Min:         121,
			Max:         130,
			LevelString: "121~130级",
		},
		&ChargeLevelDistribution{
			Min:         131,
			Max:         140,
			LevelString: "131~140级",
		},
		&ChargeLevelDistribution{
			Min:         141,
			Max:         150,
			LevelString: "141~150级",
		},
		&ChargeLevelDistribution{
			Min:         151,
			Max:         160,
			LevelString: "151~160级",
		},
		&ChargeLevelDistribution{
			Min:         161,
			Max:         170,
			LevelString: "161~170级",
		},
		&ChargeLevelDistribution{
			Min:         171,
			Max:         180,
			LevelString: "171~180级",
		},
		&ChargeLevelDistribution{
			Min:         181,
			Max:         190,
			LevelString: "181~190级",
		},
		&ChargeLevelDistribution{
			Min:         191,
			Max:         200,
			LevelString: "191~200级",
		},
		&ChargeLevelDistribution{
			Min:         201,
			Max:         250,
			LevelString: "201~250级",
		},
		&ChargeLevelDistribution{
			Min:         251,
			Max:         300,
			LevelString: "251~300级",
		},
		&ChargeLevelDistribution{
			Min:         301,
			Max:         400,
			LevelString: "301~400级",
		},
		&ChargeLevelDistribution{
			Min:         401,
			Max:         500,
			LevelString: "401~500级",
		},
		&ChargeLevelDistribution{
			Min:         501,
			Max:         10000,
			LevelString: "大于501级",
		},
	}
	maxCount := 0
	for _, e := range data {
		whereArray := make([]string, 0)
		whereArray = append(whereArray, fmt.Sprintf(" curr_level >= %d", e.Min))
		whereArray = append(whereArray, fmt.Sprintf(" curr_level <= %d", e.Max))
		whereArray = append(whereArray, fmt.Sprintf("charge_type = 99"))
		whereArray = append(whereArray, fmt.Sprintf(" part_id = '%s'", params.PlatformId))
		whereArray = append(whereArray, fmt.Sprintf(" channel in (%s) ", GetSQLWhereParam(params.ChannelList)))
		if params.IsFirst == 1 {
			whereArray = append(whereArray, fmt.Sprintf("is_first = 1"))
		}
		if params.ServerId != "" {
			whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", params.ServerId))
		}
		whereParam := strings.Join(whereArray, " and ")
		if whereParam != "" {
			whereParam = " where " + whereParam
		}
		sql := fmt.Sprintf(`SELECT count(*) as count FROM charge_info_record %s `, whereParam)
		err := DbCharge.Raw(sql).Find(&e).Error
		utils.CheckError(err)
	}
	for _, e := range data {
		maxCount += e.Count
	}
	if maxCount > 0 {
		for _, e := range data {
			e.Rate = float32(e.Count) / float32(maxCount) * 100
		}
	}

	return data
}

//获取该ip节点数量
func GetIpNodeCount(ip string) int {
	l := GetAllServerNodeList()
	count := 0
	for _, e := range l {
		thisIp := strings.Split(e.Node, "@")[1]
		//g.Log().Debug("thisIp:%+v", thisIp)
		if thisIp == ip {
			count++
		}
	}
	return count
}

//获取该ip在线人数
func GetIpOnlinePlayerCount(ip string) int {
	l := GetAllServerNodeList()
	count := 0
	for _, e := range l {
		thisIp := strings.Split(e.Node, "@")[1]
		//g.Log().Debug("thisIp:%+v", thisIp)
		if thisIp == ip {
			gameDb, err := GetGameDbByNode(e.Node)
			utils.CheckError(err)
			if err != nil {
				continue
			}
			defer gameDb.Close()
			count += GetNowOnlineCount(gameDb)
		}
	}
	return count
}

// 向中心服添加game_server
func AddGameServer(PlatformId string, sid string, desc string, node string, zoneNode string, state int, openTime int, isShow int) (string, error) {
	g.Log().Infof("向中心服添加game_server:%s %v %v %v %v %v %v %v ", PlatformId, sid, desc, node, zoneNode, state, openTime, isShow)
	out, err := utils.CenterNodeTool(
		"mod_server_mgr",
		"add_game_server",
		PlatformId,
		sid,
		desc,
		node,
		zoneNode,
		strconv.Itoa(state),
		strconv.Itoa(openTime),
		strconv.Itoa(isShow),
	)
	return out, err
}

func AddServerNode(node string, ip string, port int, webPort int, serverType int, platformId string, dbHost string, dbPort int, dbName string) (string, error) {
	g.Log().Infof("向中心服添加server_node:%v", node, ip, port, webPort, serverType, platformId, dbHost, dbPort, dbName)
	out, err := utils.CenterNodeTool(
		"mod_server_mgr",
		"add_server_node",
		node,
		ip,
		strconv.Itoa(port),
		strconv.Itoa(webPort),
		strconv.Itoa(serverType),
		platformId,
		dbHost,
		strconv.Itoa(dbPort),
		dbName,
	)
	return out, err
}

func InstallNode(node string) error {
	g.Log().Info("开始部署节点:%s......", node)
	var commandArgs []string
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	if err != nil {
		return err
	}
	app := ""
	switch serverNode.Type {
	case 0:
		app = "center"
	case 1:
		app = "game"
	case 2:
		app = "zone"
	case 4:
		app = "login_server"
	case 5:
		app = "unique_id"
	case 6:
		app = "charge"
	case 7:
		app = "war"
	case 8:
		app = "web"
	}

	version := ""

	if serverNode.PlatformId == "" || serverNode.PlatformId == "''" {
		g.Log().Warning("部署节点(%s)没有对应平台， 取默认版本库server!!!!!!!!!!!!!!!!!", node)
		version = "server"
	} else {
		platform, err := GetPlatformOne(serverNode.PlatformId)
		utils.CheckError(err)
		version = platform.Version
	}
	g.Log().Infof("版本库:%s", version)
	toolDir := utils.GetToolDir()
	shName := "do-install.sh"
	dbUserName := g.Cfg().GetString("database.default.user")
	commandArgs = []string{shName, serverNode.Node, app, serverNode.DbName, serverNode.DbHost, strconv.Itoa(serverNode.DbPort), dbUserName, version}
	//out, err := utils.Cmd("sh", commandArgs)
	out, err := utils.CmdShByDirOrParam(toolDir, commandArgs)
	utils.CheckError(err, fmt.Sprintf("部署节点失败:%v %v", node, out))
	if err != nil {
		return err
	}
	g.Log().Infof("部署节点成功:%v!!!", node)
	return nil
}

func NodeAction(nodes []string, action string) error {
	return NodeActionHandle(nodes, action, "server")
}
func NodeActionHandle(nodes []string, action, version string) error {
	g.Log().Infof("节点操作:nodes->%v, action->%v, version->%v", nodes, action, version)
	//curDir := utils.GetCurrentDirectory()
	//defer os.Chdir(curDir)
	toolDir := utils.GetToolDir()
	//err := os.Chdir(toolDir)
	//utils.CheckError(err)
	//if err != nil {
	//	return err
	//}

	if len(version) == 0 {
		version = "server"
	}

	var commandArgs []string
	for _, node := range nodes {
		switch action {
		case "start":
			commandArgs = []string{"node_tool.sh", node, action}
		case "stop":
			commandArgs = []string{"node_tool.sh", node, action}
		case "pull":
			commandArgs = []string{"node_tool.sh", node, action}
		case "hot_reload":
			commandArgs = []string{"node_hot_reload.sh", node, version}
		case "cold_reload":
			commandArgs = []string{"node_cold_reload.sh", node, version}
		}
		out, err := utils.CmdShByDirOrParam(toolDir, commandArgs)
		utils.CheckError(err, fmt.Sprintf("操作节点失败:%v %v", action, out))
		if err != nil {
			return err
		}
	}
	g.Log().Infof("节点操作成功:nodes->%v, action->%v!", nodes, action)
	return nil
}

func AfterAddGameServer() error {
	out, err := utils.CenterNodeTool(
		"mod_server_sync",
		"after_add_game_node",
	)
	utils.CheckError(err, out)
	return err
}

// 更新平台区服状态
func BatchUpdateState(platformId string, EnterState int) error {
	if platformId == "" {
		g.Log().Error("平台id不能为空")
		return gerror.New("平台id不能为空")
	}
	platform, err := GetPlatformOne(platformId)
	if err != nil {
		g.Log().Errorf("更新平台区服状态读取数据失败:%+v", platformId)
		return err
	}
	return BatchPlatformDataUpdateState(platform, EnterState)
}
func BatchPlatformDataUpdateState(platform *Platform, EnterState int) error {
	//if PlatformId == "" {
	//	g.Log().Error("平台id不能为空")
	//	return gerror.New("平台id不能为空")
	//}
	platformId := platform.Id
	out, err := utils.CenterNodeTool(
		"mod_server_mgr",
		"update_all_game_server_state",
		platformId,
		gconv.String(EnterState),
	)
	EnterStateStr := "正常"
	if EnterState == 1 {
		EnterStateStr = "维护"
	} else if EnterState == 3 {
		EnterStateStr = "火爆"
	}
	utils.CheckError(err, fmt.Sprintf("修改%s平台区服状态%v: error:%+v", platformId, EnterState, out))
	if err != nil {
		SendBackgroundMsgTemplateHandle(MSG_TEMPLATE_UPDATE_SERVER_STATE, platformId, platformId+"更新平台入口状态失败", "游戏服入口状态:"+EnterStateStr+" 操作时间:"+gtime.Datetime())
		return err
	}
	platform.EnterState = EnterState
	err = Db.Save(platform).Error
	if err == nil {
		SendBackgroundMsgTemplateHandle(MSG_TEMPLATE_UPDATE_SERVER_STATE, platformId, platformId+"更新平台入口状态成功", "游戏服入口状态:"+EnterStateStr+" 操作时间:"+gtime.Datetime())
	}
	//UpdateEtsPlatformServerClose(PlatformId , EnterState == 1)
	return err
}

func RefreshGameServer() error {
	out, err := utils.CenterNodeTool(
		"mod_server_sync",
		"push_all_login_server_node",
	)
	utils.CheckError(err, out)
	//if err == nil {
	//	RefreshEtsPlatformGameEnter("")
	//}
	return err
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

//func S() {
//	g.Log().Info("开始统计")
//	gameServerList, _ := GetAllGameServer()
//	context := ""
//	for _, e := range gameServerList {
//		var data [] struct {
//			PlayerId int
//		}
//		gameDb, err := GetGameDbByNode(e.Node)
//		utils.CheckError(err)
//		if err != nil {
//			return
//		}
//		defer gameDb.Close()
//		sql := fmt.Sprintf(
//			`select a.player_id as player_id,b.acc_id from player_platform_award AS a,player AS b where a.player_id = b.id and a.id = 1601`)
//		err = gameDb.Raw(sql).Find(&data).Error
//		utils.CheckError(err)
//		for _, i := range data {
//			player, err := GetPlayerOne(e.PlatformId, e.Sid, i.PlayerId)
//			utils.CheckError(err)
//			is2 := IsThatDayPlayerLogin(gameDb, utils.GetThatZeroTimestamp(int64(player.RegTime))+86400, player.Id)
//			is3 := IsThatDayPlayerLogin(gameDb, utils.GetThatZeroTimestamp(int64(player.RegTime))+86400*2, player.Id)
//			var moneyList [] struct {
//				Money        int
//				ChargeItemId int
//			}
//			sql := fmt.Sprintf(
//				`select money , charge_item_id from charge_info_record where player_id = %d;`, i.PlayerId)
//			err = DbCharge.Raw(sql).Find(&moneyList).Error
//			m := make([] string, 0)
//			//g.Log().Info("moneyList:%+v", moneyList)
//			for _, e := range moneyList {
//				m = append(m, strconv.Itoa(e.Money))
//			}
//
//			m1 := make([] string, 0)
//			//g.Log().Info("moneyList:%+v", moneyList)
//			for _, e := range moneyList {
//				m1 = append(m1, strconv.Itoa(e.ChargeItemId))
//			}
//			context += fmt.Sprintf("%s, %d, %s, %t, %t, [%s], [%s]\n", e.Sid, player.Id, player.AccId, is2, is3, strings.Join(m, " "), strings.Join(m1, " "))
//		}
//	}
//	utils.FilePutContext("data.txt", context)
//	g.Log().Info("统计完毕")
//}

// 获取服务器整形数据
func GetServerDataInt(db *gorm.DB, serverDataId int) int {
	var data struct {
		Data int
	}
	sql := fmt.Sprintf(
		`SELECT int_data as data FROM server_data WHERE id =  %d`, serverDataId)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Data
}

// 获取服务器字符串数据
func GetServerDataStr(db *gorm.DB, serverDataId int) string {
	var data struct {
		Data string
	}
	sql := fmt.Sprintf(
		`SELECT str_data as data FROM server_data WHERE id =  %d`, serverDataId)
	err := db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Data
}

func RRRR() {
	g.Log().Info("go")
	accountList := make([]string, 0)
	var data []struct {
		Id int
	}
	content := ""
	sql := fmt.Sprintf(
		`SELECT player_id as id FROM player_charge_info_record WHERE total_money  >=  200 and part_id = 'qq'`)
	err := DbCharge.Raw(sql).Find(&data).Error
	utils.CheckError(err)
	if err != nil {
		return
	}
	for _, e := range data {
		globalPlayer, err := GetGlobalPlayerOne(e.Id)
		utils.CheckError(err)
		if err != nil {
			return
		}
		if IsHaveArray(globalPlayer.Account, accountList) == false {
			accountList = append(accountList, globalPlayer.Account)
		}

	}
	for _, e := range accountList {

		content += fmt.Sprintf("%s\n", strings.Replace(strings.Replace(e, "ios_", "", -1), "android_", "", -1))
		//content += "nodes="
		//content += "\"[" + strings.Join(nodes, ", ") + "]\""
		//content += "\n\n"

	}
	err = utils.FilePutContext("account.txt", content)
	utils.CheckError(err)
	g.Log().Info("go finish")
}

// 是否存在列表中
func IsHaveArray(v string, array []string) bool {
	for _, e := range array {
		if e == v {
			return true
		}
	}
	return false
}

func SSS() {
	g.Log().Info("go")
	var data []struct {
		Id string
	}

	//content := ""
	sql := fmt.Sprintf(
		`select distinct account as id from  global_player`)
	err := DbCenter.Raw(sql).Find(&data).Error
	utils.CheckError(err)
	g.Log().Info("go1")
	if err != nil {
		return
	}
	args := make([]string, 0, len(data))
	for _, e := range data {
		args = append(args, e.Id)

	}
	g.Log().Info("go2")
	content := strings.Join(args, "\n")
	g.Log().Info("go3")
	err = utils.FilePutContext("all_account.txt", content)
	g.Log().Info("go4")
	utils.CheckError(err)
	g.Log().Info("go finish")
}

func QQQQ() {
	g.Log().Info("go")
	args := make([]string, 0, 100000)
	sql := fmt.Sprintf(
		`select acc_id from player , player_data where player.id = player_data.player_id and player_data.vip_level = 0 and player_data.level > 50 limit 1000;`)
	for i := 1; i <= 300; i++ {
		var data []struct {
			AccId string
		}
		gameDb, err := GetGameDbByPlatformIdAndSid("wx", fmt.Sprintf("s%d", i))
		utils.CheckError(err)
		if err != nil {
			return
		}
		defer gameDb.Close()
		err = gameDb.Raw(sql).Find(&data).Error
		utils.CheckError(err)
		for _, e := range data {
			args = append(args, e.AccId)
			if len(args) == 100000 {
				break
			}
		}
		if len(args) == 100000 {
			break
		}

	}
	g.Log().Info("go2")
	content := strings.Join(args, "\n")
	g.Log().Info("go3")
	err := utils.FilePutContext("account_20181115.txt", content)
	g.Log().Info("go4")
	utils.CheckError(err)
	g.Log().Info("go finish")
}

//
//func GetALlServerNode(platformId string) [] string{
//	GetALlServerNode()
//}

func PowerSearch(platformId string) {
	//todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	serverNodeList := GetAllGameServerNodeByPlatformId(platformId)
	//gameServerList, serverCount := GetAllGameServerDirtyByPlatfomrId(platformId)
	g.Log().Infof("PowerSearch:%v", len(serverNodeList))
	r := make(map[int]int)
	for _, serverNode := range serverNodeList {
		gameDb, err := GetGameDbByNode(serverNode.Node)
		utils.CheckError(err)
		if err != nil {
			return
		}
		defer gameDb.Close()

		for i := 10000000; i < 100000000; i += 10000000 {
			var data struct {
				Count int
			}
			sql := fmt.Sprintf(
				`select count(player_id) as count from player left join player_data on player.id = player_data.player_id and player_data.power between %d and %d and player.last_login_time > 1543852800;`, i, i+10000000-1)
			err = gameDb.Raw(sql).Scan(&data).Error
			utils.CheckError(err)
			if v, ok := r[i]; ok {
				r[i] = data.Count + v
			} else {
				r[i] = data.Count
			}
		}
		var data struct {
			Count int
		}
		sql := fmt.Sprintf(
			`select count(player_id) as count from player left join player_data on player.id = player_data.player_id and player_data.power >= 100000000 and player.last_login_time > 1543852800;`)
		err = gameDb.Raw(sql).Scan(&data).Error
		if v, ok := r[100000000]; ok {
			r[100000000] = data.Count + v
		} else {
			r[100000000] = data.Count
		}

	}
	g.Log().Info("PowerSearch完毕%+v.", r)
}
