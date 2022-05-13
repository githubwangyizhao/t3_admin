package models

import (
	"gfWeb/library/utils"
)

//平台服务器关系表
type PlatformInventorySeverRel struct {
	Id                int
	PlatformId        string
	InventoryServerId int
	//Created           time.Time
}

func (a *PlatformInventorySeverRel) TableName() string {
	return PlatformInventorySeverRelTBName()
}

func PlatformInventorySeverRelTBName() string {
	return "platform_inventory_server_rel"
}

// 更新平台平台服务器关系
func UpdatePlatformInventorySeverRelByPlatformIdList(platformId string, inventorySeverIds []int) ([]*PlatformInventorySeverRel, error) {
	RelList := make([]*PlatformInventorySeverRel, 0)
	ResultRelList := make([]*PlatformInventorySeverRel, 0)
	dataRel := &PlatformInventorySeverRel{
		PlatformId: platformId,
	}
	err := Db.Model(dataRel).Where(dataRel).Find(&RelList).Error
	if err != nil {
		return RelList, err
	}
	for _, inventoryData := range RelList {
		isDel := true
		for index, id := range inventorySeverIds {
			if id == inventoryData.InventoryServerId {
				inventorySeverIds = append(inventorySeverIds[:index], inventorySeverIds[index+1:]...)
				isDel = false
				break
			}
		}
		if isDel {
			err = Db.Delete(inventoryData).Error
			utils.CheckError(err)
			continue
		}
		ResultRelList = append(ResultRelList, inventoryData)
	}
	for _, id := range inventorySeverIds {
		relation := &PlatformInventorySeverRel{PlatformId: platformId, InventoryServerId: id}
		ResultRelList = append(ResultRelList, relation)
	}
	return ResultRelList, err
}

// 删除平台服务器关系
func DeletePlatformInventorySeverRelByPlatformIdList(PlatformIdList []string) (int, error) {
	var count int
	err := Db.Where("platform_id in (?)", PlatformIdList).Delete(&PlatformInventorySeverRel{}).Count(&count).Error
	return count, err
}

// 删除平台服务器关系
func DeletePlatformInventorySeverRelByInventoryServerIdList(InventoryServerIdList []int) (int, error) {
	var count int
	err := Db.Where("inventory_server_id in (?)", InventoryServerIdList).Delete(&PlatformInventorySeverRel{}).Count(&count).Error
	return count, err
}
