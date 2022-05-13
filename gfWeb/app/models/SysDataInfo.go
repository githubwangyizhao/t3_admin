package models

import (
	"fmt"
	"gfWeb/library/utils"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"os"
	"runtime"
	"strconv"
	"time"
)

// 系统机器数据信息
type SysRobotInfo struct {
	CpuNum          int              `json:"cpuNum"`
	CpuUsed         float64          `json:"cpuUsed"`
	CpuAvg5         float64          `json:"cpuAvg5"`
	CpuAvg15        float64          `json:"cpuAvg15"`
	MemTotal        uint64           `json:"memTotal"`
	MemUsed         uint64           `json:"memUsed"`
	MemFree         uint64           `json:"memFree"`
	MemUsage        float64          `json:"memUsage"`
	SysComputerName string           `json:"sysComputerName"`
	SysOsName       string           `json:"sysOsName"`
	SysComputerIp   string           `json:"sysComputerIp"`
	SysOsArch       string           `json:"sysOsArch"`
	GoTotal         uint64           `json:"goTotal"`
	GoUsed          uint64           `json:"goUsed"`
	GoFree          uint64           `json:"goFree"`
	GoUsage         float64          `json:"goUsage"`
	GoName          string           `json:"goName"`
	GoVersion       string           `json:"goVersion"`
	GoStartTime     string           `json:"goStartTime"`
	GoRunTime       string           `json:"goRunTime"`
	GoHome          string           `json:"goHome"`
	GoUserDir       string           `json:"goUserDir"`
	Disklist        []disk.UsageStat `json:"diskList"`
}

var StartTime = gtime.Timestamp()

func GetSysRobotInfo() *SysRobotInfo {
	cpuNum := runtime.NumCPU()      //核心数
	MB := gconv.Uint64(1024 * 1024) // 转成Mb

	var cpuUsed float64 = 0  //用户使用率
	var cpuAvg5 float64 = 0  //CPU负载5分钟平均
	var cpuAvg15 float64 = 0 //CPU负载15分钟平均

	cpuInfo, err := cpu.Percent(time.Duration(time.Second), false)
	if err == nil {
		cpuUsed, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", cpuInfo[0]), 64)
	}

	loadInfo, err := load.Avg()
	if err == nil {
		cpuAvg5, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", loadInfo.Load5), 64)
		cpuAvg15, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", loadInfo.Load15), 64)
	}

	var memTotal uint64 = 0  //总内存
	var memUsed uint64 = 0   //总内存  := 0 //已用内存
	var memFree uint64 = 0   //剩余内存
	var memUsage float64 = 0 //使用率

	v, err := mem.VirtualMemory()
	if err == nil {
		memTotal = v.Total / MB
		memUsed = v.Used / MB
		memFree = memTotal - memUsed
		memUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", v.UsedPercent), 64)
	}

	var goTotal uint64 = 0  //go分配的总内存数
	var goUsed uint64 = 0   //go使用的内存数
	var goFree uint64 = 0   //go剩余的内存数
	var goUsage float64 = 0 //使用率

	var gomem runtime.MemStats
	runtime.ReadMemStats(&gomem)
	goUsed = gomem.Sys / MB
	goUsage = gconv.Float64(fmt.Sprintf("%.2f", gconv.Float64(goUsed)/gconv.Float64(memTotal)*100))
	sysComputerIp := "" //服务器IP

	ip, err := utils.GetLocalIP()
	if err == nil {
		sysComputerIp = ip
	}

	sysComputerName := "" //服务器名称
	sysOsName := ""       //操作系统
	sysOsArch := ""       //系统架构

	sysInfo, err := host.Info()

	if err == nil {
		sysComputerName = sysInfo.Hostname
		sysOsName = sysInfo.OS
		sysOsArch = sysInfo.KernelArch
	}

	goName := "GoLang"                                   //语言环境
	goVersion := runtime.Version()                       //版本
	goStartTime := utils.TimeInt64FormDefault(StartTime) //启动时间

	goRunTime := utils.FormatTime64Second(gtime.Timestamp() - StartTime) //运行时长
	goHome := runtime.GOROOT()                                           //安装路径
	goUserDir := ""                                                      //项目路径

	curDir, err := os.Getwd()

	if err == nil {
		goUserDir = curDir
	}

	//服务器磁盘信息
	disklist := make([]disk.UsageStat, 0)
	diskInfo, err := disk.Partitions(true) //所有分区
	if err == nil {
		for _, p := range diskInfo {
			diskDetail, err := disk.Usage(p.Mountpoint)
			if err == nil {
				diskDetail.UsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", diskDetail.UsedPercent), 64)
				if diskDetail.Total == 0 {
					continue
				}
				diskDetail.Total = diskDetail.Total / MB
				diskDetail.Used = diskDetail.Used / MB
				diskDetail.Free = diskDetail.Free / MB
				disklist = append(disklist, *diskDetail)
			}
		}
	}
	data :=
		&SysRobotInfo{
			CpuNum:          cpuNum,
			CpuUsed:         cpuUsed,
			CpuAvg5:         cpuAvg5,
			CpuAvg15:        cpuAvg15,
			MemTotal:        memTotal,
			GoTotal:         goTotal,
			MemUsed:         memUsed,
			GoUsed:          goUsed,
			MemFree:         memFree,
			GoFree:          goFree,
			MemUsage:        memUsage,
			GoUsage:         goUsage,
			SysComputerName: sysComputerName,
			SysOsName:       sysOsName,
			SysComputerIp:   sysComputerIp,
			SysOsArch:       sysOsArch,
			GoName:          goName,
			GoVersion:       goVersion,
			GoStartTime:     goStartTime,
			GoRunTime:       goRunTime,
			GoHome:          goHome,
			GoUserDir:       goUserDir,
			Disklist:        disklist,
		}
	return data
}
