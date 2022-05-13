package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/app/service"
	"gfWeb/library/enums"
	"gfWeb/library/utils"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	//"errors"
)

type StatisticsController struct {
	BaseController
}

// 每日汇总
func (c *StatisticsController) DailyStatisticsList(r *ghttp.Request) {
	var params models.DailyStatisticsQueryParam
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("获取每日汇总:%+v", params)
	data, payOrder := models.GetDailyStatisticsList(&params)
	result := make(map[string]interface{})
	//result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取每日汇总", g.Map{"result": result, "order": payOrder})
}

////在线统计
//func (c *StatisticsController) OnlineStatisticsList() {
//	var params models.DailyOnlineStatisticsQueryParam
//	err := json.Unmarshal(r.GetBody(),&params)
//
//	utils.CheckError(err)
//	g.Log().Info("获取在线统计:%+v", params)
//	if params.PlatformId == "" {
//		c.CheckError(errors.New("平台ID不能为空"))
//	}
//	data := models.GetDailyOnlineStatisticsList(&params)
//	result := make(map[string]interface{})
//	//result["total"] = total
//	result["rows"] = data
//	c.HttpResult(r, enums.CodeSuccess, "获取在线统计", result)
//}
//注册统计
//func (c *StatisticsController) RegisterStatisticsList() {
//	var params models.DailyRegisterStatisticsQueryParam
//	err := json.Unmarshal(r.GetBody(),&params)
//
//	utils.CheckError(err)
//	g.Log().Info("获取注册统计:%+v", params)
//	if params.PlatformId == "" {
//		c.CheckError(errors.New("平台ID不能为空"))
//	}
//	data:= models.GetDailyRegisterStatisticsList(&params)
//	result := make(map[string]interface{})
//	//result["total"] = total
//	result["rows"] = data
//	c.HttpResult(r, enums.CodeSuccess, "获取注册统计", result)
//}

////注册统计
//func (c *StatisticsController) ActiveStatisticsList() {
//	var params models.DailyActiveStatisticsQueryParam
//	err := json.Unmarshal(r.GetBody(),&params)
//
//	utils.CheckError(err)
//	g.Log().Info("获取活跃统计:%+v", params)
//	if params.PlatformId == "" {
//		c.CheckError(errors.New("平台ID不能为空"))
//	}
//	data := models.GetDailyActiveStatisticsList(&params)
//	result := make(map[string]interface{})
//	//result["total"] = total
//	result["rows"] = data
//	c.HttpResult(r, enums.CodeSuccess, "获取活跃统计", result)
//}

//消费分析
func (c *StatisticsController) ConsumeAnalysis(r *ghttp.Request) {
	var params models.PropConsumeStatisticsQueryParam
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("获取消费分析统计:%+v", params)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, err := models.GetPropConsumeStatistics(&params)
	c.CheckError(err)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "消费分析统计", result)
}

//服务器概况
func (c *StatisticsController) GetServerGeneralize(r *ghttp.Request) {
	var params models.ServerGeneralizeQueryParam
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("查询服务器概况:%+v", params)
	params.EndTime = params.EndTime / 1000
	data, err := models.GetServerGeneralize(params.EndTime, params.PlatformId, params.ServerId, params.ChannelList)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "获取服务器概况", data)
}

func (c *StatisticsController) GetRealTimeOnline(r *ghttp.Request) {
	//g.Log().Info("GetRealTimeOnline")
	var params struct {
		PlatformId  string   `json:"platformId"`
		ServerId    string   `json:"serverId"`
		ChannelList []string `json:"channelList"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("查询实时在线统计:%+v", params)
	//platformId := c.GetString("platformId")
	//serverId := c.GetString("serverId")
	serverOnlineStatistics, err := models.GetServerOnlineStatistics(params.PlatformId, params.ServerId, params.ChannelList)
	c.CheckError(err, "查询实时在线统计")
	//g.Log().Info("查询实时在线统计成功:%+v", serverOnlineStatistics)
	c.HttpResult(r, enums.CodeSuccess, "查询实时在线统计成功", serverOnlineStatistics)
}

func (c *StatisticsController) GetChargeStatistics(r *ghttp.Request) {
	//g.Log().Info("GetRealTimeOnline")
	var params struct {
		PlatformId  string   `json:"platformId"`
		ServerId    string   `json:"serverId"`
		ChannelList []string `json:"channelList"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("查询充值对比:%+v", params)
	//platformId := c.GetString("platformId")
	//serverId := c.GetString("serverId")
	serverOnlineStatistics, err := models.GetChargeStatistics(params.PlatformId, params.ServerId, params.ChannelList)
	c.CheckError(err, "查询充值对比")
	//g.Log().Info("查询实时在线统计成功:%+v", serverOnlineStatistics)
	c.HttpResult(r, enums.CodeSuccess, "查询充值对比成功", serverOnlineStatistics)
}

func (c *StatisticsController) GetIncomeStatistics(r *ghttp.Request) {
	//g.Log().Info("GetRealTimeOnline")
	var params struct {
		PlatformId  string   `json:"platformId"`
		ServerId    string   `json:"serverId"`
		ChannelList []string `json:"channelList"`
		StartTime   int
		EndTime     int
	}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("查询总流水:%+v", params)
	//platformId := c.GetString("platformId")
	//serverId := c.GetString("serverId")
	serverOnlineStatistics := models.GetIncomeStatisticsChartData(params.PlatformId, params.ServerId, params.ChannelList, params.StartTime, params.EndTime)
	c.CheckError(err, "查询总流水")
	//g.Log().Info("查询实时在线统计成功:%+v", serverOnlineStatistics)
	c.HttpResult(r, enums.CodeSuccess, "查询总流水", serverOnlineStatistics)
}

// 获得平台订阅数据
func (c *StatisticsController) GetPlatformDingYue(r *ghttp.Request) {
	var params struct {
		models.BaseQueryParam
		PlatformId string `json:"platformId"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	g.Log().Infof("获得平台订阅数据%s:%+v", c.curUser.Name, params)
	data, dingYueCount, total := models.GetPlatformDingYue(params.PlatformId, params.BaseQueryParam)
	g.Log().Debugf("data:", data, " dingYueCount:", dingYueCount, " total:", total)
	result := make(map[string]interface{})
	result["total"] = total
	result["dingYueCount"] = dingYueCount
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "查询平台订阅数据成功", result)
}

func (c *StatisticsController) OauthList(r *ghttp.Request) {
	var params models.OauthInfoRecordQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	g.Log().Info(params)
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
	data, total, playerCount, moneyCount, rcAmount := models.GetOauthInfoRecordList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["playerCount"] = playerCount
	result["moneyCount"] = moneyCount
	result["isExceed"] = 0 //超过充值

	//是否超过充值金额
	if moneyCount > rcAmount {
		result["isExceed"] = 1
	}

	c.HttpResult(r, enums.CodeSuccess, "获取提现列表", result)
}

func (c *StatisticsController) GetRubyStatistics(r *ghttp.Request) {

	params := models.PlayerPropLogQueryParam{}
	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)
	data, count, err := service.PlayerPropLogFromGameDb(params)
	if err != nil {
		c.CheckError(err, "红宝石记录获取失败")
	}

	g.Log().Infof("count: %d", count)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data

	c.HttpResult(r, enums.CodeSuccess, "红宝石记录获取成功", result)
}

//// 获得提现信息
//func (c *StatisticsController) GetOauthOrder(r *ghttp.Request) {
//	var params struct {
//		models.BaseQueryParam
//		PlatformId string `json:"platformId"`
//	}
//	err := json.Unmarshal(r.GetBody(), &params)
//
//	c.CheckError(err)
//	g.Log().Infof("获得提现数据%s:%+v", c.curUser.Name, params)
//	data, dingYueCount, total := models.GetPlatformDingYue(params.PlatformId, params.BaseQueryParam)
//	g.Log().Debugf("data:", data, " dingYueCount:", dingYueCount, " total:", total)
//	result := make(map[string]interface{})
//	result["total"] = total
//	result["dingYueCount"] = dingYueCount
//	result["rows"] = data
//	c.HttpResult(r, enums.CodeSuccess, "获得提现数据成功", result)
//}
