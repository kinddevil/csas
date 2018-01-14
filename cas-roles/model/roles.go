package model

import (
	"time"
)

type Role struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Permissions string    `json:"permission_ids"`
	UserIds     string    `json:"user_ids"`
	Desc        string    `json:"description"`
	Type        string    `json:"type"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
	IsDeleted   bool      `json:"is_deleted"`
}
