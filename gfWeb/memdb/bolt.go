package memdb

import (
	"gfWeb/app/models"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gogf/gf/frame/g"
)

const (
	DATABASE_USER        = "user.db" // 邮箱配置库
	DATABASE_PUSH        = "push.db"
	DATABASE_GAME_USELOG = "uselog.db"

	TABLE_ACCOUNT = "account" //帐户表
	TABLE_LOG     = "log"
)

func init() {
	userDb := openDB(DATABASE_USER)
	pushDb := openDB(DATABASE_PUSH)
	uselogDb := openDB(DATABASE_GAME_USELOG)
	defer func() {
		pushDb.Close()
		userDb.Close()
		uselogDb.Close()
	}()
	initTable(userDb, TABLE_ACCOUNT)
	initTable(uselogDb, TABLE_LOG)

	// 查找所有的平台 ，即表名为 local , local2
	platforms := models.GetPlatformList()
	for _, p := range platforms {
		ScanDataInTable(pushDb, p.Id)
	}

}

func OpenDb(dbName string) *bolt.DB {
	return openDB(dbName)
}

func openDB(dbName string) *bolt.DB {
	db, err := bolt.Open(dbName, 0600, &bolt.Options{
		Timeout: time.Duration(10) * time.Second,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}

func initTable(db *bolt.DB, tableName string) {

	if db == nil {
		g.Log().Fatal("conf db is nil")
	}

	err := db.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucketIfNotExists([]byte(tableName))
		if err != nil {
			g.Log().Fatal(err.Error())
		}

		return nil
	})
	if err != nil {
		g.Log().Fatal(err.Error())
	}
}
