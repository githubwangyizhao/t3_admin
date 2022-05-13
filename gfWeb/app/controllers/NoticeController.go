// 公告管理
package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"gfWeb/memdb"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/util/gconv"
)

type NoticeController struct {
	BaseController
}

// 发送定时公告
func (c *NoticeController) SendCronNotice(r *ghttp.Request) {
	var params models.SendNoticeParams
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	//g.Log().Info("发送定时公告:%+v", params)
	c.CheckError(err)
	//if params.NoticeType != enums.NoticeTypeMoment  &&  params.NoticeType != enums.NoticeTypeClock  &&  params.NoticeType != enums.NoticeTypeLoop {
	//	c.HttpResult(r, enums.CodeFail, "公告类型错误", 0)
	//}

	err = models.CreateNoticeLog(params, c.curUser)
	c.CheckError(err, "写公告日志失败")
	c.HttpResult(r, enums.CodeSuccess, "发送公告成功", 0)
}

// 移除定时公告
func (c *NoticeController) RemoveCronNotice(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)
	g.Log().Infof("移除定时公告:%+v", params)
	c.CheckError(err)
	err = models.RemoveCronNotice(params, c.curUser)
	c.CheckError(err, "写公告日志失败")
	c.HttpResult(r, enums.CodeSuccess, "发送公告成功", 0)
}

//// 发送公告
//func (c *NoticeController) SendNotice(r *ghttp.Request) {
//	var params models.SendNoticeParams
//	err := json.Unmarshal(r.GetBody(), &params)
//
//	utils.CheckError(err)
//	g.Log().Info("发送公告:%+v", params)
//	ServerIdList, err := json.Marshal(params.ServerIdList)
//	c.CheckError(err)
//	if params.NoticeType != enums.NoticeTypeMoment && params.NoticeType != enums.NoticeTypeClock && params.NoticeType != enums.NoticeTypeLoop {
//		c.HttpResult(r, enums.CodeFail, "公告类型错误", 0)
//	}
//	err = models.CreateNoticeLog(params, c.curUser)
//	noticeLog := &models.NoticeLog{
//		Id:           params.Id,
//		PlatformId:   params.PlatformId,
//		ServerIdList: string(ServerIdList),
//		IsAllServer:  params.IsAllServer,
//		Content:      params.Content,
//		Time:         gtime.Timestamp(),
//		UserId:       c.curUser.Id,
//		NoticeType:   params.NoticeType,
//		NoticeTime:   params.NoticeTime,
//		Status:       0,
//	}
//	err = models.Db.Save(&noticeLog).Error
//	c.CheckError(err, "写公告日志失败")
//	// 异步处理日志
//	if params.NoticeType == enums.NoticeTypeMoment {
//		go models.DealNoticeLog(noticeLog.Id)
//	}
//	c.HttpResult(r, enums.CodeSuccess, "发送公告成功", 0)
//}

// 获取公告列表
func (c *NoticeController) NoticeLogList(r *ghttp.Request) {
	var params models.NoticeLogQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)

	g.Log().Infof("查询公告日志:%+v", params)
	data, total := models.GetNoticeLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取公告日志", result)
}

//删除公告
func (c *NoticeController) DelNoticeLog(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	idList := params
	utils.CheckError(err)
	g.Log().Infof("删除公告列表:%+v", idList)
	err = models.DeleteNoticeLog(idList)
	c.CheckError(err, "删除公告失败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除公告", idList)
}

// 极光推送
func (c *NoticeController) JiguangPush(r *ghttp.Request) {

	var params models.JiguangPushParam
	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)

	if params.Desc == "" {
		c.HttpResult(r, enums.CodeFail, "内容不能为空", nil)
		return
	}

	var functionConfigId = ""
	if params.IsOpenUi == 1 {
		functionConfigId = "0"
	} else {
		functionConfigId = gconv.String(params.FunctionConfigId)
	}

	// 获取registrationIds
	// var registrationIds []string

	// registrationIds = memdb.FetchDataBySids(params.ServerIdList)

	var pushList = make(map[string][]string)

	data := memdb.FetchDataBySids(params.Platform, params.ServerIdList)
	for areaSid, subData := range data {
		// fmt.Println(areaSid, item["open"])
		pushList[areaSid] = make([]string, 0)
		for _, item := range subData {
			if item["open"] == "1" && item["regis"] != "" && item["regis"] != "undefined" {
				// fmt.Println(item["regis"])
				pushList[areaSid] = append(pushList[areaSid], item["regis"])
			}
		}

	}

	if params.CronTimeStr != "" {
		g.Log().Warning("定时发送 %s ", params.CronTimeStr)
		cronFunc := func() {
			_pushArea(params.Headline, params.Desc, functionConfigId, pushList)
		}
		_, err = gcron.AddTimes(params.CronTimeStr, 1, cronFunc, "push_"+utils.RandomString(20))
		utils.CheckError(err)
	} else {
		_pushArea(params.Headline, params.Desc, functionConfigId, pushList)
	}

	utils.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功", nil)
}

func _pushArea(title, desc, functionConfigId string, pushList map[string][]string) {
	for sid, area := range pushList {
		if len(area) > 0 {
			//推送到队列
			area = append(area, title)
			area = append(area, desc)
			area = append(area, functionConfigId)
			area = append(area, sid)
			PUSHQUEUE.Push(area)
		}
	}
}

// func (c *NoticeController) doJgPush(title, desc, functionConfigId string, registrationIds []string) bool {
// 	// 推送开始
// 	// var registrationIds []string
// 	// registrationIds = append(registrationIds, "\""+"120c83f760fb51a5110"+"\"")
// 	// registrationIds = append(registrationIds, "120c83f760fb51a5110") //todo 在数据库player_global中查找到所有的
// 	var (
// 		url         = g.Cfg().GetString("jiguang.url")
// 		method      = "POST"
// 		contentType = functionConfigId
// 	)
// 	title = utils.IfStr(title, g.Cfg().GetString("jiguang.gamename")) //headline中取值

// 	var payLoadStrMap = map[string]interface{}{
// 		"platform": "all",
// 		"audience": map[string][]string{
// 			"registration_id": registrationIds,
// 		},
// 		"notification": map[string]map[string]string{
// 			"android": map[string]string{
// 				"title": title,
// 				"alert": desc,
// 			},
// 			"ios": map[string]string{
// 				"alert": desc,
// 			},
// 		},
// 		"message": map[string]string{
// 			"title":        title,
// 			"msg_content":  desc,
// 			"content_type": gconv.String(contentType),
// 		},
// 	}

// 	payLoadbyte, _ := json.Marshal(payLoadStrMap)
// 	payload := strings.NewReader(string(payLoadbyte))

// 	// fmt.Println(payLoadStr)
// 	client := &http.Client{}
// 	req, err := http.NewRequest(method, url, payload)

// 	if err != nil {
// 		utils.CheckError(err)
// 		return false
// 	}

// 	authValue := base64.StdEncoding.EncodeToString([]byte(g.Cfg().GetString("jiguang.appKey") + ":" + g.Cfg().GetString("jiguang.masterSecret")))
// 	req.Header.Add("Authorization", "Basic "+authValue)
// 	req.Header.Add("Content-Type", "application/json")

// 	res, err := client.Do(req)
// 	if err != nil {
// 		utils.CheckError(err)
// 		return false
// 	}
// 	defer res.Body.Close()

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		utils.CheckError(err)
// 		return false
// 	}

// 	type RetError struct {
// 		Code    int    `json:"code"`
// 		Message string `json:"message"`
// 	}
// 	var responseBody = struct {
// 		Error RetError `json:"error"`
// 	}{}
// 	// fmt.Println(string(body))
// 	json.Unmarshal(body, &responseBody)

// 	return responseBody.Error.Code == 0
// }

func (c *NoticeController) GetJgPushItem(r *ghttp.Request) {
	// 核对充值金额 ，获取经验值
	url := g.Cfg().GetString("game.gameCenterHost") + "/static/json/NewsPush.json"
	expDataArr := utils.HttpGetJsonSliceMap(url)

	c.HttpResult(r, enums.CodeSuccess, "成功", expDataArr)
}

//
func (c *NoticeController) GetJgPushFunction() []map[string]interface{} {
	// 核对充值金额 ，获取经验值
	url := g.Cfg().GetString("game.gameCenterHost") + "/static/json/Function.json"
	expDataArr := utils.HttpGetJsonSliceMap(url)

	return expDataArr
}
