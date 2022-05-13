package models

import (
	"fmt"
	"gfWeb/library/utils"
	"strings"

	"github.com/gogf/gf/frame/g"
)

type GlobalAccount struct {
	PlatformId       string `gorm:"primary_key" json:"platformId"`
	Account          string `gorm:"primary_key" json:"account"`
	RecentServerList string `json:"recent_server_list"`
	Type             int    `json:"type"`
	ForbidType       int    `json:"forbidType"`
	ForbidTime       int    `json:"forbidTime"`
	Promote          string `json:"promote"`
	RegistrationId   string `json:"registration_id"`
}

func (a *GlobalAccount) TableName() string {
	return "global_account"
}

func GetGlobalAccountByPromoteAndPlatformId(platformId string, promote string) ([]string, error) {
	sql := fmt.Sprintf(
		`select * from global_account where platform_id = '%s' AND  promote LIKE '%%%s%%'`,
		platformId,
		promote,
	)
	g.Log().Info("模糊查询。查看sql语句:%+v", sql)
	var accDataList []*GlobalAccount
	err := DbCenter.Raw(sql).Scan(&accDataList).Error
	utils.CheckError(err)
	var accIdList = []string{}
	for _, data := range accDataList {
		accIdList = append(accIdList, data.Account)
	}
	return accIdList, err
}

func GetGlobalAccount(platformId string, accId string) (*GlobalAccount, error) {
	globalAccount := &GlobalAccount{}
	err := DbCenter.Where(&GlobalAccount{PlatformId: platformId,
		Account: accId,
	}).First(&globalAccount).Error
	if globalAccount.Account != "" {
		utils.CheckError(err)
	}
	return globalAccount, err
}

func GetGlobalAccountByPromote(promote string) ([]string, error) {
	accDataList := make([]*GlobalAccount, 0)
	err := DbCenter.Where(&GlobalAccount{Promote: promote}).Find(&accDataList).Error
	utils.CheckError(err)
	g.Log().Info("查看数据:%v", accDataList)
	//var playerIdStr string
	//for _, e := range accDataList {
	//	playerIdStr += fmt.Sprintf(`"%v", `, e.AccId)
	//}
	//playerIdCondition = playerIdStr[0: len(playerIdStr) - 2]
	accIdList := make([]string, 0)
	for _, data := range accDataList {
		accIdList = append(accIdList, data.Account)
	}
	return accIdList, err
}

func GetGlobalAccountById(params *PlayerQueryParam, SAccountList *[]string) (*[]string, error) {

	whereArray := make([]string, 0)
	whereArray = append(whereArray, fmt.Sprintf(" id IN (%s)", strings.Join(*SAccountList, ",")))
	whereArray = append(whereArray, fmt.Sprintf(" platform_id IN (%s)", params.PlatformId))
	whereParam := strings.Join(whereArray, " and ")
	sql := fmt.Sprintf(
		`select * from global_player where  %s`,
		whereParam,
	)

	var accDataList []*GlobalAccount
	err := DbCenter.Raw(sql).Scan(&accDataList).Error
	utils.CheckError(err)
	var accIdList = []string{}
	for _, data := range accDataList {
		accIdList = append(accIdList, data.Account)
	}
	return &accIdList, err
}
