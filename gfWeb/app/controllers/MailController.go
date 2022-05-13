package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"strconv"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	//"net/http"
	//"io/ioutil"
	//"strings"
	//"encoding/base64"
)

type MailController struct {
	BaseController
}

// 获取邮件列表
func (c *MailController) MailLogList(r *ghttp.Request) {
	var params models.MailLogQueryParam
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Info("查询邮件日志:%+v", params)
	data, total := models.GetMailLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取邮件日志", result)
}

// 删除邮件
func (c *MailController) DelMailLog(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)

	idList := params
	utils.CheckError(err)
	g.Log().Info("删除邮件列表:%+v", idList)
	err = models.DeleteMailLog(idList)
	c.CheckError(err, "删除邮件失败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除邮件", idList)
}

//发送邮件
func (c *MailController) SendMail(r *ghttp.Request) {

	var params struct {
		PlatformId      string        `json:"platformId"`
		NodeList        []string      `json:"serverIdList"`
		PlayerNameList  string        `json:"playerNameList"`
		MailItemList    []models.Prop `json:"mailItemList"`
		Title           string        `json:"title"`
		Content         string        `json:"content"`
		PlayerIdList    []int         `json:"playerIdList"`
		Type            string        `json:"type"`
		ConditionsId    int           `json:"conditionsId"`
		ConditionsValue string        `json:"conditionsValue"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Info("发送邮件:%+v", params)

	if params.Type != "1" && params.Type != "2" && params.Type != "3" {
		c.HttpResult(r, enums.CodeFail, "邮件类型错误!!!", 0)
	}

	params.PlayerIdList = make([]int, 0)
	if params.Type == "1" {
		// 发送给玩家
		if params.PlayerNameList != "" && len(params.NodeList) != 1 {
			g.Log().Error("参数错误!!!!")
			c.HttpResult(r, enums.CodeFail, "参数错误!!!", 0)
		}
		if params.PlayerNameList != "" {
			params.PlayerIdList, err = models.TranPlayerNameSting2PlayerIdList(params.PlatformId, params.PlayerNameList)
			c.CheckError(err, "解析玩家失败")
		}
	} else if params.Type == "2" {
		// 发送给多个服
		if params.PlayerNameList != "" || len(params.NodeList) == 0 {
			g.Log().Error("参数错误!!!!")
			c.HttpResult(r, enums.CodeFail, "参数错误!!!", 0)
		}
	} else if params.Type == "3" {
		// 发送全服
		if params.PlayerNameList != "" || len(params.NodeList) > 0 {
			g.Log().Error("参数错误!!!!")
			c.HttpResult(r, enums.CodeFail, "参数错误!!!", 0)
		}
	}

	if params.ConditionsId == 1 { // vip等级
		_, err := strconv.Atoi(params.ConditionsValue)
		utils.CheckError(err)
		if err != nil {
			g.Log().Error("vip等级参数错误!!!!")
			c.HttpResult(r, enums.CodeFail, "vip等级参数错误!!!", 0)
		}
	} else if params.ConditionsId == 0 {
		params.ConditionsValue = ""
	}

	request, err := json.Marshal(params)
	c.CheckError(err)

	nodeList, err := json.Marshal(params.NodeList)
	c.CheckError(err)

	itemList, err := json.Marshal(params.MailItemList)
	c.CheckError(err)

	mailLog := &models.MailLog{
		PlatformId:     params.PlatformId,
		NodeList:       string(nodeList),
		Title:          params.Title,
		Content:        params.Content,
		Time:           time.Now().Unix(),
		UserId:         c.curUser.Id,
		ItemList:       string(itemList),
		PlayerNameList: params.PlayerNameList,
		Type:           params.Type,
		Status:         0,
	}

	err = models.Db.Save(&mailLog).Error
	c.CheckError(err, "写邮件日志失败")

	realNodeList := make([]string, 0)
	if params.Type == "3" {
		realNodeList = models.GetAllGameNodeByPlatformId(params.PlatformId)
	} else {
		realNodeList = models.GetNodeListByServerIdList(params.PlatformId, params.NodeList)
		//realNodeList = params.NodeList
	}
	for _, node := range realNodeList {
		url := models.GetGameURLByNode(node) + "/send_mail"
		data := string(request)
		_, err = utils.HttpRequest(url, data)
		if err != nil {
			g.Log().Errorf("发送邮件失败 区服:%s, 标题:%s, 内容:%s, 玩家:%v,条件id:%d,条件内容:%s", node, params.Title, params.Content, params.PlayerNameList, params.ConditionsId, params.ConditionsValue)
		}
		//sign := utils.String2md5(data + enums.GmSalt)
		//base64Data := base64.URLEncoding.EncodeToString([]byte(data))
		//requestBody := "data=" + base64Data + "&sign=" + sign
		//resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(requestBody))
		//utils.CheckError(err)
		//if err != nil {
		//	g.Log().Info("发送邮件失败 类型:%s 区服:%s, 标题:%s, 内容:%s, 玩家:%v", params.Type, node, params.Title, params.Content, params.PlayerNameList)
		//	continue
		//}
		//
		//defer resp.Body.Close()
		//responseBody, err := ioutil.ReadAll(resp.Body)
		//
		//if err != nil {
		//	g.Log().Info("发送邮件失败 类型:%s 区服:%s, 标题:%s, 内容:%s, 玩家:%v", params.Type, node, params.Title, params.Content, params.PlayerNameList)
		//	continue
		//}
		//
		////g.Log().Debug("result:%v", string(responseBody))
		//
		//err = json.Unmarshal(responseBody, &result)
		//
		////c.CheckError(err)
		//if result.ErrorCode != 0 || err != nil {
		//	g.Log().Info("发送邮件失败 区服:%s, 标题:%s, 内容:%s, 玩家:%v", node, params.Title, params.Content, params.PlayerNameList)
		//} else {
		//	g.Log().Info("发送邮件成功 区服:%s, 标题:%s, 内容:%s, 玩家:%v", node, params.Title, params.Content, params.PlayerNameList)
		//}

	}
	//args := fmt.Sprintf("platform_id=%d&server_id=%s&player_id=%d&type=%d&sec=%d", params.PlatformId, params.ServerId, params.PlayerId, params.Type, params.Sec)
	//sign := utils.String2md5(args + enums.GmSalt)

	//if result.ErrorCode != 0 {
	//	c.HttpResult(r, enums.CodeFail, "发送邮件失败", 0)
	//}

	//for _, node := range params.NodeList {
	//	conn, err := models.GetWsByNode(node)
	//	c.CheckError(err)
	//	defer conn.Close()
	//	request := gm.MSendMailTos{
	//		Token:          proto.String(""),
	//		Title:          proto.String(params.Title),
	//		Content:        proto.String(params.Content),
	//		PlayerNameList: proto.String(params.PlayerNameList),
	//		PropList:       params.MailItemList,
	//	}
	//	mRequest, err := proto.Marshal(&request)
	//	c.CheckError(err)
	//
	//	_, err = conn.Write(utils.Packet(9903, mRequest))
	//	c.CheckError(err)
	//	var receive = make([]byte, 100, 100)
	//	n, err := conn.Read(receive)
	//	c.CheckError(err)
	//	respone := &gm.MSendMailToc{}
	//	data := receive[5:n]
	//	err = proto.Unmarshal(data, respone)
	//	c.CheckError(err)
	//
	//	if *respone.Result == gm.MSendMailToc_success {
	//		g.Log().Info("发送邮件成功:%+v, %+v", node, request)
	//	} else {
	//		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("发送邮件失败:%+v, %+v", node, request), 0)
	//	}
	//}
	c.HttpResult(r, enums.CodeSuccess, "发送邮件成功", 0)
}
