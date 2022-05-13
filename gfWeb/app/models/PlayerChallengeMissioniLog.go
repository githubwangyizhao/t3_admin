package models

import (
	"fmt"
	"gfWeb/library/utils"
	"strconv"
)

type PlayerChallengeMissionLog struct {
	PlayerId    int    `json:"playerId"`
	PlayerName  string `json:"playerName" gorm:"-"`
	MissionType int    `json:"missionType"`
	MissionId   int    `json:"missionId"`
	Result      int    `json:"result"`
	ItemList    []Prop `json:"itemList"`
	Time        int    `json:"time"`
}

type PlayerChallengeMissionLogQueryParam struct {
	BaseQueryParam
	PlatformId  string
	ServerId    string `json:"serverId"`
	Ip          string
	PlayerId    int
	PlayerName  string
	Datetime    int `json:"datetime"`
	MissionType int `json:"missionType"`
	Result      int `json:"result"`
}

//type PlayerChallengeMissionLog struct {
//	Id          int    `json:"id"`
//	PlayerId    int    `json:"playerId"`
//	PlayerName  string `json:"playerName" gorm:"-"`
//	MissionType int    `json:"missionType"`
//	MissionId   int    `json:"missionId"`
//	Result      int    `json:"result"`
//	Time        int    `json:"time"`
//	UsedTime    int    `json:"usedTime"`
//}

//type PlayerChallengeMissionLogQueryParam struct {
//	BaseQueryParam
//	PlatformId  string
//	ServerId        string `json:"serverId"`
//	Ip          string
//	PlayerId    int
//	PlayerName  string
//	StartTime   int
//	EndTime     int
//	MissionType int
//}
//
//func GetPlayerChallengeMissionLogList(params *PlayerChallengeMissionLogQueryParam) ([]*PlayerChallengeMissionLog, int64) {
//	gameServer, err := GetGameServerOne(params.PlatformId, params.ServerId)
//	utils.CheckError(err)
//	if err != nil {
//		return nil, 0
//	}
//	node := gameServer.Node
//	gameDb, err := GetGameDbByNode(node)
//	utils.CheckError(err)
//	if err != nil {
//		return nil, 0
//	}
//	defer gameDb.Close()
//	data := make([]*PlayerChallengeMissionLog, 0)
//	var count int64
//	sortOrder := "id"
//	f := func(db *gorm.DB) *gorm.DB {
//		if params.StartTime > 0 {
//			return db.Where("time between ? and ?", params.StartTime, params.EndTime)
//		}
//		return db
//	}
//	f(gameDb.Model(&PlayerChallengeMissionLog{}).Where(&PlayerChallengeMissionLog{
//		PlayerId:    params.PlayerId,
//		MissionType: params.MissionType,
//	})).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count)
//	for _, e := range data {
//		e.PlayerName = GetPlayerName(gameDb, e.PlayerId)
//	}
//	return data, count
//}

// 获得玩家副本日志内容列表
func GetPlayerMissionLogList(params *PlayerChallengeMissionLogQueryParam) ([]*PlayerChallengeMissionLog, int) {
	grepParam := ""
	//grepParam += fmt.Sprintf(" | /usr/bin/grep \\{pid,%d\\} ", params.PlayerId)
	if params.PlayerId > 0 {
		grepParam += fmt.Sprintf(" | /usr/bin/grep \\{pid,%d\\} ", params.PlayerId)
	}
	if params.MissionType > 0 {
		grepParam += fmt.Sprintf(" | /usr/bin/grep \\{mT,%d\\} ", params.MissionType)
	}
	if params.Result > 0 {
		grepParam += fmt.Sprintf(" | /usr/bin/grep \\{r,%d\\} ", params.Result)
	}
	fileName := "player_challenge_mission_log.log"
	regStr := `(\d+):(\d+):(\d+)\s+\[{pid,(\d+)},{mT,(\d+)},{mI,(\d+)},{r,(\d+)},{award,(.*)}\]`
	err, matchArray := GameLogFileRead(params.PlatformId, params.ServerId, fileName, grepParam, regStr, params.Datetime)
	utils.CheckError(err)
	gameDb, err := GetGameDbByPlatformIdAndSid(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	data := make([]*PlayerChallengeMissionLog, 0)
	if err != nil {
		return data, 0
	}
	defer gameDb.Close()
	//reg := regexp.MustCompile(`(\d+):(\d+):(\d+)\s+\[{pid,(\d+)},{mT,(\d+)},{mI,(\d+)},{r,(\d+)},{award,(.*)}\]`)
	//matchArray := reg.FindAllStringSubmatch(out, -1)
	start := params.BaseQueryParam.Offset
	limit := params.BaseQueryParam.Limit + start
	dataLen := 0
	for _, e := range matchArray {
		if len(e) == 0 {
			continue
		}
		dataLen++
		if start > dataLen || dataLen >= limit {
			continue
		}
		h, err := strconv.Atoi(e[1])
		utils.CheckError(err)
		m, err := strconv.Atoi(e[2])
		utils.CheckError(err)
		s, err := strconv.Atoi(e[3])
		utils.CheckError(err)
		recodeChangeTime := h*60*60 + m*60 + s
		t := params.Datetime + recodeChangeTime
		playerId, err := strconv.Atoi(e[4])
		utils.CheckError(err)
		missionType, err := strconv.Atoi(e[5])
		utils.CheckError(err)
		missionId, err := strconv.Atoi(e[6])
		utils.CheckError(err)
		result, err := strconv.Atoi(e[7])
		utils.CheckError(err)
		itemList := GameLogFileItemList(e[8])
		//itemListStr := strings.Replace(e[8], "\"", "", -1)
		//var itemList []Prop
		//if itemListStr != "[]" {
		//	itemListStr1 := strings.Replace(strings.Replace(itemListStr, "[{", "", -1), "}]", "", -1)
		//	line := strings.Split(itemListStr1, "},{")
		//	for _, PropStr := range line {
		//		PropList := strings.Split(PropStr, ",")
		//		PropType, err := strconv.Atoi(PropList[0])
		//		utils.CheckError(err)
		//		PropId, err := strconv.Atoi(PropList[1])
		//		utils.CheckError(err)
		//		PropNum, err := strconv.Atoi(PropList[2])
		//		utils.CheckError(err)
		//		itemList = append(itemList, Prop{
		//			PropType: PropType,
		//			PropId:   PropId,
		//			PropNum:  PropNum,
		//		})
		//	}
		//}
		playerName := params.PlayerName
		if params.PlayerId != playerId {
			player, err := GetPlayerByDb(gameDb, playerId)
			utils.CheckError(err)
			if err != nil {
				continue
			}
			playerName = player.ServerId + "." + player.Nickname
		}
		data = append(data, &PlayerChallengeMissionLog{
			PlayerId:    playerId,
			PlayerName:  playerName,
			MissionType: missionType,
			MissionId:   missionId,
			Result:      result,
			ItemList:    itemList,
			Time:        t,
		})
	}
	//dataLen := len(data)
	//if start >= dataLen {
	//	return nil, dataLen
	//}
	//if start+limit > dataLen {
	//	limit = dataLen - start
	//}
	//g.Log().Debug(fileName, dataLen, start, limit)
	//return data[start : start+limit], dataLen
	return data, dataLen
}
