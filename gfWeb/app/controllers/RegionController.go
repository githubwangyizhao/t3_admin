package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/app/service"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"strconv"
)

type RegionController struct {
	BaseController
}

func ChkParams4Region(c *RegionController) (*models.RegionQueryParam, error) {
	params := &models.RegionQueryParam{}
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

func (c *RegionController) Regions(r *ghttp.Request) {
	params, paramErr := ChkParams4Region(c)
	if paramErr != nil {
		utils.CheckError(paramErr)
		c.HttpResult(r, enums.CodeFail, "参数错误", paramErr)
	}

	data, count := models.GetRegions(params)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data

	c.HttpResult(r, enums.CodeSuccess, "ok", result)
}

func (c *RegionController) Add(r *ghttp.Request) {
	params := &models.Region{}
	if err := json.Unmarshal(r.GetBody(), &params); err != nil {
		utils.CheckError(err)
		c.HttpResult(r, enums.CodeFail, "region add接口入参错误", nil)
	}
	params.CreatedBy = c.curUser.Id
	params.UpdatedBy = c.curUser.Id
	params.CreatedAt = int(gtime.Now().Timestamp())
	params.UpdatedAt = params.CreatedAt

	if err := service.AddRegion(params); err != nil {
		c.HttpResult(r, enums.CodeFail, "region数据修改失败", nil)
	}

	c.HttpResult(r, enums.CodeSuccess, "ok", nil)
}

func (c *RegionController) Delete(r *ghttp.Request) {
	id, err := strconv.Atoi(c.Request.Get("id").(string))
	if err != nil {
		utils.CheckError(err)
		c.HttpResult(r, enums.CodeFail, "id参数有误", err)
	}

	params := &models.Region{}
	params.Id = id

	if err := service.DeleteRegion(params); err != nil {
		utils.CheckError(err)
		c.HttpResult(r, enums.CodeFail, "region数据删除失败", err)
	}

	c.HttpResult(r, enums.CodeSuccess, "region数据删除成功", err)
}

func (c *RegionController) Edit(r *ghttp.Request) {
	id, err := strconv.Atoi(c.Request.Get("id").(string))
	if err != nil {
		utils.CheckError(err)
		c.HttpResult(r, enums.CodeFail, "id参数有误", err)
	}

	params := &models.Region{}
	if err := json.Unmarshal(r.GetBody(), &params); err != nil {
		c.CheckError(err, "region edit接口入参错误")
	}
	params.Id = id
	params.UpdatedBy = c.curUser.Id
	params.UpdatedAt = int(gtime.Now().Timestamp())

	if err := service.EditRegion(params); err != nil {
		utils.CheckError(err)
		c.HttpResult(r, enums.CodeFail, "region edit失败", nil)
	}

	c.HttpResult(r, enums.CodeSuccess, "ok", nil)
}
