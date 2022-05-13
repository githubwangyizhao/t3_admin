package models

import (
	"github.com/jinzhu/gorm"
)

type MailLog struct {
	Id             int    `json:"id"`
	PlatformId     string `json:"platformId"`
	NodeList       string `json:"serverIdList"`
	PlayerNameList string `json:"playerNameList"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	Time           int64  `json:"time"`
	UserId         int    `json:"userId"`
	ItemList       string `json:"itemList"`
	Status         int    `json:"status"`
	Type           string `json:"type"` //1:发送给多个玩家  2: 发送多个服务器 3: 全服发送
	UserName       string `json:"userName" gorm:"-"`
}

type MailLogQueryParam struct {
	BaseQueryParam
	PlatformId string
	ServerId   string
	StartTime  int
	EndTime    int
	UserId     int
}

func GetMailLogList(params *MailLogQueryParam) ([]*MailLog, int64) {
	data := make([]*MailLog, 0)
	var count int64
	sortOrder := "id"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	f := func(db *gorm.DB) *gorm.DB {
		if params.StartTime > 0 {
			return db.Where("time between ? and ?", params.StartTime, params.EndTime)
		}
		return db
	}
	f(Db.Model(&MailLog{}).Where(&MailLog{
		PlatformId: params.PlatformId,
		UserId:     params.UserId,
	}).Where("node_list LIKE ?", "%"+params.ServerId+"%")).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count)
	for _, e := range data {
		u, err := GetUserOne(e.UserId)
		if err == nil {
			e.UserName = u.Name
		}
	}
	return data, count
}

// 删除邮件日志
func DeleteMailLog(ids []int) error {
	err := Db.Where(ids).Delete(&MailLog{}).Error
	return err
}
