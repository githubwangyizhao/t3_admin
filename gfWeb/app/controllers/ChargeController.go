package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"reflect"
	"strings"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type ChargeController struct {
	BaseController
}

type ChargeRefundQueryParam struct {
	PlatformId       string
	PlatformOrderId  string
	PlatformUserId   string
	RefundType       int
	RefundReason     string
	OrderId          string
	ApplyRefundMoney int
}

// {"error_code":0,"error_msg":{"pay_info":"https://www.baidu.com","tx_orderno":"tx_orderno"}}
type payInfo struct {
	PayInfo   string `json:"pay_info"`
	TxOrderNo string `json:"tx_orderno"`
}

func (c *ChargeController) CreateOrder(r *ghttp.Request) {
	type CreateOrderParam struct {
		PlatformId string `json:"platformId"`
		PlayerId   int    `json:"playerId"`
		ServerId   string `json:"serverId"`
		ItemId     int    `json:"itemId"`
	}
	params := &CreateOrderParam{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	g.Log().Debug("param: ", params)

	url := utils.GetChargeURL() + "/customer/charge"
	g.Log().Debug("url: ", url)

	data2, err := json.Marshal(params)
	g.Log().Debug("data2: ", data2)
	utils.CheckError(err)

	jsonSlice2, err := utils.HttpRequestGetObj(url, string(data2))
	g.Log().Debug("jsonSlice2: ", jsonSlice2)
	_, ok := jsonSlice2.(string)
	if ok == true {
		c.HttpResult(r, enums.CodeFail, jsonSlice2.(string), jsonSlice2.(string))
	} else {
		PayUrl := jsonSlice2.(map[string]interface{})["pay_info"].(string)
		TxTradeNo := jsonSlice2.(map[string]interface{})["tx_orderno"].(string)
		g.Log().Debug("items: ", reflect.TypeOf(PayUrl))

		g.Log().Debug("PayInfo: ", PayUrl)

		Data := &payInfo{}
		Data.PayInfo = PayUrl
		Data.TxOrderNo = TxTradeNo
		g.Log().Debug("aaa: ", Data)

		//payInfo := "https://www.baidu.com"
		c.HttpResult(r, enums.CodeSuccess, "订单创建成功", Data)
	}
}

type ItemList struct {
	ItemId []string `json:"item_id"`
}

func (c *ChargeController) GetPaymentInfoByPlatform(r *ghttp.Request) {
	type GetPaymentInfoByPlatformParam struct {
		Platform string `json:"platform"`
	}
	params := &GetPaymentInfoByPlatformParam{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	g.Log().Debug("param: ", params)

	url := utils.GetChargeURL() + "/customer/items"
	g.Log().Debug("url: ", url)
	//url := "http://192.168.31.153:6666" + "/update_version"

	data2, err := json.Marshal(params)
	g.Log().Debug("data2: ", data2)
	utils.CheckError(err)

	jsonSlice2, err := utils.HttpRequestGetObj(url, string(data2))

	itemsIdStr := jsonSlice2.(map[string]interface{})["items"].(string)
	g.Log().Debug("items: ", reflect.TypeOf(itemsIdStr))

	Data := &ItemList{}
	Data.ItemId = strings.Split(itemsIdStr, ",")
	g.Log().Debug("fff: ", reflect.TypeOf(strings.Split(itemsIdStr, ",")))
	c.HttpResult(r, enums.CodeSuccess, "获取玩家成功", Data)
}

// 获取充值列表
func (c *ChargeController) ChargeList(r *ghttp.Request) {

	var params models.ChargeInfoRecordQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if params.PlatformId == "" {
		c.CheckError(gerror.New("平台ID不能为空"))
	}
	if params.ServerId == "" {
		c.CheckError(gerror.New("选区不能为空"))
	}

	// 20210422 sl@添加
	subDb, err := models.GetGameDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	if err != nil {
		c.CheckError(err)
	}
	defer subDb.Close()

	var hasNotPlatformOrderData = true //判断是否查询得到了PlayerChargeRecord表的第三方订单
	if params.PlatFormOrderId != "" {
		playerChargeRecord, notFound := models.FindOneByPlatformOrderId(subDb, params.PlatFormOrderId)
		if notFound {
			//无数据返回
			hasNotPlatformOrderData = false
		} else {
			params.OrderId = playerChargeRecord.OrderId
		}
	}

	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}

	firstChargePlayers, multiChargePlayers, _, firstChargeTotalMoney, multiChargeTotalMoney, _ :=
		models.GetChargeInfoBySpecifyDate(params.PlatformId, params.ServerId, params.ChannelList, params.StartTime, params.EndTime)

	data, total, playerCount, moneyCount, chargetTotalMoney, exchangeRate := models.GetChargeInfoRecordList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["playerCount"] = playerCount
	result["moneyCount"] = moneyCount
	result["chargeTotalMoney"] = chargetTotalMoney
	result["exchangeRate"] = exchangeRate
	result["firstChargePlayers"] = firstChargePlayers
	result["firstChargeMoney"] = firstChargeTotalMoney
	result["multiChargePlayers"] = multiChargePlayers
	result["multiChargeMoney"] = multiChargeTotalMoney

	if !hasNotPlatformOrderData {
		c.HttpResult(r, enums.CodeSuccess, "获取充值记录日志", result)
	}
	//添加第三方id
	var orderIds = make([]string, 0)
	if len(data) > 0 {
		for _, d := range data {
			orderIds = append(orderIds, d.OrderId)
		}
	}

	var platformOrderMap = make(map[string]string)
	if len(orderIds) > 0 {

		playerChargerRecords, err := models.GetItemByPlatformOrderId(subDb, orderIds)
		if err != nil {
			c.CheckError(err)
		}

		if len(playerChargerRecords) > 0 {
			for _, item := range playerChargerRecords {
				platformOrderMap[item.OrderId] = item.PlatformOrderId
			}
		}

	}

	if len(platformOrderMap) > 0 && len(data) > 0 {
		for key, item := range data {
			if _, ok := platformOrderMap[item.OrderId]; ok && platformOrderMap[item.OrderId] != "0" {
				data[key].PlatformOrderId = platformOrderMap[item.OrderId]
			}
		}
	}

	result["rows"] = data

	if len(data) > 0 {
		var playerIdArr []int
		for _, item := range data {
			playerIdArr = append(playerIdArr, item.PlayerId)
		}
		var sbind = make([]*models.ServerRel, 0)
		err = models.Db.Where("uid in (?)", playerIdArr).Find(&sbind).Error
		utils.CheckError(err)

		if len(sbind) > 0 {
			var mbind = map[int]int{}
			for _, bind := range sbind {
				mbind[bind.Uid] = bind.SAccount
			}

			for key, item := range data {
				data[key].SAccount = string(mbind[item.PlayerId])
			}
		}

	}

	c.HttpResult(r, enums.CodeSuccess, "获取充值记录日志", result)
}

// 下载充值数据
func (c *ChargeController) ChargeDownload(r *ghttp.Request) {
	var params models.ChargeInfoRecordQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if params.PlatformId == "" {
		c.CheckError(gerror.New("平台ID不能为空"))
	}
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data := models.GetChargeInfoDownload(&params)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "下载当前充值数据", result)
}

//获取充值排行榜
func (c *ChargeController) ChargeRankList(r *ghttp.Request) {
	var params models.PlayerChargeDataQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, total := models.GetPlayerChargeDataList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "查询充值排行", result)
}

// 充值任务分布
func (c *ChargeController) ChargeTaskDistribution(r *ghttp.Request) {
	var params models.ChargeTaskDistributionQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)
	data := models.GetChargeTaskDistribution(params)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取充充值任务分布", result)
}

// 充值活动分布
func (c *ChargeController) ChargeActivityDistribution(r *ghttp.Request) {
	var params models.ChargeActivityDistributionQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data := models.GetChargeActivityDistribution(params)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取充值活动分布", result)
}

// 充值金额分布
func (c *ChargeController) ChargeMoneyDistribution(r *ghttp.Request) {
	var params models.ChargeMoneyDistributionQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	//data := models.GetChargeMoneyDistribution(params)
	data := models.GetChargeMoneyDistributionV1(params)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取金额分布", result)
}

// 充值等级分布
func (c *ChargeController) ChargeLevelDistribution(r *ghttp.Request) {
	var params models.ChargeLevelDistributionQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	//data := models.GetChargeLevelDistribution(params)
	data := models.GetChargeLevelDistributionV1(params)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取等级分布", result)
}

// Ltv金额
func (c *ChargeController) LtvMoney(r *ghttp.Request) {
	var params models.RegRechargeParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	//data := models.GetDailyLTVList(params)
	data := models.GetLtvMoney(params)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取LTV金额", result)
}

// 每日LTV
func (c *ChargeController) GetDailyLTV(r *ghttp.Request) {
	var params models.DailyLTVQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	//data := models.GetDailyLTVList(params)
	data := models.GetDailyLTVListV1(params)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取每日LTV", result)
}

// 充值退款
func (c *ChargeController) ChargeRefundQueryParam(r *ghttp.Request) {
	var params ChargeRefundQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	models.GetChargeRefundResult(params.PlatformId, params.RefundType, params.RefundReason, params.OrderId, params.ApplyRefundMoney, params.PlatformOrderId, params.PlatformUserId)
	//result := make(map[string]interface{})
	//result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "充值退款", 0)
}

// 钻石排行榜
func (c *ChargeController) RankDiamond(r *ghttp.Request) {
	var req models.RankDiamondOrGoldCoinReq
	err := json.Unmarshal(r.GetBody(), &req)
	c.CheckError(err)

	data := models.RankDiamondOrGoldCoin(4, &req)
	c.HttpResult(r, enums.CodeSuccess, "ok", data)
}

// 金币排行榜
func (c *ChargeController) RankGoldCoin(r *ghttp.Request) {
	var req models.RankDiamondOrGoldCoinReq
	err := json.Unmarshal(r.GetBody(), &req)
	c.CheckError(err)

	data := models.RankDiamondOrGoldCoin(2, &req)
	c.HttpResult(r, enums.CodeSuccess, "ok", data)
}
