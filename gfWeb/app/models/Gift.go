package models

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
)

type Item struct {
	ItemId  int `json:"itemId"`
	ItemNum int ` json:"itemNum"`
}

type GiftCode struct {
	GiftCode     string `json:"gift_code" gorm:"gift_code"`
	GiftCodeType int    `json:"gift_code_type" gorm:"gift_code_type"`
}

type GiftCodeType struct {
	Type                 int    `json:"id" gorm:"type"`
	Name                 string `json:"name" gorm:"name"`
	PlatformId           string `json:"platformId" gorm:"platform_id"`
	ChannelList          string `json:"channelList" gorm:"channel_list"`
	AwardList            string `json:"awardList" gorm:"award_list"`
	UserId               int    `json:"userId" gorm:"user_id"`
	Kind                 int    `json:"kind" gorm:"kind"`
	Num                  int    `json:"num" gorm:"num"`
	AllowRoleRepeatedGet int    `json:"allowRoleRepeatedGet" gorm:"allow_role_repeated_get"`
	VipLimit             int    `json:"vipLimit" gorm:"vip_limit"`
	LevelLimit           int    `json:"levelLimit" gorm:"level_limit"`
	ExpireTime           int    `json:"expireTime" gorm:"expire_time"`
	UpdateTime           int    `json:"updateTime" gorm:"update_time"`
}

type Gift struct {
	Name                 string   `json:"name"`
	GiftCode             string   `json:"giftCode"`
	GiftCodeType         int      `json:"giftCodeType"`
	PlatformId           string   `json:"platformId"`
	ChannelList          []string `json:"channelList"`
	AwardList            []Item   `json:"awardList"`
	Kind                 int      `json:"kind"`
	AllowRoleRepeatedGet int      `json:"allowRoleRepeatedGet"`
	Num                  int      `json:"num"`
	ExpireTime           int      `json:"expireTime"`
	VipLimit             int      `json:"vipLimit"`
	LevelLimit           int      `json:"levelLimit"`
	UserId               int      `json:"userId"`
	UpdateTime           int      `json:"update_time"`
	UserName             string   `json:"username"`
}

type GiftRequest struct {
	Gift
	BaseQueryParam
}

type Req2Erlang4Gift struct {
	ConfigType int64 `json:"config_type"`
	ConfigId   int64 `json:"config_id"`
	Value      int64 `json:"value"`
}

var baseStr = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijklmnpqrstuvwxy"
var base = []byte(baseStr)
var BaseMap map[byte]int

func DownloadGiftCode(params *GiftRequest) ([]*GiftCode, error) {
	var (
		data = make([]*GiftCode, 0)
	)
	g.Log().Infof("下载礼包码excel:%+v", params)

	query := DbCenter
	if params.GiftCodeType != 0 {
		query = query.Where("gift_code_type = ? ", params.GiftCodeType)
	}
	err := query.Debug().Find(&data).Offset(0).Error
	if err != nil {
		return data, err
	}

	return data, err
}

func GetGiftList(params *GiftRequest) ([]*Gift, int, error) {
	var (
		count    = 0
		data     = make([]*GiftCodeType, 0)
		giftData = make([]*Gift, 0)
	)
	g.Log().Infof("获取指定礼包码的数据:%+v", params)

	query := DbCenter
	if params.PlatformId != "" {
		query = query.Where("platform_id = ? ", params.PlatformId)
	}
	if params.VipLimit != 0 {
		query = query.Where("vip_limit = ? ", params.VipLimit)
	}
	if params.LevelLimit != 0 {
		query = query.Where("level_limit = ? ", params.LevelLimit)
	}
	if params.ExpireTime > 0 {
		query = query.Where("expire_time = ? ", params.ExpireTime)
	}
	g.Log().Infof("params.ChannelList: %+v %d", params.ChannelList, len(params.ChannelList))
	//if len(params.ChannelList) > 0 {
	//query = query.Where("channel_list in (?)", GetSQLWhereParamWithoutQuotation(params.ChannelList))
	//query = query.Where("channel_list = ?", fmt.Sprintf(`[\"%s\"]`, strings.Join(params.ChannelList, "\",\"")))
	//}

	err := query.Debug().Offset(params.Offset).Limit(params.Limit).Find(&data).Offset(0).Count(&count).Error
	if err != nil {
		return giftData, count, err
	}

	for _, i := range data {
		Items := make([]Item, 0)
		err = json.Unmarshal([]byte(i.AwardList), &Items)
		channel := make([]string, 0)
		err = json.Unmarshal([]byte(i.ChannelList), &channel)
		singleGiftData := &Gift{
			GiftCodeType: i.Type,
			UserId:       i.UserId,
			PlatformId:   i.PlatformId,
			VipLimit:     i.VipLimit,
			LevelLimit:   i.LevelLimit,
			ExpireTime:   i.ExpireTime,
			Name:         i.Name,
			AwardList:    Items,
			ChannelList:  channel,
			Num:          i.Num,
			UpdateTime:   i.UpdateTime,
		}
		u, err := GetUserOne(i.UserId)
		if err == nil {
			singleGiftData.UserName = u.Name
		}
		giftData = append(giftData, singleGiftData)
	}

	return giftData, count, err
}

func InitBaseMap() {
	BaseMap = make(map[byte]int)
	for i, v := range base {
		BaseMap[v] = i
	}
}

func Base54(n uint64) []byte {
	quotient := n
	mod := uint64(0)
	l := list.New()
	for quotient != 0 {
		mod = quotient % 34
		quotient = quotient / 34
		l.PushFront(base[int(mod)])
	}
	listLen := l.Len()

	if listLen >= 6 {
		res := make([]byte, 0, listLen)
		for i := l.Front(); i != nil; i = i.Next() {
			res = append(res, i.Value.(byte))
		}
		return res
	} else {
		res := make([]byte, 0, 6)
		for i := 0; i < 6; i++ {
			if i < 6-listLen {
				res = append(res, base[0])
			} else {
				res = append(res, l.Front().Value.(byte))
				l.Remove(l.Front())
			}

		}
		return res
	}

}

func Base54ToNum(str []byte) (uint64, error) {
	if BaseMap == nil {
		return 0, errors.New("no init base map")
	}
	if str == nil || len(str) == 0 {
		return 0, errors.New("parameter is nil or empty")
	}
	var (
		res uint64 = 0
		r   uint64 = 0
	)
	for i := len(str) - 1; i >= 0; i-- {
		v, ok := BaseMap[str[i]]
		if !ok {
			fmt.Printf("")
			return 0, errors.New("character is not base")
		}
		var b uint64 = 1
		for j := uint64(0); j < r; j++ {
			b *= 34
		}
		res += b * uint64(v)
		r++
	}
	return res, nil
}
