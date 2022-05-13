package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"

	"github.com/gogf/gf/net/ghttp"
)

type ChannelController struct {
	BaseController
}

// 获取渠道列表
func (c *ChannelController) List(r *ghttp.Request) {

	list := models.GetChannelList()
	result := make(map[string]interface{})
	result["rows"] = list
	c.HttpResult(r, enums.CodeSuccess, "获取渠道列表成功", result)
}

// 编辑 添加渠道
func (c *ChannelController) Edit(r *ghttp.Request) {
	params := models.Channel{}
	err := json.Unmarshal(r.GetBody(), &params)
	m := params
	//err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	c.CheckError(err, "编辑渠道")
	err = models.Db.Save(&m).Error
	c.CheckError(err, "编辑渠道失败")

	//asyncErr := models.AsyncNoticeCenter(m.PlatformId, m.Channel, m.TrackerToken)
	asyncErr := utils.AsyncNoticeCenterRegion(m.PlatformId, m.Channel, m.TrackerToken, m.Region, m.AreaCode, m.Currency)
	c.CheckError(asyncErr)

	c.HttpResult(r, enums.CodeSuccess, "保存成功", m.Id)
}

// 删除渠道
func (c *ChannelController) Del(r *ghttp.Request) {
	var idList []int
	err := json.Unmarshal(r.GetBody(), &idList)
	c.CheckError(err)
	err = models.DeleteChannel(idList)
	c.CheckError(err, "删除渠道失败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除渠道", idList)
}
