package controllers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"gfWeb/memdb"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

type ToolController struct {
	BaseController
}

func (c *ToolController) GetJson() {
	list := models.GetPlatformList()
	//utils.CheckError(err)
	g.Log().Infof("platformList:%v", list)
	c.HttpResult(c.Request, enums.CodeSuccess, "获取平台列表成功 ", list)
}

func (c *ToolController) Action(r *ghttp.Request) {
	var params struct {
		Action     string `json:"action"`
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("Action:%+v", params)
	//node := getNode(params.ServerId)
	//gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	//c.CheckError(err)
	//node := gameServer.Node
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	commandArgs := []string{
		"nodetool",
		"-name",
		gameServer.Node,
		"-setcookie",
		"game",
		"rpc",
		"tool",
		"project_helper",
		params.Action,
	}
	out, err := utils.CmdNodetool("escript", commandArgs)

	if err != nil {
		out = strings.Replace(out, " ", "&nbsp", -1)
		out = strings.Replace(out, "\n", "<br>", -1)
		out = strings.Replace(out, "\\n", "<br>", -1)
		c.HttpResult(r, enums.CodeFail2, "失败:"+out+err.Error(), 0)
	} else {
		//if params.Action == "build_table" {
		//	commandArgs = []string{
		//		"ci",
		//		"/opt/h5/trunk/client/client_enum",
		//		"-m",
		//		"web_submit",
		//	}
		//	_, err = utils.Cmd("svn", commandArgs)
		//	c.CheckError(err, "提交客户端枚举")
		//}
		c.HttpResult(r, enums.CodeSuccess, "成功!", 0)
	}
}

func (c *ToolController) SendProp(r *ghttp.Request) {
	var params struct {
		PlayerId   int    `json:"playerId"`
		PropType   int    `json:"propType"`
		PropId     int    `json:"propId"`
		PropNum    int    `json:"propNum"`
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("Action:%+v", params)
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	//node := getNode(params.ServerId)
	//gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	//c.CheckError(err)
	//node := gameServer.Node
	commandArgs := []string{
		"nodetool",
		"-name",
		gameServer.Node,
		"-setcookie",
		"game",
		"rpc",
		"tool",
		"give_prop",
		strconv.Itoa(params.PlayerId),
		strconv.Itoa(params.PropType),
		strconv.Itoa(params.PropId),
		strconv.Itoa(params.PropNum),
	}
	out, err := utils.CmdNodetool("escript", commandArgs)

	if err != nil {
		out = strings.Replace(out, " ", "&nbsp", -1)
		out = strings.Replace(out, "\n", "<br>", -1)
		out = strings.Replace(out, "\\n", "<br>", -1)
		c.HttpResult(r, enums.CodeFail2, "失败:"+out+err.Error(), 0)
	} else {
		c.HttpResult(r, enums.CodeSuccess, "成功!", 0)
	}
}

func (c *ToolController) SetTask(r *ghttp.Request) {
	var params struct {
		PlayerId   int    `json:"playerId"`
		TaskId     int    `json:"taskId"`
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("Action:%+v", params)
	//node := getNode(params.ServerId)
	//gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	//c.CheckError(err)
	//node := gameServer.Node
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	commandArgs := []string{
		"nodetool",
		"-name",
		gameServer.Node,
		"-setcookie",
		"game",
		"rpc",
		"tool",
		"debug_set_task",
		strconv.Itoa(params.PlayerId),
		strconv.Itoa(params.TaskId),
	}
	out, err := utils.CmdNodetool("escript", commandArgs)

	if err != nil {
		out = strings.Replace(out, " ", "&nbsp", -1)
		out = strings.Replace(out, "\n", "<br>", -1)
		out = strings.Replace(out, "\\n", "<br>", -1)
		c.HttpResult(r, enums.CodeFail2, "失败:"+out+err.Error(), 0)
	} else {
		c.HttpResult(r, enums.CodeSuccess, "成功!", 0)
	}
}

func (c *ToolController) FinishBranchTask(r *ghttp.Request) {
	var params struct {
		PlayerId   int    `json:"playerId"`
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("Action:%+v", params)
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	commandArgs := []string{
		"nodetool",
		"-name",
		gameServer.Node,
		"-setcookie",
		"game",
		"rpc",
		"tool",
		"finish_branch_task",
		strconv.Itoa(params.PlayerId),
	}
	out, err := utils.CmdNodetool("escript", commandArgs)

	if err != nil {
		out = strings.Replace(out, " ", "&nbsp", -1)
		out = strings.Replace(out, "\n", "<br>", -1)
		out = strings.Replace(out, "\\n", "<br>", -1)
		c.HttpResult(r, enums.CodeFail2, "失败:"+out+err.Error(), 0)
	} else {
		c.HttpResult(r, enums.CodeSuccess, "成功!", 0)
	}
}

/**
 * 获取设置服务器时间
 * @request time format string
 *  case  -1  获取服务器时间
 *  case   0  设置北京时间
 *  case  >1  设置为time时间，此时time必须传入10位时间截字符串进行请求
 */
func (c *ToolController) ServerTime(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
		Time       string `json:"time"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)

	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	serverNode, err := models.GetServerNode(gameServer.Node)
	c.CheckError(err)
	var (
		httpUrl = "http://" + serverNode.Ip + ":" + gconv.String(serverNode.WebPort)
	)

	// 只支持10位时间截
	if len(params.Time) == 10 || params.Time == "0" {
		httpUrl += "/set_server_time"

		// 由于是form-data不传文件特例，这边不做封装，只在这里写一次
		method := "POST"

		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("time", params.Time)
		err := writer.Close()
		c.CheckError(err)

		client := &http.Client{}
		req, err := http.NewRequest(method, httpUrl, payload)
		c.CheckError(err)

		req.Header.Set("Content-Type", writer.FormDataContentType())
		res, err := client.Do(req)
		c.CheckError(err)
		defer res.Body.Close()

		_, err = ioutil.ReadAll(res.Body)
		c.CheckError(err)

		setTime := ""
		if params.Time == "0" {
			setTime = utils.TimeInt64FormDefault(time.Now().Unix())
		} else {
			setTime = utils.TimeInt64FormDefault(gconv.Int64(params.Time))
		}
		c.HttpResult(r, enums.CodeSuccess, "成功设置服务器时间为: "+setTime, "")
	} else if params.Time == "-1" {
		httpUrl += "/get_server_time"
		rs, err := utils.HttpGet(httpUrl, nil)
		c.CheckError(err)

		var reMap = make(map[string]interface{})
		err = json.Unmarshal(rs, &reMap)
		c.CheckError(err)

		// c.HttpResult(r, enums.CodeSuccess, "成功!", reMap["Msg"])
		c.HttpResult(r, enums.CodeSuccess, "成功!", utils.TimeStrToStamp(gconv.String(reMap["Msg"])))

	}

}

func (c *ToolController) ActiveFunction(r *ghttp.Request) {
	var params struct {
		PlayerId      int    `json:"playerId"`
		FunctionId    int    `json:"functionId"`
		FunctionParam int    `json:"functionParam"`
		FunctionValue int    `json:"functionValue"`
		PlatformId    string `json:"platformId"`
		ServerId      string `json:"serverId"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("Action:%+v", params)
	//node := getNode(params.ServerId)
	//gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	//c.CheckError(err)
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	//node := gameServer.Node
	commandArgs := []string{
		"nodetool",
		"-name",
		gameServer.Node,
		"-setcookie",
		"game",
		"rpc",
		"tool",
		"active_function",
		strconv.Itoa(params.PlayerId),
		strconv.Itoa(params.FunctionId),
		strconv.Itoa(params.FunctionParam),
		strconv.Itoa(params.FunctionValue),
	}
	out, err := utils.CmdNodetool("escript", commandArgs)

	if err != nil {
		out = strings.Replace(out, " ", "&nbsp", -1)
		out = strings.Replace(out, "\n", "<br>", -1)
		out = strings.Replace(out, "\\n", "<br>", -1)
		c.HttpResult(r, enums.CodeFail2, "失败:"+out+err.Error(), 0)
	} else {
		c.HttpResult(r, enums.CodeSuccess, "成功!", 0)
	}
}

func (c *ToolController) GetIpOrigin(r *ghttp.Request) {
	var params struct {
		Ip string `json:"Ip"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	o := utils.GetIpLocation(params.Ip)
	//g.Log().Debug("ip origin:%v %v", params.Ip, o)
	if err != nil {
		c.HttpResult(r, enums.CodeFail2, "失败:", "")
	} else {
		c.HttpResult(r, enums.CodeSuccess, "成功!", o)
	}
}

func (c *ToolController) Merge(r *ghttp.Request) {
	var params struct {
		PlatformId string                   `json:"platformId"`
		MergeTime  int                      `json:"mergeTime"`
		MergeList  []models.MergeServerData `json:"mergeList"`
		//MergeList [] struct {
		//	S int
		//	E int
		//} `json:"mergeList"`
	}

	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	g.Log().Infof("请求合服%s:%+v", c.curUser.Name, params)

	err = models.CreatePlatformMergeData(c.curUser.Id, c.curUser.Name, params.PlatformId, params.MergeTime, params.MergeList)
	if err != nil {
		c.HttpResult(r, enums.CodeFail2, fmt.Sprint("合服失败:%+v:", err), "")
	} else {
		c.HttpResult(r, enums.CodeSuccess, "合服成功!", "")
	}
	////g.Log().Debug("ip origin:%v %v", params.Ip, o)
	//nodeList := make([]*models.ServerNode, 0)
	//for _, e := range params.Nodes {
	//	serverNode, err := models.GetServerNode(e)
	//	c.CheckError(err)
	//	nodeList = append(nodeList, serverNode)
	//}
	//_, err = models.GetServerNode(params.ZoneNode)
	//c.CheckError(err)
	//models.UpdateMergeState(true, params.PlatformId + strconv.Itoa(len(params.MergeList )))
	//var str [] string
	//for _, e := range params.MergeList {
	//	str = append(str, strconv.Itoa(e.S) + "-" +  strconv.Itoa(e.E))
	//}
	//models.UpdateEtsPlatformServerMergeState(params.PlatformId, str, true, true)
	//IsResult := false
	//defer func() {
	//	models.UpdateEtsPlatformServerMergeState(params.PlatformId, str, false, IsResult)
	//}()
	//for _, e := range params.MergeList {
	//	g.Log().Info("执行合服:%s, %+v", params.PlatformId, e)
	//	nodeList, zoneNode, err := getMergeInfo(params.PlatformId, e.S, e.E)
	//	c.CheckError(err, "合服失败1:")
	//	if len(nodeList) < 2 {
	//		g.Log().Info("合服忽略:%s, %+v, %+v", params.PlatformId, e, nodeList)
	//		continue
	//	}
	//	for _, node := range nodeList {
	//		gameServerList := models.GetGameServerByNode(node.Node)
	//		for _, gameServer := range gameServerList {
	//			serverId, err := strconv.Atoi(models.SubString(gameServer.Sid, 1, len(gameServer.Sid)-1))
	//			c.CheckError(err, "合服失败6:")
	//			if serverId >= e.S && serverId <= e.E {
	//
	//			} else {
	//				g.Log().Error("合服配置错误, [%d-%d] node:%s", e.S, e.E, node.Node)
	//				c.HttpResult(r, enums.CodeFail2, "合服配置错误", "")
	//			}
	//		}
	//	}
	//	g.Log().Info("开始合服:%+v, %s", nodeList, zoneNode)
	//	err = merge.Merge(nodeList[1:], nodeList[0], zoneNode)
	//	c.CheckError(err, "合服失败2:")
	//}
	//if err != nil {
	//	c.HttpResult(r, enums.CodeFail2, "合服失败3:", "")
	//} else {
	//	IsResult = true
	//	c.HttpResult(r, enums.CodeSuccess, "合服成功!", "")
	//}
}

//func getMergeInfo(platformId string, s int, e int) (nodeList [] *models.ServerNode, zoneNode string, err error) {
//	//zoneNode := ""
//	nodes := make([] string, 0)
//	for i := s; i <= e; i++ {
//		gameServer, err := models.GetGameServerOne(platformId, fmt.Sprintf("s%d", i))
//		if err != nil {
//			return nodeList, zoneNode, err
//		}
//		if inArray(gameServer.Node, nodes) {
//			continue
//		}
//		nodes = append(nodes, gameServer.Node)
//		serverNode, err := models.GetServerNode(gameServer.Node)
//		if err != nil {
//			return nodeList, zoneNode, err
//		}
//		if zoneNode == "" {
//			zoneNode = serverNode.ZoneNode
//		}
//		nodeList = append(nodeList, serverNode)
//	}
//	return nodeList, zoneNode, err
//}

func inArray(v string, array []string) bool {
	for _, e := range array {
		if e == v {
			return true
		}
	}
	return false
}

func (c *ToolController) GetWeixinArgs(r *ghttp.Request) {

	c.HttpResult(r, enums.CodeSuccess, "成功!", models.GetWeiXinArgs())
}

func (c *ToolController) UpdateWeixinArgs(r *ghttp.Request) {
	params := &models.WeiXinArgs{}
	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	g.Log().Infof("UpdateWeixinArgs:%+v", params)
	err = models.UpdateWeixinArgs(params)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", 0)
}

// 活动操作
func (c *ToolController) ActivityChange(r *ghttp.Request) {
	var params struct {
		PlatformId string
		ServerId   string
		ChangeType int `json:"activityChangeType"`
		ActivityId int `json:"activityId"`
	}
	var result struct {
		Code    int
		Message string
	}
	err := json.Unmarshal(r.GetBody(), &params)

	g.Log().Infof("后台活动操作:%v", params)
	c.CheckError(err)

	args := fmt.Sprintf("activity_id=%d&activity_type=%d&gm_id=%s",
		params.ActivityId,
		params.ChangeType,
		c.curUser.Account,
	)
	sign := utils.String2md5(args + "fa9274fd68cf8991953b186507840e5e")
	g.Log().Infof("sign:%v", sign)

	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	url := models.GetGameURLByNode(gameServer.Node) + "/gm_activity?" + args + "&sign=" + sign
	g.Log().Infof("url:%v", url)
	resp, err := http.Get(url)
	c.CheckError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	c.CheckError(err)

	err = json.Unmarshal(body, &result)
	g.Log().Infof("result:%v", string(body))
	c.CheckError(err)
	g.Log().Infof("后台活动操作结果:%v", result)
	if result.Code == 0 {
		c.HttpResult(r, enums.CodeSuccess, "后台活动操作成功", 0)
	}
	c.HttpResult(r, enums.CodeFail, fmt.Sprintf("后台活动操作失败: ErrorCode: %v Messsage", result.Code, result.Message), result.Code)
}

// 查看文件
func (c *ToolController) ShowFile(r *ghttp.Request) {
	var params struct {
		FilePath string `json:"filePath"`
		Offset   int    `json:"offset"` // 文件偏移
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err, "查看文件解析数据失败")
	filePathStr := params.FilePath
	//g.Log().Debugf("查看文件/:%+v", strings.Index(filePathStr, "/"))
	//g.Log().Debugf("查看文件\\:%+v", strings.Index(filePathStr, "\\"))
	if strings.Index(filePathStr, "/") != 1 && strings.Index(filePathStr, "\\") != 1 {
		filePathStr = utils.GetShowFileDir() + filePathStr
	}
	g.Log().Debugf("查看文件:%+v", filePathStr)
	Key := "ShowFile" + "_" + filePathStr
	OldTime := utils.GetCacheInt64(Key)
	if OldTime != 0 {
		c.HttpResult(r, enums.CodeFail, "正在解析文件数据", "")
	}
	utils.SetCache(Key, gtime.Timestamp(), 0)
	defer utils.DelCache(Key)
	Offset := 0
	Str := ""
	if utils.IsFileExists(filePathStr) {
		Offset, Str = readFile(filePathStr, params.Offset)
	}
	//g.Log().Debug("查看文件:%+v", Str)

	result := make(map[string]interface{})
	result["offset"] = Offset
	result["rows"] = Str
	c.HttpResult(r, enums.CodeSuccess, "文件解析成功", result)
}

// 读文件
func readFile(path string, Offset int) (int, string) {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	// 使用bufio读取
	r := bufio.NewReader(fi)
	dataStr := ""
	OffsetIndex := 0
	for {
		// 分行读取文件  ReadLine返回单个行，不包括行尾字节(\n  或 \r\n)
		data, _, err := r.ReadLine()

		// 以分隔符形式读取,比如此处设置的分割符是\n,则遇到\n就返回,且包括\n本身 直接返回字符串
		//str, err := r.ReadString('\n')

		// 以分隔符形式读取,比如此处设置的分割符是\n,则遇到\n就返回,且包括\n本身 直接返回字节数数组
		//data, err := r.ReadBytes('\n')

		// 读取到末尾退出
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read err", err.Error())
			break
		}
		OffsetIndex++
		if OffsetIndex <= Offset {
			continue
		}
		dataStr += string(data) + "\r\n"
		// 打印出内容
		//fmt.Printf("-- %v\r\n", string(data))
	}
	return OffsetIndex, dataStr
}

// ------------------------------------------------------------------设置数据---------------------------------------------
// 获得设置数据
func (c *ToolController) GetSettingInfo(r *ghttp.Request) {
	params := &models.BaseQueryParam{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, count := models.GetSettingDataList(params)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "成功!", result)
}

// 添加 编译设置数据
func (c *ToolController) SettingDataEdit(r *ghttp.Request) {
	params := &models.SettingData{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	_, err = models.GetSettingDataOne(params.Id)
	if params.IsAdd == 1 && err == nil {
		c.HttpResult(r, enums.CodeFail, "设置数据已经存在", params.Id)
	}
	if params.IsAdd == 0 && err != nil {
		c.HttpResult(r, enums.CodeFail, "设置数据不存在", params.Id)
	}
	params.UserId = c.curUser.Id
	err = models.UpdateSettingDate(params)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

//删除设置数据
func (c *ToolController) DelSettingData(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	idList := params
	c.CheckError(err)
	err = models.DeleteSettingData(idList)
	c.CheckError(err, "删除设置数据败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除设置数据", idList)
}

// ------------------------------------------------------------------后台通知信息数据模板---------------------------------------------
// 获得后台通知信息数据模板
func (c *ToolController) GetBackgroundMsgTemplateInfo(r *ghttp.Request) {
	params := &models.BaseQueryParam{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, count := models.GetBackgroundMsgTemplateList(params)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "成功!", result)
}

// 添加 编译后台通知信息数据模板
func (c *ToolController) BackgroundMsgTemplateEdit(r *ghttp.Request) {
	params := &models.BackgroundMsgTemplate{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	OldBackgroundMsgTemplateDb, err := models.GetBackgroundMsgTemplate(params.Id)
	if params.IsAdd == 1 && err == nil {
		c.HttpResult(r, enums.CodeFail, "后台通知信息数据模板已经存在", params.Id)
	}
	if params.IsAdd == 0 && err != nil {
		c.HttpResult(r, enums.CodeFail, "后台通知信息数据模板不存在", params.Id)
	}
	params.UserId = c.curUser.Id
	err = models.UpdateBackgroundMsgTemplate(params, OldBackgroundMsgTemplateDb)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

//删除后台通知信息数据模板
func (c *ToolController) DelBackgroundMsgTemplate(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	idList := params
	c.CheckError(err)
	err = models.DeleteBackgroundMsgTemplate(idList)
	c.CheckError(err, "删除通知数据失败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除通知数据", idList)
}

// ------------------------------------------------------------------发邮件设置---------------------------------------------
// 获得发邮件设置
func (c *ToolController) GetMailDataInfo(r *ghttp.Request) {
	params := &models.BaseQueryParam{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, count := models.GetMailDataList(params)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "成功!", result)
}

// 添加 编译发邮件设置
func (c *ToolController) MailDataEdit(r *ghttp.Request) {
	params := &models.MailSmtpData{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	MailData, err := models.GetMailDataOne(params.User)
	if params.IsAdd == 1 {
		if err == nil {
			g.Log().Info(MailData)
			c.HttpResult(r, enums.CodeFail, "邮件设置已经存在", "")
		} else {
			if params.State == 1 {
				mailStart, err := models.GetMailStartData()
				if err == nil {
					c.HttpResult(r, enums.CodeFail, "邮件设置已有开启的，请先关闭"+mailStart.User, "")
				}
			}
			if len(params.Pass) == 0 {
				c.HttpResult(r, enums.CodeFail, "邮件设置已密码不能为空", "")
			}
		}
	}
	if params.IsAdd == 0 {
		if err != nil {
			c.HttpResult(r, enums.CodeFail, "邮件设置不存在", params.User)
		} else if params.State == 1 && params.State != MailData.State {
			mailStart, err := models.GetMailStartData()
			if err == nil {
				g.Log().Info(MailData)
				c.HttpResult(r, enums.CodeFail, "邮件设置已有开启的，请先关闭:"+mailStart.User, "")
			}
		}
		if params.Pass == "" {
			params.Pass = MailData.Pass
		}
	}
	err = models.UpdateMailData(params)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

//删除发邮件设置
func (c *ToolController) DelMailData(r *ghttp.Request) {
	var params []string
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	err = models.DeleteMailData(params)
	c.CheckError(err, "删除设置数据败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除设置数据", params)
}

// ------------------------------------------------------------------界面操作权限---------------------------------------------
// 获得界面操作权限
func (c *ToolController) GetPageChangeAuth(r *ghttp.Request) {
	params := &models.BaseQueryParam{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, count := models.GetPageChangeAuthList(params)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "成功!", result)
}

// 添加 编译界面操作权限
func (c *ToolController) PageChangeAuthEdit(r *ghttp.Request) {
	params := &models.PageChangeAuth{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	OldAuthData, err := models.GetPageChangeAuthOneBySign(params.Sign)
	if params.IsAdd == 1 && err == nil {
		g.Log().Info(OldAuthData)
		c.HttpResult(r, enums.CodeFail, "界面操作权限已经存在", "")
	}
	if params.IsAdd == 0 && err != nil {
		c.HttpResult(r, enums.CodeFail, "界面操作权限不存在", params.Sign)
	}
	params.UserId = c.curUser.Id
	err = models.UpdatePageChangeAuth(params, OldAuthData)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

//删除界面操作权限
func (c *ToolController) DelPageChangeAuth(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	err = models.DeletePageChangeAuth(params)
	c.CheckError(err, "删除界面操作权限失败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除界面操作权限", params)
}

// 获得界面操作权限状态
func (c *ToolController) GetChangeAuthState(r *ghttp.Request) {
	var params []string
	type resultStruct struct {
		Key    string `json:"key"`
		IsAuth bool   `json:"isAuth"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	list := make([]*resultStruct, 0)
	for _, key := range params {
		resultData := &resultStruct{
			Key:    key,
			IsAuth: models.IsPageChangeAuthDefault(c.curUser, key, false),
		}
		list = append(list, resultData)
	}
	result := make(map[string]interface{})
	result["rows"] = list
	c.HttpResult(r, enums.CodeSuccess, "成功获得界面操作权限状态", result)
}

// 获得订时器信息
func (c *ToolController) GetCronInfo(r *ghttp.Request) {
	var params struct {
		CronNameStr string `json:"cron_name_str"`
	}

	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	data, count := models.GetCronList(params.CronNameStr)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获得订时器信息", result)
}

// 获得机器系统信息
func (c *ToolController) GetSysRobotInfo(r *ghttp.Request) {
	data := models.GetSysRobotInfo()
	//result := make(map[string]interface{})
	//result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获得机器系统信息", data)
}

// ------------------------------------------------------------------新客户端版本管理---------------------------------------------

// 获得客户端版本
func (c *ToolController) GetVersionList(r *ghttp.Request) {
	data, count := models.GetVersionList()
	g.Log().Infof("获得客户端版本列表:%+v , count: %+v", data, count)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获得客户端版本列表成功", result)
}

// 获得客户端版本
func (c *ToolController) GetVersion(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:"platformId"`
	}
	g.Log().Infof("获得客户端版本:%+v", params)
	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)

	data, err := models.GetVersionOne(params.PlatformId)
	g.Log().Infof("获得客户端版本ResultVersion:%+v , err: %+v", data, err)
	utils.CheckError(err)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获得客户端版本成功", result)
}

// 设置客户端版本
func (c *ToolController) SetVersion(r *ghttp.Request) {
	params := &models.ClientVersions{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	if params.PlatformId == "" {
		c.HttpResult(r, enums.CodeFail, "PlatformId不能为空 ", params.PlatformId)
	}
	if params.FirstVersions <= 0 {
		c.HttpResult(r, enums.CodeFail, "FirstVersions错误", params.FirstVersions)
	}
	if params.Versions <= 0 {
		c.HttpResult(r, enums.CodeFail, "Versions错误", params.Versions)
	}
	_, err = models.GetPlatformOne(params.PlatformId)
	if err != nil {
		c.HttpResult(r, enums.CodeFail, "PlatformId未配置 ", params.PlatformId)
	}
	ip := r.GetClientIp()
	userId := c.curUser.Id
	err = models.SetVersion(params, ip, userId)
	if err != nil {
		c.HttpResult(r, enums.CodeFail, "设置平台版本报错 ", err)
	}
	c.HttpResult(r, enums.CodeSuccess, "设置客户端版本成功", "")
}

// --------------------------------------------渠道客户端信息---------------------------------------------
// 获得渠道客户端信息列表
func (c *ToolController) GetPlatformClientInfoList(r *ghttp.Request) {
	params := &models.PlatformClientInfoRequest{}
	g.Log().Infof("获得渠道客户端信息列表请求参数:%+v", params)

	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)

	data, count := models.GetPlatformClientInfoList(params)
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data

	g.Log().Infof("result: %+v", result)
	c.HttpResult(r, enums.CodeSuccess, "获得渠道客户端信息列表成功", result)
}

func (c *ToolController) DelStatisticRes(r *ghttp.Request) {
	params := &models.DelStatisticResReq{}
	g.Log().Infof("静态资源删除请求参数:%+v", params)

	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)

	data := &models.StatisticResource{}
	err = models.Db.Where(&models.StatisticResource{Id: params.Id}).First(&data).Error
	utils.CheckError(err)

	err = models.Db.Delete(&models.StatisticResource{Id: params.Id}).Error
	utils.CheckError(err)

	type SendData struct {
		Url     []map[string]string `json:"url"`
		AppId   string              `json:"app_id"`
		Version int                 `json:"version"`
	}

	var reqData = SendData{
		AppId:   data.AppId,
		Url:     models.StatisticRes.FetchRelUrlVer(data.AppId),
		Version: data.Version,
	}

	jsonData, _ := json.Marshal(reqData)
	// fmt.Println(string(jsonData))
	//向中心服发送数据
	var url = g.Cfg().GetString("game.gs_domain") + "/update_static_resource"
	retStr, err := utils.HttpRequest(url, string(jsonData))
	utils.CheckError(err)

	c.HttpResult(r, enums.CodeSuccess, "静态资源删除成功", retStr)
}

// 静态资源更新
func (c *ToolController) EditStatisticRes(r *ghttp.Request) {
	req := &models.UptatisticResReq{}
	g.Log().Infof("静态资源更新请求参数:%+v", req)

	err := json.Unmarshal(r.GetBody(), &req)
	utils.CheckError(err)

	var urlFix = req.Url[len(req.Url)-4:]

	if len(req.Url) > 5 && urlFix != ".zip" {
		c.HttpResult(r, enums.CodeFail, "link必须以.zip", nil)
	} else if len(req.Url) < 5 {
		c.HttpResult(r, enums.CodeFail, "过短的url", nil)
	}

	data := models.StatisticResource{
		Id:      req.Id,
		Url:     req.Url,
		AppId:   req.AppId,
		Version: req.Versoin,
		Uptime:  int(gtime.Timestamp()),
	}

	err = models.Db.Omit("add_time_str").Save(&data).Error
	c.CheckError(err, "写数据失败")

	type SendData struct {
		Url     []map[string]string `json:"url"`
		AppId   string              `json:"app_id"`
		Version int                 `json:"version"`
	}

	var reqData = SendData{
		AppId:   req.AppId,
		Url:     models.StatisticRes.FetchRelUrlVer(req.AppId),
		Version: req.Versoin,
	}

	jsonData, _ := json.Marshal(reqData)
	// fmt.Println(string(jsonData))
	//向中心服发送数据
	var url = g.Cfg().GetString("game.gs_domain") + "/update_static_resource"
	retStr, err := utils.HttpRequest(url, string(jsonData))
	utils.CheckError(err)

	c.HttpResult(r, enums.CodeSuccess, "静态资源更新成功", retStr)
}

// 获得静态资源列表
func (c *ToolController) OptStatisticRes(r *ghttp.Request) {
	params := &models.OptStatisticResReq{}
	g.Log().Infof("获得渠道客户端信息列表请求参数:%+v", params)

	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)

	count := 0
	data := make([]*models.StatisticResource, 0)
	m := models.Db
	if params.AppId != "" {
		m = m.Where(models.StatisticResource{AppId: params.AppId})
	}
	err = m.Offset(params.Offset).Limit(params.Limit).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)

	for k, v := range data {
		data[k].AddTimeStr = utils.TimeInt64FormDefault(int64(v.AddTime))
	}

	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获得静态资源列表成功", result)
}

// 添加静态资源管理
func (c *ToolController) AddStatisticRes(r *ghttp.Request) {
	req := &models.AddStatisticResReq{}
	g.Log().Infof("获得渠道客户端信息列表请求参数:%+v", req)

	err := json.Unmarshal(r.GetBody(), &req)
	utils.CheckError(err)

	var urlFix = req.Url[len(req.Url)-4:]

	if len(req.Url) > 5 && urlFix != ".zip" {
		c.HttpResult(r, enums.CodeFail, "link必须以.zip", nil)
	} else if len(req.Url) < 5 {
		c.HttpResult(r, enums.CodeFail, "过短的url", nil)
	}
	data := models.StatisticResource{
		Url:     req.Url,
		AppId:   req.AppId,
		Version: req.Versoin,
		AddTime: int(gtime.Timestamp()),
	}
	err = models.Db.Omit("add_time_str").Save(&data).Error
	c.CheckError(err, "写数据失败")

	type SendData struct {
		Url     []map[string]string `json:"url"`
		AppId   string              `json:"app_id"`
		Version int                 `json:"version"`
	}

	var reqData = SendData{
		AppId:   req.AppId,
		Url:     models.StatisticRes.FetchRelUrlVer(req.AppId),
		Version: req.Versoin,
	}

	jsonData, _ := json.Marshal(reqData)
	// fmt.Println(string(jsonData))
	//向中心服发送数据
	var url = g.Cfg().GetString("game.gs_domain") + "/update_static_resource"
	retStr, err := utils.HttpRequest(url, string(jsonData))
	utils.CheckError(err)

	c.HttpResult(r, enums.CodeSuccess, "静态资源添加成功", retStr)
}

// 获得全部渠道客户端信息列表
func (c *ToolController) GetPlatformClientInfoAllList(r *ghttp.Request) {
	data, count := models.GetPlatformClientInfoAllList()
	g.Log().Infof("获得全部渠道客户端信息列表:%+v , count: %+v", data, count)
	result := make(map[string]interface{})
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获得全部渠道客户端信息列表成功", result)
}

// 添加 编译 渠道客户端信息
func (c *ToolController) PlatformClientInfoEdit(r *ghttp.Request) {
	params := &models.PlatformClientInfo{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	// id为0的话，就代表是添加
	if params.Id == 0 {
		OldPlatformClientInfo, err := models.GetPlatformClientInfoOneByAppId(params.AppId)
		g.Log().Info("create: %+v", OldPlatformClientInfo)
		if err == nil {
			c.HttpResult(r, enums.CodeFail, "渠道客户端app_id已经存在", params.Id)
		} else {
			err = models.AddPlatformClientInfo(params)
			if err != nil {
				g.Log().Errorf("创建失败: %+v", err)
				c.HttpResult(r, enums.CodeFail, "创建失败", params.Id)
			}
		}
	} else {
		OldPlatformClientInfo, err := models.GetPlatformClientInfoOneByAppId(params.AppId)
		g.Log().Info("update: %+v", OldPlatformClientInfo)
		if err == nil {
			err = models.UpdatePlatformClientInfo(params)
			if err != nil {
				g.Log().Errorf("编辑失败: %+v", err)
				c.HttpResult(r, enums.CodeFail, "编辑失败", params.Id)
			}
		} else {
			c.HttpResult(r, enums.CodeFail, "渠道客户端信息不存在", params.Id)
		}
	}
	c.CheckError(err)
	url := utils.GetCenterURL() + "/update_version"
	//url := "http://192.168.31.153:6663" + "/update_version"
	var request struct {
		Platform           string  `json:"platform"`
		PlatformName       string  `json:"platformName"`
		AndroidDownloadUrl string  `json:"androidDownloadUrl"`
		IosDownloadUrl     string  `json:"iOSDownloadUrl"`
		Versions           string  `json:"versions"`
		FirstVersions      string  `json:"firstVersions"`
		IsCloseCharge      int     `json:"isCloseCharge"`
		ReviewingVersions  string  `json:"reviewingVersions"`
		ClientVersion      string  `json:"clientVersion"`
		AppId              string  `json:"appId"`
		FacebookAppId      string  `json:"facebookAppId"`
		ReloadUrl          string  `json:"reloadUrl"`
		NativePay          int     `json:"nativePay"`
		Stats              int     `json:"stats"`
		PayTimes           int     `json:"payTimes"`
		Region             string  `json:"region"`
		Channel            string  `json:"channel"`
		PackageSize        float64 `json:"packageSize"`
		AreaCode           string  `json:"areaCode"`
		Domain             string  `json:"domain"`
		TestDomain         string  `json:"testDomain"`
	}

	request.Platform = params.Platform
	request.PlatformName = params.PlatformRemark
	request.AndroidDownloadUrl = params.UpgradeAndroidUrl
	request.IosDownloadUrl = params.UpgradeIosUrl
	request.FirstVersions = params.FirstVersions
	request.Versions = params.Versions
	request.IsCloseCharge = 1 - params.IsChargeOpen
	request.ClientVersion = params.ClientVersion
	request.AppId = params.AppId
	request.FacebookAppId = params.FacebookAppId
	request.ReloadUrl = params.ReloadUrl
	request.NativePay = params.NativePay
	request.Stats = params.Stats
	request.PayTimes = params.PayTimes
	request.Region = params.Region
	request.Channel = params.Channel
	request.ReviewingVersions = params.ReviewingVersions
	// PackageSize
	request.PackageSize = params.PackageSize
	request.AreaCode = params.AreaCode
	// domain
	request.Domain = params.Domain
	request.TestDomain = params.TestDomain
	data2, err := json.Marshal(request)
	utils.CheckError(err)

	_, err = utils.HttpRequest(url, string(data2))
	fmt.Println(err)
	c.HttpResult(r, enums.CodeSuccess, "成功!", "")
}

// 删除渠道客户端信息
func (c *ToolController) DelPlatformClientInfo(r *ghttp.Request) {
	var params []int
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	err = models.DeletePlatformClientInfo(params)
	c.CheckError(err, "删除渠道客户端信息失败")
	c.HttpResult(r, enums.CodeSuccess, "成功删除渠道客户端信息", params)
}

// 手动生成每十分钟数据
func (c *ToolController) ManaulTenMinuteData(r *ghttp.Request) {
	next1 := gtime.Now()
	nextTimestamp := next1.Unix()
	//// 业务
	g.Log().Infof("整点10分钟定时执行:%v", next1.String())
	models.DoUpdateAllGameNodeTenMinuteStatistics(int(nextTimestamp))
}

// 手动生成时间点 怪物日志
func (c *ToolController) ManaulAddGameMonsterlog(r *ghttp.Request) {
	params := &models.ToolTest{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)

	models.InsertGameMonsterLog(time.Unix(int64(params.OptTime), 0), params.Hour, params.GenTodayAll)
}

// 手动生成时间点 游戏使用事件和物品使用日志
func (c *ToolController) ManaulAddGameUselog(r *ghttp.Request) {
	params := &models.ToolTest{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)

	models.InsertItemEventLog(time.Unix(int64(params.OptTime), 0), params.Hour, params.GenTodayAll)
}

func (c *ToolController) PlusDailyStatistics(r *ghttp.Request) {
	params := &models.PlusDailyStatistics{}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	models.PlusUpdateAllGameNodeDailyStatistics(params.StartTime, params.EndTime)
}

// 添加静态资源添加静管理列表
func (c *ToolController) Test(r *ghttp.Request) {

	// params := &models.ModifyPushReq{
	// 	Platform: "local",
	// }

	// utils.HttpPost("", json.)
	// var data struct {
	// 	Count float32
	// }
	// sql := `select sum(money) as count from charge_info_record);`
	// err := models.DbCharge.Raw(sql).Scan(&data).Error
	// utils.CheckError(err)

	// fmt.Println("count su : ", data.Count)
	// params := &models.ToolTest{}
	// err := json.Unmarshal(r.GetBody(), &params)
	// c.CheckError(err)

	// models.InsertItemEventLog(time.Unix(int64(params.OptTime), 0), params.Hour, 0)

	// data map[string]map[string]*MemItemEventLog
	// data := uselog.FetchALLDataByPS("local", "s1", "13", "1001")
	// nowHour := time.Now().Hour()
	// dataThisHour := models.GetItemEventLogList(&models.ItemEventLogQueryParam{
	// 	PlatformId: "local",
	// 	ServerId:   "s1",
	// 	LogType:    13,
	// 	Type:       1,
	// 	Datetime:   int(gtime.Timestamp()),
	// }, nowHour)

	// if data == nil && dataThisHour != nil {
	// 	data = dataThisHour
	// } else {
	// 	for timeStr, _ := range data {
	// 		data[timeStr][gconv.String(nowHour)] = dataThisHour[timeStr][gconv.String(nowHour)]
	// 	}
	// }
	// // fmt.Println(fmt.Sprintf("%v", dataThisHour))
	// // g.Dump(dataThisHour)

	// var ret = make([]*uselog.MemItemEventLog, 0)
	// for _, item := range data {
	// 	for timeStr, value := range item {
	// 		value.Number = gconv.Int(timeStr)
	// 		ret = append(ret, value)
	// 	}
	// }

	// c.HttpResult(r, enums.CodeSuccess, "ok", ret)
	// cmd := ""
	// fmt.Println(cmd)
	// test UpdateAllGameNodeDailyStatistics cron
	// models.UpdateAllGameNodeDailyStatistics()
	// c.HttpResult(r, enums.CodeSuccess, "ok", nil)

	// test push init
	// memdb.ScanDataInTable()

	// test insert or update push list
	// s1_7  Chenxizhu
	// str := memdb.InsertOrUpPushList("Chenxizhu", "s1", "registrate id ....")s
	// c.HttpResult(r, enums.CodeSuccess, "ok", str)
	// test push list
	//所有的可推送列表
	// var pushList = make(map[string][]string)

	// data := memdb.FetchDataBySids("local", []string{})
	// for areaSid, subData := range data {
	// 	// fmt.Println(areaSid, item["open"])
	// 	pushList[areaSid] = make([]string, 0)
	// 	for _, item := range subData {
	// 		if item["open"] == "1" && item["regis"] != "" {
	// 			// fmt.Println(item["regis"])
	// 			pushList[areaSid] = append(pushList[areaSid], item["regis"])
	// 		}
	// 	}
	// }
	// fmt.Println(pushList)

	// queue := gqueue.New()
	// queue.Push(data)
	// queue.Pop()
	// // g.List
	// c.HttpResult(r, enums.CodeSuccess, "ok", data)
}

func (c *ToolController) _checkPushReq(r *ghttp.Request) *models.ModifyPushReq {
	req := &models.DoPushReq{}
	err := json.Unmarshal(r.GetBody(), &req)
	utils.CheckError(err)

	//验签
	reqMap := utils.DecodeHttpRequest(req.Data, req.Sign)
	if reqMap == nil {
		c.HttpResult(r, enums.CodeFail, "fail", nil)
	}

	params := &models.ModifyPushReq{}
	if errParams := gconv.Struct(reqMap, params); errParams != nil {
		utils.CheckError(errParams)
		c.HttpResult(r, enums.CodeFail, "fail", errParams.Error())

	}
	return params
}

func (c *ToolController) GenparamsPushList(r *ghttp.Request) {
	req := &models.ModifyPushReq{}
	err := json.Unmarshal(r.GetBody(), &req)
	utils.CheckError(err)

	dataByte, err := json.Marshal(req)
	utils.CheckError(err)

	data, sign := utils.EncodeHttpRequest(string(dataByte))

	c.HttpResult(r, enums.CodeSuccess, "ok", g.Map{
		"data": data,
		"sign": sign,
	})
}

// 添加静态资源添加静管理列表
func (c *ToolController) GotPushList(r *ghttp.Request) {
	params := c._checkPushReq(r)
	if params.Platform == "" {
		c.HttpResult(r, enums.CodeFail, "fail", "platform is empty")
		return
	}

	c.HttpResult(r, enums.CodeSuccess, "ok", memdb.PushList(params.Platform))
}

//  params.Sid 为空时更新，否则插入 , 更新时自动将open状态变为1
func (c *ToolController) InsertOrUpPushList(r *ghttp.Request) {
	params := c._checkPushReq(r)

	if params.Platform == "" {
		c.HttpResult(r, enums.CodeFail, "fail", "platform is empty")
		return
	}

	var sidArr []string
	if params.Sid != "" {
		sidArr = strings.Split(params.Sid, ",")
	}
	str := memdb.InsertOrUpPushList(params.Platform, params.Account, params.RegistrationId, sidArr) // 76 , 81
	if params.Sid != "" {
		// 插入后更新所有这个帐户的区服角色的RegistrationId
		memdb.InsertOrUpPushList(params.Platform, params.Account, params.RegistrationId, []string{}) // 76 , 81
	}

	c.HttpResult(r, enums.CodeSuccess, "ok", str)
}

// {
// 	"id": 12,               数据序号
// 	"s_account": "xixi",    客服帐号
// 	"uid": 14646,           玩家ID
// 	"plat_id": "local",     平台
// 	"server_id": "s1",      区服
// 	"create_time": 1632558484,  创建时间
// 	"uptime": 0,   修改时间
// 	"nick": "卡卡罗卜"       玩家昵称
// },
//  极光推送-设置不推送帐户
func (c *ToolController) SetNopushAccount(r *ghttp.Request) {
	params := c._checkPushReq(r)

	if params.Platform == "" {
		c.HttpResult(r, enums.CodeFail, "fail", "platform is empty")
		return
	}

	retStr := memdb.SetNopushAccount(params.Account, params.Platform, nil)
	if retStr == "" {
		c.HttpResult(r, enums.CodeFail, "ok", retStr)
	} else {
		c.HttpResult(r, enums.CodeSuccess, "ok", retStr)
	}
}

func (c *ToolController) ShowBindCustomerList(r *ghttp.Request) {
	var req struct {
		models.BaseQueryParam
		PlatId   string `json:"platform_id"`
		ServerId string `json:"server_id"`
		SAccount string `json:"server_name"`
		Uid      int    `json:"player_id"`
	}
	err := json.Unmarshal(r.GetBody(), &req)
	utils.CheckError(err)

	var player = models.Player{
		Id: req.Uid,
	}
	if req.Uid > 0 && req.PlatId != "" && req.ServerId != "" {
		gameDb, err := models.GetGameDbByPlatformIdAndSid(req.PlatId, req.ServerId)
		utils.CheckError(err)
		gameDb.Where(player).First(&player)
		if player.Nickname == "" {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", []int{})
			return
		}
	}

	var count int
	var data = make([]*models.ServerRel, 0)
	whereStr := ""
	if len(req.SAccount) > 0 {
		whereStr = fmt.Sprintf(" s_account IN (%s)", req.SAccount)
	}

	err = models.Db.Debug().Offset(req.Offset).Limit(req.Limit).Where(models.ServerRel{
		Uid:      player.Id,
		ServerId: req.ServerId,
	}).Where(whereStr).
		Find(&data).Offset(0).Count(&count).Error

	utils.CheckError(err)

	for key, item := range data {
		if req.Uid > 0 {
			data[key].Id = req.Uid
		} else {
			var thePlayer = models.Player{
				Id: item.Uid,
			}
			gameDb, err := models.GetGameDbByPlatformIdAndSid(item.PlatId, item.ServerId)
			utils.CheckError(err)
			err = gameDb.Where(thePlayer).First(&thePlayer).Error
			utils.CheckError(err)

			data[key].Uname = thePlayer.Nickname

		}
		if item.SAccount > 0 {
			var SAccoutUser *models.User
			SAccoutUser, err = models.GetUserOne(item.SAccount)
			utils.CheckError(err)
			data[key].CustomerName = SAccoutUser.Name
		} else {
			data[key].CustomerName = ""
		}
	}
	result := make(map[string]interface{})
	result["total"] = count
	result["rows"] = data

	c.HttpResult(r, enums.CodeSuccess, "获取绑定列表成功", result)
}

func (c *ToolController) DelCustomerBind(r *ghttp.Request) {
	var req struct {
		Id int `json:"id"`
	}
	err := json.Unmarshal(r.GetBody(), &req)
	utils.CheckError(err)

	err = models.Db.Where("id=" + gconv.String(req.Id)).Delete(&models.ServerRel{}).Error
	utils.CheckError(err)

	c.HttpResult(r, enums.CodeSuccess, "ok", 1)
}

// 绑定客户
func (c *ToolController) BindUpCustomer(r *ghttp.Request) {
	var req struct {
		Id int `json:"id"` // id>0时修改，id=0时添加
		// PlatId   string `json:"platform_id"` // 平台id
		// ServerId string `json:"server_id"`   // 区服id
		SAccount int `json:"server_name"` // 客服帐号
		// PName    string `json:"player_name"` // 玩家昵称
		Uid    int    `json:"player_id"` // 玩家编号
		Remark string `json:"remark"`    // 备注
	}
	err := json.Unmarshal(r.GetBody(), &req)
	utils.CheckError(err)
	if req.Uid == 0 {
		c.HttpResult(r, enums.CodeFail, "玩家不存在", []int{})
		return
	}

	//
	var gplayer = models.GlobalPlayer{
		Id: req.Uid,
	}

	err = models.DbCenter.Where(gplayer).First(&gplayer).Error
	utils.CheckError(err)
	if gplayer.Nickanme == "" {
		c.HttpResult(r, enums.CodeFail, "玩家不存在", []int{})
		return
	}
	//检查客服账号
	var SAccoutUser *models.User
	SAccoutUser, err = models.GetUserOne(req.SAccount)
	utils.CheckError(err)

	if SAccoutUser.Name == "" {
		c.HttpResult(r, enums.CodeFail, "请输入正确的客服id(用户管理的id列)", []int{})
		return
	}
	var nowTime = gtime.Now().Timestamp()
	var data = models.ServerRel{
		Uid:      gplayer.Id,
		SAccount: req.SAccount,
		Remark:   req.Remark,
	}

	if req.Id > 0 {
		err = models.Db.Where(data).First(&data).Error
		utils.CheckError(err)
		data.Id = req.Id
		data.Uptime = nowTime
		err := models.Db.Debug().Table(models.ServerRelTBName()).Where("id=" + gconv.String(req.Id)).Update(&data).Error
		utils.CheckError(err)
	} else {
		data.PlatId = gplayer.PlatformId
		data.ServerId = gplayer.ServerId
		data.CreateTime = nowTime
		var inData = models.InServerRel{
			PlatId:     data.PlatId,
			ServerId:   data.ServerId,
			SAccount:   data.SAccount,
			Uid:        data.Uid,
			CreateTime: data.CreateTime,
			Remark:     data.Remark,
		}
		err = models.Db.Save(&inData).Error
		utils.CheckError(err)
		if err != nil {
			c.HttpResult(r, enums.CodeFail, "重复绑定 ", 0)
			return
		}
	}

	c.HttpResult(r, enums.CodeSuccess, "ok", 1)
}

func (c *ToolController) ToNicks(r *ghttp.Request) {
	var req struct {
		Uids string `json:"uids"`
	}
	err := json.Unmarshal(r.GetBody(), &req)
	utils.CheckError(err)

	uidArr := strings.Split(req.Uids, ",")
	utils.CheckError(err)

	type NeedData struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var players = make([]*NeedData, 0)
	sql3 := `select id,concat(server_id, ".", nickanme) name from (select "1111114444" pid`
	for _, pidStr := range uidArr {
		sql3 += ` union all select "` + pidStr + `" `
	}
	sql3 += `) as tmp, global_player t where tmp.pid=t.id;` //

	err = models.DbCenter.Raw(sql3).Scan(&players).Error
	utils.CheckError(err)

	c.HttpResult(r, enums.CodeSuccess, "ok", players)
}
