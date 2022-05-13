package models

import (
	"gfWeb/library/utils"

	"github.com/gogf/gf/os/gcron"
)

type NoticeLog struct {
	Id             int    `json:"id"`
	PlatformId     string `json:"platformId"`
	ServerIdList   string `json:"serverIdList"`
	Content        string `json:"content"`
	NoticeType     int    `json:"noticeType"`
	NoticeTime     int    `json:"noticeTime"`
	Status         int    `json:"status"`
	IsAllServer    int    `json:"isAllServer"`
	Time           int64  `json:"time"`
	UserId         int    `json:"userId"`
	UserName       string `json:"userName" gorm:"-"`
	LastSendTime   int    `json:"lastSendTime"`
	CreateCronTime int64  `json:"createCronTime"`
	CreateUserId   int    `json:"createUserId"`
	CronTimeStr    string `json:"cronTimeStr"`
	CronTimes      int    `json:"cronTimes"`
	SendTimes      int    `json:"sendTimes"`
}

type JiguangPushParam struct {
	Platform         string
	ServerIdList     []string
	Function         int
	FunctionConfigId int
	IsOpenUi         int
	Headline         string
	Desc             string
	CronTimeStr      string
}

type NoticeLogQueryParam struct {
	BaseQueryParam
	PlatformId   string
	StartTime    int
	EndTime      int
	UserId       int
	NoticeType   int
	IsShowFinish string
}

func GetNoticeLogOne(id int) (*NoticeLog, error) {
	noticeLog := &NoticeLog{
		Id: id,
	}
	err := Db.First(&noticeLog).Error
	return noticeLog, err
}

func GetAllNoticeLog() []*NoticeLog {
	data := make([]*NoticeLog, 0)
	Db.Model(&NoticeLog{}).Find(&data)
	return data
}

func GetNoticeLogList(params *NoticeLogQueryParam) ([]*NoticeLog, int64) {
	data := make([]*NoticeLog, 0)
	var count int64
	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "id")
	if params.Order == "ascending" {
		sortOrder = sortOrder + " asc"
	} else {
		sortOrder = sortOrder + " desc"
	}
	if params.IsShowFinish == "0" {
		Db.Model(&NoticeLog{}).Where(&NoticeLog{
			PlatformId: params.PlatformId,
			NoticeType: params.NoticeType,
		}).Where("status = ? ", 0).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count)
	} else {
		Db.Model(&NoticeLog{}).Where(&NoticeLog{
			PlatformId: params.PlatformId,
			NoticeType: params.NoticeType,
		}).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count)
	}

	for _, e := range data {
		u, err := GetUserOne(e.UserId)
		if err == nil {
			e.UserName = u.Name
		}
	}
	return data, count
}

// 删除公告日志
func DeleteNoticeLog(ids []int) error {
	for _, id := range ids {
		gcron.Remove(getCronNoticeName(id))
	}
	err := Db.Where(ids).Delete(&NoticeLog{}).Error
	return err
}
