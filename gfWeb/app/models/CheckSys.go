package models

// 检测系统功能
import (
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"net/url"
	"strings"
)

// 检测后台web活跃数据
func CheckBackWebActive() {
	g.Log().Debugf("检测后台web活跃数据")
	isCheckBackWeb := IsSettingOpenDefault(SETTING_DATA_IS_CHECK_BACK_WEB, false)
	//isCheckBackWeb := g.Cfg().GetBool("game.isCheckBackWeb")
	if isCheckBackWeb == true {
		PlatformType := g.Cfg().GetString("game.platform_type")
		CheckWebStr := g.Cfg().GetString("game.check_back_web")
		if CheckWebStr == "" {
			return
		}
		line := strings.Split(CheckWebStr, "|")
		for _, CheckWebInfoStr := range line {
			CheckWebInfo := strings.Split(CheckWebInfoStr, ",")
			Platform := CheckWebInfo[0]
			if PlatformType == Platform {
				continue
			}
			URL := CheckWebInfo[1]
			_, err := utils.HttpGetTimeout(URL, url.Values{}, 20)
			if err != nil {
				g.Log().Errorf("检测后台web连接错误: %v\terr:%v", CheckWebInfo, err)
				//SendSmsAdminUserMsg("SMS_164095721", Platform)
				//SendMailAdminUser("检测后台web活跃数据", Platform+"平台后台web断开连接，地址为"+URL)
				//service.SendPhoneMsgTemplate("SMS_164095721", Platform, MSG_TEMPLATE_CHECK_BACK_WEB)
				SendBackgroundMsgTemplateHandle(MSG_TEMPLATE_CHECK_BACK_WEB, Platform, "检测后台web活跃数据", Platform+"平台后台web断开连接，地址为"+URL)
				continue
			}
		}
		g.Log().Infof("PlatformType:%s, %v   ", PlatformType, line)
	}
}
