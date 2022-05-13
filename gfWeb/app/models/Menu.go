package models

import (
	"fmt"
	"gfWeb/library/utils"
	"sort"
)

func (a *Menu) TableName() string {
	return MenuTBName()
}

func MenuTBName() string {
	return TableName("menu")
}

type MenuQueryParam struct {
	BaseQueryParam
}

//Menu 权限控制表
type Menu struct {
	Id       int     `json:"id"`
	Title    string  `json:"title"`
	Name     string  `json:"name"`
	Parent   *Menu   `json:"-"`
	ParentId int     `json:"parentId"`
	Seq      int     `json:"seq"`
	IsShow   int     `json:"isShow"`
	Children []*Menu `json:"children"`
	Icon     string  `json:"icon"`
	//RoleMenuRel []*RoleMenuRel `json:"-"`
}

type menuSlice []*Menu

func (s menuSlice) Len() int      { return len(s) }
func (s menuSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s menuSlice) Less(i, j int) bool {
	if s[i].Seq == s[j].Seq {
		return s[i].Id < s[j].Id
	}
	return s[i].Seq < s[j].Seq
}

func sortMenuTree(list []*Menu) []*Menu {
	sort.Sort(menuSlice(list))
	for _, item := range list {
		if item.Children != nil {
			sortMenuTree(item.Children)
		}
	}
	return list
}

// 获取单个菜单
func GetMenuOne(id int) (*Menu, error) {
	menu := &Menu{
		Id: id,
	}
	if err := Db.First(&menu).Error; err != nil {
		return nil, err
	}
	if _, err := relateMenuParent(menu); err != nil {
		return nil, err
	}
	return menu, nil
}
func relateMenuListParent(menuList []*Menu) {
	for _, menu := range menuList {
		relateMenuParent(menu)
	}
}

func relateMenuParent(menu *Menu) (*Menu, error) {
	if menu.ParentId > 0 {
		p := &Menu{}
		if err := Db.Model(&menu).Related(&p, "ParentId").Error; err != nil {
			return menu, err
		}
		menu.Parent = p
	}
	return menu, nil
}

//获取菜单列表
func GetMenuList() []*Menu {
	data := make([]*Menu, 0)
	err := Db.Model(&Menu{}).Find(&data).Error
	utils.CheckError(err)
	relateMenuListParent(data)
	return data
}

// 删除菜单列表
func DeleteMenus(ids []int) (int64, error) {
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
	//删除菜单
	if err := Db.Where(ids).Delete(&Menu{}).Count(&count).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	// 删除角色菜单关系
	if _, err := DeleteRoleMenuRelByMenuIdList(ids); err != nil {
		tx.Rollback()
		return 0, err
	}
	return count, tx.Commit().Error
}

func MenuTreeGrid4Parent(id int) []*Menu {
	list := GetMenuList()
	tmpList := make([]*Menu, 0)
	if id > 0 {
		for _, e := range list {
			if CanParentMenu(id, e.Id) {
				tmpList = append(tmpList, e)
			}
		}
	} else {
		tmpList = list
	}
	return tmpList
}

func CanParentMenu(menuId int, parentMenuId int) bool {
	parentMenu, _ := GetMenuOne(parentMenuId)
	if (parentMenu.Parent != nil && parentMenu.Parent.Id == menuId) || menuId == parentMenuId {
		return false
	}
	if parentMenu.Parent == nil || parentMenu.Parent.Id == 0 {
		return true
	}
	return CanParentMenu(menuId, parentMenu.Parent.Id)
}

//根据用户获取有权管理的菜单列表
func GetMenuListByUserId(userId int) []*Menu {
	var list []*Menu
	user, err := GetUserOne(userId)
	utils.CheckError(err)
	if err != nil || user == nil {
		return list
	}

	var sql string
	if user.IsSuperUser() {
		list = GetMenuList()
	} else {
		sql = fmt.Sprintf(`SELECT DISTINCT T2.*
		FROM %s AS T0
		INNER JOIN %s AS T1 ON T0.role_id = T1.role_id
		INNER JOIN %s AS T2 ON T2.id = T0.menu_id
		WHERE T1.user_id = ?  Order By T2.seq asc,T2.id asc`, RoleMenuRelTBName(), RoleUserRelTBName(), MenuTBName())
		rows, err := Db.Raw(sql, userId).Rows()
		defer rows.Close()
		utils.CheckError(err)
		for rows.Next() {
			var menu Menu
			Db.ScanRows(rows, &menu)
			list = append(list, &menu)
		}
		relateMenuListParent(list)
	}
	return list
}

func TranMenuList2MenuTree(menuList []*Menu, isDealShow bool) []*Menu {
	menuTree := make([]*Menu, 0)
	for _, item := range menuList {
		if (item.Parent == nil || item.Parent.Id == 0) && (isDealShow == false || item.IsShow == 1) {
			item = TranMenuList2MenuTree_(item, menuList, isDealShow)
			menuTree = append(menuTree, item)
		}
	}
	sortMenuTree(menuTree)
	return menuTree
}

func TranMenuList2MenuTree_(cur *Menu, list []*Menu, isDealShow bool) *Menu {
	for _, item := range list {
		if (item.Parent != nil && item.Parent.Id == cur.Id) && (isDealShow == false || item.IsShow == 1) {
			item = TranMenuList2MenuTree_(item, list, isDealShow)
			cur.Children = append(cur.Children, item)
		}
	}
	return cur
}
