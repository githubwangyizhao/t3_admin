package controllers

import (
	"encoding/json"
	"fmt"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type GiftController struct {
	BaseController
}

// List 礼包列表页
func (c *GiftController) List(r *ghttp.Request) {
	params := &models.GiftRequest{}
	g.Log().Infof("获取指定礼包码的数据:%+v", params)
	err := json.Unmarshal(r.GetBody(), &params)
	if err != nil {
		g.Log().Infof("礼包码获取参数错误: %s", fmt.Sprintf("%v", err))
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("%v", err), nil)
	}

	data, count, err := models.GetGiftList(params)
	if err != nil {
		c.HttpResult(r, enums.CodeFail2, "礼包码获取失败", nil)
	}
	g.Log().Infof("data: %+v", data)

	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data

	g.Log().Infof("count: %d %d", count, len(data))
	c.HttpResult(r, enums.CodeSuccess, "礼包码创建成功", result)
}

func (c *GiftController) Delete(r *ghttp.Request) {
	params := &models.Gift{}
	g.Log().Infof("删除礼包码:%+v", params)
	err := json.Unmarshal(r.GetBody(), &params)
	if err != nil {
		g.Log().Infof("删除礼包码参数错误: %s", fmt.Sprintf("%v", err))
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("%v", err), nil)
	}

	data, _ := json.Marshal(params)
	url := utils.GetCenterURL() + "/delete_gift_code"

	_, err = utils.HttpRequest(url, string(data))
	if err != nil {
		g.Log().Infof("删除礼包码接口调用失败: %s %s %s", url, data, fmt.Sprintf("%v", err))
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("%v", err), nil)
	}

	c.HttpResult(r, enums.CodeSuccess, "礼包码删除成功", nil)
}

func (c *GiftController) Create(r *ghttp.Request) {
	params := &models.Gift{}
	g.Log().Infof("创建礼包码:%+v", params)
	err := json.Unmarshal(r.GetBody(), &params)
	if err != nil {
		g.Log().Infof("创建礼包码参数错误: %s", fmt.Sprintf("%v", err))
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("%v", err), nil)
	}

	params.UserId = c.curUser.Id
	data, _ := json.Marshal(params)
	url := utils.GetCenterURL() + "/add_gift_code"
	_, err = utils.HttpRequest(url, string(data))
	if err != nil {
		g.Log().Infof("创建礼包码接口调用失败: %s %s %s", url, data, fmt.Sprintf("%v", err))
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("%v", err), nil)
	}

	c.HttpResult(r, enums.CodeSuccess, "礼包码创建成功", nil)
}

func (c *GiftController) Update(r *ghttp.Request) {
	params := &models.Gift{}
	g.Log().Infof("追加礼包码数量:%+v", params)
	err := json.Unmarshal(r.GetBody(), &params)
	if err != nil {
		g.Log().Infof("追加礼包码数量参数错误: %s", fmt.Sprintf("%v", err))
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("%v", err), nil)
	}

	data, _ := json.Marshal(params)
	url := utils.GetCenterURL() + "/append_gift_code"
	_, err = utils.HttpRequest(url, string(data))
	if err != nil {
		g.Log().Infof("追加礼包码接口调用失败: %s %s %s", url, data, fmt.Sprintf("%v", err))
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("%v", err), nil)
	}

	c.HttpResult(r, enums.CodeSuccess, "礼包码追加成功", nil)

}

func (c *GiftController) Download(r *ghttp.Request) {
	params := &models.GiftRequest{}
	g.Log().Infof("下载礼包码excel:%+v", params)
	err := json.Unmarshal(r.GetBody(), &params)
	if err != nil {
		g.Log().Infof("礼包码获取参数错误: %s", fmt.Sprintf("%v", err))
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("%v", err), nil)
	}

	data, err := models.DownloadGiftCode(params)
	if err != nil {
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("%v", err), nil)
	}

	c.HttpResult(r, enums.CodeSuccess, "礼包码excel数据导出成功", data)
}
