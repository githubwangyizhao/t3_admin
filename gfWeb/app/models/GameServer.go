package models

import (
	"database/sql"
	"fmt"
	"gfWeb/library/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/text/gstr"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type GameServerQueryParam struct {
	BaseQueryParam
	PlatformId string `json:"platformId"`
	ServerId   string `json:"serverId"`
	Node       string `json:"node"`
}

type GameServer struct {
	PlatformId      string `gorm:"primary_key" json:"platformId"`
	Sid             string `gorm:"primary_key" json:"serverId"`
	Desc            string `json:"desc"`
	Node            string `json:"node"`
	State           int    `gorm:"-" json:"state"`
	IsShow          int    `json:"isShow"`
	Ip              string `json:"ip" gorm:"-"`
	Database        string `json:"database" gorm:"-"`
	OpenTime        int    `gorm:"-" json:"openTime"`
	ZoneNode        string `gorm:"-" json:"zoneNode"`
	IsAdd           int    `gorm:"-" json:"isAdd"`
	DbVersion       int    `json:"dbVersion" gorm:"-"`
	RunState        int    `json:"runState" gorm:"-"`
	StartTime       int    `json:"startTime" gorm:"-"`
	OnlineCount     int    `gorm:"-" json:"onlineCount"`
	CreateRoleCount int    `gorm:"-" json:"createRoleCount"`
}

func (t *GameServer) TableName() string {
	return "c_game_server"
}

//获取所有数据
func GetAllGameServerDirty() ([]*GameServer, int64) {
	data := make([]*GameServer, 0)
	var count int64
	err := DbCenter.Model(&GameServer{}).Find(&data).Count(&count).Error
	utils.CheckError(err)
	return data, count
}

// 获得当前平台所有数据
func GetAllGameServerDirtyByPlatformId(platformId string) ([]*GameServer, int64) {
	data := make([]*GameServer, 0)
	var count int64
	err := DbCenter.Model(&GameServer{}).Where(&GameServer{PlatformId: platformId}).Find(&data).Count(&count).Error
	utils.CheckError(err)
	return data, count
}

//获取所有数据
//func GetAllGameServer() ([]*GameServer, int64) {
//	var params GameServerQueryParam
//	params.Limit = -1
//	//获取数据列表和总数
//	data, total := GetGameServerList(&params)
//	return data, total
//}

type gameServerSlice []*GameServer

func (s gameServerSlice) Len() int      { return len(s) }
func (s gameServerSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s gameServerSlice) Less(i, j int) bool {
	iId, err := strconv.Atoi(SubString(s[i].Sid, 1, len(s[i].Sid)-1))
	utils.CheckError(err)
	jId, err := strconv.Atoi(SubString(s[j].Sid, 1, len(s[j].Sid)-1))
	utils.CheckError(err)
	return iId > jId
}

func sortGameServer(list []*GameServer) []*GameServer {
	sort.Sort(gameServerSlice(list))
	return list
}

//获取游戏服列表
func GetGameServerList(params *GameServerQueryParam) ([]*GameServer, int64) {
	sortOrder := "Sid"
	switch params.Sort {
	case "Sid":
		sortOrder = "Sid"
	}
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	tmpData := make([]*GameServer, 0)
	var count int64
	err := DbCenter.Model(&GameServer{}).Where(&GameServer{
		PlatformId: params.PlatformId,
		Sid:        params.ServerId,
		Node:       params.Node,
	}).Offset(params.Offset).Find(&tmpData).Offset(0).Count(&count).Error
	utils.CheckError(err)
	sortGameServer(tmpData)
	data := make([]*GameServer, 0, params.Limit)
	i := 0
	j := 0
	for _, e := range tmpData {
		if i >= params.Offset {
			if j < params.Limit {
				data = append(data, e)
				j++
			}
		}
		i++
	}
	for _, e := range data {
		serverNode, err := GetServerNode(e.Node)
		e.DbVersion = GetDbVersion(e.Node)
		utils.CheckError(err)

		if err == nil {
			e.State = serverNode.State
			e.Database = serverNode.DbName
			e.OpenTime = serverNode.OpenTime
			e.Ip = serverNode.Ip
			e.ZoneNode = serverNode.ZoneNode
			e.RunState = serverNode.RunState
			e.StartTime = GetNodeStartTime(e.Node)
			e.OnlineCount = GetNowOnlineCountByNode(e.Node)
			e.CreateRoleCount = GetTotalCreateRoleCountByNode(e.Node)
		}

	}
	return data, count
}

//func GetGameServerList(params *GameServerQueryParam) ([]*GameServer, int64) {
//	sortOrder := "Sid"
//	switch params.Sort {
//	case "Sid":
//		sortOrder = "Sid"
//	}
//	if params.Order == "descending" {
//		sortOrder = sortOrder + " desc"
//	}
//	data := make([]*GameServer, 0)
//	var count int64
//	err := DbCenter.Model(&GameServer{}).Where(&GameServer{
//		PlatformId: params.PlatformId,
//		Sid:        params.ServerId,
//		Node:       params.Node,
//	}).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Count(&count).Error
//	utils.CheckError(err)
//	for _, e := range data {
//		serverNode, err := GetServerNode(e.Node)
//		e.DbVersion = GetDbVersion(e.Node)
//		utils.CheckError(err)
//
//		if err == nil {
//			e.State = serverNode.State
//			e.OpenTime = serverNode.OpenTime
//			e.Ip = serverNode.Ip
//			e.ZoneNode = serverNode.ZoneNode
//			e.RunState = serverNode.RunState
//			e.StartTime = GetNodeStartTime(e.Node)
//			e.OnlineCount = GetNowOnlineCountByNode(e.Node)
//			e.CreateRoleCount = GetTotalCreateRoleCountByNode(e.Node)
//		}
//
//	}
//	return data, count
//}

//获取平台游戏服列表
func GetPlatformIdAllGameServerList(platformId string) []*GameServer {
	data := make([]*GameServer, 0)
	gameServer := &GameServer{
		PlatformId: platformId,
	}
	err := DbCenter.Model(&gameServer).Where(&gameServer).Find(&data).Error
	utils.CheckError(err)
	return data
}

func GetMaxGameServerByPlatformId(platformId string) (*GameServer, int, error) {
	l := GetPlatformIdAllGameServerList(platformId)
	maxId := -1
	var maxGameServer *GameServer
	for _, e := range l {
		//g.Log().Info("GetMaxGameServerByPlatformId:%+v, %+v", e, platformId)
		//sid := strings.Split(e.Sid, "@")[0]
		if e.PlatformId != platformId {
			continue
		}
		sid := e.Sid
		//g.Log().Debug("sid:%s", sid)

		//g.Log().Debug("2222:%s", SubString(sid, 1, len(sid)-1))
		id, err := strconv.Atoi(SubString(sid, 1, len(sid)-1))
		utils.CheckError(err)
		if err != nil {
			return nil, 0, err
		}
		if id > maxId {
			maxId = id
			maxGameServer = e
		}
	}
	return maxGameServer, maxId, nil
}

func GetThisIpMaxPort(ip string) int {
	l := GetAllServerNodeList()
	maxPort := g.Cfg().GetInt("game.game_min_port", 10001)
	for _, e := range l {
		nodeIp := strings.Split(e.Node, "@")[1]
		if nodeIp == ip {
			if e.Port > maxPort {
				maxPort = e.Port
			}
			if e.WebPort > maxPort {
				maxPort = e.WebPort
			}
		}
	}
	return maxPort
}

//获取最新的跨服节点
func GetLatestZone() (*ServerNode, int, error) {
	l := GetAllServerNodeByType(2)
	maxId := -1
	var maxNode *ServerNode
	for _, e := range l {
		nodeName := strings.Split(e.Node, "@")[0]

		id, err := strconv.Atoi(SubString(nodeName, 1, len(nodeName)-1))
		utils.CheckError(err)
		if err != nil {
			return nil, 0, err
		}
		if id > maxId {
			maxId = id
			maxNode = e
		}
	}
	return maxNode, maxId, nil
}

//获取最新的跨服节点
func GetLatestZoneByPlatformId(platformId string) (*ServerNode, int, error) {
	//l := GetAllServerNodeByType(2)
	l := GetPlatformIdServerNodeByType(platformId, 2)
	maxId := 0
	var maxNode *ServerNode
	for _, e := range l {
		if e.PlatformId != platformId {
			continue
		}
		nodeName := strings.Split(e.Node, "@")[0]

		if strings.Contains(nodeName, "_") {
			nodeName2 := strings.Split(nodeName, "_")[1]
			id, err := strconv.Atoi(SubString(nodeName2, 1, len(nodeName2)-1))
			utils.CheckError(err)
			if err != nil {
				return nil, 0, err
			}
			if id > maxId {
				maxId = id
				maxNode = e
			}
		} else {
			id, err := strconv.Atoi(SubString(nodeName, 1, len(nodeName)-1))
			utils.CheckError(err)
			if err != nil {
				return nil, 0, err
			}
			if id > maxId {
				maxId = id
				maxNode = e
			}
		}

	}
	return maxNode, maxId, nil
}

// 获得平台跨服
func GetFreeZoneByPlatformId(platformId string) (string, error) {
	platform, err := GetPlatformOne(platformId)
	utils.CheckError(err, "获取平台失败:"+platformId)
	if err != nil {
		return "", err
	}

	inventoryDatabase, err := GetInventoryDatabaseOne(platform.InventoryDatabaseId)
	utils.CheckError(err, "获取数据库配置失败")
	if err != nil {
		return "", err
	}

	inventoryServer, err := GetInventoryServerOneDirty(platform.ZoneInventoryServerId)
	utils.CheckError(err, "获取跨服配置失败")
	if err != nil {
		return "", err
	}

	serverNode, intZid, err := GetLatestZoneByPlatformId(platformId)
	if intZid > 0 {
		connectCount := GetZoneConnectNodeCount(serverNode.Node)
		g.Log().Infof("最新的跨服节点:%s, 连接的游戏节点个数:%d", serverNode.Node, connectCount)
		if connectCount <= 2 {
			return serverNode.Node, nil
		}
	} else {
		g.Log().Infof("初始跨服节点数据 %v%v", platformId, intZid)
	}
	//if serverNode == nil || intZid == -1 {
	//	return "", gerror.New("没有对应的跨服节点")
	//}
	//utils.CheckError(err)
	//if err != nil {
	//	return "", err
	//}

	newIntZid := intZid + 1
	newNode := fmt.Sprintf("%s_z%d@%s", platformId, newIntZid, inventoryServer.InnerIp)
	g.Log().Infof("新跨服节点:%s", newNode)
	out, err := AddServerNode(newNode, inventoryServer.InnerIp, 0, 0, 2, platformId, inventoryDatabase.Host, inventoryDatabase.Port, fmt.Sprintf("db_%s_zone_z%d", platformId, newIntZid))
	utils.CheckError(err, "新增节点失败:"+out)
	if err != nil {
		return "", err
	}

	//time.Sleep(time.Duration(5) * time.Second)

	for i := 0; i < 30; i++ {
		g.Log().Infof("等待跨服节点(%s)数据写入[%d]......", newNode, i)
		time.Sleep(time.Duration(1) * time.Second)
		isExists := IsServerNodeExists(newNode)
		if isExists == true {
			break
		}
	}
	g.Log().Infof("跨服节点(%s)数据写入成功.", newNode)

	err = InstallNode(newNode)
	utils.CheckError(err, "部署节点失败")
	if err != nil {
		return "", err
	}

	err = NodeAction([]string{newNode}, "start")
	utils.CheckError(err, "启动节点失败")
	if err != nil {
		return "", err
	}
	return newNode, nil
}

// 获得平台战区服
func GetFreeWarByPlatformId(platformId string) error {
	serverType := 7
	platform, err := GetPlatformOne(platformId)
	utils.CheckError(err, "获取平台失败:"+platformId)
	if err != nil {
		return err
	}
	warL := GetPlatformIdServerNodeByType(platformId, serverType)
	if len(warL) > 0 {
		return nil
	}
	inventoryDatabase, err := GetInventoryDatabaseOne(platform.InventoryDatabaseId)
	utils.CheckError(err, "获取数据库配置失败")
	if err != nil {
		return err
	}

	inventoryServer, err := GetInventoryServerOneDirty(platform.ZoneInventoryServerId)
	utils.CheckError(err, "获取跨服配置失败")
	if err != nil {
		return err
	}

	newNode := fmt.Sprintf("%s_war@%s", platformId, inventoryServer.InnerIp)
	g.Log().Infof("新战区服节点:%s", newNode)
	out, err := AddServerNode(newNode, inventoryServer.InnerIp, 0, 0, serverType, platformId, inventoryDatabase.Host, inventoryDatabase.Port, fmt.Sprintf("db_%s_war", platformId))
	utils.CheckError(err, "新增节点失败:"+out)
	if err != nil {
		return err
	}
	//time.Sleep(time.Duration(5) * time.Second)

	for i := 0; i < 30; i++ {
		g.Log().Infof("等待战区服节点(%s)数据写入[%d]......", newNode, i)
		time.Sleep(time.Duration(1) * time.Second)
		isExists := IsServerNodeExists(newNode)
		if isExists == true {
			break
		}
	}
	g.Log().Infof("战区服节点(%s)数据写入成功.", newNode)

	err = InstallNode(newNode)
	utils.CheckError(err, "部署节点失败")
	if err != nil {
		return err
	}

	err = NodeAction([]string{newNode}, "start")
	utils.CheckError(err, "启动节点失败")
	if err != nil {
		return err
	}
	return nil
}

// 获取跨服节点连接的数量
func GetZoneConnectNodeCount(node string) int {
	var data struct {
		Count int
	}
	sql := fmt.Sprintf(
		`SELECT count(1) as count FROM c_server_node WHERE type = 1 and zone_node = '%s'`, node)
	err := DbCenter.Raw(sql).Scan(&data).Error
	utils.CheckError(err)
	return data.Count
}
func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}
	// 返回子串
	return string(rs[begin:end])
}

func AutoCreateAndOpenServer(platformId string, isCheck bool, openServerTime int) error {
	return nil
}

// 自动人数到达开服
func AutoCreateAndCreateRoleLimit(platformId string, openServerTime int) error {
	return autoCreateAndOpenServerHandle(platformId, false, true, openServerTime)
}

// 自动定时开服
func AutoCreateAndOpenServerTime(platformId string, isCheckTime bool, openServerTime int) error {
	return autoCreateAndOpenServerHandle(platformId, isCheckTime, false, openServerTime)
}

// 自动开服操作
func autoCreateAndOpenServerHandle(platformId string, isCheckTime, isCheckRoleLimit bool, openServerTime int) error {
	Key := "AutoCreateAndOpenServerHandle" + "_" + platformId
	OldTime := utils.GetCacheInt64(Key)
	if OldTime > 0 {
		g.Log().Warningf("正在开服中:%s !!!! key:%s %d", platformId, Key, OldTime)
		if isCheckRoleLimit {
			return nil
		}
		return gerror.New("当前正在开服中")
	}

	t0 := time.Now()

	platform, err := GetPlatformOne(platformId)
	utils.CheckError(err, "获取平台失败:"+platformId)
	if err != nil {
		return err
	}

	intervalTime := 5 * 60 // 前一个服到当前服间隔几秒
	if isCheckTime || isCheckRoleLimit {
		if platform.IsAutoOpenServer == 0 {
			if isCheckTime {
				g.Log().Warningf("自动开服关闭,定时和人数开数开服无效:%+v", platformId)
			}
			//不自动开服
			return nil
		}
		if isCheckRoleLimit { // 人数限制22点后不开服
			if time.Now().Hour() >= 22 {
				//g.Log().Info("晚上10点后不自动开服")
				return nil
			}
			if platform.CreateRoleLimit < 1 {
				err = gerror.New(fmt.Sprintf("自动开服人数配置错误:%s, %d", platform.Id, platform.CreateRoleLimit))
				utils.CheckError(err)
				return err
			}
		}
	} else {
		g.Log().Infof("立即开服:%s......", platformId)
		intervalTime = 60
	}

	inventoryDatabase, err := GetInventoryDatabaseOne(platform.InventoryDatabaseId)
	utils.CheckError(err, "获取数据库配置失败")
	if err != nil {
		return err
	}

	maxGameServer, intSid, err := GetMaxGameServerByPlatformId(platformId)
	if intSid == -1 {
		g.Log().Infof("初始游戏服：%+v", platformId)
		intSid = 0
	} else if err != nil {
		utils.CheckError(err, "获取最大区服失败")
		return err
	}
	if intSid > 0 { // 有游戏服再验证
		gameDb, err := GetGameDbByNode(maxGameServer.Node)
		utils.CheckError(err, "连接游戏服数据库失败")
		if err != nil {
			return err
		}
		defer gameDb.Close()
		if isCheckRoleLimit { // 验证创角数
			createRoleCount := GetTotalCreateRoleCount(gameDb)
			if createRoleCount < platform.CreateRoleLimit {
				return nil
			}
		}
		serverNode, err := GetServerNode(maxGameServer.Node)
		utils.CheckError(err)
		if serverNode.OpenTime == 0 {
			g.Log().Warningf("最近的服务器还在开启中%s", maxGameServer.Node)
			return gerror.New(fmt.Sprintf("最近的服务器还在开启中%s", maxGameServer.Node))
		}
		if intervalTime+serverNode.OpenTime >= openServerTime {
			g.Log().Error("自动开服时间间隔小于%d秒:%s 最后的开服时间:%s 当前要开服的时间:%s", intervalTime, platformId, utils.TimeIntFormDefault(serverNode.OpenTime), utils.TimeIntFormDefault(openServerTime))
			//UpdatePlatformOpenServerTime(platformId, 0)
			return gerror.New(fmt.Sprintf("自动开服时间间隔小于%d秒", intervalTime))
		}
	}
	//g.Log().Info("最大区服:%+v", maxGameServer)
	//g.Log().Info("最大区服ID:%+v(%d)", maxGameServer.Sid, intSid)

	//g.Log().Info("最新区服:%s, 当前创角:%d, 创角临界值:%d", maxGameServer.Sid, createRoleCount, platform.CreateRoleLimit)
	newIntSid := intSid + 1
	newSid := fmt.Sprintf("s%d", newIntSid)
	//if isCheckOpenServerTake == true {
	//	serverNode, err := GetServerNode(maxGameServer.Node)
	//	utils.CheckError(err)
	//	if serverNode.OpenTime == 0 || 5*60 > serverNode.OpenTime-openServerTime && serverNode.OpenTime-openServerTime > -5*60 {
	//		g.Log().Error("定时开服时间间隔小于5分钟:%s serverNode.OpenTime:%d openServerTime:%d", platformId, serverNode.OpenTime, openServerTime)
	//		UpdatePlatformOpenServerTime(platformId, 0)
	//		return gerror.New("定时开服时间间隔小于5分钟")
	//	}
	//}
	//if isCheckTime == false || isCheckOpenServerTake || platform.IsAutoOpenServer != 0 && createRoleCount >= platform.CreateRoleLimit {
	OldTime = utils.GetCacheInt64(Key)
	if OldTime > 0 {
		return nil
	}
	utils.SetCache(Key, gtime.Timestamp(), 0)
	defer utils.DelCache(Key)
	g.Log().Infof("*************************** 开服部署新服 %s - %s *****************************%s\n", platformId, newSid, gtime.Datetime())

	//serverNode, err := GetServerNode(maxGameServer.Node)
	//utils.CheckError(err, "获取节点失败!!")
	//if err != nil {
	//	return err
	//}
	maxFreeServer, err := GetMaxFreeServerByPlatformId(platformId)
	utils.CheckError(err)
	if err != nil {
		return err
	}

	g.Log().Infof("最空闲的服务器:%+v", maxFreeServer)
	g.Log().Infof("新区服id:%s", newSid)
	newNode := fmt.Sprintf("%s_%s@%s", platformId, newSid, maxFreeServer.InnerIp)
	g.Log().Infof("新区服节点:%s", newNode)

	//maxPort := Max(serverNode.Port, serverNode.WebPort)
	maxPort := GetThisIpMaxPort(maxFreeServer.InnerIp)
	//g.Log().Info("最大端口:%d", maxPort)
	out, err := AddServerNode(newNode, maxFreeServer.Host, maxPort+1, maxPort+2, 1, platformId, inventoryDatabase.Host, inventoryDatabase.Port, fmt.Sprintf("db_%s_game_%s", platformId, newSid))
	utils.CheckError(err, "新增节点失败:"+out)
	if err != nil {
		return err
	}
	//if isCheckOpenServerTake {
	//	UpdatePlatformOpenServerTime(platformId, 0)
	//}
	g.Log().Infof("添加节点成功:%s", newNode)
	zoneNode, err := GetFreeZoneByPlatformId(platformId)
	utils.CheckError(err, "获取空闲跨服节点失败:"+out)
	if err != nil {
		return err
	}
	err = GetFreeWarByPlatformId(platformId)
	utils.CheckError(err, "获得战区服节点失败:"+platformId)
	if err != nil {
		return err
	}
	serverNameStr := fmt.Sprintf("%d区", newIntSid)
	//if strings.Index(platform.ServerAliasStr, "%d") != -1 {
	//	serverNameStr = fmt.Sprintf(platform.ServerAliasStr, newIntSid)
	//}
	if len(strings.Split(platform.ServerAliasStr, "%d")) == 2 {
		serverNameStr = fmt.Sprintf(platform.ServerAliasStr, newIntSid)
	}
	//serverNameStr := getPlatformServerEnterName(platformId, newIntSid)
	g.Log().Infof("获得平台区服入口名字: %s  %s", platformId, serverNameStr)
	out, err = AddGameServer(platformId, newSid, serverNameStr, newNode, zoneNode, 3, openServerTime, 1)

	utils.CheckError(err, "新增游戏服失败:"+out)
	if err != nil {
		return err
	}
	g.Log().Infof("添加game_server成功:%s", newSid)
	//time.Sleep(time.Duration(15) * time.Second)

	for i := 0; i < 30; i++ {
		g.Log().Infof("等待节点(%s)数据写入[%d]......", newNode, i)
		time.Sleep(time.Duration(1) * time.Second)
		isExists := IsServerNodeExists(newNode)
		if isExists == true {
			break
		}
	}
	g.Log().Infof("节点(%s)数据写入成功.", newNode)

	err = InstallNode(newNode)
	utils.CheckError(err, "部署节点失败")
	if err != nil {
		return err
	}

	err = NodeAction([]string{newNode}, "start")
	utils.CheckError(err, "启动节点失败")
	if err != nil {
		return err
	}

	err = AfterAddGameServer()
	utils.CheckError(err, "同步登录充值战区失败")
	if err != nil {
		return err
	}

	err = NodeAction([]string{zoneNode}, "pull")
	utils.CheckError(err, "更新跨服节点数据")
	if err != nil {
		return err
	}

	err = CreateAnsibleInventory()
	utils.CheckError(err, "生成ansible inventory失败")
	if err != nil {
		return err
	}
	usedTime := time.Since(t0)
	g.Log().Infof("************************ 自动开服成功:%s - %s 耗时:%s **********************", platformId, newSid, usedTime.String())
	//} else {
	//	//g.Log().Info("不满足开服条件.")
	//}
	return nil
}

func OpenServerNow(platformId string) error {
	g.Log().Infof("现在开服:%s", platformId)
	err := AutoCreateAndOpenServer(platformId, false, utils.GetTimestamp())
	utils.CheckError(err, "开服失败!!!!!!"+platformId)
	return err
}

// 开服方式
func OpenServerType(userId int, platformId string, openServerTime int, hOpenServerTime int) error {
	currDateTimeStr := gtime.Datetime()
	if openServerTime > 0 {
		g.Log().Infof("定时自动开服:%s %v", currDateTimeStr)
		err := AutoCreateAndOpenServerTime(platformId, true, openServerTime)
		if err != nil {
			utils.CheckError(err, "定时自动开服失败!!!!!!"+platformId)
			SendBackgroundMsgTemplateHandleByUser(userId, MSG_TEMPLATE_AUTO_CREATE_SERVER_FAIL, platformId, platformId+"定时自动开服失败", "详情："+utils.TimeIntFormDefault(openServerTime)+" err:"+fmt.Sprintf("%v", err))
		}
		return err
	}
	if hOpenServerTime > 0 {
		g.Log().Infof("整点自动开服:%s", currDateTimeStr)
		err := AutoCreateAndOpenServerTime(platformId, true, hOpenServerTime)
		if err != nil {
			utils.CheckError(err, "整点自动开服失败!!!!!!"+platformId)
			SendBackgroundMsgTemplateHandleByUser(userId, MSG_TEMPLATE_H_CREATE_SERVER_FAIL, platformId, platformId+"整点自动开服失败", "详情："+utils.TimeIntFormDefault(hOpenServerTime)+" err:"+fmt.Sprintf("%v", err))
		}
		return err
	}
	now := utils.GetTimestamp()
	g.Log().Infof("立即开服:%s %v", platformId, currDateTimeStr)
	err := AutoCreateAndOpenServerTime(platformId, false, now)
	if err != nil {
		utils.CheckError(err, "立即开服失败!!!!!!"+platformId)
		SendBackgroundMsgTemplateHandleByUser(userId, MSG_TEMPLATE_NOW_CREATE_SERVER_FAIL, platformId, platformId+"立即开服失败", "详情："+currDateTimeStr+" err:"+fmt.Sprintf("%v", err))
	}
	return err
}

//func AutoCreateAndOpenServer(isCheck bool) error {
//	if IsNowOpenServer == true {
//		g.Log().Warning("正在开服中!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
//		return nil
//	}
//	t0 := time.Now()
//
//
//	utils.CheckError(err, "读取自动开服人数配置失败")
//	if err != nil {
//		return err
//	}
//
//	if configDbHost == "" {
//		g.Log().Error("读取配置游戏服连接的数据库配置失败")
//		return err
//	}
//
//	if isCheck {
//		g.Log().Info("检测自动开服......")
//		// 检测是否满足开服条件
//		utils.CheckError(err, "读取是否开启自动开服配置失败")
//		if err != nil {
//			return err
//		}
//		if isAutoOpenServer == false {
//			return err
//		}
//
//		if time.Now().Hour() >= 22 {
//			g.Log().Info("晚上10点后不自动开服")
//			return err
//		}
//	} else {
//		g.Log().Info("立即开服......")
//	}
//
//	now := utils.GetTimestamp()
//
//	maxGameServer, intSid, err := GetMaxGameServer()
//	utils.CheckError(err, "获取最大区服失败")
//	if err != nil {
//		return err
//	}
//	//g.Log().Info("最大区服:%+v", maxGameServer)
//	//g.Log().Info("最大区服ID:%+v(%d)", maxGameServer.Sid, intSid)
//	gameDb, err := GetGameDbByNode(maxGameServer.Node)
//	utils.CheckError(err, "连接游戏服数据库失败")
//	if err != nil {
//		return err
//	}
//	defer gameDb.Close()
//	createRoleCount := GetTotalCreateRoleCount(gameDb)
//	g.Log().Info("最新区服:%s, 当前创角:%d, 创角临界值:%d", maxGameServer.Sid, createRoleCount, openServerRoleCount)
//
//	if isCheck == false || createRoleCount >= openServerRoleCount {
//		IsNowOpenServer = true
//		defer func() {
//			IsNowOpenServer = false
//		}()
//		g.Log().Info("*************************** 开服部署新服 *****************************\n")
//		newIntSid := intSid + 1
//		newSid := fmt.Sprintf("s%d", newIntSid)
//		serverNode, err := GetServerNode(maxGameServer.Node)
//		utils.CheckError(err, "获取节点失败!!")
//		if err != nil {
//			return err
//		}
//		maxFreeServer, err := GetMaxFreeServer()
//		utils.CheckError(err)
//		if err != nil {
//			return err
//		}
//
//		g.Log().Info("最空闲的服务器:%+v", maxFreeServer)
//		g.Log().Info("新服id:%s", newSid)
//		newNode := fmt.Sprintf("%s_%s@%s", serverNode.PlatformId, newSid, maxFreeServer.InnerIp)
//		g.Log().Info("新节点:%s", newNode)
//
//		//maxPort := Max(serverNode.Port, serverNode.WebPort)
//		maxPort := GetThisIpMaxPort(maxFreeServer.InnerIp)
//		g.Log().Info("最大端口:%d", maxPort)
//		out, err := AddServerNode(newNode, maxFreeServer.Host, maxPort+1, maxPort+2, 1, serverNode.PlatformId, configDbHost, 3306, fmt.Sprintf("db_%s_game_%s", serverNode.PlatformId, newSid))
//		utils.CheckError(err, "新增节点失败:"+out)
//		if err != nil {
//			return err
//		}
//
//		zoneNode, err := GetFreeZone()
//		utils.CheckError(err, "获取空闲跨服节点失败:"+out)
//		if err != nil {
//			return err
//	}
//
//		out, err = AddGameServer(maxGameServer.PlatformId, newSid, fmt.Sprintf("%d区", newIntSid), newNode, zoneNode, 3, now, 1)
//
//		utils.CheckError(err, "新增游戏服失败:"+out)
//		if err != nil {
//			return err
//		}
//
//		//time.Sleep(time.Duration(15) * time.Second)
//
//		for i := 0; i < 30; i++ {
//			g.Log().Info("等待节点(%s)数据写入[%d]......", newNode, i)
//			time.Sleep(time.Duration(1) * time.Second)
//			isExists := IsServerNodeExists(newNode)
//			if isExists == true {
//				break
//			}
//		}
//		g.Log().Info("节点(%s)数据写入成功.", newNode)
//
//		err = InstallNode(newNode)
//		utils.CheckError(err, "部署节点失败")
//		if err != nil {
//			return err
//		}
//
//		err = NodeAction([] string{newNode}, "start")
//		utils.CheckError(err, "启动节点失败")
//		if err != nil {
//			return err
//		}
//
//
//		err = RefreshGameServer()
//		utils.CheckError(err, "刷新区服入口失败")
//		if err != nil {
//			return err
//		}
//
//		err = NodeAction([] string{zoneNode}, "pull")
//		utils.CheckError(err, "更新跨服节点数据")
//		if err != nil {
//			return err
//		}
//
//
//
//		err = CreateAnsibleInventory()
//		utils.CheckError(err, "生成ansible inventory失败")
//		if err != nil {
//			return err
//		}
//		usedTime := time.Since(t0)
//		g.Log().Info("************************ 自动开服成功:%s 耗时:%s **********************", newSid, usedTime.String())
//	} else {
//		g.Log().Info("不满足开服条件.")
//	}
//	return nil
//}

// 获取单个游戏服
func GetGameServerOne(platformId string, id string) (*GameServer, error) {
	if gstr.Trim(platformId) == "" || gstr.Trim(id) == "" {
		return nil, sql.ErrNoRows
	}
	gameServer := &GameServer{
		Sid:        id,
		PlatformId: platformId,
	}
	err := DbCenter.First(&gameServer).Error
	return gameServer, err
}

func IsGameServerExists(platformId string, id string) bool {
	gameServer := &GameServer{
		Sid:        id,
		PlatformId: platformId,
	}
	return !DbCenter.First(&gameServer).RecordNotFound()
}

// 获取该节点关联的所有游戏服
func GetGameServerByNode(node string) []*GameServer {
	data := make([]*GameServer, 0)
	err := DbCenter.Model(&GameServer{}).Where(&GameServer{
		Node: node,
	}).Find(&data).Error
	utils.CheckError(err)
	return data
}

// 获取该平台所有游戏服
func GetGameServerByPlatformId(platformId string) []*GameServer {
	data := make([]*GameServer, 0)
	err := DbCenter.Model(&GameServer{}).Where(&GameServer{
		PlatformId: platformId,
	}).Find(&data).Error
	utils.CheckError(err)
	return data
}

func GetGameServerIdListStringByNode(node string) string {
	serverIdList := GetGameServerIdListByNode(node)
	return "'" + strings.Join(serverIdList, "','") + "'"
}
func GetGameServerIdListByNode(node string) []string {
	data := make([]*GameServer, 0)
	serverIdList := make([]string, 0)
	err := DbCenter.Model(&GameServer{}).Where(&GameServer{
		Node: node,
	}).Find(&data).Error
	utils.CheckError(err)
	for _, e := range data {
		serverIdList = append(serverIdList, e.Sid)
	}
	return serverIdList
}
