package utils

import (
	"archive/zip"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// 写入文件
func FilePutContext(filename string, context string) error {
	f, err := os.Create(filename) //创建文件
	CheckError(err)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.WriteString(f, context)
	CheckError(err)
	return err
}

// 文件夹是否存在
func EnsureDir(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0777)
		return err
	} else {
		return err
	}
}

// 文件是否存在
func IsFileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// 获得当前目录
func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		path = strings.Replace(path, "\\", "/", -1)
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		return "", gerror.New(`Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

// 压缩文件成zip
func CompressorZip(pathFileName string) (string, error) {
	bytes, err := ioutil.ReadFile(pathFileName)
	if err != nil {
		g.Log().Errorf("读文件失败:%s err:%+v", pathFileName, err)
		return "", err
	}
	//fileSuffix := path.Ext(pathFileName)               //获取文件后缀
	fileName := path.Base(pathFileName) //获取文件名带后缀
	//pathFile := strings.TrimSuffix(pathFileName, fileSuffix) // 获取路径和文件名
	zipPathFileName := pathFileName + ".zip"
	zipFileName := path.Base(zipPathFileName)
	buf, _ := os.Create(zipPathFileName)
	w := zip.NewWriter(buf)
	defer func() {
		w.Close()
		buf.Close()
	}()
	f, err := w.Create(fileName)
	if err != nil {
		g.Log().Errorf("创建zip文件包中的文件失败：%s error:%+v", fileName, err)
		return zipFileName, err
	}
	_, err = f.Write(bytes)
	if err != nil {
		g.Log().Errorf("写入zip文件包中的文件失败：%s error:%+v", fileName, err)
		return zipFileName, err
	}
	g.Log().Infof("压缩文件完成：%+v", zipPathFileName)
	return zipFileName, err
}

// 获得当前目录下的目录名称(不包含子目录)
func GetCurrDirOrAllDir(dir string) []string {
	fileDirs := make([]string, 0)
	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir(dir)
	if err != nil {
		return fileDirs
	}
	fmt.Println(len(fileInfoList))
	for i := range fileInfoList {
		f := fileInfoList[i]
		if f.IsDir() {
			fileDirs = append(fileDirs, f.Name())
		}
	}
	return fileDirs
}

// 获得当前目录下的文件名称(不包含子目录)
func GetCurrDirOrAllFile(dir string) []string {
	fileDirs := make([]string, 0)
	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir(dir)
	if err != nil {
		return fileDirs
	}
	for i := range fileInfoList {
		f := fileInfoList[i]
		if f.IsDir() {
			continue
		}
		fileDirs = append(fileDirs, f.Name())
	}
	return fileDirs
}
