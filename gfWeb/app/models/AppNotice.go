package models

import (
	"encoding/json"
	"fmt"
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
)

type AppNotice struct {
	Id          int    `json:"id" gorm:"id"`
	AppId       int    `json:"app_id" gorm:"app_id"`
	AppName     string `json:"app_name" gorm:"-"`
	Type        int    `json:"type" gorm:"type"`
	Version     string `json:"version"gorm:"version"`
	Notice      string `json:"notice" gorm:"notice"`
	Stats       int    `json:"status" gorm:"stats"`
	Repeated    int    `json:"repeated" gorm:"repeated"`
	CreatedAt   int    `json:"created_at"  gorm:"created_at"`
	CreatedBy   int    `json:"created_by" gorm:"created_by"`
	CreatedName string `json:"created_name" gorm:"-"`
	UpdatedAt   int    `json:"updated_at"  gorm:"updated_at"`
	UpdatedBy   int    `json:"updated_by" gorm:"updated_by"`
	UpdatedName string `json:"updated_name" gorm:"-"`
}

type AppNoticeQueryParam struct {
	BaseQueryParam
	AppId   int    `json:"app_id"`
	Type    int    `json:"type"`
	Version string `json:"version"`
}

func AppNoticeList4Erlang(params *AppNoticeQueryParam) ([]*AppNotice, int) {
	g.Log().Info("getList params: %+v", params)
	data := make([]*AppNotice, 0)
	count := 0

	sql := fmt.Sprintf(`SELECT a.*, p.app_id as app_name FROM app_notice AS a LEFT JOIN platform_client_info AS p ON a.app_id = p.id WHERE p.stats = 1`)
	err := Db.Raw(sql).Scan(&data).Error
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

// AppNoticeList 获取app公告列表
func AppNoticeList(params *AppNoticeQueryParam) ([]*AppNotice, int) {
	g.Log().Info("getList params: %+v", params)
	data := make([]*AppNotice, 0)
	count := 0

	sortOrder := "id"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	//err := Db.Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	sql := fmt.Sprintf(`SELECT a.*, p.app_id as app_name FROM app_notice AS a LEFT JOIN platform_client_info AS p ON a.app_id = p.id`)
	err := Db.Debug().Raw(sql).Scan(&data).Error
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

// CreateAppNotice 创建app公告
func CreateAppNotice(params AppNotice) error {
	PlatformClientInfo, err := GetPlatformClientInfoById(params.AppId)
	if err != nil {
		return err
	}
	createRes := Db.Debug().Save(&params).Error
	if createRes != nil {
		return createRes
	}

	err = AsyncNoticeCenterNode(PlatformClientInfo.AppId, params.Notice, params.Version, params.Stats, params.Repeated)
	if err != nil {
		g.Log().Errorf("notice center node failure: %+v", err)
	}
	return createRes
}

// UpdateAppNotice 编辑app公告
func UpdateAppNotice(params AppNotice) error {
	PlatformClientInfo, err := GetPlatformClientInfoById(params.AppId)
	if err != nil {
		return err
	}
	updateRes := Db.Debug().Save(&params).Error
	if updateRes != nil {
		return updateRes
	}

	err = AsyncNoticeCenterNode(PlatformClientInfo.AppId, params.Notice, params.Version, params.Stats, params.Repeated)
	if err != nil {
		g.Log().Errorf("notice center node failure: %+v", err)
	}

	return updateRes
}

func GetAppNoticeById(id int) (*AppNotice, error) {
	AppNoticeInfo := &AppNotice{}
	AppNoticeInfo.Id = id
	err := Db.Debug().First(&AppNoticeInfo).Error
	if err != nil {
		return AppNoticeInfo, err
	}
	return AppNoticeInfo, nil
}

// DeleteAppNotice 删除app公告
func DeleteAppNotice(params AppNotice) error {
	AppInfo, err := GetAppNoticeById(params.Id)
	if err != nil {
		return err
	}
	PlatformClientInfo, platformClientInfoErr := GetPlatformClientInfoById(AppInfo.AppId)
	if platformClientInfoErr != nil {
		return err
	}

	deleteRes := Db.Debug().Delete(&params).Error
	if deleteRes != nil {
		return deleteRes
	}

	err = AsyncNoticeCenterNode(PlatformClientInfo.AppId, AppInfo.Notice, AppInfo.Version, 2, AppInfo.Repeated)
	if err != nil {
		g.Log().Errorf("notice center node failure: %+v", err)
	}

	return deleteRes
}

func AsyncNoticeCenterNode(appId string, notice string, version string, status int, repeated int) error {
	//pool := grpool.New(1)
	pool := utils.GetAsyncPool()
	err := pool.Add(func() {
		var request struct {
			AppId    string `json:"app_id"`
			Version  string `json:"version"`
			Notice   string `json:"notice"`
			Type     int    `json:"type"`
			Stats    int    `json:"stats"`
			Repeated int    `json:"repeated"`
		}
		request.AppId = appId
		request.Notice = notice
		request.Version = version
		request.Stats = status
		request.Repeated = repeated

		data, err := json.Marshal(request)
		utils.CheckError(err)
		if err != nil {
			g.Log().Errorf("request failure: %+v", err)
		} else {
			url := utils.GetCenterURL() + "/set_app_notice"
			resp, _ := utils.HttpRequest(url, string(data))
			g.Log().Infof("call url: %s response: %+v", url, resp)
		}
	})
	if err != nil {
		g.Log().Errorf("error: %+v", err)
	}
	return err
}
