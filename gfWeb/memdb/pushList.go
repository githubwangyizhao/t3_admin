package memdb

import (
	"encoding/json"
	"fmt"
	"gfWeb/app/models"
	"gfWeb/library/utils"
	"regexp"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

const (
	KeySidList  = "sid_list"
	KeyInitTime = "init_time"

	maxAreaSize = 99999 //每个区服区间最多只可以存10万

)

/*** 所有表名即平台名 */

// 获取所有的表
// func
// 初次扫描推送表
func ScanDataInTable(db *bolt.DB, tableName string) {
	//判断是初始化，内在为空则继续
	if db == nil {
		g.Log().Fatal("conf db is nil")
	}
	var readyInit = true
	err := db.Update(func(t *bolt.Tx) error {

		b := t.Bucket([]byte(tableName))
		if b != nil {
			time := b.Get([]byte(KeyInitTime))

			g.Log().Warning("local push list first init time : " + string(time))
			readyInit = false
			return nil
		} else {
			_, err := t.CreateBucketIfNotExists([]byte(tableName))
			if err != nil {
				g.Log().Fatal(err.Error())
			}
		}

		return nil
	})
	if err != nil {
		g.Log().Error(err.Error())
	}

	if !readyInit {
		return
	}

	g.Log().Warning("首次启动需要初始化推送列表，请耐心等待初始化完成....... ,也可能是在扫描新的表 : " + tableName)
	/**
	以区服为key值，一个区服所包含数据量暂定为100000（十万）,
	区服s153为例:
	s153_1, s153_2, s153_3..........
	当s153_1存满100000条数据后再存s153_2,依此类推
	*/

	var (
		allSid = make(map[string]string)

		// key1 区域 key2 帐号 key3 推送条件
		sidAreaData = make(map[string]map[string]map[string]string) //这里生成的数据最终生成对应的区服
	)
	/**
	分为三类表，
	1, 包含所有分区key的表， 表中数据为 "sids" : []string{"s1_1", "s1_2", "s153_1", "s160_1",....}
	2, 具体分区数据表,表中具体数据： map[playerId]cnt
	    "s1_1" : map[string]string{"13202": `{received:1}`, "13119": `{received:0}`,......}
	*/

	var accounts = make([]models.GlobalAccount, 0)
	// err := models.DbCenter.Raw("select * from global_account").Limit(10).Scan(&accounts).Error
	err = models.DbCenter.Raw("select * from global_account").Scan(&accounts).Error

	utils.CheckError(err)

	//将数据扫入变量
	for _, account := range accounts {
		s := account.RecentServerList
		s = s[1 : len(s)-1]
		listSids := regexp.MustCompile(`\[(.*?)\]`).FindAllString(s, -1)

		for _, sidItem := range listSids {
			sidItem := sidItem[1 : len(sidItem)-1]
			sids := strings.Split(strings.ReplaceAll(sidItem, " ", ""), ",")

			var sid = ""
			for _, i := range sids {
				iii, _ := strconv.Atoi(i)
				sid += string(rune(iii))
			}

			// 获取当前区间的sid号
			sidStr, maxSidNum := _fetchCurSid(sid, allSid)

			// 添加到全局列表中
			if allSid[sidStr] == "" {
				allSid[sidStr] = "ok"
			}

			// 查找数据库这个区服是否存储了10万数据，不足继续存入内在，否则开启新区服写入全局serverId列表中并存入新区间
			if len(sidAreaData[sidStr]) > maxAreaSize {
				// 新的区域
				sidStr = sid + "_" + gconv.String(maxSidNum+1)
				allSid[sidStr] = "ok" //写入全局
			}

			// 存入
			if sidAreaData[sidStr] == nil {
				sidAreaData[sidStr] = map[string]map[string]string{}
			}
			sidAreaData[sidStr][account.Account] = map[string]string{
				"open":  "1",                    // 1可以推送, todo后边需要加判断，先写入数据
				"regis": account.RegistrationId, // 推送的设备IDK号
			}

		}
	}

	timeStr := gtime.Datetime()
	// 将变量数据扫描存入内存
	err = db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(tableName))

		err := b.Put([]byte(KeyInitTime), []byte(timeStr))
		if err != nil {
			return err
		}

		//  插入 allSidList数据
		allSidList, _ := json.Marshal(allSid)
		err = b.Put([]byte(KeySidList), allSidList)
		if err != nil {
			return err
		}

		// 插入各区服区域数据
		for sidStr, _ := range allSid {

			areaData, _ := json.Marshal(sidAreaData[sidStr])
			err = b.Put([]byte(sidStr), areaData)
			if err != nil {
				return err
			}
		}

		/** 将不可推送数据扫入 ***/
		// 1, 分别获取 Function.json和NewsPush中数据，其中和function对比have_pf_list不为1时查找对应区服数据库的player_function表中的state ！= 1的用户
		// 2, select player_id from player_client_data where id = "k_jpush_status" and value = 1 查找到的用户也不推送

		return nil
	})

	if err != nil {
		g.Log().Error(err.Error())
		return
	}

	ScanNopushAccInTable(db, tableName)

}

// 初次招描不可推送帐户入内在表
func ScanNopushAccInTable(db *bolt.DB, tableName string) {

	url := g.Cfg().GetString("game.gameCenterHost") + "/static/json/Function.json"
	functionSliceMapStr := utils.GconvToSliceMapStr(utils.HttpGetJsonSliceMap(url))

	url2 := g.Cfg().GetString("game.gameCenterHost") + "/static/json/NewsPush.json"
	newsPushSliceMapStr := utils.GconvToSliceMapStr(utils.HttpGetJsonSliceMap(url2))

	var (
		functionMap = make(map[string]map[string]string)
		newsPushMap = make(map[string]map[string]string)
	)
	for _, item := range functionSliceMapStr {
		functionMap[item["id"]] = item
	}
	for _, item := range newsPushSliceMapStr {
		newsPushMap[item["id"]] = item
	}

	//  可能不发送的function名单
	var mayFunctionList []string
	for _, item := range newsPushMap {
		if functionMap[item["function"]] == nil {
			// fmt.Println("id is nil : " + item["function"])
			mayFunctionList = append(mayFunctionList, item["function"])
		} else if functionMap[item["function"]]["have_pf_list"] != "1" {
			mayFunctionList = append(mayFunctionList, item["function"])
			// fmt.Println(" id not nil and have pf list not is 1 : " + item["function"])
		}
	}

	sql := `select player_id from (select "111111" fid`
	for _, fidStr := range mayFunctionList {
		sql += ` union all select "` + fidStr + `" `
	}
	sql += `) as tmp, player_function t where tmp.fid=t.function_id and t.state != 1;` //
	//不能推送黑名单
	var playerIdBlackList = make(map[string]string)
	gameServerList, _ := models.GetAllGameServerDirty()

	type NeedRt struct {
		PlayerId string `json:"player_id"`
		Account  string `json:"account"`
	}
	for _, server := range gameServerList {
		gamedb, err := models.GetGameDbByNode(server.Node)
		utils.CheckError(err)

		if gamedb == nil {
			continue
		}

		var ret = make([]NeedRt, 0)
		err = gamedb.Raw(sql).Scan(&ret).Error
		utils.CheckError(err)

		for _, pf := range ret {
			if playerIdBlackList[gconv.String(pf.PlayerId)] == "" {
				playerIdBlackList[gconv.String(pf.PlayerId)] = "no"
			}
		}

		// PlayerClientData 表查找
		var ret2 = make([]NeedRt, 0)
		sql2 := `select player_id from player_client_data where id = "k_jpush_status" and value = 1`
		err = gamedb.Raw(sql2).Scan(&ret2).Error
		utils.CheckError(err)
		for _, pf2 := range ret2 {
			if playerIdBlackList[gconv.String(pf2.PlayerId)] == "" {
				playerIdBlackList[gconv.String(pf2.PlayerId)] = "no"
			}
		}
	}

	// 查到不推送帐户名单
	var (
		accountBlackList = make([]NeedRt, 0)
		// accountStr       = make([]string, 0)
	)
	if len(playerIdBlackList) > 0 {
		sql3 := `select account from (select "111111" pid`
		for pidStr, _ := range playerIdBlackList {
			sql3 += ` union all select "` + pidStr + `" `
		}
		sql3 += `) as tmp, global_player t where tmp.pid=t.id;` //
		err := models.DbCenter.Raw(sql3).Scan(&accountBlackList).Error
		utils.CheckError(err)
	}

	for _, rt := range accountBlackList {
		// fmt.Println("accountBlackList : ", rt.Account)
		errStr := SetNopushAccount(rt.Account, tableName, db)
		if errStr != "" {
			g.Log().Warning(errStr)
		}
	}

	g.Log().Warning("推送表初始化成功....")
}

func _fetchCurSid(sid string, allSid map[string]string) (curSidStr string, maxSidNum int) {
	var (
		inAllSid = false
	)
	for key, _ := range allSid {
		if strings.Contains(key, sid+"_") {
			inAllSid = true
			sidArr := strings.Split(key, "_")
			if maxSidNum < gconv.Int(sidArr[1]) {
				maxSidNum = gconv.Int(sidArr[1])
			}
		}
	}
	if !inAllSid {
		curSidStr = sid + "_1"
	} else {
		curSidStr = sid + "_" + gconv.String(maxSidNum)
	}

	return
}

// 设置不可推送帐户
func SetNopushAccount(accountId, tableName string, db *bolt.DB) string {

	if len(accountId) == 0 {
		return ""
	}

	if db == nil {
		db = openDB(DATABASE_PUSH)
		defer func() {
			db.Close()
		}()
	}

	isFind := false

	err := db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(tableName))
		if b == nil {
			g.Log().Error("push list table is nil ")
			return nil
		}

		var allSid = make(map[string]string) // 当前区服所有区域

		allSidBt := b.Get([]byte(KeySidList))
		err := json.Unmarshal(allSidBt, &allSid)
		if err != nil {
			return err
		}

		for areaSidStr, _ := range allSid {
			data := b.Get([]byte(areaSidStr))

			var curData = make(map[string]map[string]string)
			json.Unmarshal(data, &curData)

			if curData[accountId] != nil {
				//
				isFind = true
				curData[accountId] = map[string]string{
					"open":  "0",                         // 1可以推送, todo后边需要加判断，先写入数据
					"regis": curData[accountId]["regis"], // 推送的设备IDK号
				}
				//更新当前区服该帐号的信息

				curDataBt, _ := json.Marshal(curData)
				err := b.Put([]byte(areaSidStr), curDataBt)
				if err != nil {
					g.Log().Error(areaSidStr+" : update push list err  -  ", err.Error())
				}
			}
		}

		return nil
	})
	if err != nil {
		g.Log().Error("set not push list err : ", err.Error())
	}

	if !isFind {
		return accountId + "用户未在推送列表中..."
	} else {
		return ""
	}
}

/**
 * 更新或增加用户
 * 如果sidList 不为空时，新用户添加
 * 如果sidList为空是，更新帐号
 */
func InsertOrUpPushList(platform, accountId, registrationId string, sidList []string) (retStr string) {
	tableName := platform

	db := openDB(DATABASE_PUSH)
	defer db.Close()

	var (
		effectAreas []string
	)

	err := db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(tableName))
		if b == nil {
			g.Log().Error("push list table is nil ")
			return nil
		}

		var allSid = make(map[string]string) // 当前区服所有区域

		allSidBt := b.Get([]byte(KeySidList))
		err := json.Unmarshal(allSidBt, &allSid)
		if err != nil {
			return err
		}

		if len(sidList) == 0 {
			// update
			for areaSidStr, _ := range allSid {
				data := b.Get([]byte(areaSidStr))

				var curData = make(map[string]map[string]string)
				json.Unmarshal(data, &curData)
				if curData[accountId] != nil {
					//
					curData[accountId] = map[string]string{
						"open":  "1",            // 1可以推送
						"regis": registrationId, // 推送的设备IDK号
					}
					//更新当前区服该帐号的信息
					effectAreas = append(effectAreas, areaSidStr)

					curDataBt, _ := json.Marshal(curData)
					err := b.Put([]byte(areaSidStr), curDataBt)
					if err != nil {
						g.Log().Error(areaSidStr+" : update push list err  -  ", err.Error())
					}
				}
			}

		} else {
			// insert
			for _, sid := range sidList {

				var (
					maxAreaIndex int = 1
				)
				for serverIdStr, _ := range allSid {
					if strings.Contains(serverIdStr, sid+"_") {
						sidArr := strings.Split(serverIdStr, "_")
						if maxAreaIndex < gconv.Int(sidArr[1]) {
							maxAreaIndex = gconv.Int(sidArr[1])
						}
					}
				}

				// 插入
				inAreaSidStr := sid + "_" + gconv.String(maxAreaIndex)
				data := b.Get([]byte(inAreaSidStr))

				var curData = make(map[string]map[string]string)
				json.Unmarshal(data, &curData)

				if len(curData) > maxAreaSize || (maxAreaIndex == 1 && len(curData) == 0) {
					// 新区间第一个元素
					if maxAreaIndex == 1 {
						maxAreaIndex--
					}
					inAreaSidStr = sid + "_" + gconv.String(maxAreaIndex+1)
					curData = map[string]map[string]string{
						accountId: {
							"open":  "1",            // 1可以推送, todo后边需要加判断，先写入数据
							"regis": registrationId, // 推送的设备IDK号
						},
					}
					allSid[inAreaSidStr] = "ok"

					allSidList, _ := json.Marshal(allSid)
					err = b.Put([]byte(KeySidList), allSidList)
					if err != nil {
						return err
					}
				} else {
					curData[accountId] = map[string]string{
						"open":  "1",            // 1可以推送, todo后边需要加判断，先写入数据
						"regis": registrationId, // 推送的设备IDK号
					}
				}

				curDataLast, _ := json.Marshal(curData)
				err := b.Put([]byte(inAreaSidStr), curDataLast)
				if err != nil {
					g.Log().Error("insert push list err : ", err.Error())
				}
			}

		}

		retStr = "插入成功"
		if len(sidList) == 0 {
			retStr = "所在区服" + strings.Join(effectAreas, ",") + "下的帐号 " + accountId + "更新成功"
		}

		return nil
	})

	if err != nil {
		g.Log().Error("insert or update push list err : ", err.Error())
	}

	return
}

// 显示内存所有推送名单
func PushList(platform string) (ret []interface{}) {
	tableName := platform

	db := openDB(DATABASE_PUSH)
	defer db.Close()

	err := db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(tableName))
		if b == nil {
			g.Log().Error("push list table is nil ")
			return nil
		}

		var (
			allSid      = make(map[string]string)
			curAreaData = make(map[string]map[string]map[string]string)
		)
		allSidBt := b.Get([]byte(KeySidList))
		err := json.Unmarshal(allSidBt, &allSid)
		if err != nil {
			return err
		}

		ret = append(ret, allSid)

		for sid, _ := range allSid {
			data := b.Get([]byte(sid))

			var curData = make(map[string]map[string]string)
			json.Unmarshal(data, &curData)

			fmt.Println(sid, len(curData))
		}

		for sid, _ := range allSid {
			data := b.Get([]byte(sid))

			var curData = make(map[string]map[string]string)
			json.Unmarshal(data, &curData)

			curAreaData[sid] = curData
		}

		ret = append(ret, curAreaData)

		// if err != nil {
		// 	g.Log().Error("fetch push list fail", err.Error())
		// }

		return nil
	})

	if err != nil {
		g.Log().Error("PushList fetch list err : ", err.Error())
	}

	return
}

// 根据sids获了推送数据 1key areaSid  2key account
func FetchDataBySids(tableName string, sids []string) (ret map[string]map[string]map[string]string) {

	ret = make(map[string]map[string]map[string]string)

	db := openDB(DATABASE_PUSH)
	defer db.Close()

	err := db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(tableName))
		if b == nil {
			g.Log().Error("push list table is nil ")
			return nil
		}

		var (
			allSid = make(map[string]string)
		)
		allSidBt := b.Get([]byte(KeySidList))
		err := json.Unmarshal(allSidBt, &allSid)
		if err != nil {
			return err
		}

		for _, sid := range sids {
			sid = strings.ReplaceAll(sid, " ", "")
			for areaSid, _ := range allSid {
				if strings.Contains(areaSid, sid+"_") {
					data := b.Get([]byte(areaSid))

					var curData = make(map[string]map[string]string)
					json.Unmarshal(data, &curData)

					ret[areaSid] = curData
				}
			}
		}

		return nil
	})

	if err != nil {
		g.Log().Error("fetch data by serverIds err : ", err.Error())
	}
	return
}
