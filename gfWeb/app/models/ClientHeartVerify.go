package models

import (
	"encoding/json"
	"errors"
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
)

type ClientHeartbeatVerify struct {
	Id       int64  `json:"id" gorm:"id"`
	Platform string `json:"platform_id" gorm:"platform"`
	//StartDate time.Time `json:"start_date" gorm:"start_date"`
	StartDate   int64  `json:"start_date" gorm:"start_date"`
	Interval    int64  `json:"expire" gorm:"interval"`
	Status      int64  `json:"status" gorm:"status"`
	CreatedBy   int    `json:"created_by" gorm:"created_by"`
	CreatedAt   int    `json:"created_at" gorm:"created_at"`
	CreatedName string `json:"createdName" gorm:"-"`
	UpdatedBy   int    `json:"updated_by" gorm:"updated_by"`
	UpdatedAt   int    `json:"updated_at" gorm:"updated_at"`
	UpdatedName string `json:"updatedName" gorm:"-"`
	ServerId    string `json:"server_id" gorm:"server_id"`
}

type ClientHeartVerifyRequest struct {
	ClientHeartbeatVerify
	BaseQueryParam
}

type Req2Erlang4Heartbeat struct {
	PlatformId string `json:"platform_id"`
	ServerId   string `json:"server_id"`
	StartDate  int64  `json:"start_date"`
	Expire     int64  `json:"expire"`
	Status     int64  `json:"status"`
}

func NoticeCenterModifyData(PlatformId string, Server string, StartDate int64, Interval int64, Status int64) error {
	var reqBody Req2Erlang4Heartbeat
	reqBody.PlatformId = PlatformId
	reqBody.ServerId = Server
	reqBody.StartDate = StartDate
	reqBody.Expire = Interval
	reqBody.Status = Status

	CenterUrl := g.Cfg().GetString("game.gs_domain")
	g.Log().Infof("CenterUrl: %+v", CenterUrl)
	url := CenterUrl + "/heartbeat_verify/setting"
	g.Log().Infof("url: %+v", url)

	data, err := json.Marshal(reqBody)
	g.Log().Infof("接口调用结果:%s", data)
	utils.CheckError(err)

	//请求
	Response, _ := utils.HttpRequest(url, string(data))
	g.Log().Infof("接口调用结果:%s", Response)
	return nil
}

func GetClientHeartVerifyById(id int64) (*ClientHeartbeatVerify, error) {
	data := &ClientHeartbeatVerify{
		Id: id,
	}
	err := Db.Where(&data).First(&data).Error
	return data, err
}

func CreOrModifyData(params *ClientHeartVerifyRequest) error {
	if params.Id > 0 {
		_, err := GetClientHeartVerifyById(params.ClientHeartbeatVerify.Id)
		if err != nil {
			return errors.New("数据不存在")
		}
	}

	g.Log().Infof("dddd: %+v", params.ClientHeartbeatVerify)
	err := Db.Save(&params.ClientHeartbeatVerify).Error
	if err != nil {
		utils.CheckError(err)
		return err
	}

	err = NoticeCenterModifyData(params.Platform, params.ServerId, params.StartDate, params.Interval, params.Status)
	if err != nil {
		utils.CheckError(err)
		return err
	}

	return nil
}

func GetClientHeartVerifyList(params *ClientHeartVerifyRequest) ([]*ClientHeartbeatVerify, int64) {
	var data = make([]*ClientHeartbeatVerify, 0)
	var count int64

	query := Db
	if params.Platform != "" {
		query = query.Where("platform = ?", params.Platform)
	}
	if params.ServerId != "" {
		query = query.Where("server_id = ?", params.ServerId)
	}
	if params.Status != 0 {
		query = query.Where("status = ?", params.Status)
	}
	sortOrder := "id desc"
	limit := 1
	if params.Limit > 0 {
		limit = params.Limit
	}

	err := query.Debug().Limit(limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Debug().Error
	if err != nil {
		utils.CheckError(err)
	}

	return data, count
}
