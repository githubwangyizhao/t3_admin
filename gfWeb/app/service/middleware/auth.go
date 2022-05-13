package middleware

import (
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"strings"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
)

//框架页面只判断是否登陆不做权限判断
//var FramePages = []string{"/", "/index", "/system/main", "/system/download", "/system/switchSkin", "/login", "/logout"}
//var FramePages = []string{"/ws", "/wss", "/system/tool/admin", "/system/tool/admin/restart", "/login"}
var FramePages = []string{"/ws", "/wss", "/system/tool/admin", "/login", "/tool/get_version_list", "/tool/get_platform_client_info_all_list", "/api/get_player_info", "/api/get_platform_info_list", "/api/get_platform_info", "/api/get_adjust_list", "/api/get_client_verify", "/api/get_app_info_list", "/api/get_all_tracker_info", "/api/get_area_code_list"}

// 白名单, 添加时防止重复,一定要用三段字符串格式
var WhiteActionList = []string{"test_account_list", "push_list"}

// 鉴权中间件，只有登录成功之后才能通过
func Auth(r *ghttp.Request) {
	path := r.URL.Path
	user := GetSessionUser(r)

	//开启过滤白名单
	for _, v := range WhiteActionList {
		if strings.Contains(path, v) {
			return
		}
	}

	if user.Id == 0 {
		g.Log().Infof("[%v]未登录接口请求method:%v\t uri:%v\t param:%+v", r.GetClientIp(), r.Method, r.GetUrl(), r.GetBodyString())
	} else {
		g.Log().Infof("[%v:%v][%v]接口请求method:%v\t uri:%v\t param:%+v", user.Id, user.Account, r.GetClientIp(), r.Method, path, r.GetBodyString())
	}
	if r.IsFileRequest() {
		r.Middleware.Next()
		return
	}
	//判断是否登陆
	if IsSignedIn(r.Session) {
		//根据url判断是否有权限
		//url := r.Request.URL
		//if !IsFramePage(url.Path) {
		//	//获取用户信息
		//	user,err := models.GetUserOne(r.GetSessionId())
		//	if err != nil {
		//		r.Response.RedirectTo("/500")
		//		return
		//	}
		//	//获取用户菜单列表
		//	menus := models.GetMenuListByUserId(user.Id)
		//	if len(menus) == 0 {
		//		r.Response.RedirectTo("/500")
		//		return
		//	}
		//
		//	hasPermission := false
		//
		//	for i := range menus {
		//		if strings.EqualFold((menus)[i].Url, url.Path) {
		//			hasPermission = true
		//			break
		//		}
		//	}
		//
		//	if !hasPermission {
		//		ajaxString := r.Request.Header.Get("X-Requested-With")
		//		if strings.EqualFold(ajaxString, "XMLHttpRequest") {
		//			models.HttpResult(r, enums.CodeFail, "您没有操作权限", 0)
		//			return
		//		} else {
		//			r.Response.RedirectTo("/403")
		//			return
		//		}
		//	}
		//}
		if path != "/login" && path != "/logout" {
			checkLogin(user, r)
		}
		r.Middleware.Next()
		return
	}

	if !IsFramePage(path) && !strings.Contains(path, "statistic_res_opt") {
		HttpResult(r, 0, enums.LoginTypeNo, enums.CodeNoLogin, "未登录223", "")
		return
	}

	r.Middleware.Next()
	//g.Log().Warning("未登录跳转到初始界面")
	//RedirectToIndex(r)
	//r.Response.RedirectTo("/login")

}

// 判断用户是否已经登录
func IsSignedIn(session *ghttp.Session) bool {
	//session.Get(models.USER_SESSION_MARK)
	return session.Contains(models.USER_SESSION_MARK)
}

//判断是否是框架页面
func IsFramePage(path string) bool {
	for _, V := range FramePages {
		if strings.Index(path, V) != -1 {
			//if strings.EqualFold(V, path) {
			return true
		}
	}
	return false
}

// 获得用户信息详情
func GetSessionUser(r *ghttp.Request) (user models.User) {
	_ = r.Session.GetStruct(models.USER_SESSION_MARK, &user)
	return user
}

//检查是否登录
func checkLogin(user models.User, r *ghttp.Request) {
	var sessionId string
	sessionKey := gconv.String(user.Id)
	currSessionId := r.Session.Id()
	_ = utils.GetCacheCode(sessionKey, &sessionId)
	IsAllowMultipleLogin := g.Cfg().GetBool("server.isAllowMultipleLogin", false)
	//g.Log().Info("checkLogin:%s, %s",sessionId, c.CruSession.SessionID())
	if len(sessionId) == 0 {
		utils.SetCacheCode(sessionKey, currSessionId, 0)
	} else if IsAllowMultipleLogin == false {
		// 检测帐号 是否在多个地方登录
		if sessionId != currSessionId {
			g.Log().Debugf(" 帐号在其他地方登录:%+v", sessionId)
			HttpResult(r, user.Id, enums.LoginTypeKick, enums.CodeLoginOther, "帐号在其他地方登录:"+user.LastLoginIp, "")
			return
		}
	}

	u, err := models.GetUserOne(user.Id)
	if err == nil {
		if u.IsAccountEnable() == false {
			HttpResult(r, user.Id, enums.LoginTypeDisabled, enums.CodeLoginDisabled, "用户被禁用，请联系管理员", "")
			return
		}
	}
}

// http调用返回
func HttpResult(r *ghttp.Request, userId, loginType int, code enums.ResultCode, msg string, data interface{}) {
	//g.Log().Infof("http调用返回::%v  %v", r.Request.URL.Path, msg)
	models.SaveRequestAdminLog(userId, loginType, code, msg+": "+r.GetClientIp()+" host:"+r.GetHost(), r)
	result := &models.Result{Code: code, Data: data, Msg: msg}
	r.Response.WriteJsonExit(result)
}

func RedirectToIndex(r *ghttp.Request) {
	r.Response.RedirectTo("/login/index")
}
