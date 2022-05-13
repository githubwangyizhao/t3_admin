package controllers

import (
	"fmt"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"strconv"
	"strings"

	"github.com/gogf/gf/container/gqueue"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
)

type BaseController struct {
	gmvc.Controller
	curUser   models.User //当前用户信息
	userAgent string      // 客户端信息
}

var PUSHQUEUE *gqueue.Queue

//func (c *BaseController) AllowCross() {
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:9528")       //允许访问源
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "*")    //允许post访问
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Token") //header的类型
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
//}

func (c *BaseController) ControllerInit(r *ghttp.Request) {
	c.Init(r)
	user := c.Session.Get(models.USER_SESSION_MARK, nil)
	if user != nil {
		gconv.Struct(user, &c.curUser)
	}
	////判断是否登录
	//c.checkLogin()
	////判断是否有权限
	//c.checkAuthor()
}

//检查错误, 失败直接终止当前请求
func (c *BaseController) CheckError(err error, msg ...string) {
	if err != nil {
		errMsg := ""
		if len(msg) == 0 {
			errMsg = fmt.Sprintf("%v", err)
		} else {
			errMsg = fmt.Sprintf("%s %v", msg, err)
		}
		g.Log().Error(errMsg)
		c.HttpResult(c.Request, enums.CodeFail, errMsg, nil)
	}
}

//是否登录
func (c *BaseController) IsLogin() bool {
	return c.curUser.Id > 0
}

//是否帐号有效
func (c *BaseController) IsAccountEnable() bool {
	return c.curUser.Status == enums.Enabled
}

// 判断某 Controller.Action 当前用户是否有权访问
func (c *BaseController) checkActionAuthor(ctrlName, ActName string) bool {
	if c.IsLogin() == false || c.IsAccountEnable() == false {
		return false
	}
	user := c.Session.Get(models.USER_SESSION_MARK)
	v, ok := user.(models.User)
	if ok {
		//如果是超级管理员，则直接通过
		if v.IsSuperUser() {
			return true
		}
		//遍历用户有权限的资源列表
		for _, v := range v.ResourceUrlForList {
			if v == ctrlName+"."+ActName || v == ctrlName+".*" {
				return true
			}
		}
	}
	return false
}

func (c *BaseController) setUser2Session(userId int) error {
	m, err := models.GetUserOne(userId)
	if err != nil {
		return err
	}
	resourceList := models.GetResourceListByUserId(userId)
	for _, item := range resourceList {
		m.ResourceUrlForList = append(m.ResourceUrlForList, strings.TrimSpace(item.UrlFor))
	}
	c.Session.Set(models.USER_SESSION_MARK, *m)
	utils.SetCacheCode(strconv.Itoa(m.Id), c.Session.Id(), 0)
	return nil
}

////请求返回json
//func (c *BaseController) Result(code enums.ResultCode, msg string, data interface{}) {
//	r := &models.Result{Code: code, Data: data, Msg: msg}
//	c.Data["json"] = r
//	c.ServeJSON()
//	c.StopRun()
//}
//请求返回json
func (c *BaseController) HttpResult(request *ghttp.Request, code enums.ResultCode, msg string, data interface{}) {
	//g.Log().Debugf("请求返回json :%v  %v", c.Request.URL.Path, msg)
	models.SaveLoginAdminLog(c.curUser, c.userAgent, request, code, msg, data)
	models.SaveRequestAdminLog(c.curUser.Id, enums.LoginTypeSuccess, code, msg, request)
	result := &models.Result{Code: code, Data: data, Msg: msg}
	//c.Data["json"] = r
	//c.Response.Write(result)
	err := request.Response.WriteJsonExit(result)
	c.CheckError(err, "请求返回json失败")
	//c.ServeJSON()
	//c.StopRun()
}
