package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"gfWeb/memdb"
	"strconv"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	//"myadmin_h5/proto"
	//"github.com/golang/protobuf/proto"
	//"encoding/base64"
	//"net/http"
	//"strings"
	//"io/ioutil"
)

type PlayerController struct {
	BaseController
}

//测试帐号显示
func (c *PlayerController) TestAccountList(r *ghttp.Request) {
	list := memdb.List()
	// r.Response.Write(list)
	c.HttpResult(r, enums.CodeSuccess, "获取测试玩家列表成功", list)
}

//测试帐号写入内存
func (c *PlayerController) AddTestAccount(r *ghttp.Request) {
	var req = struct {
		Name      string `json:"name"`
		Privilege int    `json:"privilege"`
	}{}
	err := json.Unmarshal(r.GetBody(), &req)
	utils.CheckError(err)

	errStr := memdb.InsertAccount(memdb.TestAccountObj{
		Name:      req.Name,
		Privilege: req.Privilege,
	})

	if errStr != "" {
		c.HttpResult(r, enums.CodeFail, errStr, "")
		return
	}

	//请求
	var reqBody = struct {
		Account   string `json:"account"`
		Privilege int    `json:"privilege"`
	}{
		Account:   req.Name,
		Privilege: req.Privilege,
	}
	// g.Dump(reqBody)
	data, err := json.Marshal(reqBody)
	utils.CheckError(err)

	var url = g.Cfg().GetString("game.gs_domain") + "/test_account/add"
	retStr, err := utils.HttpRequest(url, string(data))

	if err != nil {
		memdb.Delete(req.Name)
		utils.CheckError(err)
	}
	c.HttpResult(r, enums.CodeSuccess, "添加测试帐号成功", retStr)
}

func (c *PlayerController) List(r *ghttp.Request) {
	var params models.PlayerQueryParam
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Info("查询用户列表:%+v", params)
	data, total := models.GetPlayerList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data

	withdrawalList := models.GetPlayerWithdrawalInfoList()
	rechargeList := models.GetPlayerRechargeInfoList()

	for _, v := range data {
		if val, ok := withdrawalList[v.Id]; ok {
			v.TotalWithdrawalMoney = val.Sum
			v.TotalWithDrawalTimes = val.Count
		}
		if val, ok := rechargeList[v.Id]; ok {
			v.TotalChargeTime = val.Count
		}
	}

	c.HttpResult(r, enums.CodeSuccess, "获取玩家列表成功", result)
}

func (c *PlayerController) Detail(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
		PlayerId   int    `json:"playerId"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	utils.CheckError(err)
	g.Log().Infof("查询玩家详细信息:%+v", params)
	playerDetail, err := models.GetPlayerDetail(params.PlatformId, params.ServerId, params.PlayerId)
	c.CheckError(err, "查询玩家详细信息失败")
	c.HttpResult(r, enums.CodeSuccess, "获取玩家详细信息成功", playerDetail)
}

func (c *PlayerController) AccountDetail(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:"platformId"`
		Account    string `json:"account"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	g.Log().Info("查询帐号详细信息:%+v", params)
	if params.PlatformId == "" || params.Account == "" {
		c.CheckError(gerror.New("平台和帐号不能为空"), "查询帐号详细信息失败")
	}
	playerDetail, err := models.GetGlobalPlayerList(params.PlatformId, params.Account)
	c.CheckError(err, "查询帐号详细信息失败")
	c.HttpResult(r, enums.CodeSuccess, "查询帐号详细信息成功", playerDetail)
}

func (c *PlayerController) One(r *ghttp.Request) {

	platformId := r.GetQueryString("platformId")
	//serverId := c.GetString("serverId")
	playerName := r.GetQueryString("playerName")
	player, err := models.GetPlayerByPlatformIdAndNickname(platformId, playerName)
	c.CheckError(err, "查询玩家失败")

	if player.Id > 0 {
		var sbind = models.ServerRel{
			Uid: player.Id,
		}
		err := models.Db.Where(sbind).First(&sbind).Error
		utils.CheckError(err)

		player.SAccount = string(sbind.SAccount)
	}

	c.HttpResult(r, enums.CodeSuccess, "获取玩家成功", player)
}

func (c *PlayerController) GetPlayerByEncodedId(r *ghttp.Request) {
	platformId := r.GetQueryString("platformId")
	serverId := r.GetQueryString("serverId")
	playerIdEncoded := r.GetQueryString("playerIdEncoded")
	g.Log().Debug("playerIdEncoded: ", playerIdEncoded)

	url := utils.GetChargeURL() + "/customer/decodePlayerId"
	g.Log().Debug("url: ", url)
	//url := "http://192.168.31.153:6666" + "/update_version"
	var request struct {
		DecodePlayerId string `json:"decodePlayerId"`
	}

	request.DecodePlayerId = playerIdEncoded
	data2, err := json.Marshal(request)
	g.Log().Debug("data2: ", data2)
	utils.CheckError(err)

	PlayerId := ""
	PlayerId, err = utils.HttpRequest(url, string(data2))
	g.Log().Debug("PlayerId: ", PlayerId)

	intPlayerId, err := strconv.Atoi(PlayerId)
	//player, err := models.GetPlayerByPlatformIdAndNickname(platformId, playerIdEncoded)
	player, err := models.GetPlayerByPlayerId(platformId, serverId, intPlayerId)
	c.CheckError(err, "查询玩家失败")
	c.HttpResult(r, enums.CodeSuccess, "获取玩家成功", player)
}

// 设置帐号类型
func (c *PlayerController) SetAccountType(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:"platformId"`
		PlayerId   int    `json:"playerId"`
		ServerId   string `json:"serverId"`
		Type       int32  `json:"type"`
		PayTimes   int    `json:"payTimes"` //支付达到一定次数切换第三方支付
	}
	//var result struct {
	//	ErrorCode int
	//}
	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	g.Log().Info("设置帐号类型:%+v", params)

	request, err := json.Marshal(params)
	c.CheckError(err)

	_, err = models.GetPlayerOne(params.PlatformId, params.ServerId, params.PlayerId)
	c.CheckError(err)

	err = models.AddPlayerInfos(r)
	c.CheckError(err)

	url := models.GetGameURLByPlatformIdAndSid(params.PlatformId, params.ServerId) + "/set_account_type"
	data := string(request)
	resultMsg, err := utils.HttpRequest(url, data)
	c.CheckError(err, resultMsg)

	//url := models.GetGameURLByPlatformIdAndSid(params.PlatformId, params.ServerId) + "/set_account_type"
	//data := string(request)
	//g.Log().Info("url:%s", url)
	//sign := utils.String2md5(data + enums.GmSalt)
	//base64Data := base64.URLEncoding.EncodeToString([]byte(data))
	//requestBody := "data=" + base64Data+ "&sign=" + sign
	//resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(requestBody))
	//c.CheckError(err)
	//
	//defer resp.Body.Close()
	//responseBody, err := ioutil.ReadAll(resp.Body)
	//c.CheckError(err)
	//
	//g.Log().Info("result:%v", string(responseBody))
	//
	//err = json.Unmarshal(responseBody, &result)
	//
	//c.CheckError(err)
	//if result.ErrorCode != 0 {
	//	c.HttpResult(r, enums.CodeFail, "设置帐号类型失败", 0)
	//}
	c.HttpResult(r, enums.CodeSuccess, "设置帐号类型成功", 0)
	//serverId := player.ServerId
	//request := gm.MSetAccountTypeTos{Token: proto.String(""), Type: proto.Int32(params.Type), PlayerId: proto.Int32(int32(params.PlayerId))}
	//mRequest, err := proto.Marshal(&request)
	//c.CheckError(err)
	//
	//conn, err := models.GetWsByPlatformIdAndSid(params.PlatformId, params.ServerId)
	//c.CheckError(err)
	//defer conn.Close()
	//_, err = conn.Write(utils.Packet(9907, mRequest))
	//c.CheckError(err)
	//var receive = make([]byte, 100, 100)
	//n, err := conn.Read(receive)
	//c.CheckError(err)
	//response := &gm.MSetAccountTypeToc{}
	//data := receive[5:n]
	//err = proto.Unmarshal(data, response)
	//c.CheckError(err)
	//
	//if *response.Result == gm.MSetAccountTypeToc_success {
	//	c.HttpResult(r, enums.CodeSuccess, "设置帐号类型成功", 0)
	//} else {
	//	c.HttpResult(r, enums.CodeFail, "设置帐号类型失败", 0)
	//}
}
