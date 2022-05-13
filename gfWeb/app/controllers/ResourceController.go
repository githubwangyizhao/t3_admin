package controllers

import (
	"encoding/json"
	"fmt"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"strings"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type ResourceController struct {
	BaseController
}

//获取资源列表
func (c *ResourceController) List(r *ghttp.Request) {

	//获取数据列表和总数
	data := models.TranResourceList2ResourceTree(models.GetResourceList())
	result := make(map[string]interface{})
	c.UrlFor2Link(data)
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取资源列表成功", result)
}

func (c *ResourceController) ResourceTree(r *ghttp.Request) {

	c.HttpResult(r, enums.CodeSuccess, "获取资源树成功", models.TranResourceList2ResourceTree(models.GetResourceList()))
}

//获取可以成为某节点的父节点列表
func (c *ResourceController) GetParentResourceList(r *ghttp.Request) {
	var params struct {
		Id int `json:"id"`
	}

	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	tree := models.ResourceTreeGrid4Parent(params.Id)
	c.HttpResult(r, enums.CodeSuccess, "", tree)
}

// 将资源表里的UrlFor值转成LinkUrl
func (c *ResourceController) UrlFor2LinkOne(urlfor string) string {
	//if len(urlfor) == 0 {
	//	return ""
	//}
	//strs := strings.Split(urlfor, ",")
	//if len(strs) == 1 {
	//	return c.Request.URL(strs[0])
	//} else if len(strs) > 1 {
	//	var values []interface{}
	//	for _, val := range strs[1:] {
	//		values = append(values, val)
	//	}
	//	return c.URLFor(strs[0], values...)
	//}
	return ""
}

//UrlFor2Link 使用URLFor方法，批量将资源表里的UrlFor值转成LinkUrl
func (c *ResourceController) UrlFor2Link(src []*models.Resource) {
	for _, item := range src {
		item.Url = c.UrlFor2LinkOne(item.UrlFor)
		c.UrlFor2Link(item.Children)
	}
}

//编辑添加资源
func (c *ResourceController) Edit(r *ghttp.Request) {
	params := models.Resource{}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err, "编辑资源")
	g.Log().Info("编辑资源:%+v", params)
	parent := &models.Resource{}
	parentId := params.ParentId
	//获取父节点
	if parentId > 0 {
		parent, err = models.GetResourceOne(parentId)
		c.CheckError(err, "父节点无效")
		params.Parent = parent
	}
	if c.UrlFor2LinkOne(params.UrlFor) != "" || strings.Contains(params.UrlFor, ".*") {
	} else {
		c.HttpResult(r, enums.CodeFail, "控制器解析失败: "+params.UrlFor, "")
	}
	if params.Id == 0 {
		err = models.Db.Save(&params).Error
		c.CheckError(err, "添加资源失败")
		c.HttpResult(r, enums.CodeSuccess, "添加资源成功", params.Id)

	} else {
		if parentId > 0 {
			if models.CanParent(params.Id, parentId) == false {
				c.HttpResult(r, enums.CodeFail, "请重新选择父节点", "")
			}
		}
		err = models.Db.Save(&params).Error
		c.CheckError(err, "编辑资源失败")
		c.HttpResult(r, enums.CodeSuccess, "编辑资源成功", params.Id)
	}
}

//删除资源
func (c *ResourceController) Delete(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Info("删除资源:%+v", params)
	_, err = models.DeleteResources(params)
	c.CheckError(err, "删除资源失败")
	c.HttpResult(r, enums.CodeSuccess, fmt.Sprintf("删除资源成功"), 0)
}

//CheckUrlFor 填写UrlFor时进行验证
func (c *ResourceController) CheckUrlFor(r *ghttp.Request) {

	urlfor := r.GetQueryString("urlfor")
	link := c.UrlFor2LinkOne(urlfor)
	if len(link) > 0 {
		c.HttpResult(r, enums.CodeSuccess, "解析成功", link)
	} else {
		c.HttpResult(r, enums.CodeFail, "解析失败", link)
	}

	g.DB().Begin()
}
