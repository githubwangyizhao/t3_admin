package models

import (
	"fmt"
	"gfWeb/library/utils"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcron"
	"github.com/jinzhu/gorm"
	"sort"
	"strconv"
	"strings"
)

// 合服数据
type MergeServerData struct {
	S int `json:"s,omitempty"`
	E int `json:"e,omitempty"`
}

type PlatformMergeServerData struct {
	PlatformId   string `json:"platformId" gorm:"PRIMARY_KEY"`
	MergeTime    int    `json:"mergeTime" gorm:"PRIMARY_KEY"`
	MergeState   int    `json:"mergeState"`
	MergeStr     string `json:"mergeStr"`
	RequestId    int    `json:"requestId"`
	RequestTime  int    `json:"requestTime"`
	AuditId      int    `json:"auditId"`
	AuditTime    int    `json:"auditTime"`
	MergeUseTime int    `json:"mergeUseTime"`
	FailMsg      string `json:"failMsg"`
}

type MergeServerParam struct {
	BaseQueryParam
	PlatformId string // 合服平台
	MergeState int    // 合服状态
}

// 1:请求合服,3:审核通过,5:合服中,9:合服完成
var MergeState_1 = 1 // 请求合服
var MergeState_3 = 3 // 审核通过
var MergeState_4 = 4 // 合服失败
var MergeState_5 = 5 // 合服中
var MergeState_9 = 9 // 合服完成

// 按照 MergeServerData.S 从小到大排序
type mergeServerSort []MergeServerData

func (m mergeServerSort) Len() int { // 重写 Len() 方法
	return len(m)
}
func (m mergeServerSort) Swap(i, j int) { // 重写 Swap() 方法
	m[i], m[j] = m[j], m[i]
}
func (m mergeServerSort) Less(i, j int) bool { // 重写 Less() 方法， 从小到大排序
	return m[i].S < m[j].S
}

// 获得广播名字
func (m PlatformMergeServerData) getCronName() string {
	return CRON_NAME_PLATFORM_MERGE + m.PlatformId + "_" + utils.TimeIntFormDefault(m.MergeTime)
}

//func getPlatformMergeName(Id int) string {
//	return CRON_NAME_PLATFORM_MERGE + gconv.String(Id)
//}
// 创建合服数据
func CreatePlatformMergeData(userId int, userName string, PlatformId string, MergeTime int, MergeList []MergeServerData) error {
	if len(MergeList) == 0 {
		return gerror.New("合服数据不能为空")
	}
	currTime := utils.GetTimestamp()
	if MergeTime == 0 {
		MergeTime = currTime
	}
	err, mergeData := GetMergeData(PlatformId, MergeTime)
	g.Log().Info("创建合服数据:%+v", mergeData)
	requestState := MergeState_1
	if mergeData.MergeState >= requestState {
		return gerror.New("已存当前合服数据")
	}
	sort.Sort(mergeServerSort(MergeList))
	err = CheckPlatformMerge(PlatformId, MergeList)
	if err != nil {
		return err
	}
	mergeStr := mergeListToString(MergeList)
	err = Db.Save(&PlatformMergeServerData{
		PlatformId:  PlatformId,
		MergeTime:   MergeTime,
		MergeState:  requestState,
		MergeStr:    mergeStr,
		RequestId:   userId,
		RequestTime: currTime,
	}).Error
	SendBackgroundMsgTemplateHandleByUser(userId, MSG_TEMPLATE_MERGE_SERVER_REQUEST, PlatformId, PlatformId+"合服请求", userName+"申请:"+utils.TimeIntFormDefault(MergeTime)+" 时间合服,详细区服内容："+mergeStr)
	//SendMailAdminUserLevel(PlatformId+"合服请求", userName+"申请:"+utils.TimeIntFormDefault(int64(MergeTime))+" 时间合服,详细区服内容："+mergeStr, 3)
	return err
}

// 审核通过
func AuditPlatformMergeData(UserId int, PlatformId string, MergeTime int) error {
	err, mergeData := GetMergeData(PlatformId, MergeTime)
	if mergeData.MergeState != MergeState_1 && mergeData.MergeState != MergeState_4 {
		g.Log().Warningf("审核通过:当前请求合服状态%+v", mergeData.MergeState)
		return gerror.New("当前不是请求合服状态")
	}
	currTime := utils.GetTimestamp()
	mergeData.MergeState = MergeState_3
	mergeData.AuditId = UserId
	mergeData.AuditTime = currTime
	err = Db.Save(mergeData).Error
	SendBackgroundMsgTemplateHandleByUser(mergeData.RequestId, MSG_TEMPLATE_MERGE_SERVER_AUDIT, PlatformId, PlatformId+"合服审核通过", "详细内容:"+utils.TimeIntFormDefault(MergeTime)+" 时间合服,详细区服内容："+mergeData.MergeStr)
	err = StartCronMerge(PlatformId, MergeTime)
	return err
}

// 删除平台合服
func DelPlatformMergeData(PlatformId string, MergeTime int) error {
	err, mergeData := GetMergeData(PlatformId, MergeTime)
	if err != nil || mergeData.MergeState == MergeState_5 {
		return gerror.New("合服中不可删除")
	}
	err = Db.Where(mergeData).Delete(&PlatformMergeServerData{}).Error
	if err == nil {
		gcron.Remove(mergeData.getCronName())
		SendBackgroundMsgTemplateHandleByUser(mergeData.RequestId, MSG_TEMPLATE_MERGE_SERVER_AUDIT, PlatformId, PlatformId+"合服请求退回", "详细内容:"+utils.TimeIntFormDefault(MergeTime)+" 时间合服,详细区服内容："+mergeData.MergeStr)
	}
	return err
}

// 检测合服列表问题
func CheckPlatformMerge(PlatformId string, MergeList []MergeServerData) error {
	checkList := make([]MergeServerData, 0)
	isHave := false
	haveStr := ""
	for _, mE := range MergeList {
		for _, cE := range checkList {
			if mE.S < cE.S && cE.S < mE.E || mE.S < cE.E && cE.E < mE.E || mE.S < cE.S && cE.E < mE.E {
				isHave = true
				if haveStr != "" {
					haveStr += ";"
				}
				haveStr += fmt.Sprintf("%d-%d包含%d-%d", mE.S, mE.E, cE.S, cE.E)
				continue
			}
		}
		checkList = append(checkList, mE)
	}
	if isHave == true {
		return gerror.New("区服有重叠：" + haveStr)
	}

	for _, e := range MergeList {
		g.Log().Infof("检测合服:%s, %+v", PlatformId, e)
		nodeList, _, err := GetMergeInfo(PlatformId, e.S, e.E)
		if err != nil {
			return gerror.New("检测合服:获得节点数据失败")
		}
		if len(nodeList) < 2 {
			g.Log().Info("检测合服列表问题同节点合服忽略:%s, %+v, %+v", PlatformId, e, nodeList)
			continue
		}
		for _, node := range nodeList {
			gameServerList := GetGameServerByNode(node.Node)
			for _, gameServer := range gameServerList {
				serverId, err := strconv.Atoi(SubString(gameServer.Sid, 1, len(gameServer.Sid)-1))
				utils.CheckError(err, "检测合服:区服不是数字")
				if serverId < e.S || serverId > e.E {
					ErrStr := fmt.Sprintf("合服区服配置错误,%d: [%d-%d] node:%s", serverId, e.S, e.E, node.Node)
					g.Log().Error(ErrStr)
					return gerror.New(ErrStr)
				}
			}
		}
	}
	return nil
}

func GetMergeInfo(platformId string, s int, e int) (nodeList []*ServerNode, zoneNode string, err error) {
	//zoneNode := ""
	nodes := make([]string, 0)
	for i := s; i <= e; i++ {
		gameServer, err := GetGameServerOne(platformId, fmt.Sprintf("s%d", i))
		if err != nil {
			return nodeList, zoneNode, err
		}
		if utils.IsHaveArray(gameServer.Node, nodes) {
			continue
		}
		nodes = append(nodes, gameServer.Node)
		serverNode, err := GetServerNode(gameServer.Node)
		if err != nil {
			return nodeList, zoneNode, err
		}
		if zoneNode == "" {
			zoneNode = serverNode.ZoneNode
		}
		nodeList = append(nodeList, serverNode)
	}
	return nodeList, zoneNode, err
}

// 检测还没有审核的合服请求
func CheckNotAuditPlatform() {
	str := ""
	for _, NotAuditPlatformData := range getNotAuditPlatformDataList() {
		str += fmt.Sprintf(" %s合服时间:%s;", NotAuditPlatformData.PlatformId, utils.TimeIntFormDefault(NotAuditPlatformData.MergeTime))
	}
	if len(str) > 0 {
		SendBackgroundMsgTemplateMail(MSG_TEMPLATE_MERGE_SERVER_REMIND, "合服审核提醒", "需要审核内容:"+str)
		//SendMailAdminUserLevel("合服审核提醒", "需要审核内容:"+str, 3)
	}
}

// 合服列表转成字符串
func mergeListToString(MergeList []MergeServerData) string {
	str := ""
	mergeLen := len(MergeList) - 1
	for index, mergeServerData := range MergeList {
		if index == mergeLen {
			str += fmt.Sprintf("%d-%d", mergeServerData.S, mergeServerData.E)
		} else {
			str += fmt.Sprintf("%d-%d;", mergeServerData.S, mergeServerData.E)
		}
	}
	return str
}

// 字符串转成合服结构
func StringToMergeList(MergeStr string) []MergeServerData {
	var mergeList []MergeServerData
	for _, serverGroupStr := range strings.Split(MergeStr, ";") {
		serverList := strings.Split(serverGroupStr, "-")
		S, err := strconv.Atoi(serverList[0])
		utils.CheckError(err)
		E, err := strconv.Atoi(serverList[1])
		utils.CheckError(err)
		mergeList = append(mergeList, MergeServerData{
			S, E,
		})
	}
	return mergeList
}

// 是否可合服状态
func (m PlatformMergeServerData) IsCanPlatformMergeState() bool {
	return m.MergeState == MergeState_3
}

// ===============================================获得数据=============================================================

// 获得当前平台MergeTime的合服数据
func GetMergeData(PlatformId string, MergeTime int) (error, *PlatformMergeServerData) {
	mergeData := &PlatformMergeServerData{
		PlatformId: PlatformId,
		MergeTime:  MergeTime,
	}
	err := Db.First(&mergeData).Error
	return err, mergeData
}

// 获得合服数据列表
func GetPlatformMergeServerList(params MergeServerParam) ([]*PlatformMergeServerData, int64) {
	var count int64
	data := make([]*PlatformMergeServerData, 0)
	f := func(db *gorm.DB) *gorm.DB {
		if params.PlatformId != "" || params.MergeState > 0 {
			SelectData := &PlatformMergeServerData{}
			if params.MergeState > 0 {
				SelectData.MergeState = params.MergeState
			}
			if params.PlatformId != "" {
				SelectData.PlatformId = params.PlatformId
			}
			return db.Where(SelectData)
		}
		return db
	}
	f(Db.Model(&PlatformMergeServerData{})).Order("merge_state").Order("merge_time desc").Count(&count).Offset(params.Offset).Limit(params.Limit).Find(&data)
	return data, count
}

// 获得可合服时间内的平台数据
func GetCanTimeMergePlatform(currTime int) []*PlatformMergeServerData {
	data := make([]*PlatformMergeServerData, 0)
	Db.Model(&PlatformMergeServerData{}).
		Where(&PlatformMergeServerData{
			MergeState: MergeState_3,
		}).Where("merge_time <= ?", currTime).Find(&data)
	return data
}

// 获得可合服的平台数据
func GetCanMergePlatform() []*PlatformMergeServerData {
	data := make([]*PlatformMergeServerData, 0)
	err := Db.Model(&PlatformMergeServerData{}).Where(
		&PlatformMergeServerData{
			MergeState: MergeState_3,
		}).Find(&data).Error
	utils.CheckError(err, "获得可合服的平台数据")
	return data
}

// 获得未审核数据
func getNotAuditPlatformDataList() []*PlatformMergeServerData {
	data := make([]*PlatformMergeServerData, 0)
	Db.Model(&PlatformMergeServerData{}).Where(&PlatformMergeServerData{MergeState: MergeState_1, AuditTime: 0}).Find(&data)
	return data
}
