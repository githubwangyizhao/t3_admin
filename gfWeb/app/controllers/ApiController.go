package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type ApiController struct {
	BaseController
}

func (c *ApiController) GetAreaCodeList(r *ghttp.Request) {
	params := &models.RegionQueryParam{}
	params.Order = "id"
	params.Sort = "desc"
	params.Limit = 9999999
	params.Offset = 0
	data, _ := models.GetRegions(params)
	c.HttpResult(r, enums.CodeSuccess, "查询后台设置的货币单位、区号、国家/地区名称成功", data)
}

func (c *ApiController) GetAppNoticeList(r *ghttp.Request) {
	params := &models.AppNoticeQueryParam{}
	params.Order = "id"
	params.Sort = "desc"
	params.Limit = 9999999
	params.Offset = 0
	//err := json.Unmarshal(r.GetBody(), &params)
	//utils.CheckError(err)
	data, _ := models.AppNoticeList4Erlang(params)
	c.HttpResult(r, enums.CodeSuccess, "查询游戏服心跳验证规则成功", data)
}

func (c *ApiController) GetClientVerifyList(r *ghttp.Request) {
	params := &models.ClientHeartVerifyRequest{}
	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)
	data, _ := models.GetClientHeartVerifyList(params)
	c.HttpResult(r, enums.CodeSuccess, "查询游戏服心跳验证规则成功", data)
}

// GetPlayerInfo 玩家信息相关
func (c *ApiController) GetPlayerInfo(r *ghttp.Request) {
	var playerInfos models.PlayerInfos
	err := json.Unmarshal(r.GetBody(), &playerInfos)
	utils.CheckError(err)
	playerInfos = models.GetPlayerInfosFirst(playerInfos.PlatformId, playerInfos.ServerId, playerInfos.PlayerId)
	c.HttpResult(r, enums.CodeSuccess, "获取玩家信息", playerInfos)
}

// GetPlayerInfoList 玩家信息相关列表
func (c *ApiController) GetPlayerInfoList(r *ghttp.Request) {
	var playerInfos models.PlayerInfos
	err := json.Unmarshal(r.GetBody(), &playerInfos)
	utils.CheckError(err)
	playerInfosList := models.GetPlayerInfosList(playerInfos.PlatformId, playerInfos.ServerId)
	c.HttpResult(r, enums.CodeSuccess, "获取玩家信息列表", playerInfosList)
}

// GetPlatformInfo 平台信息相关
func (c *ApiController) GetPlatformInfo(r *ghttp.Request) {
	var platformClientInfo models.PlatformClientInfo
	err := json.Unmarshal(r.GetBody(), &platformClientInfo)
	utils.CheckError(err)
	data := models.GetPlatFormPayTimes(platformClientInfo.AppId, platformClientInfo)
	c.HttpResult(r, enums.CodeSuccess, "获取平台相关信息", data)
}

// GetAdjustList 获取修正值
func (c *ApiController) GetAdjustList(r *ghttp.Request) {
	var adjust models.Adjust
	err := json.Unmarshal(r.GetBody(), &adjust)
	utils.CheckError(err)
	data, _ := models.GetAdjustByPlatform(adjust.PlatformId, adjust.Server)
	c.HttpResult(r, enums.CodeSuccess, "获取修正值", data)

}

func (c *ApiController) GetTrackerInfoList(r *ghttp.Request) {
	var params struct {
		platformIdList []string
	}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("获取平台列表:%+v", params)
	//list := models.GetPlatformListByPlatformIdList(params.platformIdList)
	list := models.GetChannelList()
	g.Log().Infof("getTrackerInfoList: %+v", list)
	c.HttpResult(r, enums.CodeSuccess, "获取trackerInfoList", list)
}
