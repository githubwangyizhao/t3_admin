package models

import (
	"fmt"
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

type DailyStatistics struct {
	//Node       string `json:"node" gorm:"primary_key"`
	PlatformId string `json:"platformId" gorm:"primary_key"`
	ServerId   string `json:"serverId" gorm:"primary_key"`
	Channel    string `json:"channel" gorm:"primary_key"`
	Time       int    `json:"time" gorm:"primary_key"`
	Source     int    `json:"source" gorm:"primary_key"` // 1为谷歌支付,2为苹果支付,3为第三方支付通道

	ChargeMoney            float32 `json:"chargeMoney"`
	NewChargeMoney         float32 `json:"newChargeMoney"`
	TotalChargeMoney       float32 `json:"totalChargeMoney"`
	ChargePlayerCount      int     `json:"chargePlayerCount"`
	TotalChargePlayerCount int     `json:"totalChargePlayerCount"`
	ARPU                   float32 `json:"arpu" gorm:"-"`
	ActiveARPU             float32 `json:"active_arpu" gorm:"-"`
	NewChargePlayerCount   int     `json:"newChargePlayerCount"`
	FirstChargePlayerCount int     `json:"firstChargePlayerCount"`
	FirstChargeTotalMoney  float32 `json:"firstChargeTotalMoney"`
	//ActivePlayerCount    int     `json:"activePlayerCount" gorm:"-"`
	ActiveChargeRate float32 `json:"activeChargeRate" gorm:"-"`

	LoginTimes           int `json:"loginTimes"`
	LoginPlayerCount     int `json:"loginPlayerCount"`
	ActivePlayerCount    int `json:"activePlayerCount"`
	CreateRoleCount      int `json:"createRoleCount"`
	ShareCreateRoleCount int `json:"shareCreateRoleCount"`
	TotalCreateRoleCount int `json:"totalCreateRoleCount"`

	MaxOnlineCount int `json:"maxOnline"`
	MinOnlineCount int `json:"minOnline"`
	AvgOnlineCount int `json:"avgOnline"`
	AvgOnlineTime  int `json:"avgOnlineTime"`

	RegisterCount      int `json:"registerCount"`
	TotalRegisterCount int `json:"totalRegisterCount"`
	//CreateRoleCount int `json:"createRoleCount"`
	ValidRoleCount int `json:"validRoleCount"`
}

type DailyStatisticsQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId    string   `json:"serverId"`
	ChannelList []string `json:"channelList"`
	PromoteId   int
	StartTime   int
	EndTime     int
}

func GetIncomeStatisticsChartData(platformId string, serverId string, channelList []string, startTime int, endTime int) []map[string]string {
	today := utils.GetTodayZeroTimestamp()
	if endTime > today {
		endTime = today
	}
	if startTime > today {
		startTime = today
	}
	if endTime < startTime || startTime == 0 {
		g.Log().Error("开始结束时间错误")
		return nil
	}

	chargeData := make([]map[string]string, 0, (endTime-startTime)/86400)
	for i := startTime; i < endTime; i = i + 86400 {
		var data struct {
			ChargeCount     float32
			CreateRoleCount float32
		}
		whereArray := make([]string, 0)

		whereArray = append(whereArray, fmt.Sprintf("platform_id = '%s'", platformId))
		whereArray = append(whereArray, fmt.Sprintf("channel in(%s)", GetSQLWhereParam(channelList)))
		whereArray = append(whereArray, fmt.Sprintf("time = %d", i))
		if serverId != "" {
			whereArray = append(whereArray, fmt.Sprintf("server_id = '%s'", serverId))
		}
		whereParam := strings.Join(whereArray, " and ")
		if whereParam != "" {
			whereParam = " where " + whereParam
		}
		sql := fmt.Sprintf(
			`SELECT sum(total_charge_money) as charge_count FROM daily_statistics %s `, whereParam)
		err := Db.Raw(sql).Scan(&data).Error
		utils.CheckError(err)
		m := make(map[string]string, 2)
		m["时间"] = utils.FormatDate(int64(i))
		m["流水"] = strconv.Itoa(int(data.ChargeCount))
		//m["创角"] = strconv.Itoa(int(data.CreateRoleCount))
		chargeData = append(chargeData, m)
	}
	return chargeData
}

type NewDailyStatistics struct {
	TotalChargeMoney       string `json:"totalChargeMoney"`
	TotalChargePlayerCount string `json:"totalChargePlayerCount"`
	ChargeMoney            string `json:"chargeMoney"`
	NewChargeMoney         string `json:"newChargeMoney"`
	ChargePlayerCount      string `json:"chargePlayerCount"`
	NewChargePlayerCount   string `json:"newChargePlayerCount"`
	FirstChargePlayerCount string `json:"firstChargePlayerCount"`
	FirstChargeTotalMoney  string `json:"firstChargeTotalMoney"`

	ActivePlayerCount    int `json:"activePlayerCount"`
	LoginTimes           int `json:"loginTimes"`
	LoginPlayerCount     int `json:"loginPlayerCount"`
	MaxOnlineCount       int `json:"maxOnlineCount"`
	MinOnlineCount       int `json:"minOnlineCount"`
	AvgOnlineCount       int `json:"avgOnlineCount"`
	AvgOnlineTime        int `json:"avgOnlineTime"`
	RegisterCount        int `json:"registerCount"`
	CreateRoleCount      int `json:"createRoleCount"`
	ValidRoleCount       int `json:"validRoleCount"`
	TotalCreateRoleCount int `json:"totalCreateRoleCount"`
	ShareCreateRoleCount int `json:"shareCreateRoleCount"`
	TotalRegisterCount   int `json:"totalRegisterCount"`

	ARPU             float32 `json:"arpu"`
	ActiveARPU       float32 `json:"active_arpu"`
	ActiveChargeRate string  `json:"activeChargeRate"`
	Time             string  `json:"time"`
	ARPPU            float32 `json:"arppu"`

	Date    string `json:"date"`
	Source  string `json:"source"`
	Channel string `json:"channel"`
}

type DailyStatisticsList4Frontend struct {
	Date                 string `json:"date"`
	Channel              string `json:"channel"`
	Source               int    `json:"source"`
	RegisterCount        int    `json:"registerCount"`
	TotalRegisterCount   int    `json:"totalRegisterCount"`
	CreateRoleCount      int    `json:"createRoleCount"`
	TotalCreateRoleCount int    `json:"totalCreateRoleCount"`
	ValidRoleCount       int    `json:"validRoleCount"`
	ShareCreateRoleCount int    `json:"shareCreateRoleCount"`
	LoginPlayerCount     int    `json:"loginPlayerCount"`
	LoginTimes           int    `json:"loginTimes"`
	ActivePlayerCount    int    `json:"activePlayerCount"`
	AvgOnlineTime        int    `json:"avgOnlineTime"`
	AvgOnlineCount       int    `json:"avgOnlineCount"`

	NewChargeMoney         string `json:"newChargeMoney"`
	NewChargePlayerCount   string `json:"newChargePlayerCount"`
	ChargeMoney            string `json:"chargeMoney"`
	ChargePlayerCount      string `json:"chargePlayerCount"`
	FirstChargePlayerCount string `json:"firstChargePlayerCount"`
	FirstChargeTotalMoney  string `json:"firstChargeTotalMoney"`

	MaxOnlineCount int `json:"maxOnlineCount"`
	MinOnlineCount int `json:"minOnlineCount"`

	TotalNewChargeMoney            float64 `json:"totalNewChargeMoney"`
	TotalNewChargePlayerCount      int     `json:"totalNewChargePlayerCount"`
	TotalChargeMoney               float64 `json:"totalChargeMoney"`
	TotalChargePlayerCount         int     `json:"totalChargePlayerCount"`
	TotalFirstChargePlayerMoney    float64 `json:"totalFirstChargePlayerMoney"`
	TotalFirstChargePlayerCount    int     `json:"totalFirstChargePlayerCount"`
	TotalTotalChargeMoney          float64 `json:"totalTotalChargeMoney"`
	TotalTotalChargePlayerCount    int     `json:"totalTotalChargePlayerCount"`
	TotalTotalChargeMoneyStr       string  `json:"totalTotalChargeMoneyString"`
	TotalTotalChargePlayerCountStr string  `json:"totalTotalChargePlayerCountString"`

	ARPU             float32 `json:"arpu"`
	ActiveARPU       float32 `json:"active_arpu"`
	ActiveChargeRate string  `json:"activeChargeRate"`
	Time             string  `json:"time"`
	ARPPU            float32 `json:"arppu"`
}

func GetDailyStatisticsList(params *DailyStatisticsQueryParam) ([]*DailyStatisticsList4Frontend, int) {
	var (
		date                = "from_unixtime(time, '%Y-%m-%d') as date"
		channelSourceSelect = "channel, source"
		roleSelect          = "register_count,total_register_count, create_role_count, total_create_role_count, valid_role_count, share_create_role_count"
		onlineSelect        = "avg_online_time, avg_online_count, max_online_count, min_online_count"
		loginSelect         = "login_player_count, login_times, active_player_count"
		chargeSelect        = "total_charge_money, total_charge_player_count, charge_money, charge_player_count, new_charge_money, new_charge_player_count, first_charge_total_money, first_charge_player_count"
		whereParam          = ""
		groupBy             = "group by source, channel, date"
		orderBy             = "ORDER BY date DESC, source ASC"
		count               = 0
	)

	whereArray := make([]string, 0)
	if params.EndTime < params.StartTime {
		g.Log().Error("开始结束时间错误")
		return nil, count
	} else {
		whereArray = append(whereArray, fmt.Sprintf(`time between %d and %d`, params.StartTime, params.EndTime))
	}
	if params.PlatformId != "" {
		whereArray = append(whereArray, fmt.Sprintf(`platform_id = '%s'`, params.PlatformId))
	}
	if params.ServerId != "" {
		whereArray = append(whereArray, fmt.Sprintf(`server_id = '%s'`, params.ServerId))
	}
	realChannel := ""
	if len(params.ChannelList) > 0 {
		ChannelStr := GetSQLWhereParam(params.ChannelList)
		if ChannelStr != "'all'" && ChannelStr != "" {
			var realChannel = ChannelStr
			if n := strings.Index(ChannelStr, "all"); n != -1 {
				realChannel = strings.Replace(ChannelStr, "'all',", "", -1)
				realChannel = strings.Replace(realChannel, ",'all'", "", -1)
			}
			whereArray = append(whereArray, fmt.Sprintf(`channel in (%s)`, realChannel))
		}
	}

	if len(whereArray) > 0 {
		whereParam = " WHERE " + strings.Join(whereArray, " AND ")
	}
	data := make([]*NewDailyStatistics, 0, (params.EndTime-params.StartTime)/86400)
	sql := fmt.Sprintf(`SELECT %s, %s, %s, %s, %s, %s FROM daily_statistics %s %s %s`,
		date, channelSourceSelect, roleSelect, onlineSelect, loginSelect, chargeSelect, whereParam, groupBy, orderBy)

	err := Db.Debug().Model(&DailyStatistics{}).Raw(sql).Scan(&data).Error
	utils.CheckError(err)

	returnData := make([]*DailyStatisticsList4Frontend, count)
	matchMap := make(map[string]int, 0)

	serverNode, err := GetGameServerOne(params.PlatformId, params.ServerId)
	gameDb, err := GetGameDbByNode2(serverNode.Node, params.PlatformId, params.ServerId)
	//ActivePlayerCount := GetActivePlayersBySpecifiedDay(gameDb, params.ServerId, "local_test", params.EndTime)
	ActivePlayerCount := GetActivePlayersBySpecifiedPeriod(gameDb, params.ServerId, realChannel, params.StartTime, params.EndTime+86400)
	if err != nil {
		return returnData, 0
	}
	defer gameDb.Close()

	for idx, item := range data {
		matchKey := item.Date
		TotalNewChargeMoney, _ := strconv.ParseFloat(item.NewChargeMoney, 64)
		TotalNewChargePlayerCount, _ := strconv.Atoi(item.NewChargePlayerCount)
		TotalChargeMoney, _ := strconv.ParseFloat(item.ChargeMoney, 64)
		TotalChargePlayerCount, _ := strconv.Atoi(item.ChargePlayerCount)
		TotalFirstChargeTotalMoney, _ := strconv.ParseFloat(item.FirstChargeTotalMoney, 64)
		TotalFirstChargePlayerCount, _ := strconv.Atoi(item.FirstChargePlayerCount)
		TotalTotalChargeMoney, _ := strconv.ParseFloat(item.TotalChargeMoney, 64)
		TotalTotalChargePlayerCount, _ := strconv.Atoi(item.TotalChargePlayerCount)
		source, _ := strconv.Atoi(item.Source)

		if _, ok := matchMap[matchKey]; ok {
			i := len(returnData) - 1

			add2String := false
			if ok := returnData[i].Source != source; ok {
				add2String = true
			}

			// 新增充值玩家数(谷歌+苹果+第三方)
			returnData[i].TotalNewChargePlayerCount += TotalNewChargePlayerCount
			// 谷歌充值新增充值玩家数 - 苹果充值新增充值玩家数 - 第三方通道新增充值玩家数
			if add2String == true {
				returnData[i].NewChargePlayerCount += "-" + item.NewChargePlayerCount
			} else {
				NewChargePlayerCount, _ := strconv.Atoi(returnData[i].NewChargePlayerCount)
				returnData[i].NewChargePlayerCount = gconv.String(NewChargePlayerCount + TotalNewChargePlayerCount)
			}
			// 新增充值总金额(谷歌+苹果+第三方)
			returnData[i].TotalNewChargeMoney += TotalNewChargeMoney
			// 谷歌充值新增充值金额 - 苹果充值新增充值金额 - 第三方通道新增充值金额
			if add2String == true {
				returnData[i].NewChargeMoney += "-" + item.NewChargeMoney
			} else {
				NewChargeMoney, _ := strconv.ParseFloat(returnData[i].NewChargeMoney, 64)
				returnData[i].NewChargeMoney = gconv.String(NewChargeMoney + TotalNewChargeMoney)
			}

			// 累计付费玩家数(谷歌+苹果+第三方)
			returnData[i].TotalTotalChargePlayerCount += TotalTotalChargePlayerCount
			// 谷歌充值累计总玩家数-苹果充值累计总玩家数-第三方通道累计充值总玩家数
			if add2String == true {
				returnData[i].TotalTotalChargePlayerCountStr += "-" + item.TotalChargePlayerCount
			} else {
				TotalTotalChargePlayerCountStr, _ := strconv.Atoi(returnData[i].TotalTotalChargePlayerCountStr)
				returnData[i].TotalTotalChargePlayerCountStr = gconv.String(TotalTotalChargePlayerCountStr + TotalTotalChargePlayerCount)
			}
			// 累计付费金额(谷歌+苹果+第三方)
			returnData[i].TotalTotalChargeMoney += TotalTotalChargeMoney
			// 谷歌充值累计金额-苹果充值累计金额-第三方通道累计充值金额
			if add2String == true {
				returnData[i].TotalTotalChargeMoneyStr += "-" + item.TotalChargeMoney
			} else {
				TotalTotalChargeMoneyStr, _ := strconv.ParseFloat(returnData[i].TotalTotalChargeMoneyStr, 64)
				returnData[i].TotalTotalChargeMoneyStr = gconv.String(TotalTotalChargeMoneyStr + TotalTotalChargeMoney)
			}

			// 今日充值金额 (谷歌+苹果+第三方)
			returnData[i].TotalChargeMoney += TotalChargeMoney
			// 谷歌充值今日充值金额 - 苹果充值今日充值金额 - 第三方充值今日充值金额
			if add2String == true {
				returnData[i].ChargeMoney += "-" + item.ChargeMoney
			} else {
				ChargeMoney, _ := strconv.ParseFloat(returnData[i].ChargeMoney, 64)
				returnData[i].ChargeMoney = gconv.String(ChargeMoney + TotalChargeMoney)
			}
			// 今日充值玩家数（谷歌+苹果+第三方）
			returnData[i].TotalChargePlayerCount += TotalChargePlayerCount
			// 谷歌充值今日充值玩家数 - 苹果充值今日充值玩家数 - 第三方充值今日充值玩家数
			if add2String == true {
				returnData[i].ChargePlayerCount += "-" + item.ChargePlayerCount
			} else {
				ChargePlayerCount, _ := strconv.Atoi(returnData[i].ChargePlayerCount)
				returnData[i].ChargePlayerCount = gconv.String(ChargePlayerCount + TotalChargePlayerCount)
			}

			// 今日新增充值金额 (谷歌+苹果+第三方)
			returnData[i].TotalFirstChargePlayerMoney += TotalFirstChargeTotalMoney
			// 谷歌充值今日新增充值金额 - 苹果充值今日新增充值金额 - 第三方充值今日新增充值金额
			if add2String == true {
				returnData[i].FirstChargeTotalMoney += "-" + item.FirstChargeTotalMoney
			} else {
				FirstChargeTotalMoney, _ := strconv.ParseFloat(returnData[i].FirstChargeTotalMoney, 64)
				returnData[i].FirstChargeTotalMoney = gconv.String(FirstChargeTotalMoney + TotalFirstChargeTotalMoney)
			}
			// 今日新增充值玩家数 (谷歌+苹果+第三方)
			returnData[i].TotalFirstChargePlayerCount += TotalFirstChargePlayerCount
			// 谷歌充值今日新增充值玩家数 - 苹果充值今日新增充值玩家数 - 第三方充值今日新增充值玩家数
			if add2String == true {
				returnData[i].FirstChargePlayerCount += "-" + item.FirstChargePlayerCount
			} else {
				FirstChargePlayerCount, _ := strconv.Atoi(returnData[i].FirstChargeTotalMoney)
				returnData[i].FirstChargePlayerCount = gconv.String(FirstChargePlayerCount + TotalFirstChargePlayerCount)
			}

			if ok := returnData[i].RegisterCount == 0; ok {
				returnData[i].RegisterCount = item.RegisterCount
			}
			if ok := returnData[i].Source == source && returnData[i].Channel != item.Channel; ok {
				returnData[i].TotalRegisterCount += item.TotalRegisterCount
				returnData[i].CreateRoleCount += item.CreateRoleCount
				returnData[i].TotalCreateRoleCount += item.TotalCreateRoleCount
				returnData[i].ValidRoleCount += item.ValidRoleCount
				returnData[i].LoginPlayerCount += item.LoginPlayerCount
				returnData[i].LoginTimes += item.LoginTimes
				returnData[i].ActivePlayerCount += item.ActivePlayerCount
				returnData[i].AvgOnlineTime += item.AvgOnlineTime
				returnData[i].AvgOnlineCount += item.AvgOnlineCount
			}

			if i+1 == len(returnData) {
				calcArpu4DailyStatistics(returnData[i])
			}
		} else {
			matchMap[matchKey] = idx
			tmpData := DailyStatisticsList4Frontend{}
			tmpData.Source = source
			tmpData.Date = item.Date
			tmpData.RegisterCount = item.RegisterCount
			tmpData.TotalRegisterCount = item.TotalRegisterCount
			tmpData.CreateRoleCount = item.CreateRoleCount
			tmpData.TotalCreateRoleCount = item.TotalCreateRoleCount
			tmpData.ValidRoleCount = item.ValidRoleCount
			tmpData.ShareCreateRoleCount = item.ShareCreateRoleCount
			tmpData.LoginPlayerCount = item.LoginPlayerCount
			tmpData.LoginTimes = item.LoginTimes

			if _, ok := ActivePlayerCount[item.Date]; ok {
				tmpData.ActivePlayerCount = ActivePlayerCount[item.Date]
			}
			//tmpData.ActivePlayerCount = item.ActivePlayerCount
			tmpData.AvgOnlineTime = item.AvgOnlineTime
			tmpData.AvgOnlineCount = item.AvgOnlineCount

			tmpData.NewChargeMoney = item.NewChargeMoney
			tmpData.TotalNewChargeMoney = TotalNewChargeMoney

			tmpData.NewChargePlayerCount = item.NewChargePlayerCount
			tmpData.TotalNewChargePlayerCount = TotalNewChargePlayerCount

			tmpData.ChargeMoney = item.ChargeMoney
			tmpData.TotalChargeMoney = TotalChargeMoney

			tmpData.ChargePlayerCount = item.ChargePlayerCount
			tmpData.TotalChargePlayerCount = TotalChargePlayerCount

			tmpData.FirstChargePlayerCount = item.FirstChargePlayerCount
			tmpData.TotalFirstChargePlayerMoney = TotalFirstChargeTotalMoney

			tmpData.TotalTotalChargeMoney = TotalTotalChargeMoney
			tmpData.TotalTotalChargePlayerCount = TotalTotalChargePlayerCount

			tmpData.TotalTotalChargeMoneyStr = item.TotalChargeMoney
			tmpData.TotalTotalChargePlayerCountStr = item.TotalChargePlayerCount

			tmpData.FirstChargeTotalMoney = item.FirstChargeTotalMoney
			tmpData.TotalFirstChargePlayerCount = TotalFirstChargePlayerCount

			if len(returnData)-1 >= 0 {
				pos := len(returnData) - 1
				e := returnData[pos]
				calcArpu4DailyStatistics(e)
			}
			returnData = append(returnData, &tmpData)
		}
	}

	return returnData, count
}

func calcArpu4DailyStatistics(e *DailyStatisticsList4Frontend) {
	if e.ActivePlayerCount > 0 {
		e.ARPU = float32(e.TotalChargeMoney) / float32(e.ActivePlayerCount)
	}
	if e.TotalChargePlayerCount > 0 {
		e.ARPPU = float32(e.TotalChargeMoney) / float32(e.TotalChargePlayerCount)

		chargePlayerCountArr := strings.Split(e.ChargePlayerCount, "-")
		if len(chargePlayerCountArr) >= 3 {
			e.ActiveChargeRate = fmt.Sprintf("%f(%f-%f-%f)",
				float32(e.TotalChargePlayerCount)/float32(e.ActivePlayerCount),
				float32(gconv.Int(chargePlayerCountArr[0]))/float32(e.ActivePlayerCount),
				float32(gconv.Int(chargePlayerCountArr[1]))/float32(e.ActivePlayerCount),
				float32(gconv.Int(chargePlayerCountArr[2]))/float32(e.ActivePlayerCount))
		}
		e.ActiveARPU = float32(e.TotalChargeMoney / float64(e.TotalChargePlayerCount))
	}
	//if e.LoginPlayerCount > 0 {
	//}
}

func GetDailyStatisticsOne(platformId string, serverId string, channel []string, timestamp int) (*DailyStatistics, error) {
	data := &DailyStatistics{}
	err := Db.Model(&DailyStatistics{}).Where(&DailyStatistics{
		PlatformId: platformId,
		ServerId:   serverId,
		Time:       timestamp,
	}).Where("channel in(?)", channel).First(&data).Error
	return data, err
}
func UpdateDailyStatistics(platformId string, serverId string, channelList []*Channel, timestamp int) error {
	g.Log().Infof("UpdateDailyStatistics:%v, %v, %v, %v", platformId, serverId, len(channelList), timestamp)
	serverNode, err := GetGameServerOne(platformId, serverId)
	if err != nil {
		return err
	}
	node := serverNode.Node
	gameDb, err := GetGameDbByNode2(serverNode.Node, platformId, serverId)
	if err != nil {
		return err
	}
	defer gameDb.Close()
	for _, e := range channelList {
		//todo 也可以改一次性抛入三条，但是针对于统计实际情况，没有必要
		err = writeDailyStatisticByChargeType(e, node, gameDb, platformId, serverId, 1, timestamp)
		if err != nil {
			return err
		}
		err = writeDailyStatisticByChargeType(e, node, gameDb, platformId, serverId, 2, timestamp)
		if err != nil {
			return err
		}
		err = writeDailyStatisticByChargeType(e, node, gameDb, platformId, serverId, 3, timestamp)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeDailyStatisticByChargeType(e *Channel, node string, gameDb *gorm.DB, platformId, serverId string, source, timestamp int) error {
	channel := e.Channel
	createRoleCount := GetCreateRoleCount(gameDb, serverId, channel, timestamp, timestamp+86400)
	registerRoleCount := GetRegisterRoleCount(gameDb, serverId, channel, timestamp, timestamp+86400)
	totalChargeMoney := GetTotalChargeMoney(platformId, serverId, channel, 0, timestamp+86400, source)
	chargeMoney := GetTotalChargeMoney(platformId, serverId, channel, timestamp, timestamp+86400, source)
	//g.Log().Info("%d: %d | %d", timestamp, chargeMoney, totalChargeMoney)
	m := &DailyStatistics{
		//	Node:                   serverNode.Node,
		PlatformId:             platformId,
		ServerId:               serverId,
		Channel:                channel,
		Time:                   timestamp,
		TotalChargeMoney:       totalChargeMoney,
		ChargeMoney:            chargeMoney,
		NewChargeMoney:         GetThatDayNewChargeMoney(platformId, serverId, channel, timestamp, source),
		ChargePlayerCount:      GetThatDayServerChargePlayerCount(platformId, serverId, channel, timestamp, source),
		TotalChargePlayerCount: GetThatDayChargePlayerCount(platformId, serverId, channel, timestamp+86400, source),
		NewChargePlayerCount:   GetThadDayNewChargePlayerCount(platformId, serverId, channel, timestamp, source),
		FirstChargePlayerCount: GetThadDayServerFirstChargePlayerCount(platformId, serverId, channel, timestamp, source),
		FirstChargeTotalMoney:  GetThadDayServerFirstChargeTotalMoney(platformId, serverId, channel, timestamp, source),

		LoginTimes:        GetThatDayLoginTimes(gameDb, serverId, channel, timestamp),
		LoginPlayerCount:  GetThatDayLoginPlayerCount(gameDb, serverId, channel, timestamp),
		ActivePlayerCount: GetThatDayActivePlayerCount(gameDb, serverId, channel, timestamp),

		TotalCreateRoleCount: GetHistoryCreateRoleCount(platformId, serverId, channel, timestamp) + createRoleCount,
		TotalRegisterCount:   GetHistoryRegisterRoleCount(platformId, serverId, channel, timestamp) + registerRoleCount,
		AvgOnlineTime:        GetOnlineTime(node, serverId, channel, timestamp),
		RegisterCount:        registerRoleCount,
		CreateRoleCount:      createRoleCount,
		ShareCreateRoleCount: GetThatDayShareCreateRoleCountByChannel(gameDb, serverId, channel, timestamp),
		ValidRoleCount:       GetThatDayValidCreateRoleCountByChannel(gameDb, serverId, channel, timestamp),
		Source:               source,
	}
	err := Db.Save(&m).Error
	if err != nil {
		return err
	}

	return nil
}

func GetHistoryCreateRoleCount(platformId string, serverId string, channel string, time int) int {
	var data struct {
		Count int
	}
	//sql := fmt.Sprintf(
	//	`SELECT sum(create_role_count) as count FROM daily_statistics WHERE platform_id = '%s' and server_id = '%s' and channel = '%s' and time < %d `, platformId, serverId, channel, time)
	//sql := fmt.Sprintf(
	//	`SELECT create_role_count as count FROM daily_statistics WHERE platform_id = '%s' and server_id = '%s' and channel = '%s' and time < %d group by time`, platformId, serverId, channel, time)
	//sql := fmt.Sprintf(
	//	`SELECT create_role_count as count FROM daily_statistics WHERE platform_id = '%s' and server_id = '%s' and channel = '%s' and time between %d and %d group by time`, platformId, serverId, channel, time-86400, time)
	sql := fmt.Sprintf(
		`SELECT total_create_role_count as count FROM daily_statistics WHERE platform_id = '%s' and server_id = '%s' and channel = '%s' and time = %d group by time order by time limit 1`, platformId, serverId, channel, time-86400)
	err := Db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}
func GetHistoryRegisterRoleCount(platformId string, serverId string, channel string, time int) int {
	var data struct {
		Count int
	}
	//sql := fmt.Sprintf(
	//	`SELECT sum(register_count) as count FROM daily_statistics WHERE platform_id = '%s' and server_id = '%s' and channel = '%s' and time < %d `, platformId, serverId, channel, time)
	//sql := fmt.Sprintf(
	//	`SELECT register_count as count FROM daily_statistics WHERE platform_id = '%s' and server_id = '%s' and channel = '%s' and time < %d group by time`, platformId, serverId, channel, time)
	//sql := fmt.Sprintf(
	//	`SELECT register_count as count FROM daily_statistics WHERE platform_id = '%s' and server_id = '%s' and channel = '%s' and time between %d and %d group by time`, platformId, serverId, channel, time-86400, time)
	sql := fmt.Sprintf(
		`SELECT total_register_count as count FROM daily_statistics WHERE platform_id = '%s' and server_id = '%s' and channel = '%s' and time = %d group by time order by time limit 1`, platformId, serverId, channel, time-86400)
	err := Db.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}
