package models

import (
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type PageChangeAuth struct {
	Id         int    `json:"id"`
	Sign       string `json:"sign"`
	Name       string `json:"name"`
	State      int    `json:"state"`
	ChangeTime int64  `json:"changeTime"`
	UserId     int    `json:"userId"`
	UserList   []int  `json:"userList" gorm:"-"`
	IsAdd      int    `json:"isAdd" gorm:"-"`
}

const (
	PAGE_CHANGE_AUTH_IS_MERGE       = "merge_auth"  // 合服审核操作权限
	PAGE_CHANGE_AUTH_IS_PACK_SERVER = "pack_server" // 服务端打包
	PAGE_CHANGE_AUTH_IS_PACK_CLIENT = "pack_client" // 客户端打包
)

// 设置是否开启并有默认值
func IsPageChangeAuthDefault(u User, signStr string, defaultStr bool) bool {
	if u.IsSuperUser() {
		return true
	}
	pageChangeAuth, err := GetPageChangeAuthOneBySign(signStr)
	if err != nil {
		g.Log().Debugf("设置是否开启并有默认值:%+v", signStr)
		return defaultStr
	}
	if pageChangeAuth.State == 0 {
		g.Log().Debugf("设置是否开启State:%+v", signStr)
		return false
	}
	if !utils.IsHaveIntArray(u.Id, pageChangeAuth.UserList) {
		g.Log().Debugf("设置是未开启userId:%+v", u.Id)
		return false
	}
	return true
}

// 获取单个设置数据
func GetPageChangeAuthOneBySign(signStr string) (*PageChangeAuth, error) {
	pageChangeAuth := &PageChangeAuth{
		Sign: signStr,
	}
	err := Db.Where(pageChangeAuth).First(pageChangeAuth).Error
	if err != nil {
		return pageChangeAuth, err
	}
	pageChangeAuth.UserList = GetUserRelDataToUserIdList(USER_REL_DATA_TYPE_PAGE_CHANGE_AUTH, pageChangeAuth.Id)
	return pageChangeAuth, err
}

// 获取设置数据列表
func GetPageChangeAuthList(params *BaseQueryParam) ([]*PageChangeAuth, int64) {
	data := make([]*PageChangeAuth, 0)
	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "id")
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	err := Db.Model(&PageChangeAuth{}).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		e.UserList = GetUserRelDataToUserIdList(USER_REL_DATA_TYPE_PAGE_CHANGE_AUTH, e.Id)
	}
	return data, count
}

// 更新设置数据
func UpdatePageChangeAuth(pageChangeAuth, OldAuthData *PageChangeAuth) error {
	pageChangeAuth.ChangeTime = gtime.Timestamp()
	//err := Db.Save(pageChangeAuth).Error
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	if err := UpdateUserRelDataByUserIdList(tx, USER_REL_DATA_TYPE_PAGE_CHANGE_AUTH, pageChangeAuth.Id, OldAuthData.UserList, pageChangeAuth.UserList); err != nil {
		return err
	}
	if err := Db.Save(&pageChangeAuth).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// 删除设置数据
func DeletePageChangeAuth(ids []int) error {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := Db.Where(ids).Delete(&PageChangeAuth{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if _, err := DeleteUserRelDataByTypeIdList(USER_REL_DATA_TYPE_PAGE_CHANGE_AUTH, ids); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
