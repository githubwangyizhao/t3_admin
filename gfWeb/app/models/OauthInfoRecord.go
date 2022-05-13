package models

import (
	"fmt"
	"gfWeb/library/utils"
	"strings"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

type OauthInfoRecord struct {
	OrderId string `json:"orderId" gorm:"primary_key"`
	//ChargeType    int     `json:"chargeType"`
	Ip         string  `json:"ip"`
	PlatformId string  `json:"platformId"`
	ServerId   string  `json:"serverId"`
	Channel    string  `json:"channel"`
	AccId      string  `json:"accId"`
	Time       int     `json:"time"`
	PlayerId   int     `json:"playerId"`
	Name       string  `json:"playerName" gorm:"-"`
	Amount     float32 `json:"amount"`
	PropType   int     `json:"propType"`
	PropId     int     `json:"propId"`
	Num        int     `json:"num"`
}

type OauthInfoRecordQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId    string   `json:"serverId"`
	ChannelList []string `json:"channelList"`
	promoteId   int
	PlayerId    int
	PlayerName  string
	OrderId     string
	AccId       string
	StartTime   int
	EndTime     int
	PropType    int
	PropId      int
	Promote     string `json:"promote"`
}

func GetOauthInfoRecordList(params *OauthInfoRecordQueryParam) ([]*OauthInfoRecord, int64, int64, float64, float64) {
	g.Log().Debug("GetOauthInfoRecordList")
	data := make([]*OauthInfoRecord, 0)

	whereArray := make([]string, 0)

	if params.Promote != "" {
		AccIdList, AccErr := GetGlobalAccountByPromote(params.Promote)
		utils.CheckError(AccErr)
		g.Log().Debug("AccIdList:", params.Promote, AccIdList)

		if len(AccIdList) > 0 {
			whereArray = append(whereArray, fmt.Sprintf(" g_player.account in (%s) ", GetSQLWhereParam(AccIdList)))
		} else {
			return data, 0, 0, 0.0, 0.0
		}
	}
	//if params.Promote != "" {
	//	whereArray = append(whereArray, fmt.Sprintf(" acc_id in (%s) ", GetSQLWhereParam(AccIdList)))
	//}

	//whereArray = append(whereArray, " ")
	if params.PlayerId > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" oauth.player_id = '%d' ", params.PlayerId))
	}
	if params.PlatformId != "" {
		whereArray = append(whereArray, " g_player.platform_id =  '"+params.PlatformId+"' ")
	}
	if params.ServerId != "" {
		whereArray = append(whereArray, " g_player.server_id =  '"+params.ServerId+"' ")
	}

	whereArray = append(whereArray, " oauth.change_type = 0 ")

	if params.StartTime > 0 {
		whereArray = append(whereArray, fmt.Sprintf("oauth.create_time between %d and %d", params.StartTime, params.EndTime))
	}

	if params.PropType > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" oauth.prop_type = '%d' ", params.PropType))
	}
	if params.PropId > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" oauth.prop_id = '%d' ", params.PropId))
	}

	if len(params.ChannelList) > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" g_player.channel in (%s) ", GetSQLWhereParam(params.ChannelList)))
	}

	if params.OrderId != "" {
		whereArray = append(whereArray, " oauth.order_id =  '"+params.OrderId+"' ")
	}

	if params.AccId != "" {
		whereArray = append(whereArray, " g_player.account =  '"+params.AccId+"' ")
	}

	whereParam := strings.Join(whereArray, " and ")
	if whereParam != "" {
		whereParam = " where " + whereParam
	}

	sortOrder := " 'oauth.create_time' desc "
	if params.Order == "ascending" {
		sortOrder = " 'oauth.create_time' asc "
	}

	sql := fmt.Sprintf(
		`select g_player.platform_id as platform_id,oauth.order_id as order_id,oauth.amount as amount,
				oauth.player_id as player_id,oauth.prop_type as prop_type,oauth.prop_id as prop_id,oauth.change_num as num,
				oauth.ip as ip,oauth.create_time as time,g_player.channel as channel,g_player.server_id as server_id,
				g_player.nickanme as name,g_player.account as acc_id 
				from oauth_order_log as oauth
				inner join global_player as g_player on player_id = g_player.id
   				%s 
				order by %s 
				limit %d,%d; `,
		whereParam,
		sortOrder,
		params.Offset,
		params.Limit,
	)
	g.Log().Debug("sql:", sql)

	err := DbCenter.Raw(sql).Scan(&data).Error

	utils.CheckError(err)

	orgSql := fmt.Sprintf(
		`select g_player.platform_id as platform_id,oauth.order_id as order_id,oauth.amount as amount,oauth.player_id as player_id,oauth.prop_type as prop_type,oauth.prop_id as prop_id,oauth.change_num as num, oauth.status as status, oauth.ip as ip,oauth.create_time as time,g_player.channel as channel,g_player.server_id as server_id,g_player.nickanme as name,g_player.account as acc_id from oauth_order_log as oauth inner join global_player as g_player on player_id = g_player.id %s `,
		whereParam,
	)
	sql = fmt.Sprintf(`(select sum(amount) as money_count, count(DISTINCT player_id) as player_count, count(order_id) as total from (%s) as t );`, orgSql)
	g.Log().Debug("sql:", sql)
	var sumData struct {
		MoneyCount  float64
		PlayerCount int64
		Total       int64
	}
	g.Log().Debug("eee: ", sql)
	err = DbCenter.Raw(sql).Scan(&sumData).Error
	utils.CheckError(err)
	g.Log().Debug("aaa:", data)

	var ids = make([]string, 0)
	for _, v := range data {
		ids = append(ids, gconv.String(v.PlayerId))
	}

	ids = utils.RemoveReStrData(ids)
	var rcSum = struct {
		Sum float64 `gorm:"sum"`
	}{}
	sql = fmt.Sprintf("SELECT sum(total_money) sum FROM player_charge_info_record WHERE player_id IN(%s)", strings.Join(ids, ","))
	err = DbCharge.Raw(sql).Scan(&rcSum).Error
	utils.CheckError(err)

	fmt.Println("rc sum  : ", rcSum)

	return data, sumData.Total, sumData.PlayerCount, sumData.MoneyCount, rcSum.Sum
}

type WithdrawalGenInfo struct {
	Amount      float32 `gorm:"money_count"`
	PlayerCount int     `gorm:"player_count"`
	Times       int     `gorm:"total"`
}

// 获取提现总金额 ， 人数， 笔数
func GetWithdrawalInfo(playerIds []int) WithdrawalGenInfo {
	if len(playerIds) < 1 {
		return WithdrawalGenInfo{}
	}

	var ids = make([]string, 0)
	for _, v := range playerIds {
		ids = append(ids, gconv.String(v))
	}

	var info WithdrawalGenInfo
	sql := fmt.Sprintf("select sum(amount) as money_count, count(DISTINCT player_id) as player_count, count(order_id) as total from oauth_order_log where player_id in(%s)",
		GetSQLWhereParam(ids))

	fmt.Println(sql)
	err := DbCenter.Raw(sql).Scan(&info).Error
	utils.CheckError(err)

	return info
}
