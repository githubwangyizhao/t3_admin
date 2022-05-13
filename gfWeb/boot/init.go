package boot

import (
	"encoding/base64"
	"encoding/json"
	"gfWeb/app/controllers"
	"gfWeb/app/models"
	_ "gfWeb/app/models"
	"gfWeb/library/utils"
	_ "gfWeb/memdb"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gogf/gf/container/gqueue"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/swagger"
)

func init() {
	//s.SetIndexFolder(true)  // 是否允许列出Server主目录的文件列表（默认为false）
	//s.EnableHTTPS("/home/john/temp/server.crt", "/home/john/temp/server.key")
	//s.EnableHTTPS("", "")
	//s.EnableAdmin()

	// session改为内存方式(默认为内存+文件)
	//s.SetConfigWithMap(g.Map{
	//	"SessionStorage": gsession.NewStorageMemory(),
	//})
	s := g.Server()
	s.Plugin(&swagger.Swagger{})

	//初始化日志
	utils.InitLog()
	//初始化缓存
	utils.InitCache()
	models.InitCron()
	// 初始化异步任务协程
	utils.InitAsyncPool()

	controllers.PUSHQUEUE = gqueue.New()
	go ReadyFetchPushAreaData()

}

// 每1秒获取一次
func ReadyFetchPushAreaData() {
	g.Log().Warning("已开始准备接收推送数据....")
	for {
		select {
		case areaItem := <-controllers.PUSHQUEUE.C:
			if area, ok := areaItem.([]string); ok {
				baseInfo := area[len(area)-4:]
				registrationIds := area[:len(area)-4]

				// g.Log().Warning("baseInfo", baseInfo)
				// g.Log().Warning("len registrationIds ", len(registrationIds))

				title := baseInfo[0]
				desc := baseInfo[1]
				functionConfigId := baseInfo[2]
				sid := baseInfo[3]
				g.Log().Warningf("正在进行极光推送title : %s, desc : %s, functionConfigId : %s, registrationIds %v", title, desc, functionConfigId, registrationIds)
				pushOk := doJgPush(title, desc, functionConfigId, registrationIds)
				if !pushOk {
					g.Log().Error(sid + "推送出错")
				}

				//停1秒来保证一分钟不会超过60万
				time.Sleep(time.Duration(1 * time.Second))
				// g.Log().Warning(sid + "推送完成后一秒过后。。。。" + gtime.Datetime())
			}

		}
	}
}

func doJgPush(title, desc, functionConfigId string, registrationIds []string) bool {
	// 推送开始
	// var registrationIds []string
	// registrationIds = append(registrationIds, "\""+"120c83f760fb51a5110"+"\"")
	// registrationIds = append(registrationIds, "120c83f760fb51a5110") //todo 在数据库player_global中查找到所有的
	var (
		url         = g.Cfg().GetString("jiguang.url")
		method      = "POST"
		contentType = functionConfigId
	)
	title = utils.IfStr(title, g.Cfg().GetString("jiguang.gamename")) //headline中取值

	var payLoadStrMap = map[string]interface{}{
		"platform": "all",
		"audience": map[string][]string{
			"registration_id": registrationIds,
		},
		"notification": map[string]map[string]interface{}{
			"android": map[string]interface{}{
				"title": title,
				"alert": desc,
				"extras": map[string]string{
					"function_config": gconv.String(contentType),
				},
			},
			"ios": map[string]interface{}{
				"alert": desc,
				"extras": map[string]string{
					"function_config": gconv.String(contentType),
				},
			},
		},
		"message": map[string]string{
			"title":        title,
			"msg_content":  desc,
			"content_type": gconv.String(contentType),
		},
	}

	payLoadbyte, _ := json.Marshal(payLoadStrMap)
	g.Log().Warningf("method : %s, url : %s payLoadJson : %s", method, url, string(payLoadbyte))
	payload := strings.NewReader(string(payLoadbyte))

	// fmt.Println(payLoadStr)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		utils.CheckError(err)
		g.Log().Error("err 1 exit push")
		return false
	}

	authValue := base64.StdEncoding.EncodeToString([]byte(g.Cfg().GetString("jiguang.appKey") + ":" + g.Cfg().GetString("jiguang.masterSecret")))
	req.Header.Add("Authorization", "Basic "+authValue)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		utils.CheckError(err)
		g.Log().Error("err 2 exit push")
		return false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		utils.CheckError(err)
		g.Log().Error("err 3 exit push")
		return false
	}

	type RetError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	var responseBody = struct {
		Error RetError `json:"error"`
	}{}

	g.Log().Warning("push ret body : " + string(body))
	json.Unmarshal(body, &responseBody)

	return responseBody.Error.Code == 0
}
