package utils

import (
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/grand"
	"strings"
)

//将字符串加密成 md5
func String2md5(str string) string {
	md5Str, err := gmd5.EncryptString(str)
	if err != nil {
		g.Log().Errorf("将字符串加密成失败:", err)
		return ""
	}
	//data := []byte(str)
	//has := md5.Sum(data)
	//return fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5Str
}

//RandomString 在数字、大写字母、小写字母范围内生成num位的随机字符串
func RandomString(length int) string {
	RandStr := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return grand.Str(RandStr, length)
	//// 48 ~ 57 数字
	//// 65 ~ 90 A ~ Z
	//// 97 ~ 122 a ~ z
	//// 一共62个字符，在0~61进行随机，小于10时，在数字范围随机，
	//// 小于36在大写范围内随机，其他在小写范围随机
	//rand.Seed(time.Now().UnixNano())
	//result := make([]string, 0, length)
	//for i := 0; i < length; i++ {
	//	t := rand.Intn(62)
	//	if t < 10 {
	//		result = append(result, strconv.Itoa(rand.Intn(10)))
	//	} else if t < 36 {
	//		result = append(result, string(rand.Intn(26)+65))
	//	} else {
	//		result = append(result, string(rand.Intn(26)+97))
	//	}
	//}
	//return strings.Join(result, "")
}

// 字符串转驼峰格式
func StrToHump(s string) string {
	if len(s) == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(len(s))
	flag := false
	index := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' && index > 0 {
			flag = true
		} else if c == '_' {
			continue
		} else if flag == true && 'a' <= c && c <= 'z' {
			flag = false
			c -= 'a' - 'A'
			b.WriteByte(c)
		} else if index == 0 && 'A' <= c && c <= 'Z' {
			index++
			c += 'a' - 'A'
			b.WriteByte(c)
		} else {
			index++
			b.WriteByte(c)
		}
	}
	return b.String()
}

// 字符串驼峰转下划线小写格式
func StrHumpToUnderline(s string) string {
	return StrHumpToUnderlineDefault(s, "")
}

// 字符串驼峰转下划线小写格式并增加默认值
func StrHumpToUnderlineDefault(s string, defaultStr string) string {
	if len(s) == 0 {
		return defaultStr
	}
	var b strings.Builder
	b.Grow(len(s) * 2)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			b.WriteByte('_')
			c += 'a' - 'A'
		}
		b.WriteByte(c)
	}
	return b.String()
}

// 时间转成订时器字符串
func TimestampToCronStr(timestamp int64) string {
	//return TimeIntForm(timestamp, "05 04 15 02 01 ? 2006")
	return TimeIntForm(timestamp, "05 04 15 02 01 ?")
}
