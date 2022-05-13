package models

import (
	"gfWeb/library/utils"
)

type GlobalPlayer struct {
	Id         int    `json:"id"`
	Account    string `son:"account"`
	CreateTime int    `json:"createTime"`
	PlatformId string `json:"platformId"`
	ServerId   string `json:"serverId"`
	Channel    string `json:"channel"`
	Nickanme   string `json:"nickname"`
	Level      int    `json:"level"`
	VipLevel   int    `json:"vipLevel"`
	//TotalChargeMoney int    `json:"totalChargeMoney"`
	TotalChargeMoney float32 `json:"totalChargeMoney"`
}

func (a *GlobalPlayer) TableName() string {
	return "global_player"
}

func GetGlobalPlayerOne(playerId int) (*GlobalPlayer, error) {
	data := &GlobalPlayer{
		Id: playerId,
	}
	err := DbCenter.First(&data).Error
	return data, err
}
func GetGlobalPlayerList(platformId string, accId string) ([]*GlobalPlayer, error) {
	data := make([]*GlobalPlayer, 0)
	err := DbCenter.Model(&GlobalPlayer{}).Where(&GlobalPlayer{PlatformId: platformId,
		Account: accId,
	}).Find(&data).Error
	utils.CheckError(err)
	for _, e := range data {
		PlayerDetail, err := GetPlayerDetail(platformId, e.ServerId, e.Id)
		utils.CheckError(err)
		if err == nil {
			e.Level = PlayerDetail.Level
			e.VipLevel = PlayerDetail.VipLevel
			e.TotalChargeMoney = PlayerDetail.TotalChargeMoney
		}
	}
	return data, err
}
