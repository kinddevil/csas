package model

import (
	// "errors"
	"log"
	"time"
)

const (
	timeFmt        string = "2006-01-02 15:04:05"
	defaultExpired string = "2100-12-31 23:59:59"
)

type School struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Contact      string `json:"contact"`
	Phone        string `json:"phone"`
	Province     string `json:"province"`
	City         string `json:"city"`
	County       string `json:"county"`
	ProvinceCode string `json:"province_code"`
	CityCode     string `json:"city_code"`
	CountyCode   string `json:"conty_code"`
	Address      string `json:"address"`
	Fax          string `json:"fax"`
	Email        string `json:"email"`
	Web          string `json:"web"`
	Post         string `json:"post"`
	StartTime    string `json:"from"`
	ExpireTime   string `json:"to"`
	IsPayment    bool   `json:"is_payment"`
	TeacherNo    int    `json:"teacher"`
	StudentNo    int    `json:"student"`
	ContractId   string `json:"contract_id"`
	Contract     string `json:"contract"`
	IsLock       bool   `json:"is_lock"`
	Description  string `json:"description"`
	Status       int    `json:"status"`

	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	// Delete int `json:"is_delete"`
}

func (obj *School) CheckSchool() error {
	// if obj.Name == "" {
	// 	return errors.New("School name is null")
	// }
	if obj.StartTime == "" {
		log.Println("set default begin time as current time")
		obj.StartTime = time.Now().Format(timeFmt)
	}
	if obj.ExpireTime == "" {
		log.Println("set default expire time as", defaultExpired)
		obj.ExpireTime = defaultExpired
	}
	return nil
}
