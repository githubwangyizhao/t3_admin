package models

import (
	"fmt"
	"gfWeb/library/utils"
	"strconv"
	"strings"
)

type ServerRel struct {
	Id           int    `json:"id"`
	SAccount     int    `json:"s_account" gorm:"primary_key"`
	Uid          int    `json:"uid" gorm:"primary_key"`
	PlatId       string `json:"plat_id"`
	ServerId     string `json:"server_id"`
	CreateTime   int64  `json:"create_time"`
	Uptime       int64  `json:"uptime"`
	Remark       string `json:"remark"`
	CustomerName string `json:"customer_name"`
	Uname        string `json:"nick" gorm:"_"`
}

type InServerRel struct {
	Id         int    `json:"id"`
	SAccount   int    `json:"s_account" gorm:"primary_key"`
	Uid        int    `json:"uid" gorm:"primary_key"`
	PlatId     string `json:"plat_id"`
	ServerId   string `json:"server_id"`
	CreateTime int64  `json:"create_time"`
	Uptime     int64  `json:"uptime"`
	Remark     string `json:"remark"`
}

type AdminUser struct {
	Uid  int    `json:"uid"`
	Name string `json:"name"`
}

func ServerRelTBName() string {
	return "server_rel"
}

func (a *InServerRel) TableName() string {
	return ServerRelTBName()
}

func (a *ServerRel) TableName() string {
	return ServerRelTBName()
}

func GetGlobalAccountIdBySAccount(SAccountList []string) (*[]string, error) {
	whereArray := make([]string, 0)

	whereArray = append(whereArray, fmt.Sprintf(" s_account IN (%s)", strings.Join(SAccountList, ",")))
	whereParam := strings.Join(whereArray, " and ")
	sql := fmt.Sprintf(
		`select * from server_rel where %s`,
		whereParam,
	)
	var serverRelList []*ServerRel
	err := Db.Raw(sql).Scan(&serverRelList).Error
	utils.CheckError(err)
	var GlobalAccountId = []string{}
	for _, data := range serverRelList {
		GlobalAccountId = append(GlobalAccountId, strconv.Itoa(data.Uid))
	}
	return &GlobalAccountId, err
}

func GetCustomerByUid(UidList []string) map[int]string {
	whereArray := make([]string, 0)
	whereArray = append(whereArray, fmt.Sprintf("uid IN (%s)", strings.Join(UidList, ",")))
	whereParam := strings.Join(whereArray, " and ")
	//sql := fmt.Sprintf(`select s_account from server_rel where %s`, whereParam)
	sql := fmt.Sprintf(`select rel.uid,admin.name from myadmin_user as admin left join server_rel as rel on admin.id = rel.s_account  where %s`,
		whereParam)

	var AdminUserList []*AdminUser
	err := Db.Raw(sql).Scan(&AdminUserList).Error
	utils.CheckError(err)
	var CustomerMap = map[int]string{}
	for _, data := range AdminUserList {
		//uid => customer_name
		CustomerMap[data.Uid] = data.Name
	}
	return CustomerMap
}
