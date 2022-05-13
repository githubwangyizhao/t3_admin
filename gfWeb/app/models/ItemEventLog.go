package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	//"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"

	//"reflect"

	//"fmt"
	"gfWeb/library/utils"
	"gfWeb/memdb/uselog"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	//"regexp"
	//"strings"
	//"time"
)

type ItemEventLog struct {
	PlayerId   int     `json:"player_id"`
	PlatformId string  `json:"platform_id"`
	ServerId   string  `json:"server_id"`
	LogId      int     `json:"log_id"`
	Type       int     `json:"type"`
	Players    int     `json:"players"`
	Count      int     `json:"count"`
	Avg        float32 `json:"avg"`
	Time       int     `json:"time"`
}

type RoundBoxReq struct {
	PlatformId string `json:"platform_id"`
	ServerId   string `json:"server_id"`
	PlayerName string `json:"player_name"`
	LogType    int    `json:"log_type"`
	SceneId    int    `json:"scene_id"`
	//
	PlayerId int `json:"player_id"`
}

type ItemEventLogQueryParam struct {
	BaseQueryParam
	PlatformId string `json:"platform_id"`
	ServerId   string `json:"server_id"`
	PlayerName string `json:"player_name"`
	Datetime   int    `json:"datetime"`
	LogType    int    `json:"log_type"` // log id
	Type       int    `json:"type"`
	//
	MonsterId int `json:"monster_id"`

	// not log controller
	StartTime int `json:"start_time"`
	EndTime   int `json:"end_time"`
	Ip        string
	PlayerId  int
}

type MonsterJsonStruct struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type AdminJsonStruct struct {
	ServiceLogStruct ServiceLogStruct
	ItemStruct       ItemStruct
	SceneStruct      SceneStruct
}

// ServiceLogStruct http://后台的host/static/json/ServiceLog.json的结构体
type ServiceLogStruct struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Sign string `json:"sign"`
}

// SceneStruct http://后台的host/static/json/Scene.json的结构体
type SceneStruct struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}

// ItemStruct http://后台的host/static/json/Item.json的结构体
type ItemStruct struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	CanBeTraded int    `json:"can_be_traded"`
	CanBeLogged int    `json:"can_be_logged"`
}

func GetParams(logType string) (string, interface{}) {
	var (
		fileName       = "ServiceLog"
		responseStruct = &AdminJsonStruct{}
	)
	switch logType {
	case "item":
		fileName = "Item"
		return fileName, responseStruct.ItemStruct
	case "scene":
		fileName = "Scene"
		return fileName, responseStruct.SceneStruct
	case "monster":
		fileName = "Monster"
		return fileName, responseStruct.SceneStruct
	default:
		return fileName, responseStruct.ServiceLogStruct
	}
}

// 不用了
func ParseResponse(logType string, result interface{}) (int, int) {
	g.Log().Infof("sdffff: %+v", result)
	switch logType {
	case "item":
		if Res, ok := result.(ItemStruct); ok {
			return Res.Id, Res.CanBeTraded
		}
	default:
		if Res, ok := result.(ServiceLogStruct); ok {
			return Res.Id, 0
		}
	}
	return 0, 0
}

//func GetJson(logType string) (interface{}, interface{}, error) {
func GetJson(logType string) ([]byte, error) {
	fileName, result := GetParams(logType)
	g.Log().Infof("result: %+v", result)
	//results := result
	results := []byte("")

	var host = g.Cfg().GetString("game.gameCenterHost", "http://127.0.0.1:13000")
	// if logType == "monster" {
	// 	host = g.Cfg().GetString("game.gameCenterHostOnlyMonster", "http://127.0.0.1:13000")
	// }

	var url = fmt.Sprintf(`%s/static/json/%s.json`, host, fileName)

	g.Log().Infof("url: %+v", url)

	resp, err := http.Get(url)
	if err != nil {
		utils.CheckError(err)
		//return result, results, err
		return results, err
	}

	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)

	return responseBody, nil
	//err = json.Unmarshal(responseBody, &results)
	//if err != nil {
	//	return result, results, err
	//}
	//return result, results, nil
}

// hour 为指定时间
func InsertGameMonsterLog(optTime time.Time, hour int, genTodayAll int) error {
	platformIdList := make([]string, 0)
	platformList := GetPlatformList()
	for _, e := range platformList {
		platformIdList = append(platformIdList, e.Id)
	}
	gameServerList := GetServerList(platformIdList)
	var monsterArr = make([]*MonsterJsonStruct, 0)

	ResponseByte, itemRrr := GetJson("monster")
	if itemRrr != nil {
		utils.CheckError(itemRrr)
		return itemRrr
	}
	// 获取日志数据
	jsonErr := json.Unmarshal(ResponseByte, &monsterArr)
	if jsonErr != nil {
		utils.CheckError(jsonErr)
		return jsonErr
	}

	theHour := hour

	currentTime := optTime
	// if optTime == nil {
	// 	currentTime = time.Now()
	// }
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	currentDate := startTime.Unix()
	//for _, m := range responseStruct.([]interface{}) {
	//	if err := mapstructure.Decode(m, &respStruct); err != nil {
	//		utils.CheckError(err)
	//	}
	//
	//	Id, _ := ParseResponse("serviceLog", respStruct)

	// 如果是今天 ，只生成过去时间里的数据
	var (
		todayOptTime = time.Now()
		isToday      = false
		maxHour      = todayOptTime.Hour()
		todayTime    = time.Date(todayOptTime.Year(), todayOptTime.Month(), todayOptTime.Day(), 0, 0, 0, 0, todayOptTime.Location())
	)
	// fmt.Println(todayTime.Unix())
	// fmt.Println(startTime.Unix())
	if todayTime.Unix() == startTime.Unix() {
		isToday = true
	}

	for _, m := range monsterArr {
		var params ItemEventLogQueryParam
		for _, object := range gameServerList {
			params.PlatformId = object.PlatformId
			params.ServerId = object.Sid
			params.MonsterId = m.Id
			params.Datetime = int(currentDate)
		}

		if genTodayAll > 0 {
			for i := 1; i <= 24; i++ {
				if isToday && i > maxHour {
					break
				}
				ret := GetItemGameLogList(&params, i-1, LogMonsterType)
				uselog.InsertMonsterLog(params.PlatformId, params.ServerId, gconv.String(params.MonsterId), ret)
			}
		} else {
			ret := GetItemGameLogList(&params, theHour, LogMonsterType)
			uselog.InsertMonsterLog(params.PlatformId, params.ServerId, gconv.String(params.MonsterId), ret)
		}
	}

	return nil
}
func _doInsertEventLog(params *ItemEventLogQueryParam, genTodayAll, maxHour, theHour int, isToday bool) {
	if genTodayAll > 0 {
		for i := 1; i <= 24; i++ {
			if isToday && i > maxHour {
				break
			}
			ret := GetItemGameLogList(params, i-1, LogEventType)
			uselog.InsertEventLog(params.PlatformId, params.ServerId, gconv.String(params.LogType), gconv.String(params.Type), ret)
		}
	} else {
		ret := GetItemGameLogList(params, theHour, LogEventType)
		uselog.InsertEventLog(params.PlatformId, params.ServerId, gconv.String(params.LogType), gconv.String(params.Type), ret)
	}
}

func InsertItemEventLog(optTime time.Time, hour int, genTodayAll int) error {
	// 获取平台，区服列表
	fmt.Println("start InsertItemEventLog ... ")
	platformIdList := make([]string, 0)
	platformList := GetPlatformList()
	for _, e := range platformList {
		platformIdList = append(platformIdList, e.Id)
	}
	gameServerList := GetServerList(platformIdList)

	// 获取item.json里数据，当item.json返回结果中can_be_logged不为0时，就是对应的logid的值
	//itemRespStruct, itemResponseStruct, itemRrr := GetJson("item")
	var (
		itemResponseStruct       = make([]*ItemStruct, 0)
		SceneResponseStruct      = make([]*SceneStruct, 0)
		ServiceLogResponseStruct = make([]*ServiceLogStruct, 0)

		canBeLoggedList = make(map[int]string)
	)

	{
		ResponseByte, itemRrr := GetJson("item")
		if itemRrr != nil {
			utils.CheckError(itemRrr)
			return itemRrr
		}
		// 获取日志数据
		jsonErr := json.Unmarshal(ResponseByte, &itemResponseStruct)
		if jsonErr != nil {
			utils.CheckError(jsonErr)
			return jsonErr
		}

		// 获取scene.json里的数据，当ServiceLog.json返回结果中logid=13时
		SceneRespByte, SceneErr := GetJson("scene")
		if SceneErr != nil {
			utils.CheckError(SceneErr)
			return SceneErr
		}

		jsonErr = json.Unmarshal(SceneRespByte, &SceneResponseStruct)
		if jsonErr != nil {
			utils.CheckError(jsonErr)
			return jsonErr
		}

		//g.Log().Info("itemRespStruct: %+v %+v", itemRespStruct, reflect.TypeOf(itemRespStruct).Name())
		//g.Log().Info("itemResponseStruct: %+v", itemResponseStruct)
		for _, n := range itemResponseStruct {
			if n.CanBeLogged != 0 {
				if _, ok := canBeLoggedList[n.CanBeLogged]; ok {
					canBeLoggedList[n.CanBeLogged] += ", " + strconv.Itoa(n.Id)
				} else {
					canBeLoggedList[n.CanBeLogged] = strconv.Itoa(n.Id)
				}
			}
		}

		// 获取service_player_log.log里的所有logid的值
		//respStruct, responseStruct, err1 := GetJson("serviceLog")
		ServiceLogByte, ServiceErr := GetJson("serviceLog")
		if ServiceErr != nil {
			utils.CheckError(ServiceErr)
			return ServiceErr
		}

		jsonErr = json.Unmarshal(ServiceLogByte, &ServiceLogResponseStruct)
		if jsonErr != nil {
			utils.CheckError(jsonErr)
			return jsonErr
		}

	}

	theHour := hour

	currentTime := optTime
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	currentDate := startTime.Unix()
	// 如果是今天 ，只生成过去时间里的数据
	var (
		todayOptTime = time.Now()
		isToday      = false
		maxHour      = todayOptTime.Hour()
		todayTime    = time.Date(todayOptTime.Year(), todayOptTime.Month(), todayOptTime.Day(), 0, 0, 0, 0, todayOptTime.Location())
	)
	if todayTime.Unix() == startTime.Unix() {
		isToday = true
	}

	for _, m := range ServiceLogResponseStruct {
		var params ItemEventLogQueryParam
		for _, object := range gameServerList {
			params.PlatformId = object.PlatformId
			params.ServerId = object.Sid
			params.LogType = m.Id
			params.Datetime = int(currentDate)

			// 在指定对应shell命令时，除了传入"logid,x"外，还需要传入另一个参数"type,y"(y就是scene.json返回的id)
			if m.Id == 13 || m.Id == 14 {
				for _, o := range SceneResponseStruct {
					params.Type = o.Id
					g.Log().Info(" ---------- params: %+v", params)
					_doInsertEventLog(&params, genTodayAll, maxHour, theHour, isToday)
				}
			} else if m.Id == 15 {
				if _, ok := canBeLoggedList[m.Id]; ok {
					strList := strings.Split(canBeLoggedList[m.Id], `,`)
					// 此时在指定对应shell命令时，除了传入"logid,x"外，还需要传入另一个参数"type,y"(y就是item.json返回的id)
					for _, str := range strList {
						num, _ := strconv.ParseInt(str, 10, 64)
						params.Type = int(num)
						_doInsertEventLog(&params, genTodayAll, maxHour, theHour, isToday)
					}
				} else {
					_doInsertEventLog(&params, genTodayAll, maxHour, theHour, isToday)
				}
			} else {
				_doInsertEventLog(&params, genTodayAll, maxHour, theHour, isToday)
			}
		}
	}

	return nil
}

const (
	LogEventType   = "event"
	LogMonsterType = "monster"
)

func GetItemGameLogList(params *ItemEventLogQueryParam, hour int, logCategory string) map[string]map[string]*uselog.MemItemEventLog {

	gameServer, err := GetGameServerOne(params.PlatformId, params.ServerId)
	utils.CheckError(err)
	if err != nil {
		return nil
	}
	node := gameServer.Node
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil
	}

	var (
		grepParam  = ""
		fileName   = ""
		logDir     = g.Cfg().GetString("game.log_dir")
		scriptPath = ""
		sshKey     = g.Cfg().GetString("game.ssh_key")
		sshPort    = g.Cfg().GetString("game.ssh_port")
		nodeName   = strings.Split(serverNode.Node, "@")[0]
		nodeIp     = strings.Split(serverNode.Node, "@")[1]
	)
	if sshKey != "" {
		sshKey = "-i " + sshKey
	}

	if logCategory == LogMonsterType {
		fileName = "player_fight_log2.log"
		scriptPath = g.Cfg().GetString("game.script_path_monster")
		grepParam = fmt.Sprintf(`m,%d`, params.MonsterId)
		if params.PlayerId > 0 {
			grepParam += fmt.Sprintf(`/p,%d`, params.PlayerId)
		}
	} else if logCategory == LogEventType {
		fileName = "service_player_log.log"
		scriptPath = g.Cfg().GetString("game.script_path")
		grepParam = fmt.Sprintf(`logid,%d`, params.LogType) + "/" + fmt.Sprintf(`type,%d`, params.Type)
		g.Log().Infof("params.PlayerId: %b %d", params.PlayerId != 0, params.PlayerId)
		if params.PlayerId != 0 {
			grepParam += fmt.Sprintf(`/playerid,%d`, params.PlayerId)
		}
		g.Log().Infof("grepParam: %s", grepParam)
	}

	if g.Cfg().GetBool("server.isLocal") {
		nodeIp = "192.168.31.100" //todo 如果是别的服务器 那证书不一样会长时间断开
	}
	if sshPort != "" {
		sshPort = "-p" + sshPort
	}

	t := time.Unix(int64(params.Datetime), 0)
	Datetime := fmt.Sprintf("%d_%d_%d", t.Year(), t.Month(), t.Day())

	if strings.Contains(logDir, "server/log/game") {
		// 内网
		logDir = logDir + Datetime + "/" + fileName
	} else {
		logDir = logDir + fmt.Sprintf(`%s/%s`, nodeName, Datetime) + "/" + fileName
	}

	// logDir := g.Cfg().GetString("game.log_dir", "/opt/t3/trunk/server/log/game/") + Datetime + "/" + fileName
	var ret = map[string]map[string]*uselog.MemItemEventLog{
		Datetime: make(map[string]*uselog.MemItemEventLog),
	}

	var cmd = ""
	var scanHourDuring = gconv.String(hour)
	if len(scanHourDuring) == 1 {
		scanHourDuring = "0" + scanHourDuring
	}

	// ssh执行sh脚本获取指定日志内容
	cmd = fmt.Sprintf("ssh %s %s root@%s '/usr/bin/sh %s %s / ^%s/%s'", sshKey, sshPort, nodeIp, scriptPath, logDir, scanHourDuring, grepParam)

	if cmd == "" {
		g.Log().Errorf("cmd不能为空")
		return nil
	}
	g.Log().Infof("Cmd: %s", cmd)

	// g.Log().Warning("script shell cmd > ", cmd)
	out, err := utils.ExecShell(cmd)
	fmt.Println(err, out)
	if err != nil {
		if len(out) == 0 {
			return nil
		}
		return nil
	}

	lines := strings.Split(out, "\n")
	// fmt.Println("sum len ", len(lines))
	countLines := len(lines) - 2 // 人数
	sumline := lines[len(lines)-2]

	sumlineArr := strings.Split(sumline, "：")

	var (
		total = gconv.Int(strings.ReplaceAll(sumlineArr[1], " ", ""))
		avg   float32
	)

	if countLines == 0 {
		avg = 0
	} else {
		avg = float32(total) / float32(countLines)
	}

	g.Log().Info("sum > ", total, countLines, avg)
	ret[Datetime][gconv.String(scanHourDuring)] = &uselog.MemItemEventLog{
		PlatformId: params.PlatformId,
		ServerId:   params.ServerId,
		Avg:        float32(avg),
		Count:      total,
		LogId:      params.LogType,
		Type:       params.Type,
		Players:    countLines,
		Number:     hour,
	}

	return ret
}
