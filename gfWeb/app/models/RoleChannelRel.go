package models

//角色与渠道关系表
type RoleChannelRel struct {
	Id        int
	RoleId    int
	ChannelId int
	//Created   time.Time
}

func (a *RoleChannelRel) TableName() string {
	return RoleChannelRelTBName()
}

func RoleChannelRelTBName() string {
	return TableName("role_channel_rel")
}

// 删除角色渠道关系
func DeleteRoleChannelRelByRoleIdList(roleIdList []int) (int, error) {
	var count int
	err := Db.Where("role_id in (?)", roleIdList).Delete(&RoleChannelRel{}).Count(&count).Error
	return count, err
}

// 删除角色渠道关系
func DeleteRoleChannelRelByPlatformIdList(platformIdList []string) (int, error) {
	var count int
	err := Db.Where("platform_id in (?)", platformIdList).Delete(&RoleChannelRel{}).Count(&count).Error
	return count, err
}
