package models

import (
	"gfWeb/library/utils"
	"github.com/gogf/gf/os/gtime"
)

type SettingData struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	State      int    `json:"state"`
	ChangeTime int64  `json:"changeTime"`
	UserId     int    `json:"userId"`
	UserName   string `json:"userName" gorm:"-"`
	IsAdd      int    `json:"isAdd" gorm:"-"`
}

const (
	SETTING_DATA_IS_OPEN_MAIL      = 1 // 是否开启邮件通知
	SETTING_DATA_IS_OPEN_SMS       = 2 // 是否开启普通短信通知
	SETTING_DATA_IS_OPEN_ALL_SMS   = 3 // 是否开启全部短信通知
	SETTING_DATA_IS_CHECK_BACK_WEB = 4 // 是否开启检测连接数后台web
)

// 设置是否开启
func IsSettingOpen(id int) bool {
	return IsSettingOpenDefault(id, false)
}

// 设置是否开启并有默认值
func IsSettingOpenDefault(id int, defaultStr bool) bool {
	SettingData, err := GetSettingDataOne(id)
	if err != nil {
		return defaultStr
	}
	return SettingData.State == 1
}

// 获取单个设置数据
func GetSettingDataOne(id int) (*SettingData, error) {
	settingData := &SettingData{
		Id: id,
	}
	err := Db.Where(&settingData).First(&settingData).Error
	return settingData, err
}

// 获取设置数据列表
func GetSettingDataList(params *BaseQueryParam) ([]*SettingData, int64) {
	data := make([]*SettingData, 0)
	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "id")
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	err := Db.Model(&SettingData{}).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		e.UserName = GetUserName(e.UserId)
	}
	return data, count
}

// 更新设置数据
func UpdateSettingDate(settingData *SettingData) error {
	settingData.ChangeTime = gtime.Timestamp()
	err := Db.Save(settingData).Error
	return err
}

// 删除设置数据
func DeleteSettingData(ids []int) error {
	err := Db.Where(ids).Delete(&SettingData{}).Error
	return err
}
