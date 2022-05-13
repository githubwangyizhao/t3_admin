package models

import (
	"fmt"
	"gfWeb/library/utils"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/jinzhu/gorm"
)

type Prop struct {
	// PropType int `json:"propType,omitempty"` //t3 ver del prop type
	PropId  int `json:"propId,omitempty"`
	PropNum int ` json:"propNum,omitempty"`
}

type PlayerPropLog struct {
	Id          int    `json:"id"`
	PlayerId    int    `json:"playerId"`
	PlayerName  string `json:"playerName" gorm:"-"`
	PropType    int    `json:"propType"`
	PropId      int    `json:"propId"`
	OpType      int    `json:"opType"`
	OpTime      int    `json:"opTime"`
	ChangeValue int    `json:"changeValue"`
	NewValue    int    `json:"newValue"`
}

type PlayerPropLogQueryParam struct {
	BaseQueryParam
	PlatformId string
	ServerId   string `json:"serverId"`
	Ip         string
	PlayerId   int
	PlayerName string
	Datetime   int `json:"datetime"`
	StartTime  int
	EndTime    int
	PropType   int
	PropId     int
	OpType     int
	Type       int //1：获得 2：消耗
}

func GetPlayerPropLogList(params *PlayerPropLogQueryParam) ([]*PlayerPropLog, int64) {
	gameServer, err := GetGameServerOne(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	if err != nil {
		return nil, 0
	}
	node := gameServer.Node
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil, 0
	}
	defer gameDb.Close()
	data := make([]*PlayerPropLog, 0)
	var count int64
	sortOrder := "id"
	//if params.Order == "descending" {
	//	sortOrder = sortOrder + " desc"
	//}
	//if params.Ip != "" {
	//	gameDb = gameDb.Where("ip = ?", params.Ip)
	//}
	//if params.PlayerId != 0 {
	//	gameDb = gameDb.Where("player_id = ?", params.PlayerId)
	//}
	//if params.StartTime != 0 {
	//	gameDb = gameDb.Where("timestamp >= ?", params.StartTime)
	//}
	//if params.EndTime != 0 {
	//	gameDb = gameDb.Where("timestamp <= ?", params.EndTime)
	//}
	f := func(db *gorm.DB) *gorm.DB {
		if params.StartTime > 0 {
			return db.Where("op_time between ? and ?", params.StartTime, params.EndTime)
		}
		return db
	}
	f1 := func(db *gorm.DB) *gorm.DB {
		if params.Type == 1 {
			return db.Where("change_value > 0")
		}
		if params.Type == 2 {
			return db.Where("change_value < 0")
		}
		return db
	}
	f1(f(gameDb.Model(&PlayerPropLog{}).Where(&PlayerPropLog{
		PlayerId: params.PlayerId,
		PropType: params.PropType,
		PropId:   params.PropId,
		OpType:   params.OpType,
	}))).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count)
	for _, e := range data {
		e.PlayerName = GetPlayerName(gameDb, e.PlayerId)
	}
	return data, count
}

//
func GetPlayerPropLogList2(params *PlayerPropLogQueryParam) ([]*PlayerPropLog, int) {
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
	//exec_shell("pwd")
	//exec_shell("ssh -i /root/.ssh/thyz_87 -p22 39.108.98.87 \"pwd\""
	t := time.Unix(int64(params.Datetime), 0)
	logDir := fmt.Sprintf("%d_%d_%d", t.Year(), t.Month(), t.Day())
	g.Log().Debug("logDir:%s", logDir)

	grepParam := ""
	logFile := "player_prop_log.log"
	// 当没有传入params.PlayerId时，默认查询所有玩家的物品使用日志情况
	//grepParam += fmt.Sprintf(" | /usr/bin/grep \\{p,%d\\} ", params.PlayerId)
	if params.PropType > 0 {
		grepParam += fmt.Sprintf(" | /usr/bin/grep \\{pT,%d\\} ", params.PropType)
		if params.PropType == 7 {
			logFile = "player_special_prop_log.log"
		}
	}
	if params.PropId > 0 {
		grepParam += fmt.Sprintf(" | /usr/bin/grep \\{pI,%d\\} ", params.PropId)
	}
	if params.OpType > 0 {
		grepParam += fmt.Sprintf(" | /usr/bin/grep \\{l,%d\\} ", params.OpType)
	}
	sshKey := g.Cfg().GetString("game.ssh_key")
	sshPort := g.Cfg().GetString("game.ssh_port")
	nodeName := strings.Split(serverNode.Node, "@")[0]
	nodeIp := strings.Split(serverNode.Node, "@")[1]
	playerPropLog := g.Cfg().GetString("game.log_dir")
	cmd := fmt.Sprintf("ssh -i %s -p%s root@%s '/usr/bin/cat %s/%s/%s/%s %s'",
		sshKey, sshPort, nodeIp, playerPropLog, nodeName, logDir, logFile, grepParam)
	env := g.Cfg().GetString("game.env")
	if env == "develop" {
		cmd = fmt.Sprintf("ssh -i %s -p%s root@%s '/usr/bin/cat %s/%s/%s %s'",
			sshKey, sshPort, nodeIp, playerPropLog, logDir, logFile, grepParam)
	}
	g.Log().Debugf("cmd: %s", cmd)
	out, err := utils.ExecShell(cmd)
	utils.CheckError(err)
	//g.Log().Debug(out)

	//reg := regexp.MustCompile(`(\d+):(\d+):(\d+)\s+\[{p,(\d+)},{pT,(\d+)},{pI,(\d+)},{l,(\d+)},{c,([-\d]+)},{n,(\d+)}\]`)
	reg := regexp.MustCompile(`(\d+):(\d+):(\d+)\s+\[{p,(\d+)},{pI,(\d+)},{l,(\d+)},{c,([-\d]+)},{n,(\d+)}\]`)
	matchArray := reg.FindAllStringSubmatch(out, -1)
	g.Log().Debugf("cmd: %+v", matchArray)
	data := make([]*PlayerPropLog, 0)
	for _, e := range matchArray {
		g.Log().Debugf("cmd: %+v", e)
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
		//propType, err := strconv.Atoi(e[5])
		//utils.CheckError(err)
		propId, err := strconv.Atoi(e[5])
		utils.CheckError(err)
		logType, err := strconv.Atoi(e[6])
		utils.CheckError(err)
		change, err := strconv.Atoi(e[7])
		utils.CheckError(err)
		new, err := strconv.Atoi(e[8])
		utils.CheckError(err)
		data = append(data, &PlayerPropLog{
			PlayerId: playerId,
			//PropType:    propType,
			PropId:      propId,
			OpType:      logType,
			ChangeValue: change,
			NewValue:    new,
			OpTime:      t,
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
	//exec_shell("ssh -i /root/.ssh/thyz_87 -p22 39.108.98.87 \"cat /data/log/game/qq_s0/2018_9_30/player_prop_log.log | grep 10130\"")
	//cmdStr :=" /bin/bash -c 'ssh -i /root/.ssh/thyz_87 -p22 39.108.98.87 'cat /data/log/game/qq_s0/2018_9_30/player_prop_log.log | grep 10130''"
	//cmdStr := " ssh -i /root/.ssh/thyz_qq -p6888 10.105.54.242 \"cat /data/log/game/s1/2018_10_2/player_prop_log.log |grep 15486\""
	//commandArgs := strings.Split(cmdStr, " ")
	//out, err := utils.Cmd(commandArgs[0], commandArgs[1:])
	//utils.CheckError(err)
	//g.Log().Debug(out)
}
