package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type ClientHeartVerifyController struct {
	BaseController
}

func (c *ClientHeartVerifyController) Edit(r *ghttp.Request) {
	params := &models.ClientHeartVerifyRequest{}
	g.Log().Infof("创建、编辑游戏服心跳验证规则:%+v", params)

	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)

	result := models.CreOrModifyData(params)
	utils.CheckError(err)

	c.HttpResult(r, enums.CodeSuccess, "查询游戏服心跳验证规则成功", result)
}

func (c *ClientHeartVerifyController) List(r *ghttp.Request) {
	params := &models.ClientHeartVerifyRequest{}
	g.Log().Infof("查询游戏服心跳验证规则:%+v", params)

	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)

	//data, count := models.GetClientHeartVerifyList(params.Platform, params.Server, params.Status)
	data, count := models.GetClientHeartVerifyList(params)

	for _, e := range data {
		e.CreatedName = models.GetUserName(e.CreatedBy)
		e.UpdatedName = models.GetUserName(e.UpdatedBy)
	}

	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data

	c.HttpResult(r, enums.CodeSuccess, "查询游戏服心跳验证规则成功", result)
}
