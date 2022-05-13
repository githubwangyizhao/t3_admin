package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"io"
	"os"
	"os/exec"

	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/os/gtime"
)

//func Cmd(commandName string, params []string) (string, error) {
//	cmd := exec.Command(commandName, params...)
//	fmt.Println("Cmd", cmd.Args)
//	var out bytes.Buffer
//	cmd.Stdout = &out
//	cmd.Stderr = os.Stderr
//	err := cmd.Start()
//	//if err != nil {
//	//	panic(err)
//	//}
//	if err != nil {
//		return "", err
//	}
//	err = cmd.Wait()
//	return out.String(), err
//}
//
//func CmdAndChangeDir(dir string, commandName string, params []string) (string, error) {
//	cmd := exec.Command(commandName, params...)
//	fmt.Println("CmdAndChangeDir", dir, cmd.Args)
//	var out bytes.Buffer
//	cmd.Stdout = &out
//	cmd.Stderr = os.Stderr
//	cmd.Dir = dir
//	err := cmd.Start()
//	if err != nil {
//		return "", err
//	}
//	err = cmd.Wait()
//	return out.String(), err
//}

// GF cmd方法带脚本位置
func GfCmd(params []string) (string, error) {
	return gproc.ShellExec("", params)
	//return gproc.ShellExecDir(dir, cmd, params)
}

// GF cmd方法带脚本位置
func GfCmdByDir(dir string, cmd string) (string, error) {
	dirCmd := "cd " + dir + " && " + cmd
	return gproc.ShellExec(dirCmd)
	//return gproc.ShellExecDir(dir, cmd, params)
}

// cmd sh 实时写入文件
func CmdShShowFileByDirOrParam(fileName, dir string, params []string) error {
	return CmdAndChangeDirToFile(fileName, dir, "sh", params)
}

// cmd 目录处理并实时写入文件
func CmdAndChangeDirToFile(fileName, dir, commandName string, params []string) error {
	initTime := GetTimestamp()
	//Path += "log/cmdLog/"
	Path := GetShowFileDir()
	EnsureDir(Path)
	f, err := os.Create(Path + fileName) //创建文件
	defer f.Close()
	cmd := exec.Command(commandName, params...)
	fmt.Println("CmdAndChangeDirToFile", dir, cmd.Args)
	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("cmd.StdoutPipe: ", err)
		return err
	}
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err = cmd.Start()
	if err != nil {
		return err
	}
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		_, err = f.WriteString(line) //写入文件(字节数组)
		f.Sync()
	}
	_, err = f.WriteString("==== 处理完毕 用时: " + FormatTimeSecond(GetTimestamp()-initTime) + " 当前时间:" + gtime.Datetime() + "\n") //写入文件(字节数组)
	f.Sync()
	err = cmd.Wait()
	return err
}

// cmd执行脚本
func Cmd(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	fmt.Println("Cmd", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}

// cmd sh切换目录执行脚本
func CmdShByDirOrParam(dir string, params []string) (string, error) {
	return CmdByDirOrParam(dir, "sh", params)
}

// cmd切换目录执行脚本
func CmdByDirOrParam(dir string, commandName string, params []string) (string, error) {
	g.Log().Infof("CmdByDirOrParam:%s %s", dir, dir)
	g.Log().Info(params)
	cmd := exec.Command(commandName, params...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}
