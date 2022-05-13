package memdb

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/gogf/gf/frame/g"
)

const (
	UserErrStr = "帐户已存在"
)

func Delete(name string) (errStr string) {
	db := openDB(DATABASE_USER)
	defer db.Close()

	err := db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(TABLE_ACCOUNT))
		if b == nil {
			g.Log().Error("account table is nil ")
			return nil
		}

		err := b.Delete([]byte(name))
		if err != nil {
			g.Log().Error(err.Error())
		}

		return nil
	})

	if err != nil {
		g.Log().Error(err.Error())
	}

	return
}

type TestAccountObj struct {
	Name      string `json:"name"`
	Privilege int    `json:"privilege"`
}

func InsertAccount(account TestAccountObj) (errStr string) {
	db := openDB(DATABASE_USER)
	defer db.Close()

	err := db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(TABLE_ACCOUNT))
		if b == nil {
			g.Log().Error("account table is nil ")
			return nil
		}

		data := b.Get([]byte(account.Name))
		if data != nil {
			errStr = UserErrStr
			return nil
		}

		data, err := json.Marshal(account)
		if err != nil {
			return err
		}
		b.Put([]byte(account.Name), []byte(data))
		return nil
	})

	if err != nil {
		g.Log().Error(err.Error())
	}

	return
}

func List() (ret []TestAccountObj) {
	db := openDB(DATABASE_USER)
	defer db.Close()

	err := db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(TABLE_ACCOUNT))
		if b == nil {
			g.Log().Error("account table is nil ")
			return nil
		}

		err := b.ForEach(func(k, v []byte) error {
			// fmt.Println("foreach ...." + string(k))
			var account TestAccountObj
			if err := json.Unmarshal(v, &account); err != nil {
				return err
			}

			if v != nil {
				ret = append(ret, account)
			}
			return nil
		})
		if err != nil {
			g.Log().Error("fetch tet accounts  fail")
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return
}
