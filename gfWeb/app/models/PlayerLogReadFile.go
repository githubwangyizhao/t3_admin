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

// 游戏服日志文件读取
func GameLogFileRead(platformId, serverId, fileName, grepParam, regStr string, datetime int) (err error, matchArray [][]string) {
	gameServer, err := GetGameServerOne(platformId, serverId)
	utils.CheckError(err)
	if err != nil {
		return nil, matchArray
	}
	node := gameServer.Node
	serverNode, err := GetServerNode(node)
	utils.CheckError(err)
	if err != nil {
		return nil, matchArray
	}
	t := time.Unix(int64(datetime), 0)
	logDir := fmt.Sprintf("%d_%d_%d", t.Year(), t.Month(), t.Day())
	g.Log().Debug("logDir:%s", logDir)

	sshKey := g.Cfg().GetString("game.ssh_key")
	sshPort := g.Cfg().GetString("game.ssh_port")
	nodeName := strings.Split(serverNode.Node, "@")[0]
	nodeIp := strings.Split(serverNode.Node, "@")[1]
	cmd := fmt.Sprintf("ssh -i %s -p%s %s ' /usr/bin/cat /data/log/game/%s/%s/%s %s'", sshKey, sshPort, nodeIp, nodeName, logDir, fileName, grepParam)

	out, err := utils.ExecShell(cmd)
	if err != nil {
		if len(out) == 0 {
			return nil, matchArray
		}
		return err, matchArray
	}
	//utils.CheckError(err)
	reg := regexp.MustCompile(regStr)
	matchArray = reg.FindAllStringSubmatch(out, -1)
	return nil, matchArray
}

// 游戏服日志文件道具字段读取
func GameLogFileItemList(str string) (itemList []Prop) {
	g.Log().Debugf("道具字段读取:%s", str)
	itemListStr := strings.Replace(str, "\"", "", -1)
	if itemListStr != "[]" {
		itemListStr1 := strings.Replace(strings.Replace(itemListStr, "[{", "", -1), "}]", "", -1)
		itemListStr1 = strings.Replace(strings.Replace(itemListStr1, "[[", "", -1), "]]", "", -1)
		itemListStr1 = strings.Replace(strings.Replace(itemListStr1, "},{", "|", -1), "],[", "|", -1)
		line := strings.Split(itemListStr1, "|")
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
	return itemList
}
