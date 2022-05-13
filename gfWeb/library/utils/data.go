package utils

import (
	"strings"

	"github.com/gogf/gf/util/gconv"
)

// 查看是否在切片中
func InIntSlice(t int, s []int) bool {
	for _, v := range s {
		if t == v {
			return true
		}
	}
	return false
}

// 切片去重
func RemoveReStrData(s []string) []string {
	if len(s) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(s))
	temp := map[string]struct{}{}
	for _, item := range s {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func RemoveReIntData(s []int) []int {
	if len(s) == 0 {
		return []int{}
	}

	result := make([]int, 0, len(s))
	temp := map[int]struct{}{}
	for _, item := range s {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func IfStr(a, b string) string {
	if strings.Trim(a, " ") != "" {
		return a
	}
	return b
}

func IfInt(a, b int) int {
	if a != 0 {
		return a
	}
	return b
}

func GconvToSliceMapStr(source []map[string]interface{}) []map[string]string {

	var newSliceMapStr []map[string]string
	for _, m := range source {
		newSliceMapStr = append(newSliceMapStr, gconv.MapStrStr(m))
	}

	return newSliceMapStr

}
