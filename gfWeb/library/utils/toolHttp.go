package utils

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

// get 网络请求
func HttpGet(apiURL string, params url.Values) (rs []byte, err error) {
	return HttpGetTimeout(apiURL, params, -1)
}

// get 网络请求加超时
func HttpGetTimeout(apiURL string, params url.Values, timeoutS int) (rs []byte, err error) {
	var Url *url.URL
	Url, err = url.Parse(apiURL)
	if err != nil {
		g.Log().Errorf("解析url错误::%+v", err)
		return nil, err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	var resp *http.Response
	if timeoutS > 0 {
		httpClient := http.Client{Timeout: time.Duration(timeoutS) * time.Second}
		resp, err = httpClient.Get(Url.String())
	} else {
		resp, err = http.Get(Url.String())
	}
	if err != nil {
		g.Log().Errorf("HttpGet:%+v", err)
		return nil, err
	}
	StatusCode := resp.StatusCode
	if StatusCode != 200 {
		return nil, gerror.New(strconv.Itoa(StatusCode))
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func HttpGetJsonSliceMap(url string) []map[string]interface{} {
	client := &http.Client{}
	hReq, err := http.NewRequest(http.MethodGet, url, nil)
	CheckError(err)

	res, err := client.Do(hReq)
	CheckError(err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	CheckError(err)

	expDataArr := gconv.SliceMap(body)

	return expDataArr
}
