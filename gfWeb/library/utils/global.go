package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gfWeb/library/enums"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gogf/gf/encoding/gcharset"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gproc"
)

// 检查是否有错误
func CheckError(err error, msg ...string) {
	if err != nil {
		g.Log().Errorf("检查是否有错误:%s %v", msg, err)
	}
}

func CenterNodeTool(arg ...string) (string, error) {
	centerNode := g.Cfg().GetString("game.centerNode")
	out, err := NodeTool(centerNode, arg...)
	return out, err
}

func NodeTool(node string, arg ...string) (string, error) {
	cookie := g.Cfg().GetString("game.cookie")
	commandArgs := []string{
		"nodetool",
		"-name",
		node,
		"-setcookie",
		cookie,
		"rpc",
	}
	for _, v := range arg {
		if v == "" {
			commandArgs = append(commandArgs, "''")
		} else {
			commandArgs = append(commandArgs, v)
		}

	}
	out, err := CmdNodetool("escript", commandArgs)
	return out, err
}

// 使用nodetool脚本
func CmdNodetool(commandName string, params []string) (string, error) {
	return CmdByDirOrParam(GetToolDir(), commandName, params)
}

//封包
func Packet(methodNum int, message []byte) []byte {
	return append(append([]byte{0}, IntToBytes(methodNum)...), message...)
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// 获取gm 地址
func GetCenterURL() string {
	url := g.Cfg().GetString("game.gs_domain")
	return url
}

// 获取gm 地址
func GetChargeURL() string {
	url := g.Cfg().GetString("game.chargeUrl")
	return url
}

// 获取工具路径
func GetToolDir() string {
	ToolPath := g.Cfg().GetString("ansible.tool_path", "/data/tool/ansible/")
	return ToolPath
}

// 获取查看实时文件路径
func GetShowFileDir() string {
	Path := g.Cfg().GetString("logger.Path", "logs")
	PathAbs, _ := filepath.Abs(Path)
	return PathAbs + "/showFile/"
}

func HttpRequestGetObj(url string, data string) (interface{}, error) {
	sign := String2md5(data + enums.GmSalt)
	base64Data := base64.URLEncoding.EncodeToString([]byte(data))
	requestBody := "data=" + base64Data + "&sign=" + sign
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(requestBody))
	CheckError(err)
	//if err != nil {
	//	return "http请求失败", err
	//}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	g.Log().Debug("responseBody ", responseBody)
	CheckError(err)
	if err != nil {
		return "读取http返回内容失败", err
	}

	var jsonSlice2 = make(map[string]interface{})
	err = json.Unmarshal([]byte(fmt.Sprintf(`%s`, responseBody)), &jsonSlice2)

	return jsonSlice2["error_msg"], err

	/*
		g.Log().Debug("eee: ", jsonSlice2)
		g.Log().Debug("err: ", reflect.ValueOf(jsonSlice2))

		Reflect := reflect.ValueOf(jsonSlice2).Elem()
		//g.Log().Debug("ddd: ", Reflect.FieldByName("items"))
		g.Log().Debug("ddd: ", Reflect.FieldByName("pay_info"))

		returnString := jsonSlice2["error_msg"].(map[string]interface{})["items"].(string)
		g.Log().Debug("items: ", reflect.TypeOf(returnString))
		return returnString, err
	*/
}

func DecodeHttpRequest(data, sign string) map[string]string {
	dataBytes, _ := base64.URLEncoding.DecodeString(data)
	culSign := String2md5(string(dataBytes) + enums.GmSalt)
	if culSign != sign {
		return nil
	}

	var ret = make(map[string]string)
	json.Unmarshal(dataBytes, &ret)

	return ret
}

func EncodeHttpRequest(data string) (base64Data, sign string) {
	sign = String2md5(data + enums.GmSalt)
	base64Data = base64.URLEncoding.EncodeToString([]byte(data))
	return
}

func HttpRequest(url string, data string) (string, error) {
	var result struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	}
	sign := String2md5(data + enums.GmSalt)
	base64Data := base64.URLEncoding.EncodeToString([]byte(data))
	requestBody := "data=" + base64Data + "&sign=" + sign
	g.Log().Debugf("url:%+v", url)
	g.Log().Debugf("requestBody:%+v", requestBody)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(requestBody))
	CheckError(err)
	if err != nil {
		return "http请求失败", err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	CheckError(err)
	if err != nil {
		return "读取http返回内容失败", err
	}
	fmt.Println("responseBody str : ", string(responseBody))

	err = json.Unmarshal(responseBody, &result)
	g.Log().Info("result:%+v", result)
	CheckError(err)
	if err != nil {
		g.Log().Error("http返回不是json结构 %s", responseBody)
		return "http返回不是json结构", err
	}
	g.Log().Debugf("resultErrorMsg:%+v", result.ErrorMsg)
	if result.ErrorCode != 0 {
		return result.ErrorMsg, gerror.New(result.ErrorMsg)
	}
	return result.ErrorMsg, nil
}

// 获取ip归属地
func GetIpLocation(ip string) string {
	//url := "http://ip.taobao.com/service/getIpInfo.php?ip=" + ip
	url := "http://118.25.181.121：7060/ip?ip=" + ip
	var result struct {
		//Code int
		Ret  int
		Data struct {
			Country string
			Region  string
			City    string
			Isp     string
		}
	}
	resp, err := http.Get(url)
	CheckError(err)
	if err != nil {
		return "未知"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	CheckError(err)
	if err != nil {
		return "未知"
	}
	//g.Log().Info("result:%v", string(body))

	err = json.Unmarshal(body, &result)
	CheckError(err)
	if err != nil {
		return "未知"
	}
	if result.Ret == 0 {
		if result.Data.Country == "中国" {
			return result.Data.Region + "." + result.Data.City + " " + result.Data.Isp
		}
		return result.Data.Country + "." + result.Data.Region + "." + result.Data.City + " " + result.Data.Isp
	}
	return "未知"
}

//获取ip所属城市
func GetCityByIp(ip string) string {
	if ip == "" {
		return ""
	}
	if ip == "[::1]" || ip == "127.0.0.1" {
		return "内网IP"
	}
	url := "http://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip
	bytes := ghttp.GetBytes(url)
	src := string(bytes)
	srcCharset := "GBK"
	tmp, _ := gcharset.ToUTF8(srcCharset, src)
	json, err := gjson.DecodeToJson(tmp)
	if err != nil {
		return ""
	}
	if json.GetInt("code") == 0 {
		city := json.GetString("city")
		return city
	} else {
		return ""
	}
}

//// 获取ip归属地
//func GetIpLocation(ip string) string {
//	url := "http://int.dpool.sina.com.cn/iplookup/iplookup.php?format=json&ip=" + ip
//	var result struct {
//		Ret      int
//		Country  string
//		Province string
//		City     string
//	}
//	resp, err := http.Get(url)
//	CheckError(err)
//	if err != nil {
//		return "未知"
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	CheckError(err)
//	if err != nil {
//		return "未知"
//	}
//	//g.Log().Info("result:%v", string(body))
//
//	err = json.Unmarshal(body, &result)
//	CheckError(err)
//	if err != nil {
//		return "未知"
//	}
//	if result.Ret == 1 {
//		if result.Country == "中国" {
//			return result.Province + "." + result.City
//		}
//		return result.Country + "." +result.Province + "." + result.City
//	}
//	return "未知"
//}

//func FilePutContext(filename string, context string) error {
//	f, err := os.Create(filename) //创建文件
//	CheckError(err)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//	_, err = io.WriteString(f, context)
//	CheckError(err)
//	return err
//}

func ReportMsg(msgId string, phone string) {
	//请求地址
	juheURL := "http://v.juhe.cn/sms/send"

	//初始化参数
	param := url.Values{}

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("mobile", phone)                           //接收短信的手机号码
	param.Set("tpl_id", msgId)                           //短信模板ID，请参考个人中心短信模板设置
	param.Set("tpl_value", "")                           //变量名和变量值对。如果你的变量名或者变量值中带有#&amp;=中的任意一个特殊符号，请先分别进行urlencode编码后再传递，&lt;a href=&quot;http://www.juhe.cn/news/index/id/50&quot; target=&quot;_blank&quot;&gt;详细说明&gt;&lt;/a&gt;
	param.Set("key", "a02ffe0ce9754e8563edf690d4ebad4d") //应用APPKEY(应用详细页查询)
	param.Set("dtype", "")                               //返回数据的格式,xml或json，默认json

	//发送请求
	data, err := HttpGet(juheURL, param)
	if err != nil {
		g.Log().Error("请求失败,错误信息:\r\n%v", err)
	} else {
		var netReturn map[string]interface{}
		json.Unmarshal(data, &netReturn)
		g.Log().Infof("上报结果:%+v", netReturn)
		//if netReturn["error_code"].(float64)==0{
		//	fmt.Printf("接口返回result字段是:\r\n%v",netReturn["result"])
		//}
	}
}

//// get 网络请求
//func HttpGet(apiURL string, params url.Values) (rs []byte, err error) {
//	return HttpGetTimeout(apiURL , params , -1)
//}
//// get 网络请求加超时
//func HttpGetTimeout(apiURL string, params url.Values, timeoutS int) (rs []byte, err error) {
//	var Url *url.URL
//	Url, err = url.Parse(apiURL)
//	if err != nil {
//		g.Log().Error("解析url错误:", err)
//		return nil, err
//	}
//	//如果参数中有中文参数,这个方法会进行URLEncode
//	Url.RawQuery = params.Encode()
//	var resp *http.Response
//	if timeoutS > 0 {
//		httpClient := http.Client{Timeout: time.Duration(timeoutS) * time.Second}
//		resp, err = httpClient.Get(Url.String())
//	} else {
//		resp, err = http.Get(Url.String())
//	}
//	if err != nil {
//		g.Log().Errorf("HttpGet:", err)
//		return nil, err
//	}
//	StatusCode := resp.StatusCode
//	if StatusCode != 200 {
//		return nil, gerror.New(strconv.Itoa(StatusCode))
//	}
//	defer resp.Body.Close()
//	return ioutil.ReadAll(resp.Body)
//}

func HttpPost(url string, data string) error {
	var result struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	}
	sign := String2md5(data + enums.GmSalt)
	base64Data := base64.URLEncoding.EncodeToString([]byte(data))
	requestBody := "data=" + base64Data + "&sign=" + sign
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(requestBody))
	CheckError(err)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	CheckError(err)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &result)
	//g.Log().Info("result:%+v", result)
	CheckError(err)
	if err != nil {
		return err
	}
	if result.ErrorCode != 0 {
		return gerror.New(result.ErrorMsg)
	}
	return nil
}

func ExecShell(s string) (string, error) {
	// fmt.Println(s)
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return out.String(), err
	}
	//fmt.Printf("%s", out.String())
	return out.String(), nil
}
func GfExecShellRun(s string) error {
	return gproc.ShellRun(s)
}

//服务端ip
func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String(), nil
	}
	return
}

// 手动gc
func GC() {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	var mb uint64
	mb = 1024 * 1024
	logstr := fmt.Sprintf("手动gc 收集器内存:%vM  释放总内存:%vM  获得系统总内存:%vM  循环数:%v次", m.GCSys/mb, m.HeapReleased/mb, m.Sys/mb, m.NumGC)
	g.Log().Infof(logstr)
}
