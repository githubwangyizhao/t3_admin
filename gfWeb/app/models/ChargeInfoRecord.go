package models

import (
	"fmt"
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"

	//"os/user"
	"strings"
)

type ChargeTotalMoney struct {
	Count  int `json:"count"`
	ItemId int `json:"item_id"`
}

type ChargeInfoRecord struct {
	OrderId         string  `json:"orderId" gorm:"primary_key"`
	PlatformOrderId string  `json:"platform_order_id"`
	ChargeType      int     `json:"chargeType"`
	Ip              string  `json:"ip"`
	PartId          string  `json:"platformId"`
	ServerId        string  `json:"serverId"`
	AccId           string  `json:"accId"`
	IsFirst         int     `json:"isFirst"`
	CurrLevel       int     `json:"currLevel"`
	CurrTaskId      int     `json:"currTaskId"`
	RegTime         int     `json:"regTime"`
	FirstTime       int     `json:"firstTime"`
	CurrPower       int     `json:"currPower"`
	PlayerId        int     `json:"playerId"`
	PlayerName      string  `json:"playerName" gorm:"-"`
	LastLoginTime   int     `json:"lastLoginTime" gorm:"-"`
	Money           float32 `json:"money"`
	Ingot           int     `json:"ingot"`
	RecordTime      int     `json:"recordTime"`
	ChargeItemId    int     `json:"chargeItemId"`
	Channel         string  `json:"channel"`
	Gold            int     `json:"gold"`
	Bounty          int     `json:"bounty"`
	Coupon          int     `json:"coupon"`

	//
	SAccount string `json:"server_name" gorm:"_"`
}

type ChargeInfoRecordQueryParam struct {
	BaseQueryParam
	PlatformId      string
	PlatFormOrderId string   `json:"platformOrderId"`
	ServerId        string   `json:"serverId"`
	ChannelList     []string `json:"channelList"`
	PlayerId        int
	PlayerName      string
	OrderId         string
	AccId           string
	StartTime       int
	EndTime         int
	Promote         string   `json:"promote"`
	SAccountList    []string `json:"s_account"`
}

func GetChargeInfoBySpecifyDate(platformId string, serverId string, channel []string, startTime int, endTime int) (int, int, int, float32, float32, float32) {
	var (
		firstChargePlayers    int
		multiChargePlayers    int
		firstChargeTotalMoney float32
		multiChargeTotalMoney float32
		totalChargePlayers    int
		totalMoney            float32
		//selectSql = " count(DISTINCT player_id) AS players, sum(money) AS total_money, is_first AS is_first"
		selectSql = " player_id AS players, sum(money) AS total_money, is_first AS is_first"
		condSql   = ""
		condArr   = make([]string, 0)
		groupSql  = " GROUP BY is_first, player_id"
	)

	if platformId != "" {
		condArr = append(condArr, fmt.Sprintf(`part_id = '%s'`, platformId))
	}
	if serverId != "" {
		condArr = append(condArr, fmt.Sprintf(`server_id = '%s'`, serverId))
	}

	if len(channel) > 0 {
		ChannelStr := GetSQLWhereParam(channel)
		if ChannelStr != "'all'" && ChannelStr != "" {
			var realChannel = ChannelStr
			if n := strings.Index(ChannelStr, "all"); n != -1 {
				realChannel = strings.Replace(ChannelStr, "'all',", "", -1)
				realChannel = strings.Replace(realChannel, ",'all'", "", -1)
			}
			condArr = append(condArr, fmt.Sprintf(`channel in (%s)`, ChannelStr))
		}
	}

	if startTime != 0 && endTime != 0 && endTime > startTime {
		condArr = append(condArr, fmt.Sprintf(`record_time between %d and %d`, startTime, endTime))
	}
	if len(condArr) > 0 {
		condSql = " WHERE " + strings.Join(condArr, " AND ")
	}

	type TmpChargeStruct struct {
		Players    int
		TotalMoney float32
		IsFirst    int
	}
	data := make([]*TmpChargeStruct, 0)

	sql := fmt.Sprintf("SELECT %s FROM charge_info_record %s %s", selectSql, condSql, groupSql)

	err := DbCharge.Debug().Raw(sql).Scan(&data).Error
	utils.CheckError(err)

	matchArr := make(map[int]bool)
	//exists := make(map[int]bool)
	for _, item := range data {
		if item.IsFirst == 1 {
			firstChargePlayers += 1
			firstChargeTotalMoney += item.TotalMoney
		} else {
			multiChargePlayers += 1
			multiChargeTotalMoney += item.TotalMoney
		}
		if _, ok := matchArr[item.Players]; !ok {
			matchArr[item.Players] = true
			totalChargePlayers += 1
		}
		totalMoney += item.TotalMoney
	}

	return firstChargePlayers, multiChargePlayers, totalChargePlayers, firstChargeTotalMoney, multiChargeTotalMoney, totalMoney
}

func GetChargeInfoRecordList(params *ChargeInfoRecordQueryParam) ([]*ChargeInfoRecord, int64, int64, int64, []*ChargeTotalMoney, float32) {
	data := make([]*ChargeInfoRecord, 0)
	var count int64
	sortOrder := "record_time"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	} else if params.Order == "ascending" {
		sortOrder = sortOrder + " asc"
	} else {
		sortOrder = sortOrder + " desc"
	}
	var sumData struct {
		MoneyCount  float32
		PlayerCount int64
	}
	whereArray := make([]string, 0)
	whereArray = append(whereArray, " charge_type = 99 ")
	whereArray = append(whereArray, " part_id =  '"+params.PlatformId+"' ")
	if params.ServerId != "" {
		whereArray = append(whereArray, " server_id =  '"+params.ServerId+"' ")
		//whereArray = append(whereArray, fmt.Sprintf(" server_id in (%s) ", GetGameServerIdListStringByNode(params.Node)))
	}

	if len(params.ChannelList) > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" channel in (%s) ", GetSQLWhereParam(params.ChannelList)))
	}

	AccIdList := make([]string, 0)
	if params.Promote != "" {
		AccIdList, err := GetGlobalAccountByPromote(params.Promote)
		utils.CheckError(err)
		g.Log().Debug("AccIdList:", params.Promote, AccIdList)
	}
	if len(AccIdList) > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" acc_id in (%s) ", GetSQLWhereParam(AccIdList)))
	}
	if params.StartTime > 0 {
		whereArray = append(whereArray, fmt.Sprintf("record_time between %d and %d", params.StartTime, params.EndTime))
	}
	if len(params.SAccountList) > 0 {
		GlobalAccountId, _ := GetGlobalAccountIdBySAccount(params.SAccountList)
		whereArray = append(whereArray, fmt.Sprintf(" player_id IN (%s)", strings.Join(*GlobalAccountId, ",")))
	}

	whereParam := " where " + strings.Join(whereArray, " and ")

	sql := fmt.Sprintf(
		`select sum(money) as money_count, count(DISTINCT player_id) as player_count  from charge_info_record  %s;`, whereParam)
	g.Log().Debug("sql:", sql)
	err := DbCharge.Raw(sql).Scan(&sumData).Error
	utils.CheckError(err)
	if params.Promote != "" {
		if params.StartTime > 0 {
			ChargeModel := DbCharge.Model(&ChargeInfoRecord{}).Where(&ChargeInfoRecord{
				PartId:     params.PlatformId,
				AccId:      params.AccId,
				OrderId:    params.OrderId,
				PlayerId:   params.PlayerId,
				ServerId:   params.ServerId,
				ChargeType: 99,
			}).Where("record_time between ? and ? ", params.StartTime, params.EndTime).
				Where("channel in(?)", params.ChannelList).
				Offset(params.Offset).
				Limit(params.Limit).
				Order(sortOrder).
				Find(&data).
				Count(&count)

			if len(AccIdList) > 0 {
				ChargeModel.Where("acc_id in (?)", AccIdList)
			}
			err = ChargeModel.Error
			utils.CheckError(err)
		} else {
			ChargeModel := DbCharge.Model(&ChargeInfoRecord{}).Where(&ChargeInfoRecord{
				PartId:     params.PlatformId,
				AccId:      params.AccId,
				OrderId:    params.OrderId,
				PlayerId:   params.PlayerId,
				ServerId:   params.ServerId,
				ChargeType: 99,
			}).Where("channel in(?)", params.ChannelList).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count)
			if len(AccIdList) > 0 {
				ChargeModel.Where("acc_id in(?)", AccIdList)
			}
			err = ChargeModel.Error
			utils.CheckError(err)
		}
	} else {
		err = DbCharge.Model(&ChargeInfoRecord{}).Where(&ChargeInfoRecord{
			PartId:     params.PlatformId,
			AccId:      params.AccId,
			OrderId:    params.OrderId,
			PlayerId:   params.PlayerId,
			ServerId:   params.ServerId,
			ChargeType: 99,
		}).Where("channel in(?)", params.ChannelList).
			Where(strings.Join(whereArray, " and ")).
			Count(&count).Offset(params.Offset).
			Limit(params.Limit).Order(sortOrder).Find(&data).Error
		utils.CheckError(err)
	}

	for _, e := range data {
		e.PlayerName = GetPlayerName_2(e.PartId, e.ServerId, e.PlayerId)
		e.LastLoginTime = GetPlayerLastLoginTime(e.PartId, e.ServerId, e.PlayerId)
		//e.ChargeItemId = GetChargeItemId(e.OrderId, e.PartId, e.ServerId)
	}

	sql = fmt.Sprintf(`select count(*) as count, charge_item_id as item_id from charge_info_record %s group by charge_item_id`, whereParam)
	totalData := make([]*ChargeTotalMoney, 0)
	err = DbCharge.Raw(sql).Scan(&totalData).Error
	g.Log().Debug("eee:", sql)
	//g.Log().Debug("fff: ", totalData)

	// 转cny
	//exchangeRate := GetExchangeRate("indonesia")
	var exchangeRate float32
	exchangeRate = 1.0

	return data, count, sumData.PlayerCount, int64(sumData.MoneyCount), totalData, exchangeRate
}

//func Repair() {
//	g.Log().Info("开始修复充值数据")
//	data := make([]*ChargeInfoRecord, 0)
//	//var count int64
//	err := DbCharge.Model(&ChargeInfoRecord{}).Where(&ChargeInfoRecord{
//		ChargeType: 99,
//	}).Find(&data).Error
//	utils.CheckError(err)
//	if err != nil {
//		return
//	}
//	for _, e := range data {
//		value := GetChargeItemId(e.OrderId, e.PartId, e.ServerId)
//		g.Log().Debug("value:%v, %v", e.OrderId, value)
//		if value > 0 {
//			err = DbCharge.Model(&e).Update("charge_item_id", value).Error
//			utils.CheckError(err)
//			if err != nil {
//				return
//			}
//		}
//	}
//	g.Log().Info("修复充值数据成功")
//}
//
////获取玩家最近登录时间
//func GetChargeItemId(orderId string, platformId string, serverId string) int {
//	gameDb, err := GetGameDbByPlatformIdAndSid(platformId, serverId)
//	utils.CheckError(err)
//	if err != nil {
//		return 0
//	}
//	defer gameDb.Close()
//	var data struct {
//		ChargeItemId int
//	}
//	sql := fmt.Sprintf(
//		`SELECT charge_item_id FROM player_charge_record where order_id = %d `, orderId)
//	err = gameDb.Raw(sql).Scan(&data).Error
//	utils.CheckError(err)
//	return data.ChargeItemId
//}

type ChargeStatistics struct {
	PlatformId string `json:"platformId"`
	//ServerId                    string    `json:"serverId"`
	TodayCharge                      int `json:"todayCharge"`
	TodayChargePlayerCount           int `json:"todayChargePlayerCount"`
	YesterdayCharge                  int `json:"yesterdayCharge"`
	YesterdayChargePlayerCount       int `json:"yesterdayChargePlayerCount"`
	BeforeYesterdayCharge            int `json:"beforeYesterdayCharge"`
	BeforeYesterdayChargePlayerCount int `json:"beforeYesterdayChargePlayerCount"`

	ChargeData            []map[string]string `json:"chargeData"`
	ChargePlayerCountData []map[string]string `json:"chargePlayerCountData"`

	TodayChargeList           []string `json:"todayChargeList"`
	YesterdayChargeList       []string `json:"yesterdayChargeList"`
	BeforeYesterdayChargeList []string `json:"beforeYesterdayChargeList"`

	TodayChargePlayerCountList           []string `json:"todayChargePlayerCountList"`
	YesterdayChargePlayerCountList       []string `json:"yesterdayChargePlayerCountList"`
	BeforeYesterdayChargePlayerCountList []string `json:"beforeYesterdayChargePlayerCountList"`
}

func GetChargeStatistics(platformId string, serverId string, channelList []string) (*ChargeStatistics, error) {

	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	yesterdayZeroTimestamp := todayZeroTimestamp - 86400
	beforeYesterdayZeroTimestamp := yesterdayZeroTimestamp - 86400

	todayOnlineList, todayTotalCharge := get24hoursChargeCount(platformId, serverId, channelList, todayZeroTimestamp)
	yesterdayOnlineList, yesterdayTotalCharge := get24hoursChargeCount(platformId, serverId, channelList, yesterdayZeroTimestamp)
	beforeYesterdayOnlineList, beforeYesterdayTotalCharge := get24hoursChargeCount(platformId, serverId, channelList, beforeYesterdayZeroTimestamp)

	todayChargePlayerCountList, todayChargePlayerCount := get24hoursChargePlayerCount(platformId, serverId, channelList, todayZeroTimestamp)
	yesterdayChargePlayerCountList, yesterdayChargePlayerCount := get24hoursChargePlayerCount(platformId, serverId, channelList, yesterdayZeroTimestamp)
	beforeYesterdayChargePlayerCountList, beforeYesterdayChargePlayerCount := get24hoursChargePlayerCount(platformId, serverId, channelList, beforeYesterdayZeroTimestamp)

	chargeData := make([]map[string]string, 0, 144)
	//g.Log().Info("len:%d", len(todayOnlineList))
	for i := 0; i < 6*24; i = i + 1 {
		m := make(map[string]string, 4)
		m["时间"] = utils.FormatTime(i * 10 * 60)
		m["今日充值"] = todayOnlineList[i]
		m["昨日充值"] = yesterdayOnlineList[i]
		m["前日充值"] = beforeYesterdayOnlineList[i]
		//g.Log().Info(i)
		chargeData = append(chargeData, m)
	}

	chargePlayerCountData := make([]map[string]string, 0, 144)
	for i := 0; i < 6*24; i = i + 1 {
		m := make(map[string]string, 4)
		m["时间"] = utils.FormatTime(i * 10 * 60)
		m["今日充值人数"] = todayChargePlayerCountList[i]
		m["昨日充值人数"] = yesterdayChargePlayerCountList[i]
		m["前日充值人数"] = beforeYesterdayChargePlayerCountList[i]
		chargePlayerCountData = append(chargePlayerCountData, m)
	}

	chargeStatistics := &ChargeStatistics{
		PlatformId:             platformId,
		TodayCharge:            todayTotalCharge,
		TodayChargePlayerCount: todayChargePlayerCount,
		//TodayCreateRole: GetCreateRoleCountByChannelList(gameDb, serverId, channelList, todayZeroTimestamp, todayZeroTimestamp+86400),
		YesterdayCharge:            yesterdayTotalCharge,
		YesterdayChargePlayerCount: yesterdayChargePlayerCount,
		//MaxOnlineCount:              GetMaxOnlineCount(node),
		BeforeYesterdayCharge:            beforeYesterdayTotalCharge,
		BeforeYesterdayChargePlayerCount: beforeYesterdayChargePlayerCount,

		TodayChargeList:           todayOnlineList,
		YesterdayChargeList:       yesterdayOnlineList,
		BeforeYesterdayChargeList: beforeYesterdayOnlineList,

		TodayChargePlayerCountList:           todayChargePlayerCountList,
		YesterdayChargePlayerCountList:       yesterdayChargePlayerCountList,
		BeforeYesterdayChargePlayerCountList: beforeYesterdayChargePlayerCountList,
		ChargeData:                           chargeData,
		ChargePlayerCountData:                chargePlayerCountData,
	}
	return chargeStatistics, nil
}

// 下载当前充值数据
func GetChargeInfoDownload(params *ChargeInfoRecordQueryParam) []*ChargeInfoRecord {
	data := make([]*ChargeInfoRecord, 0)
	sortOrder := "record_time"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	} else if params.Order == "ascending" {
		sortOrder = sortOrder + " asc"
	} else {
		sortOrder = sortOrder + " desc"
	}
	if params.StartTime > 0 {
		err := DbCharge.Model(&ChargeInfoRecord{}).Where(&ChargeInfoRecord{
			PartId:     params.PlatformId,
			AccId:      params.AccId,
			PlayerId:   params.PlayerId,
			ServerId:   params.ServerId,
			ChargeType: 99,
		}).Where("record_time between ? and ? ", params.StartTime, params.EndTime).Where("channel in(?)", params.ChannelList).Order(sortOrder).Find(&data).Error
		utils.CheckError(err)
	} else {
		err := DbCharge.Model(&ChargeInfoRecord{}).Where(&ChargeInfoRecord{
			PartId:     params.PlatformId,
			AccId:      params.AccId,
			PlayerId:   params.PlayerId,
			ServerId:   params.ServerId,
			ChargeType: 99,
		}).Where("channel in(?)", params.ChannelList).Order(sortOrder).Find(&data).Error
		utils.CheckError(err)
	}
	for _, e := range data {
		e.PlayerName = GetPlayerName_2(e.PartId, e.ServerId, e.PlayerId)
	}
	return data
}

type RechargeGenInfo struct {
	FirstNums int     `gorm:"first_num"` // 首充人数 = 首充笔数
	Amount    float32 `gorm:"amount"`    // 总金额
	ReNums    int     //复充人数
	ReAcount  int     //复充笔数
	RecTotal  int     //充值笔数
}

// 按时间计算 总充值金额- 充值笔数-首充人数
func GetRechargeGeneralize(platformId, serverId string, selectZeroTime, selectEndTime int) (RechargeGenInfo, []int) {

	whereArray := make([]string, 0)
	whereArray = append(whereArray, fmt.Sprintf("part_id='%s'", platformId))
	whereArray = append(whereArray, fmt.Sprintf("server_id='%s'", serverId))
	whereArray = append(whereArray, "status=1")
	whereArray = append(whereArray, "record_time >= "+gconv.String(selectZeroTime))
	whereArray = append(whereArray, "record_time <= "+gconv.String(selectEndTime))
	whereArray = append(whereArray, "charge_type=99")

	// fields := "@num:=0 _,if(is_first=1, @num+1, @num) first_nums, count(1) count, sum(money) amount"
	// sql := fmt.Sprintf("select %s from charge_info_record where %s", fields, strings.Join(whereArray, " AND "))

	var playerRes = make([]int, 0)

	// sql1 := fmt.Sprintf("select player_id from charge_info_record where %s group by player_id", strings.Join(whereArray, " AND "))
	// err = DbCharge.Raw(sql1).Scan(&playerRes).Error
	// utils.CheckError(err)

	sql := fmt.Sprintf("select * from charge_info_record where %s", strings.Join(whereArray, " AND "))
	var res = RechargeGenInfo{}
	var data = make([]*ChargeInfoRecord, 0)
	err := DbCharge.Debug().Raw(sql).Scan(&data).Error
	utils.CheckError(err)

	var (
		mRePlayerIds = make(map[int]string) //复值总人数
		mPlayerIds   = make(map[int]string) //充值总人数
	)
	for _, item := range data {
		if item.IsFirst == 1 {
			res.FirstNums++
		} else {
			mRePlayerIds[item.PlayerId] = ""
			res.ReAcount++
		}
		res.Amount += item.Money
		mPlayerIds[item.PlayerId] = ""
	}

	res.ReNums = len(mRePlayerIds)
	res.RecTotal = len(data)

	for k, _ := range mPlayerIds {
		playerRes = append(playerRes, k)
	}

	return res, playerRes
}
