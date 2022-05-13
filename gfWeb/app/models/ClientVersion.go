package models

import (
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"net/url"
)

// 获得客户当前版本
func GetClientVersion(PlatformId, Channel string) string {
	URL := "https://web-rxxx.szfyhd.com/client_version"
	//初始化参数
	param := url.Values{}
	param.Set("platform_id", PlatformId) // 平台
	param.Set("channel", Channel)        // 渠道
	BobyBin, err := utils.HttpGet(URL, param)
	if err != nil {
		g.Log().Error("获得客户当前版本错误:err:%v", err)
		return ""
	}
	return string(BobyBin)
}

// 更新客户当前版本
func UpdateClientVersion(PlatformId, Channel, ClientVersion string) string {
	URL := "https://web-rxxx.szfyhd.com/update_client_version"
	//初始化参数
	param := url.Values{}
	param.Set("platform_id", PlatformId) // 平台
	param.Set("channel", Channel)        // 渠道
	param.Set("version", ClientVersion)  // 版本
	BobyBin, err := utils.HttpGet(URL, param)
	if err != nil {
		g.Log().Error("更新客户当前版本错误:err:%v", err)
		return ""
	}
	return string(BobyBin)
}
