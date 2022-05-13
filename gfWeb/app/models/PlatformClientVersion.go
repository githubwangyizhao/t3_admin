package models

import (
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"time"
)

type PlatformClientInfo struct {
	Id                int       `json:"id"`
	AppId             string    `json:"appId"`
	FacebookAppId     string    `json:"facebookAppId"`
	Platform          string    `json:"platform"`
	Channel           string    `json:"channel"`
	PlatformRemark    string    `json:"platformRemark"`
	ClientVersion     string    `json:"clientVersion"`
	FirstVersions     string    `json:"firstVersions"`
	ReviewingVersions string    `json:"reviewingversions"`
	Versions          string    `json:"versions"`
	IsChargeOpen      int       `json:"isChargeOpen"`
	NativePay         int       `json:"nativePay"`
	UpgradeIosUrl     string    `json:"upgradeIosUrl"`
	UpgradeAndroidUrl string    `json:"upgradeAndroidUrl"`
	ReloadUrl         string    `json:"reloadUrl"`
	Stats             int       `json:"stats"`
	CreatedAt         time.Time `json:"createdAt"`
	CreatedBy         int       `json:"createdBy"`
	CreatedName       string    `json:"createdName" gorm:"-"`
	UpdatedAt         time.Time `json:"updatedAt"`
	UpdatedBy         int       `json:"updatedBy"`
	UpdatedName       string    `json:"updatedName" gorm:"-"`
	PayTimes          int       `json:"payTimes" gorm:"pay_times"`
	Region            string    `json:"region" gorm:"region"`
	PackageSize       float64   `json:"packageSize" gorm:"packageSize"`
	AreaCode          string    `json:"areaCode" gorm:"area_code"`
	Domain            string    `json:"domain" gorm:"domain"`
	TestDomain        string    `json:"testDomain" gorm:"test_domain"`
}

type DoPushReq struct {
	Data string `json:"data"`
	Sign string `json:"sign"`
}
type ModifyPushReq struct {
	Platform       string `json:"platform"`
	Account        string `json:"account"`
	RegistrationId string `json:"registration_id"`
	Sid            string `json:"sid"` //多个区服逗号隔开
}

type PlatformClientInfoRequest struct {
	PlatformClientInfo
	BaseQueryParam
}

//// 获取单个推广员数据
//func GetPromoteDataOne(id int) (*PromoteData, error) {
//	promoteData := &PromoteData{
//		Id: id,
//	}
//	err := Db.Where(&promoteData).First(&promoteData).Error
//	return promoteData, err
//}
//
// 获取单个渠道客户端信息
func GetPlatformClientInfoOneByAppId(appId string) (*PlatformClientInfo, error) {
	platformClientInfo := &PlatformClientInfo{
		AppId: appId,
	}
	err := Db.Where(&platformClientInfo).First(&platformClientInfo).Error
	return platformClientInfo, err
}

func GetPlatformClientInfoById(id int) (*PlatformClientInfo, error) {
	platformClientInfo := &PlatformClientInfo{
		Id: id,
	}
	g.Log().Infof("platformClientInfo: %+v", platformClientInfo)
	err := Db.Debug().First(&platformClientInfo).Error
	return platformClientInfo, err
}

// 获取客户端版本列表
func GetPlatformClientInfoAllList() ([]*PlatformClientInfo, int64) {
	data := make([]*PlatformClientInfo, 0)
	var count int64
	err := Db.Find(&data).Count(&count).Error
	utils.CheckError(err)
	return data, count
}

// 获取渠道客户端信息列表
func GetPlatformClientInfoList(params *PlatformClientInfoRequest) ([]*PlatformClientInfo, int64) {
	data := make([]*PlatformClientInfo, 0)
	var count int64
	//data1 := make([]*PromoteData, 0)
	//err := Db.Find(&data1).Count(&count).Error
	//utils.CheckError(err)
	params.AppId = "%" + params.AppId + "%"
	params.Platform = "%" + params.Platform + "%"
	sortOrder := "id"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	ClientDb := Db.Where("app_id LIKE ?", params.AppId).Where("platform LIKE ?", params.Platform)
	if params.Region != "" {
		ClientDb = ClientDb.Where("region = ?", params.Region)
	}
	err := ClientDb.Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		e.CreatedName = GetUserName(e.CreatedBy)
		e.UpdatedName = GetUserName(e.UpdatedBy)
	}
	return data, count
}

// 新增渠道客户端信息数据
func AddPlatformClientInfo(platformClientInfo *PlatformClientInfo) error {
	err := Db.Save(platformClientInfo).Error
	return err
}

// 更新推广员数据
func UpdatePlatformClientInfo(platformClientInfo *PlatformClientInfo) error {
	err := Db.Save(platformClientInfo).Error
	return err
}

// DeletePlatformClientInfo 删除用户列表
func DeletePlatformClientInfo(ids []int) error {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := Db.Where(ids).Delete(&PlatformClientInfo{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error

}
