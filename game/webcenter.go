package game

import (
	"c-server/model"
	uuid "github.com/satori/go.uuid"
	"reflect"
)

//用户中心

type Webcenter struct {
	Name string
}
type WebcenterActions interface {
	Login(interface{}) []byte
	Register(interface{}) interface{}
}

func (wb Webcenter) Login(req ClientMessage) {
	res := ResponseData{}

	res.Data = make(map[string]interface{})
	loginfo, ok := req.Data["msg"].(map[string]interface{})
	if !ok {
		req.error("草拟吗啊,数据错误")
		return
	}
	var p model.Player

	if loginfo["type"] == "cm" {
		loginData := loginfo["data"].(map[string]interface{})
		accountType := loginData["accountType"]
		password := loginData["password"]
		if accountType == "email" {
			if err := model.DB.Where("email = ?", loginData["email"]).First(&p).Error; err != nil {
				req.error("草拟吗啊,账户或密码错误")
				return
			}
			if !p.CheckPassword(password.(string)) {
				req.error("草拟吗啊,账户或密码错误")
				return
			}
		} else {
			if err := model.DB.Where("phone = ?", loginData["phone"]).First(&p).Error; err != nil {
				req.error("草拟吗啊,账户或密码错误")
				return
			}
			if !p.CheckPassword(password.(string)) {
				req.error("草拟吗啊,账户或密码错误")
				return
			}
		}
	} else {
		email := loginfo["email"]
		password := loginfo["password"]
		if err := model.DB.Where("email = ?", email).First(&p).Error; err != nil {
			req.error("草拟吗啊,账户或密码错误")
			return
		}
		if !p.CheckPassword(password.(string)) {
			req.error("草拟吗啊,账户或密码错误")
			return
		}
	}
	b := p.GetNeoBalance()
	res.Put("uid", p.Id)
	res.Put("balance", b)
	jwt := uuid.NewV4()
	token := Token{token: jwt.String()}
	token.saveToken(p.Id)
	cake := p.GetCoin("coin1")
	res.Put("token", jwt.String())
	res.Put("ethChargeAddress", p.NeoChargeAddress)
	res.Put("cake", cake)

	req.response(res)
}
func (wb Webcenter) Register(b interface{}) interface{} {
	return nil
}

func (wb *Webcenter) RegisterController() map[string]reflect.Value {
	cont := make(map[string]reflect.Value)
	v := reflect.ValueOf(wb).Elem()
	t := reflect.TypeOf(wb).Elem()
	for i := 0; i < v.NumMethod(); i++ {
		action := v.Method(i)
		actionName := t.Method(i).Name
		cont[actionName] = action
	}
	return cont
}
