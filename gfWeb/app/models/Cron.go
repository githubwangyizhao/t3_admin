package models

import (
	"fmt"

	"gfWeb/library/utils"
	"math"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/util/gconv"
)

const (
	CRON_NAME_SYS_CRON                = "sys_cron_"                // 系统定时
	CRON_NAME_NOTICE_NOTICE           = "notice_notice_"           // 广播通知
	CRON_NAME_PLATFORM_MERGE          = "platform_merge_"          // 平台合服
	CRON_NAME_PLATFORM_UPDATE_VERSION = "platform_update_version_" // 平台版本更新
	CRON_NAME_PLATFORM_OPEN_SERVER    = "platform_open_server_"    // 平台定时开服
	CRON_NAME_VERSION_TOOL_CHANGE     = "version_tool_change_"     // 订时版本操作工具
)

type CronData struct {
	Name         string `json:"name"`
	RegisterTime int64  `json:"registerTime"`
	State        int    `json:"state"`
}

//初始化订时器
func InitCron() {
	fmt.Println("start cron..")
	gcron.SetLogLevel(glog.LEVEL_INFO)
	//每一个域都使用数字，但还可以出现如下特殊字符，它们的含义是：
	//(1)*：表示匹配该域的任意值，假如在Minutes域使用*, 即表示每分钟都会触发事件。
	//(2)?:只能用在DayofMonth和DayofWeek两个域。它也匹配域的任意值，但实际不会。因为DayofMonth和 DayofWeek会相互影响。例如想在每月的20日触发调度，不管20日到底是星期几，则只能使用如下写法： 13 13 15 20 * ?, 其中最后一位只能用？，而不能使用*，如果使用*表示不管星期几都会触发，实际上并不是这样。
	//(3)-:表示范围，例如在Minutes域使用5-20，表示从5分到20分钟每分钟触发一次
	//(4)/：表示起始时间开始触发，然后每隔固定时间触发一次，例如在Minutes域使用5/20,则意味着5分钟触发一次，而25，45等分别触发一次.
	//(5),:表示列出枚举值值。例如：在Minutes域使用5,20，则意味着在5和20分每分钟触发一次。
	//(6)L:表示最后，只能出现在DayofWeek和DayofMonth域，如果在DayofWeek域使用5L,意味着在最后的一个星期四触发。
	//(7)W: 表示有效工作日(周一到周五),只能出现在DayofMonth域，系统将在离指定日期的最近的有效工作日触发事件。例如：在 DayofMonth使用5W，如果5日是星期六，则将在最近的工作日：星期五，即4日触发。如果5日是星期天，则在6日(周一)触发；如果5日在星期一 到星期五中的一天，则就在5日触发。另外一点，W的最近寻找不会跨过月份
	//(8)LW:这两个字符可以连用，表示在某个月最后一个工作日，即最后一个星期五。
	//(9)#:用于确定每个月第几个星期几，只能出现在DayofMonth域。例如在4#2，表示某月的第二个星期三。

	// 0 30 * * * *
	// 秒（0~59）		Seconds:出现", - * /"四个字符，有效范围为0-59的整数
	// 分钟（0~59）		Minutes:可出现", - * /"四个字符，有效范围为0-59的整数
	// 小时（0~23）		Hours:可出现", - * /"四个字符，有效范围为0-23的整数
	// 天（月）（0~31，但是你需要考虑你月的天数）DayofMonth:可出现", - * / ? L W C"八个字符，有效范围为0-31的整数
	// 月（0~11） 		Month:可出现", - * /"四个字符，有效范围为1-12的整数或JAN-DEc
	// 天（星期）（1~7 1=SUN 或 SUN，MON，TUE，WED，THU，FRI，SAT） DayofWeek:可出现", - * / ? L C #"四个字符，有效范围为1-7的整数或SUN-SAT两个范围。1表示星期天，2表示星期一， 依次类推
	// 年份（1970－2099）Year:可出现", - * /"四个字符，有效范围为1970-2099年
	// 长度为6/7 :6位时年参数忽略

	g.Log().Info("初始化定时器")
	// 秒 分 时 天 月 周
	_, err := gcron.Add("*/10 * * * * *", tenSecondCron, getSysCronName("every10S")) // 每10秒定时器
	utils.CheckError(err)
	_, err = gcron.Add("0 * * * * *", minuteCron, getSysCronName("every1M")) // 每分钟定时器
	utils.CheckError(err)
	_, err = gcron.Add("0 */10 * * * *", tenMinuteClockCron, getSysCronName("every10M")) // 每10分钟整
	utils.CheckError(err)
	_, err = gcron.Add("0 0 * * * *", everyHourClockCron, getSysCronName("every1H")) // 每小时定时器
	utils.CheckError(err)
	_, err = gcron.Add("0 0 0 * * *", dailyZeroClockCron, getSysCronName("0H")) // 每日0点定时器
	utils.CheckError(err)
	_, err = gcron.Add("0 5 0 * * *", dailyZeroClock5MinuteCron, getSysCronName("0H5M")) // 每日0点5分定时器
	utils.CheckError(err)

	//// 10秒定时器
	//go tenSecondCron()
	//// 每分钟定时器
	//go minuteCron()
	//// 每小时定时器
	//go everyHourClockCron()
	//// 每日0点定时器
	//go dailyZeroClockCron()
	//// 每日0点5分定时器
	//go dailyZeroClock5MinuteCron()
	////整10分钟
	//go tenMinuteClockCron()

	// 定时检测开服
	go cronAutoCreateServer()
	go InitNoticeLog()           // 初始广播内容订时器
	go InitMergeCron()           // 初始合服订时器
	go InitPlatformVersionCron() // 初始平台更新订时器
	// 修复全部节点活跃留存数据
	go RepairAllGameNodeRemainActive()
	// 修复总体留存数据
	go RepireRemainTotal()

	// 初始广播内容时
	//go service.InitNoticeLog()

	// 修复ltv 留存数据 每日统计
	// go repairData()
}

func getSysCronName(name string) string {
	return CRON_NAME_SYS_CRON + name
}

// 获得订时器列表
func GetCronList(cronName string) ([]*CronData, int) {
	data := make([]*CronData, 0)
	for _, e := range gcron.Entries() {
		//g.Log().Debugf("获得订时器列表:%+v", gconv.String(runtime.FuncForPC(reflect.ValueOf(e.Job).Pointer()).Name()))
		g.Log().Debugf("获得订时器列表:%+v", runtime.FuncForPC(reflect.ValueOf(e.Job).Pointer()).Name())
		cronData := &CronData{
			Name: e.Name,
			//RegisterTime:e.Time.Format("2006-01-02 15:04:05"),
			RegisterTime: e.Time.Unix(),
			State:        e.Status(),
		}
		data = append(data, cronData)
	}
	return data, len(data)
}

// 验证订时字符串内容
func CheckCronTimeStr(TimeStr string) bool {
	TimeList := strings.Split(TimeStr, " ")
	InitValue := 0
	endValue := 60
	for Index, TimeStr1 := range TimeList {
		if Index == 0 || Index == 1 {
			InitValue = 0
			endValue = 59
		} else if Index == 2 {
			InitValue = 0
			endValue = 23
		} else if Index == 3 {
			InitValue = 1
			endValue = 31
		} else if Index == 4 {
			InitValue = 1
			endValue = 12
		} else if Index == 5 {
			InitValue = 1
			endValue = 7
		} else if Index == 6 {
			InitValue = 2020
			endValue = 2050
		}
		g.Log().Debug("CheckCronTimeStr:%d  %s", Index, TimeStr1)
		if TimeStr1 == "*" {
			continue
		} else if (Index == 3 || Index == 5) && TimeStr1 == "?" {
			continue
		} else if strings.Index(TimeStr1, "/") != -1 {
			ValueList := strings.Split(TimeStr1, "/")
			if len(ValueList) != 2 {
				g.Log().Errorf("订时字符串区间内容错误：格式*/1 %s", TimeStr)
				return false
			}
			if ValueList[0] != "*" {
				g.Log().Errorf("订时字符串区间内容错误：格式*/1 %s", TimeStr)
				return false
			} else {
				v2 := gconv.Int(ValueList[1])
				if !(InitValue <= v2 && v2 <= endValue) {
					g.Log().Errorf("订时字符串区间字段内容错误 %d-%d:%s", InitValue, endValue, TimeStr)
					return false
				}
			}

		} else if strings.Index(TimeStr1, "-") != -1 {
			ValueList := strings.Split(TimeStr1, "-")
			if len(ValueList) != 2 {
				g.Log().Errorf("订时字符串范围内容错误：格式1-3 %s", TimeStr)
				return false
			}
			v1 := gconv.Int(ValueList[0])
			v2 := gconv.Int(ValueList[1])
			if v1 >= v2 {
				g.Log().Errorf("订时字符串范围内容要从小到大:%s", TimeStr)
				return false
			}
			if v1 >= v2 {
				g.Log().Errorf("订时字符串范围内容要从小到大:%s", TimeStr)
				return false
			} else if !(InitValue <= v1 && v1 <= endValue && InitValue <= v2 && v2 <= endValue) {
				g.Log().Errorf("订时字符串字段范围内容错误 %d-%d:%s", InitValue, endValue, TimeStr)
				return false
			}
		} else if strings.Index(TimeStr1, ",") != -1 {
			ValueList := strings.Split(TimeStr1, ",")
			if len(ValueList) >= endValue {
				g.Log().Errorf("订时字符串分段内容错误：格式1-3 %s", TimeStr)
				return false
			}
			v0 := 0
			for _, VV := range ValueList {
				v1 := gconv.Int(VV)
				if v0 >= v1 {
					g.Log().Errorf("订时字符串分段内容要从小到大：格式1,2,3... %s", TimeStr)
					return false
				}
				v0 = v1
			}
		} else {
			v1 := gconv.Int(TimeStr1)
			if !(InitValue <= v1 && v1 <= endValue) {
				g.Log().Errorf("订时字符串字段内容错误 要在%d-%d之间:%s", InitValue, endValue, TimeStr)
				return false
			}
		}
	}
	return true
}

//修复ltv 留存数据 每日统计
// func repairData() {
// 	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
// 	for i := 1; i <= 5; i++ {
// 		g.Log().Infof("-----修复ltv 留存数据 每日统计: %d", i)
// 		repairTime := todayZeroTimestamp - 86400*i
// 		DoUpdateAllGameNodeDailyLTV(repairTime)
// 		DoUpdateAllGameNodeRemainTotal(repairTime)
// 		DoUpdateAllGameNodeDailyStatistics(repairTime)
// 	}
// }

func RepireRemainTotal() {
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*2)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*3)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*4)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*5)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*6)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*7)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*8)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*9)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*10)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*11)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*12)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*13)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*14)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*15)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*16)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*17)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*18)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*19)
	DoUpdateAllGameNodeRemainTotal(todayZeroTimestamp - 86400*20)
}

//每分钟执行一次
func minuteCron() {
	//timer1 := time.NewTicker(60 * time.Second)
	//for {
	//	select {
	//	case <-timer1.C:
	// 业务
	go CheckAllGameNode()
	go CheckPlatformServerClose()
	//go merge.HandleCheckCanMergePlatform()
	//	}
	//}
}

//定时检测自动开服
func cronAutoCreateServer() {
	cronSecond := g.Cfg().GetInt("game.checkOpenServerCronSecond")
	//utils.CheckError(err, "读取自动开服定时时间失败")
	//if err != nil {
	//	return
	//}
	if cronSecond <= 0 {
		return
	}
	g.Log().Infof("开服检测间隔时间:%d秒", cronSecond)
	gtimer.Add(time.Duration(cronSecond)*time.Second, func() {
		platformList := GetPlatformSimpleList()
		// 自动开服
		now := utils.GetTimestamp()
		for _, platform := range platformList {
			PlatformId := platform.Id
			if platform.CreateRoleLimit < 1 || platform.IsAutoOpenServer == 0 {
				continue
			}
			err := AutoCreateAndCreateRoleLimit(PlatformId, now)
			if err != nil {
				SendBackgroundMsgTemplateHandle(MSG_TEMPLATE_AUTO_CREATE_SERVER_FAIL, PlatformId, PlatformId+"定时自动开服失败", "自动开服详情："+utils.TimeIntFormDefault(now)+" err:"+fmt.Sprintf("%v", err))
			}
		}
	})
}

//10秒执行一次
func tenSecondCron() {
	//g.Log().Info("10秒执行一次")
	//timer1 := time.NewTicker(5 * time.Second)
	//for {
	//	select {
	//	case <-timer1.C:
	//		// 业务
	//		//go test()
	//	}
	//}
}

//每天0点执行
func dailyZeroClockCron() {
	// 业务
	g.Log().Info("0点执行定时器")
	go UpdateAllGameNodeDailyStatistics()
	go UpdateAllGameNodeRemainTotal()
	go UpdateAllGameNodeRemainActive()
	go UpdateAllGameNodeDailyLTV()
	go UpdateAllGameNodeChargeRemain()
}

//每天0点5分执行
func dailyZeroClock5MinuteCron() {
	// 业务
	g.Log().Info("0点5分执行定时器")
	hour := time.Now().Hour()
	go AutoOpenServerNow(hour)
	//}
}

//整点执行
func everyHourClockCron() {
	// 业务
	g.Log().Info("每小时执行定时器")

	//　定时ping 数据库， 防止被断开连接
	go PingDb(Db)
	go PingDb(DbCenter)
	go PingDb(DbCharge)
	hour := time.Now().Hour()
	if hour > 0 {
		go AutoOpenServerNow(hour)
	}
	go CheckNotAuditPlatform()
	go UpdateDingYueStatistics("wx", utils.GetTimestamp()-30) // 生成前一天的订阅数据

	if time.Now().Hour() == 0 {
		// 零点生成前一天23点数据
		yesTime := time.Now().AddDate(0, 0, -1)
		go InsertItemEventLog(yesTime, 23, 0)   // 物品，事件使用日志
		go InsertGameMonsterLog(yesTime, 23, 0) // 怪物事件日志

	} else {
		go InsertItemEventLog(time.Now(), time.Now().Hour()-1, 0)   // 物品，事件使用日志
		go InsertGameMonsterLog(time.Now(), time.Now().Hour()-1, 0) // 怪物事件日志
	}

}

//整10分钟执行
func tenMinuteClockCron() {
	next1 := gtime.Now()
	nextTimestamp := next1.Unix()
	//// 业务
	g.Log().Infof("整点10分钟定时执行:%v", next1.String())
	DoUpdateAllGameNodeTenMinuteStatistics(int(nextTimestamp))
	go CheckBackWebActive()
}

// 整点自动开服
func AutoOpenServerNow(hour int) {
	g.Log().Infof(">>>>>>\t整点自动开服时间:%d", hour)
	//isClockOpenServer := IsSettingOpenDefault(SETTING_DATA_IS_CHECK_AUTO_OPEN_SERVER, false)
	////isClockOpenServer := g.Cfg().GetBool("game.isClockOpenServer", false)
	//if isClockOpenServer {
	todayZeroTimestamp := utils.GetTodayZeroTimestamp()
	currTime := utils.GetTimestamp()
	platformList := GetPlatformSimpleList()
	for _, platform := range platformList {
		if platform.IntervalInitTime > currTime || len(platform.OpenServerTimeScope) == 0 {
			continue
		}
		if platform.IntervalInitTime > 0 && platform.IntervalDay > 0 {
			calcIntervalDay := gconv.Int(math.Floor(gconv.Float64(platform.IntervalInitTime-todayZeroTimestamp) / 86400))
			if calcIntervalDay%(platform.IntervalDay) != 0 {
				continue
			}
		}
		hList := strings.Split(platform.OpenServerTimeScope, ",")
		for _, hStr := range hList {
			h := gconv.Int(hStr)
			if hour == h {
				g.Log().Infof("当前整点自动开服平台:%s h:%d; hStr:%s", platform.Id, h, hStr)
				//OpenServerNow(platform.Id)
				err := OpenServerType(0, platform.Id, 0, currTime)
				utils.CheckError(err, "整点自动开服失败!!!!!!"+platform.Id)
			}
		}
		//}
	}
}
