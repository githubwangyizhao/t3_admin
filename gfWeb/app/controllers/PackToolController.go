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
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

type PackToolController struct {
	BaseController
}

// 打包工具
func (c *PackToolController) PackTool(r *ghttp.Request) {
	var params struct {
		PackServer string `json:"packServer"`
		PackType   string `json:"packType"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	g.Log().Info("打包工具:%+v", params)
	ResultStr, err := models.PackTool(params.PackServer, params.PackType)
	if err != nil {
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("%v", err), "")
		return
	}
	//ResultStr := strings.Replace(out, "正在", "\n", -1)
	c.HttpResult(r, enums.CodeSuccess, "打包成功", ResultStr)
}

// 同步工具
func (c *PackToolController) SyncTool(r *ghttp.Request) {
	var params struct {
		SyncDir      string `json:"syncDir"`
		SyncPlatform string `json:"syncPlatform"`
		ShFile       string `json:"shFile"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	g.Log().Info("同步工具:%+v", params)
	out, err := models.SyncTool(params.SyncDir, params.SyncPlatform, params.ShFile)
	if err != nil {
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("%v", err), "")
		return
	}
	c.HttpResult(r, enums.CodeSuccess, "同步成功", out)
}

// 更新平台工具
func (c *PackToolController) UpdatePlatformTool(r *ghttp.Request) {
	//var params struct {
	//	UpdateType           int `json:"updateType"` // 更新方式（0:热更;1:冷更）
	//	PlatformIdList       []string
	//	UpdateEnterStateType int `json:"updateEnterStateType"`
	//}
	//err := json.Unmarshal(r.GetBody(),&params)
	//
	//c.CheckError(err)
	////g.Log().Infof("更新平台工具:%+v", params)
	//out, err := models.UpdatePlatformVersion(c.curUser.Id, params.UpdateType, params.UpdateEnterStateType, params.PlatformIdList)
	//if err != nil {
	//	c.HttpResult(r, enums.CodeFail2, fmt.Sprintf("%v", err), "")
	//	return
	//}
	c.HttpResult(r, enums.CodeSuccess, "", "")
}

// 获取平台精简数据
func (c *PackToolController) GetPlatformSimpleList(r *ghttp.Request) {
	var platformIdList []string

	err := json.Unmarshal(r.GetBody(), &platformIdList)
	c.CheckError(err, "获取平台精简数据:")
	data := models.GetPlatformSimpleListByPlatformIdList(platformIdList)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "成功!", result)
}

// 订时更新平台版本
func (c *PackToolController) UpdatePlatformVersionCron(r *ghttp.Request) {
	var params struct {
		PlatformIdList       []string
		UpdateType           int `json:"updateType"`           // 更新方式（0:热更;1:冷更,10:开启，11:关闭）
		UpdateEnterStateType int `json:"updateEnterStateType"` // 入口方式(0:不变更，1：关入口，2：关闭更新后开启)
		CronUpdateTime       int `json:"cronUpdateTime"`       // 订时更新时间(小于当前时间，立即执行)
		RemoveUpdate         int `json:"removeUpdate"`         // 移除更新(1为移除)
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err, "订时更新平台版本解析失败:")
	err = models.UpdatePlatformVersionCron(c.curUser.Id, params.UpdateType, params.UpdateEnterStateType, params.CronUpdateTime, params.RemoveUpdate, params.PlatformIdList)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "更新平台版本操作成功", "")
}

// 操作平台工具
func (c *PackToolController) ChangePlatformTool(r *ghttp.Request) {
	var params struct {
		ChangeType     string `json:"changeType"` // 操作方式(start, stop ,restart)
		PlatformIdList []string
	}
	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	g.Log().Info("更新平台工具:%+v", params)
	if params.ChangeType != "start" && params.ChangeType != "stop" && params.ChangeType != "restart" {
		c.HttpResult(r, enums.CodeFail, "操作方式错误："+params.ChangeType, "")
	}
	out, err := models.ChangePlatformTool(params.ChangeType, params.PlatformIdList)
	if err != nil {
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("%v", err), "")
		return
	}
	c.HttpResult(r, enums.CodeSuccess, out, "")
}

// 获得版本路径详细
func (c *PackToolController) GetPlatformVersionInfo(r *ghttp.Request) {
	var params struct {
		ChangeType int `json:"changeType"` // 操作类型0:打包 1:同步 2:更新
		RobotType  int `json:"robotType"`  // 机器类型:0=客户端 1=服务器端
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data := models.GetPlatformVersionByChangeTypeOrVersionType(params.ChangeType, params.RobotType)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "成功!", result)
}

// 获得版本路径列表
func (c *PackToolController) GetBranchPath(r *ghttp.Request) {
	var params = &models.BaseQueryParam{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, count := models.GetBranchPathList(params)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "成功!", result)
}

// 获得平台版本路径列表
func (c *PackToolController) GetPlatformVersionPath(r *ghttp.Request) {
	var params = &models.BaseQueryParam{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, count := models.GetPlatformVersionPathList(params)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "成功!", result)
}

// 获得版本操作列表
func (c *PackToolController) GetVersionToolChange(r *ghttp.Request) {
	var params = &models.BaseQueryParam{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, count := models.GetVersionToolChangeList(params)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "成功!", result)
}

// 获得版本操作详情列表
func (c *PackToolController) GetVersionToolChangeInfo(r *ghttp.Request) {
	var params struct {
		ChangeType int `json:"changeType"` // 操作类型0:打包 1:同步 2:更新
		RobotType  int `json:"robotType"`  // 机器类型:0=客户端 1=服务器端
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data := models.GetVersionToolChangeInfo(params.ChangeType, params.RobotType)
	//result := make(map[string]interface{})
	//result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "成功!", data)
}

// 更新版本路径数据
func (c *PackToolController) BranchPathEdit(r *ghttp.Request) {
	var params = &models.BranchPath{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	OldData, err := models.GetBranchPathOne(params.Id)
	if params.IsAdd == 1 && err == nil {
		g.Log().Warningf("版本路径数据已经存在:%+v", OldData)
		c.HttpResult(r, enums.CodeFail, "版本路径数据已经存在", "")
	}
	if params.IsAdd == 0 && err != nil {
		c.HttpResult(r, enums.CodeFail, "版本路径数据不存在", "")
	}
	params.UserId = c.curUser.Id
	err = models.EditBranchPath(params)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

// 更新平台版本路径数据
func (c *PackToolController) PlatformVersionPathEdit(r *ghttp.Request) {
	var params = &models.PlatformVersionPath{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	OldData, err := models.GetPlatformVersionPathOne(params.Id)
	if params.IsAdd == 1 && err == nil {
		g.Log().Warningf("平台版本路径数据已经存在:%+v", OldData)
		c.HttpResult(r, enums.CodeFail, "平台版本路径数据已经存在", "")
	}
	if params.IsAdd == 0 && err != nil {
		c.HttpResult(r, enums.CodeFail, "平台版本路径数据不存在", "")
	}
	params.UserId = c.curUser.Id
	err = models.EditPlatformVersionPath(params)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

// 更新版本操作数据
func (c *PackToolController) VersionToolChangeEdit(r *ghttp.Request) {
	var params = &models.VersionToolChange{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	OldData, err := models.GetVersionToolChangeOne(params.Id)
	if params.IsAdd == 1 && err == nil {
		g.Log().Warningf("版本操作数据已经存在:%+v", OldData)
		c.HttpResult(r, enums.CodeFail, "版本操作数据已经存在", "")
	}
	if params.IsAdd == 0 && err != nil {
		c.HttpResult(r, enums.CodeFail, "版本操作数据不存在", "")
	}
	params.UserId = c.curUser.Id
	err = models.EditVersionToolChangePath(params)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

// 删除版本路径数据
func (c *PackToolController) DelBranchPath(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	idList := params
	c.CheckError(err)
	err = models.DeleteBranchPath(idList)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

// 删除平台版本路径数据
func (c *PackToolController) DelPlatformVersionPath(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	idList := params
	c.CheckError(err)
	err = models.DeletePlatformVersionPath(idList)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

// 删除版本工具数据
func (c *PackToolController) DeleteVersionToolChange(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	idList := params
	c.CheckError(err)
	err = models.DeleteVersionToolChange(idList)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

// 发送平台版本操作处理
func (c *PackToolController) SendVersionToolChange(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	err = models.SendVersionToolChange(c.curUser.Id, params)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

// 获得订时平台版本操作
func (c *PackToolController) GetVersionToolChangeCron(r *ghttp.Request) {
	var params = &models.ParamVersionToolChangeCron{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, count := models.GetVersionToolChangeCronList(params)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "成功!", result)
}

// 订时平台版本操作编辑 新增
func (c *PackToolController) VersionToolChangeCronEdit(r *ghttp.Request) {
	params := &models.VersionToolChangeCron{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	_, err = models.GetVersionToolChangeCron(params.Id)
	if params.IsAdd == 1 && err == nil {
		c.HttpResult(r, enums.CodeFail, "订时平台版本操作已经存在", params.Id)
	}
	if params.IsAdd == 0 && err != nil {
		g.Log().Errorf("订时平台版本操作不存在:%+v  err:%+v", params.Id, err)
		c.HttpResult(r, enums.CodeFail, "订时平台版本操作不存在", params.Id)
	}
	g.Log().Debugf("params.ChangeIdList:%+v", params.ChangeIdList)
	params.ChangeIdStr = strings.Join(gconv.Strings(params.ChangeIdList), ",")
	params.CronTimeStr = utils.TimestampToCronStr(params.CronTime)
	g.Log().Debugf("params.ChangeIdStr:%T", params.ChangeType)
	params.UserId = c.curUser.Id
	params.ChangeTime = gtime.Timestamp()
	err = models.Db.Save(params).Error
	c.CheckError(err)
	models.InitVersionToolChangeCron(params)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

// 删除订时平台版本操作
func (c *PackToolController) DelVersionToolChangeCron(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	idList := params
	c.CheckError(err)
	err = models.DeleteVersionToolChangeCron(idList)
	c.CheckError(err, "删除订时平台版本操作失败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除订时平台版本操作", idList)
}

// 获得发送平台版本操作日志
func (c *PackToolController) GetVersionToolChangeLog(r *ghttp.Request) {
	var params = &models.ParamsVersionTool{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, count := models.GetVersionToolLogList(params)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "成功!", result)
}
