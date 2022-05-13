package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"strconv"
	"time"
)

type AppNoticeController struct {
	BaseController
}

func ChkParams(c *AppNoticeController) (*models.AppNoticeQueryParam, error) {
	params := &models.AppNoticeQueryParam{}
	offset, err := strconv.Atoi(c.Request.Get("offset").(string))
	if err != nil {
		utils.CheckError(err)
		return params, err
	}
	limit, limitErr := strconv.Atoi(c.Request.Get("limit").(string))
	if limitErr != nil {
		utils.CheckError(limitErr)
		return params, err
	}
	sort := c.Request.Get("sortColumn").(string)
	order := c.Request.Get("sortOrder").(string)

	params.Offset = offset
	params.Limit = limit
	params.Order = order
	params.Sort = sort

	return params, nil
}

// List 获取app公告
func (c *AppNoticeController) List(r *ghttp.Request) {
	params, paramErr := ChkParams(c)
	if paramErr != nil {
		utils.CheckError(paramErr)
		c.HttpResult(r, enums.CodeFail, "参数错误", paramErr)
	}
	g.Log().Infof("params: %+v", params)
	//err := json.Unmarshal(r.GetBody(), &params)
	//utils.CheckError(err)
	data, count := models.AppNoticeList(params)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取app公告成功", result)
}

// Create 创建app公告
func (c *AppNoticeController) Create(r *ghttp.Request) {
	params := models.AppNotice{}
	err := json.Unmarshal(r.GetBody(), &params)
	if err != nil {
		utils.CheckError(err)
		c.HttpResult(r, enums.CodeFail, "创建app公告参数错误", err)
	}

	params.CreatedBy = c.curUser.Id
	params.UpdatedBy = c.curUser.Id
	params.CreatedAt = int(time.Now().Unix())
	params.UpdatedAt = params.CreatedAt
	err = models.CreateAppNotice(params)
	if err != nil {
		utils.CheckError(err)
		c.HttpResult(r, enums.CodeFail, "app公告创建失败", err)
	}
	c.HttpResult(r, enums.CodeSuccess, "app公告创建成功", err)
}

// Update 编辑app公告
func (c *AppNoticeController) Update(r *ghttp.Request) {
	id, err := strconv.Atoi(c.Request.Get("id").(string))
	if err != nil {
		utils.CheckError(err)
		c.HttpResult(r, enums.CodeFail, "id参数有误", err)
	}

	params := models.AppNotice{}
	err = json.Unmarshal(r.GetBody(), &params)
	if err != nil {
		utils.CheckError(err)
		c.HttpResult(r, enums.CodeFail, "编辑app公告参数错误", err)
	}

	params.Id = id
	params.UpdatedBy = c.curUser.Id
	params.UpdatedAt = int(time.Now().Unix())

	err = models.UpdateAppNotice(params)
	if err != nil {
		utils.CheckError(err)
		c.HttpResult(r, enums.CodeFail, "app公告编辑失败", err)
	}
	c.HttpResult(r, enums.CodeSuccess, "app公告编辑成功", err)
}

// Delete 删除app公告
func (c *AppNoticeController) Delete(r *ghttp.Request) {
	id, err := strconv.Atoi(c.Request.Get("id").(string))
	if err != nil {
		utils.CheckError(err)
		c.HttpResult(r, enums.CodeFail, "id参数有误", err)
	}

	params := models.AppNotice{}
	params.Id = id

	err = models.DeleteAppNotice(params)
	if err != nil {
		utils.CheckError(err)
		c.HttpResult(r, enums.CodeFail, "app公告删除失败", err)
	}
	c.HttpResult(r, enums.CodeSuccess, "app公告删除成功", err)
}
