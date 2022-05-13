package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"gfWeb/memdb"
	"strings"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
)

type PlatformController struct {
	BaseController
}

func (c *PlatformController) List(r *ghttp.Request) {
	var params struct {
		platformIdList []string
	}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("获取平台列表:%+v", params)
	list := models.GetPlatformListByPlatformIdList(params.platformIdList)
	result := make(map[string]interface{})
	result["rows"] = list

	c.HttpResult(r, enums.CodeSuccess, "获取平台列表成功", result)
}

// 编辑 添加平台
func (c *PlatformController) Edit(r *ghttp.Request) {
	params := models.Platform{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err, "编辑平台")
	m := params
	oldPlatform, err := models.GetPlatformOne(params.Id)
	oldPlatform.Id = gstr.Trim(oldPlatform.Id, " ")
	oldPlatform.Name = gstr.Trim(oldPlatform.Name, " ")
	if err == nil {
		m = *oldPlatform
		m.Name = gstr.Trim(params.Name, " ")
		m.Version = params.Version
		m.IsAutoOpenServer = params.IsAutoOpenServer
		m.InventoryDatabaseId = params.InventoryDatabaseId
		m.ZoneInventoryServerId = params.ZoneInventoryServerId
		m.InventorySeverIds = params.InventorySeverIds
		m.TrackerToken = params.TrackerToken
	}
	//更新平台平台服务器关系
	platformInventorySeverRelList, err := models.UpdatePlatformInventorySeverRelByPlatformIdList(m.Id, m.InventorySeverIds)
	c.CheckError(err, "更新的用户角色关系失败")
	m.PlatformInventorySeverRel = platformInventorySeverRelList
	m.Time = gtime.Timestamp()
	err = models.Db.Save(&m).Error
	c.CheckError(err, "编辑平台失败")
	if err != nil {
		pushDb := memdb.OpenDb(memdb.DATABASE_PUSH)
		memdb.ScanDataInTable(pushDb, params.Id)
	}

	//asyncErr := models.AsyncNoticeCenter(m.Id, "", m.TrackerToken)
	//utils.CheckError(asyncErr)

	c.HttpResult(r, enums.CodeSuccess, "保存成功", m.Id)
}

// 删除平台
func (c *PlatformController) Del(r *ghttp.Request) {
	var params []string
	err := json.Unmarshal(r.GetBody(), &params)

	idList := params
	utils.CheckError(err)
	g.Log().Infof("删除平台:%+v", idList)
	err = models.DeletePlatform(idList)
	c.CheckError(err, "删除平台失败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除平台", idList)
}

// 开服管理设置
func (c *PlatformController) EditOpenServerManage(r *ghttp.Request) {
	params := models.Platform{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err, "获取平台数据失败:"+params.Id)
	if len(params.OpenServerTimeScope) > 0 {
		if len(strings.Split(params.ServerAliasStr, "%d")) != 2 {
			c.CheckError(err, "区服别名内容错误:"+params.ServerAliasStr)
		}
	}
	oldPlatform, err := models.GetPlatformOne(params.Id)
	m := models.Platform{}
	m = *oldPlatform
	m.CreateRoleLimit = params.CreateRoleLimit
	m.OpenServerTakeTime = params.OpenServerTakeTime
	m.IntervalInitTime = params.IntervalInitTime
	m.IntervalDay = params.IntervalDay
	m.OpenServerTimeScope = params.OpenServerTimeScope
	m.ServerAliasStr = params.ServerAliasStr
	m.OpenServerChangeTime = gtime.Timestamp()
	err = models.SaveOpenServerManage(c.curUser.Id, oldPlatform, &m)
	c.CheckError(err)
	//err = models.SaveOpenServerManageLog(c.curUser.Id, oldPlatform, &m)
	//c.CheckError(err, "开服管理日志保存失败")
	//m.OpenServerChangeTime = gtime.Timestamp()
	//err = models.Db.Save(&m).Error
	//c.CheckError(err, "开服管理设置失败")
	//models.InitCronPlatformOpenServerTime(&m)
	c.HttpResult(r, enums.CodeSuccess, "成功开服管理设置", "")
}
