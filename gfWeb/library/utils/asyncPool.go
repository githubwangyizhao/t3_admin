package utils

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/grpool"
)

var pool *grpool.Pool

func InitAsyncPool() *grpool.Pool {
	pool = grpool.New(100)
	return pool
}

func GetAsyncPool() *grpool.Pool {
	return pool
}

func AsyncNoticeCenterRegion(platform string, channel string, trackerToken string, region string, areaCode string, currency string) error {
	pool := GetAsyncPool()

	err := pool.Add(func() {
		var request struct {
			PlatformId   string `json:"platform_id"`
			TrackerToken string `json:"tracker_token"`
			Channel      string `json:"channel"`
			Region       string `json:"region"`
			AreaCode     string `json:"area_code"`
			Currency     string `json:"currency"`
		}
		request.PlatformId = platform
		request.TrackerToken = trackerToken
		request.Channel = channel
		request.Region = region
		request.AreaCode = areaCode
		request.Currency = currency

		data, err := json.Marshal(request)
		CheckError(err)
		if err != nil {
			g.Log().Errorf("request failure: %+v", err)
		} else {
			url := GetCenterURL() + "/set_platform_tracker_token"
			resp, _ := HttpRequest(url, string(data))
			g.Log().Infof("call url: %s response: %+v", url, resp)
		}
	})
	if err != nil {
		g.Log().Errorf("error: %+v", err)
	}
	return err
}
