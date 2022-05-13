package models

import (
	"fmt"
	"gfWeb/library/utils"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
)

type PlayerMailLog struct {
	Id         int    `json:"id"`
	PlayerId   int    `json:"playerId"`
	PlayerName string `json:"playerName" gorm:"-"`
	ChangeType string `json:"changeType"`
	LogT       int    `json:"logT"`
	MailIdList string `json:"mailIdList"`
	ItemList   []Prop `json:"itemList"`
	//ItemList 	string  `json:"itemList"`
	ChangeTime int `json:"changeTime"`
}

type PlayerMailLogQueryParam struct {
	BaseQueryParam
	PlatformId string
	ServerId   string `json:"serverId"`
	Ip         string
	PlayerId   int
	PlayerName string
	Datetime   int `json:"datetime"`
	StartTime  int
	EndTime    int
	LogType    int `json:"logType"`
	Type       int //1：增加 2：删除 3：提取附件 `json:"type"`
}

func GetPlayerMailLogList(params *PlayerMailLogQueryParam) ([]*PlayerMailLog, int) {
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
	if params.Type > 0 {
		var changeTypeStr = "addMail"
		if params.Type == 3 {
			changeTypeStr = "item_mail"
		}
		grepParam += fmt.Sprintf(" | /usr/bin/grep \\{cT,%s\\} ", changeTypeStr)
	}
	if params.LogType > 0 {
		grepParam += fmt.Sprintf(" | /usr/bin/grep \\{log,%d\\} ", params.LogType)
	}
	sshKey := g.Cfg().GetString("game.ssh_key")
	sshPort := g.Cfg().GetString("game.ssh_port")
	nodeName := strings.Split(serverNode.Node, "@")[0]
	nodeIp := strings.Split(serverNode.Node, "@")[1]
	cmd := fmt.Sprintf("ssh -i %s -p%s %s ' /usr/bin/cat /data/log/game/%s/%s/player_mail_log.log %s'", sshKey, sshPort, nodeIp, nodeName, logDir, grepParam)
	//cmd := fmt.Sprintf("ssh -i /root/.ssh/thyz_87 -p22 %s ' /usr/bin/cat /data/log/game/%s/%s/player_mail_log.log %s'", nodeIp, nodeName, logDir, grepParam)
	//g.Log().Debug("-----cmd: %+v", cmd)
	out, err := utils.ExecShell(cmd)
	utils.CheckError(err)
	//g.Log().Debug(out)
	reg := regexp.MustCompile(`(\d+):(\d+):(\d+)\s+\[{pid,(\d+)},{cT,(\w+)},{log,(\d+)},{mailL,(.*)},{itemL,(.*)}\]`)
	matchArray := reg.FindAllStringSubmatch(out, -1)
	data := make([]*PlayerMailLog, 0)
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
		changType := e[5]
		logType, err := strconv.Atoi(e[6])
		utils.CheckError(err)
		mailIdList := e[7]
		itemListStr := strings.Replace(e[8], "\"", "", -1)
		var itemList []Prop
		if itemListStr != "[]" {
			itemListStr1 := strings.Replace(strings.Replace(itemListStr, "[{", "", -1), "}]", "", -1)
			line := strings.Split(itemListStr1, "},{")
			for _, PropStr := range line {
				PropList := strings.Split(PropStr, ",")
				// PropType, err := strconv.Atoi(PropList[0])
				// utils.CheckError(err)
				PropId, err := strconv.Atoi(PropList[1])
				utils.CheckError(err)
				PropNum, err := strconv.Atoi(PropList[2])
				utils.CheckError(err)
				itemList = append(itemList, Prop{
					//todo v3 del PropType: PropType,
					PropId:  PropId,
					PropNum: PropNum,
				})
			}
		}
		utils.CheckError(err)
		data = append(data, &PlayerMailLog{
			PlayerId:   playerId,
			ChangeType: changType,
			MailIdList: mailIdList,
			ItemList:   itemList,
			LogT:       logType,
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

func GetPlayerMailLogList2(params *PlayerMailLogQueryParam) ([]*PlayerMailLog, int) {
	out := "00:00:00 [{pid,10657},{cT,addMail},{log,396},{mailL,[{297,704}]},{itemL,\"[{2,10481,100},{2,10471,100},{2,40043,1}]\"}]\r\n" +
		"05:54:38 [{pid,10657},{cT,out_time},{log,379},{mailL,[{254,401}]},{itemL,\"[]\"}]\r\n" +
		"05:54:38 [{pid,10657},{cT,out_time},{log,149},{mailL,[{255,1000}]},{itemL,[]}]\r\n"
	g.Log().Debug("-----cmd:out %+v", out)
	reg := regexp.MustCompile(`(\d+):(\d+):(\d+)\s+\[{pid,(\d+)},{cT,(\w+)},{log,(\d+)},{mailL,(.*)},{itemL,(.*)}\]`)
	matchArray := reg.FindAllStringSubmatch(out, -1)
	data := make([]*PlayerMailLog, 0)
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
		changType := e[5]
		logType, err := strconv.Atoi(e[6])
		utils.CheckError(err)
		mailIdList := e[7]
		itemListStr := strings.Replace(e[8], "\"", "", -1)
		var itemList []Prop
		if itemListStr != "[]" {
			itemListStr1 := strings.Replace(strings.Replace(itemListStr, "[{", "", -1), "}]", "", -1)
			line := strings.Split(itemListStr1, "},{")
			for _, PropStr := range line {
				PropList := strings.Split(PropStr, ",")
				// PropType, err := strconv.Atoi(PropList[0])
				// utils.CheckError(err)
				PropId, err := strconv.Atoi(PropList[1])
				utils.CheckError(err)
				PropNum, err := strconv.Atoi(PropList[2])
				utils.CheckError(err)
				itemList = append(itemList, Prop{
					//todo v3 del PropType: PropType,
					PropId:  PropId,
					PropNum: PropNum,
				})
			}
		}
		utils.CheckError(err)
		data = append(data, &PlayerMailLog{
			PlayerId:   playerId,
			ChangeType: changType,
			MailIdList: mailIdList,
			ItemList:   itemList,
			LogT:       logType,
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
