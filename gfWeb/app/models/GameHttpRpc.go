package models

import (
	"encoding/json"
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
)

// 向游戏服请求rpc操作 平台，区服
func GameHttpRpcPlatform(platformId, serverId, mod, fun, args string) ([]string, string, error) {
	realNodeList := make([]string, 0)
	if serverId == "" {
		realNodeList = GetAllGameNodeByPlatformId(platformId)
	} else {
		gameServer, err := GetGameServerOne(platformId, serverId)
		utils.CheckError(err)
		if err != nil {
			//c.HttpResult(r, enums.CodeFail, "未找到服务区服", 0)
			return realNodeList, "未找到服务区服", err
		}
		realNodeList = append(realNodeList, gameServer.Node)
	}
	ResultList, resultMsg, err := GameHttpRpc(mod, fun, args, realNodeList)
	return ResultList, resultMsg, err
}

// 向游戏服请求rpc操作 节点
func GameHttpRpcNode(nodeStr, mod, fun, args string) ([]string, string, error) {
	realNodeList := make([]string, 0)
	realNodeList = append(realNodeList, nodeStr)
	ResultList, resultMsg, err := GameHttpRpc(mod, fun, args, realNodeList)
	if err != nil {
		return ResultList, "节点操作失败", err
	}
	return ResultList, resultMsg, err
}

// 向游戏服请求rpc操作
func GameHttpRpc(mod, fun, args string, realNodeList []string) ([]string, string, error) {
	var httpParams struct {
		Mod  string `json:"mod"`
		Fun  string `json:"fun"`
		Args string `json:"args"` // [a,b]
	}
	httpParams.Mod = mod
	httpParams.Fun = fun
	httpParams.Args = args
	request, err := json.Marshal(httpParams)
	g.Log().Infof("向游戏服请求rpc操作request:%+v", string(request))
	resultMsgStr := ""
	resultList := make([]string, 0)
	if err != nil {
		return resultList, resultMsgStr, err
	}
	Err := err
	for _, node := range realNodeList {
		url := GetGameURLByNode(node) + "/game_rpc"
		data := string(request)
		resultMsg, err := utils.HttpRequest(url, data)
		if err != nil {
			g.Log().Errorf("向游戏服请求rpc操作失败 node:%s, mod:%s, fun:%s args:%s  err:%s", node, mod, fun, args, err)
			resultList = append(resultList, node)
			Err = err
		}
		resultMsgStr += resultMsg
	}
	return resultList, resultMsgStr, Err
}
