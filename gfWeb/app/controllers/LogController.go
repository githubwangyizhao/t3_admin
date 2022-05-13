package controllers

import (
	"encoding/json"
	"fmt"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"gfWeb/memdb/uselog"
	"time"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
)

type LogController struct {
	BaseController
}

func (c *LogController) PlayerLoinLogList(r *ghttp.Request) {
	var params models.PlayerLoginLogQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetPlayerLoginLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取玩家登录日志", result)
}

func (c *LogController) PlayerOnlineLogList(r *ghttp.Request) {
	var params models.PlayerOnlineLogQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetPlayerOnlineLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取在线日志", result)
}

func (c *LogController) PlayerChallengeMissionLogList(r *ghttp.Request) {
	var params models.PlayerChallengeMissionLogQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetPlayerMissionLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取副本挑战日志", result)
}

func (c *LogController) PlayerPropLogList(r *ghttp.Request) {
	var params models.PlayerPropLogQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetPlayerPropLogList2(&params)
	//models.GetPlayerPropLogList2(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取道具日志", result)
}

// 查询邮件日志
func (c *LogController) PlayerMailLogList(r *ghttp.Request) {
	var params models.PlayerMailLogQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetPlayerMailLogList(&params)
	//data, total := models.GetPlayerMailLogList2(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取邮件日志", result)
}

// 查询属性日志
func (c *LogController) PlayerAttrLogList(r *ghttp.Request) {
	var params models.PlayerAttrLogQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetPlayerAttrLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取属性日志", result)
}

// 查询活动奖励日志
func (c *LogController) ActivityAwardLogList(r *ghttp.Request) {
	var params models.ActivityAwardLogQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetActivityAwardLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取活动奖励日志", result)
}

// 查询冲榜排名
func (c *LogController) ImpactRankList(r *ghttp.Request) {
	var params models.ImpactRankQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetImpactRankList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取查询冲榜排名", result)
}

//func (c *LogController) ChargeList() {
//	var params models.ChargeInfoRecordQueryParam
//	err := json.Unmarshal(r.GetBody(),&params)
//
//	utils.CheckError(err)
//	g.Log().Info("查询充值记录日志:%+v", params)
//	if params.PlayerName != "" {
//		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
//		if player == nil || err != nil {
//			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
//		}
//		params.PlayerId = player.Id
//	}
//	data, total := models.GetChargeInfoRecordList(&params)
//	result := make(map[string]interface{})
//	result["total"] = total
//	result["rows"] = data
//	c.HttpResult(r, enums.CodeSuccess, "获取充值记录日志", result)
//}
//
//

// 请求日志列表
func (c *LogController) RequestLogList(r *ghttp.Request) {
	var params models.ParamsRequestAdmin
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, total := models.GetRequestAdminLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "请求日志列表", result)
}

// 登录日志列表
func (c *LogController) LoginLogList(r *ghttp.Request) {
	var params models.ParamsLoginAdmin
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, total := models.GetLoginAdminLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "登录日志列表", result)
}

// 开服日志列表
func (c *LogController) OpenServerManageLogList(r *ghttp.Request) {
	var params models.ParamsOpenServerManage
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, total := models.GetOpenServerManageLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "开服日志列表", result)
}

// 平台版本更新日志列表
func (c *LogController) UpdatePlatformVersionLogList(r *ghttp.Request) {
	var params models.ParamsUpdatePlatformVersion
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, total := models.GetUpdatePlatformVersionLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "平台版本更新日志列表", result)
}

// 查询玩家游戏场景日志
func (c *LogController) PlayerGameSceneLogList(r *ghttp.Request) {
	var params models.PlayerGameSceneLogQueryParam
	result := make(map[string]interface{})
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data := models.GetPlayerGameSceneLogList(&params)
	result["rows"] = data
	// else {
	// 	data, total := models.GetAllGameSceneLogList(&params)
	// 	result["total"] = total
	// 	result["rows"] = data
	// }
	c.HttpResult(r, enums.CodeSuccess, "获取玩家游戏场景日志", result)
}
func (c *LogController) GameRoundBoxLog(r *ghttp.Request) {
	var params models.RoundBoxReq
	err := json.Unmarshal(r.GetBody(), &params)
	if err != nil {
		utils.CheckError(err)
		return
	}

	if params.SceneId == 0 {
		c.HttpResult(r, enums.CodeFail, "获取转轮宝箱日志失败", nil)
		return
	}

	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	g.Log().Infof("Params: %+v", params)

	gameDb, err := models.GetGameDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	if gameDb == nil {
		defer gameDb.Close()
		return
	}
	var list = make([]models.ConsumeStatistics, 0)
	whereParams := fmt.Sprintf("scene_id=%d and log_type=%d", params.SceneId, params.LogType)
	if params.PlayerId > 0 {
		whereParams = whereParams + fmt.Sprintf(" and player_id=%d", params.PlayerId)
	}
	sql := fmt.Sprintf(`select * from consume_statistics where %s; `, whereParams)

	fmt.Println("sql ", sql)
	err = gameDb.Raw(sql).Scan(&list).Error
	utils.CheckError(err)

	c.HttpResult(r, enums.CodeSuccess, "获取转轮宝箱日志成功", list)
}

// ItemEventLogList 查询怪物击杀日志
func (c *LogController) GameMonsterLogList(r *ghttp.Request) {
	var params models.ItemEventLogQueryParam
	// result := make(map[string]interface{})
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	g.Log().Infof("params: %+v", params)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	g.Log().Infof("Params: %+v", params)

	t := time.Unix(int64(params.Datetime), 0)
	timeStr := fmt.Sprintf("%d_%d_%d", t.Year(), t.Month(), t.Day())

	var ret = make([]*uselog.MemItemEventLog, 0)
	var dataThisHour = map[string]map[string]*uselog.MemItemEventLog{}
	nowHour := time.Now().Hour()
	var data = map[string]map[string]*uselog.MemItemEventLog{}

	if params.PlayerId == 0 {
		data = uselog.FechTodayMonsterDataByPS(params.PlatformId, params.ServerId, gconv.String(params.MonsterId), params.Datetime)
		dataThisHour = models.GetItemGameLogList(&models.ItemEventLogQueryParam{
			PlatformId: params.PlatformId,
			ServerId:   params.ServerId,
			MonsterId:  params.MonsterId,
			Datetime:   params.Datetime,
		}, nowHour, models.LogMonsterType)

		for ts, _ := range data {
			// if timeStr == reqIim
			if data[ts] != nil {
				data[ts][gconv.String(nowHour)] = dataThisHour[ts][gconv.String(nowHour)]
			} else {
				data[ts] = dataThisHour[ts]
			}
		}

	} else {
		// 从datetime开始计算0-23，并作为GetItemEventLogList方法的hour参数传入
		for i := params.Datetime; i <= params.Datetime+86399; i = i + 3600 {
			t := time.Unix(int64(i), 0)
			specifyHourData := models.GetItemGameLogList(&models.ItemEventLogQueryParam{
				PlatformId: params.PlatformId,
				ServerId:   params.ServerId,
				MonsterId:  params.MonsterId,
				Datetime:   params.Datetime,
				PlayerId:   params.PlayerId,
			}, t.Hour(), models.LogMonsterType)

			if data[timeStr] == nil {
				data[timeStr] = map[string]*uselog.MemItemEventLog{}
			}

			var scanHourDuring = gconv.String(t.Hour())
			if len(scanHourDuring) == 1 {
				scanHourDuring = "0" + scanHourDuring
			}
			data[timeStr][scanHourDuring] = specifyHourData[timeStr][scanHourDuring]
		}
	}

	for ts, _ := range data {
		// if timeStr == reqIim
		if data[ts] != nil {
			data[ts][gconv.String(nowHour)] = dataThisHour[ts][gconv.String(nowHour)]
		} else {
			data[ts] = dataThisHour[ts]
		}
	}

	for _, item := range data {
		for _, value := range item {
			if value != nil {
				value.MonsterId = params.MonsterId
				ret = append(ret, value)
			}
		}
	}

	c.HttpResult(r, enums.CodeSuccess, "获取怪物日志成功", ret)
}

// ItemEventLogList 查询事件场景日志
func (c *LogController) ItemEventLogList(r *ghttp.Request) {
	var params models.ItemEventLogQueryParam
	// result := make(map[string]interface{})
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	g.Log().Infof("params: %+v", params)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	g.Log().Infof("Params: %+v", params)

	t := time.Unix(int64(params.Datetime), 0)
	timeStr := fmt.Sprintf("%d_%d_%d", t.Year(), t.Month(), t.Day())

	var ret = make([]*uselog.MemItemEventLog, 0)
	var dataThisHour = map[string]map[string]*uselog.MemItemEventLog{}
	nowHour := time.Now().Hour()
	// data := uselog.FetchALLDataByPS("local", "s1", "13", "1001")
	// data map[string]map[string]*MemItemEventLog
	// data := uselog.FetchALLDataByPS("local", "s1", "13", "1001")
	var data = map[string]map[string]*uselog.MemItemEventLog{}

	if params.PlayerId == 0 {
		data = uselog.FechTodayEventDataByPS(params.PlatformId, params.ServerId, gconv.String(params.LogType), gconv.String(params.Type), params.Datetime)
		dataThisHour = models.GetItemGameLogList(&models.ItemEventLogQueryParam{
			PlatformId: params.PlatformId,
			ServerId:   params.ServerId,
			LogType:    params.LogType,
			Type:       params.Type,
			Datetime:   params.Datetime,
		}, nowHour, models.LogEventType)

		for ts, _ := range data {
			// if timeStr == reqIim
			if data[ts] != nil {
				data[ts][gconv.String(nowHour)] = dataThisHour[ts][gconv.String(nowHour)]
			} else {
				data[ts] = dataThisHour[ts]
			}
		}

	} else {
		// 从datetime开始计算0-23，并作为GetItemEventLogList方法的hour参数传入
		for i := params.Datetime; i <= params.Datetime+86399; i = i + 3600 {
			t := time.Unix(int64(i), 0)
			tHour := t.Hour()

			if tHour == nowHour {
				fmt.Println("aaa")
			}
			specifyHourData := models.GetItemGameLogList(&models.ItemEventLogQueryParam{
				PlatformId: params.PlatformId,
				ServerId:   params.ServerId,
				LogType:    params.LogType,
				Type:       params.Type,
				Datetime:   params.Datetime,
				PlayerId:   params.PlayerId,
			}, tHour, models.LogEventType)

			if data[timeStr] == nil {
				data[timeStr] = map[string]*uselog.MemItemEventLog{}
			}

			var scanHourDuring = gconv.String(t.Hour())
			if len(scanHourDuring) == 1 {
				scanHourDuring = "0" + scanHourDuring
			}
			data[timeStr][scanHourDuring] = specifyHourData[timeStr][scanHourDuring]
		}
	}

	for _, item := range data {
		for _, value := range item {
			if value != nil {
				ret = append(ret, value)
			}
		}
	}

	c.HttpResult(r, enums.CodeSuccess, "获取事件日志成功", ret)
}

// 查询玩家游戏场景日志
func (c *LogController) ClientLogList(r *ghttp.Request) {
	var params models.PlayerClientLogQueryParam
	result := make(map[string]interface{})
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data := models.ClientLogList(&params)
	result["rows"] = data
	// else {
	// 	data, total := models.GetAllGameSceneLogList(&params)
	// 	result["total"] = total
	// 	result["rows"] = data
	// }
	c.HttpResult(r, enums.CodeSuccess, "获取玩家客户端日志", result)
}

// 查询游戏物品操作相关日志
func (c *LogController) GotOptLogjson(r *ghttp.Request) {
	var params models.OptLogjsonReq
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)

	var (
		itemResponseStruct       = make([]*models.ItemStruct, 0)
		SceneResponseStruct      = make([]*models.SceneStruct, 0)
		ServiceLogResponseStruct = make([]*models.ServiceLogStruct, 0)
	)

	ResponseByte, itemRrr := models.GetJson(params.JsonName)
	if itemRrr != nil {
		utils.CheckError(itemRrr)
		c.HttpResult(r, enums.CodeFail, "物品 日志", itemResponseStruct)
		return
	}

	switch params.JsonName {
	case "item":
		jsonErr := json.Unmarshal(ResponseByte, &itemResponseStruct)
		if jsonErr != nil {
			utils.CheckError(jsonErr)
		}
		c.HttpResult(r, enums.CodeSuccess, "物品 日志", itemResponseStruct)
	case "scene":
		jsonErr := json.Unmarshal(ResponseByte, &SceneResponseStruct)
		if jsonErr != nil {
			utils.CheckError(jsonErr)
		}
		c.HttpResult(r, enums.CodeSuccess, "场景 日志", SceneResponseStruct)
	case "serviceLog":
		jsonErr := json.Unmarshal(ResponseByte, &ServiceLogResponseStruct)
		if jsonErr != nil {
			utils.CheckError(jsonErr)
		}
		c.HttpResult(r, enums.CodeSuccess, "事件 日志", ServiceLogResponseStruct)
	}

	c.HttpResult(r, enums.CodeSuccess, "查询游戏物品操作相关日志", nil)
}
