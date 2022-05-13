package controllers

import (
	"encoding/json"
	"fmt"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"

	//"strconv"
	"strconv"
	//"encoding/base64"
	//"net/http"
	//"io/ioutil"
)

type GameServerController struct {
	BaseController
}

// 获取游戏服列表
func (c *GameServerController) List(r *ghttp.Request) {
	var params models.GameServerQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	g.Log().Debug("查询游戏服列表:%+v", params)
	utils.CheckError(err)
	data, total := models.GetGameServerList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取游戏服列表成功", result)
}

// 添加 编辑 游戏服
func (c *GameServerController) Edit(r *ghttp.Request) {
	params := models.GameServer{}
	err := json.Unmarshal(r.GetBody(), &params)
	m := params
	c.CheckError(err, "编辑游戏服")
	if m.OpenTime == 0 {
		c.HttpResult(r, enums.CodeFail, "开服时间不能为0", 0)
	}
	if m.IsAdd == 1 && models.IsGameServerExists(m.PlatformId, m.Sid) {
		c.HttpResult(r, enums.CodeFail, "游戏服已经存在", m.Node)
	}
	if m.IsAdd == 0 && models.IsGameServerExists(m.PlatformId, m.Sid) == false {
		c.HttpResult(r, enums.CodeFail, "游戏服不存在", m.Node)
	}

	out, err := models.AddGameServer(m.PlatformId, m.Sid, m.Desc, m.Node, m.ZoneNode, m.State, m.OpenTime, m.IsShow)
	//out, err := utils.NodeTool(
	//	"mod_server_mgr",
	//	"add_game_server",
	//	m.PlatformId,
	//	m.Sid,
	//	m.Desc,
	//	m.Node,
	//	m.ZoneNode,
	//	strconv.Itoa(m.State),
	//	strconv.Itoa(m.OpenTime),
	//	strconv.Itoa(m.IsShow),
	//)

	c.CheckError(err, "保存游戏服失败:"+out)
	c.HttpResult(r, enums.CodeSuccess, "保存成功", m.Sid)
}

// 删除游戏服
func (c *GameServerController) Delete(r *ghttp.Request) {
	var params []struct {
		PlatformId string `json:platformId`
		ServerId   string `json:serverId`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	ids := params
	c.CheckError(err)
	for _, id := range ids {
		out, err := utils.CenterNodeTool(
			"mod_server_mgr",
			"delete_game_server",
			id.PlatformId,
			id.ServerId,
		)
		c.CheckError(err, "删除游戏服失败:"+out)
		g.Log().Info("删除游戏服:%s", id)
	}

	c.HttpResult(r, enums.CodeSuccess, fmt.Sprintf("成功删除 %d 项", len(ids)), 0)
}

// 批量修改区服状态
func (c *GameServerController) BatchUpdateState(r *ghttp.Request) {
	var params struct {
		PlatformId string   `json:platformId`
		Nodes      []string `json:node`
		State      int
	}
	err := json.Unmarshal(r.GetBody(), &params)
	param := params
	c.CheckError(err)
	if len(param.Nodes) == 0 {
		err = models.BatchUpdateState(param.PlatformId, param.State)
		if err != nil {
			c.HttpResult(r, enums.CodeFail, err.Error(), 0)
		}
		//if param.PlatformId == "" {
		//	g.Log().Error("平台id不能为空")
		//	c.HttpResult(r, enums.CodeFail, "平台id不能为空", 0)
		//}
		//out, err := utils.CenterNodeTool(
		//	"mod_server_mgr",
		//	"update_all_game_server_state",
		//	param.PlatformId,
		//	strconv.Itoa(param.State),
		//)
		//c.CheckError(err, "修改所有区服状态:"+out)
		//models.UpdateEtsPlatformServerClose(param.PlatformId , param.State == 1)
	} else {
		for _, node := range param.Nodes {
			out, err := utils.CenterNodeTool(
				"mod_server_mgr",
				"update_node_state",
				node,
				strconv.Itoa(param.State),
			)
			c.CheckError(err, "批量修改区服状态:"+out)
		}
	}

	c.HttpResult(r, enums.CodeSuccess, fmt.Sprintf("批量修改区服状态 %d 项", len(param.Nodes)), 0)
}

// 刷新区服入口
func (c *GameServerController) Refresh(r *ghttp.Request) {
	//var params struct {
	//}
	////var result struct {
	////	ErrorCode int
	////}
	//err := json.Unmarshal(r.GetBody(),&params)
	//c.CheckError(err)
	err := models.RefreshGameServer()
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "刷新区服入口成功", 0)

}

// 开服
func (c *GameServerController) OpenServer(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:platformId`
		Time       int    `json:time`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	//openServerTime := 0
	msg := "立即开服"
	if params.Time == 0 {
		// 立即开服
		//openServerTime = utils.GetTimestamp()
		err = models.OpenServerType(c.curUser.Id, params.PlatformId, 0, 0)
		c.CheckError(err, msg+"失败")
		err = models.SaveOpenServerManageLogData(c.curUser.Id, params.PlatformId, "立即开服")
		c.CheckError(err)
	} else {
		//定时开服
		//openServerTime = params.Time
		err = models.UpdatePlatformOpenServerTime(params.PlatformId, params.Time)
		msg = "定时开服设置"
		c.CheckError(err, msg+"失败")
	}
	c.HttpResult(r, enums.CodeSuccess, msg+"成功", 0)
}
