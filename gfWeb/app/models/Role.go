package models

import (
	"errors"
	"gfWeb/library/utils"
	"sort"
	"strings"

	"github.com/gogf/gf/util/gconv"
)

func (a *Role) TableName() string {
	return RoleTBName()
}

func RoleTBName() string {
	return TableName("role")
}

type RoleQueryParam struct {
	BaseQueryParam
}

//用户角色
type Role struct {
	Id              int                `json:"id"`
	Name            string             `json:"name"`
	ResourceIds     []int              `json:"resourceIds" gorm:"-"`
	MenuIds         []int              `json:"menuIds" gorm:"-"`
	ChannelIds      []int              `json:"channelIds" gorm:"-"`
	RoleResourceRel []*RoleResourceRel `json:"-"`
	RoleChannelRel  []*RoleChannelRel  `json:"-"`
	RoleMenuRel     []*RoleMenuRel     `json:"-"`
}

//获取角色列表
func GetRoleList(params *RoleQueryParam) ([]*Role, int64) {
	data := make([]*Role, 0)
	var count int64
	err := Db.Model(&Role{}).Find(&data).Count(&count).Error
	utils.CheckError(err)
	for _, v := range data {
		err = Db.Model(&v).Related(&v.RoleResourceRel).Error
		utils.CheckError(err)

		v.MenuIds = make([]int, 0)
		v.ResourceIds = make([]int, 0)
		v.ChannelIds = make([]int, 0)

		err = Db.Model(&v).Related(&v.RoleMenuRel).Error
		utils.CheckError(err)

		err = Db.Model(&v).Related(&v.RoleChannelRel).Error
		utils.CheckError(err)

		for _, e := range v.RoleResourceRel {
			v.ResourceIds = append(v.ResourceIds, e.ResourceId)
		}
		sort.Ints(v.ResourceIds)

		for _, e := range v.RoleMenuRel {
			v.MenuIds = append(v.MenuIds, e.MenuId)
		}
		sort.Ints(v.MenuIds)

		for _, e := range v.RoleChannelRel {
			v.ChannelIds = append(v.ChannelIds, e.ChannelId)
		}
		sort.Ints(v.ChannelIds)
	}
	return data, count
}

//删除角色
func DeleteRoles(ids []int) error {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := Db.Where(ids).Delete(&Role{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if _, err := DeleteRoleResourceRelByRoleIdList(ids); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := DeleteRoleChannelRelByRoleIdList(ids); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := DeleteRoleMenuRelByRoleIdList(ids); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := DeleteRoleUserRelByRoleIdList(ids); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//获取单个角色
func GetRoleOne(id int) (role *Role, err error) {
	role = &Role{
		Id: id,
	}
	err = Db.First(&role).Error
	return role, err
}

func RoleIsSuper(roleId int) (bool, error) {

	role, err := GetRoleOne(roleId)
	if err != nil {
		return false, err
	}

	if strings.Index(role.Name, "超级管理员") == -1 {
		return false, errors.New("id: " + gconv.String(roleId) + ",非超级管理员")
	}

	return true, nil
}
