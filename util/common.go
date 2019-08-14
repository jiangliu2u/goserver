package util

import (
	"math/rand"
	"reflect"
	"time"
)

// RandStringRunes 返回随机字符串
func RandStringRunes(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func HasAttr(data map[string]interface{}, key string) bool {
	if _, ok := data[key]; ok {
		return true
	} else {
		return false
	}
}

//判断注册的函数里是否有这个函数
func HasFunc(data map[string]reflect.Value, key string) bool {
	if _, ok := data[key]; ok {
		return true
	} else {
		return false
	}
}

// 按照 Int64Slice 从大到小排序
type Int64Slice [] int64

func (a Int64Slice) Len() int {    	 // 重写 Len() 方法
	return len(a)
}
func (a Int64Slice) Swap(i, j int){     // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a Int64Slice) Less(i, j int) bool {    // 重写 Less() 方法， 从大到小排序
	return a[j] < a[i]
}
