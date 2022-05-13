package utils

import (
	"bytes"
	"encoding/gob"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/util/gconv"
	"time"
)

//var cc gcache.Cache

func InitCache() {
	SetCacheCode("memory", `{"interval":60}`, 60)
}

// SetCache
func SetCache(key string, value interface{}, timeout int) {
	timeouts := time.Duration(timeout) * time.Second
	gcache.Set(key, value, timeouts)
}

func GetCache(key string) (interface{}, error) {
	data, err := gcache.Get(key)
	//if err == nil {
	//	return data, gerror.New("Cache不存在")
	//}
	return data, err
}

// SetCache
func SetCacheCode(key string, value interface{}, timeout int) error {
	data, err := Encode(value)
	if err != nil {
		return err
	}
	timeouts := time.Duration(timeout) * time.Second
	gcache.Set(key, data, timeouts)
	return nil
}
func GetCacheCode(key string, to interface{}) error {
	data ,err:= GetCache(key)
	if err != nil {
		return gerror.New("Cache不存在")
	}
	dataBytes := gconv.Bytes(data)
	if len(dataBytes) == 0 {
		to = ""
		return nil
	}
	err = Decode(dataBytes, to)
	if err != nil {
		g.Log().Errorf("GetCache失败key:%v  err:%v", key , err)
		return err
	}
	return nil
}

// 获得缓存 int64
func GetCacheInt64(key string) int64 {
	value, err := GetCache(key)
	if err == nil {
		return gconv.Int64(value)
	}
	return 0
}

// DelCache
func DelCache(key string) error {
	gcache.Remove(key)
	return nil
}

// Encode
// 用gob进行数据编码
//
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode
// 用gob进行数据解码
//
func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
