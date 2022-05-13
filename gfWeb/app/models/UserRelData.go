package models

import (
	"gfWeb/library/utils"

	"github.com/gogf/gf/os/gtime"
	"github.com/jinzhu/gorm"
)

// 用户关联数据
type UserRelData struct {
	Type       int   `json:"type" gorm:"primary_key"`
	Id         int   `json:"id" gorm:"primary_key"`
	UserId     int   `json:"userId" gorm:"primary_key"`
	CreateTime int64 `json:"createTime"`
}

func userRelDataTBName() string {
	return "user_rel_data"
}

const (
	USER_REL_DATA_TYPE_MsgTemplate_MAIL  = 1 // 后台信息模板-邮件（BackgroundMsgTemplate）
	USER_REL_DATA_TYPE_MsgTemplate_PHONE = 2 // 后台信息模板-手机（BackgroundMsgTemplate）
	USER_REL_DATA_TYPE_PAGE_CHANGE_AUTH  = 3 // 界面操作权限（PageChangeAuth）
)

// 获取用户关联数据列表
func GetUserRelDataList(relType, id int) []*UserRelData {
	data := make([]*UserRelData, 0)
	userRelData := &UserRelData{
		Type: relType,
		Id:   id,
	}
	sortOrder := "user_id"
	err := Db.Model(userRelData).Where(userRelData).Order(sortOrder).Find(&data).Error
	utils.CheckError(err)
	return data
}

// 获取用户关联数据用户列表
func GetUserRelDataToUserIdList(relType, id int) []int {
	data := make([]int, 0)
	userRelData := &UserRelData{
		Type: relType,
		Id:   id,
	}
	sortOrder := "user_id"
	err := Db.Table(userRelDataTBName()).Select([]string{"user_id"}).Where(&userRelData).Order(sortOrder).Pluck("user_id", &data).Error
	utils.CheckError(err)
	return data
}

// 更新关联用户数据
func UpdateUserRelDataByUserIdList(tx *gorm.DB, relType, id int, oldUserIdList, newUserIdList []int) error {
	delUserIdList := utils.ListMinus(oldUserIdList, newUserIdList)
	addUserIdList := utils.ListMinus(newUserIdList, oldUserIdList)
	if _, err := DeleteUserRelDataByTypeIdUserIdList(relType, id, delUserIdList); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := SaveUserRelDataByUserIdList(relType, id, addUserIdList); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// 保存用户的关联数据的用户列表
func SaveUserRelDataByUserIdList(relType, id int, userIdList []int) (int, error) {
	var count int
	for _, userId := range userIdList {
		userRelData := &UserRelData{
			Type:       relType,
			Id:         id,
			UserId:     userId,
			CreateTime: gtime.Timestamp(),
		}
		err := Db.Save(userRelData).Error
		if err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

// 删除用户的关联数据
func DeleteUserRelDataByUserIdList(userIdList []int) (int, error) {
	var count int
	err := Db.Where("user_id in (?)", userIdList).Delete(&UserRelData{}).Count(&count).Error
	return count, err
}

// 删除关联数据的数据
func DeleteUserRelDataByTypeIdList(relType int, idList []int) (int, error) {
	var count int
	userRelData := &UserRelData{
		Type: relType,
	}
	err := Db.Where(userRelData).Where("id in (?)", idList).Delete(&UserRelData{}).Count(&count).Error
	return count, err
}

// 删除关联用户数据类型的id用户列表数据
func DeleteUserRelDataByTypeIdUserIdList(relType, id int, userIdList []int) (int, error) {
	var count int
	userRelData := &UserRelData{
		Type: relType,
		Id:   id,
	}
	err := Db.Where(userRelData).Where("user_id in (?)", userIdList).Delete(&UserRelData{}).Count(&count).Error
	return count, err
}
