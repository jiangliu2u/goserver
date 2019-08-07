package game

import (
	"c-server/cache"
	"strconv"
)

type Token struct {
	token string
}

func (token *Token) saveToken(userId uint) {
	var a string
	a = "openid:webcenter:" + strconv.Itoa(int(userId)) + ":token"
	cache.RedisClient.Set(a, token.token, 0)
}
