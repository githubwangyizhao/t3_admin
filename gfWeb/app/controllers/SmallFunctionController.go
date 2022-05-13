package controllers

import (
	"encoding/json"
	"fmt"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SmallFunctionController struct {
	BaseController
}

// 邮件锁
func (c *SmallFunctionController) MailLockState(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
		State      int    `json:"state"`
	}

	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	g.Log().Infof("修改邮件锁:%+v", params)
	ResultList, ResultMsg, err := models.GameHttpRpcPlatform(params.PlatformId, params.ServerId, "mod_mail", "gm_mail_lock", "["+strconv.Itoa(params.State)+"]")
	//ResultList,err := models.GameHttpRpc("mod_mail", "gm_mail_lock", "[" + strconv.Itoa(params.State) + "]", realNodeList)
	if err != nil {
		//ResultStr := strings.Replace(strings.Trim(fmt.Sprint(ResultList), "[]"), " ", " ", -1)
		if ResultMsg == "" {
			ResultStr := strings.Trim(fmt.Sprint(ResultList), "[]")
			ResultMsg = "修改邮件锁失败列表:" + ResultStr
		}
		c.HttpResult(r, enums.CodeFail, ResultMsg, 0)
	}
	c.HttpResult(r, enums.CodeSuccess, "成功修改邮件锁", 0)
}

// 查询玩家Id平台数据
func (c *SmallFunctionController) PlayerIdPlatformList(r *ghttp.Request) {
	var params models.PlayerIdPlatformQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	g.Log().Infof("查询玩家Id平台数据:%+v", params)
	PlayerIdStr := strings.Replace(strings.Replace(params.PlayerIdStrList, "\n", "", -1), "\r\n", "", -1)
	line := strings.Split(PlayerIdStr, ",")
	len := len(line)
	if len == 0 {
		c.HttpResult(r, enums.CodeFail, "查询失败", 0)
	}
	g.Log().Infof("params.PlayerIdList:%+v", line)
	for _, PlayerId1 := range line {
		_, err := strconv.Atoi(PlayerId1)
		if err != nil {
			c.HttpResult(r, enums.CodeFail, "内容要以\"，\"分割的数字内容", 0)
		}
		//params.PlayerIdList = append(params.PlayerIdList, PlayerId)
	}
	params.PlayerIdList = line
	data, total := models.GetPlayerIdPlatformList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "查询成功", result)
}

// 发送短信
func (c *SmallFunctionController) SmeSend(r *ghttp.Request) {
	var params struct {
		PhoneStrList string `json:"phoneStrList"`
		PlatformType string `json:"platformType"`
		TemplateCode string `json:"templateCode"`
		MsgStr       string `json:"msgStr"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	g.Log().Infof("发送短信:%+v", params)
	PhoneNumberStr := strings.Replace(strings.Replace(params.PhoneStrList, "\n", "", -1), "\r\n", "", -1)
	line := strings.Split(PhoneNumberStr, ",")
	len := len(line)
	if len == 0 {
		c.HttpResult(r, enums.CodeFail, "发送短信失败", 0)
	}
	switch params.PlatformType {
	case "ALiYun":
		for _, PhoneNumber := range line {
			_, err := strconv.Atoi(PhoneNumber)
			if err != nil {
				utils.CheckError(err, "手机号只能是数字")
			}
			models.SendALiYunCode(PhoneNumber, params.TemplateCode, params.MsgStr)
		}
	case "YunPian":
		models.SendYunPian(PhoneNumberStr, params.TemplateCode, params.MsgStr)
	default:
		models.SendSms(PhoneNumberStr, params.TemplateCode, params.MsgStr)
	}
	c.HttpResult(r, enums.CodeSuccess, "发送成功", "")
}

// 获得客户端版本
func (c *SmallFunctionController) GetClientVersion(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:"platformId"`
		Channel    string `json:"channelId"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	g.Log().Infof("获得客户端版本:%+v", params)
	if params.PlatformId == "" {
		c.HttpResult(r, enums.CodeFail, "获得客户端版本失败", "平台不能为空")
	}
	ResultVersion := models.GetClientVersion(params.PlatformId, params.Channel)
	g.Log().Infof("获得客户端版本ResultVersion:%+v", ResultVersion)
	c.HttpResult(r, enums.CodeSuccess, "获得客户端版本成功", ResultVersion)
}

// 更新客户端版本
func (c *SmallFunctionController) UpdateClientVersion(r *ghttp.Request) {
	var params struct {
		PlatformId   string `json:"platformId"`
		Channel      string `json:"channelId"`
		VersionValue string `json:"versionValue"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if params.PlatformId == "" {
		c.HttpResult(r, enums.CodeFail, "获得客户端版本失败", "平台不能为空")
	}
	ResultVersion := models.UpdateClientVersion(params.PlatformId, params.Channel, params.VersionValue)
	c.HttpResult(r, enums.CodeSuccess, "获得客户端版本成功", ResultVersion)
}

// 后台发送邮件
func (c *SmallFunctionController) BackgroundMailSend(r *ghttp.Request) {
	var params struct {
		Title string `json:"title"`
		To    string `json:"to"`
		Body  string `json:"body"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	mailTo := strings.Split(params.To, ",")
	err = models.SendMailHandle(params.Title, params.Body, mailTo)
	if err != nil {
		g.Log().Errorf("邮件发送error: %v", err)
		c.HttpResult(r, enums.CodeFail, "后台发送邮件失败:"+gconv.String(err), params.To)
	}
	c.HttpResult(r, enums.CodeSuccess, "后台发送邮件成功", params.To)
}

// 更新订阅数据
func (c *SmallFunctionController) UpdatePlatformDingYueStatistics(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:"platformId"`
		UpdateTime int    `json:"updateTime"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	g.Log().Infof("更新订阅数据:%+v", params)
	models.UpdateDingYueStatistics(params.PlatformId, params.UpdateTime)
	c.HttpResult(r, enums.CodeSuccess, "更新订阅数据成功", "")
}

// 获得登陆服环境变更量
func (c *SmallFunctionController) GetEnvLoginServer(r *ghttp.Request) {
	var params struct {
		Url    string `json:"url"`
		EnvKey string `json:"envKey"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	//serverNodeList := models.GetAllServerNodeByType(4)
	if len(params.Url) == 0 {
		c.HttpResult(r, enums.CodeFail, "没有url", "")
	}
	//ServerNode := serverNodeList[0]
	args := fmt.Sprintf("evn_key=%s", params.EnvKey)
	sign := utils.String2md5(args + "fa9274fd68cf8991953b186507840e5e")
	//url := GetGameURLByNode(ServerNode.Node) + "/get_evn?sign=" + sign
	URL := params.Url + "/get_env"

	//初始化参数
	param := url.Values{}
	param.Set("env_key", params.EnvKey) // key
	param.Set("sign", sign)             // key
	BobyBin, err := utils.HttpGet(URL, param)
	if err != nil {
		g.Log().Errorf("获得登陆服环境变更量错误:", URL)
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("获得登陆服环境变更量错误:%v", err), "")
	}
	c.HttpResult(r, enums.CodeSuccess, "获得登陆服环境变更量成功", string(BobyBin))
}

// 设置登陆服环境变更量
func (c *SmallFunctionController) SetEnvLoginServer(r *ghttp.Request) {
	var params struct {
		Url         string `json:"url"`
		EnvKey      string `json:"envKey"`
		EnvValueStr string `json:"envValueStr"`
		EnvValue    string `json:"envValue"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	//serverNodeList := models.GetAllServerNodeByType(4)
	if len(params.Url) == 0 {
		c.HttpResult(r, enums.CodeFail, "没有url", "")
	}
	if params.EnvValueStr == "" {
		params.EnvValueStr = "_"
	}
	args := fmt.Sprintf("evn_key=%s&env_value_str=%s&env_value=%d", params.EnvKey, params.EnvValueStr, params.EnvValue)
	sign := utils.String2md5(args + "fa9274fd68cf8991953b186507840e5e")
	//url := GetGameURLByNode(ServerNode.Node) + "/get_evn?sign=" + sign

	//初始化参数
	param := url.Values{}
	param.Set("env_key", params.EnvKey)
	param.Set("env_value_str", params.EnvValueStr)
	//param.Set("env_value",strconv.Itoa(int(params.EnvValue)))
	param.Set("env_value", params.EnvValue)
	param.Set("sign", sign)
	//var resultStr = ""
	//for _, ServerNode := range serverNodeList  {
	//	URL := models.GetGameURLByNode(ServerNode.Node) + "/set_evn"
	//	BobyBin, err := utils.HttpGet(URL, param)
	//	if err != nil {
	//		g.Log().Error("设置登陆服环境变更量错误url:", URL, " param:", param)
	//		continue
	//	}
	//	resultStr = string(BobyBin)
	//}
	URL := params.Url + "/set_env"
	BobyBin, err := utils.HttpGet(URL, param)
	if err != nil {
		g.Log().Errorf("设置登陆服环境变更量错误:", URL)
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("设置登陆服环境变更量错误:%v", err), "")
	}
	c.HttpResult(r, enums.CodeSuccess, "设置登陆服环境变更量成功", string(BobyBin))
}

// 请求游戏服rpc
func (c *SmallFunctionController) RequestGameRpc(r *ghttp.Request) {
	var params struct {
		Type       int    `json:"type"`
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
		NodeStr    string `json:"nodeStr"`
		Mod        string `json:"mod"`
		Function   string `json:"function"`
		Args       string `json:"args"`
	}

	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	//g.Log().Infof("请求游戏服rpc:%+v", params)
	//var ResultMsg = ""
	//var ResultList [] string
	//if params.Type == 1 {
	ResultList, ResultMsg, err := models.GameHttpRpcPlatform(params.PlatformId, params.ServerId, params.Mod, params.Function, "["+params.Args+"]")
	//}else {
	//	ResultList, ResultMsg, err = models.GameHttpRpcNode(params.NodeStr, params.Mod, params.Function, "["+params.Args+"]")
	//}
	//ResultList,err := models.GameHttpRpc("mod_mail", "gm_mail_lock", "[" + strconv.Itoa(params.State) + "]", realNodeList)
	if err != nil {
		//ResultStr := strings.Replace(strings.Trim(fmt.Sprint(ResultList), "[]"), " ", " ", -1)
		//if ResultMsg == "" {
		ResultStr := strings.Trim(fmt.Sprint(ResultList), "[]")
		ResultMsg = ResultMsg + ":请求游戏服rpc失败列表:" + ResultStr
		//}
		c.HttpResult(r, enums.CodeFail, ResultMsg, "")
	}
	c.HttpResult(r, enums.CodeSuccess, "成功请求游戏服rpc", ResultMsg)
}

// 获得后台版本编号
func (c *SmallFunctionController) GetBackgroundUpdateVersion(r *ghttp.Request) {
	versionData, _ := models.GetBackgroundVersion()
	c.HttpResult(r, enums.CodeSuccess, "获得后台版本编号成功", versionData)
}

// 后台更新版本
func (c *SmallFunctionController) BackgroundUpdateVersion(r *ghttp.Request) {
	err := models.UpdateBackgroundVersion()
	c.CheckError(err, "后台更新版本失败")
	c.HttpResult(r, enums.CodeSuccess, "后台更新版本成功", "")
}

// 停止后台
func (c *SmallFunctionController) StopBackgroundUpdateVersion(r *ghttp.Request) {
	err := models.StopBackground()
	c.CheckError(err, "停止失败")
	c.HttpResult(r, enums.CodeSuccess, "停止成功", "")
}

// 从leanCloud获取最新的登陆页客服链接并更新到游戏中心服的ets中
func (c *SmallFunctionController) UpdateLoginCustomerServiceUrl(r *ghttp.Request) {
	g.Log().Info("updateLoginCustomerServiceUrl")

	CenterUrl := utils.GetCenterURL() + "/customer_service/update_login_page_url"
	g.Log().Info("url: ", CenterUrl)
	var request struct {
		Timestamp int `json:"timestamp"`
	}

	request.Timestamp = int(time.Now().Unix())
	g.Log().Info("request: ", request)
	data2, err := json.Marshal(request)
	utils.CheckError(err)

	Res, err := utils.HttpRequest(CenterUrl, string(data2))
	g.Log().Info("Res: ", Res)

	c.HttpResult(r, enums.CodeSuccess, "更新成功", "")
}
