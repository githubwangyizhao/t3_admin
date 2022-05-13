package models

import (
	"encoding/json"
	"fmt"
	"gfWeb/library/enums"
	"gfWeb/library/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"math"
	"strconv"
	"strings"
	"time"
	"xorm.io/core"
)

type db struct {
	DbHost string       `json:"db_host"`
	DbName string       `json:"db_name"`
	DbUser string       `json:"db_user"`
	DbPort int          `json:"db_port"`
	DbPwd  string       `json:"db_pwd"`
	Db     *xorm.Engine `json:"-"`
}

type dbConfig struct {
	TargetDb         *db   `json:"target_db"`
	SourceDbList     []*db `json:"source_db_list"`
	CleanLevel       int   `json:"clean_level"`
	CleanVipLevel    int   `json:"clean_vip_level"`
	CleanNoLoginDay  int   `json:"clean_no_login_day"`
	CleanLevel2      int   `json:"clean_level2"`
	CleanVipLevel2   int   `json:"clean_vip_level2"`
	CleanNoLoginDay2 int   `json:"clean_no_login_day2"`
}

type tableConfig struct {
	IgnoreList        []string        `json:"ignore_list"`
	CleanList         []string        `json:"clean_list"`
	ForeignKeyMapList []foreignKeyMap `json:"foreign_key_map_list"`
}

type foreignKeyMap struct {
	Table        string `json:"table"`
	Filed        string `json:"filed"`
	ForeignTable string `json:"foreign_table"`
	ForeignKey   string `json:"foreign_key"`
}

// 合服节点数据
type mergeNode struct {
	nodeList []*ServerNode
	zoneNode string
	s        int
	e        int
}

func startMergeHandle(sourceServerNodeList []*ServerNode, targetServerNode *ServerNode, zoneNode string) error {
	t0 := time.Now()
	//now := utils.GetTimestamp()
	allServerNodeList := make([]*ServerNode, 0)
	allServerNodeList = append(allServerNodeList, sourceServerNodeList...)
	allServerNodeList = append(allServerNodeList, targetServerNode)

	for _, e := range sourceServerNodeList {
		g.Log().Infof("合服源:%+v", e.Node)
	}
	g.Log().Infof("合服目标:%+v", targetServerNode.Node)

	if len(sourceServerNodeList) == 0 {
		return gerror.New("源节点不能为空")
	}

	zoneNodeList := make([]string, 0)

	for _, e := range allServerNodeList {
		if e.ZoneNode != "" {
			if !utils.IsHaveArray(e.ZoneNode, zoneNodeList) {
				zoneNodeList = append(zoneNodeList, e.ZoneNode)
			}
		}
	}
	if !utils.IsHaveArray(zoneNode, zoneNodeList) {
		zoneNodeList = append(zoneNodeList, zoneNode)
	}
	g.Log().Infof("跨服节点列表:%+v", zoneNodeList)

	//1.修改服务器状态
	g.Log().Info("[1].修改服务器状态...")
	for _, e := range allServerNodeList {
		out, err := utils.CenterNodeTool(
			"mod_server_mgr",
			"update_node_state",
			e.Node,
			strconv.Itoa(enums.ServerStateMaintenance),
		)
		utils.CheckError(err, "修改区服状态:"+out)
		if err != nil {
			return err
		}
	}

	//2.刷新入口
	g.Log().Info("[2].刷新入口...")
	err := RefreshGameServer()
	utils.CheckError(err)
	if err != nil {
		return err
	}

	//3.处理跨服数据
	g.Log().Info("[3].处理跨服数据...")
	for _, e := range zoneNodeList {
		out, err := utils.NodeTool(e, "fairyland_srv", "gm_settle_award")
		utils.CheckError(err, "处理跨服数据失败"+out)
	}
	//4.关闭节点
	g.Log().Info("[4].关闭节点...")
	for _, e := range allServerNodeList {
		err = NodeAction([]string{e.Node}, "stop")
		utils.CheckError(err)
		if err != nil {
			return err
		}
	}
	for _, e := range zoneNodeList {
		err = NodeAction([]string{e}, "stop")
		utils.CheckError(err)
		if err != nil {
			return err
		}
	}

	//5. 赋值 db_config
	g.Log().Info("[5].赋值 db_config...")
	gameDbPwd := g.Cfg().GetString("database.game.pass")
	cleanLevel := g.Cfg().GetInt("merge.cleanLevel", 200)
	cleanVipLevel := g.Cfg().GetInt("merge.cleanVipLevel", 0)
	cleanNoLoginDay := g.Cfg().GetInt("merge.cleanNoLoginDay", 7)
	cleanLevel2 := g.Cfg().GetInt("merge.cleanLevel2", 0)
	cleanVipLevel2 := g.Cfg().GetInt("merge.cleanVipLevel2", 0)
	cleanNoLoginDay2 := g.Cfg().GetInt("merge.cleanNoLoginDay2", 60)
	dbConfig := &dbConfig{
		CleanVipLevel:    cleanVipLevel,
		CleanLevel:       cleanLevel,
		CleanNoLoginDay:  cleanNoLoginDay,
		CleanVipLevel2:   cleanVipLevel2,
		CleanLevel2:      cleanLevel2,
		CleanNoLoginDay2: cleanNoLoginDay2,
	}
	dbUserName := g.Cfg().GetString("database.default.user")
	dbConfig.TargetDb = &db{
		DbHost: targetServerNode.DbHost,
		DbName: targetServerNode.DbName,
		DbUser: dbUserName,
		DbPort: targetServerNode.DbPort,
		DbPwd:  gameDbPwd,
	}
	for _, e := range sourceServerNodeList {
		serverNode, err := GetServerNode(e.Node)
		utils.CheckError(err)
		if err != nil {
			return err
		}
		db := &db{
			DbHost: serverNode.DbHost,
			DbName: serverNode.DbName,
			//DbUser: "root",
			DbUser: dbUserName,
			DbPort: serverNode.DbPort,
			DbPwd:  gameDbPwd,
		}
		dbConfig.SourceDbList = append(dbConfig.SourceDbList, db)
	}

	//6.备份数据库
	g.Log().Info("[6].备份数据库...")
	dbUserName = g.Cfg().GetString("database.default.user")
	_, err = BackDatabaseBase(dbUserName, dbConfig.TargetDb.DbHost, dbConfig.TargetDb.DbPwd, dbConfig.TargetDb.DbName, "")
	if err != nil {
		return err
	}
	for _, e := range dbConfig.SourceDbList {
		_, err = BackDatabaseBase(dbUserName, e.DbHost, e.DbPwd, e.DbName, "")
		if err != nil {
			return err
		}
	}

	//7.合并数据库
	g.Log().Info("[7].合并数据库...")
	err = mergeDb(dbConfig)
	utils.CheckError(err)
	if err != nil {
		return err
	}

	//8.修改区服节点映射
	g.Log().Info("[8].修改区服节点映射...")

	for _, e := range allServerNodeList {
		gameServerList := GetGameServerByNode(e.Node)
		for _, g := range gameServerList {
			out, err := AddGameServer(g.PlatformId, g.Sid, g.Desc, targetServerNode.Node, zoneNode, targetServerNode.State, targetServerNode.OpenTime, g.IsShow)
			utils.CheckError(err, out)
			if err != nil {
				return err
			}
		}
	}

	//9.删除没用的节点
	g.Log().Info("[9].删除没用的节点...")

	for _, e := range sourceServerNodeList {
		out, err := utils.CenterNodeTool(
			"mod_server_mgr",
			"delete_server_node",
			e.Node,
		)
		utils.CheckError(err, "删除游戏节点失败:"+out)
		if err != nil {
			return err
		}
	}

	for _, e := range zoneNodeList {
		if e != zoneNode && GetZoneConnectNodeCount(e) == 0 {
			_, err := utils.CenterNodeTool(
				"mod_server_mgr",
				"delete_server_node",
				e,
			)
			utils.CheckError(err, "删除跨服节点失败:"+e)
		} else {
			err = NodeAction([]string{e}, "start")
			utils.CheckError(err)
		}
	}
	//10.启动节点
	g.Log().Info("[10].启动节点...")
	err = NodeAction([]string{targetServerNode.Node}, "start")
	utils.CheckError(err)
	if err != nil {
		return err
	}
	//err = NodeAction([] string{zoneNode}, "start")
	//utils.CheckError(err)
	//if err != nil {
	//	return err
	//}

	//11.同步节点信息
	g.Log().Info("[11].同步节点信息...")
	err = AfterAddGameServer()
	utils.CheckError(err)
	if err != nil {
		return err
	}

	//12. 生成ansible
	g.Log().Info("[12]. 生成ansible...")
	err = CreateAnsibleInventory()
	utils.CheckError(err)
	if err != nil {
		return err
	}

	usedTime := time.Since(t0)
	g.Log().Infof("************************ 合服成功: 耗时:%s **********************", usedTime.String())
	return nil
}
func mergeDb(dbConfig *dbConfig) error {
	t0 := time.Now()
	fileName := g.Cfg().GetString("merge.tableConfigPath", "table_config.json")
	tableConfigFileData := gfile.GetBytes(fileName)
	//tableConfigFileData, err := ioutil.ReadFile("table_config.json")
	//utils.CheckError(err, "table_config.json失败")
	if tableConfigFileData == nil {
		g.Log().Errorf("读文件失败:%+v", fileName)
		return gerror.New("读文件失败")
	}
	tableConfig := &tableConfig{}
	err := json.Unmarshal(tableConfigFileData, tableConfig)
	utils.CheckError(err, "解析文件失败:"+fileName)
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbConfig.TargetDb.DbUser, dbConfig.TargetDb.DbPwd, dbConfig.TargetDb.DbHost, dbConfig.TargetDb.DbPort, dbConfig.TargetDb.DbName)
	fmt.Printf("目标数据库: %s:%s\n", dbConfig.TargetDb.DbHost, dbConfig.TargetDb.DbName)
	targetDb, err := xorm.NewEngine("mysql", dsn)
	utils.CheckError(err, "连接目标数据库失败:"+dsn)
	if err != nil {
		return err
	}
	defer targetDb.Close()
	_, err = targetDb.Exec("SET NAMES utf8;")
	utils.CheckError(err)
	if err != nil {
		return err
	}

	dbConfig.TargetDb.Db = targetDb
	if len(dbConfig.SourceDbList) == 0 {
		fmt.Print("[ERROR]:源数据库不能为空\n")
		return gerror.New("源数据库不能为空")
	}
	for i, e := range dbConfig.SourceDbList {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", e.DbUser, e.DbPwd, e.DbHost, e.DbPort, e.DbName)
		fmt.Printf("源数据库[%d]: %s:%s\n", i+1, e.DbHost, e.DbName)
		db, err := xorm.NewEngine("mysql", dsn)
		utils.CheckError(err, "连接源数据库失败:"+dsn)
		if err != nil {
			return err
		}
		defer db.Close()
		_, err = db.Exec("SET NAMES utf8;")
		utils.CheckError(err)
		if err != nil {
			return err
		}

		e.Db = db
	}

	g.Log().Infof("开始处理合服数据:%+v\n", gtime.Datetime())
	err = doMergeDb(dbConfig, tableConfig)
	utils.CheckError(err)
	if err != nil {
		return err
	}

	g.Log().Infof("处理公共数据...")
	sql := fmt.Sprintf("DELETE from `server_data` where id in (6,7);")
	_, err = dbConfig.TargetDb.Db.Exec(sql)
	if err != nil {
		return err
	}
	sql = fmt.Sprintf("INSERT INTO `server_data` VALUES (7,0,1,'',0),(6,0,%d,'',0);", utils.GetTimestamp())
	_, err = dbConfig.TargetDb.Db.Exec(sql)
	utils.CheckError(err, "设置目标服执行合服脚本失败:"+sql)

	if err != nil {
		return err
	}
	usedTime := time.Since(t0)
	fmt.Print("\n")
	fmt.Print("*****************************************************\n")
	fmt.Print("    数据库合并成功.\n")
	fmt.Printf("    耗时 %s. \n", usedTime.String())
	fmt.Print("*****************************************************\n\n")
	return nil
}

//func doCleanDb(dbConfig *db, cleanLevel int, cleanVipLevel int, cleanNoLoginDay int, dBMetas []*core.Table, tableConfig *tableConfig) error {
//
//}
func doCleanDb(db *db, dBMetas []*core.Table, dbConfig *dbConfig, tableConfig *tableConfig) error {
	//fmt.Printf("%s %s ", fmt.Sprintf("开始清理 %s", dbConfig.DbName), strings.Repeat(".", 50-len(fmt.Sprintf("开始清理 %s", dbConfig.DbName))))
	g.Log().Infof("开始清理数据库:%s......\n", db.DbName)

	now := utils.GetTimestamp()
	sql := fmt.Sprintf("delete player from player, player_data where player.`id` = player_data.`player_id` and player.`last_login_time` < %d and player_data.`level` <= %d and player_data.`vip_level` <= %d  ",
		now-86400*dbConfig.CleanNoLoginDay, dbConfig.CleanLevel, dbConfig.CleanVipLevel)

	if dbConfig.CleanLevel2 > 0 {
		sql = fmt.Sprintf("delete player from player, player_data where player.`id` = player_data.`player_id` and (player.`last_login_time` < %d and player_data.`level` <= %d and player_data.`vip_level` <= %d or player.`last_login_time` < %d and player_data.`level` <= %d and player_data.`vip_level` <= %d )",
			now-86400*dbConfig.CleanNoLoginDay, dbConfig.CleanLevel, dbConfig.CleanVipLevel, now-86400*dbConfig.CleanNoLoginDay2, dbConfig.CleanLevel2, dbConfig.CleanVipLevel2)
	}

	r, err := db.Db.Exec(sql)
	utils.CheckError(err, "清理玩家失败:")
	if err != nil {
		return err
	}
	cleanNum, err := r.RowsAffected()
	g.Log().Infof("清理%d个玩家.\n", cleanNum)
	utils.CheckError(err)
	if err != nil {
		return err
	}
	//删除默认外键关联数据
	g.Log().Info("开始清理默认关联表......\n")
	for _, dbMeta := range dBMetas {
		tableName := dbMeta.Name
		sql := fmt.Sprintf("desc `%s` `player_id`", tableName)
		rows, err := db.Db.QueryString(sql)
		utils.CheckError(err, "获取关联player_id 失败:"+sql)
		if err != nil {
			return err
		}
		//fmt.Printf("%s， %+v\n", tableName, rows)
		if len(rows) > 0 {
			sql := fmt.Sprintf("delete from `%s` where `player_id` NOT IN (SELECT `id` FROM `player`);", tableName)
			r, err = db.Db.Exec(sql)
			utils.CheckError(err, "清理关联表失败:"+sql)
			if err != nil {
				return err
			}
			_, err := r.RowsAffected()
			utils.CheckError(err)
			//if cleanNum > 0 {
			//	fmt.Printf("%s清理:%d\n", tableName, cleanNum)
			//}
			if err != nil {
				return err
			}
		}
	}
	g.Log().Info("清理默认关联表完毕.\n")
	//删除自定义外键关联数据
	g.Log().Info("开始清理自定义关联表......\n")
	for _, foreignKeyMap := range tableConfig.ForeignKeyMapList {
		sql := fmt.Sprintf("delete from `%s` where `%s` NOT IN (SELECT `%s` FROM `%s`);", foreignKeyMap.Table, foreignKeyMap.Filed, foreignKeyMap.ForeignKey, foreignKeyMap.ForeignTable)
		r, err = db.Db.Exec(sql)
		utils.CheckError(err, "清理关联表失败:"+sql)
		if err != nil {
			return err
		}
		cleanNum, err := r.RowsAffected()
		utils.CheckError(err)
		if cleanNum > 0 {
			fmt.Printf("%s清理:%d\n", foreignKeyMap.Table, cleanNum)
		}
	}
	fmt.Printf("清理自定义关联表完毕.\n")

	fmt.Printf("清理数据库%s完毕.\n\n\n", db.DbName)
	return nil
}
func doMergeDb(dbConfig *dbConfig, tableConfig *tableConfig) error {
	dBMetas, err := dbConfig.TargetDb.Db.DBMetas()
	utils.CheckError(err, "获取所有表失败:")
	if err != nil {
		return gerror.New(fmt.Sprintf("获取所有表失败err:%+v", err))
	}

	sqlOnceInsertMaxLimit := g.Cfg().GetInt("merge.sqlOnceInsertMaxLimit", 20000)
	mergeTablePageNum := g.Cfg().GetInt("merge.mergeTablePageNum", 3000000)
	if sqlOnceInsertMaxLimit < 1000 {
		sqlOnceInsertMaxLimit = 1000
	}

	g.Log().Info("开始清理数据库表...")
	//清理目标数据库
	err = doCleanDb(dbConfig.TargetDb, dBMetas, dbConfig, tableConfig)
	if err != nil {
		return err
	}
	//清理源数据库
	for _, s := range dbConfig.SourceDbList {
		//err = doCleanDb(s, dbConfig.CleanLevel, dbConfig.CleanVipLevel, dbConfig.CleanNoLoginDay, dBMetas, tableConfig)
		err = doCleanDb(s, dBMetas, dbConfig, tableConfig)
		if err != nil {
			return err
		}
	}

	g.Log().Info("开始合并数据库...")
	//开启事务
	session := dbConfig.TargetDb.Db.NewSession()
	defer session.Close()
	err = session.Begin()
	utils.CheckError(err, "开启事务失败:")
	if err != nil {
		return gerror.New(fmt.Sprintf("开启事务失败err:%+v", err))
	}

	//utils.GC()
	for _, dbMeta := range dBMetas {
		tableName := dbMeta.Name
		if utils.IsHaveArray(tableName, tableConfig.IgnoreList) {
			// 使用目标数据库的数据
			//fmt.Printf("%s %s [ignore]\n", tableName, strings.Repeat(".", 50-len(tableName)))
			continue
		}
		if utils.IsHaveArray(tableName, tableConfig.CleanList) {
			//清理目标数据库数据
			//fmt.Printf("%s %s ", tableName, strings.Repeat(".", 50-len(tableName)))
			sql := fmt.Sprintf("delete from %s;\n", tableName)
			//_, err := dbConfig.TargetDb.Db.Exec(sql)
			_, err := session.Exec(sql)
			utils.CheckError(err, "清空表数据失败:"+sql)
			if err != nil {
				session.Rollback()
				return gerror.New(fmt.Sprintf("清空表数据失败:%+v err:%+v", sql, err))
			}
			continue
			//fmt.Printf("[clean]\n")
		}
		// 合并各个源数据库数据到目标数据库
		//fmt.Printf("%s %s ", tableName, strings.Repeat(".", 50-len(tableName)))
		countSql := fmt.Sprintf("SELECT count(*) as totalCount FROM %s;", tableName)
		//sql := fmt.Sprintf("SELECT * FROM %s;", tableName)

		insertCount := 0 // 单次插入数量
		rowsLen := 0     // 单表数据总数

		for _, sourceDb := range dbConfig.SourceDbList {
			selectCountS, err := sourceDb.Db.QueryString(countSql)
			utils.CheckError(err, "读取源表总条数失败:"+countSql)
			if err != nil {
				return gerror.New(fmt.Sprintf("读取源表总条数失败:%+v err:%+v", countSql, err))
			}
			selectCount := gconv.Int(selectCountS[0]["totalCount"])
			g.Log().Infof("%s:读取源表总条数-----:%s\t%d", sourceDb.DbName, tableName, selectCount)
			if selectCount == 0 {
				continue
			}
			loopNum := 1
			if selectCount > mergeTablePageNum {
				loopNum = gconv.Int(math.Ceil(gconv.Float64(selectCount / mergeTablePageNum)))
			}
			for i := 0; i < loopNum; i++ {
				//utils.GC()
				g.Log().Infof("处理table:%s i:%d loopNum:%d\tnum:%d", tableName, i+1, loopNum, selectCount-i*mergeTablePageNum)
				pageSql := fmt.Sprintf("SELECT * FROM %s limit %d,%d;", tableName, i*mergeTablePageNum, mergeTablePageNum)
				rows, err := sourceDb.Db.QueryString(pageSql)
				utils.CheckError(err, "读取源表失败:"+pageSql)
				if err != nil {
					return gerror.New(fmt.Sprintf("读取源表失败:%+v err:%+v", pageSql, err))
				}
				rowsLen = len(rows)
				if rowsLen == 0 {
					continue
				}
				insertCols := make([]string, 0, sqlOnceInsertMaxLimit)
				insertCount = len(insertCols)
				//insertCols := make([]string, 0, rowsLen)
				for _, row := range rows {
					insertCol := make([]string, 0, len(dbMeta.Columns()))
					for _, c := range dbMeta.Columns() {
						col := c.Name
						value := row[col]
						if dbMeta.AutoIncrement != "" && dbMeta.AutoIncrement == col {
							insertCol = append(insertCol, "NULL")
						} else {
							insertCol = append(insertCol, "'"+Addslashes(value)+"'")
						}
					}
					insertCols = append(insertCols, "("+strings.Join(insertCol, ",")+")")
					insertCount++
					rowsLen--
					if insertCount < sqlOnceInsertMaxLimit { // 没有达到分组插入的数据条数
						continue
					}
					// 分段插入
					insertSql := fmt.Sprintf("INSERT INTO `%s` VALUES %s", tableName, strings.Join(insertCols, ","))
					//g.Log().Debug("sql:%s", insertSql)
					//_, err = dbConfig.TargetDb.Db.Exec(insertSql)
					_, err = session.Exec(insertSql)
					utils.CheckError(err, "插入数据失败:"+insertSql)
					if err != nil {
						session.Rollback()
						return gerror.New(fmt.Sprintf("插入数据失败:%+v err:%+v", insertSql, err))
					}
					if rowsLen > sqlOnceInsertMaxLimit {
						insertCols = make([]string, 0, sqlOnceInsertMaxLimit)
					} else {
						insertCols = make([]string, 0, rowsLen)
					}
					insertCount = len(insertCols)
				}
				if insertCount == 0 { // 没有可插入的数据
					continue
				}
				g.Log().Infof("meger INSERT table:%s rowLen:%d insertCount:%d insertColsLen:", tableName, rowsLen, insertCount, len(insertCols))
				insertSql := fmt.Sprintf("INSERT INTO `%s` VALUES %s", tableName, strings.Join(insertCols, ","))
				//g.Log().Debug("sql:%s", insertSql)
				//_, err = dbConfig.TargetDb.Db.Exec(insertSql)
				_, err = session.Exec(insertSql)
				utils.CheckError(err, "插入数据失败:"+insertSql)
				if err != nil {
					session.Rollback()
					return gerror.New(fmt.Sprintf("插入数据失败:%+v err:%+v", insertSql, err))
				}
			}
		}
	}
	g.Log().Info("提交事务:%s", dbConfig.TargetDb.DbName)
	err = session.Commit()
	utils.CheckError(err, "提交事务失败")
	if err != nil {
		return gerror.New(fmt.Sprintf("提交事务失败err:%+v", err))
	}
	g.Log().Info("合并数据库成功:%s", dbConfig.TargetDb.DbName)
	return nil
}

// addslashes() 函数返回在预定义字符之前添加反斜杠的字符串。
// 预定义字符是：
// 单引号（'）
// 双引号（"）
// 反斜杠（\）
func Addslashes(str string) string {
	tmpRune := []rune{}
	strRune := []rune(str)
	for _, ch := range strRune {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}
	return string(tmpRune)
}

// ================================================ 自动合服操作 ================================================
// 初始合服订时器
func InitMergeCron() {
	g.Log().Infof("初始合服订时器:%+v", gtime.Datetime())
	List := GetCanMergePlatform()
	for _, e := range List {
		if e.IsCanPlatformMergeState() == false {
			continue
		}
		g.Log().Infof("处理可合服的平台数据%s: %s %+v", e.PlatformId, utils.TimeIntFormDefault(e.MergeTime), e.MergeStr)
		err := StartCronMerge(e.PlatformId, e.MergeTime)
		if err != nil {
			g.Log().Errorf("初始合服订时器启动计时器失败:(%+v, %+v)err:%+v", e.PlatformId, e.MergeTime, err)
		}
	}
}

// 处理可合服的平台数据
func HandleCheckCanMergePlatform() {
	currTime := utils.GetTimestamp()
	List := GetCanTimeMergePlatform(currTime)
	for _, e := range List {
		g.Log().Infof("处理可合服的平台数据%s: %s %+v", e.PlatformId, utils.TimeIntFormDefault(e.MergeTime), e.MergeStr)
		go mergePlatformMergeList(e.PlatformId, e.MergeTime)
	}
}

// 启动合服订时器
func StartCronMerge(PlatformId string, MergeTime int) error {
	err, mergeData := GetMergeData(PlatformId, MergeTime)
	if err != nil {
		g.Log().Errorf("获取合服数据(%+v,%+v)失败：%+v", PlatformId, MergeTime, err)
		return err
	}
	if mergeData.IsCanPlatformMergeState() == false {
		g.Log().Warningf("启动合服订时器不是审核状态：%v：%+v", PlatformId, MergeTime)
		return gerror.New("启动合服订时器不是审核状态")
	}
	if mergeData.MergeTime <= utils.GetTimestamp() {
		go mergePlatformMergeList(PlatformId, MergeTime)
		return nil
	}
	g.Log().Infof("启动合服订时器:%+v  MergeTime:%+v", PlatformId, utils.TimeIntFormDefault(MergeTime))
	cronFun := func() {
		mergePlatformMergeList(PlatformId, MergeTime)
	}
	cronTimeStr := utils.TimestampToCronStr(gconv.Int64(MergeTime))
	g.Log().Infof("启动订时器CronTimeStr:%s  Datetime:%s", cronTimeStr, gtime.Datetime())
	_, err = gcron.AddOnce(cronTimeStr, cronFun, mergeData.getCronName())
	return err
}

// 获得平台合服key
func GetPlatformMergeKey(platformId string) string {
	return "PlatformMergeey_" + platformId
}

// 合服平台列表
func mergePlatformMergeList(PlatformId string, MergeTime int) error {
	currTime := utils.GetTimestamp()
	err, mergeData := GetMergeData(PlatformId, MergeTime)
	//utils.CheckError(err)
	if err != nil || mergeData.MergeStr == "" {
		ErrStr := fmt.Sprintf("未找到平台合服数据:%s", mergeData)
		g.Log().Error(ErrStr)
		return gerror.New(ErrStr)
	}
	if mergeData.IsCanPlatformMergeState() == false {
		ErrStr := fmt.Sprintf("不是审核通过状态合服:%s", mergeData)
		g.Log().Error(ErrStr)
		return gerror.New(ErrStr)
	}
	mergeList := StringToMergeList(mergeData.MergeStr)
	//err = json.Unmarshal([] byte(mergeData.MergeStr), &mergeList)
	//if err != nil {
	//	g.Log().Info("合服区服列表内容转换错误%s:%s", mergeData.PlatformId, mergeData.MergeStr)
	//	return err
	//}
	err = CheckPlatformMerge(PlatformId, mergeList)
	if err != nil {
		g.Log().Error(err)
		return err
	}
	Key := GetPlatformMergeKey(PlatformId)
	utils.SetCache(Key, gtime.Timestamp(), 6000)
	defer utils.DelCache(Key)

	MergeNodeList := getMergeListToNodeList(PlatformId, mergeList)
	mergeData.MergeState = MergeState_5
	err = Db.Save(mergeData).Error
	utils.CheckError(err)
	mergeLenStr := gconv.String(len(mergeList))
	SendBackgroundMsgTemplateHandle(MSG_TEMPLATE_MERGE_SERVER_CHANGE_STATE, PlatformId+mergeLenStr, PlatformId+"平台开始合服", "本次合:"+mergeLenStr+"个区服,详情："+mergeData.MergeStr)
	err = startMerge(PlatformId, MergeNodeList)
	if err == nil {
		SendBackgroundMsgTemplateHandle(MSG_TEMPLATE_MERGE_SERVER_CHANGE_STATE, PlatformId, PlatformId+"平台合服完成", "本次合服详情："+mergeData.MergeStr)
		mergeData.MergeState = MergeState_9
		mergeData.MergeUseTime = utils.GetTimestamp() - currTime
		err = Db.Save(mergeData).Error
		return err
	} else {
		SendBackgroundMsgTemplateHandle(MSG_TEMPLATE_MERGE_SERVER_CHANGE_STATE, PlatformId, PlatformId+"平台合服失败", "本次合服详情："+mergeData.MergeStr)
		mergeData.MergeState = MergeState_4
		mergeData.FailMsg = gconv.String(err)
		err = Db.Save(mergeData).Error
		return err
	}
}

// 开始合服
func startMerge(PlatformId string, MergeNodeList []mergeNode) error {
	for _, mergeNode := range MergeNodeList {
		g.Log().Infof("开始合服:%+v, %s", mergeNode.nodeList, mergeNode.zoneNode)
		err := startMergeHandle(mergeNode.nodeList[1:], mergeNode.nodeList[0], mergeNode.zoneNode)
		if err != nil {
			g.Log().Errorf("合服失败---:%+v err:%+v", mergeNode, err)
			return err
		}
	}
	return nil
}

// 获得当前区服的节点列表
func getMergeListToNodeList(PlatformId string, MergeList []MergeServerData) (MergeNodeList []mergeNode) {
	for _, e := range MergeList {
		nodeList, zoneNode, err := GetMergeInfo(PlatformId, e.S, e.E)
		utils.CheckError(err, "检测合服:获得节点数据失败")
		if len(nodeList) < 2 {
			g.Log().Info("同节点合服忽略:%s, %+v, %+v", PlatformId, e, nodeList)
			continue
		}
		MergeNodeList = append(MergeNodeList, mergeNode{
			nodeList,
			zoneNode,
			e.S,
			e.E,
		})
	}
	return MergeNodeList
}
