package model

import (
	"c-server/cache"
	"strconv"
	"time"
)

// User 用户模型
type Player struct {
	Id                 uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Town               string `gorm:"size:255;column:NAME"`
	Preference         string `gorm:"size:255"`
	Portal             string `gorm:"size:255"`
	Nickname           string `gorm:"size:255"`
	Email              string `gorm:"size:255"`
	Phone              string `gorm:"size:255"`
	Username           string `gorm:"size:255"`
	Password           string `gorm:"size:255"`
	EthChargeAddress   string `gorm:"size:255;column:ethChargeAddress"`
	EthChargePrivKey   string `gorm:"size:255;column:ethChargePrivKey"`
	EthBindingAddress  string `gorm:"size:255;column:ethBindingAddress"`
	NeoChargeAddress   string `gorm:"size:255;column:neoChargeAddress"`
	NeoBindingAddress  string `gorm:"size:255;column:neoBindingAddress"`
	NeoChargePrivKey   string `gorm:"size:255;column:neoChargePrivKey"`
	NeoChargeWIF       string `gorm:"size:255;column:neoChargeWIF"`
	NeoChargePassphase string `gorm:"size:255;column:neoChargePassphase"`
	GenesisId          string `gorm:"size:255;column:genesisId"`
	IsBanned           bool
	IsAdmin            bool
	IsActivated        bool
	Code               string `gorm:"size:255"`
	Credit             float64
	CreditGave         string `gorm:"size:255"`
	NewbieAt           time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// GetUser 用ID获取用户
func GetPlayer(ID interface{}) (Player, error) {
	var player Player
	result := DB.First(&player, ID)
	return player, result.Error
}
func GetPlayerByEmail(Email interface{}) (Player, error) {
	var player Player
	result := DB.First(&player, Email)
	return player, result.Error
}
func (p Player) CheckPassword(pwd string) bool {
	return p.Password == pwd
}
func (p Player) GetNeoBalance() (balance string) {
	var neo string
	neo = "player:" + strconv.Itoa(int(p.Id)) + ":coin9"
	n, err := cache.RedisClient.Get(neo).Result()
	if err != nil {
		return "0"
	}
	return n
}
func (p Player) GetCoin(coinName string) (coin string) {
	KEY := `player:` + strconv.Itoa(int(p.Id)) +":"+ coinName
	n, err := cache.RedisClient.Get(KEY).Result()
	if err != nil {
		return "0"
	}
	return n
}
