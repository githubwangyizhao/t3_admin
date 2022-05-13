package models

import (
	"gfWeb/library/utils"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

// 后台通知信息模板
type BackgroundMsgTemplate struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	State         int    `json:"state"`
	PhoneMsgCode  string `json:"phoneMsgCode"`
	MailUserList  []int  `json:"mailUserList" gorm:"-"`
	PhoneUserList []int  `json:"phoneUserList" gorm:"-"`
	ChangeTime    int64  `json:"changeTime"`
	UserId        int    `json:"userId"`
	UserName      string `json:"userName" gorm:"-"`
	IsAdd         int    `json:"isAdd" gorm:"-"`
}

const (
	MSG_TEMPLATE_CHECK_BACK_WEB            = 1  // 检测后台web活跃数据
	MSG_TEMPLATE_ROBOT_NOT_ENOUGH          = 2  // 机器节点不足-严重
	MSG_TEMPLATE_ROBOT_NOT_ENOUGH_WARNING  = 3  // 机器节点不足-警告
	MSG_TEMPLATE_MERGE_SERVER_CHANGE_STATE = 4  // 合服操作状态
	MSG_TEMPLATE_SERVER_CLOSE_LONG_TIME    = 5  // 平台停机维护时间过长(20分钟以上)
	MSG_TEMPLATE_SERVER_CLOSE_NUMBER       = 6  // 平台停服节点数
	MSG_TEMPLATE_PLATFORM_NOT_OPEN_ENTER   = 7  // 平台入口未开
	MSG_TEMPLATE_DEL_NOT_USE_DATA          = 8  // 平台删除没用的数据库
	MSG_TEMPLATE_AUTO_CREATE_SERVER_FAIL   = 9  // 定时自动开服失败
	MSG_TEMPLATE_MERGE_SERVER_REQUEST      = 10 // 合服请求
	MSG_TEMPLATE_MERGE_SERVER_AUDIT        = 11 // 合服审核结果
	MSG_TEMPLATE_MERGE_SERVER_REMIND       = 12 // 合服审核提醒
	MSG_TEMPLATE_UPDATE_SERVER_STATE       = 13 // 平台更新维护状态
	MSG_TEMPLATE_NOW_CREATE_SERVER_FAIL    = 14 // 立即开服失败
	MSG_TEMPLATE_H_CREATE_SERVER_FAIL      = 15 // 整点自动开服失败
	MSG_TEMPLATE_UPDATE_PLATFORM_VERSION   = 16 // 平台版本更新
)

// 获取单个后台信息数据模板
func GetBackgroundMsgTemplate(id int) (*BackgroundMsgTemplate, error) {
	data := &BackgroundMsgTemplate{
		Id: id,
	}
	err := Db.First(&data).Error
	data.MailUserList = GetUserRelDataToUserIdList(USER_REL_DATA_TYPE_MsgTemplate_MAIL, id)
	data.PhoneUserList = GetUserRelDataToUserIdList(USER_REL_DATA_TYPE_MsgTemplate_PHONE, id)
	return data, err
}

// 获取后台信息数据模板列表
func GetBackgroundMsgTemplateList(params *BaseQueryParam) ([]*BackgroundMsgTemplate, int64) {
	data := make([]*BackgroundMsgTemplate, 0)
	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "id")
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	err := Db.Model(&BackgroundMsgTemplate{}).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		//u, err := GetUserOne(e.UserId)
		//if err == nil {
		//	e.UserName = u.Name
		//}
		e.MailUserList = GetUserRelDataToUserIdList(USER_REL_DATA_TYPE_MsgTemplate_MAIL, e.Id)
		e.PhoneUserList = GetUserRelDataToUserIdList(USER_REL_DATA_TYPE_MsgTemplate_PHONE, e.Id)
	}
	return data, count
}

// 更新后台通知信息数据模板
func UpdateBackgroundMsgTemplate(data, OldBackgroundMsgTemplateDb *BackgroundMsgTemplate) error {
	data.ChangeTime = gtime.Timestamp()
	//delMailList := utils.ListMinus(OldBackgroundMsgTemplateDb.MailUserList, data.MailUserList)
	//addMailList := utils.ListMinus(data.MailUserList, OldBackgroundMsgTemplateDb.MailUserList)
	//delPhoneList := utils.ListMinus(OldBackgroundMsgTemplateDb.PhoneUserList, data.PhoneUserList)
	//addPhoneList := utils.ListMinus(data.PhoneUserList, OldBackgroundMsgTemplateDb.PhoneUserList)
	//g.Log().Debugf("delMailList:%+v addMailList:%+v delPhoneList:%+v addPhoneList:%+v ", delMailList, addMailList, delPhoneList, addPhoneList)
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	if err := UpdateUserRelDataByUserIdList(tx, USER_REL_DATA_TYPE_MsgTemplate_MAIL, data.Id, OldBackgroundMsgTemplateDb.MailUserList, data.MailUserList); err != nil {
		return err
	}
	if err := UpdateUserRelDataByUserIdList(tx, USER_REL_DATA_TYPE_MsgTemplate_PHONE, data.Id, OldBackgroundMsgTemplateDb.PhoneUserList, data.PhoneUserList); err != nil {
		return err
	}
	//if _, err := DeleteUserRelDataByTypeIdUserIdList(USER_REL_DATA_TYPE_MsgTemplate_MAIL, data.Id, delMailList); err != nil {
	//	tx.Rollback()
	//	return err
	//}
	//if _, err := DeleteUserRelDataByTypeIdUserIdList(USER_REL_DATA_TYPE_MsgTemplate_PHONE, data.Id, delPhoneList); err != nil {
	//	tx.Rollback()
	//	return err
	//}
	//if _, err := SaveUserRelDataByUserIdList(USER_REL_DATA_TYPE_MsgTemplate_MAIL, data.Id, addMailList); err != nil {
	//	tx.Rollback()
	//	return err
	//}
	//if _, err := SaveUserRelDataByUserIdList(USER_REL_DATA_TYPE_MsgTemplate_PHONE, data.Id, addPhoneList); err != nil {
	//	tx.Rollback()
	//	return err
	//}
	if err := Db.Save(&data).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// 删除后台通知信息数据模板
func DeleteBackgroundMsgTemplate(ids []int) error {
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := Db.Where(ids).Delete(&BackgroundMsgTemplate{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if _, err := DeleteUserRelDataByTypeIdList(USER_REL_DATA_TYPE_MsgTemplate_MAIL, ids); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := DeleteUserRelDataByTypeIdList(USER_REL_DATA_TYPE_MsgTemplate_PHONE, ids); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 发送给对应模板功能的成员
func SendBackgroundMsgTemplateMail(msgTemplateId int, mailTitleStr, mailBodyParam string) {
	SendBackgroundMsgTemplateHandle(msgTemplateId, "", mailTitleStr, mailBodyParam)
}
func SendBackgroundMsgTemplateHandle(msgTemplateId int, phoneMsgParam, mailTitleStr, mailBodyParam string) {
	SendBackgroundMsgTemplateHandleByUser(0, msgTemplateId, phoneMsgParam, mailTitleStr, mailBodyParam)
}
func SendBackgroundMsgTemplateHandleByUser(userId, msgTemplateId int, phoneMsgParam, mailTitleStr, mailBodyParam string) {
	msgTemplateData, err := GetBackgroundMsgTemplate(msgTemplateId)
	if err != nil {
		if userId > 0 && mailTitleStr != "" {
			sendMailUserList(mailTitleStr, mailBodyParam, utils.ArrayUnList(userId, []int{}))
		}
		g.Log().Errorf("通知信息数据模板不存:%+v", msgTemplateId)
		return
	}
	if msgTemplateData.State == 0 {
		g.Log().Warningf("通知信息数据模板关闭%+v:%+v", msgTemplateId, msgTemplateData.Name)
		return
	}
	mailUserList := gconv.Ints(msgTemplateData.MailUserList)
	if len(mailUserList) > 0 && mailTitleStr != "" {
		sendMailUserList(mailTitleStr, mailBodyParam, utils.ArrayUnList(userId, mailUserList))
	}
	if msgTemplateData.PhoneMsgCode != "" {
		phoneUserList := gconv.Ints(msgTemplateData.PhoneUserList)
		sendPhoneUserList(msgTemplateData.PhoneMsgCode, phoneMsgParam, utils.ArrayUnList(userId, phoneUserList))
	}
}

// 发送手机信息给玩家列表
func sendPhoneUserList(templateCode, templateParam string, userList []int) {
	m := gmap.New()
	for _, UserId := range userList {
		userData, err := GetUserOne(UserId)
		if err != nil || len(userData.Mobile) == 0 {
			continue
		}
		m.Set(userData.Mobile, userData.Mobile)
		//PhoneList = append(PhoneList, userData.Mobile)
	}
	if m.Size() == 0 {
		g.Log().Warningf("发送手机信息给玩家未找到玩家电话号码: %s  %s  %+v", templateCode, templateParam, userList)
		return
	}
	PhoneList := gconv.Strings(m.Keys())
	Key := "phone_" + templateCode + gconv.String(userList)
	OldTime := utils.GetCacheInt64(Key)
	var currTime = gtime.Timestamp()
	if OldTime+timeInterval > currTime {
		g.Log().Warningf("发送手机信息时间冷却时间内:%s", utils.TimeInt64FormDefault(OldTime))
		return
	}
	utils.SetCache(Key, currTime, gconv.Int(timeInterval))
	SendPhoneMsgByPhoneList(templateCode, templateParam, PhoneList)
}

// 发送邮件给玩家列表
func sendMailUserList(titleStr, bodyStr string, userList []int) error {
	m := gmap.New()
	for _, UserId := range userList {
		userData, err := GetUserOne(UserId)
		if err != nil || len(userData.MailStr) == 0 {
			continue
		}
		m.Set(userData.MailStr, userData.MailStr)
		//PhoneList = append(PhoneList, userData.Mobile)
	}
	if m.Size() == 0 {
		g.Log().Warningf("未找到发送邮件给玩家列表: %s  %+v", titleStr, userList)
		return gerror.New("未找到发送邮件给玩家列表")
	}
	MailList := gconv.Strings(m.Keys())
	Key := "mail_" + titleStr + gconv.String(userList)
	OldTime := utils.GetCacheInt64(Key)
	var currTime = gtime.Timestamp()
	if OldTime+timeInterval > currTime {
		g.Log().Warningf("邮件发送时间冷却时间内:%s", utils.TimeInt64FormDefault(OldTime))
		return gerror.New("邮件发送时间冷却时间内")
	}
	utils.SetCache(Key, currTime, gconv.Int(timeInterval))
	return SendMailHandle(titleStr, bodyStr, MailList)
}
