package models

import (
	"encoding/json"
	"gfWeb/library/utils"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

type SendNoticeParams struct {
	Id           int
	PlatformId   string
	ServerIdList []string `json:"serverIdList"`
	IsAllServer  int
	Content      string
	NoticeType   int
	NoticeTime   int
	CronTimeStr  string
	CronTimes    int
}

// 获得广播名字
func (n NoticeLog) getCronName() string {
	return CRON_NAME_NOTICE_NOTICE + gconv.String(n.Id)
}
func getCronNoticeName(Id int) string {
	return CRON_NAME_NOTICE_NOTICE + gconv.String(Id)
}

// 获得广播是否结束
func getIsNoticeClose(noticeLog *NoticeLog) bool {
	return noticeLog.Status == 1 || noticeLog.CronTimes > 0 && noticeLog.CronTimes <= noticeLog.SendTimes
}

// 初始广播内容订时器
func InitNoticeLog() {
	g.Log().Infof("初始广播内容订时器:%+v", gtime.Datetime())
	noticeLogList := GetAllNoticeLog()
	for _, noticeLog := range noticeLogList {
		if getIsNoticeClose(noticeLog) == false {
			g.Log().Infof("初始发送广播内容订时：%+v", noticeLog)
			// 处理未完成的公告
			StartNoticeLog(noticeLog.Id)
		}
	}
}

// 创建广播
func CreateNoticeLog(params SendNoticeParams, user User) error {
	if CheckCronTimeStr(params.CronTimeStr) != true {
		return gerror.New("CronTimeStr格式不对:" + params.CronTimeStr)
	}
	ServerIdList, err := json.Marshal(params.ServerIdList)
	if err != nil {
		return gerror.New("ServerIdList格式不对")
	}

	currTime := gtime.Timestamp()
	noticeLog, err := GetNoticeLogOne(params.Id)
	g.Log().Warningf("初始创建广播noticeLog：%d %+v", params.Id, noticeLog)
	if err != nil {
		g.Log().Warningf("初始创建广播：%d %+v", params.Id, err)
		noticeLog = &NoticeLog{
			PlatformId:     params.PlatformId,
			ServerIdList:   gconv.String(ServerIdList),
			IsAllServer:    params.IsAllServer,
			Content:        params.Content,
			CreateCronTime: currTime,
			Time:           currTime,
			CreateUserId:   user.Id,
			UserId:         user.Id,
			NoticeType:     params.NoticeType,
			NoticeTime:     params.NoticeTime,
			CronTimeStr:    params.CronTimeStr,
			CronTimes:      params.CronTimes,
			Status:         0,
		}
	} else {
		noticeLog.PlatformId = params.PlatformId
		noticeLog.ServerIdList = gconv.String(ServerIdList)
		noticeLog.IsAllServer = params.IsAllServer
		noticeLog.Content = params.Content
		noticeLog.Time = currTime
		noticeLog.UserId = user.Id
		noticeLog.NoticeType = params.NoticeType
		noticeLog.NoticeTime = params.NoticeTime
		noticeLog.CronTimeStr = params.CronTimeStr
		noticeLog.CronTimes = params.CronTimes
	}
	err = Db.Save(noticeLog).Error
	if err != nil {
		return err
	}
	gcron.Remove(noticeLog.getCronName())
	err = StartNoticeLog(noticeLog.Id)
	if err != nil {
		return err
	}
	return err
}

// 移除广播
func RemoveCronNotice(idList []int, user User) error {
	currTime := gtime.Timestamp()
	for _, id := range idList {
		noticeLog, err := GetNoticeLogOne(id)
		if err != nil {
			g.Log().Errorf("获取公告(%v)失败：%+v", id, err)
			continue
		}
		gcron.Remove(noticeLog.getCronName())
		if getIsNoticeClose(noticeLog) == true {
			g.Log().Warningf("公告已结束：%d", id)
			continue
		}
		noticeLog.Time = currTime
		noticeLog.UserId = user.Id
		noticeLog.Status = 1
		err = Db.Save(&noticeLog).Error
		utils.CheckError(err, "保存移除广播失败")
	}
	return nil
}

// 启动广播
func StartNoticeLog(id int) error {
	noticeLog, err := GetNoticeLogOne(id)
	if err != nil {
		g.Log().Errorf("获取公告(%v)失败：%+v", id, err)
		return err
	}
	if getIsNoticeClose(noticeLog) == true {
		g.Log().Warningf("公告已结束：%d", id)
		return gerror.New("公告已结束")
	}
	cronFun := func() {
		SendNoticeLog(id)
	}
	g.Log().Infof("启动广播CronTimes:%d CronTimeStr:%s  Datetime:%s", noticeLog.CronTimes, noticeLog.CronTimeStr, gtime.Datetime())
	if noticeLog.CronTimes > 0 {
		_, err = gcron.AddTimes(noticeLog.CronTimeStr, noticeLog.CronTimes, cronFun, noticeLog.getCronName())
		utils.CheckError(err)
		return err
	} else {
		_, err = gcron.Add(noticeLog.CronTimeStr, cronFun, noticeLog.getCronName())
		utils.CheckError(err)
		return err
	}
	return nil
}

// 发送广播
func SendNoticeLog(id int) {
	noticeLog, err := GetNoticeLogOne(id)
	if err != nil {
		g.Log().Errorf("获取公告(%v)失败：%+v", id, err)
		return
	}
	g.Log().Infof("发送广播: %d  %s", id, noticeLog.CronTimeStr)
	now := int(gtime.Timestamp())

	//g.Log().Info("处理公告:%+v", noticeLog)
	nodeList := make([]string, 0)
	if noticeLog.IsAllServer == 0 {
		serverIdList := make([]string, 0)
		err := json.Unmarshal([]byte(noticeLog.ServerIdList), &serverIdList)
		if err != nil {
			g.Log().Error("解析公告(%+v)区服列表失败%+v", id, err)
			return
		}
		nodeList = GetNodeListByServerIdList(noticeLog.PlatformId, serverIdList)
	} else {
		// 全服
		nodeList = GetAllGameNodeByPlatformId(noticeLog.PlatformId)
	}

	var request struct {
		NodeList []string `json:"nodeList"`
		Content  string   `json:"content"`
	}
	request.NodeList = nodeList
	request.Content = noticeLog.Content

	for _, node := range nodeList {
		data, err := json.Marshal(request)
		utils.CheckError(err)
		url := GetGameURLByNode(node) + "/send_notice"

		resultMsg, err := utils.HttpRequest(url, string(data))
		utils.CheckError(err, resultMsg)
		if err == nil {
			g.Log().Debugf("发送公告成功 node: %v, id: %v", node, noticeLog.Id)
		} else {
			g.Log().Errorf("发送公告失败 PlatformId: %v, node: %v, id: %v", noticeLog.PlatformId, node, noticeLog.Id)
		}
	}
	//if noticeLog.NoticeType != enums.NoticeTypeLoop {
	//	noticeLog.Status = 1
	//}
	noticeLog.LastSendTime = now
	noticeLog.SendTimes++

	if noticeLog.CronTimes > 0 && noticeLog.CronTimes <= noticeLog.SendTimes {
		noticeLog.Status = 1
	}
	err = Db.Save(&noticeLog).Error
	utils.CheckError(err, "保存公告日志失败")

}
