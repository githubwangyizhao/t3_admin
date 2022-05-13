package models

import (
	"encoding/json"
	"gfWeb/library/utils"
	"time"
)

type ToolTest struct {
	OptTime     int `json:"opt_time"`      //某一天的时间截
	Hour        int `json:"hour"`          // 统计前一小时数据，如果写10，那么会统计09-10点间的日志
	GenTodayAll int `json:"gen_today_all"` // 1生成今天所有的
}
type PlusDailyStatistics struct {
	StartTime int `json:"start_time"`
	EndTime   int `json:"end_time"`
}
type ClientVersions struct {
	PlatformId         string `json:"platformId"`
	PlatformName       string `json:"platformName"`
	AndroidDownloadUrl string `json:"androidDownloadUrl"`
	IosDownloadUrl     string `json:"iosDownloadUrl"`
	FirstVersions      int    `json:"firstVersions"`
	Versions           int    `json:"versions"`
	IsCloseCharge      int    `json:"isCloseCharge"`
	ClientVersion      string `json:"clientVersion"`
	Ip                 string
	CreatedAt          time.Time
	CreatedBy          int
	UpdatedAt          time.Time
	UpdatedBy          int
}

// 获取单个客户端版本数据
func GetVersionOne(platformId string) (*ClientVersions, error) {
	data := &ClientVersions{
		PlatformId: platformId,
	}
	err := Db.Where(&data).First(&data).Error
	return data, err
}

// 获取客户端版本列表
func GetVersionList() ([]*ClientVersions, int64) {
	data := make([]*ClientVersions, 0)
	var count int64
	err := Db.Find(&data).Count(count).Error
	utils.CheckError(err)
	return data, count
}

// 设置客户当前版本
func SetVersion(params *ClientVersions, ip string, userId int) error {

	_, err := GetVersionOne(params.PlatformId)

	if err == nil {
		err = Db.Exec("update client_versions set platform_name = ? , android_download_url = ? , ios_download_url = ? , first_versions= ? , versions = ? , ip = ? , updated_by = ? , is_close_charge = ? , client_version = ? where platform_id = ? ; ",
			params.PlatformName,
			params.AndroidDownloadUrl,
			params.IosDownloadUrl,
			params.FirstVersions,
			params.Versions,
			ip,
			userId,
			params.IsCloseCharge,
			params.ClientVersion,
			params.PlatformId).Error
	} else {
		params.CreatedBy = userId
		params.Ip = ip
		err = Db.Save(&params).Error
	}
	utils.CheckError(err)
	url := utils.GetCenterURL() + "/update_version"
	//url := "http://192.168.31.153:6666" + "/update_version"
	var request struct {
		Platform           string `json:"platform"`
		PlatformName       string `json:"platformName"`
		AndroidDownloadUrl string `json:"androidDownloadUrl"`
		IosDownloadUrl     string `json:"iosDownloadUrl"`
		Versions           int    `json:"versions"`
		FirstVersions      int    `json:"firstVersions"`
		IsCloseCharge      int    `json:"isCloseCharge"`
		ClientVersion      string `json:"clientVersion"`
	}

	request.Platform = params.PlatformId
	request.PlatformName = params.PlatformName
	request.AndroidDownloadUrl = params.AndroidDownloadUrl
	request.IosDownloadUrl = params.IosDownloadUrl
	request.FirstVersions = params.FirstVersions
	request.Versions = params.Versions
	request.IsCloseCharge = params.IsCloseCharge
	request.ClientVersion = params.ClientVersion
	data2, err := json.Marshal(request)
	utils.CheckError(err)

	_, err = utils.HttpRequest(url, string(data2))
	return err
}
