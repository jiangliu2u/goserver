package service

import (
	"c-server/model"
	"c-server/serializer"
)

type ResourceService struct {
	Resource string `form:"resource" json:"resource" binding:"required"`
}
type Result struct {
	EthBuy  float64
	EthSell float64
	Count   int
}

func (service *ResourceService) Query() (*serializer.Response, error) {
	res := model.Res{
		Resource: service.Resource,
	}
	all, err := model.GetAllRes(res.Resource)
	var r Result
	r.Count = len(all)
	for _, value := range all {
		if value.Count > 0 {
			r.EthSell = +value.Count
		} else {
			r.EthBuy = +value.Count
		}
	}

	return &serializer.Response{
		Status: 40001,
		Msg:    "两次输入的密码不相同",
		Data:   r,
	}, err
}
