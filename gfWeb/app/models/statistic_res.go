package models

import (
	"gfWeb/library/utils"

	"github.com/gogf/gf/util/gconv"
)

//"sort"

type StatisticResource struct {
	Id         int    `gorm:"primary_key" json:"id"`
	Url        string `json:"url"`
	AppId      string `json:"app_id"`
	Version    int    `json:"version"`
	AddTime    int    `json:"add_time"`
	Uptime     int    `orm:"uptime" json:"uptime"`
	AddTimeStr string `orm:"-" json:"add_times"`
}

var StatisticRes = &StatisticResource{}

func (a *StatisticResource) TableName() string {
	return "statistic_res"
}

type OptStatisticResReq struct {
	BaseQueryParam
	AppId string `json:"app_id"`
}

type UptatisticResReq struct {
	Id      int    `json:"id"`
	Url     string `json:"url"`
	AppId   string `json:"app_id"`
	Versoin int    `json:"version"`
}

type AddStatisticResReq struct {
	Url     string `json:"url"`
	AppId   string `json:"app_id"`
	Versoin int    `json:"version"`
}

type DelStatisticResReq struct {
	Id int `json:"id"`
}

func (a *StatisticResource) FetchRelUrlVer(appId string) []map[string]string {
	var data = make([]*StatisticResource, 0)
	err := Db.Where(StatisticResource{AppId: appId}).Find(&data).Error
	utils.CheckError(err)

	var ret = make([]map[string]string, 0)
	for _, v := range data {
		ret = append(ret, map[string]string{"version": gconv.String(v.Version), "download": v.Url})
	}
	return ret
}
