package models

import (
	"time"
)

//角色与资源关系表
type RoleResourceRel struct {
	Id         int
	RoleId     int
	ResourceId int
	Created    time.Time
}

func (a *RoleResourceRel) TableName() string {
	return RoleResourceRelTBName()
}

func RoleResourceRelTBName() string {
	return TableName("role_resource_rel")
}

// 删除角色资源关系
func DeleteRoleResourceRelByRoleIdList(roleIdList []int) (int, error) {
	var count int
	err := Db.Where("role_id in (?)", roleIdList).Delete(&RoleResourceRel{}).Count(&count).Error
	return count, err
}

// 删除角色资源关系
func DeleteRoleResourceRelByResourceIdList(resourceIdList []int) (int, error) {
	var count int
	err := Db.Where("resource_id in (?)", resourceIdList).Delete(&RoleResourceRel{}).Count(&count).Error
	return count, err
}
