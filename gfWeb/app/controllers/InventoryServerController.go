package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"

	"github.com/gogf/gf/net/ghttp"
)

type InventoryServerController struct {
	BaseController
}

func (c *InventoryServerController) ServerList(r *ghttp.Request) {
	var params models.InventoryServerParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, total := models.GetInventoryServerList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取资产列表成功", result)
}
func (c *InventoryServerController) AllServerList(r *ghttp.Request) {
	data := models.GetAllServerList()
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取所有资产列表成功", result)
}

// 编辑 添加服务器
func (c *InventoryServerController) EditServer(r *ghttp.Request) {
	params := models.InventoryServer{}
	err := json.Unmarshal(r.GetBody(), &params)
	m := params
	c.CheckError(err, "编辑资产")
	now := utils.GetTimestamp()
	if m.Id == 0 {
		m.AddTime = now
		m.UpdateTime = now
		err = models.Db.Save(&m).Error
		c.CheckError(err, "添加资产失败")
	} else {
		om, err := models.GetInventoryServerOne(m.Id)
		c.CheckError(err, "未找到该资产")
		m.UpdateTime = now
		m.AddTime = om.AddTime
		err = models.Db.Save(&m).Error
		c.CheckError(err, "保存资产失败")
	}
	c.HttpResult(r, enums.CodeSuccess, "保存成功", m.Id)
}

// 删除服务器
func (c *InventoryServerController) DeleteServer(r *ghttp.Request) {
	var idList []int
	err := json.Unmarshal(r.GetBody(), &idList)
	c.CheckError(err)
	err = models.DeleteInventoryServers(idList)
	c.CheckError(err, "删除资产失败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除资产", idList)
}

// 创建ansible
func (c *InventoryServerController) CreateAnsibleInventory(r *ghttp.Request) {
	err := models.CreateAnsibleInventory()
	c.CheckError(err, "生成ansible inventory 失败")
	c.HttpResult(r, enums.CodeSuccess, "生成ansible inventory成功", "")
}
