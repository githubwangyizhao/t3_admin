package models

import (
	"encoding/json"
	"gfWeb/library/utils"
	"strings"

	"github.com/gogf/gf/os/gtime"
)

type LoginNotice struct {
	Id         string `gorm:"primary_key" json:"id"`
	PlatformId string `json:"platformId"`
	ChannelId  string `json:"channelId"`
	Notice     string `json:"notice"`
	Time       int64  `json:"time"`
	UserId     int    `json:"userId"`
	UserName   string `json:"userName" gorm:"-"`
}

type LoginNoticeQueryParam struct {
	BaseQueryParam
	PlatformId string
}

func GetAllLoginNotice() []*LoginNotice {
	data := make([]*LoginNotice, 0)
	Db.Model(&LoginNotice{}).Find(&data)
	for _, e := range data {
		u, err := GetUserOne(e.UserId)
		if err == nil {
			e.UserName = u.Name
		}
	}
	return data
}

func GetLoginNoticeListByPlatformIdList(platformIdList []string) []*LoginNotice {
	data := make([]*LoginNotice, 0)
	if len(platformIdList) == 0 {
		return data
	}
	Db.Model(&LoginNotice{}).Where("platform_id in (?)", platformIdList).Find(&data)
	for _, e := range data {
		u, err := GetUserOne(e.UserId)
		if err == nil {
			e.UserName = u.Name
		}
	}
	return data
}

func BatchSetNotice(userId int, ids []string, notice string) error {
	for _, IdStr := range ids {
		list := strings.Split(IdStr, getLoginNoticeJoinStr())
		platformId := list[0]
		channelId := list[1]
		err := UpdateAndPushLoginNotice(userId, platformId, channelId, notice)
		utils.CheckError(err, "写登录公告日志失败")
		return err
	}
	return nil
}

// 删除登录公告日志
func DeleteLoginNotice(userId int, ids []string) error {
	for _, IdStr := range ids {
		list := strings.Split(IdStr, getLoginNoticeJoinStr())
		platformId := list[0]
		channelId := list[1]
		err := UpdateAndPushLoginNotice(userId, platformId, channelId, "")
		utils.CheckError(err, "写登录公告失败")
		return err
	}
	err := Db.Where("platform_id in (?)", ids).Delete(&LoginNotice{}).Error
	return err
}

func UpdateAndPushLoginNotice(userId int, platformId, channelId string, notice string) error {
	var request struct {
		PlatformId string `json:"platformId"`
		ChannelId  string `json:"channelId"`
		Notice     string `json:"notice"`
	}
	request.PlatformId = platformId
	request.Notice = notice
	data, err := json.Marshal(request)
	utils.CheckError(err)
	if err != nil {
		return err
	}
	url := utils.GetCenterURL() + "/set_login_notice"
	_, err = utils.HttpRequest(url, string(data))
	if err != nil {
		return err
	}
	noticeLog := &LoginNotice{
		Id:         platformId + getLoginNoticeJoinStr() + channelId,
		PlatformId: platformId,
		ChannelId:  channelId,
		Notice:     notice,
		Time:       gtime.Timestamp(),
		UserId:     userId,
	}
	err = Db.Save(&noticeLog).Error
	return err
}

// 登录公告连接符
func getLoginNoticeJoinStr() string {
	return "_^_"
}
