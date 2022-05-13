package controllers

import (
	"encoding/json"
	"fmt"
	"gfWeb/app/models"
	"gfWeb/app/service/middleware"
	"gfWeb/library/enums"
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Login(r *ghttp.Request) {
	var params struct {
		Account  string
		Password string
	}
	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)
	account := params.Account
	password := params.Password
	c.curUser = models.User{Account: account}
	if len(account) == 0 || len(password) == 0 {
		c.HttpResult(r, enums.CodeFail, "请输入用户名和密码", "")
	}
	password = utils.String2md5(enums.PasswordSalt + password)
	user, err := models.GetUserOneByAccount(account)
	c.CheckError(err)
	c.curUser = *user
	c.userAgent = r.Header.Get("User-Agent")
	//c.CheckError(err)
	now := utils.GetTimestamp()
	if user.CanLoginTime >= now {
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("您已经连续输错%d次密码, 请%d秒后重试", user.ContinueLoginErrorTimes, user.CanLoginTime-now), user.CanLoginTime-now)
	}

	if user.Password != password {
		user.ContinueLoginErrorTimes += 1
		if user.ContinueLoginErrorTimes >= 5 { // map[string]interface{}{"continue_login_error_times": "0", "can_login_time": now + 180}
			err = models.Db.Model(&user).Updates(map[string]interface{}{"can_login_time": now + 180}).Error
			c.CheckError(err)
		} else {
			err = models.Db.Model(&user).Updates(map[string]interface{}{"continue_login_error_times": user.ContinueLoginErrorTimes}).Error
			c.CheckError(err)
		}
		c.HttpResult(r, enums.CodeFail, "密码错误:"+gconv.String(user.ContinueLoginErrorTimes)+"次", password)
		//c.CheckError(err, "用户名或者密码错误")
	}
	if user.IsAccountEnable() == false {
		c.HttpResult(r, enums.CodeFail, "用户被禁用，请联系管理员", "")
	}
	user.LastLoginIp = r.GetClientIp()
	user.LastLoginTime = utils.GetTimestamp()
	user.ContinueLoginErrorTimes = 0
	user.LoginTimes = user.LoginTimes + 1
	err = models.Db.Save(user).Error
	//更新用户登录时间
	//err = models.Db.Model(&user).Updates(&models.User{LastLoginIp: r.GetClientIp(), LastLoginTime: int(gtime.Timestamp()), ContinueLoginErrorTimes: 0, LoginTimes:user.LoginTimes + 1}).Error
	c.CheckError(err)

	//保存用户信息到session
	c.setUser2Session(user.Id)
	g.Log().Infof("登录成功:%v, %v, %v", user.Id, r.Session.Get(models.USER_SESSION_MARK), c.curUser.Id)

	c.HttpResult(r, enums.CodeSuccess, "登录成功",
		struct {
			Token string `json:"token"`
		}{Token: c.Session.Id()})

	//if user != nil && err == nil {
	//	if user.CanLoginTime >= now {
	//		c.HttpResult(r, enums.CodeFail, "您已经连续输错5次密码, 请稍后重试", user.CanLoginTime-now)
	//	}
	//	if user.Status == enums.Disabled {
	//		c.HttpResult(r, enums.CodeFail, "用户被禁用，请联系管理员", "")
	//	}
	//	//保存用户信息到session
	//	c.setUser2Session(user.Id)
	//	g.Log().Info("登录成功:%v, %v, %v", user.Id, c.GetSession("user"), c.curUser.Id)
	//
	//	//更新用户登录时间
	//	err = models.Db.Model(&user).Updates(&models.User{LastLoginIp: c.Ctx.Input.IP(), LastLoginTime: int(time.Now().Unix())}).Error
	//	c.CheckError(err)
	//
	//	c.HttpResult(r, enums.CodeSuccess, "登录成功",
	//		struct {
	//			Token string `json:"token"`
	//		}{Token: c.CruSession.SessionID()})
	//} else {
	//	userOne, err := models.GetUserOneByAccount(account)
	//	g.Log().Debug("%+v%+v", account, userOne)
	//	if err == nil {
	//		userOne.ContinueLoginErrorTimes += 1
	//		if userOne.ContinueLoginErrorTimes > 5 {
	//			err = models.Db.Model(&userOne).Updates(&models.User{ContinueLoginErrorTimes: 0, CanLoginTime: now + 180}).Error
	//			c.CheckError(err)
	//		} else {
	//			err = models.Db.Model(&userOne).Updates(&models.User{ContinueLoginErrorTimes: userOne.ContinueLoginErrorTimes}).Error
	//			c.CheckError(err)
	//		}
	//	}
	//	c.HttpResult(r, enums.CodeFail, "用户名或者密码错误", "")
	//}
}
func (c *LoginController) Logout(r *ghttp.Request) {

	//user := models.User{}
	user := middleware.GetSessionUser(r)
	r.Session.Remove(models.USER_SESSION_MARK)
	utils.DelCache(gconv.String(user.Id))
	c.HttpResult(r, enums.CodeSuccess, "退出登录成功", "")
}
