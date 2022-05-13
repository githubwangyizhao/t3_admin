package uselog

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gogf/gf/frame/g"
)

var (
	DATABASE_GAME_USELOG = "uselog.db"
	TABLE_LOG            = "log"
)

func openDB(dbName string) *bolt.DB {
	db, err := bolt.Open(dbName, 0600, &bolt.Options{
		Timeout: time.Duration(10) * time.Second,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}

type MemItemEventLog struct {
	PlayerId   int     `json:"player_id"`
	PlatformId string  `json:"platform_id"`
	ServerId   string  `json:"server_id"`
	LogId      int     `json:"log_id"`
	Type       int     `json:"type"`
	Players    int     `json:"players"`
	Count      int     `json:"count"`
	Avg        float32 `json:"avg"`
	Time       int     `json:"time"`
	Number     int     `json:"number"`
	//
	MonsterId int `json:"monster_id"`
}

func FechTodayMonsterDataByPS(platform, serverId, monsterId string, datetime int) (ret map[string]map[string]*MemItemEventLog) {
	ret = make(map[string]map[string]*MemItemEventLog)
	t := time.Unix(int64(datetime), 0)
	timeKey := fmt.Sprintf("%d_%d_%d", t.Year(), t.Month(), t.Day())
	var key = platform + ":" + serverId + ":" + monsterId
	data := FetchALLDataByPS(key)
	ret[timeKey] = make(map[string]*MemItemEventLog)
	ret[timeKey] = data[timeKey]
	return
}

func FechTodayEventDataByPS(platform, serverId, logId, typeId string, datetime int) (ret map[string]map[string]*MemItemEventLog) {
	ret = make(map[string]map[string]*MemItemEventLog)
	t := time.Unix(int64(datetime), 0)
	timeKey := fmt.Sprintf("%d_%d_%d", t.Year(), t.Month(), t.Day())
	var key = platform + ":" + serverId + ":" + logId + ":" + typeId
	data := FetchALLDataByPS(key)
	ret[timeKey] = make(map[string]*MemItemEventLog)
	ret[timeKey] = data[timeKey]
	return
}

// 获取所有数据  key : map[data][hour]*MemItemEventLog
func FetchALLDataByPS(key string) (data map[string]map[string]*MemItemEventLog) {
	db := openDB(DATABASE_GAME_USELOG)
	defer db.Close()

	data = make(map[string]map[string]*MemItemEventLog)
	err := db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(TABLE_LOG))
		if b == nil {
			g.Log().Error(TABLE_LOG + " table is nil ")
			return nil
		}

		dataByte := b.Get([]byte(key))
		if dataByte == nil {
			g.Log().Warning(key + ", userlog data is empty")
			return nil
		}

		json.Unmarshal(dataByte, &data)
		return nil
	})

	if err != nil {
		g.Log().Error(err.Error())
	}

	return
}

func InsertMonsterLog(platform, serverId, monsterId string, data map[string]map[string]*MemItemEventLog) (errCode int, errStr string) {
	var key = platform + ":" + serverId + ":" + monsterId
	fmt.Println("monsterId key --------------- ", key)
	errCode, errStr = InsertLog(key, data)

	return
}

/**
platform:serverId 为key
data 1key 为日期， 2key为小时（1天24小时）
*/
func InsertEventLog(platform, serverId, logId, typeId string, data map[string]map[string]*MemItemEventLog) (errCode int, errStr string) {

	var key = platform + ":" + serverId + ":" + logId + ":" + typeId
	fmt.Println("key --------------- ", key)

	errCode, errStr = InsertLog(key, data)
	return
}

func InsertLog(key string, data map[string]map[string]*MemItemEventLog) (errCode int, errStr string) {
	db := openDB(DATABASE_GAME_USELOG)
	defer db.Close()

	err := db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(TABLE_LOG))
		if b == nil {
			g.Log().Error(TABLE_LOG + " table is nil ")
			return nil
		}

		dataByte := b.Get([]byte(key))

		dataJson, err := json.Marshal(data)
		if err != nil {
			return err
		}

		// 如果存在就覆盖更新并保留旧数据
		if dataByte != nil {
			var oldData = make(map[string]map[string]*MemItemEventLog)
			json.Unmarshal(dataByte, &oldData)

			for k, v := range data {
				if oldData == nil {
					oldData = make(map[string]map[string]*MemItemEventLog)
				}
				if oldData[k] == nil {
					oldData[k] = make(map[string]*MemItemEventLog)
				}
				for hour, item := range v {
					oldData[k][hour] = item
				}
			}

			dataJson, _ = json.Marshal(oldData)
		}

		b.Put([]byte(key), []byte(dataJson))
		return nil
	})

	if err != nil {
		g.Log().Error(err.Error())
	}

	return
}
