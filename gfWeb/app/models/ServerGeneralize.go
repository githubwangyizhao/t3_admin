package models

import (
	"gfWeb/library/utils"
	"time"
)

type ServerGeneralize struct {
	PlatformId              string  `json:"platformId"`
	ServerId                string  `json:"serverId"`
	OpenTime                int     `json:"openTime"`
	MergeTime               int     `json:"mergeTime"`
	Version                 int     `json:"version"`
	TotalRegister           int     `json:"totalRegister"`
	TotalCreateRole         int     `json:"totalCreateRole"`
	TodayCreateRole         int     `json:"todayCreateRole"`
	TodayRegister           int     `json:"todayRegister"`
	TotalChargeIngot        int     `json:"totalChargeIngot"`
	TotalChargeMoney        float32 `json:"totalChargeMoney"`
	TotalChargePlayerCount  int     `json:"totalChargePlayerCount"`
	SecondChargePlayerCount int     `json:"secondChargePlayerCount"`
	ChargeCount2Rate        float32 `json:"chargeCount2Rate"`
	ChargeCount3Rate        float32 `json:"chargeCount3Rate"`
	ChargeCount5Rate        float32 `json:"chargeCount5Rate"`
	//ChargeCountMore         int     `json:"chargeCountMore"`
	OnlineCount             int     `json:"onlineCount"`
	Status                  int     `json:"status"`
	ARPU                    float32 `json:"arpu"`
	ChargeRate              float32 `json:"chargeRate"`
	SecondChargeRate        float32 `json:"secondChargeRate"`
	MaxLevel                int     `json:"maxLevel"`
	MaxOnlineCount          int     `json:"maxOnlineCount"`
	YesterdayMaxOnlineCount int     `json:"yesterdayMaxOnlineCount"`
	TotalIngot              int     `json:"totalIngot"`
	TotalCoin               int     `json:"totalCoin"`
	TotalBounty             int     `json:"totalBounty"`

	TodayLoginTimes         int     `json:"todayLoginTimes"`         // 今日玩家登录次数
	FirstRechargeTimes      int     `json:"firstRechargeTimes"`      // 首充人数
	FirstRechargeCount      int     `json:"firstRechargeCount"`      // 首充笔数
	ReRechargeTimes         int     `json:"reRechargeTimes"`         // 复充人数
	ReRechargeCount         int     `json:"reRechargeCount"`         // 复充笔数
	WdlAmount               float32 `json:"withdrawalAmount"`        //提现金额
	WdlTimes                int     `json:"withdrawalTimes"`         //提现笔数
	WdlPlayerCount          int     `json:"wdlPlayerCount"`          //提现人数
	IoRate                  float32 `json:"ioRate"`                  //营收比（总充值金额-总提现金额)÷总充值
	TodayChargePlayerCount  int     `json:"todayChargePlayerCount"`  // 今日充值人数
	TodayChargeTimes        int     `json:"todayChargeTimes"`        // 今日充值笔数
	TodayFirstChargeTimes   int     `json:"todayFirstChargeTimes"`   // 今日首充笔数
	TodayChargeMoney        float32 `json:"todayChargeMoney"`        // 今日充值金额
	TodayFirstChargeMoney   float32 `json:"todayFirstChargeMoney"`   // 今日首充金额
	TodayMultiChargeMoney   float32 `json:"todayMultiChargeMoney"`   // 今日复充金额
	TodayMultiChargePlayers int     `json:"todayMultiChargePlayers"` // 今日复充人数
}

type ServerGeneralizeQueryParam struct {
	PlatformId  string
	ServerId    string   `json:"serverId"`
	ChannelList []string `json:"channelList"`

	EndTime int `json:"endTime"` // 截止日期
}

// 无法按endTime条件约束的标签： 最高等级 , 全服剩余金币， 全服剩余银币
func GetServerGeneralize(endTime int, platformId, serverId string, channelList []string) (*ServerGeneralize, error) {
	gameServer, err := GetGameServerOne(platformId, serverId)
	if err != nil {
		return nil, err
	}
	node := gameServer.Node
	gameDb, err := GetGameDbByNode(node)
	if err != nil {
		return nil, err
	}
	defer gameDb.Close()
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	// 前端入参endTime为当前时间时间戳，正确情况下前端应发送指定日期的0点0分0秒时间戳
	// 为防止前端再次发错，此处将接收到的endTime强制转换为指定日期的0点0分0秒时间戳
	t := time.Unix(int64(endTime), 0)
	SpecifyDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	realEndTime := int(SpecifyDate.Unix())

	chargeTimesList := GetServerChargeCountList(platformId, serverId, channelList, realEndTime+86400)
	//g.Log().Debug("chargeTimesList:%+v", chargeTimesList)
	chargeCount2 := 0
	chargeCount3 := 0
	chargeCount5 := 0
	//chargeCountMore := 0
	for _, e := range chargeTimesList {
		if e >= 2 {
			chargeCount2++
		}
		if e >= 3 {
			chargeCount3++
		}
		if e >= 5 {
			chargeCount5++
		}
		//if e >= 6 {
		//	chargeCountMore++
		//}
	}
	// todayZeroTime := utils.GetTodayZeroTimestamp()

	// 计算充值相关信息
	//var (
	//	selectZeroTime = utils.GetThatZeroTimestamp(int64(endTime))
	//	selectEndTime  = selectZeroTime + 3600*24 - 1
	//)
	//rechargeGenInfo, _ := GetRechargeGeneralize(platformId, serverId, realEndTime, realEndTime+86400)

	//wdlGenInfo := GetWithdrawalInfo(playerIds)

	//g.Log().Info("EndTime: %d", realEndTime)
	//g.Log().Info("EndTime+86400: %d", realEndTime + 86400)

	//今日玩家登录信息
	var todayLoginTimes = 0
	//for _, channalStr := range channelList {
	//	todayLoginTimes += GetThatDayLoginTimes(gameDb, serverId, channalStr, realEndTime)
	//}
	todayLoginTimes = getSpecifyDayLoginPlayerCount(gameDb, serverId, channelList, realEndTime)

	/*todayChargePlayerCount := 0
	if ok := 0 != len(channelList); ok {
		todayChargePlayerCount = GetChargePlayerCount(platformId, serverId, GetSQLWhereParam(channelList), realEndTime, realEndTime+86400)
	} else {
		todayChargePlayerCount = GetChargePlayerCount(platformId, serverId, "", realEndTime, realEndTime+86400)
	}*/
	firstChargePlayers, todayMultiChargePlayers, todayChargePlayerCount, firstChargeTotalMoney, todayMultiChargeMoney, todayChargeMoney :=
		GetChargeInfoBySpecifyDate(platformId, serverId, channelList, realEndTime, realEndTime+86400)

	serverGeneralize := &ServerGeneralize{
		PlatformId:      platformId,
		ServerId:        GetGameServerIdListStringByNode(node),
		OpenTime:        serverNode.OpenTime,
		Version:         GetNodeVersion(node),
		MergeTime:       GetMergeTime(node),
		Status:          serverNode.State,
		TotalRegister:   GetTotalRegisterRoleCountByChannelList(gameDb, serverId, channelList, endTime),
		TotalCreateRole: GetTotalCreateRoleCountByChannelList(gameDb, serverId, channelList, endTime),
		OnlineCount:     GetNowOnlineCount2(gameDb, serverId, channelList),
		MaxLevel:        GetMaxPlayerLevel(gameDb, serverId, channelList),
		//TodayCreateRole:         GetCreateRoleCountByChannelList(gameDb, serverId, channelList, endTime, endTime+86400),
		TodayCreateRole: GetCreateRoleCountByChannelList(gameDb, serverId, channelList, realEndTime, realEndTime+86400),
		//TodayRegister:           GetRegisterRoleCountByChannelList(gameDb, serverId, channelList, endTime, endTime+86400),
		TodayRegister:  GetRegisterRoleCountByChannelList(gameDb, serverId, channelList, realEndTime, realEndTime+86400),
		MaxOnlineCount: GetThatDayMaxOnlineCount(platformId, serverId, channelList, realEndTime, realEndTime+86400),
		//YesterdayMaxOnlineCount: GetThatDayMaxOnlineCount(platformId, serverId, channelList, endTime-86400, endTime),
		YesterdayMaxOnlineCount: GetThatDayMaxOnlineCount(platformId, serverId, channelList, realEndTime-86400, realEndTime),
		TotalChargeIngot:        GetServerTotalChargeIngot(platformId, serverId, channelList),
		TotalChargeMoney:        GetServerTotalChargeMoneyByChannelList(platformId, serverId, channelList, realEndTime+86400),
		TotalChargePlayerCount:  len(chargeTimesList), //GetServerChargePlayerCount(platformId, serverId, channelList, endTime),
		SecondChargePlayerCount: GetServerSecondChargePlayerCount(platformId, serverId, channelList, realEndTime+86400),
		//ChargeCountMore:         chargeCountMore,
		TotalIngot:  GetTotalProp(gameDb, 2, channelList),
		TotalCoin:   GetTotalProp(gameDb, 4, channelList),
		TotalBounty: GetTotalProp(gameDb, 52, channelList),

		TodayLoginTimes: todayLoginTimes,
		//FirstRechargeTimes: rechargeGenInfo.FirstNums,
		//FirstRechargeCount: rechargeGenInfo.RecTotal,
		//ReRechargeTimes:    rechargeGenInfo.ReNums,
		//ReRechargeCount:    rechargeGenInfo.ReAcount,
		//WdlAmount:          wdlGenInfo.Amount,
		//WdlTimes:           wdlGenInfo.Times,
		//WdlPlayerCount:     wdlGenInfo.PlayerCount,
		//IoRate:             CaclRate(int(rechargeGenInfo.Amount-wdlGenInfo.Amount), int(rechargeGenInfo.Amount)),
		TodayChargePlayerCount: todayChargePlayerCount,

		TodayChargeTimes:        todayChargePlayerCount,
		TodayChargeMoney:        todayChargeMoney,
		TodayFirstChargeTimes:   firstChargePlayers,
		TodayFirstChargeMoney:   firstChargeTotalMoney,
		TodayMultiChargePlayers: todayMultiChargePlayers,
		TodayMultiChargeMoney:   todayMultiChargeMoney,
	}

	serverGeneralize.ARPU = CaclRate(int(serverGeneralize.TotalChargeMoney), serverGeneralize.TotalChargePlayerCount)
	serverGeneralize.ChargeRate = CaclRate(serverGeneralize.TotalChargePlayerCount, serverGeneralize.TotalCreateRole)
	serverGeneralize.SecondChargeRate = CaclRate(serverGeneralize.SecondChargePlayerCount, serverGeneralize.TotalChargePlayerCount)
	serverGeneralize.ChargeCount2Rate = CaclRate(chargeCount2, serverGeneralize.TotalChargePlayerCount)
	serverGeneralize.ChargeCount3Rate = CaclRate(chargeCount3, serverGeneralize.TotalChargePlayerCount)
	serverGeneralize.ChargeCount5Rate = CaclRate(chargeCount5, serverGeneralize.TotalChargePlayerCount)

	return serverGeneralize, err
}
