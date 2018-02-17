package model

import (
	"time"
)

type User struct {
	Username         string    `json:"username"`
	StudentsId       int64     `json:"students_id"`
	TeacherId        int64     `json:"teacher_id"`
	AdminId          int64     `json:"admin_id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	Activated        bool      `json:"activated"`
	Activationkey    string    `json:"activationkey"`
	ResetPasswordKey string    `json:"resetpasswordkey"`
	SchoolId         int64     `json:"school_id"`
	SchoolName       string    `json:"school_name"`
	Type             string    `json:"type"`
	RoleIds          string    `json:"role_ids"`
	Phone            string    `json:"phone"`
	LastLoad         time.Time `json:"last_load"`
	ModifyTime       time.Time `json:"modify_time"`
	IsLock           bool      `json:"is_lock"`
	LoadCounter      int       `json:"load_counter"`
	CreateTime       time.Time `json:"create_time"`
	UpdateTime       time.Time `json:"update_time"`
	PassVersion      string    `json::"pass_version"`
}
