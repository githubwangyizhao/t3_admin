package models

import (
	"time"
)

//角色与用户关系
type RoleUserRel struct {
	Id      int
	UserId  int
	RoleId  int
	Created time.Time
}

func (a *RoleUserRel) TableName() string {
	return RoleUserRelTBName()
}

func RoleUserRelTBName() string {
	return TableName("role_user_rel")
}

// 删除角色用户关系
func DeleteRoleUserRelByUserIdList(userIdList []int) (int, error) {
	var count int
	err := Db.Where("user_id in (?)", userIdList).Delete(&RoleUserRel{}).Count(&count).Error
	return count, err
}

// 删除角色用户关系
func DeleteRoleUserRelByRoleIdList(roleIdList []int) (int, error) {
	var count int
	err := Db.Where("role_id in (?)", roleIdList).Delete(&RoleUserRel{}).Count(&count).Error
	return count, err
}
