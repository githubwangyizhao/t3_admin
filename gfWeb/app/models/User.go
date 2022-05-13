package models

import (
	"fmt"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/util/gconv"
	"sort"
	"strings"
)

func (a *User) TableName() string {
	return UserTBName()
}

func UserTBName() string {
	return TableName("user")
}

type UserQueryParam struct {
	BaseQueryParam
	Account string
	Role    int
}
type User struct {
	Id                      int            `json:"id"`
	Name                    string         `json:"name"`
	Account                 string         `json:"account"`
	Password                string         `json:"-"`
	IsSuper                 int            `json:"isSuper"`
	ModifyPassword          string         `json:"Password" gorm:"-"`
	Status                  int            `json:"status"`
	LoginTimes              int            `json:"loginTimes"`
	LastLoginTime           int            `json:"lastLoginTime"`
	LastLoginIp             string         `json:"lastLoginIp"`
	CanLoginTime            int            `json:"canLoginTime"`
	ContinueLoginErrorTimes int            `json:"-"`
	MailStr                 string         `json:"mailStr"`
	Mobile                  string         `json:"mobile"`
	RoleIds                 []int          `json:"roleIds" gorm:"-"`
	RoleUserRel             []*RoleUserRel `json:"-"`
	ResourceUrlForList      []string       `gorm:"-"`
}

// 精简用户数据
type UserSimpleData struct {
	Id   int    `json:"id" gorm:"-"`
	Name string `json:"name" gorm:"-"`
}

//获取用户列表
func GetUserList(params *UserQueryParam) ([]*User, int64) {
	data := make([]*User, 0)

	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "id")
	//sortOrder := "id"
	//switch params.Sort {
	//case "id":
	//	sortOrder = "id"
	//}
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	whereArray := make([]string, 0)
	if params.Role > 0 {
		whereArray = append(whereArray, fmt.Sprintf(" id in ( select user_id from myadmin_role_user_rel where role_id = %d)", params.Role))
	}
	whereParam := ""
	if len(whereArray) > 0 {
		whereParam = strings.Join(whereArray, " and ")
	}
	var count int64
	err := Db.Model(&User{}).Where(&User{
		Account: params.Account,
	}).Where(whereParam).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, v := range data {
		err = Db.Model(&v).Related(&v.RoleUserRel).Error
		utils.CheckError(err)
		roleIds := make([]int, 0)
		for _, e := range v.RoleUserRel {
			roleIds = append(roleIds, e.RoleId)
		}
		sort.Ints(roleIds)
		v.RoleIds = roleIds
	}
	return data, count
}

//获取用户精简的数据列表
func GetUserSimpleList() []*UserSimpleData {
	data := make([]*UserSimpleData, 0)
	sortOrder := "id"
	err := Db.Table(UserTBName()).Select([]string{"id", "name"}).Order(sortOrder).Find(&data).Error
	utils.CheckError(err)
	return data
}

// 获取单个用户
func GetUserOne(id int) (*User, error) {
	user := &User{
		Id: id,
	}
	err := Db.First(&user).Error
	return user, err
}

// 获取单个用户的名称
func GetUserName(userId int) string {
	u, err := GetUserOne(userId)
	if err == nil {
		return u.Name
	}
	return gconv.String(userId)
}

//是否帐号有效
func (u *User) IsAccountEnable() bool {
	return u.Status == enums.Enabled
}

// 是否是超级管理员
func (u *User) IsSuperUser() bool {
	return u.IsSuper == 1
}

// 根据用户名单条
func GetUserOneByAccount(account string) (*User, error) {
	user := &User{}
	isNotFound := Db.Where(&User{Account: account}).First(&user).RecordNotFound()
	if isNotFound {
		return nil, gerror.New("用户不存在")
	}
	return user, nil
}

// 根据用户名密码获取单条
func GetUserOneByAccountAndPassword(account, password string) (*User, error) {
	user := &User{}
	isNotFound := Db.Where(&User{Account: account, Password: password}).First(&user).RecordNotFound()
	if isNotFound {
		return nil, gerror.New("用户名或者密码错误")
	}
	return user, nil
}

// 删除用户列表
func DeleteUsers(ids []int) error {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := Db.Where(ids).Delete(&User{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if _, err := DeleteRoleUserRelByUserIdList(ids); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := DeleteUserRelDataByUserIdList(ids); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
