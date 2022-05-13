// 后台充值
package controllers

import (
	"encoding/json"
	"fmt"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"io/ioutil"
	"net/http"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type BackgroundController struct {
	BaseController
}

func (c *BackgroundController) List(r *ghttp.Request) {
	var params models.BackgroundChargeLogQueryParam
	err := json.Unmarshal(r.GetBody(), &params)

	//err := gconv.Struct(c.Ctx.Input.RequestBody, &params)
	c.CheckError(err)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.HttpResult(r, enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetBackgroundChargeLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取后台充值日志列表成功", result)
}

func (c *BackgroundController) Charge(r *ghttp.Request) {
	var params struct {
		Account     string
		Ip          string
		PlayerId    int
		PlatformId  string
		ChargeValue int
		ServerId    string
		ChargeType  string
		ItemId      int
	}
	var result struct {
		Code    int
		Message string
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)

	player, err := models.GetPlayerOne(params.PlatformId, params.ServerId, params.PlayerId)
	c.CheckError(err)

	accountType := models.GetAccountType(params.PlatformId, player.AccId)

	if accountType == 0 {
		c.HttpResult(r, enums.CodeFail, fmt.Sprintf("后台充值失败: %s.%s 不是内部帐号", params.ServerId, player.Nickname), 0)
	}

	args := fmt.Sprintf("player_id=%d&game_charge_id=0&charge_item_id=%d&item_count=%d&partid=%s&charge_type=%s&gm_id=%s",
		player.Id,
		params.ItemId,
		params.ChargeValue,
		params.PlatformId,
		params.ChargeType,
		c.curUser.Account,
	)
	sign := utils.String2md5(args + "fa9274fd68cf8991953b186507840e5e")
	g.Log().Infof("sign:%v", sign)

	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	url := models.GetGameURLByNode(gameServer.Node) + "/gm_charge?" + args + "&sign=" + sign
	g.Log().Infof("url:%v", url)
	resp, err := http.Get(url)
	c.CheckError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	c.CheckError(err)

	err = json.Unmarshal(body, &result)
	g.Log().Infof("result:%v", string(body))
	c.CheckError(err)
	g.Log().Infof("后台充值结果:%v", result)
	if result.Code == 0 {
		backgroundChargeLog := &models.BackgroundChargeLog{
			PlatformId:  params.PlatformId,
			ServerId:    string(player.ServerId),
			PlayerId:    params.PlayerId,
			Time:        utils.GetTimestamp(),
			ChargeType:  params.ChargeType,
			ChargeValue: params.ChargeValue,
			ItemId:      params.ItemId,
			UserId:      c.curUser.Id,
		}
		err = models.Db.Save(&backgroundChargeLog).Error
		c.CheckError(err, "写后台充值日志失败")
		c.HttpResult(r, enums.CodeSuccess, "后台充值成功", 0)
	}
	c.HttpResult(r, enums.CodeFail, fmt.Sprintf("后台充值失败: ErrorCode: %v Messsage", result.Code, result.Message), result.Code)
}
