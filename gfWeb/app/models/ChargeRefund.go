package models

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/url"
	"sort"

	"gfWeb/library/utils"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

type ChargeRefundResult struct {
	RefundBatchId  string
	RefundPayMoney int
}

func GetChargeRefundResult(PlatformId string, RefundType int, RefundReason string, OrderId string, ApplyRefundMoney int, PlatformOrderId string, PlatformUserId string) {
	URL := "https://etrade-api.baidu.com/cashier/applyOrderRefund"
	AppKey := "MMMzi5"
	//初始化参数
	param := url.Values{}
	if PlatformId == "baidu" {
		param.Set("orderId", PlatformOrderId)             // 百度平台订单id
		param.Set("userId", PlatformUserId)               // 百度用户id
		param.Set("refundType", gconv.String(RefundType)) //
		param.Set("refundReason", RefundReason)
		param.Set("tpOrderId", OrderId)
		param.Set("appKey", AppKey)
		param.Set("applyRefundMoney", gconv.String(ApplyRefundMoney*100))
		param.Set("bizRefundBatchId", OrderId)
	}
	dataParams := ""
	var keys []string
	for k := range param {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	//拼接
	for _, k := range keys {
		if dataParams != "" {
			dataParams += "&"
		}
		dataParams += k + "=" + param.Get(k)
	}
	fmt.Println(dataParams)
	RsaSign := getChargeRefundSign(dataParams)
	param.Set("rsaSign", RsaSign)
	_, err := utils.HttpGet(URL, param)
	if err != nil {
		g.Log().Error("申请退款失败:err:%v", err)
		return
	}
	return
}

// 获得充值
func getChargeRefundSign(str string) string {
	privateKeyBytes, err := ioutil.ReadFile("../../library/key/baidu/baidu_rsa_private_key.pem")
	if err != nil {
		return ""
	}
	block, _ := pem.Decode([]byte(privateKeyBytes))
	if block == nil {
		panic("私钥错误")
		return ""
	}
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes) //之前看java demo中使用的是pkcs8
	if err != nil {
		fmt.Print(err)
		panic("PrivateKey error")
		return ""
	}
	h := crypto.Hash.New(crypto.SHA1) //进行SHA1的散列
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	// 进行rsa加密签名
	signedData, err := rsa.SignPKCS1v15(rand.Reader, private, crypto.SHA1, hashed)
	data := base64.StdEncoding.EncodeToString(signedData)
	fmt.Println(data)
	return data
}
