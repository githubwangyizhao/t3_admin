package models

import (
	"encoding/json"
	"errors"
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

type Adjust struct {
	Id          int    `json:"id" gorm:"id"`
	PlatformId  string `json:"platformId" gorm:"platform_id"`
	Server      string `json:"server" gorm:"server"`
	Type        int64  `json:"type" gorm:"type"`
	RefId       int64  `json:"refId" gorm:"ref_id"`
	Value       int64  `json:"value" gorm:"value"`
	Status      int64  `json:"status" gorm:"status"`
	CreatedAt   int    `json:"createdAt" gorm:"created_at"`
	CreatedBy   int    `json:"createdBy" gorm:"created_by"`
	CreatedName string `json:"createdName" gorm:"-"`
	UpdatedAt   int    `json:"updatedAt" gorm:"updated_at"`
	UpdatedBy   int    `json:"updatedBy" gorm:"updated_by"`
	UpdatedName string `json:"updatedName" gorm:"-"`
}

type AdjustRequest struct {
	Adjust
	BaseQueryParam
}

type Req2Erlang4Adjust struct {
	ConfigType int64 `json:"config_type"`
	ConfigId   int64 `json:"config_id"`
	Value      int64 `json:"value"`
}

func GetAdjustByPlatform(Platform string, Server string) ([]*Adjust, int64) {
	var data = make([]*Adjust, 0)
	var count int64

	query := Db
	if Platform != "" {
		query = query.Where("platform_id = ?", Platform)
	}
	if Server != "" {
		query = query.Where("server = ?", Server)
	}
	//if Type != 0 {
	//	query = query.Where("type = ?", Type)
	//}
	sortOrder := "id desc"

	err := query.Where("status = ?", 1).Order(sortOrder).Find(&data).Offset(0).Count(&count).Debug().Error
	if err != nil {
		utils.CheckError(err)
	}
	return data, count
}

func GetAdjustDataById(id int) (*Adjust, error) {
	data := &Adjust{
		Id: id,
	}
	err := Db.Where(&data).First(&data).Error
	return data, err
}

func UpdateData2Game(PlatformId string, Server string, Type int64, RefId int64, Value int64, Status int64) error {
	if Server == "" {
		gameServerList := GetPlatformIdAllGameServerList(PlatformId)
		warNodeList := GetPlatformIdServerNodeByType(PlatformId, 7)
		if len(warNodeList) != 1 {
			return nil
		}
		WarServerNode := warNodeList[0]
		var reqBody Req2Erlang4Adjust
		reqBody.ConfigType = Type
		reqBody.ConfigId = RefId

		var (
			url      = "http://" + WarServerNode.Ip + ":" + gconv.String(WarServerNode.WebPort)
			Continue = false
		)
		if Status == 0 {
			url = url + "/delete_game_server_config"
			Continue = true
		} else if Status == 1 {
			url = url + "/set_game_server_config"
			reqBody.Value = Value
			Continue = true
		}

		if Continue == true {
			data, err := json.Marshal(reqBody)
			g.Log().Infof("接口调用结果:%s", data)
			utils.CheckError(err)

			//请求
			Response, _ := utils.HttpRequest(url, string(data))

			//if ReqErr != nil {
			//	utils.CheckError(ReqErr)
			//}
			g.Log().Infof("接口调用结果:%s", Response)
		} else {
			g.Log().Errorf("非法Status:%d", Status)
		}
		for _, v := range gameServerList {
			serverNode, err := GetServerNode(v.Node)
			utils.CheckError(err)
			if err != nil {
				return err
			}
			var reqBody Req2Erlang4Adjust
			reqBody.ConfigType = Type
			reqBody.ConfigId = RefId

			var (
				url      = "http://" + serverNode.Ip + ":" + gconv.String(serverNode.WebPort)
				Continue = false
			)
			if Status == 0 {
				url = url + "/delete_game_server_config"
				Continue = true
			} else if Status == 1 {
				url = url + "/set_game_server_config"
				reqBody.Value = Value
				Continue = true
			}

			if Continue == true {
				data, err := json.Marshal(reqBody)
				g.Log().Infof("接口调用结果:%s", data)
				utils.CheckError(err)

				//请求
				Response, _ := utils.HttpRequest(url, string(data))

				//if ReqErr != nil {
				//	utils.CheckError(ReqErr)
				//}
				g.Log().Infof("接口调用结果:%s", Response)
			} else {
				g.Log().Errorf("非法Status:%d", Status)
			}
		}
	} else {
		gameServer, err := GetGameServerOne(PlatformId, Server)
		utils.CheckError(err)
		if err != nil {
			return err
		}

		serverNode, err := GetServerNode(gameServer.Node)
		utils.CheckError(err)
		if err != nil {
			return err
		}

		var reqBody Req2Erlang4Adjust
		reqBody.ConfigType = Type
		reqBody.ConfigId = RefId

		var (
			url      = "http://" + serverNode.Ip + ":" + gconv.String(serverNode.WebPort)
			Continue = false
		)
		if Status == 0 {
			url = url + "/delete_game_server_config"
			Continue = true
		} else if Status == 1 {
			url = url + "/set_game_server_config"
			reqBody.Value = Value
			Continue = true
		}

		if Continue == true {
			data, err := json.Marshal(reqBody)
			g.Log().Infof("接口调用结果:%s", data)
			utils.CheckError(err)

			//请求
			Response, _ := utils.HttpRequest(url, string(data))

			//if ReqErr != nil {
			//	utils.CheckError(ReqErr)
			//}
			g.Log().Infof("接口调用结果:%s", Response)
		} else {
			g.Log().Errorf("非法Status:%d", Status)
		}
	}

	return nil
}

func EditAdjustData(adjust *Adjust) error {
	if adjust.Id <= 0 {
		return errors.New("id输入错误")
	}

	_, err := GetAdjustDataById(adjust.Id)
	if err != nil {
		return errors.New("数据不存在")
	}

	err = Db.Save(adjust).Error

	UpdateData2Game(adjust.PlatformId, adjust.Server, adjust.Type, adjust.RefId, adjust.Value, adjust.Status)

	return err
}

func CreateAdjustData(adjust *Adjust) error {
	err := Db.Save(adjust).Error
	return err
}

func GetAdjustDataList(params *AdjustRequest) ([]*Adjust, int64) {
	var data = make([]*Adjust, 0)
	var count int64

	query := Db
	if params.PlatformId != "" {
		query = query.Where("platform_id = ?", params.PlatformId)
	}
	if params.Server != "" {
		query = query.Where("server = ?", params.Server)
	}
	if params.RefId != 0 {
		query = query.Where("ref_id = ?", params.RefId)
	}
	if params.Status != 0 {
		query = query.Where("status = ?", params.Status)
	}
	if params.Type != 0 {
		query = query.Where("type = ?", params.Type)
	}
	sortOrder := "id"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}

	err := query.Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Debug().Error
	if err != nil {
		utils.CheckError(err)
		return data, count
	}
	for _, e := range data {
		e.CreatedName = GetUserName(e.CreatedBy)
		e.UpdatedName = GetUserName(e.UpdatedBy)
	}
	return data, count
}
