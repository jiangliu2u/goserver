package model

import "time"

type Res struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Resource  string
	Count     float64
	Player    string
	Reason    string
}

func GetAllRes(Resource interface{}) ([]Res, error) {
	var reses []Res
	result := DB.Where("resource = ?", Resource).Find(&reses)
	// result := DB.Find(&reses, Resource)
	return reses, result.Error
}
