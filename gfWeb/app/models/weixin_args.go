package models

import (
	"gfWeb/library/utils"
	"strconv"
)

type WeiXinArgs struct {
	IsOpen     int    `json:"isOpen"`
	AppId      string `json:"appId"`
	Path       string `json:"path"`
	ExtraData  string `json:"extraData"`
	EnvVersion string `json:"envVersion"`
}

func GetWeiXinArgs() *WeiXinArgs {
	data := &WeiXinArgs{
		IsOpen:     GetServerDataInt(DbCenter, 15),
		AppId:      GetServerDataStr(DbCenter, 11),
		Path:       GetServerDataStr(DbCenter, 12),
		ExtraData:  GetServerDataStr(DbCenter, 13),
		EnvVersion: GetServerDataStr(DbCenter, 14),
	}
	return data
}

func UpdateWeixinArgs(args *WeiXinArgs) error {
	out, err := utils.CenterNodeTool(
		"weixin",
		"update_navigate_args",
		strconv.Itoa(args.IsOpen),
		args.AppId,
		args.Path,
		args.ExtraData,
		args.EnvVersion,
	)
	utils.CheckError(err, out)
	return err
}
