package models

import (
	"database/sql"
	"fmt"
	"gfWeb/library/utils"

	"github.com/gogf/gf/os/gtime"
	"github.com/jinzhu/gorm"
)

type Account struct {
	AccId             string `json:"accId"`
	ServerId          string `json:"serverId"`
	Channel           string `json:"channel"`
	IsCreateRole      int    `json:"isCreateRole"`
	IsEnterGame       int    `json:"isEnterGame"`
	PlayerId          int    `json:"playerId"`
	IsFinishFirstTask int    `json:"isFinishFirstTask"`
	Time              int    `json:"time"`
}

////获取总注册人数
//func GetTotalRegister(db *gorm.DB) int {
//	var count int
//	db.Model(&Account{}).Count(&count)
//	return count
//}

//获取今日注册人数
//func GetTodayRegister(db *gorm.DB) int {
//	return GetThatDayRegister(db, utils.GetTodayZeroTimestamp())
//}

////获取某日注册人数
//func GetThatDayRegister(db *gorm.DB, zeroTimestamp int) int {
//	var data struct {
//		Count int
//	}
//	//DbCenter.Model(&CServerTraceLog{}).Where(&CServerTraceLog{Node:gameServer.Node}).Count(&count)
//	//todayZeroTimestamp := utils.GetTodayZeroTimestamp()
//	sql := fmt.Sprintf(
//		`SELECT count(1) as count FROM account WHERE time between %d and %d`, zeroTimestamp, zeroTimestamp+86400)
//	err := db.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	//g.Log().Info("ppp:%v,%v", gameServer.Node, data.Count)
//	return data.Count
//}

//获取总注册人数
func GetTotalRegisterRoleCountByChannelList(db *gorm.DB, serverId string, channelList []string, endTime int) int {
	var count int
	db.Model(&Account{}).Where(&Account{
		ServerId: serverId,
	}).Where(" channel in(?)", channelList).Where("time < ?", endTime).Count(&count)
	return count
}

//获取时间内注册人数
func GetRegisterRoleCount(db *gorm.DB, serverId string, channel string, startTime int, endTime int) int {
	var count int
	db.Model(&Account{}).Where(&Account{
		ServerId: serverId,
		Channel:  channel,
	}).Where(" time between ? and ? ", startTime, endTime).Count(&count)
	return count
}

//获取时间内注册人数
func GetRegisterRoleCountByChannelList(db *gorm.DB, serverId string, channelList []string, startTime int, endTime int) int {
	var count int
	db.Model(&Account{}).Where(&Account{
		ServerId: serverId,
	}).Where(" time between ? and ? ", startTime, endTime).Where(" channel in(?)", channelList).Count(&count)
	return count
}

////获取某日注册人数
//func GetThatDayRegisterByChannel(db *gorm.DB, serverId string, channel string, zeroTimestamp int) int {
//	var count int
//	db.Model(&Account{}).Where(&Account{
//		ServerId:serverId,
//		Channel:channel,
//	}).Where(" time between ? and ? ", zeroTimestamp, zeroTimestamp + 86400).Count(&count)
//	return count
//}

func GetTotalCreateRoleCountByNode(node string) int {
	gameDb, err := GetGameDbByNode(node)

	utils.CheckError(err)
	if err != nil {
		return -1
	}
	defer gameDb.Close()
	return GetTotalCreateRoleCount(gameDb)
}

//获取总创角色人数
func GetTotalCreateRoleCount(db *gorm.DB) int {
	var count int
	db.Model(&Account{}).Where(&Account{
		IsCreateRole: 1,
	}).Count(&count)
	return count
}

//获取总创角人数
func GetTotalCreateRoleCountByChannelList(db *gorm.DB, serverId string, channelList []string, endTime int) int {
	if endTime == 0 {
		endTime = int(gtime.Now().Timestamp())
	}
	var count int
	db.Model(&Account{}).Where(&Account{
		ServerId:     serverId,
		IsCreateRole: 1,
	}).Where(" channel in(?)", channelList).Where("time < ?", endTime).Count(&count)
	return count
}

//获取总创角色人数
//func GetTotalCreateRoleCountByChannel(db *gorm.DB, serverId string, channel string, time int) int {
//	var count int
//	db.Model(&Account{}).Where(&Account{
//		ServerId:     serverId,
//		Channel:channel,
//		IsCreateRole: 1,
//	}).Where("time < ? ", time).Count(&count)
//	return count
//}
func GetCreateRoleCount(db *gorm.DB, serverId string, channel string, startTime int, endTime int) int {
	var count int
	db.Model(&Account{}).Where(&Account{
		IsCreateRole: 1,
		ServerId:     serverId,
		Channel:      channel,
	}).Where(" time between ? and ? ", startTime, endTime).Count(&count)
	return count
}

func GetCreateRoleCountByChannelList(db *gorm.DB, serverId string, channelList []string, startTime int, endTime int) int {
	var count int
	db.Model(&Account{}).Where(&Account{
		IsCreateRole: 1,
		ServerId:     serverId,
	}).Where(" time between ? and ? ", startTime, endTime).Where("channel in(?)", channelList).Count(&count)
	return count
}

////获取今日创角人数
//func GetTodayCreateRoleCount(db *gorm.DB) int {
//	return GetThatDayCreateRoleCount(db, utils.GetTodayZeroTimestamp())
//}

//获取某日创角人数
//func GetThatDayCreateRoleCount(db *gorm.DB, zeroTimestamp int) int {
//	var data struct {
//		Count int
//	}
//	sql := fmt.Sprintf(
//		`SELECT count(1) as count FROM account WHERE is_create_role = 1 and (time between %d and %d) `, zeroTimestamp, zeroTimestamp+86400)
//	err := db.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.Count
//}

//获取该天分享创角
func GetThatDayShareCreateRoleCountByChannel(db *gorm.DB, serverId string, channel string, zeroTimestamp int) int {
	var count int
	db.Model(&Player{}).Where(&Player{
		ServerId: serverId,
		Channel:  channel,
	}).Where(" reg_time between ? and ? ", zeroTimestamp, zeroTimestamp+86400).Where(" friend_code != '' ").Count(&count)
	return count
	//var data struct {
	//	Count int
	//}
	//sql := fmt.Sprintf(
	//	`SELECT count(1) as count FROM account WHERE is_create_role = 1 and (time between %d and %d) and server_id = '%s'`, zeroTimestamp, zeroTimestamp+86400, serverId)
	//err := db.Raw(sql).Scan(&data).Error
	//utils.CheckError(err)
	//return data.Count
}

//获取当前全服资源总数
func GetTotalProp(db *gorm.DB, propId int, channelList []string) int {
	var data struct {
		Count int
	}
	sqlStatement := fmt.Sprintf(
		`SELECT sum(num) as count FROM player_prop WHERE prop_id = '%d' and player_id in(select id from player where channel in(%s))`, propId, GetSQLWhereParam(channelList))
	err := db.Raw(sqlStatement).Scan(&data).Error
	if err != sql.ErrNoRows {
		utils.CheckError(err)
	}

	return data.Count
}

//获取某日有效角色人数
func GetThatDayValidCreateRoleCountByChannel(db *gorm.DB, serverId string, channel string, zeroTimestamp int) int {
	var count int
	db.Model(&Account{}).Where(&Account{
		IsCreateRole:      1,
		IsFinishFirstTask: 1,
		ServerId:          serverId,
		Channel:           channel,
	}).Where("time between ? and ?", zeroTimestamp, zeroTimestamp+86400).Count(&count)
	return count
}
