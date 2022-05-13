package models

import (
	"github.com/jinzhu/gorm"
)

type PlayerChargeRecord struct {
	OrderId         string `json:"orderId" gorm:"primary_key"`
	PlatformOrderId string `json:"platform_order_id"`
}

func GetItemByPlatformOrderId(db *gorm.DB, orderIds []string) (PlayerChargeRecords []*PlayerChargeRecord, err error) {
	err = db.Where("order_id in(?)", orderIds).Find(&PlayerChargeRecords).Error
	return
}

func FindOneByPlatformOrderId(db *gorm.DB, platformOrderId string) (*PlayerChargeRecord, bool) {
	var pcr = &PlayerChargeRecord{}
	isNotFound := db.Where(&PlayerChargeRecord{PlatformOrderId: platformOrderId}).First(&pcr).RecordNotFound()
	return pcr, isNotFound
}
