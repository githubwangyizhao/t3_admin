package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type MenuController struct {
	BaseController
}

//获取菜单列表
func (c *MenuController) List(r *ghttp.Request) {

	data := models.TranMenuList2MenuTree(models.GetMenuList(), false)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取菜单列表成功", result)
}

func (c *MenuController) MenuTree(r *ghttp.Request) {

	c.HttpResult(r, enums.CodeSuccess, "获取菜单树成功", models.TranMenuList2MenuTree(models.GetMenuList(), false))
}

//获取可以成为某节点的父节点列表
func (c *MenuController) GetParentMenuList(r *ghttp.Request) {
	var params struct {
		Id int `json:"id"`
	}

	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Debug("获取可以成为某节点的父节点列表:%+v", params)
	tree := models.MenuTreeGrid4Parent(params.Id)
	c.HttpResult(r, enums.CodeSuccess, "", tree)
}

//编辑添加菜单
func (c *MenuController) Edit(r *ghttp.Request) {
	params := models.Menu{}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err, "编辑菜单")
	g.Log().Info("编辑菜单:%+v", params)
	parentId := params.ParentId
	//获取父节点
	if parentId > 0 {
		parent, err := models.GetMenuOne(parentId)
		c.CheckError(err, "父节点无效")
		params.Parent = parent
	}
	if params.Id == 0 {
		err = models.Db.Save(&params).Error
		c.CheckError(err, "添加菜单失败")
		c.HttpResult(r, enums.CodeSuccess, "添加菜单成功", params.Id)
	} else {
		if parentId > 0 {
			if models.CanParentMenu(params.Id, parentId) == false {
				c.HttpResult(r, enums.CodeFail, "请重新选择父节点", "")
			}
		}
		err = models.Db.Save(&params).Error
		c.CheckError(err, "编辑菜单失败")
		c.HttpResult(r, enums.CodeSuccess, "编辑菜单成功", params.Id)
	}
}

// 删除菜单
func (c *MenuController) Delete(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Info("删除菜单:%+v", params)
	_, err = models.DeleteMenus(params)
	c.CheckError(err, "删除菜单失败")
	c.HttpResult(r, enums.CodeSuccess, "删除菜单成功", "")
}
