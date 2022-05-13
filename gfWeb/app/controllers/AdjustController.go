package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type AdjustController struct {
	BaseController
}

func (c *AdjustController) EditAdjust(r *ghttp.Request) {
	params := &models.Adjust{}
	g.Log().Infof("编辑修正值:%+v", params)
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	utils.CheckError(err, "编辑修正值")

	params.UpdatedBy = c.curUser.Id

	err = models.EditAdjustData(params)
	if err != nil {
		c.CheckError(err, "编辑修正值错误")
	}
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

func (c *AdjustController) CreateAdjust(r *ghttp.Request) {
	params := &models.Adjust{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)

	g.Log().Infof("创建修正值:%+v", params)
	utils.CheckError(err, "创建修正值")

	params.CreatedBy = c.curUser.Id
	params.UpdatedBy = c.curUser.Id
	err = models.CreateAdjustData(params)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

// SelectAdjust 查询场景修正值
func (c *AdjustController) SelectAdjust(r *ghttp.Request) {
	params := &models.AdjustRequest{}
	g.Log().Infof("查询场景修正值:%+v", params)

	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)

	data, count := models.GetAdjustDataList(params)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data

	c.HttpResult(r, enums.CodeSuccess, "查询修正值成功", result)
}
