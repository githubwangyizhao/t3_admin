// 登录公告管理
package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type LoginNoticeController struct {
	BaseController
}

// 设置登录公告
func (c *LoginNoticeController) SetNotice(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:"platformId"`
		ChannelId  string `json:"channelId"`
		Notice     string `json:"notice"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	g.Log().Infof("设置登录公告:%+v", params)
	//c.CheckError(err)
	//
	//data, err := json.Marshal(params)
	//utils.CheckError(err)
	//
	//url := utils.GetCenterURL() + "/set_login_notice"
	//err = utils.HttpRequest(url, string(data))
	////out, err := utils.NodeTool(
	////	"mod_login_notice",
	////	"update_login_notice",
	////	params.PlatformId,
	////	params.Notice,
	////)
	//c.CheckError(err, "设置中心服登录公告:" + url)
	//noticeLog := &models.LoginNotice{
	//	PlatformId:  params.PlatformId,
	//	Notice:     params.Notice,
	//	Time:        time.Now().Unix(),
	//	UserId:      c.curUser.Id,
	//}
	//err = models.Db.Save(&noticeLog).Error
	err = models.UpdateAndPushLoginNotice(c.curUser.Id, params.PlatformId, params.ChannelId, params.Notice)
	c.CheckError(err, "写登录公告日志失败")
	c.HttpResult(r, enums.CodeSuccess, "设置登录公告成功", 0)
}

// 批量设置登录公告
func (c *LoginNoticeController) BatchSetNotice(r *ghttp.Request) {
	var params struct {
		Id     []string
		Notice string
	}
	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	g.Log().Infof("批量设置登录公告:%+v", params)
	err = models.BatchSetNotice(c.curUser.Id, params.Id, params.Notice)
	c.CheckError(err, "写登录公告日志失败")
	//for _, PlatformId := range params.PlatformIdList {
	//	//var request struct {
	//	//	PlatformId string `json:"platformId"`
	//	//	Notice    string	`json:"notice"`
	//	//}
	//	//request.PlatformId = PlatformId
	//	//request.Notice = params.Notice
	//	//data, err := json.Marshal(request)
	//	//utils.CheckError(err)
	//	//url := utils.GetCenterURL() + "/set_login_notice"
	//	//err = utils.HttpRequest(url, string(data))
	//	////out, err := utils.NodeTool(
	//	////	"mod_login_notice",
	//	////	"update_login_notice",
	//	////	PlatformId,
	//	////	params.Notice,
	//	////)
	//	//c.CheckError(err, "批量中心服登录公告")
	//	//noticeLog := &models.LoginNotice{
	//	//	PlatformId:  PlatformId,
	//	//	Notice:     params.Notice,
	//	//	Time:        time.Now().Unix(),
	//	//	UserId:      c.curUser.Id,
	//	//}
	//	//err = models.Db.Save(&noticeLog).Error
	//
	//	err = models.UpdateAndPushLoginNotice(c.curUser.Id, PlatformId, params.Notice)
	//	c.CheckError(err, "写登录公告日志失败")
	//}
	c.HttpResult(r, enums.CodeSuccess, "批量设置登录公告成功", 0)
}

// 获取登录公告列表
func (c *LoginNoticeController) LoginNoticeList(r *ghttp.Request) {
	//g.Log().Infof("查询公告日志:%+v", params)
	var params struct {
		PlatformIdList []string
	}
	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	g.Log().Infof("获取登录公告列表:%+v", params)
	data := models.GetLoginNoticeListByPlatformIdList(params.PlatformIdList)
	result := make(map[string]interface{})
	//result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取登录公告列表", result)
}

//删除登录公告
func (c *LoginNoticeController) DelLoginNotice(r *ghttp.Request) {
	var idList []string
	err := json.Unmarshal(r.GetBody(), &idList)
	c.CheckError(err)
	g.Log().Infof("删除登录公告列表:%+v", idList)
	err = models.DeleteLoginNotice(c.curUser.Id, idList)
	c.CheckError(err, "删除登录公告失败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除登录公告", idList)
}
