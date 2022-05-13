package models

import (
	"gfWeb/library/utils"
	//"sort"
)

type Channel struct {
	Id           int    `json:"id"`
	PlatformId   string `json:"platformId"`
	Channel      string `json:"channel"`
	Name         string `json:"name"`
	TrackerToken string `json:"trackerToken"gorm:"tracker_token"`
	Region       string `json:"region" gorm:"region"`
	AreaCode     string `json:"areaCode" gorm:"area_code"`
	Currency     string `json:"currency" gorm:"currency"`
}

func (a *Channel) TableName() string {
	return ChannelDatabaseTBName()
}

func ChannelDatabaseTBName() string {
	return "channel"
}

type ChannelParam struct {
	BaseQueryParam
}

//获取渠道列表
func GetChannelListByPlatformIdList(platformIdList []string) []*Channel {
	data := make([]*Channel, 0)
	err := Db.Model(&Channel{}).Where("platform_id in (?)", platformIdList).Find(&data).Error
	utils.CheckError(err)
	return data
}

func GetChannelList() []*Channel {
	data := make([]*Channel, 0)
	err := Db.Model(&Channel{}).Find(&data).Error
	utils.CheckError(err)
	return data
}

//获取渠道列表
func GetChannelListByPlatformId(platformId string) []*Channel {
	data := make([]*Channel, 0)
	err := Db.Model(&Channel{}).Where(&Channel{PlatformId: platformId}).Find(&data).Error
	utils.CheckError(err)
	//for _, e := range  data {
	//platform, err := GetPlatformOne(e.PlatformId)
	//utils.CheckError(err)
	//e.Name = platform.Name + "-" + e.Name
	//}
	return data
}

//获取单个平台
func GetChannelOne(id int) (*Channel, error) {
	r := &Channel{
		Id: id,
	}
	err := Db.First(&r).Error
	return r, err
}

// 删除渠道列表
func DeleteChannel(ids []int) error {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := Db.Where(ids).Delete(&Channel{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// 删除平台渠道关系
func DeleteChannelRelByPlatformIdList(platformIdList []string) (int, error) {
	var count int
	err := Db.Where("platform_id in (?)", platformIdList).Delete(&Channel{}).Count(&count).Error
	return count, err
}
