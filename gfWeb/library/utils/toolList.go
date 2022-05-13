package utils

import "github.com/gogf/gf/container/garray"

// 差集
func ListMinus(aL, bL []int) []int {
	var r []int
	if len(aL) == 0 {
		return r
	}
	flag := false
	for _, a := range aL {
		for ib, b := range bL {
			if b == a {
				bL = append(bL[:ib], bL[ib+1:]...)
				flag = true
				break
			}
		}
		if flag == false {
			r = append(r, a)
		}
		flag = false
	}
	return r
}

// 是否存在列表中
func IsHaveArray(v string, array []string) bool {
	for _, e := range array {
		if e == v {
			return true
		}
	}
	return false
}

// 是否存在列表中
func IsHaveIntArray(v int, array []int) bool {
	for _, e := range array {
		if e == v {
			return true
		}
	}
	return false
}

// 列表去重 v要大于0
func ArrayUnList(v int, l []int) []int {
	if v <= 0 {
		return l
	}
	list := garray.NewIntArrayFrom(l)
	list.Append(v)
	return list.Unique().Slice()
}

// mapList增加值
func MapListAdd(mapList map[string][]string, key string, value string) map[string][]string {
	tmpL, ok := mapList[key]
	if ok {
		mapList[key] = append(tmpL, value)
	} else {
		mapList[key] = append([]string{}, value)
	}
	return mapList
}
