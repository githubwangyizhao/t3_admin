package models

import (
	"fmt"
	"gfWeb/library/utils"
)

func (a *Resource) TableName() string {
	return ResourceTBName()
}

func ResourceTBName() string {
	return TableName("resource")
}

type ResourceQueryParam struct {
	BaseQueryParam
}

//资源
type Resource struct {
	Id       int         `json:"id"`
	Name     string      `json:"name"`
	Parent   *Resource   `json:"-"`
	ParentId int         `json:"parentId"`
	Children []*Resource `json:"children"`
	UrlFor   string      `json:"urlFor"`
	Url      string      `json:"url" gorm:"-"`
	//RoleResourceRel []*RoleResourceRel `json:"-"`
}

func relateResourceListParent(resourceList []*Resource) {
	for _, resource := range resourceList {
		relateResourceParent(resource)
	}
}

func relateResourceParent(resource *Resource) (*Resource, error) {
	if resource.ParentId > 0 {
		p := &Resource{}
		if err := Db.Model(&resource).Related(&p, "ParentId").Error; err != nil {
			return resource, err
		}
		resource.Parent = p
	}
	return resource, nil
}

//获取单个资源
func GetResourceOne(id int) (*Resource, error) {
	resource := &Resource{
		Id: id,
	}
	if err := Db.First(&resource).Error; err != nil {
		return nil, err
	}
	if _, err := relateResourceParent(resource); err != nil {
		return nil, err
	}
	return resource, nil
}

//获取资源列表
func GetResourceList() []*Resource {
	data := make([]*Resource, 0)
	err := Db.Model(&Resource{}).Find(&data).Error
	utils.CheckError(err)
	relateResourceListParent(data)
	return data
}

// 删除资源列表
func DeleteResources(ids []int) (int64, error) {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return 0, tx.Error
	}
	var count int64
	// 删除资源
	if err := Db.Where(ids).Delete(&Resource{}).Count(&count).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	//删除角色资源关系
	if _, err := DeleteRoleResourceRelByResourceIdList(ids); err != nil {
		tx.Rollback()
		return 0, err
	}
	return count, tx.Commit().Error
}

//获取可以成为某个节点父节点的列表
func ResourceTreeGrid4Parent(id int) []*Resource {
	list := GetResourceList()
	tmpList := make([]*Resource, 0)
	if id > 0 {
		for _, e := range list {
			if CanParent(id, e.Id) {
				tmpList = append(tmpList, e)
			}
		}
	} else {
		tmpList = list
	}
	return tmpList
}

func CanParent(resourceId int, parentResourceId int) bool {
	parentResource, _ := GetResourceOne(parentResourceId)
	if (parentResource.Parent != nil && parentResource.Parent.Id == resourceId) || resourceId == parentResourceId {
		//g.Log().Debug("CanParent:%+v %+v %+v", resourceId, thisResourceId, thisResource.Parent.Id)
		return false
	}
	if parentResource.Parent == nil || parentResource.Parent.Id == 0 {
		return true
	}
	return CanParent(resourceId, parentResource.Parent.Id)
}

//根据用户获取有权管理的资源列表
func GetResourceListByUserId(userId int) []*Resource {
	var list []*Resource
	user, err := GetUserOne(userId)
	utils.CheckError(err)
	if err != nil || user == nil {
		return list
	}

	var sql string
	if user.IsSuperUser() {
		list = GetResourceList()
	} else {
		sql = fmt.Sprintf(`SELECT DISTINCT T2.*
		FROM %s AS T0
		INNER JOIN %s AS T1 ON T0.role_id = T1.role_id
		INNER JOIN %s AS T2 ON T2.id = T0.resource_id
		WHERE T1.user_id = %d `, RoleResourceRelTBName(), RoleUserRelTBName(), ResourceTBName(), userId)
		rows, err := Db.Raw(sql).Rows()
		defer rows.Close()
		utils.CheckError(err)
		for rows.Next() {
			var resource Resource
			Db.ScanRows(rows, &resource)
			list = append(list, &resource)
		}
		relateResourceListParent(list)
	}
	return list
}

func TranResourceList2ResourceTree(resourceList []*Resource) []*Resource {
	resourceTree := make([]*Resource, 0)
	for _, item := range resourceList {
		if item.Parent == nil || item.Parent.Id == 0 {
			item = TranResourceList2ResourceTree_(item, resourceList)
			resourceTree = append(resourceTree, item)
		}
	}
	//g.Log().Debug("TranResourceList2ResourceTree:%+v", resourceTree)
	return resourceTree
}
func TranResourceList2ResourceTree_(cur *Resource, list []*Resource) *Resource {
	for _, item := range list {
		if item.Parent != nil && item.Parent.Id == cur.Id {
			item = TranResourceList2ResourceTree_(item, list)
			cur.Children = append(cur.Children, item)
		}
	}
	return cur
}
