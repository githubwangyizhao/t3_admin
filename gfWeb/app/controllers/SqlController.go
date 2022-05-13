package controllers

import (
	"encoding/json"
	"gfWeb/app/models"
	"gfWeb/library/enums"
	"strings"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type SqlController struct {
	BaseController
}

func (c *SqlController) SqlQuery(r *ghttp.Request) {
	type col struct {
		Label string `json:"label"`
		Prop  string `json:"prop"`
	}
	var params struct {
		Sql        string `json:"sql"`
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	nodeStrList := make([]string, 0)
	if len(params.ServerId) == 0 {
		nodeStrList = models.GetAllGameNodeByPlatformId(params.PlatformId)
	} else {
		gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
		c.CheckError(err)
		nodeStrList = append(nodeStrList, gameServer.Node)
	}
	c.nodeSqlQuery(nodeStrList, params.Sql)
	//gameDb, err := models.GetVisitorGameDbByNode(gameServer.Node)
	//c.CheckError(err)
	//defer gameDb.Close()
	//rows, err := gameDb.DB().Query(params.Sql)
	//c.CheckError(err)
	//columns, _ := rows.Columns()
	//scanArgs := make([]interface{}, len(columns))
	//values := make([]interface{}, len(columns))
	//for i := range values {
	//	scanArgs[i] = &values[i]
	//}
	//records := make([]map[string]string, 0)
	//cols := make([]*col, 0)
	//for _, e := range columns {
	//	cols = append(cols, &col{
	//		Label: e,
	//		Prop:  e,
	//	})
	//}
	//for rows.Next() {
	//	//将行数据保存到record字典
	//	err = rows.Scan(scanArgs...)
	//	record := make(map[string]string, 0)
	//	for i, col := range values {
	//		if col != nil {
	//			//record = append(record, string(col.([]byte)))
	//			record[columns[i]] = string(col.([]byte))
	//		}
	//	}
	//	records = append(records, record)
	//	//fmt.Println(record)
	//}
	//if len(records) > 200 {
	//	c.CheckError(errors.New("返回结果超过200行， 请加限制语句！"))
	//}
	//result := make(map[string]interface{})
	//result["cols"] = cols
	//result["rows"] = records
	//c.HttpResult(r, enums.CodeSuccess, "成功!", result)
}

func (c *SqlController) nodeSqlQuery(nodeStrList []string, sql string) {
	type col struct {
		Label string `json:"label"`
		Prop  string `json:"prop"`
	}
	cols := make([]*col, 0)
	records := make([]map[string]string, 0)
	g.Log().Debugf("nodeStrList:%+v", nodeStrList)
	for _, nodeStr := range nodeStrList {
		gameDb, err := models.GetVisitorGameDbByNode(nodeStr)
		if err != nil {
			g.Log().Warningf("查询sql节点数据失败：%+v", nodeStr)
			continue
		}
		//c.CheckError(err)
		//defer gameDb.Close()
		rows, err := gameDb.DB().Query(sql)
		if err != nil {
			gameDb.Close()
			g.Log().Warningf("查询sql连接数据库失败：%+v  %s", nodeStr, sql)
			continue
		}
		//c.CheckError(err)
		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		if len(cols) == 0 {
			for _, e := range columns {
				cols = append(cols, &col{
					Label: e,
					Prop:  e,
				})
			}
		}
		for rows.Next() {
			//将行数据保存到record字典
			err = rows.Scan(scanArgs...)
			record := make(map[string]string, 0)
			for i, col := range values {
				if col != nil {
					//record = append(record, string(col.([]byte)))
					record[columns[i]] = string(col.([]byte))
				}
			}
			records = append(records, record)
			//fmt.Println(record)
		}
		gameDb.Close()
		if len(records) > 20000 {
			c.CheckError(gerror.New("返回结果超过20000行， 请加限制语句！"))
		}
	}
	result := make(map[string]interface{})
	result["total"] = len(records)
	result["cols"] = cols
	result["rows"] = records
	c.HttpResult(c.Request, enums.CodeSuccess, "成功!", result)
}

// 删除平台没有使用的数据库
func (c *SqlController) DelPlatformNotDatabase(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:"platformId"`
		DataName   string `json:"dataName"`
		DelState   int    `json:"delState"`
	}
	err := json.Unmarshal(r.GetBody(), &params)

	c.CheckError(err)
	g.Log().Infof("删除平台没用的数据库:%+v", params)
	if err != nil {
		g.Log().Error("删除平台没有使用的数据库失败: %+v", err)
		c.HttpResult(r, enums.CodeFail, "删除平台没有使用的数据库解析数据失败", "")
	}
	DataName := params.DataName
	PlatformId := params.PlatformId
	if len(PlatformId) == 0 {
		c.HttpResult(r, enums.CodeFail, "删除平台数据库没有选择平台", "")
	}
	//if strings.Index( DataName, "db_") == -1 || strings.Index( DataName, "_game_") == -1  {
	if strings.Index(DataName, "db_") == -1 || strings.Index(DataName, "_game_") == strings.Index(DataName, "_zone_") || strings.Index(DataName, params.PlatformId) == -1 {
		c.HttpResult(r, enums.CodeFail, "删除平台数据库参数格式错误 db_xx_game_:"+params.PlatformId+" 参数内容: "+params.DataName, "")
	}
	serverType := 1
	if strings.Index(DataName, "_game_") == -1 {
		serverType = 2
	}
	DelDataNameList, err := models.DelPlatformNotUseData(c.curUser.Id, params.PlatformId, DataName, params.DelState, serverType)
	if err != nil {
		c.HttpResult(r, enums.CodeFail, err.Error(), "")
	}
	Str := ""
	if len(DelDataNameList) == 0 {
		Str = "没有可以删除的数据库"
	} else {
		for i, DelDataName := range DelDataNameList {
			Str = Str + " | " + DelDataName
			if (i+1)%5 == 0 {
				Str += "\r\n"
			}
		}
	}
	c.HttpResult(r, enums.CodeSuccess, Str, "")
}

// 获得数据名字
func (c *SqlController) GetDatabaseName(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	serverNode, err := models.GetServerNode(gameServer.Node)
	c.CheckError(err)
	result := make(map[string]interface{})
	result["rows"] = serverNode.DbName
	c.HttpResult(r, enums.CodeSuccess, "成功", result)
}

// 打包数据库
func (c *SqlController) PackDatabase(r *ghttp.Request) {
	var params struct {
		PlatformId string `json:"platformId"`
		ServerId   string `json:"serverId"`
	}
	err := json.Unmarshal(r.GetBody(), &params)
	c.CheckError(err)
	gameServer, err := models.GetGameServerOne(params.PlatformId, params.ServerId)
	c.CheckError(err)
	serverNode, err := models.GetServerNode(gameServer.Node)
	c.CheckError(err)
	zipPathFileName, err := models.BackDatabaseBaseAndZip("root", serverNode.DbHost, "", serverNode.DbName)
	c.CheckError(err)
	result := make(map[string]interface{})
	result["res_url"] = g.Cfg().GetString("database.sqlDownloadPath", "/mysql_back") + "/" + zipPathFileName
	c.HttpResult(r, enums.CodeSuccess, "成功", result)
}
