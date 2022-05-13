package models

import (
	"fmt"
	"gfWeb/library/enums"
	"gfWeb/library/utils"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/jinzhu/gorm"

	//"golang.org/x/net/websocket"
	"strings"
)

//SESSION前缀
const (
	USER_SESSION_MARK = "user"
)

type Result struct {
	Code enums.ResultCode `json:"code"`
	Msg  string           `json:"msg"`
	Data interface{}      `json:"data"`
}

type BaseQueryParam struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

// 通过平台id和区服id获取 gorm.DB 实例
func GetGameDbByPlatformIdAndSid(platformId string, Sid string) (gameDb *gorm.DB, err error) {
	gameServer, err := GetGameServerOne(platformId, Sid)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	serverNode, err := GetServerNode(gameServer.Node)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	return GetGameDbByNode(serverNode.Node)
}

// 通过平台id和区服id获取 gorm.DB 实例
func GetGameURLByNode(node string) string {
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	if err != nil {
		return ""
	}
	url := fmt.Sprintf("http://%s:%d", serverNode.Ip, serverNode.WebPort)
	return url
}

// 通过平台id和区服id获取 gorm.DB 实例
func GetGameURLByPlatformIdAndSid(platformId string, Sid string) string {
	gameServer, err := GetGameServerOne(platformId, Sid)
	utils.CheckError(err)
	if err != nil {
		return ""
	}
	serverNode, err := GetServerNode(gameServer.Node)
	utils.CheckError(err)
	if err != nil {
		return ""
	}
	url := fmt.Sprintf("http://%s:%d", serverNode.Ip, serverNode.WebPort)
	return url
}

// 通过node获取 gorm.DB 实例
func GetGameDbByNode2(node string, platformId string, sid string) (gameDb *gorm.DB, err error) {
	if node == "" {
		return nil, gerror.New("节点不能未空")
	}
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	array := strings.Split(serverNode.Node, "@")
	if len(array) != 2 {
		return nil, gerror.New("解析节点名字失败:" + serverNode.Node)
	}
	//gameDbName := "game_" + array[0]
	gameDbName := serverNode.DbName
	gameDbHost := serverNode.DbHost
	if platformId == "qq" && (sid == "s1" || sid == "s2" || sid == "s3" || sid == "s4" || sid == "s5" || sid == "s6") {
		gameDbName = "test_" + sid
		gameDbHost = "10.66.253.43"
	}

	gameDbPort := serverNode.DbPort
	gameDbPwd := g.Cfg().GetString("database.game.pass")
	dbUserName := g.Cfg().GetString("database.default.user")
	dbArgs := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", dbUserName, gameDbPwd, gameDbHost, gameDbPort, gameDbName)
	//dbArgs := "root:game1234@tcp(" + serverNode.Ip + ":3306)/" + gameDbName +"?charset=utf8&parseTime=True&loc=Local"
	gameDb, err = gorm.Open("mysql", dbArgs)
	if err != nil {
		g.Log().Error("连接节点(%v)数据库失败:%v", node, dbArgs)
		return nil, err
	}
	gameDb.SingularTable(true)
	return gameDb, err
}

// 通过node获取 gorm.DB 实例
func GetGameDbByNode(node string) (gameDb *gorm.DB, err error) {
	if node == "" {
		return nil, gerror.New("节点不能未空")
	}
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	array := strings.Split(serverNode.Node, "@")
	if len(array) != 2 {
		return nil, gerror.New("解析节点名字失败:" + serverNode.Node)
	}
	//gameDbName := "game_" + array[0]
	gameDbName := serverNode.DbName
	gameDbHost := serverNode.DbHost
	gameDbPort := serverNode.DbPort
	gameDbPwd := g.Cfg().GetString("database.game.pass")
	dbUserName := g.Cfg().GetString("database.default.user")
	dbArgs := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", dbUserName, gameDbPwd, gameDbHost, gameDbPort, gameDbName)
	//dbArgs := fmt.Sprintf("root:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", gameDbPwd, gameDbHost, gameDbPort, gameDbName)
	//dbArgs := "root:game1234@tcp(" + serverNode.Ip + ":3306)/" + gameDbName +"?charset=utf8&parseTime=True&loc=Local"
	gameDb, err = gorm.Open("mysql", dbArgs)
	if err != nil {
		g.Log().Errorf("连接节点(%v)数据库失败:%v", node, dbArgs)
		return nil, err
	}
	gameDb.SingularTable(true)
	return gameDb, err
}

// 通过平台id获取 gorm.DB 实例
func GetGameDbByPlatformId(PlatformId string) (gameDb *gorm.DB, err error) {
	platform, err := GetPlatformOne(PlatformId)
	utils.CheckError(err, "获取平台失败:"+PlatformId)
	if err != nil {
		return nil, err
	}
	inventoryDatabase, err := GetInventoryDatabaseOne(platform.InventoryDatabaseId)
	utils.CheckError(err, "获取数据库配置失败")
	if err != nil {
		return nil, err
	}
	gameDbHost := inventoryDatabase.Host
	gameDbPort := inventoryDatabase.Port
	gameDbPwd := g.Cfg().GetString("database.game.pass")
	dbUserName := g.Cfg().GetString("database.default.user")
	dbArgs := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql?charset=utf8&parseTime=True&loc=Local", dbUserName, gameDbPwd, gameDbHost, gameDbPort)
	//dbArgs := fmt.Sprintf("root:%s@tcp(%s:%d)/mysql?charset=utf8&parseTime=True&loc=Local", gameDbPwd, gameDbHost, gameDbPort)
	g.Log().Debug(dbArgs)
	gameDb, err = gorm.Open("mysql", dbArgs)
	if err != nil {
		g.Log().Error("连接平台(%v)游戏服数据库失败:%v", PlatformId, dbArgs)
		return nil, err
	}
	gameDb.SingularTable(true)
	return gameDb, err
}

// 通过node获取 gorm.DB 实例
func GetVisitorGameDbByNode(node string) (gameDb *gorm.DB, err error) {
	if node == "" {
		return nil, gerror.New("节点不能未空")
	}
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	//gameDbName := "game_" + array[0]
	gameDbName := serverNode.DbName
	gameDbHost := serverNode.DbHost
	gameDbPort := serverNode.DbPort
	//gameDbPwd := g.Cfg().GetString("database.game.pass")
	gameDbPwd := g.Cfg().GetString("database.game.visitor_pass")
	dbArgs := fmt.Sprintf("visitor:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", gameDbPwd, gameDbHost, gameDbPort, gameDbName)
	//dbArgs := "root:game1234@tcp(" + serverNode.Ip + ":3306)/" + gameDbName +"?charset=utf8&parseTime=True&loc=Local"
	gameDb, err = gorm.Open("mysql", dbArgs)
	if err != nil {
		g.Log().Error("连接节点(%v)数据库失败:%v", node, dbArgs)
		return nil, err
	}
	gameDb.SingularTable(true)
	return gameDb, err
}

// 通过平台id和区服id 获取ip地址和端口
func GetIpAndPortByPlatformIdAndSid(platformId string, Sid string) (string, int, error) {
	gameServer, err := GetGameServerOne(platformId, Sid)
	utils.CheckError(err)
	serverNode, err := GetServerNode(gameServer.Node)
	utils.CheckError(err)
	return serverNode.Ip, serverNode.Port, err
}

//
//func GetWsByPlatformIdAndSid(platformId string, Sid string) (*websocket.Conn, error) {
//	ip, port, err := GetIpAndPortByPlatformIdAndSid(platformId, Sid)
//	if err != nil {
//		return nil, err
//	}
//	wsUrl := fmt.Sprintf("ws://%s:%d", ip, port)
//	ws, err := websocket.Dial(wsUrl, "", wsUrl)
//	return ws, err
//}
//
//func GetWsByNode(node string) (*websocket.Conn, error) {
//	serverNode, err := GetServerNode(node)
//	if err != nil {
//		return nil, err
//	}
//	wsUrl := fmt.Sprintf("ws://%s:%d", serverNode.Ip, serverNode.Port)
//	ws, err := websocket.Dial(wsUrl, "", wsUrl)
//	return ws, err
//}

type Server struct {
	PlatformId string `json:"platformId"`
	Sid        string `json:"serverId"`
	Desc       string `json:"desc"`
	Node       string `json:"node"`
}

func GetServerList(platformIdList []string) []*Server {
	serverList := make([]*Server, 0)
	//gameServerNodeList := GetGameServerByPlatformId()
	for _, platformId := range platformIdList {
		gameServerList := GetGameServerByPlatformId(platformId)
		for _, gameServer := range gameServerList {
			server := &Server{
				PlatformId: gameServer.PlatformId,
				Sid:        gameServer.Sid,
				Desc:       gameServer.Desc,
				Node:       gameServer.Node,
			}
			serverList = append(serverList, server)
		}
	}
	return serverList
}

func TranPlayerNameSting2PlayerIdList(platformId string, playerName string) ([]int, error) {
	playerIdList := make([]int, 0)
	nameList := strings.Split(playerName, ",")
	for _, name := range nameList {
		player, err := GetPlayerByPlatformIdAndNickname(platformId, name)
		if err != nil {
			return nil, err
		}
		playerIdList = append(playerIdList, player.Id)
	}
	return playerIdList, nil
}
