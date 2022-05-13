package models

import (
	"encoding/json"
	"gfWeb/library/utils"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var smsTimeInterval = make(map[string]int64) // 同数据内容的短信发送间隙

var timeInterval = int64(600) // 间隔10分钟

// 管理员短信
func SendSmsAdminUser(TemplateCode, TemplateParam string) {
	PlatformType := g.Cfg().GetString("game.platform_type")
	IsOpenSms := IsSettingOpenDefault(SETTING_DATA_IS_OPEN_SMS, false)
	g.Log().Infof("成员短信code:%s 参数%s 平台:%s", TemplateCode, TemplateParam, PlatformType)
	if IsOpenSms == false {
		g.Log().Infof("普通短信通知关闭!! %s:%s", TemplateCode, TemplateParam)
		return
	}
	SendSmsAdminUserStr(TemplateCode, TemplateParam)
	//if PlatformType == "qq" || PlatformType == "wx" || PlatformType == "jkzx" || PlatformType == "gat" {
	//}
}

// 管理员短信直接发送
func SendSmsAdminUserMsg(TemplateCode, TemplateParam string) {
	g.Log().Infof("管理员短信直接发送code:%s >>%s", TemplateCode, TemplateParam)
	SendSmsAdminUserStr(TemplateCode, TemplateParam)
}
func SendSmsAdminUserStr(TemplateCode, TemplateParam string) {
	PhoneNumberStr := g.Cfg().GetString("game.sms_phone_number")
	Key := TemplateCode + TemplateParam
	var currTime = time.Now().Unix()
	OldTime, ok := smsTimeInterval[Key]
	if ok {
		if OldTime+timeInterval > currTime {
			return
		}
	}
	smsTimeInterval[Key] = currTime
	SendSms(PhoneNumberStr, TemplateCode, TemplateParam)
}

// 发送短信 PhoneNumberStr:手机号（多个请有","分割）
// TemplateCode:模版   TemplateParam:模板参数
func SendSms(PhoneNumberStr, TemplateCode, TemplateParam string) {
	PhoneList := strings.Split(PhoneNumberStr, ",")
	SendPhoneMsgByPhoneList(TemplateCode, TemplateParam, PhoneList)
}

// 发送短信手机号列表
func SendPhoneMsgByPhoneList(TemplateCode, TemplateParam string, PhoneList []string) {
	IsOpenAllSms := IsSettingOpenDefault(SETTING_DATA_IS_OPEN_ALL_SMS, false)
	if IsOpenAllSms == false {
		g.Log().Infof("短信通知全部关闭!! %s:%s", TemplateCode, TemplateParam)
		return
	}
	//PhoneList := strings.Split(PhoneNumberStr, ",")
	for _, PhoneNumber := range PhoneList {
		_, err := strconv.Atoi(PhoneNumber)
		if err != nil {
			utils.CheckError(err, "手机号只能是数字")
		}
		// 要分条发
		err = SendALiYunCode(PhoneNumber, TemplateCode, TemplateParam)
		if err != nil {
			utils.CheckError(err, "发送短信失败err:%s:%s", TemplateCode, TemplateParam)
			continue
		}
		//SendJuHe(TemplateCode, PhoneNumber)
	}
	// 不用分条发
	//SendYunPian(PhoneNumberStr, TemplateCode, TemplateParam)
}

// 阿里云发数短信code
// PhoneNumbers 支持对多个手机号码发送短信，手机号码之间以英文逗号（,）分隔。上限为1000个手机号码。批量调用相对于单条调用及时性稍有延迟。
func SendALiYunCode(PhoneNumbers, TemplateCode, TemplateParam string) error {
	client, err := sdk.NewClientWithAccessKey("", "", "")
	if err != nil {
		utils.CheckError(err, "阿里云sdk失败")
		return err
	}
	TemplateParam = "{\"code\":\"" + TemplateParam + "\"}"
	g.Log().Infof("阿里云短信参数:%+v", TemplateParam)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = PhoneNumbers
	request.QueryParams["TemplateParam"] = TemplateParam
	request.QueryParams["TemplateCode"] = TemplateCode
	request.QueryParams["SignName"] = "修仙"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		utils.CheckError(err, "阿里云发送请求失败")
		return err
	}
	ResultJson := response.GetHttpContentString()
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(ResultJson), &dat); err != nil {
		g.Log().Infof("阿里云发数短信json失败:%+v", ResultJson)
		return err
	}
	Code := dat["Code"]
	if Code != "OK" {
		g.Log().Error("%s:阿里云发数短信失败%s msg:%+v", PhoneNumbers, Code, dat["Message"])
		return gerror.New("发送失败")
	} else {
		g.Log().Infof("阿里云发数短信成功:%s", PhoneNumbers)
	}
	return nil
}

// 聚合短信
func SendJuHe(msgId string, phone string) error {
	//请求地址
	juheURL := "http://v.juhe.cn/sms/send"
	//初始化参数
	param := url.Values{}
	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("mobile", phone)                           //接收短信的手机号码
	param.Set("tpl_id", msgId)                           //短信模板ID，请参考个人中心短信模板设置
	param.Set("tpl_value", "")                           //变量名和变量值对。如果你的变量名或者变量值中带有#&amp;=中的任意一个特殊符号，请先分别进行urlencode编码后再传递，&lt;a href=&quot;http://www.juhe.cn/news/index/id/50&quot; target=&quot;_blank&quot;&gt;详细说明&gt;&lt;/a&gt;
	param.Set("key", "") //应用APPKEY(应用详细页查询)
	param.Set("dtype", "")                               //返回数据的格式,xml或json，默认json

	//发送请求
	data, err := utils.HttpGet(juheURL, param)
	if err != nil {
		g.Log().Error("聚合短信请求失败,错误信息:\r\n%v", err)
		return err
	} else {
		var netReturn map[string]interface{}
		json.Unmarshal(data, &netReturn)
		g.Log().Infof("聚合短信结果:%+v", netReturn)
	}
	return nil
}

// 云片短信  	PhoneNumbers:要发送的手机号码，多个号码用逗号隔开
func SendYunPian(PhoneNumbers, TemplateCode, TemplateParam string) error {
	//请求地址
	smsURL := "https://sms.yunpian.com/v2/sms/tpl_single_send.json"
	// 修改为您的apikey(https://www.yunpian.com)登录官网后获取
	apikey := ""
	tpl_value := url.Values{"#code#": {TemplateParam}}.Encode()
	param := url.Values{"apikey": {apikey}, "mobile": {PhoneNumbers},
		"tpl_id": {TemplateCode}, "tpl_value": {tpl_value}}
	g.Log().Error("云片短信参数信息:%v", param)
	////发送请求
	data, err := HttpPostForm(smsURL, param)
	if err != nil {
		g.Log().Error("云片短信请求失败,错误信息:\r\n%v", err)
		return err
	} else {
		var resultJson map[string]interface{}
		json.Unmarshal(data, &resultJson)
		Code := resultJson["code"]
		Phones := resultJson["mobile"]
		if Code == 0. {
			g.Log().Infof("云片短信发送成功:%+v", Phones)
		} else {
			msg := resultJson["msg"]
			g.Log().Infof("云片短信发送失败Code:%+v,msg:%s, phones:%s", Code, msg, Phones)
			return gerror.New("云片短信发送失败")
		}
	}
	return nil
}

// http post参数请求
func HttpPostForm(url string, requestBody url.Values) ([]byte, error) {
	var responseByte []byte
	resp, err := http.PostForm(url, requestBody)
	utils.CheckError(err)
	if err != nil {
		return responseByte, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	utils.CheckError(err)
	if err != nil {
		return responseBody, err
	}
	return responseBody, err
}

// 检测平台停机维护 20分钟未开启机器时,警报
func CheckPlatformServerClose() {
	currTimeMicro := utils.GetTimestampMicro()
	var CheckTime = currTimeMicro - 1200*1000
	for _, platformData := range GetPlatformSimpleList() {
		platformId := platformData.Id
		key := GetUpdatePlatformVersionCacheKey(platformId)
		platformValue := utils.GetCacheInt64(key)
		if platformValue > 0 && platformValue < CheckTime {
			SendBackgroundMsgTemplateHandle(MSG_TEMPLATE_SERVER_CLOSE_LONG_TIME, platformId, platformId+"平台停机维护过长", "游戏服已经维护时间："+utils.FormatTimeSecond(gconv.Int((currTimeMicro-platformValue)/1000)))
		}
	}
}
