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

type Advertising struct {
	// title, province, city, startTime, expireTime, schoolIds string, isAnonymous, isSchool, isTeacher, isStudent bool
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Province      string `json:"province"`
	City          string `json:"city""`
	StartTime     string `json:"from"`
	ExpireTime    string `json:"to"`
	SchoolIds     string `json:"school_ids"`
	ImageIds      string `json:"image_ids"`
	ImageNames    string `json:"image_names"`
	ImageLinks    string `json:"image_links"`
	PreviewUrls   string `json:"preview_urls"`
	IsLoginPage   bool   `json:"on_login_page"`
	IsSchoolPage  bool   `json:"on_school_page"`
	IsTeacherPage bool   `json:"on_teacher_page"`
	IsStudentPage bool   `json:"on_student_page"`
	DisplayPages  int    `json:"display_pages"`
	Pending       int    `json:"is_lock"`
	Click         int    `json:"click_count"`
	View          int    `json:"view_count"`

	CreateTime int `json:"create_time"`
	UpdateTime int `json:"update_time"`
	Delete     int `json:"delete"`
	EnableType int `json:"enable_type"`
}

func (ad *Advertising) CheckAd() error {
	// if ad.Title == "" {
	// 	return errors.New("Ad title is null")
	// }
	if ad.StartTime == "" {
		log.Println("set default begin time as current time")
		ad.StartTime = time.Now().Format(timeFmt)
	}
	if ad.ExpireTime == "" {
		log.Println("set default expire time as", defaultExpired)
		ad.ExpireTime = defaultExpired
	}
	return nil
}
