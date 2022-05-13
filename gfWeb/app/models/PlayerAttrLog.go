package models

import (
	"fmt"
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PlayerAttrLog struct {
	Id         int    `json:"id"`
	PlayerId   int    `json:"playerId"`
	PlayerName string `json:"playerName" gorm:"-"`
	FunctionId int    `json:"functionId"`
	Power      int    `json:"power"`
	ChangeTime int    `json:"changeTime"`
	AddPower   int    `json:"addPower"`
}

type PlayerAttrLogQueryParam struct {
	BaseQueryParam
	PlatformId string
	ServerId   string `json:"serverId"`
	Ip         string
	PlayerId   int
	PlayerName string
	Datetime   int `json:"datetime"`
	StartTime  int
	EndTime    int
	FunctionId int `json:"functionId"`
	Type       int // 1：减少 `json:"type"`
}

func GetPlayerAttrLogList(params *PlayerAttrLogQueryParam) ([]*PlayerAttrLog, int) {
	gameServer, err := GetGameServerOne(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	if err != nil {
		return nil, 0
	}
	node := gameServer.Node
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil, 0
	}
	t := time.Unix(int64(params.Datetime), 0)
	logDir := fmt.Sprintf("%d_%d_%d", t.Year(), t.Month(), t.Day())
	g.Log().Debug("logDir:%s", logDir)

	grepParam := ""
	grepParam += fmt.Sprintf(" | /usr/bin/grep \\{pid,%d\\} ", params.PlayerId)
	if params.FunctionId > 0 {
		grepParam += fmt.Sprintf(" | /usr/bin/grep \\{fid,%d\\} ", params.FunctionId)
	}
	if params.Type > 0 {
		grepParam += fmt.Sprintf(" | /usr/bin/grep \\{add,- ")
	}
	sshKey := g.Cfg().GetString("game.ssh_key")
	sshPort := g.Cfg().GetString("game.ssh_port")
	nodeName := strings.Split(serverNode.Node, "@")[0]
	nodeIp := strings.Split(serverNode.Node, "@")[1]
	cmd := fmt.Sprintf("ssh -i %s -p%s %s ' /usr/bin/cat /data/log/game/%s/%s/player_attr_log.log %s'", sshKey, sshPort, nodeIp, nodeName, logDir, grepParam)
	out, err := utils.ExecShell(cmd)
	utils.CheckError(err)
	//g.Log().Debug(out)
	reg := regexp.MustCompile(`(\d+):(\d+):(\d+)\s+\[{pid,(\d+)},{fid,(\d+)},{p,(\d+)},{add,([-\d]+)},{at,(\d+)},{de,(\d+)}\]`)
	matchArray := reg.FindAllStringSubmatch(out, -1)
	//g.Log().Debug("matchArray:%+v", matchArray)
	data := make([]*PlayerAttrLog, 0)
	for _, e := range matchArray {
		h, err := strconv.Atoi(e[1])
		utils.CheckError(err)
		m, err := strconv.Atoi(e[2])
		utils.CheckError(err)
		s, err := strconv.Atoi(e[3])
		utils.CheckError(err)

		time := h*60*60 + m*60 + s
		if params.StartTime > 0 && params.EndTime > 0 {
			if time < params.StartTime || time > params.EndTime {
				continue
			}
		}
		t := params.Datetime + time
		playerId, err := strconv.Atoi(e[4])
		utils.CheckError(err)
		fId, err := strconv.Atoi(e[5])
		utils.CheckError(err)
		power, err := strconv.Atoi(e[6])
		utils.CheckError(err)
		addPower, err := strconv.Atoi(e[7])
		utils.CheckError(err)
		data = append(data, &PlayerAttrLog{
			PlayerId:   playerId,
			FunctionId: fId,
			Power:      power,
			AddPower:   addPower,
			ChangeTime: t,
		})
	}
	len := len(data)
	limit := params.BaseQueryParam.Limit
	start := params.BaseQueryParam.Offset
	if start >= len {
		return nil, len
	}
	if start+limit > len {
		limit = len - start
	}
	g.Log().Debug(len, start, limit)
	return data[start : start+limit], len
}
