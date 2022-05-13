package controllers

import (
	"encoding/json"
	"fmt"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type MergeController struct {
	BaseController
}

// 平台合服
func (c *MergeController) PlatformMerge(r *ghttp.Request) {
	var params struct {
		PlatformId string                   `json:"platformId"`
		MergeTime  int                      `json:"mergeTime"`
		MergeList  []models.MergeServerData `json:"mergeList"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	g.Log().Infof("请求合服%s:%+v", c.curUser.Name, params)

	err = models.CreatePlatformMergeData(c.curUser.Id, c.curUser.Name, params.PlatformId, params.MergeTime, params.MergeList)
	if err != nil {
		c.HttpResult(r, enums.CodeFail2, fmt.Sprintf("%+v", err), "")
	} else {
		c.HttpResult(r, enums.CodeSuccess, "合服成功!", "")
	}
}

// 获得平台合服信息
func (c *MergeController) GetPlatformMerge(r *ghttp.Request) {
	var params models.MergeServerParam
	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	g.Log().Infof("获得平台合服信息%s:%+v", c.curUser.Name, params)
	data, total := models.GetPlatformMergeServerList(params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	result["isAuth"] = c.isMergeAuth()
	c.HttpResult(r, enums.CodeSuccess, "查询成功", result)
}

// 审核通过
func (c *MergeController) AuditPlatformMerge(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:"platformId"`
		MergeTime  int    `json:"mergeTime"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if c.isMergeAuth() == false {
		c.CheckError(gerror.New("没有权限操作"))
	}
	g.Log().Infof("审核通过%s:%+v", c.curUser.Name, params)
	err = models.AuditPlatformMergeData(c.curUser.Id, params.PlatformId, params.MergeTime)
	if err != nil {
		c.HttpResult(r, enums.CodeFail2, fmt.Sprint("%+v:", err), "")
	} else {
		c.HttpResult(r, enums.CodeSuccess, "", "")
	}
}

// 删除合服平台
func (c *MergeController) DelPlatformMerge(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:"platformId"`
		MergeTime  int    `json:"mergeTime"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if c.isMergeAuth() == false {
		c.CheckError(gerror.New("没有权限操作"))
	}
	g.Log().Infof("删除合服平台%s:%+v", c.curUser.Name, params)
	err = models.DelPlatformMergeData(params.PlatformId, params.MergeTime)
	if err != nil {
		c.HttpResult(r, enums.CodeFail2, fmt.Sprint("删除合服平台失败:%+v:", err), "")
	} else {
		c.HttpResult(r, enums.CodeSuccess, "删除合服平台成功!", "")
	}
}

// 合服审核权限
func (c *MergeController) isMergeAuth() bool {
	if c.curUser.IsSuperUser() {
		return true
	}
	return models.IsPageChangeAuthDefault(c.curUser, models.PAGE_CHANGE_AUTH_IS_MERGE, false)
	//return true
	//url := r.URL.Path
	//return strings.Index(url, "audit_platform_merge") != -1
}
