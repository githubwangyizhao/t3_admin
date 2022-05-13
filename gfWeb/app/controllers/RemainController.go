package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type RemainController struct {
	BaseController
}

// 总体留存
func (c *RemainController) GetTotalRemain(r *ghttp.Request) {
	var params models.TotalRemainQueryParam
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("查询总体留存:+", params)
	data, total := models.GetRemainTotalList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "查询总体留存成功", result)
}

// 付费留存
func (c *RemainController) GetChargeRemain(r *ghttp.Request) {
	var params models.RemainChargeQueryParam
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("查询付费留存:+", params)
	data := models.GetRemainChargeList(&params)
	result := make(map[string]interface{})
	result["total"] = 100
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "查询付费留存成功", result)
}

// 活跃留存
func (c *RemainController) GetActiveRemain(r *ghttp.Request) {
	var params models.ActiveRemainQueryParam
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("查询活跃留存:+", params)
	data, total := models.GetRemainActiveList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "查询活跃留存成功", result)
}

// 任务留存
func (c *RemainController) GetTaskRemain(r *ghttp.Request) {
	var params struct {
		PlatformId     string   `json:"platformId"`
		ServerId       string   `json:"serverId"`
		ChannelList    []string `json:"channelList"`
		IsChargePlayer int      `json:"isChargePlayer"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	g.Log().Infof("查询任务留存:%+v", params)
	utils.CheckError(err)
	data := models.GetRemainTask(params.PlatformId, params.ServerId, params.ChannelList, params.IsChargePlayer)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "查询任务留存成功", result)
}

// 等级留存
func (c *RemainController) GetLevelRemain(r *ghttp.Request) {
	var params struct {
		PlatformId     string   `json:"platformId"`
		ServerId       string   `json:"serverId"`
		ChannelList    []string `json:"channelList"`
		StartTime      int      `json:"startTime"`
		EndTime        int      `json:"endTime"`
		IsChargePlayer int      `json:"isChargePlayer"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("查询等级留存:%+v", params)
	data := models.GetRemainLevel(params.PlatformId, params.ServerId, params.ChannelList, params.StartTime, params.EndTime, params.IsChargePlayer)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "查询等级留存成功", result)
}

//时长留存
func (c *RemainController) GetTimeRemain(r *ghttp.Request) {
	var params struct {
		PlatformId     string   `json:"platformId"`
		ServerId       string   `json:"serverId"`
		ChannelList    []string `json:"channelList"`
		StartTime      int      `json:"startTime"`
		EndTime        int      `json:"endTime"`
		IsChargePlayer int      `json:"isChargePlayer"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("查询时长留存:%+v", params)
	data := models.GetRemainTime(params.PlatformId, params.ServerId, params.ChannelList, params.StartTime, params.EndTime, params.IsChargePlayer)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "查询时长留存成功", result)
}
