package models

import (
	"gfWeb/library/utils"
	"time"
)

type PromoteData struct {
	Id          int       `json:"id"`
	Promote     string    `json:"promote"`
	Name        string    `json:"name"`
	State       int       `json:"state"`
	CreatedAt   time.Time `json:"createdAt"`
	CreatedBy   int       `json:"createdBy"`
	CreatedName string    `json:"createdName" gorm:"-"`
	UpdatedAt   time.Time `json:"updatedAt"`
	UpdatedBy   int       `json:"updatedBy"`
	UpdatedName string    `json:"updatedName" gorm:"-"`
}

type PromoteRequest struct {
	PromoteData
	BaseQueryParam
}

// 获取单个推广员数据
func GetPromoteDataOne(id int) (*PromoteData, error) {
	promoteData := &PromoteData{
		Id: id,
	}
	err := Db.Where(&promoteData).First(&promoteData).Error
	return promoteData, err
}

// 获取单个推广员数据
func GetPromoteDataOneByPromote(promote string) (*PromoteData, error) {
	promoteData := &PromoteData{
		Promote: promote,
	}
	err := Db.Where(&promoteData).First(&promoteData).Error
	return promoteData, err
}

// 获取推广员数据列表
func GetPromoteDataList(params *PromoteRequest) ([]*PromoteData, int64) {
	data := make([]*PromoteData, 0)
	var count int64
	if params.Id > 0 {
		promoteData, err := GetPromoteDataOne(params.Id)
		utils.CheckError(err)
		if err != nil {
			return data, count
		}
		count = 1
		data = append(data, promoteData)
	} else {
		//data1 := make([]*PromoteData, 0)
		//err := Db.Find(&data1).Count(&count).Error
		//utils.CheckError(err)
		params.Promote = "%" + params.Promote + "%"
		params.Name = "%" + params.Name + "%"
		sortOrder := "id"
		if params.Order == "descending" {
			sortOrder = sortOrder + " desc"
		}
		err := Db.Where("promote LIKE ?", params.Promote).Where("name LIKE ?", params.Name).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
		utils.CheckError(err)
		for _, e := range data {
			e.CreatedName = GetUserName(e.CreatedBy)
			e.UpdatedName = GetUserName(e.UpdatedBy)
		}
	}
	return data, count
}

// 新增推广员数据
func AddPromoteDate(promoteData *PromoteData) error {
	err := Db.Save(promoteData).Error
	return err
}

// 更新推广员数据
func UpdatePromoteDate(promoteData *PromoteData) error {
	err := Db.Save(promoteData).Error
	return err
}

// 删除用户列表
func DeletePromoteDataList(ids []int) error {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := Db.Where(ids).Delete(&PromoteData{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
