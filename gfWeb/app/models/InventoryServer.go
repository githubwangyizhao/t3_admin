package models

import (
	"fmt"
	"gfWeb/library/utils"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"path"
	"strings"
	"time"
)

func (a *InventoryServer) TableName() string {
	return InventoryServerTBName()
}

func InventoryServerTBName() string {
	return TableName("inventory_server")
}

type InventoryServerParam struct {
	BaseQueryParam
	Type int
}
type InventoryServer struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	OutIp             string `json:"outIp"`
	InnerIp           string `json:"innerIp"`
	Host              string `json:"host"`
	Type              int    `json:"type"`
	MaxNodeCount      int    `json:"maxNodeCount"`
	NodeCount         int    `json:"nodeCount" gorm:"-"`
	OnlinePlayerCount int    `json:"onlinePlayerCount" gorm:"-"`
	AddTime           int    `json:"addTime"`
	UpdateTime        int    `json:"updateTime"`
}

//获取用户列表
func GetInventoryServerList(params *InventoryServerParam) ([]*InventoryServer, int64) {
	data := make([]*InventoryServer, 0)
	sortOrder := "id"
	switch params.Sort {
	case "id":
		sortOrder = "id"
	}
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	err := Db.Model(&InventoryServer{}).Where(&InventoryServer{
		Type: params.Type,
	}).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		e.NodeCount = GetIpNodeCount(e.InnerIp)
		e.OnlinePlayerCount = GetIpOnlinePlayerCount(e.InnerIp)
	}
	return data, count
}

func GetAllServerListOfGame() []*InventoryServer {
	data := make([]*InventoryServer, 0)
	err := Db.Model(&InventoryServer{}).Where(&InventoryServer{
		Type: 4,
	}).Find(&data).Error
	utils.CheckError(err)
	for _, e := range data {
		e.NodeCount = GetIpNodeCount(e.InnerIp)
		e.OnlinePlayerCount = GetIpOnlinePlayerCount(e.InnerIp)
	}
	return data
}

func GetAllServerList() []*InventoryServer {
	data := make([]*InventoryServer, 0)
	err := Db.Model(&InventoryServer{}).Find(&data).Error
	utils.CheckError(err)
	for _, e := range data {
		e.NodeCount = GetIpNodeCount(e.InnerIp)
		e.OnlinePlayerCount = GetIpOnlinePlayerCount(e.InnerIp)
	}
	return data
}

func GetAllServerListDirty() []*InventoryServer {
	data := make([]*InventoryServer, 0)
	err := Db.Model(&InventoryServer{}).Find(&data).Error
	utils.CheckError(err)
	return data
}

func GetMaxFreeServerByPlatformId(platformId string) (*InventoryServer, error) {
	platform, err := GetPlatformOne(platformId)
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}

	err = Db.Model(&platform).Related(&platform.PlatformInventorySeverRel).Error
	utils.CheckError(err)
	if err != nil {
		return nil, err
	}
	l := make([]*InventoryServer, 0)
	for _, v := range platform.PlatformInventorySeverRel {
		thisInventoryServer, err := GetInventoryServerOne(v.InventoryServerId)
		utils.CheckError(err)
		if err != nil {
			return nil, err
		}
		l = append(l, thisInventoryServer)
	}
	if len(l) == 0 {
		return nil, gerror.New("没有空闲服务器")
	}
	minCount := -1
	nodeCountStr := ""
	inventoryServer := &InventoryServer{}
	for _, e := range l {
		g.Log().Debugf("server:%s, nodeCount:%d, onlineCount:%d ,MaxNodeCount:%d", e.Name, e.NodeCount, e.OnlinePlayerCount, e.MaxNodeCount)
		if e.NodeCount >= e.MaxNodeCount {
			// 一个服务器最多安装50个节点
			continue
		}
		//CalcNodeCountStr(platformId, e.Name, e.OutIp, e.MaxNodeCount, e.NodeCount)
		if e.MaxNodeCount-e.NodeCount <= 3 {
			nodeCountStr += fmt.Sprintf("名称：%s ip:%s 最大节点数:%d 当前节点数:%d \r\n", e.Name, e.OutIp, e.MaxNodeCount, e.NodeCount)
		}
		// 一个节点当作25个在线来计算
		thisValue := e.OnlinePlayerCount + e.NodeCount*25

		if minCount == -1 {
			minCount = thisValue
			inventoryServer = e
		} else {

			if thisValue < minCount {
				minCount = thisValue
				inventoryServer = e
			}
		}
	}
	//go NoticeNodeCountStr(platformId, minCount)
	if minCount == -1 {
		SendBackgroundMsgTemplateHandle(MSG_TEMPLATE_ROBOT_NOT_ENOUGH, platformId, platformId+"机器没有空闲节点", "节点不足平台:"+platformId)
		return nil, gerror.New("没有空闲的服务器")
	}
	if nodeCountStr != "" {
		SendBackgroundMsgTemplateHandle(MSG_TEMPLATE_ROBOT_NOT_ENOUGH_WARNING, platformId, platformId+"机器节点不足-警报", "节点不足详细内容:"+nodeCountStr)
	}
	return inventoryServer, nil
}

func GetMaxFreeServer() (*InventoryServer, error) {
	l := GetAllServerListOfGame()
	if len(l) == 0 {
		return nil, gerror.New("没有空闲服务器")
	}
	minCount := -1
	inventoryServer := &InventoryServer{}
	for _, e := range l {
		g.Log().Debug("server:%s, nodeCount:%d, onlineCount:%d", e.Name, e.NodeCount, e.OnlinePlayerCount)
		if e.NodeCount >= 33 {
			// 一个服务器最多安装33个节点
			continue
		}
		// 一个节点当作10个在线来计算
		thisValue := e.OnlinePlayerCount + e.NodeCount*10

		if minCount == -1 {
			minCount = thisValue
			inventoryServer = e
		} else {

			if thisValue < minCount {
				minCount = thisValue
				inventoryServer = e
			}
		}
	}
	if minCount == -1 {
		return nil, gerror.New("没有空闲的服务器")
	}
	return inventoryServer, nil
}

// 获取单个服务器
func GetInventoryServerOneDirty(id int) (*InventoryServer, error) {
	r := &InventoryServer{
		Id: id,
	}
	err := Db.First(&r).Error
	return r, err
}

// 获取单个服务器
func GetInventoryServerOne(id int) (*InventoryServer, error) {
	r := &InventoryServer{
		Id: id,
	}
	err := Db.First(&r).Error
	r.NodeCount = GetIpNodeCount(r.InnerIp)
	r.OnlinePlayerCount = GetIpOnlinePlayerCount(r.InnerIp)
	return r, err
}

// 删除服务器列表
func DeleteInventoryServers(ids []int) error {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := Db.Where(ids).Delete(&InventoryServer{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if _, err := DeletePlatformInventorySeverRelByInventoryServerIdList(ids); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func CreateAnsibleInventory() error {
	g.Log().Info("开始创建ansible inventory 文件.....")
	ansibleInventoryFile := g.Cfg().GetString("ansible.ansible_inventory_file")
	if ansibleInventoryFile == "" {
		g.Log().Error("读取配置ansible_inventory_file失败")
		return gerror.New("读取配置ansible_inventory_file失败")
	}
	//mapList := make(map[string][]string, 0)
	//serverNodeList := GetAllServerNodeList()
	//for _, e := range serverNodeList {
	//	array := strings.Split(e.Node, "@")
	//	if len(array) != 2 {
	//		return gerror.New("解析节点名字失败:" + e.Node)
	//	}
	//	nodeName := array[0]
	//	nodeIp := array[1]
	//	//g.Log().Info("nodeName:%+v", nodeName)
	//	if v, ok := mapList[nodeIp]; ok {
	//		v = append(v, "'"+nodeName+"'")
	//		mapList[nodeIp] = v
	//	} else {
	//		mapList[nodeIp] = append(make([] string, 0), "'"+nodeName+"'")
	//	}
	//}
	serverList := GetAllServerListDirty()
	//for _, e := range serverOfGameList {
	//	if _, ok := mapList[e.InnerIp]; ok {
	//	} else {
	//		g.Log().Info("e:%+v", e)
	//		mapList[e.InnerIp] = make([] string, 0)
	//	}
	//}
	//g.Log().Info("serverOfGameList:%+v", mapList)
	//g.Log().Info("mapList:%+v", mapList)
	content := "## Generated automatically, no need to modify.\n"
	content += fmt.Sprintf("## Auto Created :%s\n\n", time.Now().String())
	//for ip, _ := range mapList {
	//	content += fmt.Sprintf("%s\n", ip)
	//}
	//content += "\n\n\n"

	for _, server := range serverList {
		content += fmt.Sprintf("%s    \n\n", server.InnerIp)
		//content += "nodes="
		//content += "\"[" + strings.Join(nodes, ", ") + "]\""
		//content += "\n\n"

	}

	err := utils.FilePutContext(ansibleInventoryFile, content)

	if err != nil {
		return err
	}
	g.Log().Info("创建ansible inventory 文件(%s)成功.", ansibleInventoryFile)

	err = CreateNodes()
	utils.CheckError(err, "生成ansible inventory2 失败")

	return nil
}

//func CreateNginxConf(platformId string) error {
//	g.Log().Info("开始创建nginx config 文件.....")
//	mapList := make(map[string][]string, 0)
//	serverNodeList := GetAllGameServerNodeByPlatformId(platformId)
//	for _, e := range serverNodeList {
//		array := strings.Split(e.Node, "@")
//		if len(array) != 2 {
//			return gerror.New("解析节点名字失败:" + e.Node)
//		}
//		nodeName := array[0]
//		nodeIp := array[1]
//		//g.Log().Info("nodeName:%+v", nodeName)
//		if v, ok := mapList[nodeIp]; ok {
//			v = append(v, "'"+nodeName+"'")
//			mapList[nodeIp] = v
//		} else {
//			mapList[nodeIp] = append(make([] string, 0), "'"+nodeName+"'")
//		}
//	}
//	serverList := GetAllServerList()
//	//for _, e := range serverOfGameList {
//	//	if _, ok := mapList[e.InnerIp]; ok {
//	//	} else {
//	//		g.Log().Info("e:%+v", e)
//	//		mapList[e.InnerIp] = make([] string, 0)
//	//	}
//	//}
//	//g.Log().Info("serverOfGameList:%+v", mapList)
//	//g.Log().Info("mapList:%+v", mapList)
//	content := "## Generated automatically, no need to modify.\n"
//	content += fmt.Sprintf("## Auto Created :%s\n\n", time.Now().String())
//	//for ip, _ := range mapList {
//	//	content += fmt.Sprintf("%s\n", ip)
//	//}
//	//content += "\n\n\n"
//
//	for _, server := range serverList {
//		content += fmt.Sprintf("%s    \n\n", server.InnerIp)
//		//content += "nodes="
//		//content += "\"[" + strings.Join(nodes, ", ") + "]\""
//		//content += "\n\n"
//
//	}
//
//	err := utils.FilePutContext(ansibleInventoryFile, content)
//
//	if err != nil {
//		return err
//	}
//	g.Log().Info("创建ansible inventory 文件(%s)成功.", ansibleInventoryFile)
//
//	err = CreateNodes()
//	utils.CheckError(err, "生成ansible inventory2 失败")
//
//	return nil
//}

type nodeType struct {
	typeId             int
	name               string
	isDivisionPlatform bool
}

func CreateNodes() error {
	g.Log().Info("开始创建nodes文件.....")
	ansibleInventoryDir := g.Cfg().GetString("ansible.ansible_nodes_dir")
	if ansibleInventoryDir == "" {
		g.Log().Error("ansible_nodes_dir")
		return gerror.New("读取配置ansible_nodes_dir失败")
	}

	serverNodeList := GetAllServerNodeList()
	typeList := []nodeType{
		//nodeType{
		//	typeId:             0,
		//	name:               "center",
		//	isDivisionPlatform: false,
		//},
		nodeType{
			typeId:             1,
			name:               "game",
			isDivisionPlatform: true,
		},
		nodeType{
			typeId:             2,
			name:               "zone",
			isDivisionPlatform: true,
		},
		nodeType{
			typeId:             4,
			name:               "login_server",
			isDivisionPlatform: false,
		},
		nodeType{
			typeId:             5,
			name:               "unique_id",
			isDivisionPlatform: false,
		},
		nodeType{
			typeId:             6,
			name:               "charge",
			isDivisionPlatform: false,
		},
		nodeType{
			typeId:             7,
			name:               "war",
			isDivisionPlatform: true,
		},
	}
	platformMapList := make(map[string]bool, 0)
	platformList := make([]string, 0)
	//g.Log().Info("platformList0:%+v", platformList)
	for _, e := range serverNodeList {
		//g.Log().Info("a:%++v, %++v, %++v", e.PlatformId, e.PlatformId == "", len(e.PlatformId))
		if e.PlatformId == "" {
			continue
		}

		if _, ok := platformMapList[e.PlatformId]; ok {
		} else {
			platformMapList[e.PlatformId] = true
			//g.Log().Info("add:%s, %+v", e.PlatformId, e.PlatformId == "")
			platformList = append(platformList, e.PlatformId)
		}
	}

	g.Log().Info("platformList:%+v", platformList)

	for _, t := range typeList {
		var ansibleInventoryFile string
		if t.isDivisionPlatform == true {
			for _, platform := range platformList {
				ansibleInventoryFile = path.Join(ansibleInventoryDir) + "/" + t.name + "_" + platform
				err := doCreateNode(ansibleInventoryFile, t.typeId, serverNodeList, platform)
				utils.CheckError(err)
			}
		} else {
			ansibleInventoryFile = path.Join(ansibleInventoryDir) + "/" + t.name
			err := doCreateNode(ansibleInventoryFile, t.typeId, serverNodeList, "")
			utils.CheckError(err)
		}
	}

	return nil
}

func doCreateNode(ansibleInventoryFile string, typeId int, serverNodeList []*ServerNode, platformId string) error {
	//g.Log().Info("doCreateNode:%s", ansibleInventoryFile)
	mapList := make(map[string][]string, 0)
	for _, e := range serverNodeList {
		if e.Type != typeId {
			continue
		}
		if platformId != "" {
			if e.PlatformId == "" {
				g.Log().Warning("节点平台未配置:%d, %+s", typeId, e)
				continue
			}
			if e.PlatformId != platformId {
				continue
			}
		}
		array := strings.Split(e.Node, "@")
		if len(array) != 2 {
			return gerror.New("解析节点名字失败:" + e.Node)
		}
		nodeName := array[0]
		nodeIp := array[1]
		//g.Log().Info("nodeName:%+v", nodeName)
		if v, ok := mapList[nodeIp]; ok {
			v = append(v, "'"+nodeName+"'")
			mapList[nodeIp] = v
		} else {
			mapList[nodeIp] = append(make([]string, 0), "'"+nodeName+"'")
		}
	}
	//serverOfGameList := GetAllServerList()
	//for _, e := range serverOfGameList {
	//	if _, ok := mapList[e.InnerIp]; ok {
	//	} else {
	//		//g.Log().Info("e:%+v", e)
	//		mapList[e.InnerIp] = make([] string, 0)
	//	}
	//}
	//g.Log().Info("serverOfGameList:%+v", mapList)
	//g.Log().Info("mapList:%+v", mapList)
	content := "## Generated automatically, no need to modify.\n"
	content += fmt.Sprintf("## Auto Created :%s\n\n", time.Now().String())
	//for ip, _ := range mapList {
	//	content += fmt.Sprintf("%s\n", ip)
	//}
	//content += "\n\n\n"

	for ip, nodes := range mapList {
		content += fmt.Sprintf("%s    ", ip)
		content += "nodes="
		content += "\"[" + strings.Join(nodes, ", ") + "]\""
		content += "\n\n"

	}

	err := utils.FilePutContext(ansibleInventoryFile, content)

	if err != nil {
		return err
	}
	return nil
}
