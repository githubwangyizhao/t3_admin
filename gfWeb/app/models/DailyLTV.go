package models

import (
	"fmt"
	"gfWeb/library/utils"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

type DailyRegisterPlayer struct {
	RegisterRole int    `json:"registerRole" gorm:"-"`
	CreateRole   int    `json:"createRole" gorm:"-"`
	Channel      string `json:"channel" gorm:"primary_key"`
	Date         string `json:"date"`
}

//注册充值 表 charge_info_record
type LtvMoney struct {
	//Node         string  `json:"node" gorm:"primary_key"`
	// Time       string  `json:"time" gorm:"column:time"`
	// PlatformId string `json:"platformId" gorm:"primary_key"`
	ServerId string `json:"serverId" gorm:"primary_key"`
	// Channel    string `json:"channel" gorm:"primary_key"`
	Time         int     `json:"time" gorm:"primary_key"`
	RegisterRole int     `json:"registerRole" gorm:"-"`
	CreateRole   int     `json:"createRole" gorm:"-"`
	PlatformId   string  `json:"platform_id" gorm:"plat_id"`
	Money        string  `json:"money"`
	Days         float32 `json:"days"` //相差多少天，最多不超过120天
	RegTime      int     `json:"reg_time" gorm:"reg_time"`
	RecordTime   string  `json:"record_time" gorm:"record_time"`
	RegDate      string  `json:"reg_data" gorm:"reg_data"`

	Sum0   float32 `json:"sum0"  gorm:"-"`
	Sum1   float32 `json:"sum1"  gorm:"-"`
	Sum2   float32 `json:"sum2" gorm:"-"`
	Sum3   float32 `json:"sum3" gorm:"-"`
	Sum4   float32 `json:"sum4" gorm:"-"`
	Sum5   float32 `json:"sum5" gorm:"-"`
	Sum6   float32 `json:"sum6" gorm:"-"`
	Sum7   float32 `json:"sum7" gorm:"-"`
	Sum8   float32 `json:"sum8" gorm:"-"`
	Sum9   float32 `json:"sum9" gorm:"-"`
	Sum10  float32 `json:"sum10" gorm:"-"`
	Sum11  float32 `json:"sum11" gorm:"-"`
	Sum12  float32 `json:"sum12" gorm:"-"`
	Sum13  float32 `json:"sum13" gorm:"-"`
	Sum14  float32 `json:"sum14" gorm:"-"`
	Sum15  float32 `json:"sum15" gorm:"-"`
	Sum30  float32 `json:"sum30" gorm:"-"`
	Sum60  float32 `json:"sum60" gorm:"-"`
	Sum90  float32 `json:"sum90" gorm:"-"`
	Sum120 float32 `json:"sum120" gorm:"-"`
	Sum    float32 `json:"sum" gorm:"-"`
}

type DailyLTVV1 struct {
	//Node         string  `json:"node" gorm:"primary_key"`
	PlatformId string `json:"platformId" gorm:"primary_key"`
	ServerId   string `json:"serverId" gorm:"primary_key"`
	Channel    string `json:"channel" gorm:"primary_key"`
	Date       string `json:"date"`
	C1         int    `json:"c1"  gorm:"column:c1"`
	C2         int    `json:"c2" gorm:"column:c2"`
	C3         int    `json:"c3" gorm:"column:c3"`
	C7         int    `json:"c7" gorm:"column:c7"`
	C14        int    `json:"c14" gorm:"column:c14"`
	C30        int    `json:"c30" gorm:"column:c30"`
	C60        int    `json:"c60" gorm:"column:c60"`
	C90        int    `json:"c90" gorm:"column:c90"`
	C120       int    `json:"c120" gorm:"column:c120"`
}

type DailyLTV struct {
	//Node         string  `json:"node" gorm:"primary_key"`
	PlatformId   string  `json:"platformId" gorm:"primary_key"`
	ServerId     string  `json:"serverId" gorm:"primary_key"`
	Channel      string  `json:"channel" gorm:"primary_key"`
	Time         int     `json:"time" gorm:"primary_key"`
	RegisterRole int     `json:"registerRole" gorm:"-"`
	CreateRole   int     `json:"createRole" gorm:"-"`
	C1           int     `json:"c1"  gorm:"column:c1"`
	C2           int     `json:"c2" gorm:"column:c2"`
	C3           int     `json:"c3" gorm:"column:c3"`
	C7           int     `json:"c7" gorm:"column:c7"`
	C14          int     `json:"c14" gorm:"column:c14"`
	C30          int     `json:"c30" gorm:"column:c30"`
	C60          int     `json:"c60" gorm:"column:c60"`
	C90          int     `json:"c90" gorm:"column:c90"`
	C120         int     `json:"c120" gorm:"column:c120"`
	LTV1         float32 `json:"ltv1"  gorm:"-"`
	LTV2         float32 `json:"ltv2" gorm:"-"`
	LTV3         float32 `json:"ltv3" gorm:"-"`
	LTV7         float32 `json:"ltv7" gorm:"-"`
	LTV14        float32 `json:"ltv14" gorm:"-"`
	LTV30        float32 `json:"ltv30" gorm:"-"`
	LTV60        float32 `json:"ltv60" gorm:"-"`
	LTV90        float32 `json:"ltv90" gorm:"-"`
	LTV120       float32 `json:"ltv120" gorm:"-"`
}

type DailyLTVQueryParam struct {
	BaseQueryParam
	PlatformId string
	//ServerId   string
	ServerId    string   `json:"serverId"`
	ChannelList []string `json:"channelList"`
	Channel     string
	StartTime   int
	EndTime     int
	Promote     string `json:"promote"`
}

type RegRechargeParam struct {
	BaseQueryParam
	PlatformId string
	//ServerId   string
	ServerId    string   `json:"serverId"`
	ChannelList []string `json:"channelList"`
	Channel     string
	StartTime   int
	EndTime     int
	Promote     string `json:"promote"`
}

func GetLtvMoney(params RegRechargeParam) []*LtvMoney {
	if params.EndTime < params.StartTime {
		g.Log().Error("开始结束时间错误")
		return nil
	}

	realData := make([]*LtvMoney, 0)

	data := make([]*LtvMoney, 0)
	channelLen := len(params.ChannelList)

	whereArray := make([]string, 0)
	whereArray1 := make([]string, 0)

	if params.ServerId != "" {
		whereArray = append(whereArray, " server_id =  '"+params.ServerId+"' ")
		whereArray1 = append(whereArray1, " server_id =  '"+params.ServerId+"' ")
	}

	if channelLen > 0 {
		whereArray = append(whereArray, fmt.Sprintf(` channel in (%s) `, GetSQLWhereParam(params.ChannelList)))
		whereArray1 = append(whereArray1, fmt.Sprintf(` channel in (%s) `, GetSQLWhereParam(params.ChannelList)))
	}

	whereArray = append(whereArray, fmt.Sprintf("reg_time between %d and %d", params.StartTime, params.EndTime))
	whereArray1 = append(whereArray1, fmt.Sprintf("time between %d and %d", params.StartTime, params.EndTime))

	if params.PlatformId != "" {
		whereArray = append(whereArray, fmt.Sprintf(" part_id = '%s' ", params.PlatformId))
		whereArray1 = append(whereArray1, fmt.Sprintf(" platform_id = '%s' ", params.PlatformId))
	}

	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}

	whereParam1 := strings.Join(whereArray1, " and ")
	if whereParam1 != "" {
		whereParam1 = " where " + whereParam1
	}

	var (
		date       = "FROM_UNIXTIME(time, '%Y-%m-%d')"
		recordDate = "FROM_UNIXTIME(record_time, '%Y-%m-%d')"
		regDate    = "FROM_UNIXTIME(reg_time, '%Y-%m-%d')"
	)

	sql := fmt.Sprintf("select sum(money) as money, UNIX_TIMESTAMP(%s) as reg_time, ((UNIX_TIMESTAMP(%s) - UNIX_TIMESTAMP(%s))/86400) as days from charge_info_record %s"+
		" and (record_time - reg_time) < (121 * 86400) group by days, reg_time",
		regDate, recordDate, regDate, whereParam)
	err := DbCharge.Raw(sql).Find(&data).Error
	utils.CheckError(err)

	registerCount := make([]*DailyRegisterPlayer, (params.EndTime-params.StartTime)/86400)
	whereParam1 += " and source = 1"
	sql = fmt.Sprintf(
		`select server_id, platform_id, %s as date, SUM(create_role_count) as create_role, SUM(register_count) as register_role from daily_statistics %s group by date, server_id, platform_id`,
		date,
		whereParam1,
	)
	err = Db.Raw(sql).Find(&registerCount).Error
	utils.CheckError(err)
	//exchangeRate := GetExchangeRate(params.PlatformId)
	// 转cny
	//cnyExchangeRate := GetExchangeRate("")

	for i := params.StartTime; i <= params.EndTime; i++ {

		matchDate := time.Unix(int64(i), 0).Format("2006-01-02")

		realSingleData := LtvMoney{
			Time: int(i),
		}
		timeLayout := "2006-01-02"           //转化所需模板
		loc, _ := time.LoadLocation("Local") //重要：获取时区
		theTime, _ := time.ParseInLocation(timeLayout, matchDate, loc)
		matchTimes := int(theTime.Unix())

		realSingleData.RegDate = matchDate

		for _, v := range registerCount {
			if v.Date == matchDate {
				realSingleData.CreateRole = v.CreateRole
				realSingleData.RegisterRole = v.RegisterRole
			}
		}
		for _, e := range data {
			// realSingleData.PlatformId = e.PlatformId
			// realSingleData.ServerId = e.ServerId
			// realSingleData.RegTime = e.RegTime
			// realSingleData.RecordTime = e.RecordTime

			if e.RegTime == matchTimes {
				mny := gconv.Float32(e.Money)
				switch e.Days {
				case 0:
					//realSingleData.Sum0 = mny
					realSingleData.Sum1 = mny
				case 1:
					//realSingleData.Sum1 = mny
					realSingleData.Sum2 = mny
				case 2:
					//realSingleData.Sum2 = mny
					realSingleData.Sum3 = mny
				case 3:
					//realSingleData.Sum3 = mny
					realSingleData.Sum4 = mny
				case 4:
					//realSingleData.Sum4 = mny
					realSingleData.Sum5 = mny
				case 5:
					//realSingleData.Sum5 = mny
					realSingleData.Sum6 = mny
				case 6:
					//realSingleData.Sum6 = mny
					realSingleData.Sum7 = mny
				case 7:
					realSingleData.Sum8 = mny
				case 8:
					realSingleData.Sum9 = mny
				case 9:
					realSingleData.Sum10 = mny
				case 10:
					realSingleData.Sum11 = mny
				case 11:
					realSingleData.Sum12 = mny
				case 12:
					realSingleData.Sum13 = mny
				case 13:
					realSingleData.Sum14 = mny
				case 14:
					realSingleData.Sum15 = mny
				case 30:
					realSingleData.Sum30 = mny
				case 60:
					realSingleData.Sum60 = mny
				case 120:
					realSingleData.Sum120 = mny
				}

				realSingleData.Sum = realSingleData.Sum0 + realSingleData.Sum1 + realSingleData.Sum2 + realSingleData.Sum3 + realSingleData.Sum4 +
					realSingleData.Sum5 + realSingleData.Sum6 + realSingleData.Sum7 + realSingleData.Sum14 + realSingleData.Sum30 + realSingleData.Sum60 +
					realSingleData.Sum120
			}

		}

		realData = append(realData, &realSingleData)
		i += 86399
	}

	return realData
}

func GetDailyLTVListV1(params DailyLTVQueryParam) []*DailyLTV {
	if params.EndTime <= params.StartTime {
		g.Log().Error("开始结束时间错误")
		return nil
	}

	realData := make([]*DailyLTV, 0, (params.EndTime-params.StartTime)/86400)

	data := make([]*DailyLTVV1, 0, (params.EndTime-params.StartTime)/86400)
	channelLen := len(params.ChannelList)

	whereArray := make([]string, 0)
	if params.PlatformId != "" {
		whereArray = append(whereArray, fmt.Sprintf(" platform_id = '%s' ", params.PlatformId))
	}

	if params.ServerId != "" {
		whereArray = append(whereArray, " server_id =  '"+params.ServerId+"' ")
	}

	if channelLen > 0 {
		whereArray = append(whereArray, fmt.Sprintf(` channel in (%s) `, GetSQLWhereParam(params.ChannelList)))
	}

	whereArray = append(whereArray, fmt.Sprintf("time between %d and %d", params.StartTime, params.EndTime))

	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}

	date := "FROM_UNIXTIME(time, '%Y-%m-%d')"
	sql := fmt.Sprintf(
		`select server_id, platform_id, %s as date, SUM(c1) As c1, SUM(c2) As c2, SUM(c3) As c3, SUM(c7) As c7, SUM(c14) As c14, SUM(c30) As c30, SUM(c60) As c60, SUM(c90) As c90, SUM(c120) As c120 from daily_ltv %s group by date, server_id, platform_id order by date asc`,
		date,
		whereParam,
	)

	err := Db.Raw(sql).Find(&data).Error
	utils.CheckError(err)

	whereArray = append(whereArray, "source = 1")
	whereParam = strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}
	registerCount := make([]*DailyRegisterPlayer, (params.EndTime-params.StartTime)/86400)
	sql = fmt.Sprintf(
		`select server_id, platform_id, %s as date, SUM(create_role_count) as create_role, SUM(register_count) as register_role from daily_statistics %s group by date, server_id, platform_id`,
		date,
		whereParam,
	)
	g.Log().Infof("dailyStatistics sql: %s", sql)
	err = Db.Raw(sql).Find(&registerCount).Error
	utils.CheckError(err)

	//exchangeRate := GetExchangeRate(params.PlatformId)
	// 转cny
	//cnyExchangeRate := GetExchangeRate("")

	for i := params.StartTime; i <= params.EndTime; i++ {
		i += 86399
		matchDate := time.Unix(int64(i), 0).Format("2006-01-02")

		realSingleData := DailyLTV{
			Time: int(i),
		}

		for _, v := range registerCount {
			if v.Date == matchDate {
				realSingleData.CreateRole = v.CreateRole
				realSingleData.RegisterRole = v.RegisterRole
			}
		}

		for _, e := range data {
			if e.Date == matchDate {
				realSingleData.C1 = e.C1
				realSingleData.C2 = e.C1
				realSingleData.C3 = e.C3
				realSingleData.C7 = e.C7
				realSingleData.C14 = e.C14
				realSingleData.C30 = e.C30
				realSingleData.C60 = e.C60
				realSingleData.C90 = e.C90
				realSingleData.C120 = e.C120
				if realSingleData.CreateRole > 0 {
					realSingleData.LTV1 = float32(e.C1) / float32(realSingleData.CreateRole)     // / exchangeRate
					realSingleData.LTV2 = float32(e.C2) / float32(realSingleData.CreateRole)     // / exchangeRate
					realSingleData.LTV3 = float32(e.C3) / float32(realSingleData.CreateRole)     // / exchangeRate
					realSingleData.LTV7 = float32(e.C7) / float32(realSingleData.CreateRole)     // / exchangeRate
					realSingleData.LTV14 = float32(e.C14) / float32(realSingleData.CreateRole)   // / exchangeRate
					realSingleData.LTV30 = float32(e.C30) / float32(realSingleData.CreateRole)   // / exchangeRate
					realSingleData.LTV60 = float32(e.C60) / float32(realSingleData.CreateRole)   // / exchangeRate
					realSingleData.LTV90 = float32(e.C90) / float32(realSingleData.CreateRole)   // / exchangeRate
					realSingleData.LTV120 = float32(e.C120) / float32(realSingleData.CreateRole) // / exchangeRate
				}
			}
		}

		realData = append(realData, &realSingleData)
	}

	return realData
}

// 获取总体留存
func GetDailyLTVList(params DailyLTVQueryParam) []*DailyLTV {
	if params.EndTime < params.StartTime {
		g.Log().Error("开始结束时间错误")
		return nil
	}
	data := make([]*DailyLTV, 0, (params.EndTime-params.StartTime)/86400)
	channelLen := len(params.ChannelList)

	for i := params.StartTime; i <= params.EndTime; i = i + 86400 {
		tmpData := make([]*DailyLTV, 0, channelLen)
		err := Db.Model(&DailyLTV{}).Where(&DailyLTV{PlatformId: params.PlatformId, ServerId: params.ServerId, Time: i}).Where("channel in (?)", params.ChannelList).Find(&tmpData).Error
		//g.Log().Debug("ok")
		utils.CheckError(err)
		if len(tmpData) > 0 {
			tmpE := &DailyLTV{
				PlatformId: params.PlatformId,
				ServerId:   params.ServerId,
				Time:       i,
			}
			//g.Log().Debug("1")
			for _, e := range tmpData {
				dailyStatistics, err := GetDailyStatisticsOne(e.PlatformId, e.ServerId, []string{e.Channel}, e.Time)
				//g.Log().Debug("11")
				utils.CheckError(err)
				tmpE.RegisterRole += dailyStatistics.RegisterCount
				tmpE.CreateRole += dailyStatistics.CreateRoleCount
				if e.C1 > 0 {
					tmpE.C1 += e.C1
				}
				if e.C2 > 0 {
					tmpE.C2 += e.C2
				}
				if e.C3 > 0 {
					tmpE.C3 += e.C3
				}
				if e.C7 > 0 {
					tmpE.C7 += e.C7
				}
				if e.C14 > 0 {
					tmpE.C14 += e.C14
				}
				if e.C30 > 0 {
					tmpE.C30 += e.C30
				}
				if e.C60 > 0 {
					tmpE.C60 += e.C60
				}
				if e.C90 > 0 {
					tmpE.C90 += e.C90
				}
				if e.C120 > 0 {
					tmpE.C120 += e.C120
				}
			}
			//g.Log().Debug("2")
			data = append(data, tmpE)
		}
	}

	var exchangeRate float32
	switch params.PlatformId {
	case "indonesia":
		exchangeRate = 14410
	default:
		exchangeRate = 8
	}

	for _, e := range data {
		if e.CreateRole > 0 {
			e.LTV1 = float32(e.C1) / float32(e.CreateRole) / exchangeRate
			e.LTV2 = float32(e.C2) / float32(e.CreateRole) / exchangeRate
			e.LTV3 = float32(e.C3) / float32(e.CreateRole) / exchangeRate
			e.LTV7 = float32(e.C7) / float32(e.CreateRole) / exchangeRate
			e.LTV14 = float32(e.C14) / float32(e.CreateRole) / exchangeRate
			e.LTV30 = float32(e.C30) / float32(e.CreateRole) / exchangeRate
			e.LTV60 = float32(e.C60) / float32(e.CreateRole) / exchangeRate
			e.LTV90 = float32(e.C90) / float32(e.CreateRole) / exchangeRate
			e.LTV120 = float32(e.C120) / float32(e.CreateRole) / exchangeRate

		}
	}
	return data
	//data := make([]*DailyLTV, 0)
	//var count int64
	//f := func(db *gorm.DB) *gorm.DB {
	//	if params.StartTime > 0 {
	//		return db.Where("time between ? and ?", params.StartTime, params.EndTime)
	//	}
	//	return db
	//}
	//err := f(Db.Model(&DailyLTV{}).Where(&DailyLTV{PlatformId: params.PlatformId, ServerId: params.ServerId, Channel: params.Channel})).Where(" channel in (?)", params.ChannelList).Count(&count).Offset(params.Offset).Limit(params.Limit).Find(&data).Error
	//utils.CheckError(err)
	//for _, e := range data {
	//	dailyStatistics, err := GetDailyStatisticsOne(e.PlatformId, e.ServerId, e.Channel, e.Time)
	//	//g.Log().Info("dailyStatistics:%+v", dailyStatistics)
	//	utils.CheckError(err)
	//	e.CreateRole = dailyStatistics.CreateRoleCount
	//	e.RegisterRole = dailyStatistics.RegisterCount
	//}
	//return data, count
}

//更新 每日ltv
func UpdateDailyLTV(platformId string, serverId string, channelList []*Channel, timestamp int) error {
	g.Log().Infof("更新每日ltv:%v, %v, %v, %v", platformId, serverId, len(channelList), timestamp)
	gameServer, err := GetGameServerOne(platformId, serverId)
	if err != nil {
		return err
	}
	node := gameServer.Node
	serverNode, err := GetServerNode(node)
	if err != nil {
		return err
	}
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return err
	}
	defer gameDb.Close()
	openDayZeroTimestamp := utils.GetThatZeroTimestamp(int64(serverNode.OpenTime))

	for _, e := range channelList {
		channel := e.Channel
		for i := 1; i <= 120; i++ {
			if i == 1 || i == 2 || i == 3 || i == 7 || i == 14 || i == 30 || i == 60 || i == 90 || i == 120 {

			} else {
				continue
			}
			thatDayZeroTimestamp := timestamp - i*86400
			if openDayZeroTimestamp > thatDayZeroTimestamp {
				continue
			}

			//dailyStatistics, err := GetDailyStatisticsOne(platformId, serverId, channel, thatDayZeroTimestamp)

			//registerNum := dailyStatistics.RegisterCount
			//createNum := dailyStatistics.CreateRoleCount
			totalCharge := GetTotalChargeMoneyByRegisterTime(platformId, serverId, channel, 0, timestamp, thatDayZeroTimestamp)
			if totalCharge > 0 {
				rate := totalCharge
				//if createNum > 0 {
				//	rate = int(float32(totalCharge) / float32(createNum))
				//}

				m := &DailyLTV{
					//Node:       node,
					PlatformId: platformId,
					ServerId:   serverId,
					Channel:    channel,
					Time:       thatDayZeroTimestamp,
					C1:         0,
					C2:         0,
					C3:         0,
					C7:         0,
					C14:        0,
					C30:        0,
					C60:        0,
					C90:        0,
					C120:       0,
				}
				err = Db.FirstOrCreate(&m).Error
				if err != nil {
					return err
				}
				switch i {
				case 1:
					err = Db.Model(&m).Update("c1", rate).Error
				case 2:
					err = Db.Model(&m).Update("c2", rate).Error
				case 3:
					err = Db.Model(&m).Update("c3", rate).Error
				case 7:
					err = Db.Model(&m).Update("c7", rate).Error
				case 14:
					err = Db.Model(&m).Update("c14", rate).Error
				case 30:
					err = Db.Model(&m).Update("c30", rate).Error
				case 60:
					err = Db.Model(&m).Update("c60", rate).Error
				case 90:
					err = Db.Model(&m).Update("c90", rate).Error
				case 120:
					err = Db.Model(&m).Update("c120", rate).Error
				}
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}
