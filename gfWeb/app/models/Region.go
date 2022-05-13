package models

import (
	"gfWeb/library/utils"
	"github.com/gogf/gf/frame/g"
)

type Region struct {
	Id          int    `json:"id" gorm:"id"`
	Region      string `json:"region" gorm:"region"`
	Currency    string `json:"currency" gorm:"currency"`
	AreaCode    string `json:"area_code" gorm:"area_code"`
	CreatedBy   int    `json:"created_by" gorm:"created_by"`
	CreatedAt   int    `json:"created_at" gorm:"created_at"`
	UpdatedBy   int    `json:"updated_by" gorm:"updated_by"`
	UpdatedAt   int    `json:"updated_at" gorm:"updated_at"`
	CreatedName string `json:"created_name" gorm:"-"`
	UpdatedName string `json:"updated_name" gorm:"-"`
}

type RegionQueryParam struct {
	BaseQueryParam
	AppId   int    `json:"app_id"`
	Type    int    `json:"type"`
	Version string `json:"version"`
}

func GetRegions(params *RegionQueryParam) ([]*Region, int) {
	g.Log().Info("getList params: %+v", params)
	data := make([]*Region, 0)
	count := 0
	sortOrder := "id"
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	err := Db.Debug().Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		e.CreatedName = GetUserName(e.CreatedBy)
		e.UpdatedName = GetUserName(e.UpdatedBy)
	}
	return data, count
}

func AddRegion(params *Region) error {
	return Db.Debug().Save(params).Error
}

func EditRegion(params *Region) error {
	return Db.Debug().Save(params).Error
}

func DeleteRegion(params *Region) error {
	return Db.Debug().Delete(&params).Error
}

func GetRegionById(id int) *Region {
	data := &Region{
		Id: id,
	}
	err := Db.Debug().First(&data).Error
	if err != nil {
		return nil
	}
	return data
}
