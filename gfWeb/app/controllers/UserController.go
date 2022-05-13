package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type UserController struct {
	BaseController
}

func (c *UserController) Info(r *ghttp.Request) {
	//c.ControllerInit(r)
	m := c.curUser
	platformList := models.GetPlatformListByUserId(m.Id)
	platformIdList := make([]string, 0)
	for _, e := range platformList {
		platformIdList = append(platformIdList, e.Id)
	}
	//utils.CheckError(err)
	gameServerList := models.GetServerList(platformIdList)
	channelList := models.GetChannelListByPlatformIdList(platformIdList)
	c.HttpResult(r, enums.CodeSuccess, "获取用户信息成功",
		struct {
			Name           string         `json:"name"`
			ResourceTree   []*models.Menu `json:"menuTree"`
			PlatformList   []*models.Platform
			ChannelList    []*models.Channel
			GameServerList []*models.Server
			IsSuper        int
		}{
			Name:           m.Name,
			ResourceTree:   models.TranMenuList2MenuTree(models.GetMenuListByUserId(m.Id), true),
			PlatformList:   platformList,
			GameServerList: gameServerList,
			ChannelList:    channelList,
			IsSuper:        m.IsSuper,
		})
}

// 获取用户列表
func (c *UserController) List(r *ghttp.Request) {
	var params models.UserQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	//c.ControllerInit(r)
	utils.CheckError(err)
	g.Log().Infof("查询用户列表:%+v", params)
	data, total := models.GetUserList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取用户列表成功", result)
}

// 获取精简的用户列表
func (c *UserController) SimpleList(r *ghttp.Request) {
	data := models.GetUserSimpleList()
	c.HttpResult(r, enums.CodeSuccess, "获取精简的用户列表", data)
}

// 编辑 添加用户
func (c *UserController) Edit(r *ghttp.Request) {
	params := models.User{}
	err := json.Unmarshal(r.GetBody(), &params)
	//c.ControllerInit(r)
	g.Log().Infof("编辑用户:%+v", params)
	utils.CheckError(err, "编辑用户")
	changeUserId := c.curUser.Id

	//删除旧的用户角色关系
	_, err = models.DeleteRoleUserRelByUserIdList([]int{params.Id})
	c.CheckError(err, "删除旧的用户角色关系失败")
	for _, roleId := range params.RoleIds {

		if params.IsSuper == 1 {
			// 如果不是超级管理员，报错
			if ok, err := models.RoleIsSuper(roleId); !ok {
				c.CheckError(err, "非超级管理员")
			}
		}
		theTime := time.Now()
		relation := models.RoleUserRel{UserId: params.Id, RoleId: roleId, Created: theTime}
		params.RoleUserRel = append(params.RoleUserRel, &relation)
	}
	if params.Id == 0 {
		params.Password = utils.String2md5(enums.PasswordSalt + params.ModifyPassword)
		err = models.Db.Save(&params).Error
		c.CheckError(err, "添加用户失败")
	} else {
		oM, err := models.GetUserOne(params.Id)
		c.CheckError(err, "未找到该用户")
		if changeUserId == params.Id && params.Status == 0 {
			c.HttpResult(r, enums.CodeFail2, "不能禁用自己", "")
		}
		params.Password = strings.TrimSpace(params.ModifyPassword)
		if len(params.Password) == 0 {
			//密码为空则不修改
			//params.Password = oM.Password
		} else {
			oM.Password = utils.String2md5(enums.PasswordSalt + params.Password)
		}
		oM.Name = params.Name
		oM.IsSuper = params.IsSuper
		oM.Status = params.Status
		oM.Mobile = params.Mobile
		oM.MailStr = params.MailStr
		oM.RoleUserRel = params.RoleUserRel
		err = models.Db.Save(&oM).Error
		c.CheckError(err, "保存用户失败")
	}
	c.HttpResult(r, enums.CodeSuccess, "保存成功", params.Id)
}

// 删除用户
func (c *UserController) Delete(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	//c.ControllerInit(r)
	utils.CheckError(err)
	userIdList := params
	changeUserId := c.curUser.Id
	for _, checkUserId := range userIdList {
		if changeUserId == checkUserId {
			c.CheckError(err, "删除用户中不能包含自己")
		}
	}
	g.Log().Infof("删除用户:%+v", userIdList)
	err = models.DeleteUsers(userIdList)
	c.CheckError(err, "删除用户失败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除用户", userIdList)
}

// 解除用户状态
func (c *UserController) RemoveState(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	//c.ControllerInit(r)
	utils.CheckError(err)
	userIdList := params
	flag := false
	for _, userId := range userIdList {
		oM, err := models.GetUserOne(userId)
		if err != nil {
			g.Log().Warningf("解除用户状态失败:%d err:%+v", userId, err)
			continue
		}
		oM.CanLoginTime = 0
		oM.ContinueLoginErrorTimes = 0
		err = models.Db.Save(&oM).Error
		if err != nil {
			g.Log().Warningf("解除用户状态保存失败:%d err:%+v", userId, err)
			continue
		}
		flag = true
	}
	if flag == false {
		c.CheckError(err, "解除用户状态失败")
	}
	c.HttpResult(r, enums.CodeSuccess, "成功解除用户状态", "")
}

//修改密码
func (c *UserController) ChangePassword(r *ghttp.Request) {
	var params struct {
		OldPwd string `json:"oldPwd"`
		NewPwd string `json:"newPwd"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	//c.ControllerInit(r)
	utils.CheckError(err)
	g.Log().Infof("修改密码:%+v", params)
	Id := c.curUser.Id
	user, err := models.GetUserOne(Id)
	c.CheckError(err, "未找到该用户")
	md5str := utils.String2md5(enums.PasswordSalt + params.OldPwd)
	if user.Password != md5str {
		c.HttpResult(r, enums.CodeFail, "原密码错误", "")
	}
	if len(params.NewPwd) == 0 {
		c.HttpResult(r, enums.CodeFail, "请输入新密码", "")
	}
	user.Password = utils.String2md5(enums.PasswordSalt + params.NewPwd)
	err = models.Db.Model(&user).Updates(models.User{Password: user.Password}).Error
	c.CheckError(err, "保存失败")
	c.setUser2Session(Id)
	c.HttpResult(r, enums.CodeSuccess, "保存成功", user.Id)
}
