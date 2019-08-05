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
