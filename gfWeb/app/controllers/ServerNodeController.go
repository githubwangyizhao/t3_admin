package controllers

import (
	"encoding/json"
	"fmt"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	//"os"
)

type ServerNodeController struct {
	BaseController
}

//获取节点列表
func (c *ServerNodeController) List(r *ghttp.Request) {
	var params models.ServerNodeQueryParam
	err := json.Unmarshal(r.GetBody(), &params)
	utils.CheckError(err)
	data, total := models.ServerNodePageList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.HttpResult(r, enums.CodeSuccess, "获取节点列表成功", result)
}

//  添加 编辑 节点
func (c *ServerNodeController) Edit(r *ghttp.Request) {
	params := models.ServerNode{}
	err := json.Unmarshal(r.GetBody(), &params)

	m := params
	g.Log().Debug("编辑 节点:%v", m)
	utils.CheckError(err, "编辑节点")
	if m.IsAdd == 1 && models.IsServerNodeExists(m.Node) {
		c.HttpResult(r, enums.CodeFail, "节点已经存在", m.Node)
	}
	if m.IsAdd == 0 && models.IsServerNodeExists(m.Node) == false {
		c.HttpResult(r, enums.CodeFail, "节点不存在", m.Node)
	}

	out, err := models.AddServerNode(m.Node, m.Ip, m.Port, m.WebPort, m.Type, m.PlatformId, m.DbHost, m.DbPort, m.DbName)

	//out, err := utils.NodeTool(
	//	"mod_server_mgr",
	//	"add_server_node",
	//	m.Node,
	//	m.Ip,
	//	strconv.Itoa(m.Port),
	//	strconv.Itoa(m.WebPort),
	//	strconv.Itoa(m.Type),
	//	m.PlatformId,
	//	m.DbHost,
	//	strconv.Itoa(m.DbPort),
	//	m.DbName,
	//)
	c.CheckError(err, "保存节点失败:"+out)
	c.HttpResult(r, enums.CodeSuccess, "保存成功", m.Node)
}

// 删除节点
func (c *ServerNodeController) Delete(r *ghttp.Request) {
	var params []string
	err := json.Unmarshal(r.GetBody(), &params)

	ids := params
	utils.CheckError(err)
	g.Log().Infof("删除节点:%+v", ids)
	for _, str := range ids {
		out, err := utils.CenterNodeTool(
			"mod_server_mgr",
			"delete_server_node",
			str,
		)
		c.CheckError(err, "删除节点失败:"+out)
	}
	c.HttpResult(r, enums.CodeSuccess, fmt.Sprintf("成功删除 %d 项", len(ids)), 0)
}

////  启动 节点
//func (c *ServerNodeController) Start() {
//	var params struct {
//		Node     string
//	}
//	err := json.Unmarshal(r.GetBody(),&params)
//
//	g.Log().Debug("启动 节点:%v",params )
//	c.CheckError(err, "启动节点")
//	curDir := utils.GetCurrentDirectory()
//	defer os.Chdir(curDir)
//	toolDir := utils.GetToolDir()
//	err = os.Chdir(toolDir)
//
//	c.CheckError(err, "启动节点")
//	commandArgs := []string{
//		"node_tool.sh",
//		params.Node,
//		"start",
//	}
//	out, err := utils.Cmd("sh", commandArgs)
//	c.CheckError(err, "启动节点失败:"+ out)
//	c.HttpResult(r, enums.CodeSuccess, "启动成功", params.Node)
//}
//
////  停止 节点
//func (c *ServerNodeController) Stop() {
//	var params struct {
//		Node     string
//	}
//	err := json.Unmarshal(r.GetBody(),&params)
//
//	g.Log().Debug("停止节点:%v",params )
//	c.CheckError(err, "停止节点")
//	curDir := utils.GetCurrentDirectory()
//	defer os.Chdir(curDir)
//	toolDir := utils.GetToolDir()
//	err = os.Chdir(toolDir)
//	c.CheckError(err, "停止节点")
//	commandArgs := []string{
//		"node_tool.sh",
//		params.Node,
//		"stop",
//	}
//	out, err := utils.Cmd("sh", commandArgs)
//	c.CheckError(err, "停止节点失败:"+ out)
//	c.HttpResult(r, enums.CodeSuccess, "停止成功", params.Node)
//}

func (c *ServerNodeController) Action(r *ghttp.Request) {
	var params struct {
		Nodes  []string `json:"nodes"`
		Action string
	}

	err := json.Unmarshal(r.GetBody(), &params)

	g.Log().Debugf("节点操作Action:%v", params)
	c.CheckError(err)
	if len(params.Nodes) == 1 {
		nodeStr := params.Nodes[0]
		serverNode, err := models.GetServerNode(nodeStr)
		c.CheckError(err)
		switch serverNode.Type {
		case 1, 2, 7:
			platform := models.GetPlatformSimpleOne(serverNode.PlatformId)
			err = models.NodeActionHandle(params.Nodes, params.Action, platform.Version)
			c.CheckError(err)
		default:
			err = models.NodeAction(params.Nodes, params.Action)
			c.CheckError(err)
		}
	} else {
		err = models.NodeAction(params.Nodes, params.Action)
		c.CheckError(err)
	}

	//curDir := utils
	// .GetCurrentDirectory()
	//defer os.Chdir(curDir)
	//toolDir := utils.GetToolDir()
	//err = os.Chdir(toolDir)
	//c.CheckError(err)
	//var commandArgs []string
	//for _, node := range params.Nodes {
	//	switch params.Action {
	//	case "start":
	//		commandArgs = []string{"node_tool.sh", node, params.Action,}
	//	case "stop":
	//		commandArgs = []string{"node_tool.sh", node, params.Action,}
	//	case "hot_reload":
	//		commandArgs = []string{"node_hot_reload.sh", node, "server",}
	//	case "cold_reload":
	//		commandArgs = []string{"node_cold_reload.sh", node, "server",}
	//	}
	//	out, err := utils.Cmd("sh", commandArgs)
	//	c.CheckError(err, fmt.Sprintf("操作节点失败:%v %v", params, out))
	//}

	c.HttpResult(r, enums.CodeSuccess, "操作节点成功", "")
}

func (c *ServerNodeController) Install(r *ghttp.Request) {
	var params struct {
		Node string `json:"node"`
	}

	err := json.Unmarshal(r.GetBody(), &params)

	g.Log().Infof("部署节点:%v", params)
	c.CheckError(err)
	//curDir := utils.GetCurrentDirectory()
	//defer os.Chdir(curDir)
	//toolDir := utils.GetToolDir()
	//err = os.Chdir(toolDir)
	//c.CheckError(err)
	//var commandArgs []string
	//serverNode, err := models.GetServerNode(params.Node)
	//c.CheckError(err, "节点不存在:"+params.Node)
	//app := ""
	//switch serverNode.Type {
	//case 0:
	//	app = "center"
	//case 1:
	//	app = "game"
	//case 2:
	//	app = "zone"
	//case 4:
	//	app = "login_server"
	//case 5:
	//	app = "unique_id"
	//case 6:
	//	app = "charge"
	//}
	//
	//commandArgs = []string{"/data/tool/ansible/do-install.sh", serverNode.Node, app,serverNode.DbName, serverNode.DbHost, strconv.Itoa(serverNode.DbPort), "root"}
	//out, err := utils.Cmd("sh", commandArgs)
	//c.CheckError(err, fmt.Sprintf("操作节点失败:%v %v", params, out))
	//g.Log().Info("部署节点成功:%v", params)
	err = models.InstallNode(params.Node)
	c.CheckError(err)
	c.HttpResult(r, enums.CodeSuccess, "部署节点成功", params.Node)
}
