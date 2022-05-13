package models

//角色菜单关系表
type RoleMenuRel struct {
	Id     int
	RoleId int
	MenuId int
	//Created time.Time
}

func (a *RoleMenuRel) TableName() string {
	return RoleMenuRelTBName()
}

func RoleMenuRelTBName() string {
	return TableName("role_menu_rel")
}

// 删除角色菜单关系
func DeleteRoleMenuRelByRoleIdList(roleIdList []int) (int, error) {
	var count int
	err := Db.Where("role_id in (?)", roleIdList).Delete(&RoleMenuRel{}).Count(&count).Error
	return count, err
}

// 删除角色菜单关系
func DeleteRoleMenuRelByMenuIdList(menuIdList []int) (int, error) {
	var count int
	err := Db.Where("menu_id in (?)", menuIdList).Delete(&RoleMenuRel{}).Count(&count).Error
	return count, err
}
