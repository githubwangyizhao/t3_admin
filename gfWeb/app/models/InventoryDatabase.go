package models

import (
	"gfWeb/library/utils"
)

func (a *InventoryDatabase) TableName() string {
	return InventoryDatabaseTBName()
}

func InventoryDatabaseTBName() string {
	return TableName("inventory_database")
}

type InventoryDatabaseParam struct {
	BaseQueryParam
}
type InventoryDatabase struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	User       string `json:"user"`
	Port       int    `json:"port"`
	Host       string `json:"host"`
	AddTime    int    `json:"addTime"`
	UpdateTime int    `json:"updateTime"`
}

//获取用户列表
func GetInventoryDatabaseList(params *InventoryDatabaseParam) ([]*InventoryDatabase, int64) {
	data := make([]*InventoryDatabase, 0)
	sortOrder := "id"
	switch params.Sort {
	case "id":
		sortOrder = "id"
	}
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	err := Db.Model(&InventoryDatabase{}).Where(&InventoryDatabase{}).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	return data, count
}

//获取所有数据库
func GetAllInventoryDatabaseList() []*InventoryDatabase {
	data := make([]*InventoryDatabase, 0)
	err := Db.Model(&InventoryDatabase{}).Find(&data).Error
	utils.CheckError(err)
	return data
}

// 获取单个用户
func GetInventoryDatabaseOne(id int) (*InventoryDatabase, error) {
	r := &InventoryDatabase{
		Id: id,
	}
	err := Db.First(&r).Error
	return r, err
}

// 删除用户列表
func DeleteInventoryDatabases(ids []int) error {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := Db.Where(ids).Delete(&InventoryDatabase{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
