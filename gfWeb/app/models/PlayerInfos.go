package models

import (
	"encoding/json"
	"gfWeb/library/utils"
	"time"

	"github.com/gogf/gf/net/ghttp"
)

type PlayerInfos struct {
	Id         int       `json:"id"`
	ServerId   string    `json:"serverId"`
	PlayerId   int       `json:"playerId"`
	PlatformId string    `json:"platformId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	PayTimes   int       `json:"payTimes"`
}

func AddPlayerInfos(r *ghttp.Request) error {
	var playerInfos PlayerInfos
	err := json.Unmarshal(r.GetBody(), &playerInfos)
	utils.CheckError(err)
	playerInfosItem := GetPlayerInfosFirst(playerInfos.PlatformId, playerInfos.ServerId, playerInfos.PlayerId)
	if playerInfosItem.Id != 0 {
		playerInfos.Id = playerInfosItem.Id
	}
	err = Db.Save(&playerInfos).Error
	utils.CheckError(err)
	return err
}

func GetPlayerInfosFirst(PlatformId string, ServerId string, PlayerId int) PlayerInfos {
	playerInfos := PlayerInfos{
		PlatformId: PlatformId,
		ServerId:   ServerId,
		PlayerId:   PlayerId,
	}
	err := Db.Where(&playerInfos).Select("id,platform_id,server_id,player_id,pay_times").First(&playerInfos).Error
	utils.CheckError(err)
	return playerInfos
}

func GetPlayerInfosList(PlatformId string, ServerId string) []map[string]interface{} {
	data := []PlayerInfos{}

	playerInfos := PlayerInfos{
		PlatformId: PlatformId,
		ServerId:   ServerId,
	}
	err := Db.Where(&playerInfos).Select("id,platform_id,server_id,player_id,pay_times").Find(&data).Error
	utils.CheckError(err)
	dataList := []map[string]interface{}{}
	for _, v := range data {
		dataList = append(dataList, map[string]interface{}{"player_id": v.PlayerId, "pay_times": v.PayTimes})
	}
	utils.CheckError(err)
	return dataList
}
