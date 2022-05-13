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

type PlayerActivityAwardLog struct {
	ActivityId int    `json:"activityId"`
	Rank       int    `json:"rank"`
	PlayerId   int    `json:"playerId"`
	PlayerName string `json:"playerName" gorm:"-"`
	Value      int    `json:"value"`
	ChangeTime int    `json:"changeTime"`
	AwardId    int    `json:"awardId"`
}

type ActivityAwardLogQueryParam struct {
	BaseQueryParam
	PlatformId string
	ServerId   string `json:"serverId"`
	ServerType int
	PlayerId   int    `json:"playerId"`
	PlayerName string `json:"playerName" gorm:"-"`
	Datetime   int    `json:"datetime"`
	ActivityId int    `json:"activityId"`
}

// 查询活动奖励日志
func GetActivityAwardLogList(params *ActivityAwardLogQueryParam) ([]*PlayerActivityAwardLog, int) {
	var node string
	var dirType string
	if params.ServerType == 7 {
		WarNodeList := GetPlatformIdServerNodeByType(params.PlatformId, params.ServerType)
		if len(WarNodeList) != 1 {
			return nil, 0
		}
		ServerNode := WarNodeList[0]
		node = ServerNode.Node
		dirType = "zone"
	} else {
		gameServer, err := GetGameServerOne(params.PlatformId, params.ServerId)
		utils.CheckError(err)
		if err != nil {
			return nil, 0
		}
		node = gameServer.Node
		dirType = "game"
	}
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil, 0
	}
	t := time.Unix(int64(params.Datetime), 0)
	logDir := fmt.Sprintf("%d_%d_%d", t.Year(), t.Month(), t.Day())
	g.Log().Debug("logDir:%s", logDir)

	grepParam := ""
	if params.ActivityId > 0 {
		grepParam += fmt.Sprintf(" | /usr/bin/grep \\{aId,%d\\} ", params.ActivityId)
	}
	sshKey := g.Cfg().GetString("game.ssh_key")
	sshPort := g.Cfg().GetString("game.ssh_port")
	nodeName := strings.Split(serverNode.Node, "@")[0]
	nodeIp := strings.Split(serverNode.Node, "@")[1]
	cmd := fmt.Sprintf("ssh -i %s -p%s %s ' /usr/bin/cat /data/log/%s/%s/%s/activity_award_log.log %s'", sshKey, sshPort, nodeIp, dirType, nodeName, logDir, grepParam)
	out, err := utils.ExecShell(cmd)
	utils.CheckError(err)
	reg := regexp.MustCompile(`(\d+):(\d+):(\d+)\s+\[{aId,(\d+)},{logT,(\d+)},{data,(.*)}\]`)
	matchArray := reg.FindAllStringSubmatch(out, -1)
	//g.Log().Debug("matchArray:%+v", matchArray)
	data := make([]*PlayerActivityAwardLog, 0)
	for _, e := range matchArray {
		playerStr := strings.Replace(e[6], "\"", "", -1)
		if playerStr == "[]" {
			continue
		}
		h, err := strconv.Atoi(e[1])
		utils.CheckError(err)
		m, err := strconv.Atoi(e[2])
		utils.CheckError(err)
		s, err := strconv.Atoi(e[3])
		utils.CheckError(err)
		time := h*60*60 + m*60 + s
		t := params.Datetime + time
		activityId, err := strconv.Atoi(e[4])
		utils.CheckError(err)
		//logType, err := strconv.Atoi(e[5])
		//utils.CheckError(err)
		//playerStr := e[6]
		playerStr1 := strings.Replace(strings.Replace(playerStr, "[{", "", -1), "}]", "", -1)
		line := strings.Split(playerStr1, "},{")
		for _, PlayerAwardStr := range line {
			AwardInfo := strings.Split(PlayerAwardStr, ",")
			PlayerId, err := strconv.Atoi(AwardInfo[0])
			utils.CheckError(err)
			if params.PlayerId > 0 {
				if params.PlayerId != PlayerId {
					continue
				}
			}
			Rank, err := strconv.Atoi(AwardInfo[1])
			utils.CheckError(err)
			Value, err := strconv.Atoi(AwardInfo[2])
			utils.CheckError(err)
			AwardId, err := strconv.Atoi(AwardInfo[3])
			utils.CheckError(err)
			data = append(data, &PlayerActivityAwardLog{
				PlayerId:   PlayerId,
				ActivityId: activityId,
				Rank:       Rank,
				Value:      Value,
				AwardId:    AwardId,
				ChangeTime: t,
			})
		}
	}
	len := len(data)
	if params.Order == "ascending" {
		for i, j := 0, len-1; i < j; i, j = i+1, j-1 {
			data[i], data[j] = data[j], data[i]
		}
	}
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
func GetActivityAwardLogList2(params *ActivityAwardLogQueryParam) ([]*PlayerActivityAwardLog, int) {
	out := "23:59:58 [{aId,6027},{logT,387},{data,[{3743857,50,500,34428},{3978435,49,500,34428},{3870938,48,510,34428},{3865647,47,510,34428},{3870592,46,510,34428},{3577830,45,520,34428},{2666245,44,520,34428},{3707223,43,520,34428},{3947470,42,530,34428},{3697628,41,540,34428},{3063224,40,540,34427},{915023,39,540,34427},{2039233,38,540,34427},{536240,37,550,34427},{1622174,36,550,34427},{2719887,35,550,34427},{1916846,34,550,34427},{3254512,33,570,34427},{2266235,32,590,34427},{3134384,31,600,34427},{3701200,30,620,34426},{3817337,29,620,34426},{1103881,28,630,34426},{2629226,27,640,34426},{2675485,26,640,34426},{508661,25,700,34426},{1633935,24,700,34426},{1740074,23,700,34426},{2795540,22,700,34426},{1570591,21,700,34426},{88621,20,720,34425},{1558451,19,800,34425},{3751816,18,820,34425},{3636641,17,830,34425},{1741813,16,840,34425},{3248306,15,840,34425},{2341604,14,850,34425},{3505870,13,1160,34425},{3931708,12,1170,34425},{4115442,11,1180,34425},{1640192,10,1560,34424},{3761120,9,1570,34424},{3992045,8,1600,34424},{2300280,7,1790,34424},{3974319,6,1800,34424},{3869684,5,2130,34424},{3837696,4,2130,34424},{2418110,3,2280,34423},{3942070,2,3040,34422},{3546054,1,3470,34421}]}]\r\n" +
		"23:59:58 [{aId,5012},{logT,210},{data,[{1475355,50,7500,78348},{2273540,49,7500,78348},{2483581,48,7500,78348},{2546041,47,7800,78348},{2303347,46,7800,78348},{1724732,45,7800,78348},{2424512,44,7800,78348},{2139903,43,8100,78348},{2394942,42,8100,78348},{2483343,41,8100,78348},{2279726,40,8100,78348},{2292349,39,8400,78348},{2541895,38,8400,78348},{2167718,37,9000,78348},{1278604,36,9000,78348},{1939942,35,9000,78348},{1580445,34,9300,78348},{2220900,33,9600,78348},{2158161,32,9900,78348},{1006003,31,9900,78348},{2297066,30,10200,78347},{877164,29,10200,78347},{2132861,28,10200,78347},{2628607,27,10200,78347},{1667490,26,10200,78347},{2015990,25,11100,78347},{1534281,24,11400,78347},{2138433,23,11700,78347},{2275888,22,12000,78347},{2107747,21,12000,78347},{2269205,20,12000,78346},{2636512,19,13500,78346},{1532817,18,13500,78346},{2197678,17,13500,78346},{1572592,16,14100,78346},{2141033,15,14400,78346},{1564303,14,15000,78346},{1597123,13,15900,78346},{2214501,12,18000,78346},{1606343,11,18000,78346},{1981626,10,18600,78345},{800860,9,18600,78345},{2082904,8,18600,78345},{1793485,7,18900,78345},{1686960,6,20700,78345},{2255231,5,21000,78344},{1667327,4,24000,78344},{1589345,3,27000,78343},{2695665,2,28200,78342},{2089915,1,29100,78341}]}]"
	reg := regexp.MustCompile(`(\d+):(\d+):(\d+)\s+\[{aId,(\d+)},{logT,(\d+)},{data,(.*)}\]`)
	matchArray := reg.FindAllStringSubmatch(out, -1)
	//g.Log().Debug("matchArray:%+v", matchArray)
	data := make([]*PlayerActivityAwardLog, 0)
	for _, e := range matchArray {
		playerStr := strings.Replace(e[6], "\"", "", -1)
		if playerStr == "[]" {
			continue
		}
		h, err := strconv.Atoi(e[1])
		utils.CheckError(err)
		m, err := strconv.Atoi(e[2])
		utils.CheckError(err)
		s, err := strconv.Atoi(e[3])
		utils.CheckError(err)
		time := h*60*60 + m*60 + s
		t := params.Datetime + time
		activityId, err := strconv.Atoi(e[4])
		utils.CheckError(err)
		//logType, err := strconv.Atoi(e[5])
		//utils.CheckError(err)
		//playerStr := e[6]
		playerStr1 := strings.Replace(strings.Replace(playerStr, "[{", "", -1), "}]", "", -1)
		line := strings.Split(playerStr1, "},{")
		for _, PlayerAwardStr := range line {
			AwardInfo := strings.Split(PlayerAwardStr, ",")
			PlayerId, err := strconv.Atoi(AwardInfo[0])
			utils.CheckError(err)
			if params.PlayerId > 0 {
				if params.PlayerId != PlayerId {
					continue
				}
			}
			Rank, err := strconv.Atoi(AwardInfo[1])
			utils.CheckError(err)
			Value, err := strconv.Atoi(AwardInfo[2])
			utils.CheckError(err)
			AwardId, err := strconv.Atoi(AwardInfo[3])
			utils.CheckError(err)
			data = append(data, &PlayerActivityAwardLog{
				PlayerId:   PlayerId,
				ActivityId: activityId,
				Rank:       Rank,
				Value:      Value,
				AwardId:    AwardId,
				ChangeTime: t,
			})
		}
	}
	len := len(data)
	if params.Order == "ascending" {
		for i, j := 0, len-1; i < j; i, j = i+1, j-1 {
			data[i], data[j] = data[j], data[i]
		}
	}
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
