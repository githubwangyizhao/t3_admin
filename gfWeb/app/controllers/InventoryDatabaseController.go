package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"github.com/gogf/gf/net/ghttp"
)

type InventoryDatabaseController struct {
	BaseController
}

func (c *InventoryDatabaseController) DatabaseList(r *ghttp.Request) {
	var params models.InventoryDatabaseParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, total := models.GetInventoryDatabaseList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取资产列表成功", result)
}
func (c *InventoryDatabaseController) AllDatabaseList(r *ghttp.Request) {
	data := models.GetAllInventoryDatabaseList()
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取资产列表成功", result)
}

// 编辑 添加数据库
func (c *InventoryDatabaseController) EditDatabase(r *ghttp.Request) {
	params := models.InventoryDatabase{}
	err := json.Unmarshal(r.GetBody(), &params)
	m := params
	c.CheckError(err, "编辑资产")
	now := utils.GetTimestamp()
	if m.Id == 0 {
		m.AddTime = now
		m.UpdateTime = now
		err = models.Db.Save(&m).Error
		c.CheckError(err, "添加数据库资产失败")
	} else {
		om, err := models.GetInventoryDatabaseOne(m.Id)
		c.CheckError(err, "未找到该数据库资产")
		m.UpdateTime = now
		m.AddTime = om.AddTime
		err = models.Db.Save(&m).Error
		c.CheckError(err, "保存数据库资产失败")
	}
	c.HttpResult(r, enums.CodeSuccess, "保存成功", m.Id)
}

// 删除数据库
func (c *InventoryDatabaseController) DeleteDatabase(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	idList := params
	c.CheckError(err)
	err = models.DeleteInventoryDatabases(idList)
	c.CheckError(err, "删除数据库失败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除数据库", idList)
}
