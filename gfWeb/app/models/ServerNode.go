package models

import (
	"fmt"
	"gfWeb/library/utils"
	"strings"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

type ServerNodeQueryParam struct {
	BaseQueryParam
	Type       int
	Node       string
	PlatformId string `json:"platformId"`
}

type ServerNode struct {
	Node          string `gorm:"primary_key" json:"node"`
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	WebPort       int    `json:"webPort"`
	DbHost        string `json:"dbHost"`
	DbPort        int    `json:"dbPort"`
	DbName        string `json:"dbName"`
	Type          int    `json:"type"`
	ZoneNode      string `json:"zoneNode"`
	ServerVersion int    `json:"serverVersion" gorm:"-"`
	DbVersion     int    `json:"dbVersion" gorm:"-"`
	IsAdd         int    `json:"isAdd" gorm:"-"`
	//ClientVersion string `json:"clientVersion"`
	OpenTime  int `json:"openTime"`
	StartTime int `json:"startTime" gorm:"-"`
	//IsTest        int    `json:"isTest"`
	PlatformId string `json:"platformId"`
	State      int    `json:"state"`
	RunState   int    `json:"runState"`
}

func (t *ServerNode) TableName() string {
	return "c_server_node"
}

var delPlatformNotUseDataState = 0 // 删除没用的数据库状态

//获取分页数据
func ServerNodePageList(params *ServerNodeQueryParam) ([]*ServerNode, int64) {
	data := make([]*ServerNode, 0)
	//默认排序
	sortOrder := "node"
	switch params.Sort {
	case "node":
		sortOrder = "node"
	}
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	err := DbCenter.Model(&ServerNode{}).Where(&ServerNode{
		Type:       params.Type,
		Node:       params.Node,
		PlatformId: params.PlatformId,
	}).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	if err == nil {
		for _, e := range data {
			e.ServerVersion = GetNodeVersion(e.Node)
			e.DbVersion = GetDbVersion(e.Node)
			e.StartTime = GetNodeStartTime(e.Node)
		}
	}
	return data, count
}

func GetAllServerNodeList() []*ServerNode {
	data := make([]*ServerNode, 0)
	err := DbCenter.Model(&ServerNode{}).Find(&data).Error
	utils.CheckError(err)
	return data
}

func GetServerNode(node string) (*ServerNode, error) {
	serverNode := &ServerNode{
		Node: node,
	}
	err := DbCenter.First(&serverNode).Error
	return serverNode, err
}

func IsServerNodeExists(node string) bool {
	serverNode := &ServerNode{
		Node: node,
	}
	return !DbCenter.First(&serverNode).RecordNotFound()
}

// 获取节点版本
func GetNodeVersion(node string) int {
	//return "nullddd"
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return 0
	}
	defer gameDb.Close()
	var data struct {
		Version int
	}

	sql := fmt.Sprintf(
		`SELECT int_data as version FROM server_data where id = 0 `)
	err = gameDb.Raw(sql).Scan(&data).Error
	utils.CheckError(err, node)
	if err != nil {
		return 0
	}
	return data.Version
}

//获取节点合服时间
func GetMergeTime(node string) int {
	//return "nullddd"
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return 0
	}
	defer gameDb.Close()
	var data struct {
		Time int
	}

	sql := fmt.Sprintf(
		`SELECT int_data as time FROM server_data where id = 6 `)
	err = gameDb.Raw(sql).Scan(&data).Error
	utils.CheckError(err, "获取节点合服时间")
	if err != nil {
		return 0
	}
	return data.Time
}

func GetDbVersion(node string) int {
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return -1
	}
	defer gameDb.Close()
	var data struct {
		Version int
	}

	sql := fmt.Sprintf(
		`SELECT version FROM db_version`)
	err = gameDb.Raw(sql).Scan(&data).Error
	if data.Version != 0 {
		utils.CheckError(err)
	}

	if err != nil {
		return -1
	}
	return data.Version
}

func GetNodeStartTime(node string) int {
	//return "nullddd"
	gameDb, err := GetGameDbByNode(node)
	utils.CheckError(err)
	if err != nil {
		return 0
	}
	defer gameDb.Close()
	var data struct {
		Time int
	}

	sql := fmt.Sprintf(
		`SELECT int_data as time FROM server_data where id = 3 `)
	err = gameDb.Raw(sql).Scan(&data).Error
	utils.CheckError(err, "获取节点启动时间")
	if err != nil {
		return 0
	}
	return data.Time
}

// 获取所有游戏节点
//func GetAllGameServerNode() []*ServerNode {
//	//data := make([]*ServerNode, 0)
//	//err := DbCenter.Model(&ServerNode{}).Where(&ServerNode{
//	//	Type: 1,
//	//}).Find(&data).Error
//	//utils.CheckError(err, "查询所有游戏节点失败")
//	return GetAllServerNodeByType(1)
//}

// 获取节点列表
func GetAllServerNodeByType(nodeType int) []*ServerNode {
	data := make([]*ServerNode, 0)
	err := DbCenter.Model(&ServerNode{}).Where(&ServerNode{
		Type: nodeType,
	}).Find(&data).Error
	utils.CheckError(err, "获取节点列表失败")
	return data
}

// 获取平台类型节点列表
func GetPlatformIdServerNodeByType(PlatformId string, nodeType int) []*ServerNode {
	data := make([]*ServerNode, 0)
	err := DbCenter.Model(&ServerNode{}).Where(&ServerNode{
		Type:       nodeType,
		PlatformId: PlatformId,
	}).Find(&data).Error
	utils.CheckError(err, "获取平台类型节点列表失败")
	return data
}

// 获取所有游戏节点数据
func GetAllGameServerNodeByPlatformId(platformId string) []*ServerNode {
	data := make([]*ServerNode, 0)
	err := DbCenter.Model(&ServerNode{}).Where(&ServerNode{
		Type:       1,
		PlatformId: platformId,
	}).Find(&data).Error
	utils.CheckError(err)
	return data
}

// 获取所有游戏节点
func GetAllGameNodeByPlatformId(platformId string) []string {
	data := make([]string, 0)
	serverNodeList := GetAllGameServerNodeByPlatformId(platformId)
	for _, e := range serverNodeList {
		data = append(data, e.Node)
	}
	return data
}

// 获取登录节点 14101 11101
//func GetLoginServerNode() (*ServerNode, error) {
//	serverNode := &ServerNode{}
//	err := DbCenter.Where(&ServerNode{
//		Type: 4,
//	}).First(&serverNode).Error
//	return serverNode, err
//}

func CheckAllGameNode() {
	isCheck := g.Cfg().GetBool("game.is_check_node", false)
	if isCheck {
		now := utils.GetTimestamp()
		serverNodeList := GetAllServerNodeList()
		num := 0
		mapList := make(map[string][]string)
		for _, e := range serverNodeList {
			if e.Type == 1 && e.OpenTime+120 > now {
				g.Log().Debug("开服时间间隔小于检测时间忽略:%s", e.Node)
				continue
			}
			if e.RunState == 0 {
				g.Log().Info("节点未开启:~p", e.Node)
				//UpdateEtsPlatformServerCloseCount(e.PlatformId, e.Node)
				mapList = utils.MapListAdd(mapList, e.PlatformId, e.Node)
				num++
			}
		}

		if num > 0 {
			g.Log().Error("节点:%d个未开启!!!!!!!!!!!!", num)
			//utils.ReportMsg("105138", "13616005067")
			//utils.ReportMsg("105138", "19905929917")
			for mapPlatformId, mapNodeList := range mapList {
				platform := GetPlatformSimpleOne(mapPlatformId)
				platformMergeTime := utils.GetCacheInt64(GetPlatformMergeKey(mapPlatformId))
				if platform.UpdateState == 2 || platformMergeTime > 0 {
					continue
				}
				SendBackgroundMsgTemplateHandle(MSG_TEMPLATE_SERVER_CLOSE_NUMBER, mapPlatformId+gconv.String(len(mapNodeList)), mapPlatformId+"平台停服节点", "停服节点详情:"+strings.Join(mapNodeList, ";"))

			}
			//NoticeEtsPlatformServerCloseCount()
			//if isMerge == false {
			//	SendSmsAdminUser("SMS_162735700", platformIdStr + strconv.Itoa(num))
			//}
		}
	}
}

// 删除平台没用的数据库
func DelPlatformNotUseData(userId int, PlatformId, DataNameLike string, IsDelState int, serverType int) ([]string, error) {
	var DatabaseNameList []string
	if delPlatformNotUseDataState == 1 {
		return DatabaseNameList, gerror.New("正在删除数据库中")
	}
	initTime := utils.GetTimestamp()
	gameDb, err := GetGameDbByPlatformId(PlatformId)
	utils.CheckError(err)
	if err != nil {
		return DatabaseNameList, err
	}
	defer gameDb.Close()
	//gameServerList, _ := GetAllGameServerDirtyByPlatformId(PlatformId)
	//serverNodeList := GetAllGameServerNodeByPlatformId(PlatformId)
	serverNodeList := GetPlatformIdServerNodeByType(PlatformId, serverType)
	sql := "show databases;"
	//g.Log().Debug("sql:", sql)
	//err := Db.Debug().Raw(sql).Scan(&data).Error
	rows, err := gameDb.DB().Query(sql)
	utils.CheckError(err)
	columns, _ := rows.Columns()

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		for _, col := range values {
			if col != nil {
				DatabaseNameList = append(DatabaseNameList, string(col.([]byte)))
			}
		}
	}
	var NameList []string
	for _, DataName := range DatabaseNameList {
		if strings.Index(DataName, DataNameLike) == -1 {
			continue
		} // 不存在DataNameLike的移除
		IsHaveName := 0
		for _, gameServerNode := range serverNodeList {
			if gameServerNode.DbName == DataName {
				IsHaveName = 1
				break
			} // 如果当前数据库有在使用移除
		}
		if IsHaveName == 1 {
			continue
		}
		NameList = append(NameList, DataName)
	}
	//g.Log().Debug(NameList)
	if IsDelState == 1 {
		delPlatformNotUseDataState = 1
		for _, DelDatabaseName := range NameList {
			delSql := fmt.Sprintf("DROP DATABASE IF EXISTS %s;", DelDatabaseName)
			g.Log().Debug("delSql:", delSql)
			//result := Db.Exec(delSql)
			//result, err := Db.Exec(delSql)
			result, err := gameDb.DB().Exec(delSql)
			utils.CheckError(err)
			delNum, err := result.RowsAffected()
			utils.CheckError(err)
			g.Log().Info("DelDatabaseName:%s ;del num:%d", DelDatabaseName, delNum)
			//g.Log().Info("DelDatabaseName:", DelDatabaseName, ">> ", result)
		}
		delPlatformNotUseDataState = 0
		SendBackgroundMsgTemplateHandleByUser(userId, MSG_TEMPLATE_DEL_NOT_USE_DATA, PlatformId+gconv.String(len(NameList)), PlatformId+"平台删除没用的数据库", "本次删除:"+gconv.String(len(NameList))+"个数据库,用时:"+utils.FormatTimeSecond(utils.GetTimestamp()-initTime)+";详情："+strings.Join(NameList, ";"))
		//SendMailAdminUser(PlatformId+"平台删除没用的数据库", "本次删除:"+strconv.Itoa(len(NameList))+"个数据库,用时:"+utils.FormatTimeSecond(utils.GetTimestamp()-initTime)+";详情："+strings.Join(NameList, ";"))
	}
	return NameList, nil
}
